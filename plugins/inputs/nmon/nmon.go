package nmon

import (
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
)

func init() {
	inputs.Add("nmon", newNmonServer)
}

//
var simpleConfig = `
	# mode: http/socket/proxy, default value: http
	mode = "http"
	
	# ip address and port will be listened, default value: 0.0.0.0:12345
	listen = "0.0.0.0:12345"
	
	# data client report is increment part or full, values can be increment
	# or full, default value: increment
	report_mode = increment
	
	# format , default value: nmon
	data_format = "nmon"
	
	# http mode, options you could set.
	# max body size, unit: byte, default value: 100*1024*1024
	max_body_size = 100000
`

// configuration for nmon server
type config struct {
	Mode        string `toml:"mode"`
	Listen      string `toml:"listen"`
	ReportMode  string `toml:"report_mode"`
	DataFormat  string `toml:"data_format"`
	MaxBodySize int64  `toml:"max_body_size"`
}

//
type NmonServer struct {
	Mode        string `toml:"mode"`
	Listen      string `toml:"listen"`
	ReportMode  string `toml:"report_mode"`
	DataFormat  string `toml:"data_format"`
	MaxBodySize int64  `toml:"max_body_size"`

	srv server
}

func newNmonServer() telegraf.Input {
	return &NmonServer{
		Mode:       httpMode,
		Listen:     "0.0.0.0:12345",
		DataFormat: "nmon",
	}
}

// implement ServiceInput interface
func (p *NmonServer) SampleConfig() string {
	return simpleConfig
}

//
func (p *NmonServer) Description() string {
	return "receive nmon message from other agent send to."
}

//
func (p *NmonServer) Gather(_ telegraf.Accumulator) error {
	return nil
}

//
func (p *NmonServer) SetParser(_ parsers.Parser) {
}

//
func (p *NmonServer) Start(acc telegraf.Accumulator) error {
	cfg := config{}
	cfg.Listen = p.Listen
	cfg.Mode = p.Mode
	cfg.ReportMode = p.ReportMode
	cfg.DataFormat = p.DataFormat
	cfg.MaxBodySize = p.MaxBodySize
	p.srv = newServer(cfg)
	p.srv.SetAccumulator(acc)
	return p.srv.Start()
}

//
func (p *NmonServer) Stop() {
	p.srv.Stop()
}
