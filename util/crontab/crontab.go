package crontab

import (
	"bigagent/scrape/machine"
	"bigagent/util/logger"
	"github.com/robfig/cron/v3"
)

// CronTask Crontab执行的任务列表
func cronTask() {
	//开始采集
	machine.Ma = machine.NewMachine()
	//更新通知
	machine.NotifyMachineAddressChange()
}

// ScrapeCrontab 初始化采集crontab任务
func ScrapeCrontab() {
	crontabRule := "@every 900s"
	c := cron.New()
	c.Start()

	addFunc, err := c.AddFunc(crontabRule, cronTask)
	if err != nil {
		logger.DefaultLogger.Error("定时任务启动异常：", err)
		return
	}
	logger.DefaultLogger.Info("定时任务启动成功,EntryID：", addFunc)
}
