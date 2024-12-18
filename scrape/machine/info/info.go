package info

import (
	machine "github.com/super-l/machine-code"
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
	addr, err := machine.GetIpAddr()
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
	uuid := machine.Machine

	return uuid.UUID
}

// 对外接口
func NewInfo() *Info {
	uuid := machine.Machine
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	addr, err := machine.GetIpAddr()
	if err != nil {
		log.Fatal(err)
	}
	return &Info{
		Uuid:     uuid.UUID,
		Hostname: hostname,
		IPv4:     addr,
		Time:     time.Now(),
	}
}
