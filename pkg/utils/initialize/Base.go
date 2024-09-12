package initialize

import (
	"cmdb_mini_agent/pkg/utils/base"
	"cmdb_mini_agent/pkg/utils/loggers"
)

func InitBase() {
	base.GetDir()
	base.GetMachine()
	loggers.DefaultLogger.Info("基础信息获取完毕...")
}
