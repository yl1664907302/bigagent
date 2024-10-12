package inits

import (
	"bigagent/web/handers"
	"log"
	"net/http"
)

func Hander(port string) {
	http.HandleFunc("/bigagent", handers.ApiStandApp.ShowData)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
