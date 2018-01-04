package nmon

import (
	"fmt"
	"strconv"
	"strings"
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

// measurement
var (
	aixMeasurementPrefix    = "nmon"
	aixMEMMeasurement       = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyMEM)
	aixMEMUSEMeasurement    = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyMEMUSE)
	aixMEMNEWMeasurement    = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyMEMNEW)
	aixPAGEMeasurement      = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyPAGE)
	aixPROCMeasurement      = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyPROC)
	aixFILEMeasurement      = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyFILE)
	aixNETMeasurement       = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyNET)
	aixCPUALLMeasurement    = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyCPUALL)
	aixCPUMeasurement       = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyCPU)
	aixSCPUALLMeasurement   = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeySCPUALL)
	aixPCPUALLMeasurement   = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyPCPUALL)
	aixPCPUMeasurement      = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyPCPU)
	aixSCPUMeasurement      = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeySCPU)
	aixTOPMeasurement       = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyTOP)
	aixPOOLSMeasurement     = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyPOOLS)
	aixLPARMeasurement      = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyLPAR)
	aixNETPACKETMeasurement = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyNETPACKET)
	aixNETSIZEMeasurement   = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyNETSIZE)
	aixNETERRORMeasurement  = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyNETERROR)
	aixDISKMeasurement      = fmt.Sprintf("%s", aixMeasurementPrefix)
	aixJFSMeasurement       = fmt.Sprintf("%s", aixMeasurementPrefix)
	aixIOADAPTMeasurement   = fmt.Sprintf("%s_%s", aixMeasurementPrefix, aixKeyIOADAPT)
)

// errors
var (
	aixKeyNotTopError = fmt.Errorf("line not top line")
	aixFieldLenError  = fmt.Errorf("field count not match")
)

// common struct for aix cpu/pcpu/scpu
// CPU01,CPU 1 rcvioc03,User%,Sys%,Wait%,Idle%
// PCPU01,PCPU 1 rcvioc03,User ,Sys ,Wait ,Idle
// SCPU01,SCPU 1 rcvioc03,User ,Sys ,Wait ,Idle
type aixCPU struct {
	User float64
	Sys  float64
	Wait float64
	Idle float64
}

//
func (p *aixCPU) Measurement() string {
	return aixCPUMeasurement
}

//
func (p *aixCPU) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 6 {
		err = aixFieldLenError
		return
	}

	p.User, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.Sys, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.Wait, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.Idle, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"user": p.User,
			"sys":  p.Sys,
			"wait": p.Wait,
			"idle": p.Idle,
		},
	}, nil
}

//
type aixPCPU struct {
	aixCPU
}

//
func (p *aixPCPU) Measurement() string {
	return aixPCPUMeasurement
}

//
type aixSCPU struct {
	aixCPU
}

//
func (p *aixSCPU) Measurement() string {
	return aixSCPUMeasurement
}

// struct for cpu_all, like this:
// CPU_ALL,CPU Total rcvioc03,User%,Sys%,Wait%,Idle%,Busy,PhysicalCPUs
type aixCPUALL struct {
	aixCPU
	PhysicalCPUS int64
}

//
func (p *aixCPUALL) Measurement() string {
	return aixCPUALLMeasurement
}

//
func (p *aixCPUALL) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 8 {
		err = aixFieldLenError
		return
	}

	p.User, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.Sys, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.Wait, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.Idle, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}
	// this field is empty
	//	p.Busy, err = strconv.ParseFloat(array[6], 10)
	//	if err != nil {
	//		return
	//	}

	p.PhysicalCPUS, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"user":         p.User,
			"sys":          p.Sys,
			"wait":         p.Wait,
			"idle":         p.Idle,
			"physicalCPUS": p.PhysicalCPUS,
		},
	}, nil

}

// pcpu_all
// PCPU_ALL,PCPU Total rcvioc03,User  ,Sys  ,Wait  ,Idle  , Entitled Capacity
type aixPCPUALL struct {
	aixCPU
	EntitledCapacity float64
}

