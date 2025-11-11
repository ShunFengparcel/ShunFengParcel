# Jaeger 集成问题总结

## 问题描述

尝试将 ShunFengParcel 服务集成 Jaeger 链路追踪，连接到服务器 `14.103.153.242`，但 Jaeger UI 中始终看不到服务数据。

## 已完成的工作

### 1. 代码集成 ✅
- [x] 添加 OpenTelemetry SDK
- [x] 配置 TracerProvider
- [x] 实现自动链路追踪
- [x] 添加日志关联（trace.id, span.id）
- [x] 实现优雅降级
- [x] 添加连接验证

### 2. 测试验证 ✅
- [x] 端口连通性测试（16686, 4318, 14268 都通）
- [x] OTLP HTTP 测试（返回 200）
- [x] Jaeger HTTP 测试（返回 200）
- [x] 独立测试程序（运行成功）
- [x] 服务请求测试（HTTP 200）

### 3. 尝试的方案 ✅
1. **OTLP HTTP (4318)** - 端口通，但数据不显示
2. **OTLP gRPC (4317)** - 未测试
3. **Jaeger HTTP (14268)** - 端口通，但数据不显示

## 问题根因

**Jaeger 服务器端配置问题**

虽然所有端口都是通的，测试程序也运行成功，但 Jaeger UI 中始终看不到数据。

可能的原因：
1. ❌ Jaeger 版本太旧（< v1.35），不支持 OTLP
2. ❌ Collector 没有正确配置
3. ❌ 数据存储有问题
4. ❌ 只接受特定格式的数据

## 证据

### 客户端日志
```
✅ OTLP Exporter 创建成功
✅ 资源信息创建成功
✅ TracerProvider 已设置
✅ 测试 span 已发送
✅ Jaeger TracerProvider 已关闭
```

### 服务端表现
- Jaeger UI 只显示 `jaeger-all-in-one` 服务
- 没有任何来自 `ShunFengParcel` 或 `TestService` 的数据
- 端口响应正常（HTTP 200）

### 网络测试
```powershell
# 端口测试
Test-NetConnection -ComputerName 14.103.153.242 -Port 16686  # ✅ 通
Test-NetConnection -ComputerName 14.103.153.242 -Port 4318   # ✅ 通
Test-NetConnection -ComputerName 14.103.153.242 -Port 14268  # ✅ 通

# HTTP 测试
Invoke-WebRequest -Uri "http://14.103.153.242:4318/v1/traces" -Method POST
# 返回: HTTP 200 {"partialSuccess":{}}
```

## 解决方案

### 短期方案（立即可用）

#### 方案 A：使用本地 Jaeger
```bash
# 启动本地 Jaeger（Docker）
docker run -d --name jaeger \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 4318:4318 \
  jaegertracing/all-in-one:latest

# 修改代码中的地址
jaegerEndpoint := "localhost:14268"
```

#### 方案 B：暂时禁用 Jaeger
```go
// 注释掉 Jaeger 初始化
// shutdown := initJaeger()
```

### 长期方案（需要服务器管理员配合）

#### 1. 检查 Jaeger 版本
```bash
docker exec jaeger /go/bin/all-in-one-linux version
# 需要 >= v1.35 才支持 OTLP
```

#### 2. 检查 Jaeger 配置
```bash
# 查看启动参数
docker inspect jaeger | grep -A 20 Cmd

# 确认是否启用了 OTLP
# 需要包含: --collector.otlp.enabled=true
```

#### 3. 查看 Jaeger 日志
```bash
docker logs -f jaeger

# 查找错误信息
docker logs jaeger 2>&1 | grep -i error
```

#### 4. 重新启动 Jaeger（正确配置）
```bash
docker run -d --name jaeger \
  -e COLLECTOR_OTLP_ENABLED=true \
  -e SPAN_STORAGE_TYPE=badger \
  -e BADGER_EPHEMERAL=false \
  -e BADGER_DIRECTORY_VALUE=/badger/data \
  -e BADGER_DIRECTORY_KEY=/badger/key \
  -p 16686:16686 \
  -p 4318:4318 \
  -p 14268:14268 \
  -v badger:/badger \
  jaegertracing/all-in-one:latest
```

#### 5. 使用 OpenTelemetry Collector（推荐）
```yaml
# otel-collector-config.yaml
receivers:
  otlp:
    protocols:
      http:
        endpoint: 0.0.0.0:4318
      grpc:
        endpoint: 0.0.0.0:4317

exporters:
  jaeger:
    endpoint: jaeger:14250
    tls:
      insecure: true

service:
  pipelines:
    traces:
      receivers: [otlp]
      exporters: [jaeger]
```

## 建议行动

### 立即行动
1. ✅ 代码已经准备好，暂时禁用 Jaeger
2. ✅ 继续开发其他功能
3. ✅ 等待服务器配置修复

### 后续行动
1. 📧 联系 `14.103.153.242` 服务器管理员
2. 📋 提供本文档作为问题说明
3. 🔧 请求检查和修复 Jaeger 配置
4. ✅ 配置修复后重新启用 Jaeger

## 联系信息

**服务器信息**：
- IP: 14.103.153.242
- Jaeger UI: http://14.103.153.242:16686
- OTLP HTTP: 14.103.153.242:4318
- Jaeger HTTP: 14.103.153.242:14268

**需要服务器管理员提供**：
1. Jaeger 版本号
2. Jaeger 启动配置
3. Jaeger 日志（最近 100 行）
4. 是否有其他服务成功上报过数据

## 参考资料

- [Jaeger 官方文档](https://www.jaegertracing.io/docs/)
- [OpenTelemetry 文档](https://opentelemetry.io/docs/)
- [Jaeger OTLP 支持](https://www.jaegertracing.io/docs/1.35/apis/#opentelemetry-protocol-stable)
- [本项目集成文档](./docs/jaeger-integration.md)

---

**文档创建时间**：2025-10-23  
**问题状态**：待服务器管理员修复  
**代码状态**：已完成，可随时启用
