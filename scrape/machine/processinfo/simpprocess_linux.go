//go:build linux
// +build linux

package processinfo

import (
	"bigagent/scrape/machine/formatsize"
	"encoding/json"
	"log"
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
		log.Fatal(err)
	}

	for _, info := range psinfo {
		name, err := info.Name()
		if err != nil {
			log.Print(err)
		}
		username, err := info.Username()
		if err != nil {
			log.Print(err)
		}
		cpupercent, err := info.CPUPercent()
		if err != nil {
			log.Print(err)
		}
		mempercent, err := info.MemoryPercent()
		if err != nil {
			log.Print(err)
		}
		mem, err := info.MemoryInfo()
		if err != nil {
			log.Print(err)
		}
		tty, err := info.Terminal()
		if err != nil {
			log.Print(err)
		}
		stat, err := info.Status()
		if err != nil {
			log.Print(err)
		}
		start, err := info.CreateTime()
		if err != nil {
			log.Print(err)
		}
		t := time.Unix(start/1000, 0)
		starts := t.Format("2006-01-02 15:04:05")
		cmd, err := info.Cmdline()
		if err != nil {
			log.Print(err)
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
		log.Fatal(err)
	}
	return string(b)
}
