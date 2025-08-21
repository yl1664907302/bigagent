package inits

import (
	"bigagent/internal/config/global"
	model "bigagent/internal/model/k8s"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB() error {
	conf := &global.MySQLConf{
		Name:  global.V.GetString("mysql_info.name"),
		Addr:  global.V.GetString("mysql_info.addr"),
		Max:   global.V.GetInt("mysql_info.max"),
		Idle:  global.V.GetInt("mysql_info.idel"), // 注意：config.yml 中是 idel
		Debug: global.V.GetBool("mysql_info.debug"),
	}
	if conf.Addr == "" {
		return fmt.Errorf("mysql dsn 为空: mysql_info.addr")
	}
	return model.InitDB(conf)
}
