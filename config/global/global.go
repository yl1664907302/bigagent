package global

import (
	"bigagent/config"
	"github.com/spf13/viper"
)

var (
	V    *viper.Viper
	CONF *config.Server // 全局配置
)
