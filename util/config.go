package utils

import (
	"path"
	"runtime"
)

var (
	RootPath string
)

func init() {
	RootPath = path.Dir(GetCurrentPath()+"..") + "/"
}

// 获取此函数的代码路径
func GetCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
