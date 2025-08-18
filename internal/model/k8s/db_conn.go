package model

import (
	"bigagent/internal/config/global"
	"fmt"
	"time"
	"xorm.io/xorm"
	xlog "xorm.io/xorm/log"
)

// 根据配置文件中的mysql info 去连接
// 每一种orm 都有独特的init方法

// 定义1个全局变量 类型是engine 指针
// 其他模块使用 models.DB.xxx
var (
	DB *xorm.Engine
)

func InitDB(conf *global.MySQLConf) error {
	db, err := xorm.NewEngine("mysql", conf.Addr)
	if err != nil {
		fmt.Printf("[init.mysql.error][cannot connect to mysql][addr:%v][err:%v]\n", conf.Addr, err)
		return err
	}
	db.SetMaxIdleConns(conf.Idle)
	db.SetMaxOpenConns(conf.Max)
	db.SetConnMaxLifetime(time.Hour)
	db.ShowSQL(conf.Debug)
	db.Logger().SetLevel(xlog.LOG_INFO)
	DB = db
	return nil
}