//
func (p *aixPCPUALL) Measurement() string {
	return aixPCPUALLMeasurement
}

//
func (p *aixPCPUALL) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 7 {
		err = aixFieldLenError
		return
	}

	p.User, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.Sys, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.Wait, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.Idle, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}

	p.EntitledCapacity, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"user":             p.User,
			"sys":              p.Sys,
			"wait":             p.Wait,
			"idle":             p.Idle,
			"entitledCapacity": p.EntitledCapacity,
		},
	}, nil

}

// scpu_all
type aixSCPUALL struct {
	aixCPU
}

//
func (p *aixSCPUALL) Measurement() string {
	return aixSCPUALLMeasurement
}

// MEM,Memory rcvioc03,Real Free %,Virtual free %,Real free(MB),Virtual free(MB),
// Real total(MB),Virtual total(MB
type aixMEM struct {
	RealFreePrecent    float64
	VirtualFreePrecent float64
	RealFreeMB         float64
	VirtualFreeMB      float64
	RealTotalMB        float64
	VirtualTotalMB     float64
}

//
func (p *aixMEM) Measurement() string {
	return aixMEMMeasurement
}

//
func (p *aixMEM) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 8 {
		err = aixFieldLenError
		return
	}

	p.RealFreePrecent, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.VirtualFreePrecent, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.RealFreeMB, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.VirtualFreeMB, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}

	p.RealTotalMB, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return
	}

	p.VirtualTotalMB, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"RealFreePrecent":    p.RealFreePrecent,
			"VirtualFreePrecent": p.VirtualFreePrecent,
			"RealFreeMB":         p.RealFreeMB,
			"VirtualFreeMB":      p.VirtualFreeMB,
			"RealTotalMB":        p.RealTotalMB,
			"VirtualTotalMB":     p.VirtualTotalMB,
		},
	}, nil

}

// MEMNEW,Memory New rcvioc03,Process%,FScache%,System%,Free%,Pinned%,User%
type aixMEMNEW struct {
	Process float64
	FSCache float64
	System  float64
	Free    float64
	Pinned  float64
	User    float64
}

func (p *aixMEMNEW) Measurement() string {
	return aixMEMNEWMeasurement
}

func (p *aixMEMNEW) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 8 {
		err = aixFieldLenError
		return
	}

	p.Process, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.FSCache, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.System, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.Free, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}

	p.Pinned, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return
	}

	p.User, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"process": p.Process,
			"fscache": p.FSCache,
			"system":  p.System,
			"free":    p.Free,
			"pinned":  p.Pinned,
			"usr":     p.User,
		},
	}, nil

}

// MEMUSE,Memory Use rcvioc03,%numperm,%minperm,%maxperm,minfree,maxfree,
// %numclient,%maxclient, lruable pages
type aixMEMUSE struct {
	NumPerm      float64
	MinPerm      float64
	MaxPerm      float64
	MinFree      float64
	MaxFree      float64
	NumClient    float64
	MaxClient    float64
	LrualbePages float64
}

//
func (p *aixMEMUSE) Measurement() string {
	return aixMEMUSEMeasurement
}

//
func (p *aixMEMUSE) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 10 {
		err = aixFieldLenError
		return
	}

	p.NumPerm, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.MinPerm, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.MaxPerm, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.MinFree, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}

	p.MaxFree, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return
	}

	p.NumClient, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return
	}

	p.MaxClient, err = strconv.ParseFloat(array[8], 10)
	if err != nil {
		return
	}

	p.LrualbePages, err = strconv.ParseFloat(strings.TrimSpace(array[9]), 10)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"NumPerm":      p.NumPerm,
			"MinPerm":      p.MinPerm,
			"MaxPerm":      p.MaxPerm,
			"MinFree":      p.MinFree,
			"MaxFree":      p.MaxFree,
			"NumClient":    p.NumClient,
			"MaxClient":    p.MaxClient,
			"LrualbePages": p.LrualbePages,
		},
	}, nil

}

