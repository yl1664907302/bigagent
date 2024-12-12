package global

import (
	"github.com/spf13/viper"
)

/*
	    V为全局变量，用于读取配置文件
		ASTATUS为全局变量，用于记录agent状态
*/

var (
	V             *viper.Viper
	ACTION_DETAIL string
	ASTATUS       string
)
