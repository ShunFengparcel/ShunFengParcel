# Jaeger 分布式链路追踪实战案例

## 📖 使用场景模拟

### 业务背景
顺丰快递管理系统（ShunFengParcel）是一个微服务架构的物流管理平台。

某天凌晨 2 点，运维收到告警：**用户下单接口响应时间从 200ms 飙升到 3000ms**。

运维困境：
- ❌ 不知道哪个环节慢了
- ❌ 无法定位具体服务
- ❌ 日志分散难以关联
- ❌ 无法重现问题

---

## 🎯 一、讲需求

### 核心需求
需要一个能够**追踪完整请求链路**的系统。

#### 功能需求
- ✅ 全链路追踪：从请求到响应
- ✅ 性能分析：定位每个环节耗时
- ✅ 依赖关系：可视化服务调用
- ✅ 错误定位：快速找到问题
- ✅ 日志关联：通过 trace.id 关联

#### 非功能需求
- ⚡ 低侵入性：不修改业务代码
- ⚡ 低性能开销：< 2%
- ⚡ 高可用性：追踪故障不影响业务
- ⚡ 易于使用：快速上手

### 具体场景

**场景：用户下单慢**
```
用户下单 (3200ms)
  ├─ 查询用户 (50ms)
  ├─ 计算运费 (30ms)
  ├─ 调用支付 (2800ms) ← 慢！
  │   └─ 第三方支付 API (2600ms) ← 问题在这！
  └─ 创建订单 (320ms)
```

**需求**：精确定位到"第三方支付 API"慢。

---

## 🔥 二、讲难点

### 难点 1：分布式上下文传播
**问题**：如何在多个服务间传递 trace 信息？

```
服务 A (trace.id=abc123)
  ↓ HTTP
服务 B (如何知道 trace.id?)
  ↓ gRPC
服务 C (如何知道 trace.id?)
```

### 难点 2：性能开销
每个请求需要：
1. 创建 span 对象
2. 记录时间
3. 收集属性
4. 序列化数据
5. 网络传输

**挑战**：如何减少开销？

### 难点 3：采样策略
```
日请求量：1000 万
每个 trace：2KB
每天数据：20GB
每月数据：600GB
```

**挑战**：如何在保证可追踪的前提下降低数据量？

### 难点 4：连接可靠性
**问题**：Jaeger 故障不能影响业务

```go
if jaegerDown {
    // 业务必须继续运行
    // 但链路数据会丢失
}
```

### 难点 5：时钟同步
分布式系统中各服务器时钟可能不同步。

---

## 💡 三、讲解决

### 3.1 技术选型

#### 为什么选择 Jaeger？

| 对比项 | Jaeger | Zipkin | SkyWalking |
|--------|--------|--------|------------|
| 协议 | ✅ OTLP | Zipkin | 私有 |
| 性能 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐ |
| UI | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| 学习成本 | 低 | 低 | 中 |

**选择理由**：
1. ✅ 支持 OpenTelemetry 标准
2. ✅ 性能优秀
3. ✅ UI 友好
4. ✅ CNCF 毕业项目

### 3.2 核心实现

#### 初始化 TracerProvider

```go
func initJaeger() func(context.Context) error {
    // 1. 创建 OTLP exporter
    exp, _ := otlptracehttp.New(ctx,
        otlptracehttp.WithInsecure(),
        otlptracehttp.WithEndpoint("14.103.153.242:4318"),
        otlptracehttp.WithURLPath("/v1/traces"),
    )
    
    // 2. 创建资源信息
    res, _ := resource.New(ctx,
        resource.WithAttributes(
            semconv.ServiceName("ShunFengParcel"),
        ),
    )
    
    // 3. 创建 TracerProvider
    tp := tracesdk.NewTracerProvider(
        tracesdk.WithBatcher(exp),      // 批量上报
        tracesdk.WithResource(res),
        tracesdk.WithSampler(tracesdk.AlwaysSample()),
    )
    
    otel.SetTracerProvider(tp)
    return tp.Shutdown
}
```

**关键点**：
- **WithBatcher**：批量上报，减少网络开销
- **WithSampler**：采样策略（开发 100%，生产 10%）

#### 自动追踪

Kratos 框架自动集成 OpenTelemetry：

```go
// HTTP Server 自动创建 span
srv := http.NewServer(opts...)
// 每个请求自动追踪
```

**自动收集**：
- http.method
- http.url
- http.status_code
- 请求/响应大小

