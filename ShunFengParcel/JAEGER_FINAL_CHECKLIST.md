# ✅ Jaeger 集成最终检查清单

## 📋 已完成的修改

### 1. ✅ admin.go - 追踪代码实现
**文件**: `ShunFengParcel/internal/service/admin.go`

- [x] 在 `FindAdminByUsername` 方法开始时创建 Span
- [x] 使用 `defer span.End()` 确保 Span 正确结束
- [x] 添加 Span 属性 (username, operation)
- [x] 添加事件标记 (密码验证失败/成功)
- [x] 导入必要的包 (`go.opentelemetry.io/otel/attribute`)

### 2. ✅ main.go - TracerProvider 配置
**文件**: `ShunFengParcel/cmd/ShunFengParcel/main.go`

- [x] 创建 `initTracerProvider()` 函数
- [x] 配置 OTLP gRPC 导出器 (14.103.153.242:4317)
- [x] 设置服务名称为 "ShunFengParcel"
- [x] 配置采样策略 (AlwaysSample)
- [x] 在 main() 中调用初始化函数
- [x] 正确处理 shutdown 清理

### 3. ✅ http.go - HTTP Server 中间件
**文件**: `ShunFengParcel/internal/server/http.go`

- [x] 已配置 `tracing.Server()` 中间件
- [x] 中间件顺序正确 (recovery → tracing)

### 4. ✅ otel-collector-config.yaml - Collector 配置
**文件**: `ShunFengParcel/deploy/otel-collector-config.yaml`

- [x] 添加 Jaeger 导出器配置
- [x] 配置 Jaeger 端点 (14.103.153.242:14250)
- [x] 在 traces pipeline 中添加 jaeger 导出器

### 5. ✅ 编译和构建
- [x] 应用程序已成功编译
- [x] 可执行文件：`bin/ShunFengParcel.exe`

## ⚠️ 待执行操作

### 🔴 关键步骤（必须执行）

**在服务器 14.103.153.242 上重启 OpenTelemetry Collector**

```bash
# 1. SSH 到服务器
ssh user@14.103.153.242

# 2. 上传新配置
# 将 ShunFengParcel/deploy/otel-collector-config.yaml 上传到服务器

# 3. 重启 Collector (根据您的部署方式选择)

## Docker 方式：
docker stop otel-collector
docker rm otel-collector
docker run -d --name otel-collector \
  --restart=always \
  -p 4317:4317 -p 4318:4318 -p 14250:14250 \
  -v /path/to/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  otel/opentelemetry-collector-contrib:latest \
  --config=/etc/otel-collector-config.yaml

## Kubernetes 方式：
kubectl apply -f otel-collector-config.yaml
kubectl rollout restart deployment/otel-collector

## 系统服务方式：
sudo cp otel-collector-config.yaml /etc/otel-collector/config.yaml
sudo systemctl restart otel-collector

# 4. 检查日志
docker logs -f otel-collector  # Docker
kubectl logs -f deployment/otel-collector  # Kubernetes
sudo journalctl -u otel-collector -f  # 系统服务
```

## 🧪 测试流程

### 步骤 1：重启应用程序

```powershell
cd E:\gowork\src\ShunFengParcel\ShunFengParcel
Stop-Process -Name "ShunFengParcel" -Force
.\bin\ShunFengParcel.exe -conf .\configs
```

### 步骤 2：调用接口

```powershell
curl.exe -X POST http://localhost:8000/admin/findadminbyusername `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"test123"}'
```

### 步骤 3：查看 Jaeger UI

1. 打开浏览器访问：http://14.103.153.242:16686/search
2. 在 **Service** 下拉框中应该能看到 `ShunFengParcel`
3. 选择服务后点击 **Find Traces**
4. 查看追踪详情和 Span 信息

## 📊 预期结果

### Jaeger UI 中应该显示：

#### Service 列表
```
- jaeger-all-in-one
- ShunFengParcel  ← 新增的服务
```

#### Trace 详情
- **Service**: ShunFengParcel
- **Operation**: FindAdminByUsername
- **Duration**: ~10-50ms
- **Spans**: 1-2 个
- **Tags**:
  - service.name: ShunFengParcel
  - service.version: 1.0.0
  - username: testuser
  - operation: find_admin_by_username
- **Events/Logs**:
  - "密码验证失败" 或 "用户查询成功"

## 🔍 故障排查

### 问题：Jaeger UI 中看不到 ShunFengParcel 服务

#### 检查 1：Collector 是否收到数据
```bash
# 查看 Collector 日志
docker logs otel-collector | grep ShunFengParcel

