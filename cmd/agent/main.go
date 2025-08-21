package main

import (
	"bigagent/inits"
	"bigagent/internal/config/global"
	"bigagent/internal/utils"
	"context"
	"os"
)

func init() {
	// 创建一个通道来接收系统信号
	ctx := context.Background()
	sigs := make(chan os.Signal, 1)
	inits.InitCmd(sigs)
	inits.LoggerInit()
	inits.InstallIfNotExists([]string{"git", "wget", "curl", "gcc", "make", "jq", "https://pkg.osquery.io/rpm/osquery-5.11.0-1.linux.x86_64.rpm"})
	inits.AgentRegister()
	err := inits.InitDB()
	if err != nil {
		utils.DefaultLogger.WithError(err).Error("数据库初始化失败")
	}
	//inits.InitOsqueryClient()
	inits.Crontab(ctx)
	inits.ListerChannel()
	utils.DefaultLogger.WithField("version", "20250730").Info("Agent 版本")
}

func main() {
	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
