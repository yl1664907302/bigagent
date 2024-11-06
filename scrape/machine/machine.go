package machine

import (
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/meminfo"
	"bigagent/scrape/machine/netinfo"
	"bigagent/scrape/machine/processinfo"
)

// Machine 存放所有的采集层数据，被懒汉式创建
type Machine struct {
	Info    *info.Info           `json:"info"`
	Cpu     *cpuinfo.Cpus        `json:"cpu"`
	Disk    *diskinfo.DISK       `json:"disk"`
	Memory  *meminfo.Memory      `json:"memory"`
	Net     *netinfo.Net         `json:"network"`
	Process *processinfo.PROCESS `json:"process"`
}

func NewMachine() *Machine {
	return &Machine{
		Info:    info.NewInfo(),
		Cpu:     cpuinfo.NewCPU(),
		Disk:    diskinfo.NewDISK(),
		Memory:  meminfo.NewMemory(),
		Net:     netinfo.NewNET(),
		Process: processinfo.NewProcess(),
	}
}

func NotifyMachineAddressChange() {
	// 非阻塞写入 MachineCh 通道，通知监听者地址变化
	select {
	case MachineCh <- struct{}{}:
		// 通知成功
	default:
		// 通道已满，跳过通知（避免阻塞）
	}
}

var (
	Ma        *Machine
	MachineCh = make(chan struct{}, 1)
)
