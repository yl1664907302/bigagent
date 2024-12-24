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
}
