package register

import (
	strategy2 "bigagent/internal/strategy"
	router2 "bigagent/internal/web/router"
	"log"
)

// Stand1Register  策略注册,open push值是否开启push, only push是否只开启push（关闭api）
func Stand1Register(grpcHost string, token string, openPush bool, onlyPush bool) {
	if grpcHost == "" {
		return
	}
	agent := strategy2.NewAgent()
	if !router2.StandRouterApp.K {
		if onlyPush {
			switch token {
			case "":
				log.Println("请配置push操作的token值")
			default:
				agent.SetPushStrategy(&strategy2.StandardStrategy{K: token, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch token {
				case "":
					agent.SetApiStrategy(&strategy2.StandardStrategy{})
					router2.StandRouterApp.A = agent
					router2.StandRouterApp.K = true
				default:
					agent.SetApiStrategy(&strategy2.StandardStrategy{})
					agent.SetPushStrategy(&strategy2.StandardStrategy{K: token, G: grpcHost})
					router2.StandRouterApp.A = agent
					router2.StandRouterApp.K = true
				}
			default:
				agent.SetApiStrategy(&strategy2.StandardStrategy{KeyUse: false})
				router2.StandRouterApp.A = agent
				router2.StandRouterApp.K = true
			}
		}
	} else {
		if onlyPush {
			switch token {
			case "":
				log.Println("请配置push操作的token值")
			default:
				agent.SetPushStrategy(&strategy2.StandardStrategy{K: token, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch token {
				case "":
					log.Println("请配置push操作的token值")
				default:
					agent.SetPushStrategy(&strategy2.StandardStrategy{K: token, G: grpcHost})
				}
			default:
				agent.SetPushStrategy(&strategy2.StandardStrategy{K: token, G: grpcHost})
			}
		}
	}
	strategy2.Agents = append(strategy2.Agents, *agent)
}

// Stand2Register  策略注册,open push值是否开启push, only push是否只开启push（关闭api）
func Stand2Register(grpcHost string, token string, openPush bool, onlyPush bool) {
	if grpcHost == "" {
		return
	}
	agent := strategy2.NewAgent()
	if !router2.StandRouterApp2.K {
		if onlyPush {
			switch token {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy2.StandardStrategy2{K: token, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch token {
				case "":
					agent.SetApiStrategy(&strategy2.StandardStrategy2{})
					router2.StandRouterApp2.A = agent
					router2.StandRouterApp2.K = true
				default:
					agent.SetApiStrategy(&strategy2.StandardStrategy2{})
					agent.SetPushStrategy(&strategy2.StandardStrategy2{K: token, G: grpcHost})
					router2.StandRouterApp2.A = agent
					router2.StandRouterApp2.K = true
				}
			default:
				agent.SetApiStrategy(&strategy2.StandardStrategy2{KeyUse: false})
				router2.StandRouterApp2.A = agent
				router2.StandRouterApp2.K = true
			}
		}
	} else {
		if onlyPush {
			switch token {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy2.StandardStrategy2{K: token, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch token {
				case "":
					log.Println("请配置push操作的host值")
				default:
					agent.SetPushStrategy(&strategy2.StandardStrategy2{K: token, G: grpcHost})
				}
			default:
				agent.SetPushStrategy(&strategy2.StandardStrategy2{K: token, G: grpcHost})
			}
		}
	}
	strategy2.Agents = append(strategy2.Agents, *agent)
}

// VeopsRegister 策略注册,openpush值是否开启push, onlypush是否只开启push（关闭api）
func VeopsRegister(host string, openpush bool, onlypush bool) {
	if host == "" {
		return
	}
	agent := strategy2.NewAgent()
	if !router2.StandRouterApp.K {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy2.VeopsStrategy{host, true})
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					agent.SetApiStrategy(&strategy2.VeopsStrategy{})
					router2.StandRouterApp.A = agent
				default:
					agent.SetApiStrategy(&strategy2.VeopsStrategy{})
					agent.SetPushStrategy(&strategy2.VeopsStrategy{host, true})
					router2.StandRouterApp.A = agent
				}
			default:
				agent.SetApiStrategy(&strategy2.VeopsStrategy{})
				router2.StandRouterApp.A = agent
			}
		}
	} else {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy2.VeopsStrategy{host, true})
				router2.StandRouterApp.A = agent
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					log.Println("请配置push操作的host值")
				default:
					agent.SetPushStrategy(&strategy2.VeopsStrategy{host, true})
					router2.StandRouterApp.A = agent
				}
			default:
				agent.SetPushStrategy(&strategy2.VeopsStrategy{host, true})
				router2.StandRouterApp.A = agent
			}
		}
	}
	strategy2.Agents = append(strategy2.Agents, *agent)
}
