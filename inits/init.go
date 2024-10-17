package inits

import (
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

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// AgentRegister agent注册
func AgentRegister() {
	register.StandRegister("", false, false)
}

// Crontab 执行定时任务
func Crontab() {
	crontab.ScrapeCrontab()
}

// Channel
func ListerChannel() {
	go func() {
		for range machine.MachineCh {
			for _, agent := range web.Agents {
				err := agent.ExecutePush()
				if err != nil {
					logger.DefaultLogger.Error("Agent execute push error:", err)
				}
			}
		}
	}()
}

func LoggerInit() {
	logger.InitLogger(viper.GetString("logger.runtimeLogFile"), viper.GetString("logger.level"), viper.GetString("logger.format"), true)
}
