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
	utils.DefaultLogger.Info("当前代码版本为：", "20250728")
	inits.InstallIfNotExists([]string{"git", "wget", "curl", "gcc", "make", "jq", "https://pkg.osquery.io/rpm/osquery-5.11.0-1.linux.x86_64.rpm"})
	inits.AgentRegister()
	inits.InitOsqueryClient()
	inits.Crontab()
	inits.ListerChannel()
}

func main() {
	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
