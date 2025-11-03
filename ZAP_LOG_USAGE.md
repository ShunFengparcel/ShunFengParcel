# Zap 日志使用指南

## 概述

本项目使用 `go.uber.org/zap` 作为日志框架，提供高性能、结构化的日志记录功能。

---

## 功能特性

✅ **多级别日志**：Debug、Info、Warn、Error、Fatal  
✅ **多输出目标**：控制台 + 文件（分离普通日志和错误日志）  
✅ **日志轮转**：自动按大小和时间轮转，压缩旧日志  
✅ **结构化日志**：支持字段化日志，便于查询和分析  
✅ **Kratos集成**：完全兼容 Kratos 框架  
✅ **便捷API**：提供简化的日志函数和业务场景封装

---

## 日志配置

### 日志输出位置

- **控制台输出**：实时显示所有日志（彩色格式）
- **app.log**：记录 Info 及以上级别的所有日志（JSON格式）
- **error.log**：单独记录 Error 和 Fatal 级别日志（JSON格式）

### 日志轮转配置

```go
MaxSize:    100,  // 单个文件最大100MB
MaxBackups: 30,   // 最多保留30个备份文件
MaxAge:     7,    // 保留7天
Compress:   true, // 压缩旧日志文件
```

---

## 基础使用

### 1. 初始化日志系统

在 `main.go` 中初始化（已集成）：

```go
import "ShunFengParcel/utils"

func main() {
    // 初始化日志系统
    utils.InitZap()
    defer utils.Sync() // 程序退出时刷新缓冲区
    
    // ... 其他代码
}
```

### 2. 基础日志记录

#### 使用格式化函数（推荐用于简单日志）

```go
import "ShunFengParcel/utils"

// Info级别
utils.Infof("用户登录成功: userId=%d, username=%s", 123, "张三")

// Error级别
utils.Errorf("订单处理失败: orderNo=%s, error=%v", "SF123456", err)

// Debug级别
utils.Debugf("调试信息: %s", debugData)

// Warn级别
utils.Warnf("警告: 库存不足, productId=%d, stock=%d", productId, stock)
```

#### 使用结构化日志（推荐用于复杂数据）

```go
import (
    "ShunFengParcel/utils"
    "go.uber.org/zap"
)

// Info级别
utils.Info("用户登录",
    zap.Int64("user_id", 123),
    zap.String("username", "张三"),
    zap.String("ip", "192.168.1.1"),
)

// Error级别
utils.Error("数据库操作失败",
    zap.String("operation", "insert"),
    zap.String("table", "sf_orders"),
    zap.Error(err),
)
```

---

## 业务场景封装

项目提供了专门的业务日志函数，开箱即用：

### 1. 支付相关错误

```go
import "ShunFengParcel/utils"

func ProcessPayment(orderNo string) error {
    err := doPayment(orderNo)
    if err != nil {
        // 记录支付错误
        utils.LogPaymentError(orderNo, "支付宝回调验签失败", err)
        return err
    }
    return nil
}
```

### 2. 订单相关错误

```go
// 记录订单操作错误
utils.LogOrderError("SF123456", "更新订单状态", err)
```

### 3. 数据库错误

```go
// 记录数据库错误
utils.LogDatabaseError("INSERT", "sf_orders", err)
utils.LogDatabaseError("UPDATE", "sf_payments", err)
```

### 4. API调用错误

```go
// 记录第三方API调用错误
utils.LogAPIError("https://api.alipay.com/gateway", "POST", 500, err)
```

### 5. 支付宝专用日志

```go
// 记录支付宝操作错误（带详细参数）
params := map[string]interface{}{
    "app_id": "2021000148652076",
    "sign_type": "RSA2",
    "total_amount": "39.00",
}
utils.LogAlipayError("验签", "2025103122001465680508923", err, params)
```

### 6. 业务信息日志

```go
// 记录业务操作信息
data := map[string]interface{}{
    "order_no": "SF123456",
    "user_id": 1001,
    "amount": 39.00,
}
utils.LogBusinessInfo("订单模块", "创建订单", data)
```

### 7. HTTP请求日志

```go
import "time"

// 记录HTTP请求
utils.LogRequest("POST", "/api/v1/orders", "192.168.1.1", 150*time.Millisecond, 200)
```

---

## 实际应用示例

### payment.go 中使用

```go
package service

import (
    "ShunFengParcel/utils"
    "go.uber.org/zap"
)

func (s *PaymentService) HandleAlipayCallback(ctx context.Context, req *pb.UpdatePaymentRequest) (*pb.UpdatePaymentReply, error) {
    // 解析请求参数
    rawReq := req.GetReq().(*http.Request)
    if err := rawReq.ParseForm(); err != nil {
        utils.Error("解析表单参数失败",
            zap.String("remote_addr", rawReq.RemoteAddr),
            zap.Error(err),
        )
        return &pb.UpdatePaymentReply{Result: "fail"}, nil
    }
    
    // 验签
    params := rawReq.PostForm
    if err := s.alipayClient.VerifySign(params); err != nil {
        paramMap := make(map[string]interface{})
        for k, v := range params {
            if len(v) > 0 {
                paramMap[k] = v[0]
            }
        }
        utils.LogAlipayError("回调验签", params.Get("out_trade_no"), err, paramMap)
        return &pb.UpdatePaymentReply{Result: "fail"}, nil
    }
    
    // 记录成功信息
    utils.LogBusinessInfo("支付模块", "支付宝回调成功", map[string]interface{}{
        "order_no": params.Get("out_trade_no"),
        "trade_no": params.Get("trade_no"),
        "amount": params.Get("total_amount"),
        "status": params.Get("trade_status"),
    })
    
    return &pb.UpdatePaymentReply{Result: "success"}, nil
}
```

