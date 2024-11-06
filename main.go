package main

import "bigagent/inits"

func init() {
	inits.LoggerInit()
	inits.AgentRegister()
	inits.Crontab()
	inits.ListerChannel()
	inits.Hander("8080")

}
func main() {

}
