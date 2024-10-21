package data

import (
	"bigagent/util"
	"net/http"
	"time"
)

type DataApi struct{}

func (*DataApi) ShowData(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			util.Log.Println(err)
		}
	}()
	time.Sleep(150 * time.Millisecond)
	w.Write([]byte("hi girl"))
}
