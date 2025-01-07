package crontab

import (
	"bigagent/internal/config/global"
	"bigagent/internal/scrape/machine"
	"bigagent/internal/util"
	"github.com/robfig/cron/v3"
)

// CronTask Crontab执行的任务列表
func cronTask() {
	//开始采集
	machine.SmpMa = machine.NewSmpMachine()
	//更新通知
	machine.NotifySmpMachineAddressChange()
}

// ScrapeCrontab 初始化采集crontab任务
func ScrapeCrontab() {
	var crontabRule string
	if global.V.GetString("collection_frequency") == "" {
		crontabRule = "@every 3s"
	}
	crontabRule = "@every " + global.V.GetString("collection_frequency")
	c := cron.New()
	c.Start()

	addFunc, err := c.AddFunc(crontabRule, cronTask)
	if err != nil {
		utils.DefaultLogger.Error("定时任务启动异常：", err)
		return
	}
	utils.DefaultLogger.Info("定时任务启动成功,EntryID：", addFunc)
}
