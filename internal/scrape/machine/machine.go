package machine

import (
	"bigagent/internal/scrape/machine/cpuinfo"
	"time"
)

// Machine 存放所有的采集层数据，被懒汉式创建
type Machine struct {
	Uuid       string            `json:"uuid"`
	Platform   string            `json:"platform"`
	Os         string            `json:"os"`
	Kernel     string            `json:"kernel"`
	Hostname   string            `json:"hostname"`
	IPv4       string            `json:"ipv4"`
	Arch       string            `json:"arch"`
	Machine    string            `json:"machine"`
	Disk_use   map[string]string `json:"disk_use"`
	Memory_use string            `json:"memory_use"`
	Cpu_use    string            `json:"cpu_use"`
	Time       time.Time         `json:"time"`
	Cpu        *cpuinfo.SmpCpu   `json:"cpu"`
}

func NewMachine() *Machine {
	//info := info.NewInfo()
	//
	//disk_use := make(map[string]string)
	//
	//disks := diskinfo.NewSmpDisk()
	//
	//for _, v := range *disks {
	//	disk_use[v.Device] = v.UsedPercent
	//}
	return &Machine{}
}

func NotifyMachineAddressChange() {
	select {
	case MachineCh <- true:
	default:
	}
}

var (
	Ma        = NewMachine()
	MachineCh = make(chan bool, 1)
)
