package nmon

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// key name
const (
	aixKeyCPU        = "CPU"
	aixKeyPCPU       = "PCPU"
	aixKeySCPU       = "SCPU"
	aixKeyCPUALL     = "CPU_ALL"
	aixKeyPCPUALL    = "PCPU_ALL"
	aixKeySCPUALL    = "SCPU_ALL"
	aixKeyLPAR       = "LPAR"
	aixKeyPOOLS      = "POOLS"
	aixKeyMEM        = "MEM"
	aixKeyMEMNEW     = "MEMNEW"
	aixKeyMEMUSE     = "MEMUSE"
	aixKeyPAGE       = "PAGE"
	aixKeyPROC       = "PROC"
	aixKeyFILE       = "FILE"
	aixKeyNET        = "NET"
	aixKeyNETPACKET  = "NETPACKET"
	aixKeyNETSIZE    = "NETSIZE"
	aixKeyNETERROR   = "NETERROR"
	aixKeyIOADAPT    = "IOADAPT"
	aixKeyJFSFILE    = "JFSFILE"
	aixKeyJFSINODE   = "JFSINODE"
	aixKeyDISKBUSY   = "DISKBUSY"
	aixKeyDISKREAD   = "DISKREAD"
	aixKeyDISKWRITE  = "DISKWRITE"
	aixKeyDISKXFER   = "DISKXFER"
	aixKeyDISKRXFER  = "DISKRXFER"
	aixKeyDISKBSIZE  = "DISKBSIZE"
	aixKeyDISKRIO    = "DISKRIO"
	aixKeyDISKWIO    = "DISKWIO"
	aixKeyDISKAVGRIO = "DISKAVGRIO"
	aixKeyDISKAVGWIO = "DISKAVGWIO"
	aixKeyTOP        = "TOP"
)

var shortMonthNames = map[string]string{
	"JAN": "Jan",
	"FEB": "Feb",
	"MAR": "Mar",
	"APR": "Apr",
	"MAY": "May",
	"JUN": "Jun",
	"JUL": "Jul",
	"AUG": "Aug",
	"SEP": "Sep",
	"OCT": "Oct",
	"NOV": "Nov",
	"DEC": "Dec",
}

var (
	fieldSep = ","

	// errors
	aixKeyNotTopError = fmt.Errorf("line not top line")
	aixFieldLenError  = fmt.Errorf("field count not match")
)

type metric struct {
	Instance string      `json:"a_id"`
	Counter  string      `json:"tag"`
	Value    interface{} `json:"value"`
}

type perfObject struct {
	Name      string    `json:"device_type"`
	Timestamp time.Time `json:"t"`
	Metrics   []metric  `json:"c"`
}

//
func newCPUPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "cpu",
		Metrics: make([]metric, 0),
	}
}

//
func newMemPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "mem",
		Metrics: make([]metric, 0),
	}
}

//
func newProcessPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "process",
		Metrics: make([]metric, 0),
	}
}

//
func newFilePerfOjbect() *perfObject {
	return &perfObject{
		Name:    "file",
		Metrics: make([]metric, 0),
	}
}

//
func newNetPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "net",
		Metrics: make([]metric, 0),
	}
}

//
func newDiskPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "disk",
		Metrics: make([]metric, 0),
	}
}

//
func newDiskAdapterPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "san",
		Metrics: make([]metric, 0),
	}
}

//
func newLparPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "lpar",
		Metrics: make([]metric, 0),
	}
}

//
func newPoolsPerfOjbect() *perfObject {
	return &perfObject{
		Name:    "pools",
		Metrics: make([]metric, 0),
	}
}

//
type nmonMetric struct {
	name      string
	tags      map[string]string
	fields    map[string]interface{}
	timestamp time.Time
}

//
func (p *nmonMetric) Tags() map[string]string {
	return p.tags
}

func (p *nmonMetric) Fields() map[string]interface{} {
	return p.fields
}

func (p *nmonMetric) Measurement() string {
	return p.name
}

func (p *nmonMetric) Time() time.Time {
	return p.timestamp
}

func (p *nmonMetric) AddTags(tags map[string]string) {
	for k, v := range tags {
		p.tags[k] = v
	}
}

func (p *nmonMetric) AddTag(k, v string) {
	p.tags[k] = v
}

//
type nmonParser struct {
	os          string
	model       string
	sn          string
	hostname    string
	netcards    []string
	disks       []string
	diskAdapter []string
	fs          []string
	ctime       string
	cdate       string
	lpar        map[string]string // name:id
	version     string
	perfObjects []*perfObject
}

//
func (p *nmonParser) SetVersion(v string) {
	p.version = v
}

//
func (p *nmonParser) Json() ([]byte, error) {
	return json.Marshal(p.perfObjects)
}

