package main

import "bigagent/inits"

func init() {
	inits.AgentRegister()
	inits.Crontab()
	inits.Hander("8080")
}
func main() {

}
