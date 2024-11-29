package main

import (
	"bigagent/config/global"
	"bigagent/inits"
)

func init() {
	inits.Viper()
	inits.LoggerInit()
	inits.RunG()
	inits.Crontab()
	inits.ListerChannel()
	inits.Hander(global.V.GetString("system.addr"))
}
func main() {

}
