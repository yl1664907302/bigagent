package model

import (
	"bigagent/config/global"
	"bigagent/scrape/machine"
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/meminfo"
	"encoding/json"
	"log"
	"time"
)

// StandData 暴露原生utils数据
type StandData struct {
	Serct    string         `json:"serct"`
	Uuid     string         `json:"uuid"`
	Hostname string         `json:"hostname"`
	IPv4     string         `json:"ipv4"`
	Time     time.Time      `json:"time"`
	Info     info.Info      `json:"info"`
	Cpu      cpuinfo.Cpus   `json:"cpu"`
	Disk     diskinfo.DISK  `json:"disk"`
	Memory   meminfo.Memory `json:"memory"`
	//Net      netinfo.Net    `json:"net"`
	//Process  processinfo.PROCESS `json:"process"`
}

func NewStandData() *StandData {
	s := global.V.GetString("system.serct")
	u := machine.Ma.Uuid
	h := machine.Ma.Hostname
	i := machine.Ma.IPv4
	t := machine.Ma.Time
	c := machine.Ma.Cpu
	d := machine.Ma.Disk
	m := machine.Ma.Memory
	//n := machine.Ma.Net
	//p := machine.Ma.Process

	return &StandData{
		Serct:    s,
		Uuid:     u,
		Hostname: h,
		IPv4:     i,
		Time:     t,
		Cpu:      *c,
		Disk:     *d,
		Memory:   *m,
		//Net:     *n,
		//Process:  *p,
	}
}

func NewStandDataApi() *StandData {
	s := global.V.GetString("system.serct")
	u := info.NewInfo().Uuid
	h := info.NewInfo().Hostname
	i := info.NewInfo().IPv4
	t := info.NewInfo().Time
	c := cpuinfo.NewCPU()
	d := diskinfo.NewDISK()
	m := meminfo.NewMemory()
	//n := netinfo.NewNET()
	//p := processinfo.NewProcess()
	return &StandData{
		Serct:    s,
		Uuid:     u,
		Hostname: h,
		IPv4:     i,
		Time:     t,
		Cpu:      *c,
		Disk:     *d,
		Memory:   *m,
		//Net:      *n,
		//Process:  *p,
	}
}

func (d *StandData) ToString() string {
	s, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}
	return string(s)
}