// PAGE,Paging rcvioc03,faults,pgin,pgout,pgsin,pgsout,reclaims,scans,cycles
type aixPAGE struct {
	Faults   float64
	PgIn     float64
	PgOut    float64
	PgsIn    float64
	PgsOut   float64
	Reclaims float64
	Scans    float64
	Cycles   float64
}

//
func (p *aixPAGE) Measurement() string {
	return aixPAGEMeasurement
}

//
func (p *aixPAGE) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 10 {
		err = aixFieldLenError
		return
	}

	p.Faults, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.PgIn, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.PgOut, err = strconv.ParseFloat(array[4], 10)
	if err != nil {
		return
	}

	p.PgsIn, err = strconv.ParseFloat(array[5], 10)
	if err != nil {
		return
	}

	p.PgsOut, err = strconv.ParseFloat(array[6], 10)
	if err != nil {
		return
	}

	p.Reclaims, err = strconv.ParseFloat(array[7], 10)
	if err != nil {
		return
	}

	p.Scans, err = strconv.ParseFloat(array[8], 10)
	if err != nil {
		return
	}

	p.Cycles, err = strconv.ParseFloat(array[9], 10)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"faults":   p.Faults,
			"pgin":     p.PgIn,
			"pgout":    p.PgOut,
			"pgsin":    p.PgsIn,
			"pgsout":   p.PgsOut,
			"reclaims": p.Reclaims,
			"scans":    p.Scans,
			"cycles":   p.Cycles,
		},
	}, nil

}

// PROC,Processes rcvioc03,Runnable,Swap-in,pswitch,syscall,read,write,fork,
// exec,sem,msg,asleep_bufio,asleep_rawio,asleep_diocio
type aixPROC struct {
	Runnable     float64
	SwapIn       float64
	Pswitch      int64
	Syscall      int64
	Read         int64
	Write        int64
	Fork         int64
	Exec         int64
	Sem          int64
	Msg          int64
	ASleepBufIO  int64
	ASleepRawIO  int64
	ASleepDiocIO int64
}

//
func (p *aixPROC) Measurement() string {
	return aixPROCMeasurement
}

//
func (p *aixPROC) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 15 {
		err = aixFieldLenError
		return
	}

	p.Runnable, err = strconv.ParseFloat(array[2], 10)
	if err != nil {
		return
	}

	p.SwapIn, err = strconv.ParseFloat(array[3], 10)
	if err != nil {
		return
	}

	p.Pswitch, err = strconv.ParseInt(array[4], 10, 64)
	if err != nil {
		return
	}

	p.Syscall, err = strconv.ParseInt(array[5], 10, 64)
	if err != nil {
		return
	}

	p.Read, err = strconv.ParseInt(array[6], 10, 64)
	if err != nil {
		return
	}

	p.Write, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return
	}

	p.Fork, err = strconv.ParseInt(array[8], 10, 64)
	if err != nil {
		return
	}

	p.Exec, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return
	}

	p.Sem, err = strconv.ParseInt(array[10], 10, 64)
	if err != nil {
		return
	}

	p.Msg, err = strconv.ParseInt(array[11], 10, 64)
	if err != nil {
		return
	}

	p.ASleepBufIO, err = strconv.ParseInt(array[12], 10, 64)
	if err != nil {
		return
	}

	p.ASleepRawIO, err = strconv.ParseInt(array[13], 10, 64)
	if err != nil {
		return
	}

	p.ASleepDiocIO, err = strconv.ParseInt(array[14], 10, 64)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"runnable":     p.Runnable,
			"swapin":       p.SwapIn,
			"pswithc":      p.Pswitch,
			"syscall":      p.Syscall,
			"read":         p.Read,
			"write":        p.Write,
			"fork":         p.Fork,
			"exec":         p.Exec,
			"sem":          p.Sem,
			"msg":          p.Msg,
			"ASleepBufIO":  p.ASleepBufIO,
			"ASleepRawIO":  p.ASleepRawIO,
			"ASleepDiocIO": p.ASleepDiocIO,
		},
	}, nil

}

