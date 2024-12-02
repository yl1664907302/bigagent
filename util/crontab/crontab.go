package crontab

import (
	"bigagent/scrape/machine"
	utils "bigagent/util"
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
	crontabRule := "@every 3s"
	c := cron.New()
	c.Start()

	addFunc, err := c.AddFunc(crontabRule, cronTask)
	if err != nil {
		utils.DefaultLogger.Error("定时任务启动异常：", err)
		return
	}
	utils.DefaultLogger.Info("定时任务启动成功,EntryID：", addFunc)
}
