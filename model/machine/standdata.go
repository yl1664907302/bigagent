package model

import (
	"bigagent/scrape/machine"
	"bigagent/scrape/machine/info"
	"bigagent/scrape/machine/memory"
	"encoding/json"
	"log"
)

// StandData 暴露原生utils数据
type StandData struct {
	Memory memory.Memory `json:"memory"`
	Info   info.Info     `json:"info"`
}

func NewStandData() *StandData {
	m := machine.Ma.M
	i := machine.Ma.I
	return &StandData{Memory: *m, Info: *i}
}

func NewStandDataApi() *StandData {
	m := memory.NewMemory()
	i := info.NewInfo()
	return &StandData{Memory: *m, Info: *i}
}

func (d *StandData) ToString() string {
	s, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}
	return string(s)
}
