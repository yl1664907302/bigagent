package meminfo

import (
	"github.com/shirou/gopsutil/v4/mem"
	"log"
)

// 内存信息
type MemInfo struct {
	VirMem  *mem.VirtualMemoryStat
	SwapMem *mem.SwapMemoryStat
}

func GetVirMem() *mem.VirtualMemoryStat {
	virtualMemory, err := mem.VirtualMemory()
	if err != nil {
		log.Println(err)
		return nil
	}
	return virtualMemory
}

func GetSwapMem() *mem.SwapMemoryStat {
	SwapMemory, err := mem.SwapMemory()
	if err != nil {
		log.Println(err)
		return nil
	}
	return SwapMemory
}
