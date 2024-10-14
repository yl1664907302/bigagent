package memory

import (
	"github.com/shirou/gopsutil/v4/mem"
	"log"
)

type Memory struct {
	V1 mem.VirtualMemoryStat `json:"strategy"`
}

func NewMemory() *Memory {
	memory, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
	}
	return &Memory{V1: *memory}
}

func (m *Memory) Total() uint64 {
	return m.V1.Total
}

func (m *Memory) Used() uint64 {
	return m.V1.Used
}

func (m *Memory) Free() uint64 {
	return m.V1.Free
}

func (m *Memory) Available() uint64 {
	return m.V1.Available
}

func (m *Memory) UsedPercent() float64 {
	return m.V1.UsedPercent
}
