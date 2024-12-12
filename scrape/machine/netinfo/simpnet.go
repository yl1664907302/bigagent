package netinfo

import (
	grpc_server "bigagent/grpcs/server"
	"encoding/json"
	"log"

	"github.com/shirou/gopsutil/v4/net"
)

type smpInfo struct {
	Name string `json:"name"`
	Mtu  int    `json:"mtu"`
	Mac  string `json:"mac"`
	IP   string `json:"ip"`
}

type SmpNet map[string]smpInfo

// var wg sync.RWMutex

func NewSmpNet() *SmpNet {
	// smpnet := make(map[string]map[string]any)
	smpnet := make(SmpNet)
	n, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range n {
		ip := ""
		if len(i.Addrs) > 0 {
			ip = i.Addrs[0].Addr
		}

		smpinfo := smpInfo{
			Name: i.Name,
			Mtu:  i.MTU,
			Mac:  i.HardwareAddr,
			IP:   ip,
		}
		smpnet[i.Name] = smpinfo
	}

	return &smpnet
}

func NewSmpNetGrpc() *map[string]*grpc_server.SmpNetInfo {
	smpnet := make(map[string]*grpc_server.SmpNetInfo)
	n, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range n {
		ip := ""
		if len(i.Addrs) > 0 {
			ip = i.Addrs[0].Addr
		}

		smpinfo := &grpc_server.SmpNetInfo{
			Name: i.Name,
			Mtu:  int64(i.MTU),
			Mac:  i.HardwareAddr,
			Ip:   ip,
		}
		smpnet[i.Name] = smpinfo
	}

	return &smpnet
}

func (n *SmpNet) ToString() string {
	b, err := json.MarshalIndent(n, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}