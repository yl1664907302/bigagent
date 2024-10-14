package inits

import (
	"bigagent/web/router"
	"net/http"
)

type StandRouterGroup struct {
}

func (r *StandRouterGroup) StandRouter() {
	http.HandleFunc("/bigagent/showdata", router.StandRouterApp.ShowData)
}

var StandRouterGroupApp = new(StandRouterGroup)
