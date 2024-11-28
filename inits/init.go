package inits

import (
	"bigagent/config/global"
	"bigagent/register"
	"bigagent/scrape/machine"
	"bigagent/strategy"
	"bigagent/util/crontab"
	"bigagent/util/logger"
	"log"
	"net/http"
)

// Hander 启动http服务
func Hander(port string) {
	StandRouterGroupApp.StandRouter()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// AgentRegister agent注册
func AgentRegister() {
	register.StandRegister("127.0.0.1:8080", global.CONF.System.Client_port, true, false)
}

// Crontab 执行定时任务
func Crontab() {
	crontab.ScrapeCrontab()
}

// ListerChannel 监听channel
func ListerChannel() {
	go func() {
		for range machine.MachineCh {
			temp := <-machine.MachineCh
			if temp {
				logger.DefaultLogger.Info("数据更新，执行推送")
				machine.MachineCh <- false
				if strategy.Agents == nil || len(strategy.Agents) == 0 {
					logger.DefaultLogger.Warn("web.Agents is nil or empty")
					return
				}
				for _, agent := range strategy.Agents {
					err := agent.ExecutePush()
					if err != nil {
						logger.DefaultLogger.Error("数据推送异常:", err)
					}
				}
			}
		}
	}()
}

func LoggerInit() {
	logger.InitLogger(global.CONF.System.Logfile, "info", "json", true)
}