// FILE,File I/O rcvioc03,iget,namei,dirblk,readch,writech,ttyrawch,
// ttycanch,ttyoutch
type aixFILE struct {
	Iget     int64
	NameI    int64
	DirBLK   int64
	ReadCH   int64
	WriteCH  int64
	TtyRawCH int64
	TtyCanCH int64
	TtyOutCH int64
}

//
func (p *aixFILE) Measurement() string {
	return aixFILEMeasurement
}

//
func (p *aixFILE) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)
	if len(array) != 10 {
		err = aixFieldLenError
		return
	}

	p.Iget, err = strconv.ParseInt(array[2], 10, 64)
	if err != nil {
		return
	}

	p.NameI, err = strconv.ParseInt(array[3], 10, 64)
	if err != nil {
		return
	}

	p.DirBLK, err = strconv.ParseInt(array[4], 10, 64)
	if err != nil {
		return
	}
	p.ReadCH, err = strconv.ParseInt(array[5], 10, 64)
	if err != nil {
		return
	}

	p.WriteCH, err = strconv.ParseInt(array[6], 10, 64)
	if err != nil {
		return
	}

	p.TtyRawCH, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return
	}
	p.TtyCanCH, err = strconv.ParseInt(array[8], 10, 64)
	if err != nil {
		return
	}
	p.TtyOutCH, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		array[0]: map[string]interface{}{
			"iget":     p.Iget,
			"namei":    p.NameI,
			"dirblk":   p.DirBLK,
			"readch":   p.ReadCH,
			"writech":  p.WriteCH,
			"ttyrawch": p.TtyRawCH,
			"ttycanch": p.TtyCanCH,
			"ttyoutch": p.TtyOutCH,
		},
	}, nil

}

// TOP,%CPU Utilisation
// TOP,+PID,Time,%CPU,%Usr,%Sys,Threads,Size,ResText,ResData,CharIO,%RAM,Paging,Command,WLMclass
type aixTOP struct {
	PID      int64
	Time     string
	CPU      float64
	Usr      float64
	Sys      float64
	Threads  int64
	Size     int64
	ResText  int64
	ResData  int64
	CharIO   int64
	RAM      int64
	Paging   int64
	Command  string
	WLWClass string
}

//
func (p *aixTOP) Measurement() string {
	return aixTOPMeasurement
}

//
func (p *aixTOP) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)

	if len(array) != 15 {
		err = aixFieldLenError
		return
	}

	p.PID, err = strconv.ParseInt(array[1], 10, 64)
	if err != nil {
		return
	}
	p.Time = array[2]
	p.CPU, err = strconv.ParseFloat(array[3], 64)
	if err != nil {
		return
	}
	p.Usr, err = strconv.ParseFloat(array[4], 64)
	if err != nil {
		return
	}
	p.Sys, err = strconv.ParseFloat(array[5], 64)
	if err != nil {
		return
	}
	p.Threads, err = strconv.ParseInt(array[6], 10, 64)
	if err != nil {
		return
	}
	p.Size, err = strconv.ParseInt(array[7], 10, 64)
	if err != nil {
		return
	}
	p.ResText, err = strconv.ParseInt(array[8], 10, 64)
	if err != nil {
		return
	}
	p.ResData, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return
	}
	p.CharIO, err = strconv.ParseInt(array[10], 10, 64)
	if err != nil {
		return
	}
	p.RAM, err = strconv.ParseInt(array[11], 10, 64)
	if err != nil {
		return
	}
	p.Paging, err = strconv.ParseInt(array[12], 10, 64)
	if err != nil {
		return
	}
	p.Command = array[13]
	p.WLWClass = array[14]

	return map[string]map[string]interface{}{
		strconv.Itoa(int(p.PID)): map[string]interface{}{
			"cpu":     p.CPU,
			"usr":     p.Usr,
			"sys":     p.Sys,
			"size":    p.Size,
			"restext": p.ResText,
			"resdata": p.ResData,
			"chario":  p.CharIO,
			"ram":     p.RAM,
			"paging":  p.Paging,
		},
	}, nil
}

