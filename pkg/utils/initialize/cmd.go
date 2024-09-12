package initialize

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"strings"
)

// 标准化参数名称（容错：下划线、点号、减号、等等）
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return pflag.NormalizedName(name)
}

// InitCmd 初始化环境
func InitCmd() {

	/*
	 * 命令行参数（启动时可以通过 --xxx=aaa 的方式来调用，优先级最高，可以覆盖config/default.yaml中的变量）
	 */
	var machineType = pflag.StringP("machine", "m", "", "指定主机类型（virtual、physical、hwy、aly）")
	var env = pflag.StringP("env", "e", "", "指定主机所属的环境信息（offline、prod）")

	//设置命令行参数标准化兼容函数（防止用户手滑填写错误，可以兼容：下划线、点号、减号、等等）
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)

	//解析命令行参数
	pflag.Parse()
	if len(*machineType) != 0 {
		//判断--machine是否为空，如果不为空则覆盖到viper中
		//因为main.go的init中先执行InitConfig，后执行InitEnv，随后才陆续启动其它组件。
		//因此在这里拿到pflag参数覆盖到viper中时；后续调用viper中的变量，就是已经被pflag覆盖后的值。
		viper.Set("global.machineType", *machineType)
	}
	if len(*env) != 0 {
		//判断--machine是否为空，如果不为空则覆盖到viper中
		//因为main.go的init中先执行InitConfig，后执行InitEnv，随后才陆续启动其它组件。
		//因此在这里拿到pflag参数覆盖到viper中时；后续调用viper中的变量，就是已经被pflag覆盖后的值。
		viper.Set("global.env", *env)
	}
}