# 或查看导出的文件
tail -f /var/log/otel-collector/traces.json | grep ShunFengParcel
```

#### 检查 2：Collector 到 Jaeger 的连接
```bash
# 在 Collector 服务器上测试
telnet 14.103.153.242 14250
```

#### 检查 3：应用程序日志
应该能看到：
```
INFO msg=✅ OTLP Tracer initialized
INFO msg=   Endpoint: 14.103.153.242:4317 (OpenTelemetry Collector)
触发 Trace 埋点，用户名: testuser
```

#### 检查 4：网络连接
```powershell
Test-NetConnection -ComputerName 14.103.153.242 -Port 4317
```

### 问题：Collector 日志中有错误

#### 常见错误 1：无法连接到 Jaeger
```
Error: failed to export to jaeger: connection refused
```
**解决**：检查 Jaeger 是否在 14250 端口监听

#### 常见错误 2：配置文件格式错误
```
Error: failed to load config: yaml: ...
```
**解决**：检查 YAML 文件格式，确保缩进正确

## 📁 相关文档

- **详细设置指南**: `JAEGER_SETUP_GUIDE.md`
- **完整总结**: `JAEGER_INTEGRATION_SUMMARY.md`
- **测试脚本**: `test-jaeger-integration.ps1`
- **Collector 配置**: `deploy/otel-collector-config.yaml`

## 🎯 核心修改点总结

### 代码修改（已完成）

1. **admin.go** - 第 33-57 行
   ```go
   ctx, span := otel.Tracer("admin-service").Start(ctx, "FindAdminByUsername")
   defer span.End()
   span.SetAttributes(...)
   span.AddEvent(...)
   ```

2. **main.go** - 第 55-100 行
   ```go
   func initTracerProvider() func(context.Context) error {
       exporter, _ := otlptracegrpc.New(...)
       tracerProvider := trace.NewTracerProvider(...)
       otel.SetTracerProvider(tracerProvider)
   }
   ```

### 配置修改（已完成）

3. **otel-collector-config.yaml** - 第 27-31 行
   ```yaml
   exporters:
     jaeger:
       endpoint: 14.103.153.242:14250
       tls:
         insecure: true
   ```

4. **otel-collector-config.yaml** - 第 46 行
   ```yaml
   exporters: [logging, file, jaeger]  # 添加了 jaeger
   ```

### 部署操作（待执行）

5. **重启 OpenTelemetry Collector** ← **这是唯一需要手动执行的步骤！**

## ✅ 完成确认

执行完 Collector 重启后，请确认：

- [ ] Collector 日志中没有错误
- [ ] Collector 日志中能看到 "Jaeger exporter enabled"
- [ ] 应用程序成功启动
- [ ] 调用接口返回正常响应
- [ ] Jaeger UI 中能看到 "ShunFengParcel" 服务
- [ ] 点击 Find Traces 能看到追踪数据
- [ ] Span 详情中能看到设置的属性和事件

## 🎉 成功标志

当您在 Jaeger UI 中看到类似下面的内容时，说明集成成功：

```
Service: ShunFengParcel
  └─ Trace ID: abc123...
      └─ Span: FindAdminByUsername
          Duration: 15.2ms
          Tags:
            - service.name: ShunFengParcel
            - username: testuser
            - operation: find_admin_by_username
          Events:
            - 用户查询成功
```

---

**如有任何问题，请参考 `JAEGER_SETUP_GUIDE.md` 中的详细说明。**