// POOLS,Multiple CPU Pools rcvioc03,shcpus_in_sys,max_pool_capacity,
// entitled_pool_capacity,pool_max_time,pool_busy_time,shcpu_tot_time,
// shcpu_busy_time,Pool_id,entitled
type aixPOOL struct {
	ShCPUSInSys          float64
	MaxPoolCapacity      float64
	EntitledPoolCapacity float64
	PoolMaxTime          float64
	PoolBusyTime         float64
	ShCPUTotTime         float64
	ShCPUBusyTime        float64
	PoolId               int64
	Entitled             float64
}

//
func (p *aixPOOL) Measurement() string {
	return aixPOOLSMeasurement
}

//
func (p *aixPOOL) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)

	if len(array) != 11 {
		err = aixFieldLenError
		return
	}

	p.ShCPUSInSys, err = strconv.ParseFloat(array[2], 64)
	if err != nil {
		return
	}

	p.MaxPoolCapacity, err = strconv.ParseFloat(array[3], 64)
	if err != nil {
		return
	}
	p.EntitledPoolCapacity, err = strconv.ParseFloat(array[4], 64)
	if err != nil {
		return
	}
	p.PoolMaxTime, err = strconv.ParseFloat(array[5], 64)
	if err != nil {
		return
	}
	p.PoolBusyTime, err = strconv.ParseFloat(array[6], 64)
	if err != nil {
		return
	}
	p.ShCPUTotTime, err = strconv.ParseFloat(array[7], 64)
	if err != nil {
		return
	}
	p.ShCPUBusyTime, err = strconv.ParseFloat(array[8], 64)
	if err != nil {
		return
	}
	p.PoolId, err = strconv.ParseInt(array[9], 10, 64)
	if err != nil {
		return
	}

	p.Entitled, err = strconv.ParseFloat(array[8], 64)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		strconv.Itoa(int(p.PoolId)): map[string]interface{}{
			"ShCPUSInSys":          p.ShCPUSInSys,
			"MaxPoolCapacity":      p.MaxPoolCapacity,
			"EntitledPoolCapacity": p.EntitledPoolCapacity,
			"PoolMaxTime":          p.PoolMaxTime,
			"PoolBusyTime":         p.PoolBusyTime,
			"ShCPUTotTime":         p.ShCPUTotTime,
			"ShCPUBusyTime":        p.ShCPUBusyTime,
			"Entitled":             p.Entitled,
		},
	}, nil

}

// LPAR,Logical Partition rcvioc03,PhysicalCPU,virtualCPUs,logicalCPUs,poolCPUs,
// entitled,weight,PoolIdle,usedAllCPU%,usedPoolCPU%,SharedCPU,Capped,EC_User%,
// EC_Sys%,EC_Wait%,EC_Idle%,VP_User%,VP_Sys%,VP_Wait%,VP_Idle%,Folded,Pool_id
type aixLPAR struct {
	PhysicalCPU float64
	VirtualCPUS float64
	LogicalCPUS float64
	PoolCPUS    float64
	Entitled    float64
	Weight      float64
	PoolIdle    float64
	UsedAllCPU  float64
	UsedPoolCPU float64
	SharedCPU   float64
	Capped      float64
	ECUser      float64
	ECSys       float64
	ECWait      float64
	ECIdle      float64
	VPUser      float64
	VPSys       float64
	VPWait      float64
	VPIdle      float64
	Folded      float64
	PoolId      int64
}

//
func (p *aixLPAR) Measurement() string {
	return aixLPARMeasurement
}

