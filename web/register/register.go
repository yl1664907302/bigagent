package register

import (
	"bigagent/web"
	"bigagent/web/router"
	"bigagent/web/strategy"
)

// 策略注册
func StandRegister(a *router.StandRouter) {
	a.A = web.NewAgent()
	a.A.SetApiStrategy(&strategy.StandardStrategy{})
	a.A.SetPushStrategy(&strategy.StandardStrategy{})
}
