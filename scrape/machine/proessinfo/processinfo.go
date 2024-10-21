package proessinfo

import (
	"github.com/shirou/gopsutil/v4/process"
	"log"
)

// proessinfo
type ProcessInfo struct {
	Process []*process.Process
}

func GetProcessInfo() []*process.Process {
	processes, err := process.Processes()
	if err != nil {
		log.Println(err)
		return nil
	}
	return processes
}
