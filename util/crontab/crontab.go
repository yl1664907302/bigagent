package crontab

import (
	"bigagent/scrape/machine"
	"bigagent/util/log"
	"github.com/robfig/cron/v3"
)

// CronTask Crontab执行的任务列表
func cronTask() {
	//开始采集
	machine.Ma = machine.NewMachine()
}

// ScrapeCrontab 初始化采集crontab任务
func ScrapeCrontab() {
	crontabRule := "@every 900s"
	c := cron.New()
	c.Start()

	addFunc, err := c.AddFunc(crontabRule, cronTask)
	if err != nil {
		log.DefaultLogger.Error("定时任务启动异常：", err)
		return
	}
	log.DefaultLogger.Info("定时任务启动成功,EntryID：", addFunc)
}
