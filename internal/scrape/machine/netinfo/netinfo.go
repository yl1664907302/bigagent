package netinfo

import (
	"github.com/shirou/gopsutil/v4/net"
	"log"
)

type Net struct {
	IA net.InterfaceStatList `json:"IA"`
	IO []net.IOCountersStat  `json:"io"`
}

type Options func(n *Net)

func WithIA(ia net.InterfaceStatList, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(n *Net) {
		n.IA = ia
	}
}

func NewNET(options ...Options) *Net {
	interfaces, err := net.Interfaces()

	counters, err := net.IOCounters(true)
	if err != nil {
		log.Fatal(err)
	}
	n := &Net{
		IA: interfaces,
		IO: counters,
	}
	for _, option := range options {
		option(n)
	}
	return n
}
