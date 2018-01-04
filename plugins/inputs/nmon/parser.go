package nmon

import (
	"fmt"
	"log"
	"strings"

	"github.com/influxdata/telegraf"
)

//
const (
	lineSep  = "\n"    // nmon context lines sep
	infoSep  = "#SEP#" // baseinfo line sep
	pairSep  = "="     // baseinfo keyvalue pair sep
	fieldSep = ","
)

//
const (
	nmonFormat = "nmon"
	jsonFormat = "json"
)

// define error info for procee data
var (
	baseInfoFormatError = fmt.Errorf("the baseinfo format is wrong")
	notnMonContextError = fmt.Errorf("not nmon context")
	fieldsCountError    = fmt.Errorf("number of field is not match")
	notZZZZLineError    = fmt.Errorf("zzzz line  is not match")
	keyNotMatchError    = fmt.Errorf("the key does not match")
)

// metricer
type metricer interface {
	Measurement() string
	Fields() map[string]interface{}
	Tags() map[string]string
}

type nmonMetric struct {
	measurement string
	tags        map[string]string
	fields      map[string]interface{}
}

func newNmonMetric() *nmonMetric {
	return &nmonMetric{
		tags: make(map[string]string),
	}
}

//
func (p *nmonMetric) Measurement() string {
	return p.measurement
}

//
func (p *nmonMetric) Tags() map[string]string {
	return p.tags
}

//
func (p *nmonMetric) AddTag(k, v string) {
	p.tags[k] = v
}

//
func (p *nmonMetric) Fields() map[string]interface{} {
	return p.fields
}

//
type parser interface {
	AddTags(map[string]string)
	Parse(map[string]string, string) error
	Metrics() []metricer
}

func newParser(format string) parser {
	var p parser
	switch strings.ToLower(format) {
	case jsonFormat:
	default:
		p = newNmonParser()
	}

	return p
}

// nmonParser
type nmonParser struct {
	measurement string
	tags        map[string]string
	fields      map[string]map[string]interface{}
}

// newNmonParser create and init a nmonParser
func newNmonParser() *nmonParser {
	return &nmonParser{
		tags:   make(map[string]string),
		fields: make(map[string]map[string]interface{}),
	}
}

//
func (p *nmonParser) init() {
	p.measurement = ""
	p.tags = make(map[string]string)
	p.fields = make(map[string]map[string]interface{})
}

//
func (p *nmonParser) AddTags(tags map[string]string) {
	for k, v := range tags {
		p.tags[k] = v
	}
}

//
func (p *nmonParser) Parse(baseinfo map[string]string, line string) error {
	fields := strings.Split(line, fieldSep)
	if len(fields) < 3 {
		return fieldsCountError
	}

	p.init()

	var err error
	if strings.HasPrefix(fields[0], "CPU") {
		if fields[0] != "CPU_ALL" {
			m := new(aixCPU)
			p.measurement = m.Measurement()
			p.fields, err = m.Parse(line)
			return err
		}
	}

	if strings.HasPrefix(fields[0], "PCPU") {
		if fields[0] != "PCPU_ALL" {
			m := new(aixPCPU)
			p.measurement = m.Measurement()
			p.fields, err = m.Parse(line)
			return err
		}
	}

	if strings.HasPrefix(fields[0], "SCPU") {
		if fields[0] != "SCPU_ALL" {
			m := new(aixSCPU)
			p.measurement = m.Measurement()
			p.fields, err = m.Parse(line)
			return err
		}
	}

	switch strings.ToUpper(fields[0]) {
	case aixKeyCPUALL:
		m := new(aixCPUALL)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyPCPUALL:
		m := new(aixPCPUALL)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeySCPUALL:
		m := new(aixSCPUALL)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyMEM:
		m := new(aixMEM)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyMEMNEW:
		m := new(aixMEMNEW)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyMEMUSE:
		m := new(aixMEMUSE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyPAGE:
		m := new(aixPAGE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyPROC:
		m := new(aixPROC)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyFILE:
		m := new(aixFILE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyTOP:
		m := new(aixTOP)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyPOOLS:
		m := new(aixPOOL)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyLPAR:
		m := new(aixLPAR)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(line)
		return err
	case aixKeyNET:
		m := new(aixNET)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyNETPACKET:
		m := new(aixNETPACKET)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyNETERROR:
		m := new(aixNETERROR)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyNETSIZE:
		m := new(aixNETSIZE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKBUSY:
		m := new(aixDISKBUSY)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKREAD:
		m := new(aixDISKREAD)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKWRITE:
		m := new(aixDISKWRITE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKXFER:
		m := new(aixDISKXFER)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKRXFER:
		m := new(aixDISKRXFER)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKBSIZE:
		m := new(aixDISKBSIZE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKRIO:
		m := new(aixDISKRIO)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKWIO:
		m := new(aixDISKWIO)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKAVGRIO:
		m := new(aixDISKAVGRIO)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyDISKAVGWIO:
		m := new(aixDISKAVGWIO)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyJFSFILE:
		m := new(aixJFSFILE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyJFSINODE:
		m := new(aixJFSINODE)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	case aixKeyIOADAPT:
		m := new(aixIOADAPT)
		p.measurement = m.Measurement()
		p.fields, err = m.Parse(baseinfo, line)
		return err
	default:
		err = keyNotMatchError
		return err
	}
	return nil

}

//
func (p *nmonParser) Metrics() []metricer {
	var metrics = make([]metricer, 0)
	for k, v := range p.fields {
		metric := newNmonMetric()
		metric.measurement = p.measurement
		for k, v := range p.tags {
			metric.AddTag(k, v)
		}
		metric.AddTag("object", k)
		metric.fields = v
		metrics = append(metrics, metric)
	}
	return metrics
}

//
func processData(acc telegraf.Accumulator, src, dataformat string, data []byte) error {
	lines := strings.Split(string(data), "\n")
	if len(lines) < 2 {
		return notnMonContextError
	}

	baseinfo, err := baseInfo(lines[0])
	if err != nil {
		return err
	}

	if !strings.HasPrefix(lines[1], "ZZZZ") {
		return notZZZZLineError
	}

	var (
		nmonparser = newParser(dataformat)
		defaultTag = make(map[string]string)
	)

	defaultTag["ip"] = src

	var metrics = make([]metricer, 0)
	for _, line := range lines[2:] {
		if strings.TrimSpace(line) == "" {
			continue
		}

		err := nmonparser.Parse(baseinfo, line)
		if err != nil {
			log.Println(line, err)
			continue
		}

		nmonparser.AddTags(defaultTag)
		tm := nmonparser.Metrics()
		metrics = append(metrics, tm...)
	}

	for _, metric := range metrics {
		acc.AddGauge(metric.Measurement(), metric.Fields(), metric.Tags())
	}

	return nil
}

// process the first line which client report, it contains some info for
// net or disk or others
func baseInfo(line string) (baseinfo map[string]string, err error) {
	infoPairs := strings.Split(line, infoSep)
	if len(infoPairs) == 0 {
		err = baseInfoFormatError
		return
	}
	baseinfo = make(map[string]string)

	for _, infoPair := range infoPairs {
		kv := strings.Split(infoPair, pairSep)
		if len(kv) != 2 {
			err = baseInfoFormatError
			return
		}
		baseinfo[strings.ToUpper(kv[0])] = kv[1]
	}

	return baseinfo, nil
}
