package machine

import (
	grpc_server "bigagent/grpcs/server"
	"bigagent/scrape/machine/cpuinfo"
	"bigagent/scrape/machine/diskinfo"
	"bigagent/scrape/machine/kmodule"
	kmodules "bigagent/scrape/machine/kmodule"
	"bigagent/scrape/machine/meminfo"
	"bigagent/scrape/machine/netinfo"
	"bigagent/scrape/machine/processinfo"
	"log"
	"os"
	"time"

	"github.com/super-l/machine-code/machine"
)

type SmpMachine struct {
	Uuid     string             `json:"uuid"`
	Hostname string             `json:"hostname"`
	IPv4     string             `json:"ipv4"`
	Time     time.Time          `json:"time"`
	Cpu      *cpuinfo.SmpCpu    `json:"cpu"`
	Disk     *diskinfo.SmpDisk  `json:"disk"`
	Memory   *meminfo.SmpMemory `json:"memory"`
	Kmodules *kmodule.Kmodules  `json:"kernel_modules"`
	Net      *netinfo.SmpNet    `json:"network"`
	Process  *processinfo.SmpPs `json:"process"`
}

type SmpMachineGrpc struct {
	Uuid     string                                     `json:"uuid"`
	Hostname string                                     `json:"hostname"`
	IPv4     string                                     `json:"ipv4"`
	Time     time.Time                                  `json:"time"`
	Cpu      *cpuinfo.SmpCpu                            `json:"cpu"`
	Disk     map[string]*grpc_server.SmpDisk            `json:"disk"`
	Memory   *meminfo.SmpMemory                         `json:"memory"`
	Kmodules map[string]*grpc_server.Win32_SystemDriver `json:"kernel_modules"`
	Net      map[string]*grpc_server.SmpNetInfo         `json:"network"`
	Process  map[string]*grpc_server.SmPsInfo           `json:"process"`
}

var (
	SmpMa        = NewSmpMachine()
	SmpMaGrpc    = NewSmpMachineGrpc()
	MachineChSmp = make(chan bool, 1)
)

// Machine 存放所有的采集层数据，被懒汉式创建
func NewSmpMachine() *SmpMachine {
	uuid := machine.GetMachineData()
	addr, err := machine.GetLocalIpAddr()
	if err != nil {
		log.Fatal(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	return &SmpMachine{
		Uuid:     uuid.PlatformUUID,
		Hostname: hostname,
		IPv4:     addr,
		Time:     time.Now(),
		Cpu:      cpuinfo.NewSmpCpu(),
		Disk:     diskinfo.NewSmpDisk(),
		Memory:   meminfo.NewSmpMem(),
		Kmodules: kmodules.NewKmodules(),
		Net:      netinfo.NewSmpNet(),
		Process:  processinfo.NewSmpPs(),
	}
}

func NewSmpMachineGrpc() *SmpMachineGrpc {
	uuid := machine.GetMachineData()
	addr, err := machine.GetLocalIpAddr()
	if err != nil {
		log.Fatal(err)
	}
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	return &SmpMachineGrpc{
		Uuid:     uuid.PlatformUUID,
		Hostname: hostname,
		IPv4:     addr,
		Time:     time.Now(),
		Cpu:      cpuinfo.NewSmpCpu(),
		Disk:     *diskinfo.NewSmpDiskGrpc(),
		Memory:   meminfo.NewSmpMem(),
		Kmodules: *kmodules.NewKmodulesGrpc(),
		Net:      *netinfo.NewSmpNetGrpc(),
		Process:  *processinfo.NewSmpPsGrpc(),
	}
}

func NotifySmpMachineAddressChange() {
	select {
	case MachineCh <- true:
	default:
	}
}