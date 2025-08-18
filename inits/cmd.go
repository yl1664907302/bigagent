package inits

import (
	utils "bigagent/internal/utils"
	"github.com/spf13/pflag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const pidFile = "agent.pid"

var runInfo = "接收到停止信号，正在优雅退出..."

func writePID() error {
	pid := os.Getpid()
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
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

func fileExists(filePath string) bool {
	// 检查文件是否存在
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false // 文件不存在
		}
		return false // 其他错误（如权限问题）
	}

	// 检查文件是否为空
	if info.Size() == 0 {
		return false // 文件存在但为空
	}

	return true // 文件存在且有内容
}

func Stop(re bool) {
	pidData, err := ioutil.ReadFile(pidFile)
	if err != nil {
		log.Printf("无法读取 PID 文件: %v", err)
	}

	pid := strings.TrimSpace(string(pidData)) // 去除空白字符
	if pid == "" {
		log.Printf("PID 文件为空")
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows 使用 taskkill 命令
		cmd = exec.Command("taskkill", "/F", "/PID", pid)
	} else {
		// 类 Unix 系统使用 kill 命令
		cmd = exec.Command("kill", "-9", pid)
	}

	log.Println(runInfo)
	time.Sleep(6 * time.Second)
	cmd.Run()

	if err := os.Remove(pidFile); err != nil {
		log.Printf("删除 PID 文件失败: %v", err)
	}
	if !re {
		os.Exit(0)
	}
}

// InitCmd 启动时初始化环境以及信号
func InitCmd(sigs chan os.Signal) {
	// 监听捕获信号
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// 解析命令行参数
	var daemon = pflag.BoolP("daemon", "d", false, "是否以守护进程方式运行")
	var env = pflag.StringP("server", "s", "start", "指定agent默认启动操作: start")
	var conf = pflag.StringP("config", "c", "", "指定agent配置文件路径: /path/config.yaml")
	// 设置命令行参数标准化兼容函数
	pflag.CommandLine.SetNormalizeFunc(wordSepNormalizeFunc)
	// 解析命令行参数
	pflag.Parse()

	// 规则校准
	if *conf != "" && !fileExists(*conf) {
		runInfo = "未指定配置文件或配置文件不存在！"
		log.Println(runInfo)
		os.Exit(0)
	}

	if *daemon && *env == "restart" {
		log.Println("无法在守护进程模式下执行重启操作！")
		os.Exit(0)
	}

	if *daemon && *env == "stop" {
		log.Println("无法在守护进程模式下执行停止操作！")
		os.Exit(0)
	}
	// 运行时处理信号
	go func() {
		<-sigs
		os.Remove(pidFile)
		utils.DefaultLogger.Info(runInfo) // 直接停止进程
		os.Exit(0)
	}()

	if *daemon {
		// 检查是否已有进程在运行
		if fileExists(pidFile) {
			log.Println("程序已在运行,请勿重复启动！")
			os.Exit(0)
		}
		args := []string{}
		for _, arg := range os.Args {
			if arg != "-d" {
				args = append(args, arg)
			}
		}
		cmd := exec.Command(os.Args[0], args...)
		if err := cmd.Start(); err != nil {
			log.Fatalf("Failed to start daemon: %v", err)
		}
		//log.Println("已进入后台运行，PID为:", cmd.Process.Pid)
		os.Exit(0)
	} else {

		if len(*env) == 0 && len(*conf) == 0 {
			log.Println("必须提供操作参数,请输入--help查看帮助")
			log.Println(runInfo)
			os.Exit(0)
			return
		}
		if len(*env) != 0 {
			switch *env {
			case "start":
				if len(*conf) != 0 {
					if fileExists(*conf) {
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
				} else {
					Viper("config.yml")
				}
				break
			case "stop":
				Stop(false)
				os.Exit(0)
			case "restart":
				Stop(true)
				args := []string{}
				args = append(args, "-s", "start")
				cmd := exec.Command(os.Args[0], args...)
				if err := cmd.Start(); err != nil {
					log.Fatalf("Failed to start daemon: %v", err)
				}
				os.Exit(0)
			default:
				log.Println("无效的操作参数")
				log.Println(runInfo)
				os.Exit(0)
				return
			}
		}

		if *env != "stop" {
			err := writePID()
			if err != nil {
				utils.DefaultLogger.Error("无法写入 PID 文件:", err)
				return
			}
		}
	}
}
