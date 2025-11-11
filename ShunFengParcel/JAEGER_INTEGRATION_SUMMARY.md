# Jaeger 集成完成总结

## ✅ 已完成的工作

### 1. 应用程序追踪代码 ✅

**文件**: `ShunFengParcel/internal/service/admin.go`

在 `FindAdminByUsername` 方法中添加了完整的追踪逻辑：

```go
func (s *AdminService) FindAdminByUsername(ctx context.Context, req *pb.FindAdminByUsernameReq) (*pb.FindAdminByUsernameResp, error) {
    // 创建 Span（在方法开始时）
    ctx, span := otel.Tracer("admin-service").Start(ctx, "FindAdminByUsername")
    defer span.End()

    // 添加 Span 属性
    span.SetAttributes(
        attribute.String("username", req.Username),
        attribute.String("operation", "find_admin_by_username"),
    )

    // 记录日志
    log.Printf("触发 Trace 埋点，用户名: %s", req.Username)

    // 业务逻辑
    var a config.SysAdmin
    inits.DB.Where("username = ?", req.Username).Find(&a)
    
    if a.Password != pkg.Md5(req.Password) {
        span.AddEvent("密码验证失败")  // 记录事件
        return nil, nil
    }

    span.AddEvent("用户查询成功")  // 记录事件
    return &pb.FindAdminByUsernameResp{Id: a.Id}, nil
}
```

**关键点**：
- ✅ Span 在方法开始时创建
- ✅ 使用 `defer span.End()` 确保 Span 结束
- ✅ 添加了属性记录关键信息
- ✅ 添加了事件标记业务节点

### 2. TracerProvider 配置 ✅

**文件**: `ShunFengParcel/cmd/ShunFengParcel/main.go`

配置了 OTLP gRPC 导出器，连接到 OpenTelemetry Collector：

```go
func initTracerProvider() func(context.Context) error {
    // 创建 OTLP gRPC 导出器
    exporter, err := otlptracegrpc.New(ctx,
        otlptracegrpc.WithEndpoint("14.103.153.242:4317"),
        otlptracegrpc.WithInsecure(),
    )
    
    // 创建资源信息
    res, err := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName("ShunFengParcel"),
            semconv.ServiceVersion("1.0.0"),
            semconv.ServiceInstanceID(id),
        ),
    )
    
    // 创建 TracerProvider
    tracerProvider := trace.NewTracerProvider(
        trace.WithBatcher(exporter),
        trace.WithResource(res),
        trace.WithSampler(trace.AlwaysSample()),
    )
    otel.SetTracerProvider(tracerProvider)
    
    return shutdownFunc
}
```

**配置说明**：
- ✅ 使用 OTLP gRPC 协议
- ✅ 连接到 Collector (14.103.153.242:4317)
- ✅ 服务名称设置为 "ShunFengParcel"
- ✅ 采样策略设置为总是采样

### 3. HTTP Server 中间件 ✅

**文件**: `ShunFengParcel/internal/server/http.go`

HTTP Server 已配置追踪中间件：

```go
func NewHTTPServer(...) *http.Server {
    var opts = []http.ServerOption{
        http.Middleware(
            recovery.Recovery(),
            tracing.Server(),  // ✅ 追踪中间件
        ),
        http.Filter(corsFilter()),
    }
    // ...
}
```

### 4. OpenTelemetry Collector 配置 ✅

**文件**: `ShunFengParcel/deploy/otel-collector-config.yaml`

已更新配置，添加了 Jaeger 导出器：

```yaml
exporters:
  logging:
    loglevel: info
  file:
    path: /var/log/otel-collector/traces.json
  jaeger:  # ✅ 新增
    endpoint: 14.103.153.242:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging, file, jaeger]  # ✅ 添加了 jaeger
```

## ⚠️ 需要您执行的操作

### 关键步骤：重启 OpenTelemetry Collector

**这是唯一需要您手动执行的步骤！**

在服务器 **14.103.153.242** 上执行以下操作：

#### 方式 1：如果使用 Docker

```bash
# SSH 到服务器
ssh user@14.103.153.242

# 停止现有的 collector
docker stop otel-collector
docker rm otel-collector

# 上传新配置文件
# 将 ShunFengParcel/deploy/otel-collector-config.yaml 上传到服务器

# 启动 collector
docker run -d --name otel-collector \
  --restart=always \
  -p 4317:4317 \
  -p 4318:4318 \
  -p 14250:14250 \
  -v /path/to/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  otel/opentelemetry-collector-contrib:latest \
  --config=/etc/otel-collector-config.yaml

# 查看日志确认启动成功
docker logs -f otel-collector
```

#### 方式 2：如果使用 Kubernetes

```bash
# 更新 ConfigMap
kubectl create configmap otel-collector-config \
  --from-file=otel-collector-config.yaml \
  --dry-run=client -o yaml | kubectl apply -f -

# 重启 collector
kubectl rollout restart deployment/otel-collector

# 查看日志
kubectl logs -f deployment/otel-collector
```

