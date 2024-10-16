package inits

import (
	"bigagent/register"
	"bigagent/util/crontab"
	"log"
	"net/http"
)

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

// Crontab 执行定时任务
func Crontab() {
	crontab.ScrapeCrontab()
}
