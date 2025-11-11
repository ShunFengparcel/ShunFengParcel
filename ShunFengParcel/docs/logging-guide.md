# 日志系统使用指南

## 概述

本项目使用结构化日志系统，支持 JSON 格式输出，便于 ELK 栈集成和日志分析。

## 日志配置

### 配置文件 (configs/config.yaml)

```yaml
log:
  level: info              # 日志级别: debug, info, warn, error
  format: json             # 日志格式: json 或 text
  output_path: ./logs/app.log  # 日志文件路径
  max_size: 100            # 单个日志文件最大大小（MB）
  console: true            # 是否同时输出到控制台
```

## 日志级别

- **DEBUG**: 详细的调试信息
- **INFO**: 一般信息，如请求日志
- **WARN**: 警告信息，不影响运行
- **ERROR**: 错误信息，需要关注
- **FATAL**: 致命错误，程序退出

## 使用示例

### 基础日志

```go
import (
    "github.com/go-kratos/kratos/v2/log"
)

// 在服务中使用
func (s *GreeterService) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    log.Context(ctx).Infow("msg", "received hello request", "name", req.Name)
    
    // 业务逻辑
    
    log.Context(ctx).Infow("msg", "hello request processed", "name", req.Name)
    return &pb.HelloReply{Message: "Hello " + req.Name}, nil
}
```

### 错误日志

```go
if err != nil {
    log.Context(ctx).Errorw(
        "msg", "failed to process request",
        "error", err,
        "user_id", userID,
        "operation", "create_order",
    )
    return nil, err
}
```

### 带上下文的日志

```go
// 添加请求 ID
ctx = context.WithValue(ctx, "request_id", uuid.New().String())

log.Context(ctx).Infow(
    "msg", "processing order",
    "order_id", orderID,
    "amount", amount,
)
```

## JSON 日志格式

```json
{
  "timestamp": "2025-10-23T10:30:45+08:00",
  "level": "INFO",
  "message": "received hello request",
  "service_id": "shunfeng-parcel-001",
  "trace_id": "abc123def456",
  "span_id": "789ghi",
  "caller": "service/greeter.go:25",
  "fields": {
    "name": "张三",
    "operation": "say_hello",
    "duration_ms": 15
  }
}
```

## 日志字段说明

| 字段 | 说明 | 示例 |
|------|------|------|
| timestamp | 时间戳 | 2025-10-23T10:30:45+08:00 |
| level | 日志级别 | INFO, ERROR |
| message | 日志消息 | received hello request |
| service_id | 服务实例 ID | shunfeng-parcel-001 |
| trace_id | 链路追踪 ID | abc123def456 |
| span_id | Span ID | 789ghi |
| caller | 调用位置 | service/greeter.go:25 |
| fields | 自定义字段 | {"user_id": "123"} |

## 日志轮转

日志文件会自动轮转：
- 当文件大小超过配置的 `max_size` 时自动轮转
- 旧文件会添加时间戳后缀，如：`app.log.20251023-103045`

## 日志目录结构

```
logs/
├── app.log                    # 当前日志文件
├── app.log.20251023-103045   # 历史日志文件
├── app.log.20251023-150230
└── error.log                  # 错误日志（可选）
```

## 性能优化建议

1. **生产环境使用 INFO 级别**：避免过多 DEBUG 日志影响性能
2. **异步写入**：日志写入不阻塞业务逻辑
3. **合理使用字段**：避免记录大对象或敏感信息
4. **定期清理**：设置日志保留策略，避免磁盘占满

## 敏感信息处理

**不要记录以下信息：**
- 密码
- Token
- 身份证号
- 银行卡号
- 完整手机号（可脱敏：138****1234）

```go
// ❌ 错误示例
log.Infow("msg", "user login", "password", password)

// ✅ 正确示例
log.Infow("msg", "user login", "username", username)
```

## 告警规则示例

### 错误率告警
```
当 5 分钟内错误日志数量 > 100 时触发告警
```

### 慢请求告警
```
当请求耗时 > 3000ms 时触发告警
```

### 服务异常告警
```
当出现 FATAL 级别日志时立即告警
```

## 下一步：ELK 集成

参考 [ELK 部署指南](./elk-deployment.md) 完成日志收集和分析系统的搭建。