//
func (p *nmonParser) Metrics() []*nmonMetric {
	ctags := p.Tags()
	var ms = make([]*nmonMetric, 0)
	for _, v := range p.perfObjects {
		m := new(nmonMetric)
		m.fields = make(map[string]interface{})
		m.tags = make(map[string]string)
		m.name = fmt.Sprintf("%s_%s", "nmon", v.Name)
		m.AddTags(ctags)
		m.timestamp = v.Timestamp
		m.AddTag("device_type", v.Name)
		m.AddTag("t", fmt.Sprintf("%s", v.Timestamp.Unix()))

		for _, vv := range v.Metrics {
			m.fields[vv.Counter] = vv.Value
			m.AddTag("a_id", vv.Instance)
		}

		ms = append(ms, m)
	}

	return ms

}

//
func (p *nmonParser) Tags() map[string]string {
	tf, err := timeTranslate(p.ctime, p.cdate)
	if err != nil {
		return nil
	}
	tags := map[string]string{
		"s":   fmt.Sprintf("%s*%s", p.model, p.sn),
		"t_f": fmt.Sprintf("%d", tf),
		"v":   p.version,
	}
	for k, v := range p.lpar {
		tags["l"] = k
		tags["l_id"] = v
	}

	return tags
}

//
func (p *nmonParser) parseAAA(lines []string) error {
	if len(lines) < 3 {
		return fmt.Errorf("AAA format wrong")
	}

	p.lpar = make(map[string]string)

	for _, line := range lines {
		if strings.HasPrefix(line, "AAA,build,") {
			p.os = strings.Split(line, ",")[2]
			continue
		}
		if strings.HasPrefix(line, "AAA,host,") {
			p.hostname = strings.Split(line, ",")[2]
			continue
		}
		if strings.HasPrefix(line, "AAA,time,") {
			p.ctime = strings.Split(line, ",")[2]
			continue
		}
		if strings.HasPrefix(line, "AAA,date,") {
			p.cdate = strings.Split(line, ",")[2]
			continue
		}
		// AAA,SerialNumber,06B86A1
		if strings.HasPrefix(line, "AAA,SerialNumber,") {
			p.sn = strings.Split(line, ",")[2]
			continue
		}
		// AAA,LPARNumberName,3,rc_06B86A1_VIOC3
		if strings.HasPrefix(line, "AAA,LPARNumberName,") {
			array := strings.Split(line, ",")
			p.lpar[array[3]] = array[2]
			continue
		}
		// AAA,MachineType,IBM,9117-MMA
		if strings.HasPrefix(line, "AAA,MachineType,") {
			p.model = strings.Split(line, ",")[3]
			continue
		}

	}
	return nil

}

// BBBN,000,NetworkName,MTU,Mbits,Name
// BBBN,001,en0,1500,10240,Standard Ethernet Network Interface
// BBBN,002,lo0,16896,0,Loopback Network Interface
func (p *nmonParser) parseBBBN(lines []string) error {
	if len(lines) < 2 {
		return fmt.Errorf("BBBN format wrong")
	}

	var err error
	p.netcards = make([]string, 0)

	for _, line := range lines {
		if strings.HasPrefix(line, "BBBN,000,") {
			continue
		}
		array := strings.Split(line, ",")
		if len(array) != 6 {
			err = fmt.Errorf("BBBN line format wrong")
			continue
		}
		p.netcards = append(p.netcards, array[2])
	}
	return err
}

// DISKBUSY,Disk %Busy rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDisks(line string) error {
	if !strings.HasPrefix(line, "DISKBUSY,") {
		return fmt.Errorf("not diskbusy line")
	}
	array := strings.Split(line, ",")
	if len(array) < 3 {
		return fmt.Errorf("diskbusy line format wrong")
	}

	p.disks = make([]string, 0)

	for _, line := range array[2:] {
		p.disks = append(p.disks, line)
	}
	return nil
}

// IOADAPT,Disk Adapter rcvioc03,vscsi0_read-KB/s,vscsi0_write-KB/s,vscsi0_xfer-tps
func (p *nmonParser) parseDiskAdapter(line string) error {
	if !strings.HasPrefix(line, "IOADAPT,Disk Adapter") {
		return fmt.Errorf("not IOADAPT line")
	}
	array := strings.Split(line, ",")
	alen := len(array)
	if len(array) < 3 || (alen-2)%3 != 0 {
		return fmt.Errorf("IOADAPT line format wrong")
	}

	step := (alen - 2) / 3
	p.diskAdapter = make([]string, 0)

	for i := 2; i < 2+step; i++ {
		ay := strings.Split(array[i], "_")
		p.diskAdapter = append(p.diskAdapter, ay[0])
	}
	return nil

}

// JFSINODE,JFS Inode %Used rcvioc03,/,/home,/usr,/var,/tmp,/admin,/opt,/var/adm/ras/livedump
func (p *nmonParser) parseFS(line string) error {
	if !strings.HasPrefix(line, "JFSINODE,") {
		return fmt.Errorf("not JFSINODE line")
	}
	array := strings.Split(line, ",")
	if len(array) < 3 {
		return fmt.Errorf("JFSINODE line format wrong")
	}

	p.fs = make([]string, 0)

	for _, v := range array[2:] {
		p.fs = append(p.fs, v)
	}

	return nil
}

