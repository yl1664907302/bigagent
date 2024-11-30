package main

import (
	"bigagent/scrape/machine"
	"encoding/json"
	"fmt"
)

//	func init() {
//		inits.Viper()
//		inits.LoggerInit()
//		inits.AgentRegister()
//		inits.Crontab()
//		inits.ListerChannel()
//		inits.Hander(global.CONF.System.Addr)
//	}
func main() {
	smpm := machine.NewSmpMachine()
	b, _ := json.MarshalIndent(smpm, "", " ")
	// b, _ := json.Marshal(smpm)
	fmt.Println(string(b))
	// km := kmodules.NewKmodules()
	// fmt.Printf("km.ToString(): %v\n", km.ToString())

}
