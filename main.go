package main

import (
	model "bigagent/model/machine"
	"encoding/json"
	"fmt"
)

//	func init() {
//		inits.Viper()
//		inits.LoggerInit()
//		inits.Crontab()
//		inits.AgentRegister()
//		inits.ListerChannel()
//	}
func main() {
	// inits.RunG()
	// inits.Hander(global.V.GetString("system.addr"))
	// smma := machine.NewSmpMachine()
	smma := model.NewSmpData()
	b, _ := json.MarshalIndent(smma, "", " ")
	fmt.Println(string(b))
	// sma := machine.NewSmpMachine()
	// fmt.Printf("sma.Memory.Vmem.UsedPercent: %v\n", sma.Memory.Vmem.UsedPercent)
	// fmt.Println("-------------------------------")
	// fmt.Printf("sma.Cpu.Usage: %v\n", sma.Cpu.Usage)

	// fmt.Printf("meminfo.NewMemory().V.UsedPercent: %v\n", meminfo.NewMemory().V.UsedPercent)
}