### 数据库操作日志

```go
func CreateOrder(order *config.SfOrders) error {
    if err := inits.DB.Create(order).Error; err != nil {
        utils.LogDatabaseError("CREATE", "sf_orders", err)
        return err
    }
    
    utils.Infof("订单创建成功: orderNo=%s, userId=%d", order.OrderNo, order.UserId)
    return nil
}
```

### HTTP Handler 中使用

```go
func OrderHandler(w http.ResponseWriter, r *http.Request) {
    startTime := time.Now()
    
    // 处理请求...
    
    // 记录请求日志
    utils.LogRequest(
        r.Method,
        r.URL.Path,
        r.RemoteAddr,
        time.Since(startTime),
        200,
    )
}
```

---

## 日志级别说明

### Debug
- **用途**：调试信息
- **示例**：变量值、函数调用栈、临时数据
- **是否记录到文件**：默认不记录（可配置）

### Info
- **用途**：正常业务流程信息
- **示例**：用户登录、订单创建、支付成功
- **是否记录到文件**：✅ 记录到 app.log

### Warn
- **用途**：警告信息，不影响正常运行
- **示例**：库存不足、超时重试、配置缺失
- **是否记录到文件**：✅ 记录到 app.log

### Error
- **用途**：错误信息，影响功能但不导致崩溃
- **示例**：数据库错误、API调用失败、验签失败
- **是否记录到文件**：✅ 记录到 app.log 和 error.log

### Fatal
- **用途**：致命错误，记录后程序退出
- **示例**：初始化失败、配置错误
- **是否记录到文件**：✅ 记录到 app.log 和 error.log，然后程序退出

---

## 日志查询

### 查看实时日志

```bash
# 查看所有日志
tail -f logs/app.log

# 查看错误日志
tail -f logs/error.log

# 过滤特定订单
tail -f logs/app.log | grep "SF123456"
```

### JSON日志查询

使用 `jq` 工具查询JSON格式日志：

```bash
# 查询所有错误日志
cat logs/error.log | jq 'select(.level=="error")'

# 查询特定订单的日志
cat logs/app.log | jq 'select(.order_no=="SF123456")'

# 统计错误数量
cat logs/error.log | jq -s 'length'

# 查询最近10条错误
tail -n 10 logs/error.log | jq '.'
```

---

## 性能建议

### 1. 使用结构化日志

❌ **不推荐**：
```go
utils.Infof("用户%d执行了%s操作，结果：%v", userId, action, result)
```

✅ **推荐**：
```go
utils.Info("用户操作",
    zap.Int64("user_id", userId),
    zap.String("action", action),
    zap.Any("result", result),
)
```

### 2. 避免在循环中频繁打印Debug日志

```go
// 生产环境建议关闭Debug级别
for i := 0; i < 10000; i++ {
    // ❌ 会产生大量日志
    utils.Debugf("处理第%d条数据", i)
}
```

### 3. 合理使用日志级别

- **开发环境**：Debug/Info
- **测试环境**：Info/Warn
- **生产环境**：Info/Warn/Error

---

## 配置修改

如需调整日志级别，修改 `utils/zap.go`：

```go
// 修改默认日志级别
atomicLevel.SetLevel(zap.DebugLevel) // 改为Debug级别
atomicLevel.SetLevel(zap.WarnLevel)  // 改为Warn级别
```

如需调整日志文件路径：

```go
// 修改日志目录
logDir := "/var/log/shunfeng" // 自定义路径
```

---

## 日志目录结构

```
ShunFengParcel/
├── logs/
│   ├── app.log              # 所有日志（Info及以上）
│   ├── app-2025-10-30.log.gz  # 自动压缩的旧日志
│   ├── error.log            # 错误日志（Error及以上）
│   └── error-2025-10-30.log.gz
```

---

## 常见问题

### Q1: 日志文件没有生成？

**A**: 检查程序是否有权限创建 `./logs` 目录，或者修改 `logDir` 为绝对路径。

### Q2: 日志文件过大怎么办？

**A**: 已配置自动轮转，超过100MB会自动分割。可以调整 `MaxSize` 参数。

### Q3: 如何禁用控制台输出？

**A**: 在 `InitZap()` 中注释掉控制台输出的 `zapcore.NewCore`。

### Q4: 如何输出到 syslog 或远程日志服务？

**A**: 可以添加新的 `zapcore.Core`，配置相应的 `WriteSyncer`。

---

## 总结

✅ **推荐使用场景专用函数**：如 `LogPaymentError`、`LogOrderError`  
✅ **生产环境使用 Info 级别**：避免过多Debug日志  
✅ **使用结构化日志**：便于后续查询分析  
✅ **定期清理旧日志**：虽然有自动清理，但建议定期检查

**最后**：记得在 `main.go` 退出时调用 `utils.Sync()` 确保日志完整写入！

