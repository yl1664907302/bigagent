package machine

import (
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/memory"
)

// Machine 存放所有的采集层数据，被懒汉式创建
type Machine struct {
	I *info.Info     `json:"i"`
	M *memory.Memory `json:"m"`
}

func NewMachine() *Machine {
	return &Machine{I: info.NewInfo(), M: memory.NewMemory()}
}

var Ma *Machine
