package model

import (
	"bigagent/scrape/machine"
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/kmodule"
	"bigagent/scrape/machine/meminfo"
	"bigagent/scrape/machine/netinfo"
	"bigagent/scrape/machine/processinfo"
	"encoding/json"
	"log"
	"time"
)

// StandData 暴露原生utils数据
type SmpData struct {
	// Serct    string            `json:"serct"`
	Uuid     string            `json:"uuid"`
	Hostname string            `json:"hostname"`
	IPv4     string            `json:"ipv4"`
	Time     time.Time         `json:"time"`
	Cpu      cpuinfo.SmpCpu    `json:"cpu"`
	Disk     diskinfo.SmpDisk  `json:"disk"`
	Memory   meminfo.SmpMemory `json:"memory"`
	Kmodules kmodule.Kmodules  `json:"kernel_module"`
	Net      netinfo.SmpNet    `json:"net"`
	Process  processinfo.SmpPs `json:"process"`
}

func NewSmpData() *SmpData {
	if machine.SmpMa == nil {
		log.Fatal("machine.SmpMa is nil!")
	}

	// s := global.CONF.System.Serct
	u := machine.SmpMa.Uuid
	h := machine.SmpMa.Hostname
	i := machine.SmpMa.IPv4
	t := machine.SmpMa.Time
	c := machine.SmpMa.Cpu
	d := machine.SmpMa.Disk
	m := machine.SmpMa.Memory
	k := machine.SmpMa.Kmodules
	n := machine.SmpMa.Net
	p := machine.SmpMa.Process

	return &SmpData{
		// Serct:    s,
		Uuid:     u,
		Hostname: h,
		IPv4:     i,
		Time:     t,
		Cpu:      *c,
		Disk:     *d,
		Memory:   *m,
		Kmodules: *k,
		Net:      *n,
		Process:  *p,
	}
}

func NewSmpDataApi() *SmpData {
	// s := global.CONF.System.Serct
	u := machine.SmpMa.Uuid
	h := machine.SmpMa.Hostname
	i := machine.SmpMa.IPv4
	t := machine.SmpMa.Time
	c := cpuinfo.NewSmpCpu()
	d := diskinfo.NewSmpDisk()
	m := meminfo.NewSmpMem()
	k := kmodule.NewKmodules()
	n := netinfo.NewSmpNet()
	p := processinfo.NewSmpPs()
	return &SmpData{
		// Serct:    s,
		Uuid:     u,
		Hostname: h,
		IPv4:     i,
		Time:     t,
		Cpu:      *c,
		Disk:     *d,
		Memory:   *m,
		Kmodules: *k,
		Net:      *n,
		Process:  *p,
	}
}

func (d *SmpData) ToString() string {
	s, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}
	return string(s)
}
