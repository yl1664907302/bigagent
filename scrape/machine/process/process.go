package proessinfo

import (
	"context"
	"encoding/json"
	"github.com/shirou/gopsutil/v4/process"
	"log"
	"strings"
)

// ProcessInfo 进程结构体
type ProcessInfo struct {
	Process []*process.Process
}

// GetProcessInfo 获取进行信息
func (p *ProcessInfo) GetProcessInfo() ([]*process.Process, error) {
	processes, err := process.ProcessesWithContext(context.Background())
	return processes, err
}

// 获取进程所有信息
func (p *ProcessInfo) ToString() (string, error) {
	var v []string
	processes, err := p.GetProcessInfo()
	if err != nil {
		log.Println("get process info error", err)
		return "", err
	}
	for _, process := range processes {
		v = append(v, process.String())
	}
	ps := strings.Join(v, "")
	marshal, err := json.Marshal(ps)
	return string(marshal), err
}

// NewProcessInfo 进程对外接口
func NewProcessInfo() *ProcessInfo {
	return &ProcessInfo{}
}
