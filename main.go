package main

import (
	model "bigagent/model/machine"
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
	fmt.Printf("model.NewSmpDataApi().ToString(): %v\n", model.NewSmpDataApi().ToString())
	fmt.Printf("model.NewSmpData().ToString(): %v\n", model.NewSmpData().ToString())
	// fmt.Printf("model.NewStandData().ToString(): %v\n", model.NewStandData().ToString())
	// fmt.Println("---------------------------------------------------------------------")
	// fmt.Printf("model.NewStandDataApi().ToString(): %v\n", model.NewStandDataApi().ToString())
}
