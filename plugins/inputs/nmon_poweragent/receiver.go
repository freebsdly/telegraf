package nmon

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"net/http"

	"github.com/boltdb/bolt"
	"github.com/jlaffaye/ftp"
)

//
var (
	version = "1.0"
)

//
type ftpReceiver struct {
	sync.RWMutex
	client   *ftp.ServerConn
	addr     string
	username string
	password string
	dirPath  string
	period   int

	db     *bolt.DB
	bucket string

	httpsrv     *http.ServeMux
	debugListen string

	deletelist  []string
	processlist []string
	last        map[string]int64

	wg *sync.WaitGroup

	dataChan chan<- []byte
}

//
func newftpReceiver(cfg receiverConfig) *ftpReceiver {
	return &ftpReceiver{
		addr:        cfg.FtpServer,
		username:    cfg.FtpUsername,
		password:    cfg.FtpPassword,
		dirPath:     cfg.FtpDirPath,
		period:      cfg.FtpPullPeriod,
		last:        make(map[string]int64),
		bucket:      "ftpfilefetcher",
		httpsrv:     http.DefaultServeMux,
		debugListen: cfg.HttpListen,
		wg:          new(sync.WaitGroup),
		dataChan:    cfg.WriteDataChan,
		db:          cfg.DB,
	}
}

//
func (p *ftpReceiver) Start() error {
	var err error

	tx, err := p.db.Begin(true)
	if err != nil {
		return fmt.Errorf("open bolt tx failed. error: %s", err)
	}

	bkt, err := tx.CreateBucketIfNotExists([]byte(p.bucket))
	if err != nil {
		return fmt.Errorf("create bucket failed. error: %s", err)
	}

	bkt.ForEach(p.foreach)

	err = tx.Commit()
	if err != nil {
		if tx.Rollback() != nil {
			log.Printf("tx rollback failed. error: %s\n", err)
		}
		return fmt.Errorf("tx commit failed. error: %s", err)
	}

	// start debug
	p.httpsrv.HandleFunc("/debug", p.debugHandle)
	go func() {
		if err := http.ListenAndServe(p.debugListen, p.httpsrv); err != nil {
			log.Printf("start debug http server failed. error: %s\n", err)
		}
	}()

	go p.getFileList()

	return nil
}

// foreach conv kev/value to cache
func (p *ftpReceiver) foreach(key, value []byte) error {
	iv, err := strconv.ParseInt(string(value), 10, 64)
	if err != nil {
		return err
	}
	p.last[string(key)] = iv
	return nil
}

// debugHandle print cache/processlist/deletelist
func (p *ftpReceiver) debugHandle(resp http.ResponseWriter, req *http.Request) {
	var data = make([]byte, 0)

	data = append(data, []byte("cache: ")...)
	data1, _ := json.Marshal(p.last)
	data = append(data, data1...)
	data = append(data, []byte("\n")...)

	data = append(data, []byte("processlist: ")...)
	data1, _ = json.Marshal(p.processlist)
	data = append(data, data1...)
	data = append(data, []byte("\n")...)

	data = append(data, []byte("deletelist: ")...)
	data1, _ = json.Marshal(p.deletelist)
	data = append(data, data1...)
	resp.Write(data)
}

// stop ftpFileFetcher
func (p *ftpReceiver) Stop() error {
	err := p.db.Close()
	if err != nil {
		log.Printf("close db failed. error: %s\n", err)
	}
	return p.client.Logout()
}

// presist store the last cache into db
func (p *ftpReceiver) presist() {
	tx, err := p.db.Begin(true)
	if err != nil {
		log.Printf("open bolt tx failed. error: %s\n", err)
		return
	}

	bkt := tx.Bucket([]byte(p.bucket))
	for k, v := range p.last {
		bkt.Put([]byte(k), []byte(fmt.Sprintf("%d", v)))
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("tx commit failed. error: %s\n", err)
		if err = tx.Rollback(); err != nil {
			log.Printf("tx rollback failed. error: %s\n", err)
		}
	}
}

