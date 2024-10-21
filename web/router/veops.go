package router

import (
	"bigagent/web"
	"bigagent/web/response"
	"log"
	"net/http"
)

type VeopsRouter struct {
	K bool
	A *web.Agent
}

func (r *VeopsRouter) ShowData(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		data, err := r.A.ExecuteApi("xxx")
		if err != nil {
			log.Println(err)
		}
		response.SuccessWithDetailed(w, data)
	}
}

var VeopsRouterApp = &VeopsRouter{false, nil}
