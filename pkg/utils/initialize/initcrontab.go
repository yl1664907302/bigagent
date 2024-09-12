package initialize

import (
	"cmdb_mini_agent/pkg/utils/crontab"
	"cmdb_mini_agent/pkg/utils/loggers"
)

// InitCrontab 初始化定时性计划任务
func InitCrontab() {
	loggers.DefaultLogger.Info("定时性计划任务初始化成功...")
	crontab.InitCrontab()
}
