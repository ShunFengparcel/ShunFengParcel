# OpenTelemetry Collector 配置更新指南

## 📋 您的环境信息

根据您的 1Panel 截图：

- **Jaeger 容器**: `1Panel-jaeger-SBwi` (172.18.0.8)
  - 端口: 14268/tcp, 16686/tcp
  
- **OpenTelemetry Collector 容器**: `otel-collector` (172.17.0.5)
  - 端口: 4317-4318/tcp, 55679/tcp

## 🎯 需要做什么

更新 `otel-collector` 容器的配置，添加 Jaeger 导出器，让它能将追踪数据转发到 Jaeger。

## 方法 1：通过 1Panel 界面更新（最简单）⭐

### 步骤 1：准备配置文件

配置文件已经准备好了：`E:\gowork\src\ShunFengParcel\ShunFengParcel\deploy\otel-collector-config.yaml`

### 步骤 2：在 1Panel 中操作

1. **打开 1Panel 界面**
2. **找到 `otel-collector` 容器**
3. **点击"终端"按钮**进入容器终端
4. **查看当前配置路径**：
   ```bash
   # 查找配置文件
   find / -name "*otel*config*" 2>/dev/null
   # 或者查看容器启动命令
   ps aux | grep otel
   ```

5. **记下配置文件路径**（通常是 `/etc/otel-collector-config.yaml` 或 `/etc/otelcol/config.yaml`）

### 步骤 3：上传新配置

**方式 A：通过 1Panel 文件管理**
1. 在 1Panel 中点击容器的"文件"按钮
2. 导航到配置文件目录（如 `/etc/`）
3. 备份原配置文件（重命名为 `otel-collector-config.yaml.bak`）
4. 上传新的 `otel-collector-config.yaml`

**方式 B：通过命令行**
```bash
# 在服务器上执行（SSH 到 14.103.153.242）
docker cp /path/to/otel-collector-config.yaml otel-collector:/etc/otel-collector-config.yaml
```

### 步骤 4：重启容器

在 1Panel 界面中：
1. 选中 `otel-collector` 容器
2. 点击"重启"按钮
3. 等待容器重启完成（约 5-10 秒）

### 步骤 5：检查日志

在 1Panel 中点击容器的"日志"按钮，应该看到：
```
Jaeger exporter enabled
Exporter is starting...
Everything is ready. Begin running and processing data.
```

## 方法 2：使用自动化脚本（需要 SSH 访问）

### 在服务器上执行：

```bash
# SSH 到服务器
ssh user@14.103.153.242

# 上传脚本和配置文件到服务器
# 然后执行：
cd /path/to/ShunFengParcel
chmod +x update-otel-collector.sh
./update-otel-collector.sh
```

### 或在 Windows 上通过 SSH 执行：

```powershell
# 确保能访问服务器的 Docker
ssh user@14.103.153.242 "cd /path/to/ShunFengParcel && ./update-otel-collector.sh"
```

## 方法 3：手动更新（详细步骤）

### 1. 连接到服务器

```bash
ssh user@14.103.153.242
```

### 2. 备份当前配置

```bash
docker exec otel-collector cat /etc/otel-collector-config.yaml > otel-config-backup.yaml
```

### 3. 创建新配置文件

```bash
cat > otel-collector-config-new.yaml << 'EOF'
receivers:
  otlp:
    protocols:
      grpc:
        endpoint: 0.0.0.0:4317
      http:
        endpoint: 0.0.0.0:4318

processors:
  batch:
    timeout: 10s
    send_batch_size: 1024
  
  memory_limiter:
    check_interval: 1s
    limit_mib: 512

exporters:
  logging:
    loglevel: info
  
  file:
    path: /var/log/otel-collector/traces.json
  
  # 添加 Jaeger 导出器
  jaeger:
    endpoint: 172.18.0.8:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging, file, jaeger]
    
    metrics:
      receivers: [otlp]
      processors: [memory_limiter, batch]
      exporters: [logging, file]
  
  telemetry:
    logs:
      level: info
EOF
```

