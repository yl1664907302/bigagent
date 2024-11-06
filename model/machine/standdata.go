package model

import (
	"bigagent/scrape/machine"
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/meminfo"
	"bigagent/scrape/machine/netinfo"
	"bigagent/scrape/machine/processinfo"
	"encoding/json"
	"log"
)

// StandData 暴露原生utils数据
type StandData struct {
	Info    info.Info           `json:"info"`
	Cpu     cpuinfo.Cpus        `json:"cpu"`
	Disk    diskinfo.DISK       `json:"disk"`
	Memory  meminfo.Memory      `json:"memory"`
	Net     netinfo.Net         `json:"net"`
	Process processinfo.PROCESS `json:"process"`
}

func NewStandData() *StandData {
	i := machine.Ma.Info
	c := machine.Ma.Cpu
	d := machine.Ma.Disk
	m := machine.Ma.Memory
	n := machine.Ma.Net
	p := machine.Ma.Process

	return &StandData{
		Info:    *i,
		Cpu:     *c,
		Disk:    *d,
		Memory:  *m,
		Net:     *n,
		Process: *p,
	}
}

func NewStandDataApi() *StandData {
	i := info.NewInfo()
	c := cpuinfo.NewCPU()
	d := diskinfo.NewDISK()
	m := meminfo.NewMemory()
	n := netinfo.NewNET()
	p := processinfo.NewProcess()
	return &StandData{
		Info:    *i,
		Cpu:     *c,
		Disk:    *d,
		Memory:  *m,
		Net:     *n,
		Process: *p,
	}
}

func (d *StandData) ToString() string {
	s, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}
	return string(s)
}
