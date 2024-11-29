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
func AgentRegister(i int) {
	register.StandRegister("127.0.0.1:8080", global.V.GetString("system.client_port"), true, false, i)
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
				for i, agent := range strategy.Agents {
					err := agent.ExecutePush()
					if err != nil {
						logger.DefaultLogger.Errorf("agent策略序号：%d,数据推送异常: %s", i, err)
					}

				}
			}
		}
	}()
}

func LoggerInit() {
	logger.InitLogger(global.V.GetString("system.logfile"), "info", "json", true)
}
