package decision

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

type HostDoing struct{}

// Exec 执行任意宿主机命令，返回 stdout/stderr 合并输出
func (h *HostDoing) Exec(ctx context.Context, cmd []string) (string, error) {
	if len(cmd) == 0 {
		return "", fmt.Errorf("empty command")
	}
	return run(ctx, cmd)
}

// Service 跨平台服务控制抽象：start|stop|restart|status
// - Linux: 优先 systemctl，缺省回退 service
// - Windows: 使用 PowerShell Start/Stop/Restart/Get-Service
func (h *HostDoing) Service(ctx context.Context, name string, action string) (string, error) {
	action = strings.ToLower(strings.TrimSpace(action))
	switch action {
	case "start", "stop", "restart", "status":
	default:
		return "", fmt.Errorf("unsupported action: %s", action)
	}
	if strings.TrimSpace(name) == "" {
		return "", fmt.Errorf("service name is required")
	}

	if runtime.GOOS == "windows" {
		cmd := winServiceCmd(action, name)
		return run(ctx, cmd)
	}

	// Linux / Unix-like
	// 1) systemctl
	if _, err := exec.LookPath("systemctl"); err == nil {
		unit := normalizeUnitName(name)
		out, err := run(ctx, []string{"systemctl", action, unit})
		if err == nil {
			return out, nil
		}
		// 如果 systemctl 不可用或执行失败，继续尝试 service
		if !isNotFoundErr(err) && !strings.Contains(strings.ToLower(out), "systemctl: command not found") {
			return out, err
		}
	}

	// 2) service 回退
	if _, err := exec.LookPath("service"); err == nil {
		return run(ctx, []string{"service", name, action})
	}
	return "", fmt.Errorf("neither systemctl nor service found on this system")
}

func run(ctx context.Context, args []string) (string, error) {
	c := exec.CommandContext(ctx, args[0], args[1:]...)
	var buf bytes.Buffer
	c.Stdout = &buf
	c.Stderr = &buf
	err := c.Run()
	return buf.String(), err
}

func winServiceCmd(action, name string) []string {
	ps := []string{"powershell", "-NoProfile", "-NonInteractive", "-Command"}
	switch action {
	case "start":
		return append(ps, fmt.Sprintf("Start-Service -Name '%s'", name))
	case "stop":
		return append(ps, fmt.Sprintf("Stop-Service -Name '%s' -Force", name))
	case "restart":
		return append(ps, fmt.Sprintf("Restart-Service -Name '%s' -Force", name))
	case "status":
		return append(ps, fmt.Sprintf("(Get-Service -Name '%s').Status", name))
	default:
		return ps // 不会到达
	}
}

func normalizeUnitName(unit string) string {
	u := strings.TrimSpace(unit)
	if u == "" {
		return ""
	}
	if !strings.Contains(u, ".") {
		u += ".service"
	}
	return u
}

func isNotFoundErr(err error) bool {
	var ee *exec.Error
	if errors.As(err, &ee) {
		return true
	}
	es := strings.ToLower(err.Error())
	return strings.Contains(es, "executable file not found") ||
		strings.Contains(es, "file does not exist") ||
		strings.Contains(es, "not found")
}
