package router

import (
	utils "bigagent/internal/util"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

const pidFile = "agent.pid"

var runInfo = "接收到停止信号，正在优雅退出..."

type SysRouter struct{}

// 定义命令请求结构体
type CommandRequest struct {
	Command string `json:"command"` // 命令类型（如 "start", "stop"）
	Params  string `json:"params"`  // 命令参数（可选）
}

func (r *SysRouter) Cmd(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "仅支持POST请求", http.StatusMethodNotAllowed)
		return
	}
	var cmdReq CommandRequest
	err := json.NewDecoder(req.Body).Decode(&cmdReq)
	if err != nil {
		utils.DefaultLogger.Info("解析请求体失败: %v", err)
		http.Error(w, "无效的请求体", http.StatusBadRequest)
		return
	}

	if cmdReq.Command == "" {
		utils.DefaultLogger.Error("命令不能为空")
		http.Error(w, "命令不能为空", http.StatusBadRequest)
		return
	}

	utils.DefaultLogger.Info("接收到命令: %s", cmdReq.Command)
	switch cmdReq.Command {
	case "stop":
		if err := os.Remove(pidFile); err != nil {
			log.Printf("删除 PID 文件失败: %v", err)
		}
		os.Exit(0)
	case "restart":
		Restart()
	default:
		utils.DefaultLogger.Error("无效的命令: %s", cmdReq.Command)
		http.Error(w, "无效的命令", http.StatusBadRequest)
		return
	}
}

var SysRouterApp = &SysRouter{}

func Restart() {
	args := []string{
		"-s", "start", // 默认启动操作
		//"-c", "config.yml",
	}

	cmd := exec.Command(os.Args[0], args...)
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to restart: %v", err)
	}

	if err := writePID(); err != nil {
		log.Fatalf("无法写入 PID 文件: %v", err)
	}

	log.Println("重启成功，新进程 PID 为:", cmd.Process.Pid)
}

func writePID() error {
	pid := os.Getpid()
	return ioutil.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}
