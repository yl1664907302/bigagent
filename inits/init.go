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

// 用于描述一条注册计划
type registrationSpec struct {
	kind     string // stand1 | stand2 | veops | stand3(预留)
	host     string
	token    string
	openPush bool
	onlyPush bool
}

type regKey struct {
	kind     string
	host     string
	token    string
	openPush bool
	onlyPush bool
}

// 根据配置编译出所有需要执行的注册计划（显式 + 自动发现）
func compileRegistrationSpecs() []registrationSpec {
	specs := make([]registrationSpec, 0, 8)

	// 显式：server 端（stand1）
	specs = append(specs, registrationSpec{
		kind:     "stand1",
		host:     global.V.GetString("system.grpc_server"),
		token:    global.V.GetString("system.serct"),
		openPush: true,
		onlyPush: false,
	})

	// 显式：仅开启 API（stand2）
	if global.V.GetString("system.api") == "1" {
		specs = append(specs, registrationSpec{
			kind:     "stand2",
			host:     "占位符，只开启api",
			token:    global.V.GetString("system.serct"),
			openPush: false,
			onlyPush: false,
		})
	}

	// 显式：维易 cmdb
	specs = append(specs, registrationSpec{
		kind:     "veops",
		host:     global.V.GetString("veops.address"),
		openPush: true,
		onlyPush: true,
	})

	// 自动发现：grpc_cmdbX_standY 与其 _token
	configs := global.V.AllSettings()
	hostByBase := make(map[string]string)
	tokenByBase := make(map[string]string)
	for key, value := range configs {
		valStr, ok := value.(string)
		if !ok {
			continue
		}
		if matches := cmdbPattern.FindStringSubmatch(key); len(matches) == 3 {
			hostByBase[key] = valStr
			continue
		}
		if strings.HasSuffix(key, "_token") {
			base := strings.TrimSuffix(key, "_token")
			if matches := cmdbPattern.FindStringSubmatch(base); len(matches) == 3 {
				tokenByBase[base] = valStr
			}
		}
	}
	for base, host := range hostByBase {
		matches := cmdbPattern.FindStringSubmatch(base)
		if len(matches) != 3 {
			continue
		}
		standNum := matches[2]
		token := tokenByBase[base]
		switch standNum {
		case "1":
			specs = append(specs, registrationSpec{kind: "stand1", host: host, token: token, openPush: true, onlyPush: false})
		case "2":
			specs = append(specs, registrationSpec{kind: "stand2", host: host, token: token, openPush: true, onlyPush: false})
		case "3":
			// 预留 stand3
		default:
			utils.DefaultLogger.Error("未识别的标准数据类型 序号: %s", standNum)
		}
	}
	return specs
}

// 统一执行注册计划（包含去重）
func applyRegistrationSpecs(specs []registrationSpec) {
	seen := make(map[regKey]struct{}, len(specs))
	for _, s := range specs {
		k := regKey{s.kind, s.host, s.token, s.openPush, s.onlyPush}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}

		switch s.kind {
		case "stand1":
			register.Stand1Register(s.host, s.token, s.openPush, s.onlyPush)
		case "stand2":
			register.Stand2Register(s.host, s.token, s.openPush, s.onlyPush)
		case "veops":
			register.VeopsRegister(s.host, s.openPush, s.onlyPush)
		default:
			utils.DefaultLogger.Error("未知的注册类型: %s", s.kind)
		}
	}
}

// AgentRegister agent注册
func AgentRegister() {
	// 重置已注册的 agent 策略
	strategy.Agents = nil
	// 编译并统一执行注册计划
	specs := compileRegistrationSpecs()
	applyRegistrationSpecs(specs)
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
