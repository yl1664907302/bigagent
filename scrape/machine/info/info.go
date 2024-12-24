package info

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/elastic/go-sysinfo"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/super-l/machine-code/machine"
)

// 定义系统类型
type Info struct {
	Uuid     string    `json:"uuid"`
	Platform string    `json:"platform"`
	Os       string    `json:"os"`
	Kernel   string    `json:"kernel"`
	Arch     string    `json:"arch"`
	Hostname string    `json:"hostname"`
	IPv4     string    `json:"ipv4"`
	Virtual  string    `json:"virtual_machine" yaml:"virtual_machine"`
	Time     time.Time `json:"time"`
}

func (i *Info) GetArch() string {
	return i.Arch
}

func (i *Info) GetPlatform() string {
	return runtime.GOOS
}

// GetIPv4 获取IP地址
func (i *Info) GetIPv4() string {
	addr, err := machine.GetLocalIpAddr()
	if err != nil {
		log.Fatal(err)
	}
	return addr
}

// GetHostname 获取hostname
func (i *Info) GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	return hostname
}

// GetUuid 获得系统uuid
func (i *Info) GetUuid() string {
	uuid := machine.GetMachineData()

	return uuid.PlatformUUID
}

func (i *Info) GetVirtual() string {
	infos, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}

	if infos.VirtualizationSystem == "" {
		return "Physical"
	}

	return infos.VirtualizationSystem
}

// 对外接口
func NewInfo() *Info {
	uuid := machine.GetMachineData()
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	addr, err := machine.GetLocalIpAddr()
	if err != nil {
		log.Fatal(err)
	}
	infos, err := host.Info()
	if err != nil {
		log.Fatal(err)
	}

	if infos.VirtualizationSystem == "" {
		infos.VirtualizationSystem = "physical"
	}

	// 获取当前主机信息
	host, err := sysinfo.Host()
	if err != nil {
		log.Fatal(err)
	}

	os := host.Info()

	os_version := os.OS.Name + "" + os.OS.Version
	return &Info{
		Uuid:     uuid.PlatformUUID,
		Platform: runtime.GOOS,
		Os:       os_version,
		Kernel:   os.KernelVersion,
		Arch:     host.Info().Architecture,
		Hostname: hostname,
		IPv4:     addr,
		Virtual:  infos.VirtualizationSystem,
		Time:     time.Now(),
	}
}
