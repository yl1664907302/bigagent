package inits

import (
	"bigagent/config/global"
	"bigagent/register"
	"bigagent/scrape/machine"
	"bigagent/util/crontab"
	logger "bigagent/util/logger"
	"bigagent/web"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

// Handler http
func Hander(port string) {
	StandRouterGroupApp.StandRouter()
	//VeopsRouterGroupApp.VeopsRouter()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// AgentRegister agent注册
func AgentRegister() {
	register.StandRegister("127.0.0.1:8080", global.CONF.System.Grpc, true, false)
	//register.VeopsRegister("192.x.x.1", true, true)
}

// Crontab 执行定时任务
func Crontab() {
	crontab.ScrapeCrontab()
}

// Channel
func ListerChannel() {
	go func() {
		for range machine.MachineCh {
			temp := <-machine.MachineCh
			if temp {
				logger.DefaultLogger.Info("数据更新，执行推送")
				machine.MachineCh <- false
				if web.Agents == nil || len(web.Agents) == 0 {
					logger.DefaultLogger.Warn("web.Agents is nil or empty")
					return
				}
				for _, agent := range web.Agents {
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
	logger.InitLogger(viper.GetString(global.CONF.System.Logfile), viper.GetString("info"), viper.GetString("json"), true)
}
