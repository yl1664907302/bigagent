package cpuinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v4/cpu"
	"log"
	"strings"
)

// CpuInfo 定义cpu信息类型
type CpuInfo struct {
	CpuInfo *cpu.InfoStat
	TimeCpu *cpu.TimesStat
}

// GetCpuInfo 获取cpu信息
func (c *CpuInfo) GetCpuInfo() ([]cpu.InfoStat, error) {
	infoWithContext, err := cpu.InfoWithContext(context.Background())
	if err != nil {
		log.Println("get cpu info failed:", err)
		return nil, err
	}
	return infoWithContext, err
}

// GetTimeCpu 获取cpu time 相关信息
func (c *CpuInfo) GetTimeCpu() ([]cpu.TimesStat, error) {
	timesWithContext, err := cpu.TimesWithContext(context.Background(), true)
	if err != nil {
		log.Println("get cpu times info failed:", err)
	}
	return timesWithContext, err
}

// ToString cpu相关总数据
func (c *CpuInfo) ToString() (string, error) {
	var v []string

	infoStats, err := c.GetCpuInfo()
	if err != nil {
		log.Println("get cpu info failed:", err)
		return "", err
	}
	for _, cpus := range infoStats {
		v = append(v, cpus.String())
	}
	timesStats, err := c.GetTimeCpu()
	if err != nil {
		log.Println("get cpu times info failed:", err)
		return "", err
	}
	for _, times := range timesStats {
		v = append(v, times.String())
	}
	// v构造成 string
	infos := strings.Join(v, ",")
	fmt.Println(infos)
	// 返回json
	indent, err := json.MarshalIndent(infos, " ", " ")
	return string(indent), err
}

// NewCpuInfo 对外暴露接口
func NewCpuInfo() *CpuInfo {
	return &CpuInfo{}
}
