package router

import (
	"bigagent/internal/k8s"
	"bigagent/internal/web/response"
	"net/http"
)

type K8sRouter struct {
	S k8s.Service
}

func (r *K8sRouter) Ping(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	response.SuccessWithDetailed(w, map[string]interface{}{"code": 0, "data": "pong"})
}

func (r *K8sRouter) Info(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	data, _ := r.S.ClusterInfo()
	response.SuccessWithDetailed(w, data)
}

var K8sRouterApp = &K8sRouter{}
