package model

import (
	"bigagent/scrape/machine"
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
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

	// 使用空指针检查，避免 nil 解引用
	var cValue cpuinfo.SmpCpu
	if machine.SmpMa.Cpu != nil {
		cValue = *machine.SmpMa.Cpu
	}

	var dValue diskinfo.SmpDisk
	if machine.SmpMa.Disk != nil {
		dValue = *machine.SmpMa.Disk
	}

	var mValue meminfo.SmpMemory
	if machine.SmpMa.Memory != nil {
		mValue = *machine.SmpMa.Memory
	}

	var nValue netinfo.SmpNet
	if machine.SmpMa.Net != nil {
		nValue = *machine.SmpMa.Net
	}

	var pValue processinfo.SmpPs
	if machine.SmpMa.Process != nil {
		pValue = *machine.SmpMa.Process
	}
	return &SmpData{
		// Serct:    s,
		Uuid:     u,
		Hostname: h,
		IPv4:     i,
		Time:     t,
		Cpu:      cValue,
		Disk:     dValue,
		Memory:   mValue,
		Net:      nValue,
		Process:  pValue,
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
