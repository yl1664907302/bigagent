package HostInfo

import (
	"cmdb_mini_agent/apps"
	"fmt"
	"os"
	"runtime"
)

func hostBasic() {
	fmt.Println("系统类型：", runtime.GOOS)
	fmt.Println("系统架构：", runtime.GOARCH)
	fmt.Println("CPU核数：", runtime.GOMAXPROCS(0))
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	println(`电脑名称：`, name)
	ip, err := apps.GetIp()
	if err != nil {
		panic(err)
	}
	fmt.Println("IP：", ip)
}
