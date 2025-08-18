package model

import (
	"bigagent/internal/scrape/machine"
	"bigagent/internal/scrape/machine/cpuinfo"
	"bigagent/internal/scrape/machine/diskinfo"
	"bigagent/internal/scrape/machine/meminfo"
	"bigagent/internal/utils"
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
	smpMachine := machine.NewSmpMachine()
	t := smpMachine.Time
	c := smpMachine.Cpu
	d := smpMachine.Disk
	m := smpMachine.Memory

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
