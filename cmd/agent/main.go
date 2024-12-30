package main

import (
	"bigagent/inits"
	"bigagent/internal/config/global"
	utils "bigagent/internal/util"
	"os"
)

func init() {
	// 创建一个通道来接收系统信号
	sigs := make(chan os.Signal, 1)
	inits.InitCmd(sigs)
	inits.LoggerInit()
	utils.DefaultLogger.Info("当前代码版本为：", "20241230")
	inits.Crontab()
	inits.AgentRegister()
	inits.ListerChannel()
}

func main() {
	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
