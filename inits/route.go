package inits

import (
	"bigagent/route"
	"net/http"
)

type wr func(w http.ResponseWriter, r *http.Request)

func Router() {
	router := route.RouterGroupApp.DataRouter.AllRouter()
	iterator := router.CreateIterator()
	entry := iterator.Next()
	for entry != nil {
		value := entry.Value.(wr)
		http.HandleFunc(entry.Key, value)
		entry = iterator.Next()
	}
	http.ListenAndServe(":8080", nil)
}
