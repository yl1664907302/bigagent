package cpuinfo

import (
	"bigagent/internal/scrape/machine/formatsize"
	utils "bigagent/internal/util"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type SmpCpu struct {
	Name  string `json:"name"`
	Core  int64  `json:"core"`
	Usage string `json:"usage"`
}

func GetCpuPercent() float64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return percent[0]
}

func NewSmpCpu() *SmpCpu {
	defer func() {
		if r := recover(); r != nil {
			utils.DefaultLogger.Error(fmt.Sprintf("捕获到 panic: %v", r))
		}
	}()
	c, err := cpu.Info()
	if err != nil {
		utils.DefaultLogger.Error(err)
	}

	usage, _ := cpu.Percent(time.Second, false)
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	cores := runtime.NumCPU()

	return &SmpCpu{
		Name: c[0].ModelName,
		Core: int64(cores),
		// Usage: fmt.Sprintf("%.2f%%", GetCpuPercent()),
		Usage: formatsize.FormatPercent(usage[0]),
	}
}

func (s *SmpCpu) ToString() string {
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return string(b)
}
