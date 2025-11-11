# APM 配置从 HTTP 迁移到 gRPC

## 修改内容

### 1. 导入包变更

**之前（HTTP 方式）：**
```go
"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
```

**现在（gRPC 方式）：**
```go
"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
```

### 2. 端口变更

**之前（HTTP）：**
```go
apmEndpoint := "14.103.153.242:4318"  // HTTP 端口
```

**现在（gRPC）：**
```go
apmEndpoint := "14.103.153.242:4317"  // gRPC 端口
```

### 3. Trace 导出器配置

**之前：**
```go
traceExporter, err := otlptracehttp.New(ctx,
    otlptracehttp.WithInsecure(),
    otlptracehttp.WithEndpoint(apmEndpoint),
    otlptracehttp.WithHeaders(map[string]string{
        "X-ByteAPIM-AppKey": apmAppKey,
    }),
)
```

**现在：**
```go
traceExporter, err := otlptracegrpc.New(ctx,
    otlptracegrpc.WithInsecure(),
    otlptracegrpc.WithEndpoint(apmEndpoint),
    otlptracegrpc.WithHeaders(map[string]string{
        "X-ByteAPIM-AppKey": apmAppKey,
    }),
)
```

### 4. Metrics 导出器配置

**之前：**
```go
metricsExporter, err := otlpmetrichttp.New(ctx,
    otlpmetrichttp.WithInsecure(),
    otlpmetrichttp.WithEndpoint(apmEndpoint),
    otlpmetrichttp.WithHeaders(map[string]string{
        "X-ByteAPIM-AppKey": apmAppKey,
    }),
)
```

**现在：**
```go
metricsExporter, err := otlpmetricgrpc.New(ctx,
    otlpmetricgrpc.WithInsecure(),
    otlpmetricgrpc.WithEndpoint(apmEndpoint),
    otlpmetricgrpc.WithHeaders(map[string]string{
        "X-ByteAPIM-AppKey": apmAppKey,
    }),
)
```

## gRPC vs HTTP 对比

| 特性 | HTTP | gRPC |
|------|------|------|
| 端口 | 4318 | 4317 |
| 协议 | HTTP/1.1 | HTTP/2 |
| 性能 | 较低 | 更高 |
| 连接 | 短连接 | 长连接 |
| 二进制传输 | JSON | Protobuf |
| 适用场景 | 简单部署 | 高性能生产环境 |

## 优势

1. **更高性能**：gRPC 使用 HTTP/2 和 Protobuf，传输效率更高
2. **更低延迟**：长连接减少了连接建立的开销
3. **更小带宽**：二进制序列化比 JSON 更紧凑
4. **更好的流控**：HTTP/2 原生支持流控和多路复用

## 注意事项

1. 确保 APM 服务端支持 gRPC（端口 4317）
2. 如果有防火墙，需要开放 4317 端口
3. gRPC 需要 HTTP/2 支持
4. 如果遇到连接问题，可以回退到 HTTP 方式（端口 4318）

## 测试

重启服务后，检查日志确认连接成功：
```bash
# 查看是否有 trace/metrics 导出错误
# 如果没有错误日志，说明连接成功
```

## 回退方案

如果需要回退到 HTTP 方式：
1. 将端口改回 `4318`
2. 将导入包改回 `otlptracehttp` 和 `otlpmetrichttp`
3. 运行 `go mod tidy`
