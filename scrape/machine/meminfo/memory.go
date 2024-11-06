package meminfo

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v4/mem"
	"log"
)

type Memory struct {
	V  *mem.VirtualMemoryStat `json:"virtual_memory"`
	S  *mem.SwapMemoryStat    `json:"swap_memory"`
	SD []*mem.SwapDevice      `json:"swap_devices"`
}

type Opthions func(m *Memory)

func (m *Memory) ToString() string {
	indent, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		return err.Error()
	}
	return string(indent)
}

func WithSD(sd []*mem.SwapDevice) Opthions {
	sDevices, err := mem.SwapDevices()
	if err != nil {
		log.Fatal(err)
	}

	return func(m *Memory) {
		m.SD = sDevices
	}
}

func WithS(swap *mem.VirtualMemoryStat) Opthions {
	swapMemory, err := mem.SwapMemory()
	if err != nil {
		log.Fatal(err)
	}
	return func(m *Memory) {
		m.S = swapMemory
	}
}

func WithV(vri *mem.VirtualMemoryStat, err error) Opthions {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}
	return func(m *Memory) {
		m.V = memory
	}
}

func NewMemory(optins ...Opthions) *Memory {
	vrim, err := mem.VirtualMemory()
	if err != nil {
		log.Fatal(err)
	}

	sm, err := mem.SwapMemory()
	if err != nil {
		log.Fatal(err)
	}

	devices, err := mem.SwapDevices()
	if err != nil {
		log.Fatal(err)
	}
	m := &Memory{
		V:  vrim,
		S:  sm,
		SD: devices,
	}
	for _, o := range optins {
		o(m)
	}
	return m
}
