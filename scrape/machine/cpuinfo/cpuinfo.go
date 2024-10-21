package cpuinfo

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"log"
)

type CpuInfo struct {
	CpuInfo []cpu.InfoStat
	TimeCpu []cpu.TimesStat
}

func GetCpuInfo() []cpu.InfoStat {
	info, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}
	return info
}

func GetTimeCpu() []cpu.TimesStat {
	times, err := cpu.Times(true)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return times
}
