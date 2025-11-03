# 顺丰快递支付系统 - 部署指南

## 应用启动

### 编译
```bash
cd ShunFengParcel
go build ./cmd/ShunFengParcel
```

### 运行方式

#### 方式一：使用 PowerShell 脚本（推荐）
```powershell
# 在项目根目录运行
.\run.ps1
```

#### 方式二：使用 Batch 脚本
```cmd
run.bat
```

#### 方式三：手动设置环境变量后运行
```powershell
# PowerShell
$env:GOMAXPROCS=16
kratos run

# 或
go run ./cmd/ShunFengParcel
```

#### 方式四：直接运行（默认）
```bash
kratos run
# 或
go run ./cmd/ShunFengParcel
```

## 已知问题与解决方案

### 1. 支付宝公钥未找到错误
**错误信息**：`ERROR msg=签名验证失败:alipay: alipay public key not found`

**解决方案**：
- ✅ 已在 `configs/config.yaml` 中添加 `PublicKey` 字段
- ✅ 已在 `internal/service/payment.go` 中加载公钥
- ✅ 已在 `utils/alipay.go` 中加载公钥

需要确保将实际的支付宝公钥替换到配置文件中。

### 2. maxprocs: Leaving GOMAXPROCS-16: CPU quota undefined
**警告级别**：INFO（不影响应用运行）

**原因**：应用在检测 CPU 配额时出现问题（通常在容器环境中）

**解决方案**：
- 这是 `go.uber.org/automaxprocs` 库的正常行为
- 不会影响应用功能
- 如需消除此警告，可：
  1. 在容器中明确设置 CPU 限制
  2. 或手动设置 `GOMAXPROCS` 环境变量

**推荐方式（最简单）**：
使用启动脚本自动处理，脚本会自动设置 `GOMAXPROCS` 环境变量：

```powershell
# PowerShell - 推荐
.\run.ps1

# 或 Batch
run.bat

# 或手动设置
$env:GOMAXPROCS=16; kratos run
```

**验证**：
启动应用后，如果不再看到 CPU quota 的警告，表示问题已解决。

## 应用配置

### 服务端口
- **HTTP**: `0.0.0.0:8080`
- **gRPC**: `0.0.0.0:9000`

### 数据库配置
- **MySQL**: `14.103.136.136:3306`
- **Redis**: `14.103.136.136:6379`

配置文件位置：`configs/config.yaml`

## 支付宝配置清单

- [ ] 从支付宝开放平台获取应用公钥
- [ ] 将公钥添加到 `config.yaml` 的 `Alipay.PublicKey` 字段
- [ ] 验证私钥（Akey）正确性
- [ ] 验证应用 ID（AppId）正确性
- [ ] 测试支付流程确认签名验证成功

## 日志输出说明

应用启动时如果看到以下日志，表示支付宝配置正确：
```
初始化支付宝客户端成功
加载支付宝公钥成功
```

如果看到错误日志：
```
加载支付宝公钥失败: alipay public key not found
```

需要检查：
1. 公钥是否正确复制到配置文件
2. 公钥格式是否正确（应为 Base64 编码字符串）
3. 公钥是否被截断或包含换行符

## 故障排查步骤

1. **检查编译**
   ```bash
   go build ./cmd/ShunFengParcel
   ```

2. **检查配置文件**
   ```bash
   cat configs/config.yaml
   ```

3. **查看应用日志**
   ```bash
   go run ./cmd/ShunFengParcel 2>&1 | tee app.log
   ```

4. **验证 HTTP 服务**
   ```bash
   curl http://localhost:8080/health
   ```

5. **验证 gRPC 服务**
   ```bash
   grpcurl -plaintext localhost:9000 list
   ```

## 生产环境建议

1. **密钥安全**
   - 不要在版本控制中提交密钥
   - 使用环境变量存储敏感信息
   - 定期轮换密钥

2. **日志管理**
   - 配置日志级别为 INFO 或 WARN
   - 使用日志聚合系统（ELK、Splunk 等）
   - 定期清理日志文件

3. **监控告警**
   - 监控 HTTP/gRPC 服务可用性
   - 监控数据库连接状态
   - 监控支付宝回调验证失败率

4. **容器部署**
   - 设置合理的资源限制
   - 配置优雅关闭超时
   - 实现健康检查端点

## 相关文件修改记录

| 文件 | 修改内容 | 日期 |
|------|--------|------|
| `configs/config.yaml` | 添加 PublicKey 字段 | 2025-10-30 |
| `internal/service/payment.go` | 加载支付宝公钥 | 2025-10-30 |
| `utils/alipay.go` | 加载支付宝公钥 | 2025-10-30 |
| `ALIPAY_CONFIG.md` | 支付宝配置说明 | 2025-10-30 |
| `DEPLOYMENT_GUIDE.md` | 部署指南（本文档） | 2025-10-30 |

## 联系与支持

如遇到问题，请检查：
- [支付宝开放平台文档](https://open.alipay.com)
- [Kratos 框架文档](https://go-kratos.dev)
- [本项目 README](README.md)
