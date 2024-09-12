package crontab

import (
	"cmdb_mini_agent/apps"
	"cmdb_mini_agent/pkg/utils/loggers"
	"github.com/robfig/cron/v3"
)

// CronTask Crontab执行的任务列表
func CronTask() {
	// 创建唯一的UUID
	apps.CreateMachineUUID()
	//	更新机器的基本信息
	apps.UpdateMachineData()
}

// InitCrontab 初始化一个crontab任务
func InitCrontab() {
	crontabRule := "@every 900s"
	c := cron.New()
	c.Start()

	addFunc, err := c.AddFunc(crontabRule, CronTask)
	if err != nil {
		loggers.DefaultLogger.Error("定时任务启动异常：", err)
		return
	}
	loggers.DefaultLogger.Info("定时任务启动成功,EntryID：", addFunc)
}
