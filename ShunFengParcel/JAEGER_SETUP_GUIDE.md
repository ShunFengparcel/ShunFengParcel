# Jaeger 集成配置指南

## 问题诊断

经过测试发现，您的 **OpenTelemetry Collector 没有将追踪数据转发到 Jaeger**。

原有配置只将追踪数据输出到日志和文件，没有配置 Jaeger 导出器。

## 解决方案

### 1. 更新 OpenTelemetry Collector 配置

已更新 `deploy/otel-collector-config.yaml`，添加了 Jaeger 导出器：

```yaml
exporters:
  jaeger:
    endpoint: 14.103.153.242:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging, file, jaeger]  # 添加了 jaeger
```

### 2. 重启 OpenTelemetry Collector

您需要在部署 OpenTelemetry Collector 的服务器上执行以下操作：

#### 如果使用 Docker 部署：

```bash
# 停止现有的 collector
docker stop otel-collector

# 使用新配置启动
docker run -d --name otel-collector \
  -p 4317:4317 \
  -p 4318:4318 \
  -v /path/to/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  otel/opentelemetry-collector:latest \
  --config=/etc/otel-collector-config.yaml
```

#### 如果使用 Kubernetes 部署：

```bash
# 更新 ConfigMap
kubectl apply -f otel-collector-config.yaml

# 重启 collector pod
kubectl rollout restart deployment/otel-collector
```

#### 如果使用系统服务：

```bash
# 复制新配置
sudo cp deploy/otel-collector-config.yaml /etc/otel-collector/config.yaml

# 重启服务
sudo systemctl restart otel-collector
```

### 3. 验证配置

重启 Collector 后，检查日志确认 Jaeger 导出器已启动：

```bash
# Docker
docker logs otel-collector

# Kubernetes  
kubectl logs -f deployment/otel-collector

# 系统服务
sudo journalctl -u otel-collector -f
```

应该看到类似的日志：
```
Jaeger exporter enabled
Exporter is starting...
```

### 4. 测试应用程序

```bash
# 重新编译
cd ShunFengParcel
go build -o bin/ShunFengParcel.exe ./cmd/ShunFengParcel

# 重启服务
# Windows PowerShell:
Stop-Process -Name "ShunFengParcel" -Force
Start-Process -FilePath ".\bin\ShunFengParcel.exe" -ArgumentList "-conf", ".\configs"

# 调用接口产生追踪数据
curl -X POST http://localhost:8000/admin/findadminbyusername \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"test123"}'
```

### 5. 查看 Jaeger UI

访问：http://14.103.153.242:16686/search

- 在 **Service** 下拉框中选择 `ShunFengParcel`
- 点击 **Find Traces** 按钮
- 应该能看到追踪数据和 Span 信息

## 架构说明

```
ShunFengParcel 应用
    ↓ (OTLP gRPC)
    ↓ 端口 4317
OpenTelemetry Collector
    ↓ (Jaeger gRPC)
    ↓ 端口 14250
Jaeger Backend
    ↓
Jaeger UI (端口 16686)
```

## 当前配置

- **应用程序**: 使用 OTLP gRPC 协议发送到 Collector (14.103.153.242:4317)
- **Collector**: 接收 OTLP 数据，转发到 Jaeger (14.103.153.242:14250)
- **Jaeger UI**: http://14.103.153.242:16686/search

## admin.go 中的追踪代码

```go
func (s *AdminService) FindAdminByUsername(ctx context.Context, req *pb.FindAdminByUsernameReq) (*pb.FindAdminByUsernameResp, error) {
    // 创建 Span
    ctx, span := otel.Tracer("admin-service").Start(ctx, "FindAdminByUsername")
    defer span.End()

    // 添加属性
    span.SetAttributes(
        attribute.String("username", req.Username),
        attribute.String("operation", "find_admin_by_username"),
    )

    // 业务逻辑...
    var a config.SysAdmin
    inits.DB.Where("username = ?", req.Username).Find(&a)
    
    if a.Password != pkg.Md5(req.Password) {
        span.AddEvent("密码验证失败")
        return nil, nil
    }

    span.AddEvent("用户查询成功")
    return &pb.FindAdminByUsernameResp{Id: a.Id}, nil
}
```

## 故障排查

### 如果还是看不到追踪数据：

1. **检查 Collector 是否运行**：
   ```bash
   # 检查端口
   netstat -an | grep 4317
   ```

2. **检查 Collector 日志**：查看是否有错误信息

3. **检查 Jaeger 端口 14250**：
   ```bash
   telnet 14.103.153.242 14250
   ```

4. **测试 Collector 连接**：
   ```bash
   curl -X POST http://14.103.153.242:4318/v1/traces \
     -H "Content-Type: application/json" \
     -d '{"resourceSpans":[]}'
   ```

5. **查看 Collector 的追踪数据**：
   ```bash
   # 如果配置了文件导出
   tail -f /var/log/otel-collector/traces.json
   ```

## 需要您执行的操作

**最重要的一步**：请在部署 OpenTelemetry Collector 的服务器 (14.103.153.242) 上：

1. 使用更新后的 `deploy/otel-collector-config.yaml` 配置文件
2. 重启 OpenTelemetry Collector 服务
3. 确认 Collector 日志中没有错误

完成后，重新运行应用程序并调用接口，就能在 Jaeger UI 中看到追踪数据了！

