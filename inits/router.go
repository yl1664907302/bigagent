package inits

import (
	"bigagent/util/logger"
	"bigagent/web/router"
	"net/http"
	"time"
)

type StandRouterGroup struct {
}
type VeopsRouterGroup struct{}

// 需要优化，自动带上“bigagent”前缀
func (r *StandRouterGroup) StandRouter() {
	http.Handle("/bigagent/showdata", loggingMiddleware(http.HandlerFunc(router.StandRouterApp.ShowData)))
}

func (r *VeopsRouterGroup) VeopsRouter() {
	http.Handle("/veops/showdata", loggingMiddleware(http.HandlerFunc(router.VeopsRouterApp.ShowData)))
}

var StandRouterGroupApp = new(StandRouterGroup)
var VeopsRouterGroupApp = new(VeopsRouterGroup)

// 日志中间件，记录每次请求的访问日志(闭包和适配器（可接口匿名实现）的骚操作)
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r) //使用适配器方法，等用于ShowData方法，闭包可以直接使用
		logger.DefaultLogger.Println(
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(time.Now()),
		)
	})
}