//ZZZZ,T0955,15:59:34,24-JAN-2018
//CPU01,T0955,2.1,2.1,0.0,95.7
//CPU02,T0955,0.2,0.0,0.0,99.8
//CPU03,T0955,0.0,0.0,0.0,100.0
//CPU04,T0955,0.0,0.0,0.0,100.0
//PCPU01,T0955,0.03,0.01,0.00,0.00
//PCPU02,T0955,0.00,0.00,0.00,0.00
//PCPU03,T0955,0.00,0.00,0.00,0.00
//PCPU04,T0955,0.00,0.00,0.00,0.00
//SCPU01,T0955,0.03,0.01,0.00,0.00
//SCPU02,T0955,0.00,0.00,0.00,0.00
//SCPU03,T0955,0.00,0.00,0.00,0.00
//SCPU04,T0955,0.00,0.00,0.00,0.00
//CPU_ALL,T0955,2.9,1.5,0.3,95.3,,4
//PCPU_ALL,T0955,0.03,0.01,0.0,0.00,1.00
//SCPU_ALL,T0955,0.03,0.01,0.0,0.00
//LPAR,T0955,0.046,2,4,4,1.00,128,0.00,1.14,1.14,1,0,2.89,1.49,0.00,0.20,1.45,0.74,0.00,0.10,0,0
//POOLS,T0955,4,4.00,3.10,0.00,0.00,0.00,0.00,0,1.00
//MEM,T0955,1.4,97.5,56.3,499.1,4096.0,512.0
//MEMNEW,T0955,19.7,55.9,23.0,1.4,23.5,72.1
//MEMUSE,T0955,55.9,3.0,90.0,960,1088,55.9,90.0, 1006896.0
//PAGE,T0955,1078.7,0.0,0.3,0.0,0.0,0.0,0.0,0.0
//PROC,T0955,2.38,0.02,233,1572,83,5,2,3,0,0,0,0,0
//FILE,T0955,0,164,0,278437,4461,0,0,0
//NET,T0955,0.2,0.1,0.1,0.1
//NETPACKET,T0955,4.1,0.8,0.5,0.8
//NETSIZE,T0955,51.0,67.6,243.1,67.6
//NETERROR,T0955,0.0,0.0,0.0,0.0,0.0,0.0
//IOADAPT,T0955,35.2,8.1,2.3
//JFSFILE,T0955,29.8,6.0,27.8,87.9,28.4,0.3,50.5,0.1
//JFSINODE,T0955,1.1,0.2,2.3,33.7,1.2,0.0,1.1,0.0
//DISKBUSY,T0955,1.3,0.0,0.0
//DISKREAD,T0955,17.5,0.0,17.5
//DISKWRITE,T0955,8.1,0.0,0.0
//DISKXFER,T0955,2.0,0.0,0.2
//DISKRXFER,T0955,0.2,0.0,0.2
//DISKBSIZE,T0955,12.5,0.0,87.7
//DISKRIO,T0955,0.2,0.0,0.2
//DISKWIO,T0955,1.8,0.0,0.0
//DISKAVGRIO,T0955,87.7,0.0,87.7
//DISKAVGWIO,T0955,4.4,0.0,0.0
//TOP,8519710,T0955,0.17,0.11,0.06,77,32344,3440,27888,2558,1,54,kuxagent,Unclassified
//TOP,7340108,T0955,0.11,0.05,0.06,9,12772,216,12448,762,0,34,aixdp_daemon,Unclassified
func (p *nmonParser) parseZZZZ(lines []string) error {
	if len(lines) < 2 {
		return fmt.Errorf("ZZZZ format wrong")
	}

	if p.perfObjects == nil {
		p.perfObjects = make([]*perfObject, 0)
	}

	var (
		err       error
		timestamp time.Time
	)

	for _, line := range lines {
		if strings.HasPrefix(line, "ZZZZ") {
			timestamp, err = zzzzTimeTranslate(line)
			if err != nil {
				return err
			}
			continue
		}

		fields := strings.Split(line, ",")

		if strings.HasPrefix(fields[0], "CPU") {
			// for cpu
			if fields[0] != "CPU_ALL" {
				err = p.parseCPU(timestamp, line)
				if err != nil {
					return err
				}
			} else {
				// for cpu_all
				err = p.parseCPUALL(timestamp, line)
				if err != nil {
					return err
				}
			}
			continue
		}

		if strings.HasPrefix(fields[0], "PCPU") {
			if fields[0] != "PCPU_ALL" {
				err = p.parsePCPU(timestamp, line)
				if err != nil {
					return err
				}
			} else {
				err = p.parsePCPUALL(timestamp, line)
				if err != nil {
					return err
				}
			}
			continue
		}

		if strings.HasPrefix(fields[0], "SCPU") {
			if fields[0] != "SCPU_ALL" {
				err = p.parseSCPU(timestamp, line)
				if err != nil {
					return err
				}
			} else {
				err = p.parseSCPUALL(timestamp, line)
				if err != nil {
					return err
				}
			}
			continue
		}

		switch strings.ToUpper(fields[0]) {
		case aixKeyMEM:
			err = p.parseMEM(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyMEMNEW:
			err = p.parseMEMNEW(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyMEMUSE:
			err = p.parseMEMUSE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyPAGE:
			err = p.parsePAGE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyPROC:
			err = p.parsePROC(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyFILE:
			err = p.parseFILE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyTOP:
			err = p.parseTOP(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyPOOLS:
			err = p.parsePOOLS(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyLPAR:
			err = p.parseLPAR(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyNET:
			err = p.parseNET(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyNETPACKET:
			err = p.parseNETPACKET(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyNETERROR:
			err = p.parseNETERROR(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyNETSIZE:
			err = p.parseNETSIZE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKBUSY:
			err = p.parseDISKBUSY(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKREAD:
			err = p.parseDISKREAD(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKWRITE:
			err = p.parseDISKWRITE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKXFER:
			err = p.parseDISKXFER(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKRXFER:
			err = p.parseDISKRXFER(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKBSIZE:
			err = p.parseDISKBSIZE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKRIO:
			err = p.parseDISKRIO(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKWIO:
			err = p.parseDISKWIO(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKAVGRIO:
			err = p.parseDISKAVGRIO(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyDISKAVGWIO:
			err = p.parseDISKAVGWIO(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyJFSFILE:
			err = p.parseJFSFILE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyJFSINODE:
			err = p.parseJFSINODE(timestamp, line)
			if err != nil {
				return err
			}
		case aixKeyIOADAPT:
			err = p.parseIOADPAT(timestamp, line)
			if err != nil {
				return err
			}
		default:
			break
		}
	}

	return nil
}

//ZZZZ,T0952,15:56:34,24-JAN-2018
func zzzzTimeTranslate(line string) (time.Time, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 4 {
		return time.Time{}, fmt.Errorf("zzzz line format wrong")
	}

	array := strings.Split(fields[3], "-")
	//"Mon, 02 Jan 2006 15:04:05 -0700"
	s := fmt.Sprintf("%s, %s %s %s %s %s", "Mon", array[0], shortMonthNames[array[1]],
		array[2], fields[2], "+0800")

	return time.Parse(time.RFC1123Z, s)

}

//
func timeTranslate(t, d string) (int64, error) {
	array := strings.Split(d, "-")
	//"Mon, 02 Jan 2006 15:04:05 -0700"
	s := fmt.Sprintf("%s, %s %s %s %s %s", "Mon", array[0], shortMonthNames[array[1]],
		array[2], t, "+0800")

	tt, err := time.Parse(time.RFC1123Z, s)
	if err != nil {
		return 0, err
	}
	return tt.Unix(), nil
}

//
func (p *nmonParser) Parse(data []byte) error {
	var err error
	// split AAA data
	end := bytes.LastIndex(data, []byte("AAA,"))
	if end <= 0 {
		return fmt.Errorf("no AAA string in data")
	}

	lines := strings.Split(string(data[0:end]), "\n")
	err = p.parseAAA(lines)
	if err != nil {
		return fmt.Errorf("parse AAA failed. %s", err)
	}

	// split BBBN data
	start := bytes.Index(data, []byte("BBBN,"))
	if start < 0 {
		return fmt.Errorf("no BBBN string in data")
	}
	end = bytes.LastIndex(data, []byte("BBBN,"))
	if start >= end {
		return fmt.Errorf("only one BBBN string in data")
	}

	last := bytes.Index(data[end:], []byte("\n"))
	if last < 0 {
		return fmt.Errorf("BBBN string format wrong")
	}

	lines = strings.Split(string(data[start:end+last]), "\n")
	err = p.parseBBBN(lines)
	if err != nil {
		return fmt.Errorf("parse netcards failed. %s", err)
	}

	// split DISKBUSY LINE
	start = bytes.Index(data, []byte("DISKBUSY,"))
	if start < 0 {
		return fmt.Errorf("not DISKBUSY line in data")
	}
	end = bytes.Index(data[start:], []byte("\n"))
	if end < 0 {
		return fmt.Errorf("DISKBUSY line format wrong")
	}

	err = p.parseDisks(string(data[start : start+end]))
	if err != nil {
		return fmt.Errorf("parse disks faled. %s\n", err)
	}

	// split IOADAPT line
	start = bytes.Index(data, []byte("IOADAPT,"))
	if start < 0 {
		return fmt.Errorf("not IOADAPT line in data")
	}

	end = bytes.Index(data[start:], []byte("\n"))
	if end < 0 {
		return fmt.Errorf("IOADAPT line format wrong")
	}

	err = p.parseDiskAdapter(string(data[start : start+end]))
	if err != nil {
		return fmt.Errorf("parse disk adapter failed. %s", err)
	}

	// split JFSNODE line
	start = bytes.Index(data, []byte("JFSINODE,"))
	if start < 0 {
		return fmt.Errorf("not JFSINODE line in data")
	}

	end = bytes.Index(data[start:], []byte("\n"))
	if end < 0 {
		return fmt.Errorf("JFSINODE line format wrong")
	}

	err = p.parseFS(string(data[start : start+end]))
	if err != nil {
		return fmt.Errorf("parse fs failed. %s", err)
	}

	// split ZZZZ data
	max := len(data) - 1
	for {
		i := bytes.LastIndex(data[0:max], []byte("ZZZZ,"))
		if i < 0 {
			break
		}
		lines := strings.Split(string(data[i:max]), "\n")
		p.parseZZZZ(lines)
		max = i - 1

	}

	return nil
}

//CPU01,CPU 1 rcvioc03,User%,Sys%,Wait%,Idle%
//CPU02,CPU 2 rcvioc03,User%,Sys%,Wait%,Idle%
//CPU03,CPU 3 rcvioc03,User%,Sys%,Wait%,Idle%
//CPU04,CPU 4 rcvioc03,User%,Sys%,Wait%,Idle%
func (p *nmonParser) parseCPU(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "CPU") || strings.HasPrefix(line, "CPU_ALL") {
		return fmt.Errorf("not cpu line")
	}

	array := strings.Split(line, fieldSep)
	if len(array) != 6 {
		return aixFieldLenError
	}

	var (
		obj = newCPUPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "wait"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "idle"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil

}

//PCPU01,PCPU 1 rcvioc03,User ,Sys ,Wait ,Idle
//PCPU02,PCPU 2 rcvioc03,User ,Sys ,Wait ,Idle
//PCPU03,PCPU 3 rcvioc03,User ,Sys ,Wait ,Idle
//PCPU04,PCPU 4 rcvioc03,User ,Sys ,Wait ,Idle
func (p *nmonParser) parsePCPU(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "PCPU") || strings.HasPrefix(line, "PCPU_ALL") {
		return fmt.Errorf("not cpu line")
	}

	array := strings.Split(line, fieldSep)
	if len(array) != 6 {
		return aixFieldLenError
	}

	var (
		obj = newCPUPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "wait"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "idle"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil

}

//SCPU01,SCPU 1 rcvioc03,User ,Sys ,Wait ,Idle
//SCPU02,SCPU 2 rcvioc03,User ,Sys ,Wait ,Idle
//SCPU03,SCPU 3 rcvioc03,User ,Sys ,Wait ,Idle
//SCPU04,SCPU 4 rcvioc03,User ,Sys ,Wait ,Idle
func (p *nmonParser) parseSCPU(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "SCPU") || strings.HasPrefix(line, "SCPU_ALL") {
		return fmt.Errorf("not cpu line")
	}

	array := strings.Split(line, fieldSep)
	if len(array) != 6 {
		return aixFieldLenError
	}

	var (
		obj = newCPUPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "wait"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "idle"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil

}

//CPU_ALL,CPU Total rcvioc03,User%,Sys%,Wait%,Idle%,Busy,PhysicalCPUs
func (p *nmonParser) parseCPUALL(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "CPU_ALL") {
		return fmt.Errorf("not cpu_all line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 8 {
		return aixFieldLenError
	}

	var (
		obj = newCPUPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "wait"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "idle"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "physicalCPUS"
	m.Value, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)

	return nil
}

//PCPU_ALL,PCPU Total rcvioc03,User  ,Sys  ,Wait  ,Idle  , Entitled Capacity
func (p *nmonParser) parsePCPUALL(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "PCPU_ALL") {
		return fmt.Errorf("not cpu_all line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 7 {
		return aixFieldLenError
	}

	var (
		obj = newCPUPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "wait"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "idle"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "entitledCapacity"
	m.Value, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)

	return nil
}

//SCPU_ALL,SCPU Total rcvioc03,User  ,Sys  ,Wait  ,Idle
func (p *nmonParser) parseSCPUALL(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "SCPU_ALL") {
		return fmt.Errorf("not cpu line")
	}

	array := strings.Split(line, fieldSep)
	if len(array) != 6 {
		return aixFieldLenError
	}

	var (
		obj = newCPUPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "wait"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "idle"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//MEM,Memory rcvioc03,Real Free %,Virtual free %,Real free(MB),Virtual free(MB),Real total(MB),Virtual total(MB)
func (p *nmonParser) parseMEM(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "MEM") {
		return fmt.Errorf("not mem line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 8 {
		return aixFieldLenError
	}

	var (
		obj = newMemPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "RealFree"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "VirtualFree"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "RealFreeMB"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "VirtualFreeMB"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "RealTotalMB"
	m.Value, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "VritualTotalMB"
	m.Value, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//MEMNEW,Memory New rcvioc03,Process%,FScache%,System%,Free%,Pinned%,User%
func (p *nmonParser) parseMEMNEW(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "MEMNEW") {
		return fmt.Errorf("not memnew line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 8 {
		return aixFieldLenError
	}

	var (
		obj = newMemPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "process"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "fscache"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "system"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "free"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pinned"
	m.Value, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "user"
	m.Value, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//MEMUSE,Memory Use rcvioc03,%numperm,%minperm,%maxperm,minfree,maxfree,%numclient,%maxclient, lruable pages
func (p *nmonParser) parseMEMUSE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "MEMUSE") {
		return fmt.Errorf("not memuse line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 10 {
		return aixFieldLenError
	}

	var (
		obj = newMemPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "numperm"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "minperm"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "maxperm"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "minfree"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "maxfree"
	m.Value, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "numclient"
	m.Value, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "maxclient"
	m.Value, err = strconv.ParseFloat(array[8], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//PAGE,Paging rcvioc03,faults,pgin,pgout,pgsin,pgsout,reclaims,scans,cycles
func (p *nmonParser) parsePAGE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "PAGE") {
		return fmt.Errorf("not page line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 10 {
		return aixFieldLenError

	}

	var (
		obj = newMemPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "faults"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pgin"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pgout"
	m.Value, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pgsin"
	m.Value, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pgsout"
	m.Value, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "reclaims"
	m.Value, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "scans"
	m.Value, err = strconv.ParseFloat(array[8], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "cycles"
	m.Value, err = strconv.ParseFloat(array[9], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//PROC,Processes rcvioc03,Runnable,Swap-in,pswitch,syscall,read,write,fork,exec,sem,msg,asleep_bufio,asleep_rawio,asleep_diocio
func (p *nmonParser) parsePROC(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "PROC") {
		return fmt.Errorf("not proc line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 15 {
		return aixFieldLenError

	}

	var (
		obj = newProcessPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "runnable"
	m.Value, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "swapin"
	m.Value, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pswitch"
	m.Value, err = strconv.ParseInt(array[4], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "syscall"
	m.Value, err = strconv.ParseInt(array[5], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "read"
	m.Value, err = strconv.ParseInt(array[6], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "write"
	m.Value, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "fork"
	m.Value, err = strconv.ParseInt(array[8], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "exec"
	m.Value, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sem"
	m.Value, err = strconv.ParseInt(array[10], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "msg"
	m.Value, err = strconv.ParseInt(array[11], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "asleepbufio"
	m.Value, err = strconv.ParseInt(array[12], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "asleeprawio"
	m.Value, err = strconv.ParseInt(array[13], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "asleepdiocio"
	m.Value, err = strconv.ParseInt(array[14], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//FILE,File I/O rcvioc03,iget,namei,dirblk,readch,writech,ttyrawch,ttycanch,ttyoutch
func (p *nmonParser) parseFILE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "FILE") {
		return fmt.Errorf("not file line")
	}
	array := strings.Split(line, fieldSep)
	if len(array) != 10 {
		return aixFieldLenError

	}

	var (
		obj = newFilePerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "iget"
	m.Value, err = strconv.ParseInt(array[2], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "namei"
	m.Value, err = strconv.ParseInt(array[3], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "dirblk"
	m.Value, err = strconv.ParseInt(array[4], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "readch"
	m.Value, err = strconv.ParseInt(array[5], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "writech"
	m.Value, err = strconv.ParseInt(array[6], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ttyrawch"
	m.Value, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ttycanch"
	m.Value, err = strconv.ParseInt(array[8], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ttyoutch"
	m.Value, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//NET,Network I/O rcvioc03,en0-read-KB/s,lo0-read-KB/s,en0-write-KB/s,lo0-write-KB/s
func (p *nmonParser) parseNET(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "NET") {
		return fmt.Errorf("not net line")
	}

	fields := strings.Split(line, fieldSep)
	step := len(p.netcards)
	if len(fields) != (2 + 2*step) {
		return aixFieldLenError
	}

	var (
		obj = newNetPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.netcards {
		m.Instance = v

		m.Counter = "read"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "write"
		m.Value, err = strconv.ParseFloat(fields[i+2+step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)
	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//NETPACKET,Network Packets rcvioc03,en0-reads/s,lo0-reads/s,en0-writes/s,lo0-writes/s
func (p *nmonParser) parseNETPACKET(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "NETPACKET") {
		return fmt.Errorf("not netpacket line")
	}

	fields := strings.Split(line, fieldSep)
	step := len(p.netcards)
	if len(fields) != (2 + 2*step) {
		return aixFieldLenError
	}

	var (
		obj = newNetPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.netcards {
		m.Instance = v

		m.Counter = "packetread"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "packetwrite"
		m.Value, err = strconv.ParseFloat(fields[i+2+step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)
	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//NETSIZE,Network Size rcvioc03,en0-readsize,lo0-readsize,en0-writesize,lo0-writesize
func (p *nmonParser) parseNETSIZE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "NETSIZE") {
		return fmt.Errorf("not netsize line")
	}

	fields := strings.Split(line, fieldSep)
	step := len(p.netcards)
	if len(fields) != (2 + 2*step) {
		return aixFieldLenError
	}

	var (
		obj = newNetPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.netcards {
		m.Instance = v

		m.Counter = "readsize"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "writesize"
		m.Value, err = strconv.ParseFloat(fields[i+2+step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)
	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//NETERROR,Network Errors rcvioc03,en0-ierrs,lo0-ierrs,en0-oerrs,lo0-oerrs,en0-collisions,lo0-collisions
func (p *nmonParser) parseNETERROR(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "NETERROR") {
		return fmt.Errorf("not neterror line")
	}

	fields := strings.Split(line, fieldSep)
	step := len(p.netcards)
	if len(fields) != (2 + 3*step) {
		return aixFieldLenError
	}

	var (
		obj = newNetPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.netcards {
		m.Instance = v

		m.Counter = "ierrs"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "oerrs"
		m.Value, err = strconv.ParseFloat(fields[i+2+step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "collisions"
		m.Value, err = strconv.ParseFloat(fields[i+2+2*step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)
	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKBUSY,Disk %Busy rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKBUSY(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKBUSY") {
		return fmt.Errorf("not diskbusy line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "busy"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKREAD,Disk Read KB/s rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKREAD(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKREAD") {
		return fmt.Errorf("not diskread line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "read"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKWRITE,Disk Write KB/s rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKWRITE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKWRITE") {
		return fmt.Errorf("not diskwrite line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "write"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKXFER,Disk transfers per second rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKXFER(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKXFER") {
		return fmt.Errorf("not DISKXFER line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "xfer"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKRXFER,Transfers from disk (reads) per second rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKRXFER(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKRXFER") {
		return fmt.Errorf("not DISKRXFER line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "rxfer"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKBSIZE,Disk Block Size rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKBSIZE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKBSIZE") {
		return fmt.Errorf("not DISKBSIZE line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "bsize"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKRIO,Disk IO Reads per second rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKRIO(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKRIO") {
		return fmt.Errorf("not DISKRIO line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "rio"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKWIO,Disk IO Writes per second rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKWIO(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKWIO") {
		return fmt.Errorf("not DISKWIO line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "wio"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKAVGRIO,Disk IO Average Reads per second rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKAVGRIO(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKAVGRIO") {
		return fmt.Errorf("not DISKAVGRIO line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "avgrio"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//DISKAVGWIO,Disk IO Average Writes per second rcvioc03,hdisk0,cd0,hdisk1
func (p *nmonParser) parseDISKAVGWIO(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "DISKAVGWIO") {
		return fmt.Errorf("not DISKAVGWIO line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.disks)) {
		return aixFieldLenError
	}

	var (
		obj = newDiskPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.disks {
		m.Instance = v

		m.Counter = "avgwio"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//IOADAPT,Disk Adapter rcvioc03,vscsi0_read-KB/s,vscsi0_write-KB/s,vscsi0_xfer-tps
func (p *nmonParser) parseIOADPAT(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "IOADAPT") {
		return fmt.Errorf("not IOADAPT line")
	}
	fields := strings.Split(line, fieldSep)
	step := len(p.diskAdapter)
	if len(fields) != (2 + 3*step) {
		return aixFieldLenError
	}

	var (
		obj = newDiskAdapterPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.diskAdapter {
		m.Instance = v

		m.Counter = "read"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "write"
		m.Value, err = strconv.ParseFloat(fields[i+2+step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

		m.Counter = "xfer-tps"
		m.Value, err = strconv.ParseFloat(fields[i+2+2*step], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//LPAR,Logical Partition rcvioc03,PhysicalCPU,virtualCPUs,logicalCPUs,poolCPUs,entitled,weight,PoolIdle,usedAllCPU%,usedPoolCPU%,SharedCPU,Capped,EC_User%,EC_Sys%,EC_Wait%,EC_Idle%,VP_User%,VP_Sys%,VP_Wait%,VP_Idle%,Folded,Pool_id
func (p *nmonParser) parseLPAR(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "LPAR") {
		return fmt.Errorf("not LPAR line")
	}
	array := strings.Split(line, fieldSep)

	if len(array) != 23 {
		return aixFieldLenError

	}

	var (
		obj = newLparPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[0]

	m.Counter = "physicalcpu"
	m.Value, err = strconv.ParseFloat(array[2], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "virtualcpus"
	m.Value, err = strconv.ParseFloat(array[3], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "logicalcpus"
	m.Value, err = strconv.ParseFloat(array[4], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "poocpus"
	m.Value, err = strconv.ParseFloat(array[5], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "entitled"
	m.Value, err = strconv.ParseFloat(array[6], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "weight"
	m.Value, err = strconv.ParseFloat(array[7], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "poolidle"
	m.Value, err = strconv.ParseFloat(array[8], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "usedallcpu"
	m.Value, err = strconv.ParseFloat(array[9], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "usedpoolcpu"
	m.Value, err = strconv.ParseFloat(array[10], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sharedcpu"
	m.Value, err = strconv.ParseFloat(array[11], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "capped"
	m.Value, err = strconv.ParseFloat(array[12], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ecuser"
	m.Value, err = strconv.ParseFloat(array[13], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ecsys"
	m.Value, err = strconv.ParseFloat(array[14], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ecwait"
	m.Value, err = strconv.ParseFloat(array[15], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ecidle"
	m.Value, err = strconv.ParseFloat(array[16], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "vpuser"
	m.Value, err = strconv.ParseFloat(array[17], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "vpsys"
	m.Value, err = strconv.ParseFloat(array[18], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "vpwait"
	m.Value, err = strconv.ParseFloat(array[19], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "vpidle"
	m.Value, err = strconv.ParseFloat(array[20], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "folded"
	m.Value, err = strconv.ParseFloat(array[21], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "poolid"
	m.Value, err = strconv.ParseInt(array[22], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil

}

//POOLS,Multiple CPU Pools rcvioc03,shcpus_in_sys,max_pool_capacity,entitled_pool_capacity,pool_max_time,pool_busy_time,shcpu_tot_time,shcpu_busy_time,Pool_id,entitled
func (p *nmonParser) parsePOOLS(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "POOLS") {
		return fmt.Errorf("not POOLS line")
	}
	array := strings.Split(line, fieldSep)

	if len(array) != 11 {
		return aixFieldLenError
	}

	var (
		obj = newPoolsPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[9]

	m.Counter = "shcpu_in_sys"
	m.Value, err = strconv.ParseFloat(array[2], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "max_pool_capacity"
	m.Value, err = strconv.ParseFloat(array[3], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "entitled_pool_capacity"
	m.Value, err = strconv.ParseFloat(array[4], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pool_max_time"
	m.Value, err = strconv.ParseFloat(array[5], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "pool_busy_time"
	m.Value, err = strconv.ParseFloat(array[6], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "shcpu_tot_time"
	m.Value, err = strconv.ParseFloat(array[7], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "shcpu_busy_time"
	m.Value, err = strconv.ParseFloat(array[8], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "entitled"
	m.Value, err = strconv.ParseFloat(array[10], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//JFSFILE,JFS Filespace %Used rcvioc03,/,/home,/usr,/var,/tmp,/admin,/opt,/var/adm/ras/livedump
func (p *nmonParser) parseJFSFILE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "JFSFILE") {
		return fmt.Errorf("not JFSFILE line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.fs)) {
		return aixFieldLenError
	}

	var (
		obj = newFilePerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.fs {
		m.Instance = v

		m.Counter = "u_p"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//JFSINODE,JFS Inode %Used rcvioc03,/,/home,/usr,/var,/tmp,/admin,/opt,/var/adm/ras/livedump
func (p *nmonParser) parseJFSINODE(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "JFSINODE") {
		return fmt.Errorf("not JFSINODE line")
	}
	fields := strings.Split(line, fieldSep)
	if len(fields) != (2 + len(p.fs)) {
		return aixFieldLenError
	}

	var (
		obj = newFilePerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp

	for i, v := range p.fs {
		m.Instance = v

		m.Counter = "inode"
		m.Value, err = strconv.ParseFloat(fields[i+2], 10)
		if err != nil {
			return err
		}
		obj.Metrics = append(obj.Metrics, m)

	}

	p.perfObjects = append(p.perfObjects, obj)
	return nil
}

//TOP,%CPU Utilisation
//TOP,+PID,Time,%CPU,%Usr,%Sys,Threads,Size,ResText,ResData,CharIO,%RAM,Paging,Command,WLMclass
//TOP,8519710,T0955,0.17,0.11,0.06,77,32344,3440,27888,2558,1,54,kuxagent,Unclassified
//TOP,7340108,T0955,0.11,0.05,0.06,9,12772,216,12448,762,0,34,aixdp_daemon,Unclassified
func (p *nmonParser) parseTOP(timestamp time.Time, line string) error {
	if !strings.HasPrefix(line, "TOP") {
		return fmt.Errorf("not TOP line")
	}
	array := strings.Split(line, fieldSep)

	if len(array) != 15 {
		return aixFieldLenError
	}

	var (
		obj = newProcessPerfOjbect()
		m   metric
		err error
	)

	obj.Timestamp = timestamp
	m.Instance = array[1]

	m.Counter = "cpu"
	m.Value, err = strconv.ParseFloat(array[3], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "usr"
	m.Value, err = strconv.ParseFloat(array[4], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "sys"
	m.Value, err = strconv.ParseFloat(array[5], 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "threads"
	m.Value, err = strconv.ParseInt(array[6], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "size"
	m.Value, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "restext"
	m.Value, err = strconv.ParseInt(array[8], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "resdata"
	m.Value, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "chario"
	m.Value, err = strconv.ParseInt(array[10], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "ram"
	m.Value, err = strconv.ParseInt(array[11], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	m.Counter = "paging"
	m.Value, err = strconv.ParseInt(array[12], 10, 64)
	if err != nil {
		return err
	}
	obj.Metrics = append(obj.Metrics, m)

	p.perfObjects = append(p.perfObjects, obj)
	return nil

}
