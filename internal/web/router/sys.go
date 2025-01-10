package router

import (
	utils "bigagent/internal/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type SysRouter struct{}

// 定义命令请求结构体
type CommandRequest struct {
	Command string `json:"command"` // 命令类型（如 "start", "stop"）
	Params  string `json:"params"`  // 命令参数（可选）
}

func (r *SysRouter) Cmd(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "仅支持 POST 请求", http.StatusMethodNotAllowed)
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
		Stop()
	default:
		utils.DefaultLogger.Error("无效的命令: %s", cmdReq.Command)
		http.Error(w, "无效的命令", http.StatusBadRequest)
		return
	}
}

var SysRouterApp = &SysRouter{}

func Stop() error {
	pidData, err := ioutil.ReadFile("agent.pid")
	if err != nil {
		return fmt.Errorf("无法读取 PID 文件: %v", err)
	}
	pid := string(pidData)
	if pid == "" {
		return fmt.Errorf("PID 文件为空")
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("taskkill", "/PID", pid, "/F")
	case "linux", "darwin":
		cmd = exec.Command("kill", "-9", pid)
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止进程失败: %v", err)
	}

	if err := os.Remove("agent.pid"); err != nil {
		return fmt.Errorf("删除 PID 文件失败: %v", err)
	}

	return nil
}
