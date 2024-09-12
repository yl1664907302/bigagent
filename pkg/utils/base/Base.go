package base

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	FullDir *string
	OsType  *string
)

// GetDir 获取执行文件所在的目录，返回执行文件所在目录的绝对路径
func GetDir() {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(ex)
	FullDir = &dir
	//fmt.Println(*FullDir)
}

func GetMachine() {
	data := runtime.GOOS
	OsType = &data
}