**注意**：Jaeger 端点使用 `172.18.0.8:14250`（您的 Jaeger 容器 IP）

### 4. 复制配置到容器

```bash
docker cp otel-collector-config-new.yaml otel-collector:/etc/otel-collector-config.yaml
```

### 5. 重启容器

```bash
docker restart otel-collector
```

### 6. 检查日志

```bash
docker logs -f otel-collector
```

应该看到：
- `Jaeger exporter enabled`
- `Everything is ready`
- 没有错误信息

## ✅ 验证配置是否生效

### 1. 重启您的应用程序

```powershell
cd E:\gowork\src\ShunFengParcel\ShunFengParcel
.\stop.ps1
.\start.ps1
```

### 2. 调用接口产生追踪数据

```powershell
curl.exe -X POST http://localhost:8000/admin/findadminbyusername `
  -H "Content-Type: application/json" `
  -d '{"username":"testuser","password":"test123"}'
```

### 3. 查看 Jaeger UI

访问：http://14.103.153.242:16686/search

在 **Service** 下拉框中应该能看到：
- ✅ `ShunFengParcel` （新增的服务）
- `jaeger-all-in-one`

### 4. 查看追踪详情

1. 选择 `ShunFengParcel` 服务
2. 点击 **Find Traces** 按钮
3. 应该能看到追踪记录
4. 点击追踪记录查看详情：
   - Span 名称: `FindAdminByUsername`
   - Tags: `username`, `operation`
   - Events: `用户查询成功` 或 `密码验证失败`

## 🔍 故障排查

### 问题 1：配置文件路径不对

**症状**：复制配置文件失败

**解决**：
```bash
# 查找配置文件位置
docker exec otel-collector find / -name "*config*" 2>/dev/null | grep otel

# 或查看容器启动命令
docker inspect otel-collector | grep -A 10 "Cmd"
```

### 问题 2：Jaeger 端点连接失败

**症状**：Collector 日志显示 `connection refused`

**解决**：
1. 检查 Jaeger 容器 IP：
   ```bash
   docker inspect 1Panel-jaeger-SBwi | grep IPAddress
   ```
2. 更新配置文件中的 `endpoint` 为正确的 IP
3. 确保端口是 14250（Jaeger gRPC 端口）

### 问题 3：容器重启后配置丢失

**症状**：重启后又看不到追踪数据

**解决**：需要使用 volume 挂载配置文件
```bash
docker run -d --name otel-collector \
  -v /path/to/otel-collector-config.yaml:/etc/otel-collector-config.yaml \
  ...其他参数...
```

### 问题 4：还是看不到 ShunFengParcel 服务

**检查清单**：
1. ✅ Collector 配置已更新
2. ✅ Collector 已重启
3. ✅ Collector 日志没有错误
4. ✅ 应用程序已重启
5. ✅ 已调用接口产生追踪数据
6. ✅ 等待了 10-30 秒

如果都确认了还是不行，查看：
```bash
# Collector 日志
docker logs otel-collector | grep -i error

# 查看导出的文件
docker exec otel-collector cat /var/log/otel-collector/traces.json | grep ShunFengParcel
```

## 📝 配置文件说明

关键修改点：

```yaml
exporters:
  jaeger:  # 新增
    endpoint: 172.18.0.8:14250  # Jaeger 容器的 IP:端口
    tls:
      insecure: true

service:
  pipelines:
    traces:
      exporters: [logging, file, jaeger]  # 添加 jaeger
```

## 🎯 总结

**最简单的方法**：
1. 在 1Panel 界面找到 `otel-collector` 容器
2. 通过文件管理上传新配置文件
3. 点击重启按钮
4. 查看日志确认成功
5. 测试应用程序

完成后，您就能在 Jaeger UI 中看到完整的追踪数据了！🎉

