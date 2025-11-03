# 启动应用完整指南

## 🚀 快速启动（推荐 - 无需安装任何工具）

### 方式一：使用 Go 命令直接运行（最简单）

```powershell
# 进入项目目录
cd D:\gowork\src\ShunFengParcel\ShunFengParcel

# 直接运行应用
go run ./cmd/ShunFengParcel
```

**优点**：
- ✅ 无需安装额外工具
- ✅ 自动编译和运行
- ✅ 代码修改后无需重新编译
- ✅ 立即生效

---

## 问题排查

### 问题1: `kratos run` 命令不识别

**原因**：kratos 工具未安装或未在 PATH 中

**解决方案**：

#### 安装 kratos CLI（可选）
```powershell
# 使用 go install 安装 kratos
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# 验证安装
kratos --version
```

#### 或者直接使用 go run（推荐）
```powershell
go run ./cmd/ShunFengParcel
```

---

## 所有启动方式对比

| 方式 | 命令 | 需要安装 | 推荐度 |
|------|------|---------|--------|
| Go 直接运行 | `go run ./cmd/ShunFengParcel` | ❌ 无 | ⭐⭐⭐⭐⭐ |
| 编译后运行 | `go build ./cmd/ShunFengParcel && ./ShunFengParcel` | ❌ 无 | ⭐⭐⭐⭐ |
| Kratos 运行 | `kratos run` | ✅ 需要 | ⭐⭐⭐ |
| Kratos 创建项目 | `kratos new` | ✅ 需要 | ⭐⭐ |

---

## 详细启动步骤

### 步骤1：进入项目目录
```powershell
cd D:\gowork\src\ShunFengParcel\ShunFengParcel
```

### 步骤2：运行应用
```powershell
go run ./cmd/ShunFengParcel
```

### 步骤3：验证应用是否启动
应该能看到类似的日志输出：
```
ts=2025-10-30T20:13:00.000+08:00 level=info caller=... HTTP server started on...
ts=2025-10-30T20:13:00.000+08:00 level=info caller=... gRPC server started on...
```

---

## 编译后运行（适合生产环境）

### 步骤1：编译应用
```powershell
cd D:\gowork\src\ShunFengParcel\ShunFengParcel

# 编译为可执行文件
go build -o ShunFengParcel.exe ./cmd/ShunFengParcel
```

### 步骤2：运行可执行文件
```powershell
# 直接运行编译好的文件
.\ShunFengParcel.exe

# 或在后台运行（需要保持终端打开）
Start-Process -NoNewWindow .\ShunFengParcel.exe
```

### 步骤3：停止应用
```powershell
# 查找进程image.png
Get-Process ShunFengParcel

# 停止进程
Stop-Process -Name ShunFengParcel -Force
```

---

## 常见错误及解决方案

### 错误1: "找不到 head 命令"
```
head : 无法将"head"项识别为 cmdlet、函数、脚本文件或可运行程序的名称
```

**原因**：PowerShell 中没有 `head` 命令（Linux 命令）

**解决方案**：
```powershell
# ❌ 错误做法
go run ./cmd/ShunFengParcel 2>&1 | head -30

# ✅ 正确做法
go run ./cmd/ShunFengParcel

# 或使用 PowerShell 命令
go run ./cmd/ShunFengParcel 2>&1 | Select-Object -First 30
```

### 错误2: "找不到模块"
```
cannot find module for path
```

**解决方案**：
```powershell
# 进入正确的项目目录
cd D:\gowork\src\ShunFengParcel\ShunFengParcel

# 下载依赖
go mod download

# 整理依赖
go mod tidy

# 运行应用
go run ./cmd/ShunFengParcel
```

### 错误3: "端口被占用"
```
bind: address already in use
```

**解决方案**：
```powershell
# 查看占用 8080 端口的进程
netstat -ano | findstr "8080"

# 或查看所有 Go 应用
Get-Process | Where-Object {$_.Name -like "*ShunFengParcel*"}

# 停止进程（使用 PID 或名称）
Stop-Process -Id <PID> -Force
# 或
Stop-Process -Name ShunFengParcel -Force
```

---

## 环境变量设置（可选）

