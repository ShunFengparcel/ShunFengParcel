# 快速启动指南

## 🚀 最快的启动方式

### Windows PowerShell（推荐）
```powershell
# 方法1：使用启动脚本
.\run.ps1

# 方法2：手动启动（无警告）
$env:GOMAXPROCS=16; kratos run

# 方法3：简单启动（可能有CPU quota警告）
kratos run
```

### Windows 命令行
```cmd
run.bat
```

### Linux/Mac
```bash
export GOMAXPROCS=16
kratos run
```

## 📝 常见命令

| 命令 | 说明 |
|------|------|
| `kratos run` | 启动应用 |
| `go build ./cmd/ShunFengParcel` | 编译应用 |
| `go test ./...` | 运行测试 |
| `go mod tidy` | 清理依赖 |

## 🔍 验证应用是否运行正常

### 1. 检查 HTTP 服务
```bash
curl http://localhost:8080/health
```

### 2. 检查 gRPC 服务
```bash
grpcurl -plaintext localhost:9000 list
```

### 3. 查看日志
应用启动时应该能看到类似的日志：
```
ts=2025-10-30T20:08:50.000+08:00 level=info caller=... msg=HTTP server started on...
ts=2025-10-30T20:08:50.000+08:00 level=info caller=... msg=gRPC server started on...
```

## ⚠️ 常见问题

### Q: 看到 "CPU quota undefined" 警告怎么办？
A: 这不是错误，只是信息提示。使用 `$env:GOMAXPROCS=16` 可以消除。

### Q: 无法连接到数据库？
A: 检查 `config.yaml` 中的数据库配置是否正确：
```yaml
data:
  database:
    source: root:mysql_8tDfmn@tcp(14.103.136.136:3306)/shunfeng
```

### Q: 支付宝回调验证失败？
A: 确保在 `config.yaml` 中配置了正确的支付宝公钥：
```yaml
Alipay:
  PublicKey: "MIIBIj..."
```

## 📚 完整文档

- 详细配置说明：[DEPLOYMENT_GUIDE.md](DEPLOYMENT_GUIDE.md)
- 支付宝集成指南：[ALIPAY_CONFIG.md](ALIPAY_CONFIG.md)
- 项目 README：[README.md](ShunFengParcel/README.md)

## 💡 开发者提示

### 修改代码后重新启动
```bash
# 方法1：直接运行（自动重新编译）
kratos run

# 方法2：手动编译后运行
go build ./cmd/ShunFengParcel
./ShunFengParcel.exe
```

### 查看所有端口
```bash
netstat -ano | findstr "8080\|9000"
```

### 杀死进程
```bash
# 查找占用 8080 端口的进程
Get-Process | Where-Object {$_.Name -eq "ShunFengParcel"}

# 关闭进程
Stop-Process -Name ShunFengParcel -Force
```

---
**最后更新**：2025-10-30  
**应用版本**：v1.0.0  
**Go 版本**：1.21+
