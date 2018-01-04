package nmon

import (
	"github.com/influxdata/telegraf"
)

//
type socketServer struct {
	cfgs config
	acc  telegraf.Accumulator
}

//
func (p *socketServer) Start(cfg config) error {
	return nil
}

//
func (p *socketServer) Stop() error {
	return nil
}

//
func (p *socketServer) SetAccumulator(acc telegraf.Accumulator) {
	p.acc = acc
}
