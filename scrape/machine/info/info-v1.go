package info

import (
	"github.com/shirou/gopsutil/v4/host"
	"log"
)

type Info struct {
	V1 host.InfoStat `json:"strategy"`
}

func NewInfo() *Info {
	info, err := host.Info()
	if err != nil {
		log.Println(err)
	}
	return &Info{V1: *info}
}

func (info *Info) Platform() string {
	return info.V1.Platform
}

func (info *Info) PlatformFamily() string {
	return info.V1.PlatformFamily
}

func (Info *Info) PlatformVersion() string {
	return Info.V1.PlatformVersion
}

func (info *Info) Hostname() string {
	return info.V1.Hostname
}

func (info *Info) KernelVersion() string {
	return info.V1.KernelVersion
}

func (info *Info) OS() string {
	return info.V1.OS
}

func (info *Info) VirtualizationSystem() string {
	return info.V1.VirtualizationSystem
}

func (info *Info) HostID() string {
	return info.V1.HostID
}
