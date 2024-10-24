package main

import (
	"bigagent/config/global"
	"bigagent/inits"
)

func init() {
	inits.LoggerInit()
	inits.AgentRegister()
	inits.Crontab()
	inits.ListerChannel()
	inits.Hander(global.CONF.System.Addr)
}
func main() {

}