#### 方式 3：如果使用系统服务

```bash
# 备份旧配置
sudo cp /etc/otel-collector/config.yaml /etc/otel-collector/config.yaml.bak

# 复制新配置
sudo cp otel-collector-config.yaml /etc/otel-collector/config.yaml

# 重启服务
sudo systemctl restart otel-collector

# 查看状态
sudo systemctl status otel-collector

# 查看日志
sudo journalctl -u otel-collector -f
```

## 🧪 测试步骤

完成 Collector 重启后，执行以下测试：

### 1. 重启应用程序

```powershell
# 在项目目录下
cd E:\gowork\src\ShunFengParcel\ShunFengParcel

# 停止旧服务
Stop-Process -Name "ShunFengParcel" -Force

# 启动新服务
Start-Process -FilePath ".\bin\ShunFengParcel.exe" -ArgumentList "-conf", ".\configs"
```

### 2. 调用接口

```powershell
# 调用接口产生追踪数据
curl.exe -X POST http://localhost:8000/admin/findadminbyusername `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"test123"}'
```

### 3. 查看 Jaeger UI

1. 访问：http://14.103.153.242:16686/search
2. 在 **Service** 下拉框中选择 `ShunFengParcel`
3. 点击 **Find Traces** 按钮
4. 查看追踪详情

## 📊 预期结果

在 Jaeger UI 中应该能看到：

### Service 列表
- ✅ `ShunFengParcel` 服务

### Trace 详情
- **Operation**: `FindAdminByUsername`
- **Span 名称**: `FindAdminByUsername`
- **Tags**:
  - `service.name`: `ShunFengParcel`
  - `username`: `testuser`
  - `operation`: `find_admin_by_username`
- **Events**:
  - `密码验证失败` 或 `用户查询成功`

### Span 层次结构
```
HTTP Request (由 Kratos tracing 中间件创建)
  └─ FindAdminByUsername (由 admin.go 中的代码创建)
       └─ 数据库查询 (如果配置了 GORM tracing)
```

## 🔍 故障排查

### 如果还是看不到 ShunFengParcel 服务：

1. **检查 Collector 日志**：
   ```bash
   # Docker
   docker logs otel-collector | grep -i error
   
   # Kubernetes
   kubectl logs deployment/otel-collector | grep -i error
   
   # 系统服务
   sudo journalctl -u otel-collector | grep -i error
   ```

2. **检查 Collector 是否收到数据**：
   在 Collector 日志中应该能看到类似信息：
   ```
   Trace received from ShunFengParcel
   ```

3. **检查 Jaeger 连接**：
   ```bash
   # 在 Collector 服务器上测试
   telnet 14.103.153.242 14250
   ```

4. **查看 Collector 导出的文件**：
   ```bash
   tail -f /var/log/otel-collector/traces.json
   ```
   应该能看到 ShunFengParcel 的追踪数据

5. **检查应用程序日志**：
   应该能看到：
   ```
   触发 Trace 埋点，用户名: testuser
   ```

## 📁 相关文件

- ✅ `ShunFengParcel/internal/service/admin.go` - 追踪代码
- ✅ `ShunFengParcel/cmd/ShunFengParcel/main.go` - TracerProvider 配置
- ✅ `ShunFengParcel/deploy/otel-collector-config.yaml` - Collector 配置
- ✅ `ShunFengParcel/JAEGER_SETUP_GUIDE.md` - 详细设置指南
- ✅ `ShunFengParcel/test-jaeger-integration.ps1` - 测试脚本

## 🎯 架构图

```
┌─────────────────────────┐
│  ShunFengParcel 应用     │
│  (localhost:8000)       │
│                         │
│  - HTTP Server          │
│  - tracing.Server()     │
│  - admin.go Spans       │
└────────────┬────────────┘
             │ OTLP gRPC
             │ (端口 4317)
             ↓
┌─────────────────────────┐
│ OpenTelemetry Collector │
│ (14.103.153.242:4317)   │
│                         │
│  - 接收 OTLP 数据        │
│  - 批处理                │
│  - 转发到 Jaeger         │
└────────────┬────────────┘
             │ Jaeger gRPC
             │ (端口 14250)
             ↓
┌─────────────────────────┐
│   Jaeger Backend        │
│ (14.103.153.242:14250)  │
│                         │
│  - 存储追踪数据          │
│  - 提供查询 API          │
└────────────┬────────────┘
             │
             ↓
┌─────────────────────────┐
│     Jaeger UI           │
│ (14.103.153.242:16686)  │
│                         │
│  - 可视化追踪数据        │
│  - 搜索和分析            │
└─────────────────────────┘
```

## ✅ 总结

**应用程序端的所有代码已完成并正确配置！**

唯一需要的操作是：**在服务器 14.103.153.242 上使用更新后的配置重启 OpenTelemetry Collector**。

完成后，您就能在 Jaeger UI 中看到 `ShunFengParcel` 服务和所有 Span 的追踪情况了！🎉

