package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
)

var DefaultLogger = logrus.New()

func InitLogger(logFile, logLevel, logFormat string, openConsole bool) {
	//确保存放日志文件的目录始终存在
	logPathDir := path.Dir(logFile)                              //返回路径中除去最后一个元素的剩余部分，也就是路径最后一个元素所在的目录
	if err := os.MkdirAll(logPathDir, os.ModePerm); err != nil { //创建目录类似于（mkdir -p /aaa/bbb的效果）
		fmt.Println("创建日志目录失败：", err)
		os.Exit(3)
	}

	//指定日志输出方式
	outConsole := os.Stdout                                                       //标准输出
	outFile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644) //写入到文件中；此处使用追加写入，每一次都是追加写入。
	if err != nil {
		fmt.Println("打开日志文件失败: ", err)
	}

	//设定日志输出位置
	if openConsole {
		DefaultLogger.SetOutput(io.MultiWriter(outConsole, outFile))
	} else {
		DefaultLogger.SetOutput(io.MultiWriter(outFile))
	}

	//设定输出日志中是否要携带上文件名与行号
	DefaultLogger.SetReportCaller(false)

	//设定日志等级
	setLogLevel(DefaultLogger, logLevel)

	//设定日志输出格式
	setLogFormat(logFormat)

	//如果开启es的日志投放功能，则加载对应的钩子
	//if viper.GetBool("l.log2es") {
	//	l.AddHook(es.NewRuntimeToEsHook())
	//}

	DefaultLogger.Info("log模块初始化完成...")
}

// 设定日志等级
func setLogLevel(l *logrus.Logger, level string) {
	switch level {
	case "debug":
		l.SetLevel(logrus.DebugLevel)

	case "info":
		l.SetLevel(logrus.InfoLevel)

	default:
		l.SetLevel(logrus.DebugLevel)
	}
	/*
	   logrus有一个日志级别，默认的级别为InfoLevel。如果为了能看到Trace和Debug日志，需要在设置日志级别为TraceLevel。
	   DefaultLogger.Trace("trace msg") //很细粒度的信息，一般用不到
	   DefaultLogger.Debug("debug msg") //一般程序中输出的调试信息
	   DefaultLogger.Info("info msg")   //关键操作，核心流程的日志
	   DefaultLogger.Warn("warn msg")   //警告信息，提醒程序员注意
	   DefaultLogger.Error("error msg") //错误日志，需要查看原因
	   DefaultLogger.Panic("panic msg") //记录日志，然后panic
	   DefaultLogger.Fatal("fatal msg") //致命错误，出现错误时，程序无法正常运转，输出日志后，程序退出
	*/
}

// 设定日志格式
func setLogFormat(format string) {
	//指定日志输出格式
	switch format {
	case "json":
		//设置日志的输出格式（json格式）
		DefaultLogger.SetFormatter(
			&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05.000 -0700 MST",
			},
		)
	case "text":
		//设置日志的输出格式（Text格式）
		DefaultLogger.SetFormatter(
			&logrus.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05.000 -0700 MST",
			},
		)
	default:
		//设置日志的输出格式（Text格式）
		DefaultLogger.SetFormatter(
			&logrus.TextFormatter{
				TimestampFormat: "2006-01-02 15:04:05.000 -0700 MST",
			},
		)
	}
}
