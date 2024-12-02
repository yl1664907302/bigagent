package inits

import (
	"bigagent/config/global"
	utils "bigagent/util"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Viper() {
	initConfig()
	go dynamicConfig()
}

func initConfig() {
	// 创建新的Viper实例
	global.V = viper.New()

	// 设置配置文件路径和类型
	global.V.SetConfigFile("config.yml")
	global.V.SetConfigType("yaml")

	// 读取配置文件
	err := global.V.ReadInConfig()
	if err != nil {
		utils.DefaultLogger.Error("Error reading config file: " + err.Error())
		return
	}
}

// viper支持应用程序在运行中实时读取配置文件的能力。确保在调用 WatchConfig()之前添加所有的configPaths。
func dynamicConfig() {
	global.V.WatchConfig()
	global.V.OnConfigChange(func(event fsnotify.Event) {
		fmt.Printf("发现配置信息发生变化: %s\n", event.String())
	})
}