//
func (p *aixLPAR) Parse(line string) (data map[string]map[string]interface{}, err error) {
	array := strings.Split(line, fieldSep)

	if len(array) != 23 {
		err = aixFieldLenError
		return
	}

	p.PhysicalCPU, err = strconv.ParseFloat(array[2], 64)
	if err != nil {
		return
	}
	p.VirtualCPUS, err = strconv.ParseFloat(array[3], 64)
	if err != nil {
		return
	}
	p.LogicalCPUS, err = strconv.ParseFloat(array[4], 64)
	if err != nil {
		return
	}
	p.PoolCPUS, err = strconv.ParseFloat(array[5], 64)
	if err != nil {
		return
	}
	p.Entitled, err = strconv.ParseFloat(array[6], 64)
	if err != nil {
		return
	}
	p.Weight, err = strconv.ParseFloat(array[7], 64)
	if err != nil {
		return
	}
	p.PoolIdle, err = strconv.ParseFloat(array[8], 64)
	if err != nil {
		return
	}
	p.UsedAllCPU, err = strconv.ParseFloat(array[9], 64)
	if err != nil {
		return
	}
	p.UsedPoolCPU, err = strconv.ParseFloat(array[10], 64)
	if err != nil {
		return
	}
	p.SharedCPU, err = strconv.ParseFloat(array[11], 64)
	if err != nil {
		return
	}
	p.Capped, err = strconv.ParseFloat(array[12], 64)
	if err != nil {
		return
	}
	p.ECUser, err = strconv.ParseFloat(array[13], 64)
	if err != nil {
		return
	}
	p.ECSys, err = strconv.ParseFloat(array[14], 64)
	if err != nil {
		return
	}
	p.ECWait, err = strconv.ParseFloat(array[15], 64)
	if err != nil {
		return
	}
	p.ECIdle, err = strconv.ParseFloat(array[16], 64)
	if err != nil {
		return
	}
	p.VPUser, err = strconv.ParseFloat(array[17], 64)
	if err != nil {
		return
	}
	p.VPSys, err = strconv.ParseFloat(array[18], 64)
	if err != nil {
		return
	}
	p.VPWait, err = strconv.ParseFloat(array[19], 64)
	if err != nil {
		return
	}
	p.VPIdle, err = strconv.ParseFloat(array[20], 64)
	if err != nil {
		return
	}
	p.Folded, err = strconv.ParseFloat(array[21], 64)
	if err != nil {
		return
	}
	p.PoolId, err = strconv.ParseInt(array[22], 10, 64)
	if err != nil {
		return
	}

	return map[string]map[string]interface{}{
		strconv.Itoa(int(p.PoolId)): map[string]interface{}{
			"PhysicalCPU": p.PhysicalCPU,
			"VirtualCPUS": p.VirtualCPUS,
			"LogicalCPUS": p.LogicalCPUS,
			"PoolCPUS":    p.PoolCPUS,
			"Entitled":    p.Entitled,
			"Weight":      p.Weight,
			"PoolIdle":    p.PoolIdle,
			"UsedAllCPU":  p.UsedAllCPU,
			"UsedPoolCPU": p.UsedPoolCPU,
			"SharedCPU":   p.SharedCPU,
			"capped":      p.Capped,
			"ecuser":      p.ECUser,
			"ecsys":       p.ECSys,
			"ecwait":      p.ECWait,
			"ecidle":      p.ECIdle,
			"vpusr":       p.VPUser,
			"vpsys":       p.VPSys,
			"vpwait":      p.VPWait,
			"vpidle":      p.VPIdle,
			"folded":      p.Folded,
		},
	}, nil

}

// NET,Network I/O rcvioc03,en0-read-KB/s,lo0-read-KB/s,
// en0-write-KB/s,lo0-write-KB/s
type aixNET struct {
}

//
func (p *aixNET) Measurement() string {
	return aixNETMeasurement
}

//
func (p *aixNET) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyNET], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	var (
		step = (klen - 2) / 2
		nd   = make([]string, 0)
	)

	for _, v := range karray[2 : 2+step] {
		array := strings.Split(v, "-")
		if len(array) != 3 {
			err = aixFieldLenError
			return
		}
		nd = append(nd, array[0])
	}

	var (
		read, write float64
	)

	data = make(map[string]map[string]interface{})

	for i, v := range nd {

		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		write, err = strconv.ParseFloat(varray[i+2+step], 10)
		if err != nil {
			return
		}
		i++

		data[v] = map[string]interface{}{
			"read":  read,
			"write": write,
		}
	}

	return data, nil

}

