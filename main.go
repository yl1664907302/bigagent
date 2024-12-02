package main

import (
	"bigagent/config/global"
	"bigagent/inits"
)

func init() {
	inits.Viper()
	inits.LoggerInit()
	inits.Crontab()
	inits.ListerChannel()
}
func main() {
	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
