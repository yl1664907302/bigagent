package main

import (
	"bigagent/config/global"
	"bigagent/inits"
)

func init() {
	inits.Viper()
	inits.LoggerInit()
	inits.AgentRegister()
}
func main() {
	inits.RunG()
	inits.Crontab()
	inits.ListerChannel()
	inits.Hander(global.CONF.System.Addr)
}
