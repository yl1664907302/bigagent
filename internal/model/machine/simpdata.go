package model

import (
	"bigagent/internal/scrape/machine"
	"bigagent/internal/scrape/machine/cpuinfo"
	"bigagent/internal/scrape/machine/diskinfo"
	"bigagent/internal/scrape/machine/meminfo"
	"bigagent/internal/scrape/machine/netinfo"
	"bigagent/internal/scrape/machine/processinfo"
	"bigagent/internal/util"
	"bigagent/internal/web/grpcs/server"
	"encoding/json"
	"time"
)

// StandData 暴露原生utils数据
type SmpData struct {
	Serct      string            `json:"serct"`
	Uuid       string            `json:"uuid"`
	Os         string            `json:"os"`
	Kernel     string            `json:"kernel"`
	Platform   string            `json:"system_platform"`
	Hostname   string            `json:"hostname"`
	IPv4       string            `json:"ipv4"`
	Arch       string            `json:"arch"`
	Virtual    string            `json:"virtual_machine"`
	Disk_use   map[string]string `json:"disk_use"`
	Memory_use string            `json:"memory_use"`
	Cpu_use    string            `json:"cpu_use"`
	Time       time.Time         `json:"time"`
	Cpu        cpuinfo.SmpCpu    `json:"cpu"`
	Disk       diskinfo.SmpDisk  `json:"disk"`
	Memory     meminfo.SmpMemory `json:"memory"`
	//Kmodules   kmodule.Kmodules  `json:"kernel_module"`
	Net     netinfo.SmpNet    `json:"net"`
	Process processinfo.SmpPs `json:"process"`
}

type SmpDataGrpc struct {
	Uuid       string                          `json:"uuid"`
	Os         string                          `json:"os"`
	Kernel     string                          `json:"kernel"`
	Platform   string                          `json:"system_platform"`
	Hostname   string                          `json:"hostname"`
	IPv4       string                          `json:"ipv4"`
	Arch       string                          `json:"arch"`
	Virtual    string                          `json:"virtual_machine"`
	Disk_use   map[string]string               `json:"disk_use"`
	Memory_use string                          `json:"memory_use"`
	Cpu_use    string                          `json:"cpu_use"`
	Time       time.Time                       `json:"time"`
	Cpu        cpuinfo.SmpCpu                  `json:"cpu"`
	Disk       map[string]*grpc_server.SmpDisk `json:"disk"`
	Memory     meminfo.SmpMemory               `json:"memory"`
	//Kmodules   map[string]*grpc_server.Win32_SystemDriver `json:"kernel_module"`
	Net     map[string]*grpc_server.SmpNetInfo `json:"net"`
	Process map[string]*grpc_server.SmPsInfo   `json:"process"`
}

func NewSmpData() *SmpData {
	if machine.SmpMa == nil {
		utils.DefaultLogger.Error("machine.SmpMa is nil!")
	}

	// s := global.CONF.System.Serct
	u := machine.SmpMa.Uuid
	z := machine.SmpMa.Platform
	o := machine.SmpMa.Os
	ker := machine.SmpMa.Kernel
	h := machine.SmpMa.Hostname
	i := machine.SmpMa.IPv4
	arh := machine.SmpMa.Arch
	v := machine.SmpMa.Machine
	du := machine.SmpMa.Disk_use
	mu := machine.SmpMa.Memory.Vmem.UsedPercent
	cu := machine.SmpMa.Cpu.Usage
	t := machine.SmpMa.Time
	c := machine.SmpMa.Cpu
	d := machine.SmpMa.Disk
	m := machine.SmpMa.Memory
	//k := machine.SmpMa.Kmodules
	n := machine.SmpMa.Net
	p := machine.SmpMa.Process

	return &SmpData{
		// Serct:      s,
		Uuid:       u,
		Platform:   z,
		Os:         o,
		Kernel:     ker,
		Hostname:   h,
		IPv4:       i,
		Arch:       arh,
		Virtual:    v,
		Disk_use:   du,
		Memory_use: mu,
		Cpu_use:    cu,
		Time:       t,
		Cpu:        *c,
		Disk:       *d,
		Memory:     *m,
		//Kmodules:   *k,
		Net:     *n,
		Process: *p,
	}
}

func NewSmpDataGrpc() *SmpDataGrpc {
	if machine.SmpMaGrpc == nil {
		utils.DefaultLogger.Error("machine.SmpMa is nil!")
	}

	// s :=glc.CONF.System.Serct
	u := machine.SmpMaGrpc.Uuid
	z := machine.SmpMaGrpc.Platform
	o := machine.SmpMaGrpc.Os
	ker := machine.SmpMaGrpc.Kernel
	h := machine.SmpMaGrpc.Hostname
	i := machine.SmpMaGrpc.IPv4
	arh := machine.SmpMaGrpc.Arch
	v := machine.SmpMaGrpc.Machine
	du := machine.SmpMaGrpc.Disk_use
	mu := machine.SmpMaGrpc.Memory_use
	cu := machine.SmpMaGrpc.Cpu_use
	t := machine.SmpMaGrpc.Time
	c := machine.SmpMaGrpc.Cpu
	d := machine.SmpMaGrpc.Disk
	m := machine.SmpMaGrpc.Memory
	//k := machine.SmpMaGrpc.Kmodules
	n := machine.SmpMaGrpc.Net
	p := machine.SmpMaGrpc.Process

	return &SmpDataGrpc{
		// Serct:    s,
		Uuid:       u,
		Platform:   z,
		Os:         o,
		Kernel:     ker,
		Hostname:   h,
		IPv4:       i,
		Arch:       arh,
		Virtual:    v,
		Disk_use:   du,
		Memory_use: mu,
		Cpu_use:    cu,
		Time:       t,
		Cpu:        *c,
		Disk:       d,
		Memory:     *m,
		//Kmodules:   k,
		Net:     n,
		Process: p,
	}
}

func NewSmpDataApi() *SmpData {

	disk_use := make(map[string]string)
	smpma := machine.NewSmpMachine()

	for _, v := range *smpma.Disk {
		disk_use[v.Device] = v.UsedPercent
	}

	// s :=config.CONF.System.Serct
	u := machine.SmpMa.Uuid
	z := machine.SmpMa.Platform
	o := machine.SmpMa.Os
	ker := machine.SmpMa.Kernel
	h := machine.SmpMa.Hostname
	i := machine.SmpMa.IPv4
	arh := machine.SmpMa.Arch
	v := machine.SmpMa.Machine
	du := machine.SmpMa.Disk_use
	mu := machine.SmpMa.Memory.Vmem.UsedPercent
	cu := machine.SmpMa.Cpu.Usage
	t := machine.SmpMa.Time
	c := cpuinfo.NewSmpCpu()
	d := diskinfo.NewSmpDisk()
	m := meminfo.NewSmpMem()
	//k := kmodule.NewKmodules()
	n := netinfo.NewSmpNet()
	p := processinfo.NewSmpPs()
	return &SmpData{
		// Serct:      s,
		Uuid:       u,
		Platform:   z,
		Os:         o,
		Kernel:     ker,
		Hostname:   h,
		IPv4:       i,
		Arch:       arh,
		Virtual:    v,
		Disk_use:   du,
		Memory_use: mu,
		Cpu_use:    cu,
		Time:       t,
		Cpu:        *c,
		Disk:       *d,
		Memory:     *m,
		//Kmodules:   *k,
		Net:     *n,
		Process: *p,
	}
}

func (d *SmpData) ToString() string {
	s, err := json.Marshal(d)
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return string(s)
}
