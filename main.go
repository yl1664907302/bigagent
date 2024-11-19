package main

import (
	"bigagent/config/global"
	"bigagent/inits"
)

func init() {
	inits.Viper()
	//inits.LoggerInit()
	inits.AgentRegister()
	inits.Crontab()
	inits.ListerChannel()
	inits.Hander(global.CONF.System.Addr)
}
func main() {

}
