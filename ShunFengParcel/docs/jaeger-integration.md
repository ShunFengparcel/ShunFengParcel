# Jaeger 链路追踪集成文档

## 概述

本项目已集成 Jaeger 分布式链路追踪系统，用于监控和排查微服务调用链路。

## 配置信息

### Jaeger 服务器
- **服务器地址**: `14.103.153.242`
- **UI 端口**: `16686` (http://14.103.153.242:16686)
- **OTLP HTTP 端口**: `4318` (数据上报)
- **服务名称**: `ShunFengParcel`

### 协议说明
使用 **OTLP HTTP** 协议上报链路数据到 Jaeger。OTLP (OpenTelemetry Protocol) 是 OpenTelemetry 的标准协议，Jaeger 从 v1.35+ 开始原生支持。

## 架构

```
┌─────────────────┐
│  ShunFengParcel │
│   (Go Service)  │
└────────┬────────┘
         │ OTLP HTTP
         │ Port: 4318
         ▼
┌─────────────────┐
│  Jaeger Server  │
│  14.103.153.242 │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│   Jaeger UI     │
│   Port: 16686   │
└─────────────────┘
```

## 功能特性

### 1. 自动链路追踪
- HTTP 请求自动生成 trace
- gRPC 调用自动生成 trace
- 自动记录请求耗时
- 自动记录错误信息

### 2. 日志关联
日志中包含 trace 信息：
```
trace.id=xxx span.id=xxx
```

### 3. 采样策略
当前配置：**AlwaysSample**（100% 采样）
- 开发环境：建议 100% 采样
- 生产环境：建议调整为 10-20% 采样

## 使用方法

### 1. 启动服务
```bash
cd ShunFengParcel
go run ./cmd/ShunFengParcel
```

启动日志会显示：
```
INFO Jaeger tracing initialized successfully
```

### 2. 访问 Jaeger UI
打开浏览器访问：
```
http://14.103.153.242:16686
```

### 3. 查看链路
1. 在 Jaeger UI 左侧选择服务：`ShunFengParcel`
2. 点击 "Find Traces" 查看链路列表
3. 点击具体的 trace 查看详细调用链

### 4. 搜索功能
- **按操作名搜索**: 如 `/admin/register`
- **按标签搜索**: 如 `http.status_code=200`
- **按时间范围**: 选择时间范围查看历史数据
- **按耗时**: 查找慢请求

## 链路信息说明

### Span 标签
每个 span 包含以下信息：
- `service.name`: 服务名称
- `service.version`: 服务版本
- `http.method`: HTTP 方法
- `http.url`: 请求 URL
- `http.status_code`: 响应状态码
- `error`: 是否有错误

### 示例链路
```
ShunFengParcel
  └─ POST /admin/register (150ms)
      ├─ MySQL Query (50ms)
      ├─ Redis Get (10ms)
      └─ External API Call (80ms)
```

## 性能影响

### 资源消耗
- **CPU**: < 1% 额外开销
- **内存**: ~10MB 额外内存
- **网络**: 每个请求约 1-2KB 数据上报

### 优化建议
1. **生产环境降低采样率**:
   ```go
   tracesdk.WithSampler(tracesdk.TraceIDRatioBased(0.1)) // 10% 采样
   ```

2. **批量上报**:
   ```go
   tracesdk.WithBatcher(exp,
       tracesdk.WithMaxExportBatchSize(512),
       tracesdk.WithBatchTimeout(5 * time.Second),
   )
   ```

## 故障排查

### 1. 无法连接到 Jaeger
**症状**: 启动时报错 `failed to create OTLP exporter`

**解决方案**:
- 检查网络连接: `ping 14.103.153.242`
- 检查端口开放: `telnet 14.103.153.242 4318`
- 检查防火墙规则

### 2. UI 中看不到数据
**可能原因**:
- 服务未产生请求
- 采样率设置为 0
- 时间范围选择不正确

**解决方案**:
- 发送测试请求
- 检查采样配置
- 调整 UI 时间范围

### 3. 数据延迟
**正常情况**: 数据会有 5-10 秒延迟
**原因**: 批量上报机制

## 端口说明

| 端口 | 协议 | 用途 | 状态 |
|------|------|------|------|
| 16686 | HTTP | Jaeger UI | ✅ 开放 |
| 4318 | HTTP | OTLP HTTP Collector | ✅ 开放 |
| 4317 | gRPC | OTLP gRPC Collector | ❌ 未测试 |
| 14268 | HTTP | Jaeger HTTP Collector (旧) | ❌ 关闭 |
| 14250 | gRPC | Jaeger gRPC Collector (旧) | ❌ 关闭 |
| 6831 | UDP | Jaeger Agent (旧) | ❌ 关闭 |

## 最佳实践

### 1. 为关键操作添加自定义 Span
```go
import "go.opentelemetry.io/otel"

func criticalOperation(ctx context.Context) {
    tracer := otel.Tracer("ShunFengParcel")
    ctx, span := tracer.Start(ctx, "critical-operation")
    defer span.End()
    
    // 业务逻辑
}
```

### 2. 添加自定义标签
```go
span.SetAttributes(
    attribute.String("user.id", userID),
    attribute.Int("order.count", count),
)
```

### 3. 记录错误
```go
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, err.Error())
}
```

## 相关链接

- [Jaeger 官方文档](https://www.jaegertracing.io/docs/)
- [OpenTelemetry Go SDK](https://opentelemetry.io/docs/instrumentation/go/)
- [OTLP 协议规范](https://opentelemetry.io/docs/specs/otlp/)

## 维护信息

- **集成时间**: 2025-10-23
- **协议版本**: OTLP HTTP 1.0
- **OpenTelemetry SDK**: v1.38.0
- **维护人员**: [Your Name]
