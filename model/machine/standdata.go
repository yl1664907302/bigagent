package model

import (
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/disk"
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/memory"
	"bigagent/scrape/machine/process"
	"encoding/json"
	"strings"
)

// StandData 暴露原生utils数据
type StandData struct {
	Uuid    *info.Info              `json:"uuid"`
	Cpu     *cpuinfo.CpuInfo        `json:"cpu_info"`
	Mem     *meminfo.MemInfo        `json:"mem_info"`
	Disk    *diskinfo.DiskInfo      `json:"disk_info"`
	Process *proessinfo.ProcessInfo `json:"progress"`
}

func NewStandData() *StandData {
	c := cpuinfo.NewCpuInfo()
	m := meminfo.NewMemInfo()
	d := diskinfo.NewDiskInfo()
	p := proessinfo.NewProcessInfo()
	i := info.NewInfo()
	return &StandData{
		Uuid:    i,
		Cpu:     c,
		Mem:     m,
		Disk:    d,
		Process: p,
	}
}

func NewStandDataApi() *StandData {
	c := cpuinfo.NewCpuInfo()
	m := meminfo.NewMemInfo()
	d := diskinfo.NewDiskInfo()
	p := proessinfo.NewProcessInfo()
	i := info.NewInfo()
	return &StandData{
		Uuid:    i,
		Cpu:     c,
		Mem:     m,
		Disk:    d,
		Process: p,
	}
}

func (d *StandData) ToString() (string, error) {
	var v []string
	uuid := d.Uuid.GetUuid()
	cpus, _ := d.Cpu.ToString()
	mems, _ := d.Mem.ToString()
	disks, _ := d.Disk.ToString()
	ps, _ := d.Process.ToString()
	v = append(v, uuid)
	v = append(v, cpus)
	v = append(v, mems)
	v = append(v, disks)
	v = append(v, ps)
	standData := strings.Join(v, "")
	marshal, err := json.Marshal(standData)

	return string(marshal), err
}
