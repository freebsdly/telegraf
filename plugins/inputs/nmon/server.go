package nmon

import (
	"strings"

	"github.com/influxdata/telegraf"
)

// server mode
const (
	httpMode   string = "restful"
	socketMode string = "socket"
	proxyMode  string = "proxy"
)

// server is interface , you can use this
// opreate socket/http server
type server interface {
	Start() error
	Stop() error
	SetAccumulator(acc telegraf.Accumulator)
}

//
func newServer(cfg config) server {
	var srv server
	switch strings.ToLower(cfg.Mode) {
	case socketMode:
	case proxyMode:
	default:
		srv = newRestfulApiServer(cfg)
	}
	return srv
}