### 消除 CPU quota 警告
```powershell
# 方式1：临时设置（仅当前会话）
$env:GOMAXPROCS = "16"
go run ./cmd/ShunFengParcel

# 方式2：永久设置（需要重启）
[Environment]::SetEnvironmentVariable("GOMAXPROCS", "16", "User")
```

### 设置其他环境变量
```powershell
# 设置日志级别
$env:LOG_LEVEL = "info"

# 设置配置文件路径
$env:CONFIG_PATH = "../../configs"

# 然后运行应用
go run ./cmd/ShunFengParcel
```

---

## 验证应用状态

### 1. 检查 HTTP 服务
```powershell
# 使用 curl（需要安装 curl）
curl http://localhost:8080/health

# 或使用 PowerShell 内置命令
Invoke-WebRequest -Uri http://localhost:8080/health
```

### 2. 检查 gRPC 服务
```bash
grpcurl -plaintext localhost:9000 list
```

### 3. 查看应用日志
```powershell
# 直接查看终端输出（应用运行中）
# 日志会实时显示在终端中

# 保存日志到文件
go run ./cmd/ShunFengParcel 2>&1 | Tee-Object app.log
```

---

## 后台运行应用

### 方式1：使用 PowerShell 后台任务
```powershell
# 在后台启动应用
Start-Job -ScriptBlock {
    cd D:\gowork\src\ShunFengParcel\ShunFengParcel
    go run ./cmd/ShunFengParcel
}

# 查看后台任务
Get-Job

# 停止后台任务
Stop-Job -Name <JobName>
Remove-Job -Name <JobName>
```

### 方式2：使用 nohup（类 Unix 系统）
```bash
nohup go run ./cmd/ShunFengParcel > app.log 2>&1 &
```

### 方式3：编译后使用 Windows 任务计划
```powershell
# 创建一个 VBScript 或 PowerShell 脚本来无窗口运行应用
```

---

## 开发工作流

### 快速开发模式
```powershell
# 1. 启动应用
go run ./cmd/ShunFengParcel

# 2. 修改代码
# 编辑文件...

# 3. 停止应用 (Ctrl+C)

# 4. 重新运行
go run ./cmd/ShunFengParcel
```

### 使用热重载工具（可选）
```powershell
# 安装 air（自动重新编译）
go install github.com/cosmtrek/air@latest

# 使用 air 运行（代码修改时自动重新编译）
air
```

---

## 生产环境部署

### 交叉编译为其他平台
```powershell
# Windows 编译为 Linux
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o ShunFengParcel ./cmd/ShunFengParcel

# 恢复为 Windows
$env:GOOS = "windows"
$env:GOARCH = "amd64"
```

### Docker 部署
```dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o ShunFengParcel ./cmd/ShunFengParcel

EXPOSE 8080 9000

CMD ["./ShunFengParcel"]
```

---

## 快速参考命令

```powershell
# 最常用的命令
go run ./cmd/ShunFengParcel                    # 运行应用
go build ./cmd/ShunFengParcel                  # 编译应用
go test ./...                                   # 运行测试
go mod tidy                                     # 清理依赖
go mod download                                 # 下载依赖

# 调试命令
netstat -ano | findstr "8080"                  # 检查端口
Get-Process | Where-Object {$_.Name -like "*ShunFengParcel*"}  # 查看进程
taskkill /PID <PID> /F                         # 关闭进程
```

---

## 常见问题 FAQ

**Q: 为什么我的代码修改没有生效？**
A: 使用 `go run` 时会自动编译，确保文件已保存。如果使用编译的可执行文件，需要重新编译。

**Q: 应用启动很慢，怎么办？**
A: 第一次启动会下载依赖，比较慢。可以使用 `go mod download` 预下载依赖。

**Q: 如何调试应用？**
A: 可以使用 Delve 调试器：
```powershell
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug ./cmd/ShunFengParcel
```

**Q: 应用占用了很多 CPU/内存，怎么办？**
A: 可以使用 `pprof` 进行性能分析：
```powershell
import _ "net/http/pprof"
# 然后访问 http://localhost:6060/debug/pprof/
```

---

**最后更新**：2025-10-30  
**Go 版本**：1.21+  
**Kratos 版本**：v2.x
