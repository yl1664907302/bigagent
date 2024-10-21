package model

import (
	"encoding/json"
	"github.com/super-l/machine-code/machine"
	"log"
	"machine/cpuinfo"
	"machine/diskinfo"
	"machine/meminfo"
	"machine/proessinfo"
)

var uuid string

type Model struct {
	Uuid    string                 `json:"uuid"`
	Cpu     cpuinfo.CpuInfo        `json:"cpu_info"`
	Mem     meminfo.MemInfo        `json:"mem_info"`
	Disk    diskinfo.DiskInfo      `json:"disk_info"`
	Process proessinfo.ProcessInfo `json:"progress"`
}

// GetUuid 获得系统uuid
func (m *Model) GetUuid() string {
	data := machine.GetMachineData()
	m.Uuid = data.PlatformUUID
	return m.Uuid
}

// 获得系统所有信息已json格式
func (m *Model) String() string {
	s, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(s)
}

// NewModel 系统信息对外接口
func NewModel() *Model {
	return &Model{
		Uuid: uuid,
		Cpu: cpuinfo.CpuInfo{
			CpuInfo: cpuinfo.GetCpuInfo(),
			TimeCpu: cpuinfo.GetTimeCpu(),
		},
		Mem: meminfo.MemInfo{
			VirMem:  meminfo.GetVirMem(),
			SwapMem: meminfo.GetSwapMem(),
		},
		Disk: diskinfo.DiskInfo{
			Usage:      diskinfo.GetUsage(),
			Partition:  diskinfo.GetPartition(),
			IOCounters: diskinfo.GetIOCounters(),
		},
		Process: proessinfo.ProcessInfo{
			Process: proessinfo.GetProcessInfo(),
		},
	}
}
