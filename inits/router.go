package inits

import (
	"bigagent/internal/config/global"
	utils "bigagent/internal/util"
	router2 "bigagent/internal/web/router"
	"net/http"
	"time"
)

type SysRouterGroup struct {
	Prefix string
}
type StandRouterGroup struct {
	Prefix string
}

type StandRouterGroup2 struct {
	Prefix string
}

type VeopsRouterGroup struct {
	Prefix string
}

type K8sRouterGroup struct {
	Prefix string
}

var StandRouterGroupApp = &StandRouterGroup{Prefix: "/stand1"}
var StandRouterGroupApp2 = &StandRouterGroup2{Prefix: "/stand2"}
var SysRouterGroupApp = &SysRouterGroup{Prefix: "/sys"}
var K8sRouterGroupApp = &K8sRouterGroup{Prefix: "/k8s"}

// var VeopsRouterGroupApp = &VeopsRouterGroup{Prefix: "/veops"}

// StandRouter 添加路由，自动带上前缀
func (r *StandRouterGroup) StandRouter() {
	http.Handle(r.Prefix+"/showdata", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router2.StandRouterApp.ShowData))))
	http.Handle(r.Prefix+"/patroldata", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router2.StandRouterApp.PatrolData))))
}

// Stand2Router 添加路由，自动带上前缀
func (r *StandRouterGroup2) StandRouter() {
	http.Handle(r.Prefix+"/showdata", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router2.StandRouterApp2.ShowData))))
}

// K8sRouter 添加路由，自动带上前缀
func (r *K8sRouterGroup) K8sRouter() {
	http.Handle(r.Prefix+"/ping", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router2.K8sRouterApp.Ping))))
	http.Handle(r.Prefix+"/info", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router2.K8sRouterApp.Info))))
}

func (r *SysRouterGroup) SysRouter() {
	http.Handle(r.Prefix+"/cmd", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router2.SysRouterApp.Cmd))))
}

// VeopsRouter 添加路由，自动带上前缀
// func (r *VeopsRouterGroup) VeopsRouter() {
// 	http.Handle(r.Prefix+"/showdata", loggingMiddleware(AuthMiddleware(http.HandlerFunc(router.VeopsRouterApp.ShowData))))
// }

// 日志中间件，记录每次请求的访问日志
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r) // 处理请求
		duration := time.Since(startTime)
		utils.DefaultLogger.Printf("Method: %s, URI: %s, RemoteAddr: %s, Duration: %v",
			r.Method, r.RequestURI, r.RemoteAddr, duration)
	})
}

// AuthMiddleware 验证请求头中的 Token
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != global.V.GetString("system.serct") { // 验证 Token 是否匹配
			http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r) // 验证通过，继续处理请求
	})
}
