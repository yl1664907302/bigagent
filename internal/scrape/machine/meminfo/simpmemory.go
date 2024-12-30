package meminfo

import (
	"bigagent/internal/scrape/machine/formatsize"
	utils "bigagent/internal/util"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v4/mem"
)

type SmpMem struct {
	Total       string `json:"total"`
	Used        string `json:"used"`
	Free        string `json:"free"`
	UsedPercent string `json:"usedpercent"`
}

type SmpSwap struct {
	Total       string `json:"total"`
	Used        string `json:"used"`
	Free        string `json:"free"`
	UsedPercent string `json:"usedpercent"`
}

type SmpMemory struct {
	Vmem *SmpMem  `json:"virtual_memory"`
	Swap *SmpSwap `json:"swap_memory"`
}

// NewSmpMem 创建并返回一个新的 SmpMemory 实例
func NewSmpMem() *SmpMemory {
	// 获取虚拟内存信息
	m, err := mem.VirtualMemory()
	if err != nil {
		utils.DefaultLogger.Error(err)
	}

	// 获取交换区内存信息
	sm, err := mem.SwapMemory()
	if err != nil {
		utils.DefaultLogger.Error(err)
	}

	// 获取物理内存信息

	// 构建 SmpMemory
	return &SmpMemory{

		Vmem: &SmpMem{
			Total:       formatsize.FormatSize(m.Total),
			Used:        formatsize.FormatSize(m.Used),
			Free:        formatsize.FormatSize(m.Free),
			UsedPercent: fmt.Sprintf("%.2f%%", m.UsedPercent),
		},
		Swap: &SmpSwap{
			Total:       formatsize.FormatSize(sm.Total),
			Used:        formatsize.FormatSize(sm.Used),
			Free:        formatsize.FormatSize(sm.Free),
			UsedPercent: fmt.Sprintf("%.2f%%", sm.UsedPercent),
		},
	}
}

func (m *SmpMemory) ToString() string {
	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return string(b)
}
