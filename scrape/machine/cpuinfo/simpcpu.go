package cpuinfo

import (
	"bigagent/scrape/machine/formatsize"
	"encoding/json"
	"log"

	"github.com/shirou/gopsutil/v4/cpu"
)

type SmpCpu struct {
	Name  string `json:"name"`
	Core  int    `json:"core"`
	Usage string `json:"usage"`
}

func NewSmpCpu() *SmpCpu {

	c, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}

	usage, _ := cpu.Percent(0, false)
	if err != nil {
		log.Fatal(err)
	}

	return &SmpCpu{
		Name:  c[0].ModelName,
		Core:  len(c),
		Usage: formatsize.FormatPercent(usage[0]),
	}
}

func (s *SmpCpu) ToString() string {
	b, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	return string(b)
}
