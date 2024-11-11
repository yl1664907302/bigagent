package info

import (
	"github.com/super-l/machine-code/machine"
	"log"
	"os"
	"time"
)

// 定义系统类型
type Info struct {
	Uuid     string    `json:"uuid"`
	Hostname string    `json:"hostname"`
	IPv4     string    `json:"ipv4"`
	Time     time.Time `json:"time"`
}

// GetIPv4 获取IP地址
func (i Info) GetIPv4() string {
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
	return &Info{
		Uuid:     uuid.PlatformUUID,
		Hostname: hostname,
		IPv4:     addr,
		Time:     time.Now(),
	}
}
