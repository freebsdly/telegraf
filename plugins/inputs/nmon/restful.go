package nmon

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/influxdata/telegraf"
)

// endpoints
var endpoint = "/metrics"

const (
	// DEFAULT_MAX_BODY_SIZE is the default maximum request body size, in bytes.
	// if the request body is over this size, we will return an HTTP 413 error.
	// 500 MB
	DEFAULT_MAX_BODY_SIZE int64 = 100 * 1024 * 1024
)

// return messages
var (
	maxBodySizeError []byte = []byte("too large http body size")
	methodError      []byte = []byte("only accept post method")
	readBodyError    []byte = []byte("read request body failed")
	processBodyError []byte = []byte("process body data failed")
	okStatus         []byte = []byte("ok")
)

// restfulApiServer implement http server
// it provide a callback function interface to process data
type restfulApiServer struct {
	cfgs config
	mux  *http.ServeMux
	srv  *http.Server
	acc  telegraf.Accumulator
}

// create a new restful api server
func newRestfulApiServer(cfg config) *restfulApiServer {
	return &restfulApiServer{
		cfgs: cfg,
		mux:  http.NewServeMux(),
		srv:  &http.Server{Addr: cfg.Listen},
	}
}

// start the restful api server
func (p *restfulApiServer) Start() error {
	if p.cfgs.MaxBodySize == 0 {
		p.cfgs.MaxBodySize = DEFAULT_MAX_BODY_SIZE
	}

	p.mux.Handle(endpoint, p)
	p.srv.Handler = p.mux
	go func(acc telegraf.Accumulator) {
		if err := p.srv.ListenAndServe(); err != nil {
			acc.AddError(err)
		}
	}(p.acc)
	return nil
}

// stop the restful api server
func (p *restfulApiServer) Stop() error {
	return p.srv.Shutdown(nil)
}

//
func (p *restfulApiServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		rw.Write(methodError)
		return
	}

	if req.ContentLength > p.cfgs.MaxBodySize {
		rw.Write(maxBodySizeError)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		rw.Write(readBodyError)
		return
	}

	err = processData(p.acc, strings.Split(req.RemoteAddr, ":")[0], p.cfgs.DataFormat, body)
	if err != nil {
		log.Println(err)
		rw.Write(processBodyError)
		return
	}
	rw.Write(okStatus)
}

// callback function
func (p *restfulApiServer) SetAccumulator(acc telegraf.Accumulator) {
	p.acc = acc
}
