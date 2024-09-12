package initialize

import (
	"cmdb_mini_agent/pkg/utils/CMDBClient"
	"cmdb_mini_agent/pkg/utils/loggers"
)

// InitCMDB	初始化CMDB客户端
func InitCMDB() {
	CMDBClient.CheckClient()
	loggers.DefaultLogger.Info("CMDB客户端初始化完成...")
}
