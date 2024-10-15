package inits

import (
	"bigagent/web"
	"bigagent/web/register"
	"log"
	"net/http"
)

var Agents []web.Agent

// Handler http
func Hander(port string) {
	StandRouterGroupApp.StandRouter()

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// AgentRegister agent注册
func AgentRegister() {
	register.StandRegister("", false, false)
}
