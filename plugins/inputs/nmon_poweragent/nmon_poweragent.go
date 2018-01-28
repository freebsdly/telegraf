// telegraf plugins need you implement a struct which include configuration options
// the struct must implement methods:
// SimpleConfig() string                    SimpleConfig method return a config
//                                          simple, used when run telegraf with
//                                          --config option
//
// Description() string                     Description method return a message
//                                          descript this plugin
//
// Gather(_ telegraf.Accumulator) error     if telegraf as a agent collect metric
//                                          on it's host, you should implement
//                                          collection method in this method
//
// SetParser(_ parsers.Parser)              if you want to use telegraf's Parser
//                                          interface, you can implement it.
//
// Start(acc telegraf.Accumulator) error    if implement a daemon plugins, you
//                                          should implement this method.
//
// Stop()                                   if implement a daemon plugins, you
//                                          should implement this method.

package nmon

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
	"github.com/influxdata/telegraf/plugins/parsers"
)

// when this plugins be loaded, package init method will register this plugin to
// inputs module.
func init() {
	inputs.Add("nmon_poweragent", newNmonServer)
}

var simpleConfig = `
	##################################
	#     receiver configuration     #
	##################################	
	# ip address and port will be listened, default value: 0.0.0.0:12345
	# ftp debug use it.
	http_listen = "0.0.0.0:12345"
	
	# ftp  username and password 
	ftp_username = "admin"
	ftp_password = "admin"
	
	# the ftp server which will connect to
	ftp_server = "10.10.10.10:21"
	ftp_dirpath = "/upload"
	
	# pull nmon perfdata files period, unit: second
	ftp_pullperiod = 60
	
	
	##################################
	#    processer configuration     #
	##################################
	
	# format , default value: influx
	data_format = "influx"
	
	# the data channal size which used to send data by receiver to processer
	data_chansize = 2000
	
	# how many data processers will be start to process data
	data_threads = 100
`

// http receiver configuration
type receiverConfig struct {
	FtpUsername   string
	FtpPassword   string
	FtpServer     string
	FtpDirPath    string
	FtpPullPeriod int
	HttpListen    string
	WriteDataChan chan<- []byte
	DB            *bolt.DB
}

// NmonServer implement the plugins interface
type NmonServer struct {
	FtpUsername   string `toml:"ftp_username"`
	FtpPassword   string `toml:"ftp_password"`
	FtpServer     string `toml:"ftp_server"`
	FtpDirPath    string `toml:"ftp_dirpath"`
	FtpPullPeriod int    `toml:"ftp_pullperiod"`

	HttpListen string `toml:"http_listen"`

	// common options
	DataFormat   string `toml:"data_format"`
	DataChanSize int    `toml:"data_chansize"`
	DataThreads  int    `toml:"data_threads"`

	receiver   *ftpReceiver // nmon perf data receive server
	processers []*processer

	DBFile string `toml:"db_file"`
	db     *bolt.DB

	dataChan chan []byte
}

// newNmonServer create a new NmonServer and return it
func newNmonServer() telegraf.Input {
	return &NmonServer{
		DataChanSize: 2000,
		DataThreads:  10,
		DBFile:       "nmon_poweragent.db",
	}
}

// implement ServiceInput interface
func (p *NmonServer) SampleConfig() string {
	return simpleConfig
}

// implement ServiceInput interface
func (p *NmonServer) Description() string {
	return "receive nmon message from other agent send to."
}

// implement ServiceInput interface
func (p *NmonServer) Gather(_ telegraf.Accumulator) error {
	return nil
}

// implement ServiceInput interface
func (p *NmonServer) SetParser(_ parsers.Parser) {
}

// implement ServiceInput interface
// nmon plugin does not support running http/ftp mode together currently.
// if want, need make the configration support.
func (p *NmonServer) Start(acc telegraf.Accumulator) error {

	var err error
	// init boltdb
	p.db, err = bolt.Open(p.DBFile, os.ModePerm, &bolt.Options{
		ReadOnly: false,
		Timeout:  time.Duration(10) * time.Second,
	})
	if err != nil {
		return fmt.Errorf("open db %s failed. error: %s", p.DBFile, err)
	}

	p.dataChan = make(chan []byte, p.DataChanSize)
	p.processers = make([]*processer, 0)
	for i := 1; i <= p.DataThreads; i++ {
		ps := newProcesser(i, p.dataChan, acc)
		err := ps.Start()
		if err != nil {
			log.Println(err)
			continue
		}
		p.processers = append(p.processers, ps)
	}

	p.receiver = newftpReceiver(p.ReceiverConfig())
	return p.receiver.Start()

}

// implement ServiceInput interface
func (p *NmonServer) Stop() {
	for _, ps := range p.processers {
		ps.Stop()
	}
	p.receiver.Stop()
}

// Config method return a config struct. it will be used to run a
// server implementation
func (p *NmonServer) ReceiverConfig() receiverConfig {
	return receiverConfig{
		FtpUsername:   p.FtpUsername,
		FtpPassword:   p.FtpPassword,
		FtpServer:     p.FtpServer,
		FtpDirPath:    p.FtpDirPath,
		FtpPullPeriod: p.FtpPullPeriod,
		HttpListen:    p.HttpListen,
		WriteDataChan: p.dataChan,
		DB:            p.db,
	}
}
