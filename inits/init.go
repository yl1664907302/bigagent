package inits

import (
	"bigagent/internal/config/global"
	"bigagent/internal/register"
	"bigagent/internal/scrape/machine"
	"bigagent/internal/scrape/osquery"
	"bigagent/internal/strategy"
	utils "bigagent/internal/util"
	"bigagent/internal/util/crontab"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
	"time"
)

var (
	cmdbPattern = regexp.MustCompile(`grpc_cmdb(\d+)_stand(\d+)`)
	//cmdbPattern = regexp.MustCompile(`(\w+)`)
)

// Hander 启动http服务
func Hander(port string) {
	StandRouterGroupApp.StandRouter()
	StandRouterGroupApp2.StandRouter()
	SysRouterGroupApp.SysRouter()
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*
   关于AgentRegister
   参数 host，填充的是http推送端口，目前未启用，如下仅为占位符
   参数 grpc_host，读取的是配置文中的grpc服务器地址，由server端发送的配置进行热加载
   参数 openpush，是否开启推送
   参数 onlypush，是否只开启推送

   此外，在进行agent注册的时候，每种标准数据类型的api功能只会注册一次
*/

// AgentRegister agent注册
func AgentRegister() {
	strategy.Agents = nil
	//注册server端
	register.Stand1Register(global.V.GetString("system.grpc_server"), global.V.GetString("system.serct"), true, false)
	//注册api功能
	if global.V.GetString("system.api") == "1" {
		register.Stand2Register("占位符，只开启api", global.V.GetString("system.serct"), false, false)
	}
	//注册维易cmdb端
	register.VeopsRegister(global.V.GetString("veops.address"), true, true)
	//自动注册cmdb端
	configs := global.V.AllSettings()
	for key, value := range configs {
		matches := cmdbPattern.FindStringSubmatch(key)
		if len(matches) == 3 {
			standNum := matches[2]
			// 根据stand序号选择对应的注册函数
			switch standNum {
			case "1":
				for k, v := range configs {
					if k == key+"_token" {
						register.Stand1Register(value.(string), v.(string), true, false)
					}
				}
			case "2":
				for k, v := range configs {
					if k == key+"_token" {
						register.Stand2Register(value.(string), v.(string), true, false)
					}
				}
			case "3":
				//register.Stand3Register("127.0.0.1:8080", value.(string), true, false)
			// 可以继续添加更多的 case 以支持更多的 stand类型
			default:
				utils.DefaultLogger.Error("未识别的标准数据类型 序号: %s", standNum)
			}
		}
	}
}

// Crontab 执行定时任务
func Crontab() {
	crontab.ScrapeCrontab()
}

// ListerChannel 监听channel
func ListerChannel() {
	go func() {
		// 使用 select 实现更好的通道处理
		for {
			select {
			case signal := <-machine.MachineCh:
				if !signal {
					continue
				}
				// 使用非阻塞方式发送 false
				select {
				case machine.MachineCh <- false:
				default:
				}
				//utils.DefaultLogger.Info("数据更新，执行推送")

				// 异步执行注册和推送
				go func() {
					AgentRegister()

					agents := strategy.Agents
					if len(agents) == 0 {
						utils.DefaultLogger.Warn("strategy.Agents 为空")
						global.ASTATUS = "无可用agent策略"
						return
					}

					// 并发执行推送
					errChan := make(chan error, len(agents))
					for i, agent := range agents {
						go func(index int, a strategy.Agent) {
							if err := a.ExecutePush(); err != nil {
								utils.DefaultLogger.Errorf("agent策略序号：%d,数据推送异常: %s", index, err)
								errChan <- err
							} else {
								errChan <- nil
							}
						}(i, agent)
					}

					// 收集错误结果
					var e_len int
					for i := 0; i < len(agents); i++ {
						if err := <-errChan; err != nil {
							global.ASTATUS = "部分数据推送异常"
							e_len++
						}
					}
					if e_len == 1 {
						global.ASTATUS = "数据推送成功"
					}
				}()
			}
		}
	}()
}

func LoggerInit() {
	utils.InitLogger(global.V.GetString("system.logfile"), "info", "json", true)
}

// InstallIfNotExists 必要软件包预检
func InstallIfNotExists(pkgs []string) {
	if runtime.GOOS == "windows" {
		utils.DefaultLogger.Info("当前为 Windows 系统，跳过安装检测与安装逻辑。")
		return
	}
	for _, pkg := range pkgs {
		if isInstalled(pkg) {
			utils.DefaultLogger.Info("已安装：%s，跳过安装。\n", pkg)
		} else {
			err := InstallSomeThing(pkg)
			if err != nil {
				utils.DefaultLogger.Info("❌ 安装 %s 失败：%v\n", pkg, err)
				panic("安装失败，请检查软件包是否存在或网络连接是否正常。")
			} else {
				utils.DefaultLogger.Info("✅ 成功安装：%s\n", pkg)
			}
		}
	}
}

// 从URL中提取rpm包名
func extractRPMName(url string) string {
	// 获取URL的最后一部分
	parts := strings.Split(url, "/")
	filename := parts[len(parts)-1]

	// 移除.rpm扩展名和版本号等
	name := filename
	if idx := strings.LastIndex(filename, ".rpm"); idx > 0 {
		name = filename[:idx]
	}
	// 移除版本号（形如 -1.2.3-1.el7）
	if idx := strings.LastIndex(name, "-"); idx > 0 {
		name = name[:idx]
	}
	return name
}

// 检查包是否已安装
func isInstalled(pkg string) bool {
	// 如果是远程包，提取包名
	if strings.Contains(pkg, "http") {
		pkg = extractRPMName(pkg)
	}

	// 使用rpm -q检查包是否安装
	cmd := exec.Command("rpm", "-q", pkg)
	err := cmd.Run()
	return err == nil
}

// InstallSomeThing 安装指定的包（使用 yum 命令）
func InstallSomeThing(pkg string) error {
	utils.DefaultLogger.Info("正在安装 %s...\n", pkg)

	cmd := exec.Command("yum", "install", "-y", pkg)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// 初始化osquery客户端
func InitOsqueryClient() {
	queried, err := osquery.NewQuerier(global.V.GetString("system.empath"), 10*time.Second)
	if err != nil {
		utils.DefaultLogger.Error("osqueryd套接字连接失败", err)
		panic("osqueryd套接字连接失败")
	}
	osquery.OQry = queried
}
