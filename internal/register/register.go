package register

import (
	"bigagent/internal/strategy"
	"bigagent/internal/web/router"
	"log"
)

// 抽取通用注册流程
func registerCommon(hostEmpty bool, token string, openPush bool, onlyPush bool,
	missingMsg string,
	pushFactory func() strategy.PushStrategy,
	apiFactory func() strategy.ApiStrategy,
	attach func(*strategy.Agent),
) {
	if hostEmpty {
		return
	}
	agent := strategy.NewAgent()

	if onlyPush {
		if token == "" {
			log.Println(missingMsg)
			// 仅提示，不附加到router
		} else {
			agent.SetPushStrategy(pushFactory())
		}
	} else if openPush {
		if token == "" {
			log.Println(missingMsg)
		} else {
			agent.SetPushStrategy(pushFactory())
			attach(agent)
		}
	} else {
		agent.SetApiStrategy(apiFactory())
		attach(agent)
	}

	strategy.Agents = append(strategy.Agents, *agent)
}

func attachStand1(agent *strategy.Agent) {
	router.StandRouterApp.A = agent
	router.StandRouterApp.K = true
}

func attachStand2(agent *strategy.Agent) {
	router.StandRouterApp2.A = agent
	router.StandRouterApp2.K = true
}

// Stand1Register  策略注册,open push值是否开启push, only push是否只开启push（关闭api）
func Stand1Register(grpcHost string, token string, openPush bool, onlyPush bool) {
	registerCommon(
		grpcHost == "",
		token,
		openPush,
		onlyPush,
		"请配置push操作的token值",
		func() strategy.PushStrategy { return &strategy.StandardStrategy{K: token, G: grpcHost, KeyUse: true} },
		func() strategy.ApiStrategy { return &strategy.StandardStrategy{KeyUse: false} },
		attachStand1,
	)
}

// Stand2Register  策略注册,open push值是否开启push, only push是否只开启push（关闭api）
func Stand2Register(grpcHost string, token string, openPush bool, onlyPush bool) {
	registerCommon(
		grpcHost == "",
		token,
		openPush,
		onlyPush,
		"请配置push操作的host值",
		func() strategy.PushStrategy {
			return &strategy.StandardStrategy2{K: token, G: grpcHost, KeyUse: true}
		},
		func() strategy.ApiStrategy { return &strategy.StandardStrategy2{KeyUse: false} },
		attachStand2,
	)
}

// VeopsRegister 策略注册,openpush值是否开启push, onlypush是否只开启push（关闭api）
func VeopsRegister(host string, openpush bool, onlypush bool) {
	registerCommon(
		host == "",
		// Veops 不使用 token，沿用提示语义
		host,
		openpush,
		onlypush,
		"请配置push操作的host值",
		func() strategy.PushStrategy { return &strategy.VeopsStrategy{host, true} },
		func() strategy.ApiStrategy { return &strategy.VeopsStrategy{} },
		attachStand1,
	)
}
