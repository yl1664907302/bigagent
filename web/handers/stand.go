package handers

import (
	"bigagent/web"
	"bigagent/web/response"
	"bigagent/web/strategy"
	"log"
	"net/http"
)

type ApiStand struct{}

func (a *ApiStand) ShowData(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		agent := &web.Agent{}
		agent.SetApiStrategy(&strategy.StandardStrategy{})
		data, err := agent.ExecuteApi()
		if err != nil {
			log.Println(err)
		}
		response.SuccessWithDetailed(w, "获取主机信息成功", data)
	}
}

var ApiStandApp = new(ApiStand)
