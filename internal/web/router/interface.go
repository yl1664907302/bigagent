package router

import "net/http"

type Router interface {
	ShowData(w http.ResponseWriter, req *http.Request)
}
