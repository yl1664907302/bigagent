package register

import (
	"bigagent/web"
	"bigagent/web/router"
	"bigagent/web/strategy"
	"log"
)

// StandRegister 策略注册,openpush值是否开启push, onlypush是否只开启push（关闭api）
func StandRegister(host string, grpc_host string, openpush bool, onlypush bool) {
	agent := web.NewAgent()
	if !router.StandRouterApp.K {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy{host, grpc_host})
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					agent.SetApiStrategy(&strategy.StandardStrategy{})
					router.StandRouterApp.A = &agent
				default:
					agent.SetApiStrategy(&strategy.StandardStrategy{})
					agent.SetPushStrategy(&strategy.StandardStrategy{host, grpc_host})
					router.StandRouterApp.A = &agent
				}
			default:
				agent.SetApiStrategy(&strategy.StandardStrategy{})
				router.StandRouterApp.A = &agent
			}
		}
	} else {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy{host, grpc_host})
				router.StandRouterApp.A = &agent
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					log.Println("请配置push操作的host值")
				default:
					agent.SetPushStrategy(&strategy.StandardStrategy{host, grpc_host})
					router.StandRouterApp.A = &agent
				}
			default:
				agent.SetPushStrategy(&strategy.StandardStrategy{host, grpc_host})
				router.StandRouterApp.A = &agent
			}
		}
	}
	web.Agents = append(web.Agents, agent)
}

// VeopsRegister 策略注册,openpush值是否开启push, onlypush是否只开启push（关闭api）
func VeopsRegister(host string, openpush bool, onlypush bool) {
	agent := web.NewAgent()
	if !router.StandRouterApp.K {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.VeopsStrategy{host})
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					agent.SetApiStrategy(&strategy.VeopsStrategy{})
					router.StandRouterApp.A = &agent
				default:
					agent.SetApiStrategy(&strategy.VeopsStrategy{})
					agent.SetPushStrategy(&strategy.VeopsStrategy{host})
					router.StandRouterApp.A = &agent
				}
			default:
				agent.SetApiStrategy(&strategy.VeopsStrategy{})
				router.StandRouterApp.A = &agent
			}
		}
	} else {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				agent.SetPushStrategy(&strategy.VeopsStrategy{host})
				router.StandRouterApp.A = &agent
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					log.Println("请配置push操作的host值")
				default:
					agent.SetPushStrategy(&strategy.VeopsStrategy{host})
					router.StandRouterApp.A = &agent
				}
			default:
				agent.SetPushStrategy(&strategy.VeopsStrategy{host})
				router.StandRouterApp.A = &agent
			}
		}
	}
	web.Agents = append(web.Agents, agent)
}
