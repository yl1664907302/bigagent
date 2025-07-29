package machine

import (
	"bigagent/internal/scrape/machine/cpuinfo"
	"bigagent/internal/scrape/machine/diskinfo"
	"bigagent/internal/scrape/machine/info"
	"bigagent/internal/scrape/machine/meminfo"
	"bigagent/internal/scrape/machine/netinfo"
	grpc_server "bigagent/internal/web/grpcs/server"
	"runtime"
	"strings"
	"time"
)

var (
	MachineCh = make(chan bool, 1)
	SmpMa     *SmpMachine
)

type SmpMachine struct {
	Uuid       string                          `json:"uuid"`
	Platform   string                          `json:"platform"`
	Os         string                          `json:"os"`
	Kernel     string                          `json:"kernel"`
	Hostname   string                          `json:"hostname"`
	IPv4       string                          `json:"ipv4"`
	Arch       string                          `json:"arch"`
	Machine    string                          `json:"machine"`
	Disk_use   map[string]string               `json:"disk_use"`
	Memory_use string                          `json:"memory_use"`
	Cpu_use    string                          `json:"cpu_use"`
	Time       time.Time                       `json:"time"`
	Cpu        *cpuinfo.SmpCpu                 `json:"cpu"`
	Disk       *diskinfo.SmpDisk               `json:"disk"`
	Disk_g     map[string]*grpc_server.SmpDisk `json:"disk_g"`
	Memory     *meminfo.SmpMemory              `json:"memory"`
	//Kmodules   *kmodules.Kmodules `json:"kernel_modules"`
	Net      *netinfo.SmpNet                    `json:"network"`
	Net_g    map[string]*grpc_server.SmpNetInfo `json:"network_g"`
	BiosData string                             `json:"bios_data"`
	//Process    *processinfo.SmpPs `json:"process"`
}

type SmpMachineGrpc struct {
	Uuid       string                          `json:"uuid"`
	Platform   string                          `json:"platform"`
	Os         string                          `json:"os"`
	Kernel     string                          `json:"kernel"`
	Hostname   string                          `json:"hostname"`
	IPv4       string                          `json:"ipv4"`
	Arch       string                          `json:"arch"`
	Machine    string                          `json:"machine"`
	Disk_use   map[string]string               `json:"disk_use"`
	Memory_use string                          `json:"memory_use"`
	Cpu_use    string                          `json:"cpu_use"`
	Time       time.Time                       `json:"time"`
	Cpu        *cpuinfo.SmpCpu                 `json:"cpu"`
	Disk       map[string]*grpc_server.SmpDisk `json:"disk"`
	Memory     *meminfo.SmpMemory              `json:"memory"`
	//Kmodules   map[string]*grpc_server.Win32_SystemDriver `json:"kernel_modules"`
	Net map[string]*grpc_server.SmpNetInfo `json:"network"`
	//Process    map[string]*grpc_server.SmPsInfo           `json:"process"`
}

// NewSmpMachine 存放所有的采集层数据，被懒汉式创建
func NewSmpMachine() *SmpMachine {
	info := info.NewInfo()

	disk_use := make(map[string]string)

	disks := diskinfo.NewSmpDisk()
	// 创建一个新的 SmpDisk 来存储过滤后的结果
	filteredDisks := &diskinfo.SmpDisk{}
	for _, d := range *disks {
		// Windows 系统下保留所有设备
		if runtime.GOOS == "windows" {
			(*filteredDisks)[d.Path] = d
			disk_use[d.Device] = d.UsedPercent
			continue
		}

		// Linux 系统下过滤掉特定挂载点
		if runtime.GOOS == "linux" {
			if strings.Contains(d.Path, "/proc") ||
				strings.Contains(d.Path, "/docker") ||
				strings.Contains(d.Path, "/kubelet") ||
				strings.Contains(d.Path, "/run") ||
				strings.Contains(d.Path, "/tmp") ||
				strings.Contains(d.Path, "/var") ||
				strings.Contains(d.Path, "/sys") {
				continue
			}
			// 只添加需要保留的设备
			(*filteredDisks)[d.Path] = d
			disk_use[d.Device] = d.UsedPercent
		}
	}

	// 使用过滤后的结果替换原始数据
	disks = filteredDisks
	cpu := cpuinfo.NewSmpCpu()
	memory := meminfo.NewSmpMem()
	disk_g := *diskinfo.NewSmpDiskGrpc()
	net := netinfo.NewSmpNet()
	net_g := *netinfo.NewSmpNetGrpc()
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
		Memory_use: memory.Vmem.UsedPercent,
		Cpu_use:    cpu.Usage,
		Time:       time.Now(),
		Cpu:        cpu,
		Disk:       disks,
		Disk_g:     disk_g,
		Memory:     memory,
		//Kmodules:   kmodules.NewKmodules(),
		Net:      net,
		Net_g:    net_g,
		BiosData: "",
		//Process: processinfo.NewSmpPs(),
	}
}

func NotifySmpMachineAddressChange() {
	select {
	case MachineCh <- true:
	default:
	}
}
