package inits

import (
	"bigagent/web/router"
	"net/http"
)

type StandRouterGroup struct {
}

// 需要优化，自动带上“bigagent”前缀
func (r *StandRouterGroup) StandRouter() {
	http.HandleFunc("/bigagent/showdata", router.StandRouterApp.ShowData)
}

var StandRouterGroupApp = new(StandRouterGroup)
