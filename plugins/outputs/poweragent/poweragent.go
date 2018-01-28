package poweragent

import (
	ejson "encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/internal"
	"github.com/influxdata/telegraf/plugins/outputs"
	"github.com/influxdata/telegraf/plugins/serializers"
)

type SocketWriter struct {
	Address         string
	KeepAlivePeriod *internal.Duration

	serializers.Serializer

	net.Conn
}

func (sw *SocketWriter) Description() string {
	return "poweragent socket writer capable of handling tcp4 socket type."
}

func (sw *SocketWriter) SampleConfig() string {
	return `
  ## URL to connect to
  # address = "tcp://127.0.0.1:8094"
  # address = "tcp://example.com:http"
  # address = "tcp4://127.0.0.1:8094"

  ## Period between keep alive probes.
  ## Only applies to TCP sockets.
  ## 0 disables keep alive probes.
  ## Defaults to the OS configuration.
  # keep_alive_period = "5m"

`
}

func (sw *SocketWriter) SetSerializer(_ serializers.Serializer) {
	sw.Serializer = newJsonSerializer()
}

func (sw *SocketWriter) Connect() error {
	spl := strings.SplitN(sw.Address, "://", 2)
	if len(spl) != 2 {
		return fmt.Errorf("invalid address: %s", sw.Address)
	}

	c, err := net.Dial(spl[0], spl[1])
	if err != nil {
		return err
	}

	if err := sw.setKeepAliveByHeartbeat(c); err != nil {
		log.Printf("unable to configure keep alive (%s): %s", sw.Address, err)
	}

	sw.Conn = c
	return nil
}

func (sw *SocketWriter) setKeepAlive(c net.Conn) error {
	if sw.KeepAlivePeriod == nil {
		return nil
	}
	tcpc, ok := c.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("cannot set keep alive on a %s socket", strings.SplitN(sw.Address, "://", 2)[0])
	}
	if sw.KeepAlivePeriod.Duration == 0 {
		return tcpc.SetKeepAlive(false)
	}
	if err := tcpc.SetKeepAlive(true); err != nil {
		return err
	}
	return tcpc.SetKeepAlivePeriod(sw.KeepAlivePeriod.Duration)
}

func (sw *SocketWriter) setKeepAliveByHeartbeat(c net.Conn) error {
	if sw.KeepAlivePeriod == nil {
		return nil
	}

	ticker := time.NewTicker(sw.KeepAlivePeriod.Duration)

	go func(c net.Conn) {
		var (
			err  error
			lock sync.RWMutex
		)
		for {
			lock.Lock()
			_, err = c.Write([]byte("test\n"))
			if err != nil {
				log.Printf("send heartbeat failed. %s\n", err)
			} else {
				log.Printf("send heartbeat successful\n")
			}
			lock.Unlock()
			<-ticker.C
		}
	}(c)

	return nil

}

// Write writes the given metrics to the destination.
// If an error is encountered, it is up to the caller to retry the same write again later.
// Not parallel safe.
func (sw *SocketWriter) Write(metrics []telegraf.Metric) error {
	if sw.Conn == nil {
		// previous write failed with permanent error and socket was closed.
		if err := sw.Connect(); err != nil {
			return err
		}
	}

	for _, m := range metrics {
		bs, err := sw.Serialize(m)
		if err != nil {
			//TODO log & keep going with remaining metrics
			return err
		}
		if _, err := sw.Conn.Write(bs); err != nil {
			//TODO log & keep going with remaining strings
			if err, ok := err.(net.Error); !ok || !err.Temporary() {
				// permanent error. close the connection
				log.Printf("write error: %s\n", err)
				sw.Close()
				sw.Conn = nil
			}
			return err
		}
	}

	return nil
}

// Close closes the connection. Noop if already closed.
func (sw *SocketWriter) Close() error {
	if sw.Conn == nil {
		return nil
	}
	err := sw.Conn.Close()
	sw.Conn = nil
	return err
}

func newSocketWriter() *SocketWriter {
	s, _ := serializers.NewInfluxSerializer()
	return &SocketWriter{
		Serializer: s,
	}
}

type JsonSerializer struct {
	TimestampUnits time.Duration
}

func newJsonSerializer() *JsonSerializer {
	return new(JsonSerializer)
}

func (p *JsonSerializer) Serialize(metric telegraf.Metric) ([]byte, error) {
	m := make(map[string]interface{})
	units_nanoseconds := p.TimestampUnits.Nanoseconds()
	// if the units passed in were less than or equal to zero,
	// then serialize the timestamp in seconds (the default)
	if units_nanoseconds <= 0 {
		units_nanoseconds = 1000000000
	}
	tags := metric.Tags()
	fields := metric.Fields()
	m["t"] = metric.UnixNano() / units_nanoseconds
	for k, v := range tags {
		if k == "a_id" {
			continue
		}
		m[k] = v
	}

	var s = make([]map[string]interface{}, 0)
	for k, v := range fields {
		ss := make(map[string]interface{})
		ss["a_id"] = tags["a_id"]
		ss["tag"] = k
		ss["value"] = v
		s = append(s, ss)
	}
	m["c"] = s

	serialized, err := ejson.Marshal(m)
	if err != nil {
		return []byte{}, err
	}
	serialized = append(serialized, '\n')

	return serialized, nil
}

func init() {
	outputs.Add("poweragent", func() telegraf.Output { return newSocketWriter() })
}
