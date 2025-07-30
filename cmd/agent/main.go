package main

import (
	"bigagent/inits"
	"bigagent/internal/config/global"
	"bigagent/internal/scrape/osquery"
	utils "bigagent/internal/util"
	"context"
	"os"
)

func init() {
	// 创建一个通道来接收系统信号
	sigs := make(chan os.Signal, 1)
	inits.InitCmd(sigs)
	inits.LoggerInit()
	inits.InstallIfNotExists([]string{"git", "wget", "curl", "gcc", "make", "jq", "https://pkg.osquery.io/rpm/osquery-5.11.0-1.linux.x86_64.rpm"})
	inits.AgentRegister()
	inits.InitOsqueryClient()
	inits.Crontab()
	inits.ListerChannel()
	utils.DefaultLogger.Info("当前代码版本为：", "20250730")
}

func main() {
	ctx := context.Background()
	query, err := osquery.OQry.Query(ctx, "select * from os_version")
	if err != nil {
		utils.DefaultLogger.Errorf("osquery查询失败: %v", err)
		return
	}
	utils.DefaultLogger.Infof("osquery查询结果: %v", query)
	inits.RunG()
	inits.Hander(global.V.GetString("system.addr"))
}
