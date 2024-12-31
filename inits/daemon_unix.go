//go:build !windows

package inits

import "syscall"

func daemonAttr() *syscall.SysProcAttr {
	return &syscall.SysProcAttr{Setsid: true}
}
