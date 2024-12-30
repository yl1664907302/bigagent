package processinfo

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v4/process"
	"log"
)

type PROCESS struct {
	P []*process.Process `json:"process_info"`
}

type Options func(p *PROCESS)

func (p *PROCESS) ToString() string {
	indent, err := json.MarshalIndent(p, "", " ")
	if err != nil {
		return err.Error()
	}
	return string(indent)
}

func WithP(ps []*process.Process, err error) Options {
	if err != nil {
		log.Fatal(err)
	}
	return func(p *PROCESS) {
		p.P = ps
	}
}

func NewProcess(options ...Options) *PROCESS {
	pros, err := process.Processes()
	if err != nil {
		log.Fatal(err)
	}
	p := &PROCESS{
		P: pros,
	}
	for _, option := range options {
		option(p)
	}
	return p
}
