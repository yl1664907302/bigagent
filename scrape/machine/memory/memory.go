package memory

import (
	"github.com/shirou/gopsutil/v4/mem"
)

// Memory 定义Memory 类型
type Memory struct {
	VirtueMemory mem.VirtualMemoryStat `json:"virtue_memory"`
	SwapMemory   mem.SwapMemoryStat    `json:"swap_memory"`
}

// NewMemory 内存对象
func NewMemory() *Memory {
	virtualMemory, _ := mem.VirtualMemory()
	swapMemory, _ := mem.SwapMemory()
	return &Memory{
		VirtueMemory: *virtualMemory,
		SwapMemory:   *swapMemory,
	}
}

// GetVirtualMemoryStat 虚拟内存
func (memory *Memory) GetVirtualMemoryStat() mem.VirtualMemoryStat {
	return memory.VirtueMemory
}

// GetSwapStat 交换内存
func (memory *Memory) GetSwapStat() mem.SwapMemoryStat {
	return memory.SwapMemory
}
