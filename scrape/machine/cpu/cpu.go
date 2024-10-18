package cpu

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

// CpuInfo 定义Cpu类型
type CpuInfo struct {
	InfoCpu   []cpu.InfoStat  `json:"cpu_info"`
	TimesStat []cpu.TimesStat `json:"times_cpu"`
}

// NewCpuInfo cpuInfo对象
func NewCpuInfo() *CpuInfo {
	cpuInfo, _ := cpu.Info()
	times, _ := cpu.Times(true)
	return &CpuInfo{
		InfoCpu:   cpuInfo,
		TimesStat: times,
	}
}

// GetCpuInfo cpu信息
func (cpuInfo *CpuInfo) GetCpuInfo() []cpu.InfoStat {
	return cpuInfo.InfoCpu
}

// GetTimesStat cpu 时间信息
func (cpuInfo *CpuInfo) GetTimesStat() []cpu.TimesStat {

	return cpuInfo.TimesStat
}
