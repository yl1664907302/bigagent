package inits

import (
	"bigagent/web/register"
	"bigagent/web/router"
	"log"
	"net/http"
)

// http监听器
func Hander(port string) {
	StandRouterGroupApp.StandRouter()
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// agent注册
func AgentRegister() {
	register.StandRegister(router.StandRouterApp)
}
