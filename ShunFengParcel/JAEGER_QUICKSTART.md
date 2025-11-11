# Jaeger 链路追踪快速开始

## 🎯 已完成配置

✅ Jaeger 客户端已集成  
✅ 使用 OTLP HTTP 协议  
✅ 连接到服务器: `14.103.153.242:4318`  
✅ 自动追踪 HTTP/gRPC 请求  
✅ 日志包含 trace.id 和 span.id  

## 🚀 快速开始

### 1. 启动服务
```bash
cd ShunFengParcel
go run ./cmd/ShunFengParcel
```

看到以下日志说明连接成功：
```
========================================
🎉 Jaeger 链路追踪初始化成功！
   - 采样率: 100% (所有请求都会被追踪)
   - 查看链路: http://14.103.153.242:16686
   - 服务名称: ShunFengParcel
   - 测试 span: jaeger-connection-test
========================================
💡 提示: 等待 5-10 秒后在 Jaeger UI 中查看测试链路
```

如果看到 ❌ 错误信息，请查看 [验证指南](docs/jaeger-verification.md)

### 2. 发送测试请求
```bash
curl -X POST http://localhost:8000/admin/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'
```

### 3. 验证连接（重要！）
**等待 5-10 秒**，然后打开浏览器访问：
```
http://14.103.153.242:16686
```

在 Jaeger UI 中验证：
1. Service 下拉框选择：`ShunFengParcel`
2. Operation 下拉框选择：`jaeger-connection-test`
3. 点击 `Find Traces` 按钮
4. ✅ 如果看到测试链路，说明连接成功！

### 4. 查看真实请求链路
发送请求后，在 Jaeger UI 中：
1. Service 选择：`ShunFengParcel`
2. Operation 选择：`/admin/register` 或其他接口
3. 点击 `Find Traces`
4. 查看完整的调用链路和性能数据

## 🧪 自动化测试

### Windows (PowerShell)
```powershell
.\scripts\test-jaeger.ps1
```

### Linux/Mac (Bash)
```bash
chmod +x scripts/test-jaeger.sh
./scripts/test-jaeger.sh
```

## 📊 查看效果

在 Jaeger UI 中你可以看到：
- 请求的完整调用链路
- 每个步骤的耗时
- HTTP 状态码、方法、URL
- 错误信息（如果有）
- 服务依赖关系图

## 🔧 配置说明

### 当前配置
- **服务名**: ShunFengParcel
- **采样率**: 100% (AlwaysSample)
- **协议**: OTLP HTTP
- **端点**: 14.103.153.242:4318

### 修改配置
编辑 `cmd/ShunFengParcel/main.go` 中的 `initJaeger()` 函数：

```go
// 修改服务器地址
jaegerEndpoint := "your-server:4318"

// 修改采样率（生产环境建议 10-20%）
tracesdk.WithSampler(tracesdk.TraceIDRatioBased(0.1)) // 10% 采样
```

## 📖 详细文档

查看完整文档：[docs/jaeger-integration.md](docs/jaeger-integration.md)

## ⚠️ 注意事项

1. **生产环境**: 建议降低采样率到 10-20%
2. **性能影响**: 链路追踪会增加约 1-2% 的性能开销
3. **数据延迟**: 链路数据会有 5-10 秒的上报延迟
4. **网络要求**: 确保服务器能访问 Jaeger 服务器

## 🐛 故障排查

### 如何判断是否连接成功？
查看 [完整验证指南](docs/jaeger-verification.md)

### 快速检查清单
- ✅ 启动日志显示 "🎉 Jaeger 链路追踪初始化成功！"
- ✅ Jaeger UI 中能看到 "jaeger-connection-test" 链路
- ✅ 请求日志包含 trace.id 和 span.id
- ✅ Jaeger UI 中能搜索到对应的 trace

### 看不到链路数据？
1. **等待 5-10 秒**（数据上报有延迟）
2. 刷新 Jaeger UI 页面
3. 检查时间范围设置（建议 "Last 15 minutes"）
4. 确认 Service 选择的是 "ShunFengParcel"

### 连接失败？
```powershell
# Windows: 测试端口连通性
Test-NetConnection -ComputerName 14.103.153.242 -Port 4318

# 如果端口不通，检查：
# 1. 网络连接
# 2. 防火墙设置
# 3. Jaeger 服务器状态
```

## 📞 支持

如有问题，请查看：
- [Jaeger 官方文档](https://www.jaegertracing.io/docs/)
- [OpenTelemetry Go 文档](https://opentelemetry.io/docs/instrumentation/go/)