#### 日志关联

```go
logger := log.With(log.NewStdLogger(os.Stdout),
    "trace.id", tracing.TraceID(),
    "span.id", tracing.SpanID(),
)
```

**效果**：
```
trace.id=abc123 span.id=xyz789 msg="处理订单"
```

在 Jaeger 中搜索 trace.id 即可看到完整链路。

#### 优雅降级

```go
if err != nil {
    log.Errorf("Jaeger 连接失败，服务继续运行")
    return func(ctx context.Context) error { return nil }
}
```

**策略**：
- ✅ Jaeger 故障不影响启动
- ✅ 上报失败不阻塞请求
- ✅ 自动重试连接

#### 连接验证

```go
// 启动时发送测试 span
testSpan := tracer.Start(ctx, "jaeger-connection-test")
testSpan.End()
tp.ForceFlush(ctx)
```

在 Jaeger UI 查找 `jaeger-connection-test` 验证连接。

### 3.3 实际效果

#### 快速定位慢请求

**Jaeger 显示**：
```
POST /order/create (3200ms)
  ├─ 验证用户 (50ms)
  ├─ 调用支付 (2800ms) ← 慢！
  │   └─ 微信支付 API (2600ms) ← 问题！
  └─ 创建订单 (320ms)
```

**结论**：微信支付 API 慢，需要优化或增加超时。

#### 错误追踪

```
POST /admin/register ❌
  └─ 数据库插入 ❌
      Error: Duplicate entry 'admin'
```

**结论**：用户名重复，需要增加唯一性检查。

---

## 🚀 四、讲拓展

### 4.1 高级功能

#### 自定义 Span

```go
func ProcessOrder(ctx context.Context, orderID string) error {
    tracer := otel.Tracer("ShunFengParcel")
    ctx, span := tracer.Start(ctx, "process-order")
    defer span.End()
    
    span.SetAttributes(
        attribute.String("order.id", orderID),
        attribute.Float64("amount", 99.99),
    )
    
    if err != nil {
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
    }
    
    return nil
}
```

#### 智能采样

```go
// 错误和慢请求 100%，其他 10%
type SmartSampler struct{}

func (s *SmartSampler) ShouldSample(p SamplingParameters) SamplingResult {
    if hasError(p) || isSlow(p) {
        return RecordAndSample
    }
    return baseSampler.ShouldSample(p)
}
```

#### 多后端支持

```go
// 同时上报到 Jaeger 和 Prometheus
exporters := []SpanExporter{
    jaegerExporter,
    prometheusExporter,
}
```

#### 告警集成

```go
// 监控慢请求并告警
if span.Duration() > 3*time.Second {
    alertManager.Send(Alert{
        Level: "warning",
        Message: "慢请求",
        TraceID: span.TraceID(),
    })
}
```

### 4.2 最佳实践

#### 采样率配置

| 环境 | 采样率 |
|------|--------|
| 开发 | 100% |
| 测试 | 100% |
| 预发布 | 50% |
| 生产（低流量） | 20% |
| 生产（高流量） | 5-10% |

#### Span 命名规范

```go
// ✅ 好的命名
"GET /api/v1/orders"
"mysql.query.orders"
"redis.get.user:123"

// ❌ 不好的命名
"handler"
"process"
"func1"
```

---

## 📊 五、总结

### 价值体现

| 维度 | 改进前 | 改进后 | 提升 |
|------|--------|--------|------|
| 故障定位 | 2-4 小时 | 5-10 分钟 | 95% ↓ |
| 性能优化 | 1 周 | 1 天 | 85% ↑ |
| 可观测性 | 20% | 95% | 375% ↑ |

### 投入产出

**投入**：
- 开发：2 天
- 学习：1 天
- 成本：$50/月
- 性能开销：< 2%

**产出**：
- 故障定位效率 ↑ 95%
- 减少损失：$10,000/月
- ROI：**200:1**

### 经验教训

**成功经验**：
1. ✅ 选择标准化协议（OTLP）
2. ✅ 优雅降级保证可用性
3. ✅ 自动化测试验证
4. ✅ 详细文档

**踩过的坑**：
1. ❌ 100% 采样数据爆炸
2. ❌ 没设超时服务卡死
3. ❌ 忘记 URLPath 连接失败
4. ❌ 时钟不同步顺序错乱

---

**作者**：ShunFengParcel 团队  
**日期**：2025-10-23  
**版本**：v1.0
