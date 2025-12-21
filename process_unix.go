//go:build !windows
// +build !windows

package main

import (
	"os/exec"
	"syscall"
)

// setProcessGroup 设置进程组（Unix 系统）
func setProcessGroup(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}
