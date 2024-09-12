package main

import (
	"cmdb_mini_agent/apps"
	"cmdb_mini_agent/pkg/utils/initialize"
	"cmdb_mini_agent/pkg/utils/loggers"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	//	初始化基本信息
	initialize.InitBase()
	//	初始化配置文件
	initialize.InitConfig()
	//	初始化环境变量
	initialize.InitCmd()
	//	初始化日志
	loggers.Init()
	//	初始化CMDB
	initialize.InitCMDB()
	//	初始化定时性计划任务
	initialize.InitCrontab()
}

func main() {
	loggers.DefaultLogger.Info("当前代码版本为：", "20240730")
	//	程序每次运行之后都会先获取一次信息
	// 创建唯一的UUID
	apps.CreateMachineUUID()
	//	更新机器的基本信息
	apps.UpdateMachineData()

	// 创建一个通道来接收系统信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞主线程，直到接收到终止信号
	<-sigs

	// 在这里可以添加清理代码
	loggers.DefaultLogger.Info("程序正在优雅退出...")

}
