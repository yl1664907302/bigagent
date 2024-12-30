package cpuinfo

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v4/cpu"
	"log"
)

type Cpus struct {
	C []cpu.InfoStat  `json:"cpu_info"`
	T []cpu.TimesStat `json:"cpu_times"`
}

type Options func(c *Cpus)

func (cpus *Cpus) ToString() string {
	indent, err := json.MarshalIndent(cpus, "", " ")
	if err != nil {
		return err.Error()
	}
	return string(indent)
}

func WithC(info []cpu.InfoStat, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(c *Cpus) {
		c.C = info
	}
}

func WithT(times []cpu.TimesStat, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(c *Cpus) {
		c.T = times
	}
}

func NewCPU(options ...Options) *Cpus {
	info, err := cpu.Info()
	if err != nil {
		log.Fatal(err)
	}
	times, err := cpu.Times(true)
	if err != nil {
		log.Fatal(err)
	}

	cpus := &Cpus{
		C: info,
		T: times,
	}
	for _, option := range options {
		option(cpus)
	}
	return cpus
}
