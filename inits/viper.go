package inits

import (
	"bigagent/internal/config/global"
	"bigagent/internal/utils"
	"github.com/spf13/viper"
)

func Viper(path string) {
	initConfig(path)
	go dynamicConfig()
}

func initConfig(path string) {
	// 创建新的Viper实例
	global.V = viper.New()

	// 设置配置文件路径和类型
	global.V.SetConfigFile(path)
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
	//global.V.OnConfigChange(func(event fsnotify.Event) {
	//	fmt.Printf("发现配置信息发生变化: %s\n", event.String())
	//})
}
