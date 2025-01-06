package inits

import (
	utils "bigagent/internal/util"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/pflag"
)

const pidFile = "agent.pid"

var runInfo = "接收到停止信号，正在优雅退出..."

func writePID() error {
	pid := os.Getpid()
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}

func readPID() (int, error) {
	data, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(data))
}

func isRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return process.Signal(syscall.Signal(0)) == nil
}

// 标准化参数名称（容错：下划线、点号、减号、等等）
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return pflag.NormalizedName(name)
}

func sendSignal(pid int, sig syscall.Signal) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(sig)
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// InitCmd 启动时初始化环境以及信号
func InitCmd(sigs chan os.Signal) {
	// 监听捕获信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// 检查是否已有进程在运行
	pid, err := readPID()
	if err == nil && isRunning(pid) {
		utils.DefaultLogger.Error("程序已在运行，PID:", pid)
		return
	}

	// 写入当前进程的 PID
	if err := writePID(); err != nil {
		utils.DefaultLogger.Error("无法写入 PID 文件:", err)
		return
	}
	defer os.Remove(pidFile) // 确保在退出时删除 PID 文件
	// 解析命令行参数
	var env = pflag.StringP("server", "s", "", "指定agent默认启动操作: start")
	var conf = pflag.StringP("config", "c", "", "指定agent配置文件路径: /path/config.yaml")
	// 设置命令行参数标准化兼容函数
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	// 解析命令行参数
	pflag.Parse()

	// 运行时处理信号
	go func() {
		<-sigs
		utils.DefaultLogger.Info(runInfo) // 直接停止进程
		os.Exit(0)
	}()

	if len(*env) == 0 && len(*conf) == 0 {
		log.Println("必须提供操作参数,请输入--help查看帮助")
		log.Println(runInfo)
		os.Exit(0)
		return
	}
	if len(*env) != 0 {
		switch *env {
		case "start":
			Viper("config.yml")
			log.Println("启动完成")
			return
		default:
			log.Println("无效的操作参数")
			log.Println(runInfo)
			os.Exit(0)
			return
		}
	}

	if len(*conf) != 0 && fileExists(*conf) {
		switch *conf {
		default:
			Viper(*conf)
			log.Println("启动完成")
		}
	} else {
		runInfo = "未指定配置文件或配置文件不存在！"
		log.Println(runInfo)
		os.Exit(0)
	}
}
