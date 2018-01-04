package nmon

import (
	"github.com/influxdata/telegraf"
)

//
type proxyServer struct {
	cfgs config
	acc  telegraf.Accumulator
}

//
func (p *proxyServer) Start(cfg config) error {
	return nil
}

//
func (p *proxyServer) Stop() error {
	return nil
}

//
func (p *proxyServer) SetAccumulator(acc telegraf.Accumulator) {
	p.acc = acc
}
