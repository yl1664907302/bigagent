package machine

import (
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/meminfo"
	"bigagent/scrape/machine/netinfo"
	"bigagent/scrape/machine/processinfo"
	"time"
)

// Machine 存放所有的采集层数据，被懒汉式创建
type Machine struct {
	Uuid     string               `json:"uuid"`
	Hostname string               `json:"hostname"`
	IPv4     string               `json:"ipv4"`
	Time     time.Time            `json:"time"`
	Cpu      *cpuinfo.Cpus        `json:"cpu"`
	Disk     *diskinfo.DISK       `json:"disk"`
	Memory   *meminfo.Memory      `json:"memory"`
	Net      *netinfo.Net         `json:"network"`
	Process  *processinfo.PROCESS `json:"process"`
}

func NewMachine() *Machine {
	return &Machine{
		Uuid:     info.NewInfo().Uuid,
		Hostname: info.NewInfo().Hostname,
		IPv4:     info.NewInfo().IPv4,
		Time:     time.Now(),
		Cpu:      cpuinfo.NewCPU(),
		Disk:     diskinfo.NewDISK(),
		Memory:   meminfo.NewMemory(),
		Net:      netinfo.NewNET(),
		Process:  processinfo.NewProcess(),
	}
}

func NotifyMachineAddressChange() {
	select {
	case MachineCh <- true:
	default:
	}
}

var (
	Ma        *Machine
	MachineCh = make(chan bool, 1)
)
