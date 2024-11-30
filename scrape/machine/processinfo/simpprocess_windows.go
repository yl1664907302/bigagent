//go:build windows
// +build windows

package processinfo

import (
	"bigagent/scrape/machine/formatsize"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/process"
)

type SmpInfo struct {
	Name       string `json:"name"`
	User       string `json:"user"`
	Pid        int32  `json:"pid"`
	CpuPercent string `json:"cpu_percent"`
	MemPercent string `json:"mem_percent"`
	Vsz        string `json:"vritual_memory_size"`
	Rss        string `json:"Resident_Set_Size"`
	Start      string `json:"start_time"`
	Cmd        string `json:"cmd"`
}

type SmpPs map[string]*SmpInfo

var wg sync.RWMutex

func NewSmpPs() *SmpPs {
	smpps := make(map[string]*SmpInfo)
	// psinfo, err := process.Processes()
	psinfo, err := process.Processes()
	if err != nil {
		log.Fatal(err)
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
			log.Printf("get ps cpupercent err: %s", err)
			continue
		}
		mempercent, err := info.MemoryPercent()
		if err != nil {
			log.Printf("get ps mempercent err: %s", err)
			continue
		}
		mem, err := info.MemoryInfo()
		if err != nil {
			log.Printf("get ps meminfo err: %s", err)
			continue
		}

		start, err := info.CreateTime()
		if err != nil {
			log.Printf("get ps start err: %s", err)
			continue
		}
		t := time.Unix(start/1000, 0)
		starts := t.Format("2006-01-02 15:04:05")
		cmd, err := info.Cmdline()
		if err != nil {
			log.Print(err)
		}
		wg.Lock()
		smpps[name] = &SmpInfo{
			Name:       name,
			User:       username,
			Pid:        info.Pid,
			CpuPercent: formatsize.FormatPercent(cpupercent),
			MemPercent: formatsize.FormatPercent(float64(mempercent)),
			Vsz:        formatsize.FormatSize(mem.VMS),
			Rss:        formatsize.FormatSize(mem.RSS),
			Start:      starts,
			Cmd:        cmd,
		}
		wg.Unlock()
	}

	return (*SmpPs)(&smpps)
}

func (p *SmpPs) ToString() string {
	b, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
