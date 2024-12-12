package inits

import (
	"bigagent/config/global"
	"bigagent/register"
	"bigagent/scrape/machine"
	"bigagent/strategy"
	utils "bigagent/util"
	"bigagent/util/crontab"
	"log"
	"net/http"
	"regexp"
)

var (
	cmdbPattern = regexp.MustCompile(`grpc_cmdb(\d+)_stand(\d+)`)
)

// Hander 启动http服务
func Hander(port string) {
	StandRouterGroupApp.StandRouter()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*
   关于AgentRegister
   参数 host，填充的是http推送端口，目前未启用，如下仅为占位符
   参数 grpc_host，读取的是配置文中的grpc服务器地址，由server端发送的配置进行热加载
   参数 openpush，是否开启推送
   参数 onlypush，是否只开启推送

   此外，在进行agent注册的时候，每种标准数据类型的api功能只会注册一次
*/

// AgentRegister agent注册
func AgentRegister() {
	strategy.Agents = nil
	//注册server端
	register.Stand1Register("127.0.0.1:8080", global.V.GetString("system.grpc_server"), true, false)
	//注册cmdb端
	configs := global.V.AllSettings()
	for key, value := range configs {
		matches := cmdbPattern.FindStringSubmatch(key)
		if len(matches) == 3 {
			standNum := matches[2]
			// 根据stand序号选择对应的注册函数
			switch standNum {
			case "1":
				register.Stand1Register("127.0.0.1:8080", value.(string), true, false)
			case "2":
				//register.Stand2Register("127.0.0.1:8080", value.(string), true, false)
			case "3":
				//register.Stand3Register("127.0.0.1:8080", value.(string), true, false)
			// 可以继续添加更多的 case 以支持更多的 stand类型
			default:
				utils.DefaultLogger.Error("未识别的标准数据类型 序号: %s", standNum)
			}
		}
	}
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
				utils.DefaultLogger.Info("数据更新，执行推送")
				machine.MachineCh <- false
				//热加载
				AgentRegister()
				if strategy.Agents == nil || len(strategy.Agents) == 0 {
					utils.DefaultLogger.Warn("web.Agents is nil or empty")
					return
				}
				for i, agent := range strategy.Agents {
					err := agent.ExecutePush()
					if err != nil {
						utils.DefaultLogger.Errorf("agent策略序号：%d,数据推送异常: %s", i, err)
						global.ASTATUS = "部分数据推送异常"
					}
				}
				global.ASTATUS = "数据推送正常"
			}
		}
	}()
}

func LoggerInit() {
	utils.InitLogger(global.V.GetString("system.logfile"), "info", "json", true)
}
