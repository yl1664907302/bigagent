package router

import (
	"bigagent/strategy"
	"bigagent/web/response"
	"log"
	"net/http"
)

type StandRouter struct {
	K bool
	A *strategy.Agent
}

func (r *StandRouter) ShowData(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		data, err := r.A.ExecuteApi("bigagent")
		if err != nil {
			log.Println(err)
		}
		response.SuccessWithDetailed(w, data)
	}
}

var StandRouterApp = &StandRouter{false, nil}
