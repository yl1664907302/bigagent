package model

import (
	"encoding/json"
	"newmachine/cpu"
	"newmachine/disk"
	"newmachine/memory"
	"newmachine/psinfo"
	"sync"
)

// StandData Data 暴露原生utils数据
type StandData struct {
	Cpu     []*cpu.CpuInfo   `json:"cpu"`
	Memory  []*memory.Memory `json:"memory"`
	Dsk     []*disk.DiskInfo `json:"disk"`
	PsInfo  []*psinfo.PsInfo `json:"progress"`
	maplock sync.RWMutex
}

// model对外接口获取全部数据
func NewStandData() *StandData {
	cpuInfo := []*cpu.CpuInfo{}
	cpuInfo = append(cpuInfo, cpu.NewCpuInfo())
	memoryInfo := []*memory.Memory{}
	memoryInfo = append(memoryInfo, memory.NewMemory())
	diskInfo := []*disk.DiskInfo{}
	diskInfo = append(diskInfo, disk.NewDiskInfo())
	psInfo := []*psinfo.PsInfo{}
	psInfo = append(psInfo, psinfo.NewProgressInfo())
	return &StandData{
		Cpu:     cpuInfo,
		Memory:  memoryInfo,
		Dsk:     diskInfo,
		PsInfo:  psInfo,
		maplock: sync.RWMutex{},
	}
}

// 将所有数据解析成json格式
func (standData *StandData) String() string {
	s, _ := json.MarshalIndent(standData, "", "  ")
	return string(s)
}
