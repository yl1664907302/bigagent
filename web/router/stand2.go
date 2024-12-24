package router

import (
	"bigagent/strategy"
	"bigagent/web/response"
	"log"
	"net/http"
)

type StandRouter2 struct {
	K bool
	A *strategy.Agent
}

func (r *StandRouter2) ShowData(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		data, err := r.A.ExecuteApi("showdata")
		if err != nil {
			log.Println(err)
		}
		response.SuccessWithDetailed(w, data)
	}
}

var StandRouterApp2 = &StandRouter2{false, nil}
