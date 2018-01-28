package nmon

import (
	"context"
	"log"

	"github.com/influxdata/telegraf"
)

// processer parse data dan send metrics to sender
type processer struct {
	id       int                // id identity a processer to make diffrent with others processers
	dataChan <-chan []byte      // dataChan used by receiver to send data
	cancel   context.CancelFunc // cancel stop all the processer's jobs before it's stop

	acc telegraf.Accumulator
}

// start a processer with cancel
func (p *processer) Start() error {
	var ctx context.Context
	ctx, p.cancel = context.WithCancel(context.Background())
	go func(ctx context.Context) {
		log.Printf("starting porcesser %d\n", p.id)
		for {
			select {
			case data := <-p.dataChan:
				err := p.parseNmonData(data)
				if err != nil {
					log.Printf("parse nmon data failed. %s\n", err)
				}
				break
			case <-ctx.Done():
				return
			}
		}
	}(ctx)
	return nil
}

//
func (p *processer) Stop() error {
	log.Printf("stopping processer %d\n", p.id)
	p.cancel()
	return nil
}

func (p *processer) parseNmonData(data []byte) error {
	parser := new(nmonParser)
	err := parser.Parse(data)
	if err != nil {
		return err
	}
	parser.SetVersion("1.0")

	for _, v := range parser.Metrics() {
		p.acc.AddGauge(v.Measurement(), v.Fields(), v.Tags(), v.Time())
	}
	return nil
}

//
func newProcesser(id int, ch <-chan []byte, acc telegraf.Accumulator) *processer {
	return &processer{
		id:       id,
		dataChan: ch,
		acc:      acc,
	}
}
