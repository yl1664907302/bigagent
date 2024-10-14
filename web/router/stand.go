package router

import (
	"bigagent/web"
	"bigagent/web/response"
	"log"
	"net/http"
)

type StandRouter struct {
	A *web.Agent
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

var StandRouterApp = new(StandRouter)
