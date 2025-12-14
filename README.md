# ShellGo - Web Hook 服务

一个基于 Go 的轻量级 Web Hook 服务，通过 HTTP 请求触发执行对应的 Shell 脚本。

## 功能特性

- 接收 HTTP GET 请求，通过 `tag` 参数指定要执行的脚本
- 自动执行对应的 Shell 脚本并返回执行结果
- 记录详细的执行日志，包括执行时间、结果状态和输出内容
- 支持日志分类存储，按 tag 和日期组织日志文件

## 编译说明

### 单平台编译

在当前系统平台编译可执行文件：

```bash
# 编译为当前平台的可执行文件
go build -o shellgo main.go

# 或者使用默认名称（生成 shellgo 或 shellgo.exe）
go build main.go
```

编译后的可执行文件可以直接运行：

```bash
# Linux/macOS
./shellgo

# Windows
shellgo.exe
```

### 多平台交叉编译

项目提供了 `build.sh` 脚本，可以一次性编译多个平台的可执行文件：

```bash
# 赋予执行权限（首次使用）
chmod +x build.sh

# 执行编译脚本
./build.sh
```

编译脚本会生成以下平台的可执行文件，并保存在 `bin/` 目录下：

- **Windows**: `shellgo_windows_amd64.exe`, `shellgo_windows_386.exe`, `shellgo_windows_arm64.exe`
- **Linux**: `shellgo_linux_amd64`, `shellgo_linux_386`, `shellgo_linux_arm64`, `shellgo_linux_arm`
- **macOS**: `shellgo_darwin_amd64`, `shellgo_darwin_arm64`

### 手动指定平台编译

如果需要手动指定目标平台进行交叉编译：

```bash
# Linux 64位
GOOS=linux GOARCH=amd64 go build -o shellgo_linux_amd64 main.go

# Windows 64位
GOOS=windows GOARCH=amd64 go build -o shellgo_windows_amd64.exe main.go

# macOS ARM64 (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o shellgo_darwin_arm64 main.go

# macOS Intel
GOOS=darwin GOARCH=amd64 go build -o shellgo_darwin_amd64 main.go
```

### 编译选项

可以使用以下编译选项优化可执行文件：

```bash
# 去除调试信息，减小文件大小
go build -ldflags="-s -w" -o shellgo main.go

# 静态链接，生成独立可执行文件（不依赖系统库）
CGO_ENABLED=0 go build -o shellgo main.go

# 组合使用：静态链接 + 去除调试信息
CGO_ENABLED=0 go build -ldflags="-s -w" -o shellgo main.go
```

## 使用方法

### 启动服务

开发环境直接运行：

```bash
go run main.go
```

或者使用编译后的可执行文件：

```bash
./shellgo
```

### 调用示例

```
http://localhost:8088?tag=deploy
```

这会查找并执行 `shell/deploy.sh` 脚本文件。

### 返回内容

- 脚本执行状态（成功/失败）
- 脚本输出内容
- 日志文件路径（总日志和详细日志）

## 脚本文件命名规则

脚本文件应放在 `shell/` 目录下，命名格式为：`{tag}.sh`

例如：

- `shell/deploy.sh` - 对应 `tag=deploy`
- `shell/backup.sh` - 对应 `tag=backup`
- `shell/test.sh` - 对应 `tag=test`

## 日志系统

服务会自动记录每次脚本执行的日志：

1. **总日志** (`{LOG_DIR}/{tag_dir}/{tag}.log`)
   - 格式：`时间|耗时|结果|MD5值`
   - 用于快速查看执行历史

2. **详细日志** (`{LOG_DIR}/{tag_dir}/detail/{日期}/{MD5值}.log`)
   - 存储完整的脚本输出内容
   - 相同输出的请求会共享同一个日志文件（通过 MD5 值判断）

## 环境变量

- `PORT`: 服务监听端口（默认: 8088）
- `LOG_DIR`: 日志存储目录（默认: logs）

## 配置说明

在项目根目录创建 `.env` 文件来配置环境变量：

```env
PORT=8088
LOG_DIR=logs
```