//
func (p *ftpReceiver) getFileList() {
	var err error
	ticker := time.NewTicker(time.Duration(p.period) * time.Second)
	for {
		// connect to ftp
		p.client, err = ftp.Connect(p.addr)
		if err != nil {
			log.Printf("connect to ftp failed. error: %s", err)
		} else {
			err = p.client.Login(p.username, p.password)
			if err != nil {
				log.Printf("login to ftp failed. error: %s", err)
			} else {
				lst, err := p.client.List(p.dirPath)
				if err != nil {
					log.Printf("list dir %s failed. error: %s\n", p.dirPath, err)
				} else {
					p.parseFilesList(lst)
					p.presist()
					p.processFiles()
					p.deleteFiles()
				}
			}
		}
		<-ticker.C
	}
}

// if there are more than two tar.gz file of one host,
// check all the files's timestamp, if timestamp less than or equal the timestamp
// cahced, will not process this file, otherwise will process it.
// then delete it.

// 9117-MMA*06B86A1-rc_06B86A1_VIOC3-1516690201.tar.gz
// 9117-MMA is machine model
// rc_06B86A1_VIOC3 is LPARNumberName
// 06B86A1 is serial number
// 1516690201 is report time
func (p *ftpReceiver) parseFilesList(lst []*ftp.Entry) {
	var (
		lparname string
	)
	p.deletelist = make([]string, 0)
	p.processlist = make([]string, 0)

	for _, v := range lst {
		p.deletelist = append(p.deletelist, v.Name)

		array := strings.Split(v.Name, "*")
		if len(array) != 2 {
			log.Printf("filename %s field format wrong when split it by sep *\n", v.Name)
			continue
		}

		array1 := strings.Split(array[1], "-")
		if len(array1) != 3 {
			log.Printf("filename %s field format wrong when split it by sep -\n", v.Name)
			continue
		}

		lparname = array1[1]

		lastTimestamp, exist := p.last[lparname]
		if exist {
			if v.Time.Unix() <= lastTimestamp {
				log.Printf("file %s current timestamp(%d) less than or equal last one(%d), skip it.\n", v.Name, v.Time.Unix(), lastTimestamp)
				continue
			}
		}

		p.last[lparname] = v.Time.Unix()
		p.processlist = append(p.processlist, v.Name)
	}

}

// delete tarball from ftp server
func (p *ftpReceiver) deleteFiles() {
	var err error
	for _, v := range p.deletelist {
		err = p.client.Delete(filepath.Join(p.dirPath, v))
		if err != nil {
			log.Printf("delete file %s from server failed. error: %s\n", v, err)
		}
	}
}

// call goruntines to process multi hosts's nmon tarball
func (p *ftpReceiver) processFiles() {
	for _, v := range p.processlist {
		resp, err := p.client.Retr(filepath.Join(p.dirPath, v))
		if err != nil {
			log.Printf("get file %s from server failed. error: %s\n", v, err)
			continue
		}

		data, err := ioutil.ReadAll(resp)
		if err != nil {
			log.Printf("read from resp of file: %s failed. error: %s\n", v, err)
		}
		go func(name string, d []byte) {
			p.wg.Add(1)
			r, err := nmonTgzFileReader(d)
			if err != nil {
				log.Printf("extract file %s failed, error: %s\n", name, err)
			} else {
				// parse to json format and send to socket
				p.dataChan <- r
			}
			p.wg.Done()
		}(v, data)

		resp.Close()
	}
	p.wg.Wait()
}

// expression the nmon tar.gz file
// there is only one file in the xxx.tar.gz
func nmonTgzFileReader(r []byte) (data []byte, err error) {
	gr, err := gzip.NewReader(bytes.NewReader(r))
	if err != nil {
		err = fmt.Errorf("gzip new reader failed. %s", err)
		return
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	for _, err = tr.Next(); err != io.EOF; _, err = tr.Next() {
		if err != nil {
			return
		}

		data, err = ioutil.ReadAll(tr)
		if err != nil {
			return
		}
	}

	return data, nil
}
