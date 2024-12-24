package register

import (
	"bigagent/strategy"
	"bigagent/web/router"
	"log"
)

// Stand1Register  策略注册,open push值是否开启push, only push是否只开启push（关闭api）
func Stand1Register(host string, grpcHost string, openPush bool, onlyPush bool) {
	if grpcHost == "" {
		return
	}
	agent := strategy.NewAgent()
	if !router.StandRouterApp.K {
		if onlyPush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy{H: host, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch host {
				case "":
					agent.SetApiStrategy(&strategy.StandardStrategy{})
					router.StandRouterApp.A = agent
					router.StandRouterApp.K = true
				default:
					agent.SetApiStrategy(&strategy.StandardStrategy{})
					agent.SetPushStrategy(&strategy.StandardStrategy{H: host, G: grpcHost})
					router.StandRouterApp.A = agent
					router.StandRouterApp.K = true
				}
			default:
				agent.SetApiStrategy(&strategy.StandardStrategy{})
				router.StandRouterApp.A = agent
				router.StandRouterApp.K = true
			}
		}
	} else {
		if onlyPush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy{H: host, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch host {
				case "":
					log.Println("请配置push操作的host值")
				default:
					agent.SetPushStrategy(&strategy.StandardStrategy{H: host, G: grpcHost})
				}
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy{H: host, G: grpcHost})
			}
		}
	}
	strategy.Agents = append(strategy.Agents, *agent)
}

// Stand2Register  策略注册,open push值是否开启push, only push是否只开启push（关闭api）
func Stand2Register(host string, grpcHost string, openPush bool, onlyPush bool) {
	if grpcHost == "" {
		return
	}
	agent := strategy.NewAgent()
	if !router.StandRouterApp2.K {
		if onlyPush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy2{H: host, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch host {
				case "":
					agent.SetApiStrategy(&strategy.StandardStrategy2{})
					router.StandRouterApp2.A = agent
					router.StandRouterApp2.K = true
				default:
					agent.SetApiStrategy(&strategy.StandardStrategy2{})
					agent.SetPushStrategy(&strategy.StandardStrategy2{H: host, G: grpcHost})
					router.StandRouterApp2.A = agent
					router.StandRouterApp2.K = true
				}
			default:
				agent.SetApiStrategy(&strategy.StandardStrategy2{})
				router.StandRouterApp2.A = agent
				router.StandRouterApp2.K = true
			}
		}
	} else {
		if onlyPush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy2{H: host, G: grpcHost})
			}
		} else {
			switch openPush {
			case true:
				switch host {
				case "":
					log.Println("请配置push操作的host值")
				default:
					agent.SetPushStrategy(&strategy.StandardStrategy2{H: host, G: grpcHost})
				}
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy2{H: host, G: grpcHost})
			}
		}
	}
	strategy.Agents = append(strategy.Agents, *agent)
}

// VeopsRegister 策略注册,openpush值是否开启push, onlypush是否只开启push（关闭api）
//func VeopsRegister(host string, openpush bool, onlypush bool) {
//	agent := strategy.NewAgent()
//	if !router.StandRouterApp.K {
//		if onlypush {
//			switch host {
//			case "":
//				log.Println("请配置push操作的host值")
//			default:
//				agent.SetPushStrategy(&strategy.VeopsStrategy{host})
//			}
//		} else {
//			switch openpush {
//			case true:
//				switch host {
//				case "":
//					agent.SetApiStrategy(&strategy.VeopsStrategy{})
//					router.StandRouterApp.A = &agent
//				default:
//					agent.SetApiStrategy(&strategy.VeopsStrategy{})
//					agent.SetPushStrategy(&strategy.VeopsStrategy{host})
//					router.StandRouterApp.A = &agent
//				}
//			default:
//				agent.SetApiStrategy(&strategy.VeopsStrategy{})
//				router.StandRouterApp.A = &agent
//			}
//		}
//	} else {
//		if onlypush {
//			switch host {
//			case "":
//				log.Println("请配置push操作的host值")
//			default:
//				agent.SetPushStrategy(&strategy.VeopsStrategy{host})
//				router.StandRouterApp.A = &agent
//			}
//		} else {
//			switch openpush {
//			case true:
//				switch host {
//				case "":
//					log.Println("请配置push操作的host值")
//				default:
//					agent.SetPushStrategy(&strategy.VeopsStrategy{host})
//					router.StandRouterApp.A = &agent
//				}
//			default:
//				agent.SetPushStrategy(&strategy.VeopsStrategy{host})
//				router.StandRouterApp.A = &agent
//			}
//		}
//	}
//	strategy.Agents = append(strategy.Agents, agent)
//}
