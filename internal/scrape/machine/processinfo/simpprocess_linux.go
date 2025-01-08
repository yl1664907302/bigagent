//go:build linux
// +build linux

package processinfo

import (
	"bigagent/internal/scrape/machine/formatsize"
	utils "bigagent/internal/util"
	grpc_server "bigagent/internal/web/grpcs/server"
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type smpInfo struct {
	Name       string   `json:"name"`
	User       string   `json:"user"`
	Pid        int32    `json:"pid"`
	CpuPercent string   `json:"cpu_percent"`
	MemPercent string   `json:"mem_percent"`
	Vsz        string   `json:"vritual_memory_size"`
	Rss        string   `json:"Resident_Set_Size"`
	Tty        string   `json:"tty"`
	Stat       []string `json:"stat"`
	Start      string   `json:"start_time"`
	Cmd        string   `json:"cmd"`
}

type SmpPs map[string]smpInfo

func NewSmpPs() *SmpPs {
	smpps := make(SmpPs)
	psinfo, err := process.Processes()
	if err != nil {
		utils.DefaultLogger.Error(err)
	}

	for _, info := range psinfo {
		name, err := info.Name()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		username, err := info.Username()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		cpupercent, err := info.CPUPercent()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		mempercent, err := info.MemoryPercent()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		mem, err := info.MemoryInfo()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		tty, err := info.Terminal()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		stat, err := info.Status()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		start, err := info.CreateTime()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}
		t := time.Unix(start/1000, 0)
		starts := t.Format("2006-01-02 15:04:05")
		cmd, err := info.Cmdline()
		if err != nil {
			utils.DefaultLogger.Error(err)
		}

		smpinfo := smpInfo{
			Name:       name,
			User:       username,
			Pid:        info.Pid,
			CpuPercent: formatsize.FormatPercent(cpupercent),
			MemPercent: formatsize.FormatPercent(float64(mempercent)),
			Vsz:        formatsize.FormatSize(mem.VMS),
			Rss:        formatsize.FormatSize(mem.RSS),
			Tty:        tty,
			Stat:       stat,
			Start:      starts,
			Cmd:        cmd,
		}
		smpps[name] = smpinfo
	}

	return &smpps
}

func (p *SmpPs) ToString() string {
	b, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		utils.DefaultLogger.Error(err)
	}
	return string(b)
}

func NewSmpPsGrpc() *map[string]*grpc_server.SmPsInfo {
	smpps := make(map[string]*grpc_server.SmPsInfo)
	psinfo, err := process.Processes()
	if err != nil {
		utils.DefaultLogger.Error(err)
	}

	for _, info := range psinfo {
		name, err := info.Name()
		if err != nil || name == "" {
			continue
		}
		username, err := info.Username()
		if err != nil || username == "" {
			continue
		}
		cpupercent, err := info.CPUPercent()
		if err != nil {
			utils.DefaultLogger.Errorf("get ps cpupercent err: %s", err)
			continue
		}
		mempercent, err := info.MemoryPercent()
		if err != nil {
			utils.DefaultLogger.Errorf("get ps mempercent err: %s", err)
			continue
		}
		mem, err := info.MemoryInfo()
		if err != nil {
			utils.DefaultLogger.Errorf("get ps meminfo err: %s", err)
			continue
		}

		start, err := info.CreateTime()
		if err != nil {
			utils.DefaultLogger.Errorf("get ps start err: %s", err)
			continue
		}
		t := time.Unix(start/1000, 0)
		starts := t.Format("2006-01-02 15:04:05")
		cmd, err := info.Cmdline()
		if err != nil {
			utils.DefaultLogger.Errorf(err.Error())
		}
		// wg.Lock()
		smpps[name] = &grpc_server.SmPsInfo{
			Name:              name,
			User:              username,
			Pid:               string(info.Pid),
			CpuPercent:        formatsize.FormatPercent(cpupercent),
			MemPercent:        formatsize.FormatPercent(float64(mempercent)),
			VritualMemorySize: formatsize.FormatSize(mem.VMS),
			ResidentSetSize:   formatsize.FormatSize(mem.RSS),
			StartTime:         starts,
			Cmd:               cmd,
		}
		// wg.Unlock()
	}

	return &smpps
}