// NETPACKET,Network Packets rcvioc03,en0-reads/s,lo0-reads/s,
// en0-writes/s,lo0-writes/s
type aixNETPACKET struct {
	aixNET
}

//
func (p *aixNETPACKET) Measurement() string {
	return aixNETPACKETMeasurement
}

// NETSIZE,Network Size rcvioc03,en0-readsize,lo0-readsize,
// en0-writesize,lo0-writesize
type aixNETSIZE struct {
	aixNET
}

//
func (p *aixNETSIZE) Measurement() string {
	return aixNETSIZEMeasurement
}

// NETERROR,Network Errors rcvioc03,en0-ierrs,lo0-ierrs,en0-oerrs,lo0-oerrs,
// en0-collisions,lo0-collisions
type aixNETERROR struct {
}

//
func (p *aixNETERROR) Measurement() string {
	return aixNETERRORMeasurement
}

//
//
func (p *aixNETERROR) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyNETERROR], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	var (
		step = (klen - 2) / 3
		nd   = make([]string, 0)
	)

	for _, v := range karray[2 : 2+step] {
		array := strings.Split(v, "-")
		if len(array) != 2 {
			err = aixFieldLenError
			return
		}
		nd = append(nd, array[0])
	}

	var (
		read, write, collisions float64
	)

	data = make(map[string]map[string]interface{})

	for i, v := range nd {
		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		write, err = strconv.ParseFloat(varray[i+2+step], 10)
		if err != nil {
			return
		}
		collisions, err = strconv.ParseFloat(varray[i+2+2*step], 10)
		if err != nil {
			return
		}

		i++

		data[v] = map[string]interface{}{
			"read":       read,
			"write":      write,
			"collisions": collisions,
		}
	}

	return data, nil

}

//
// DISKBUSY,Disk %Busy rcvioc03,hdisk0,cd0,hdisk1
type aixDISKBUSY struct {
}

//
func (p *aixDISKBUSY) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKBUSY) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var busy float64
	for i, _ := range karray[2:] {
		busy, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKBUSY: busy,
		}
	}

	return data, nil

}

// DISKREAD,Disk Read KB/s rcvioc03,hdisk0,cd0,hdisk1
type aixDISKREAD struct {
}

//
func (p *aixDISKREAD) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKREAD) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var read float64
	for i, _ := range karray[2:] {
		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKREAD: read,
		}
	}

	return data, nil

}

// DISKWRITE,Disk Write KB/s rcvioc03,hdisk0,cd0,hdisk1
type aixDISKWRITE struct{}

//
func (p *aixDISKWRITE) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKWRITE) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var write float64
	for i, _ := range karray[2:] {
		write, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKWRITE: write,
		}
	}

	return data, nil

}

// DISKXFER,Disk transfers per second rcvioc03,hdisk0,cd0,hdisk1
type aixDISKXFER struct{}

//
func (p *aixDISKXFER) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKXFER) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var write float64
	for i, _ := range karray[2:] {
		write, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKXFER: write,
		}
	}

	return data, nil

}

// DISKRXFER,Transfers from disk (reads) per second rcvioc03,hdisk0,cd0,hdisk1
type aixDISKRXFER struct{}

//
func (p *aixDISKRXFER) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKRXFER) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var read float64
	for i, _ := range karray[2:] {
		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKRXFER: read,
		}
	}

	return data, nil

}

// DISKBSIZE,Disk Block Size rcvioc03,hdisk0,cd0,hdisk1
type aixDISKBSIZE struct{}

//
func (p *aixDISKBSIZE) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKBSIZE) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var bsize float64
	for i, _ := range karray[2:] {
		bsize, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKBSIZE: bsize,
		}
	}

	return data, nil

}

