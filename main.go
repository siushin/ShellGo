package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

// 用于跟踪每个 tag 对应的正在执行的进程
var (
	runningProcesses = make(map[string]*exec.Cmd)
	processMutex     sync.Mutex
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Printf("警告: 无法加载 .env 文件: %v", err)
	}
	
	http.HandleFunc("/", handleHook)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8088"
	}
	
	log.Printf("服务器启动在端口 %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleHook(w http.ResponseWriter, r *http.Request) {
	// 记录请求开始时间
	startTime := time.Now()
	requestTime := startTime.Format("2006-01-02 15:04:05")
	
	// 获取 tag 参数
	tag := r.URL.Query().Get("tag")
	if tag == "" {
		http.Error(w, "缺少 tag 参数", http.StatusBadRequest)
		return
	}
	
	// 直接使用 tag 作为脚本文件名
	scriptName := fmt.Sprintf("%s.sh", tag)
	scriptPath := filepath.Join("shell", scriptName)
	
	// 检查文件是否存在
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		http.Error(w, fmt.Sprintf("脚本文件不存在: %s", scriptName), http.StatusNotFound)
		return
	}
	
	// 检查是否有相同 tag 的进程正在执行，如果有则终止
	processMutex.Lock()
	if oldCmd, exists := runningProcesses[tag]; exists {
		log.Printf("检测到相同 tag (%s) 的进程正在执行，正在终止旧进程...", tag)
		// 终止进程及其子进程
		if oldCmd.Process != nil && oldCmd.Process.Pid > 0 {
			pid := oldCmd.Process.Pid
			if runtime.GOOS == "windows" {
				// Windows 系统：直接终止进程
				oldCmd.Process.Kill()
			} else {
				// Unix 系统：使用负 PID 来终止整个进程组（包括子进程）
				// 先尝试优雅终止（SIGTERM）
				killCmd := exec.Command("kill", "-TERM", fmt.Sprintf("-%d", pid))
				killCmd.Run()
				// 等待一下让进程有机会优雅退出
				time.Sleep(300 * time.Millisecond)
				// 检查进程是否还在运行（使用 kill -0 检查进程组）
				checkCmd := exec.Command("kill", "-0", fmt.Sprintf("-%d", pid))
				if checkCmd.Run() == nil {
					// 进程组还在运行，强制杀死（SIGKILL）
					log.Printf("进程组仍在运行，强制终止...")
					killCmd = exec.Command("kill", "-KILL", fmt.Sprintf("-%d", pid))
					killCmd.Run()
				}
				// 也直接终止主进程（双重保险）
				oldCmd.Process.Kill()
			}
		}
		delete(runningProcesses, tag)
		log.Printf("旧进程已终止")
	}
	processMutex.Unlock()
	
	// 执行 shell 脚本
	cmd := exec.Command("sh", scriptPath)
	// 设置进程组，以便能够终止子进程（仅在 Unix 系统上）
	setProcessGroup(cmd)
	
	// 记录新进程
	processMutex.Lock()
	runningProcesses[tag] = cmd
	processMutex.Unlock()
	
	// 执行脚本
	output, err := cmd.CombinedOutput()
	
	// 执行完成后清理进程记录
	processMutex.Lock()
	delete(runningProcesses, tag)
	processMutex.Unlock()
	
	// 计算耗时
	duration := time.Since(startTime)
	
	// 计算输出内容的 MD5 值
	md5Hash := md5.Sum(output)
	md5Value := fmt.Sprintf("%x", md5Hash)
	
	// 获取日志目录
	logDir := os.Getenv("LOG_DIR")
	if logDir == "" {
		logDir = "logs" // 默认日志目录
	}
	
	// 将 tag 中的点转成下划线作为目录名
	tagDir := strings.ReplaceAll(tag, ".", "_")
	tagLogDir := filepath.Join(logDir, tagDir)
	
	// 确保 tag 日志目录存在
	if err := os.MkdirAll(tagLogDir, 0755); err != nil {
		log.Printf("创建 tag 日志目录失败: %v", err)
	}
	
	// 第一个日志：总的信息（tag.log）- 一行存储，用 | 分隔
	summaryLogPath := filepath.Join(tagLogDir, fmt.Sprintf("%s.log", tag))
	result := "成功"
	if err != nil {
		result = fmt.Sprintf("失败: %v", err)
	}
	summaryContent := fmt.Sprintf("%s|%v|%s|%s\n", 
		requestTime, 
		duration, 
		result,
		md5Value)
	
	// 追加模式写入
	file, err := os.OpenFile(summaryLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("打开总日志文件失败: %v", err)
	} else {
		if _, err := file.WriteString(summaryContent); err != nil {
			log.Printf("写入总日志文件失败: %v", err)
		}
		file.Close()
	}
	
	// 第二个日志：详细日志（detail/年-月-日/md5值.log）
	now := time.Now()
	dateDir := now.Format("2006-01-02")
	detailDir := filepath.Join(tagLogDir, "detail", dateDir)
	
	// 确保详细日志目录存在
	if err := os.MkdirAll(detailDir, 0755); err != nil {
		log.Printf("创建详细日志目录失败: %v", err)
	}
	
	detailLogPath := filepath.Join(detailDir, fmt.Sprintf("%s.log", md5Value))
	detailContent := string(output)
	
	if err := os.WriteFile(detailLogPath, []byte(detailContent), 0644); err != nil {
		log.Printf("写入详细日志文件失败: %v", err)
	}
	
	if err != nil {
		log.Printf("执行脚本失败: %v, 输出: %s", err, string(output))
		http.Error(w, fmt.Sprintf("执行脚本失败: %v\n输出: %s", err, string(output)), http.StatusInternalServerError)
		return
	}
	
	// 返回执行结果
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "脚本执行成功: %s\n\n输出:\n%s\n\n总日志: %s\n详细日志: %s", 
		scriptName, 
		string(output), 
		summaryLogPath,
		detailLogPath)
}

