package register

import (
	"bigagent/inits"
	"bigagent/web"
	"bigagent/web/router"
	"bigagent/web/strategy"
	"log"
)

// StandRegister 策略注册,openpush值是否开启push, onlypush是否只开启push（关闭api）
func StandRegister(host string, openpush bool, onlypush bool) {
	agent := web.NewAgent()
	if !router.StandRouterApp.K {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				router.StandRouterApp.A = &agent
				router.StandRouterApp.A.SetPushStrategy(&strategy.StandardStrategy{host})
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					router.StandRouterApp.A = &agent
					router.StandRouterApp.A.SetApiStrategy(&strategy.StandardStrategy{})
				default:
					router.StandRouterApp.A = &agent
					router.StandRouterApp.A.SetApiStrategy(&strategy.StandardStrategy{})
					router.StandRouterApp.A.SetPushStrategy(&strategy.StandardStrategy{host})
				}
			default:
				router.StandRouterApp.A = &agent
				router.StandRouterApp.A.SetApiStrategy(&strategy.StandardStrategy{})
			}
		}
	} else {
		if onlypush {
			switch host {
			case "":
				log.Println("请配置push操作的host值")
			default:
				router.StandRouterApp.A = &agent
				router.StandRouterApp.A.SetPushStrategy(&strategy.StandardStrategy{host})
			}
		} else {
			switch openpush {
			case true:
				switch host {
				case "":
					log.Println("请配置push操作的host值")
				default:
					router.StandRouterApp.A = &agent
					router.StandRouterApp.A.SetPushStrategy(&strategy.StandardStrategy{host})
				}
			default:
				router.StandRouterApp.A = &agent
				router.StandRouterApp.A.SetPushStrategy(&strategy.StandardStrategy{host})
			}
		}
	}
	inits.Agents = append(inits.Agents, agent)
}
