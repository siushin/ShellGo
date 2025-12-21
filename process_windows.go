//go:build windows
// +build windows

package main

import "os/exec"

// setProcessGroup 设置进程组（Windows 系统，无需设置）
func setProcessGroup(cmd *exec.Cmd) {
	// Windows 不需要设置进程组
}
