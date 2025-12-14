# ShellGo - Web Hook 服务

一个基于 Go 的轻量级 Web Hook 服务，通过 HTTP 请求触发执行对应的 Shell 脚本。

## 功能特性

- 接收 HTTP GET 请求，通过 `tag` 参数指定要执行的脚本
- 自动执行对应的 Shell 脚本并返回执行结果
- 记录详细的执行日志，包括执行时间、结果状态和输出内容
- 支持日志分类存储，按 tag 和日期组织日志文件

## 使用方法

### 启动服务

```bash
go run main.go
```

或者编译后运行：

```bash
go build -o shellgo
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
