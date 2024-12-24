package model

import (
	"bigagent/scrape/machine"
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/meminfo"
	utils "bigagent/util"
	"encoding/json"
	"time"
)

type PatrolData struct {
	Cpu    cpuinfo.SmpCpu    `json:"cpu"`
	Disk   diskinfo.SmpDisk  `json:"disk"`
	Memory meminfo.SmpMemory `json:"memory"`
	Time   time.Time         `json:"time"`
}

func NewPatrolData() *PatrolData {
	if machine.SmpMa == nil {
		utils.DefaultLogger.Error("machine.SmpMa is nil!")
	}

	t := machine.SmpMa.Time
	c := machine.SmpMa.Cpu
	d := machine.SmpMa.Disk
	m := machine.SmpMa.Memory

	return &PatrolData{
		Disk:   *d,
		Memory: *m,
		Cpu:    *c,
		Time:   t,
	}
}

func (d *PatrolData) ToString() string {
	s, err := json.Marshal(d)
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return string(s)
}