// DISKRIO,Disk IO Reads per second rcvioc03,hdisk0,cd0,hdisk1
type aixDISKRIO struct{}

//
func (p *aixDISKRIO) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKRIO) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var read float64
	for i, _ := range karray[2:] {
		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKRIO: read,
		}
	}

	return data, nil

}

// DISKWIO,Disk IO Writes per second rcvioc03,hdisk0,cd0,hdisk1
type aixDISKWIO struct{}

//
func (p *aixDISKWIO) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKWIO) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var write float64
	for i, _ := range karray[2:] {
		write, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKWIO: write,
		}
	}

	return data, nil

}

// DISKAVGRIO,Disk IO Average Reads per second rcvioc03,hdisk0,cd0,hdisk1
type aixDISKAVGRIO struct{}

//
func (p *aixDISKAVGRIO) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKAVGRIO) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var read float64
	for i, _ := range karray[2:] {
		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKAVGRIO: read,
		}
	}

	return data, nil

}

// DISKAVGWIO,Disk IO Average Writes per second rcvioc03,hdisk0,cd0,hdisk1
type aixDISKAVGWIO struct{}

//
func (p *aixDISKAVGWIO) Measurement() string {
	return aixDISKMeasurement
}

//
func (p *aixDISKAVGWIO) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyDISKBUSY], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var write float64
	for i, _ := range karray[2:] {
		write, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyDISKAVGWIO: write,
		}
	}

	return data, nil

}

// JFSFILE,JFS Filespace %Used rcvioc03,/,/home,/usr,/var,/tmp,/admin,/opt,/var/adm/ras/livedump
type aixJFSFILE struct {
}

//
func (p *aixJFSFILE) Measurement() string {
	return aixJFSMeasurement
}

//
func (p *aixJFSFILE) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyJFSFILE], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var tmp float64
	for i, _ := range karray[2:] {
		tmp, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyJFSFILE: tmp,
		}
	}

	return data, nil

}

// JFSINODE,JFS Inode %Used rcvioc03,/,/home,/usr,/var,/tmp,/admin,/opt,/var/adm/ras/livedump
type aixJFSINODE struct {
}

//
func (p *aixJFSINODE) Measurement() string {
	return aixJFSMeasurement
}

//
func (p *aixJFSINODE) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyJFSINODE], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	data = make(map[string]map[string]interface{})

	var tmp float64
	for i, _ := range karray[2:] {
		tmp, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		data[karray[i+2]] = map[string]interface{}{
			aixKeyJFSINODE: tmp,
		}
	}

	return data, nil

}

//
type aixIOADAPT struct{}

//
func (p *aixIOADAPT) Measurement() string {
	return aixIOADAPTMeasurement
}

//
//
func (p *aixIOADAPT) Parse(datainfo map[string]string, line string) (data map[string]map[string]interface{}, err error) {
	karray := strings.Split(datainfo[aixKeyIOADAPT], fieldSep)
	varray := strings.Split(line, fieldSep)

	klen := len(karray)
	vlen := len(varray)

	if klen != vlen {
		err = aixFieldLenError
		return
	}

	var (
		step = (klen - 2) / 3
		nd   = make([]string, 0)
	)

	for _, v := range karray[2 : 2+step] {
		array := strings.Split(v, "_")
		if len(array) != 2 {
			err = aixFieldLenError
			return
		}
		nd = append(nd, array[0])
	}

	var (
		read, write, xfertps float64
	)

	data = make(map[string]map[string]interface{})

	for i, v := range nd {
		read, err = strconv.ParseFloat(varray[i+2], 10)
		if err != nil {
			return
		}
		write, err = strconv.ParseFloat(varray[i+2+step], 10)
		if err != nil {
			return
		}
		xfertps, err = strconv.ParseFloat(varray[i+2+2*step], 10)
		if err != nil {
			return
		}

		i++

		data[v] = map[string]interface{}{
			"read":    read,
			"write":   write,
			"xfertps": xfertps,
		}
	}

	return data, nil

}
