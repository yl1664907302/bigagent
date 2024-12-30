package machine

import (
	"bigagent/internal/scrape/machine/cpuinfo"
	"bigagent/internal/scrape/machine/diskinfo"
	"bigagent/internal/scrape/machine/info"
	kmodules "bigagent/internal/scrape/machine/kmodule"
	"bigagent/internal/scrape/machine/meminfo"
	"bigagent/internal/scrape/machine/netinfo"
	"bigagent/internal/scrape/machine/processinfo"
	"bigagent/internal/web/grpcs/server"
	"time"
)

type SmpMachine struct {
	Uuid       string             `json:"uuid"`
	Platform   string             `json:"platform"`
	Os         string             `json:"os"`
	Kernel     string             `json:"kernel"`
	Hostname   string             `json:"hostname"`
	IPv4       string             `json:"ipv4"`
	Arch       string             `json:"arch"`
	Machine    string             `json:"machine"`
	Disk_use   map[string]string  `json:"disk_use"`
	Memory_use string             `json:"memory_use"`
	Cpu_use    string             `json:"cpu_use"`
	Time       time.Time          `json:"time"`
	Cpu        *cpuinfo.SmpCpu    `json:"cpu"`
	Disk       *diskinfo.SmpDisk  `json:"disk"`
	Memory     *meminfo.SmpMemory `json:"memory"`
	Kmodules   *kmodules.Kmodules `json:"kernel_modules"`
	Net        *netinfo.SmpNet    `json:"network"`
	Process    *processinfo.SmpPs `json:"process"`
}

type SmpMachineGrpc struct {
	Uuid       string                                     `json:"uuid"`
	Platform   string                                     `json:"platform"`
	Os         string                                     `json:"os"`
	Kernel     string                                     `json:"kernel"`
	Hostname   string                                     `json:"hostname"`
	IPv4       string                                     `json:"ipv4"`
	Arch       string                                     `json:"arch"`
	Machine    string                                     `json:"machine"`
	Disk_use   map[string]string                          `json:"disk_use"`
	Memory_use string                                     `json:"memory_use"`
	Cpu_use    string                                     `json:"cpu_use"`
	Time       time.Time                                  `json:"time"`
	Cpu        *cpuinfo.SmpCpu                            `json:"cpu"`
	Disk       map[string]*grpc_server.SmpDisk            `json:"disk"`
	Memory     *meminfo.SmpMemory                         `json:"memory"`
	Kmodules   map[string]*grpc_server.Win32_SystemDriver `json:"kernel_modules"`
	Net        map[string]*grpc_server.SmpNetInfo         `json:"network"`
	Process    map[string]*grpc_server.SmPsInfo           `json:"process"`
}

var (
	SmpMa        = NewSmpMachine()
	SmpMaGrpc    = NewSmpMachineGrpc()
	MachineChSmp = make(chan bool, 1)
)

// Machine 存放所有的采集层数据，被懒汉式创建
func NewSmpMachine() *SmpMachine {
	info := info.NewInfo()

	disk_use := make(map[string]string)

	disks := diskinfo.NewSmpDisk()

	for _, v := range *disks {
		disk_use[v.Device] = v.UsedPercent
	}
	return &SmpMachine{
		Uuid:       info.Uuid,
		Os:         info.Os,
		Kernel:     info.Kernel,
		Platform:   info.Platform,
		Hostname:   info.Hostname,
		IPv4:       info.IPv4,
		Arch:       info.Arch,
		Machine:    info.Virtual,
		Disk_use:   disk_use,
		Memory_use: meminfo.NewSmpMem().Vmem.UsedPercent,
		Cpu_use:    cpuinfo.NewSmpCpu().Usage,
		Time:       time.Now(),
		Cpu:        cpuinfo.NewSmpCpu(),
		Disk:       disks,
		Memory:     meminfo.NewSmpMem(),
		Kmodules:   kmodules.NewKmodules(),
		Net:        netinfo.NewSmpNet(),
		Process:    processinfo.NewSmpPs(),
	}
}

func NewSmpMachineGrpc() *SmpMachineGrpc {
	info := info.NewInfo()

	disk_use := make(map[string]string)

	disks := diskinfo.NewSmpDisk()

	for _, v := range *disks {
		disk_use[v.Device] = v.UsedPercent
	}
	return &SmpMachineGrpc{
		Uuid:       info.Uuid,
		Os:         info.Os,
		Kernel:     info.Kernel,
		Platform:   info.Platform,
		Hostname:   info.Hostname,
		IPv4:       info.IPv4,
		Arch:       info.Arch,
		Machine:    info.Virtual,
		Disk_use:   disk_use,
		Memory_use: meminfo.NewSmpMem().Vmem.UsedPercent,
		Cpu_use:    cpuinfo.NewSmpCpu().Usage,
		Time:       time.Now(),
		Cpu:        cpuinfo.NewSmpCpu(),
		Disk:       *diskinfo.NewSmpDiskGrpc(),
		Memory:     meminfo.NewSmpMem(),
		// Kmodules:   kmodules.NewKmodules(),
		Net:     *netinfo.NewSmpNetGrpc(),
		Process: *processinfo.NewSmpPsGrpc(),
	}
}

func NotifySmpMachineAddressChange() {
	select {
	case MachineChSmp <- true:
	default:
	}
}
