# Jaeger 追踪集成 - 快速使用指南

## ✅ 问题已解决

启动报错的问题已修复：
- **问题**：配置文件路径错误 (`../../configs`)
- **修复**：已改为 `./configs`
- **状态**：✅ 服务现在可以正常启动

## 🚀 快速启动

### 方式 1：使用启动脚本（推荐）

```powershell
cd E:\gowork\src\ShunFengParcel\ShunFengParcel
.\start.ps1
```

### 方式 2：手动启动

```powershell
cd E:\gowork\src\ShunFengParcel\ShunFengParcel
.\bin\ShunFengParcel.exe -conf .\configs
```

### 方式 3：使用 go run

```powershell
cd E:\gowork\src\ShunFengParcel\ShunFengParcel
go run ./cmd/ShunFengParcel
```

## 🛑 停止服务

```powershell
.\stop.ps1
# 或
Stop-Process -Name "ShunFengParcel" -Force
```

## 📊 当前状态

### ✅ 已完成
- [x] 追踪代码实现 (`admin.go`)
- [x] TracerProvider 配置 (`main.go`)
- [x] 配置文件路径修复
- [x] 服务成功启动
- [x] 接口正常响应

### ⚠️ 待完成（需要在服务器上操作）
- [ ] 在服务器 14.103.153.242 上重启 OpenTelemetry Collector
- [ ] 使用更新后的配置文件 `deploy/otel-collector-config.yaml`

## 🧪 测试追踪功能

### 1. 调用接口产生追踪数据

```powershell
# 调用 FindAdminByUsername 接口
curl.exe -X POST http://localhost:8000/admin/findadminbyusername `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"test123"}'

# 预期响应
{"id":"11"}
```

### 2. 查看 Jaeger UI

访问：http://14.103.153.242:16686/search

**当前状态**：
- Service 列表中只有 `jaeger-all-in-one`
- 还看不到 `ShunFengParcel` 服务

**原因**：
OpenTelemetry Collector 还没有配置 Jaeger 导出器

## 🔧 完成 Jaeger 集成的最后一步

### 在服务器 14.103.153.242 上执行：

```bash
# 1. SSH 到服务器
ssh user@14.103.153.242

# 2. 上传配置文件
# 将 ShunFengParcel/deploy/otel-collector-config.yaml 上传到服务器

# 3. 重启 Collector（根据部署方式选择）

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

### 完成后测试：

```powershell
# 1. 调用接口
curl.exe -X POST http://localhost:8000/admin/findadminbyusername `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"test123"}'

# 2. 等待几秒

# 3. 检查 Jaeger
# 访问 http://14.103.153.242:16686/search
# 应该能在 Service 列表中看到 "ShunFengParcel"
```

## 📁 服务信息

- **HTTP 服务**: http://localhost:8000
- **gRPC 服务**: localhost:9000
- **Jaeger UI**: http://14.103.153.242:16686/search
- **OpenTelemetry Collector**: 14.103.153.242:4317

## 📚 详细文档

- **JAEGER_FINAL_CHECKLIST.md** - 最终检查清单
- **JAEGER_INTEGRATION_SUMMARY.md** - 完整集成总结
- **JAEGER_SETUP_GUIDE.md** - 详细设置指南

## 🔍 追踪代码位置

### admin.go 中的追踪实现

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

## 🎯 预期的 Jaeger 追踪结果

完成 Collector 配置后，在 Jaeger UI 中应该看到：

```
Service: ShunFengParcel
  └─ Trace
      └─ Span: FindAdminByUsername
          Duration: ~10-50ms
          Tags:
            - service.name: ShunFengParcel
            - service.version: 1.0.0
            - username: testuser
            - operation: find_admin_by_username
          Events:
            - 用户查询成功 / 密码验证失败
```

## ❓ 常见问题

### Q: 服务启动报错 "cannot find the file specified"
**A**: 已修复。确保使用最新编译的版本：
```powershell
go build -o bin/ShunFengParcel.exe ./cmd/ShunFengParcel
.\bin\ShunFengParcel.exe -conf .\configs
```

### Q: Jaeger UI 中看不到 ShunFengParcel 服务
**A**: 需要在服务器上重启 OpenTelemetry Collector 并使用更新后的配置文件。

### Q: 如何查看服务日志
**A**: 服务在后台运行，可以通过以下方式查看：
```powershell
Get-Process ShunFengParcel
```

## ✅ 总结

**应用程序端的工作已全部完成！**

- ✅ 代码修改完成
- ✅ 配置文件路径修复
- ✅ 服务可以正常启动
- ✅ 接口正常响应
- ✅ 追踪数据正在发送到 Collector

**下一步**：在服务器 14.103.153.242 上重启 OpenTelemetry Collector，就能在 Jaeger UI 中看到完整的追踪数据了！🎉

