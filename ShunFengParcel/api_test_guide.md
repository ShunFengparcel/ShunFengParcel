# 顺丰快递系统 API 测试指南

## 服务器信息
- **HTTP服务地址**: http://127.0.0.1:8888
- **gRPC服务地址**: 127.0.0.1:8089
- **服务状态**: 正在运行 ✅

## 主要接口列表

### 1. 创建订单
**接口**: `POST /createOrder`
**功能**: 创建新的快递订单

**请求示例**:
```bash
curl -X POST http://127.0.0.1:8888/createOrder \
  -H "Content-Type: application/json" \
  -d '{
    "senderName": "张三",
    "receiverName": "李四", 
    "senderPhone": "13800138000",
    "receiverPhone": "13900139000",
    "senderAddress": "{\"city\": \"深圳\", \"detail\": \"南山区科技园\", \"district\": \"南山\", \"province\": \"广东\"}",
    "receiverAddress": "{\"city\": \"北京\", \"detail\": \"朝阳区CBD\", \"district\": \"朝阳\", \"province\": \"北京\"}",
    "serviceType": "same_day",
    "productType": "document", 
    "estimatedWeight": 1.0,
    "isInsured": 0,
    "paymentMethod": "sender_pay",
    "userId": "1"
  }'
```

**PowerShell示例**:
```powershell
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/createOrder" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"senderName": "张三", "receiverName": "李四", "senderPhone": "13800138000", "receiverPhone": "13900139000", "senderAddress": "{\"city\": \"深圳\", \"detail\": \"南山区科技园\", \"district\": \"南山\", \"province\": \"广东\"}", "receiverAddress": "{\"city\": \"北京\", \"detail\": \"朝阳区CBD\", \"district\": \"朝阳\", \"province\": \"北京\"}", "serviceType": "same_day", "productType": "document", "estimatedWeight": 1.0, "isInsured": 0, "paymentMethod": "sender_pay", "userId": "1"}').Content
```

### 2. 查看订单列表
**接口**: `POST /orders`
**功能**: 获取所有订单列表

**请求示例**:
```bash
curl -X POST http://127.0.0.1:8888/orders \
  -H "Content-Type: application/json" \
  -d '{}'
```

**PowerShell示例**:
```powershell
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/orders" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{}').Content
```

### 3. 快递员接单
**接口**: `POST /addTask`
**功能**: 快递员从任务队列中接取任务

**请求示例**:
```bash
curl -X POST http://127.0.0.1:8888/addTask \
  -H "Content-Type: application/json" \
  -d '{"courierId": 1}'
```

**PowerShell示例**:
```powershell
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/addTask" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"courierId": 1}').Content
```

### 4. 查看任务列表
**接口**: `POST /taskList`
**功能**: 获取所有快递员任务列表

**请求示例**:
```bash
curl -X POST http://127.0.0.1:8888/taskList \
  -H "Content-Type: application/json" \
  -d '{}'
```

**PowerShell示例**:
```powershell
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/taskList" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{}').Content
```

### 5. 异常处理（重点功能）
**接口**: `POST /handleException`
**功能**: 处理订单异常情况，如用户取消、快递员取消等

**用户取消订单示例**:
```bash
curl -X POST http://127.0.0.1:8888/handleException \
  -H "Content-Type: application/json" \
  -d '{
    "order_id": 1,
    "exception_type": 1,
    "reason": "用户主动取消订单",
    "courier_id": 1
  }'
```

**PowerShell示例**:
```powershell
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/handleException" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"order_id": 1, "exception_type": 1, "reason": "用户主动取消订单", "courier_id": 1}').Content
```

**异常类型说明**:
- `exception_type: 1` - 用户取消订单
- `exception_type: 2` - 快递员取消订单  
- `exception_type: 3` - 其他异常

## 测试流程建议

### 完整测试流程:
1. **创建订单** → 获得订单ID
2. **查看订单列表** → 确认订单状态为"pending"
3. **快递员接单** → 任务状态变为"accepted"
4. **测试异常处理**:
   - 未出发场景：直接取消pending状态的订单
   - 已出发场景：取消accepted状态的订单（需要协商空跑费）

### 快速测试命令（PowerShell）:
```powershell
# 1. 创建订单
$orderId = ((Invoke-WebRequest -Uri "http://127.0.0.1:8888/createOrder" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"senderName": "测试用户", "receiverName": "收件人", "senderPhone": "13800000000", "receiverPhone": "13900000000", "senderAddress": "{\"city\": \"深圳\", \"detail\": \"测试地址\", \"district\": \"南山\", \"province\": \"广东\"}", "receiverAddress": "{\"city\": \"北京\", \"detail\": \"测试地址\", \"district\": \"朝阳\", \"province\": \"北京\"}", "serviceType": "same_day", "productType": "document", "estimatedWeight": 1.0, "isInsured": 0, "paymentMethod": "sender_pay", "userId": "1"}').Content | ConvertFrom-Json).orderId

# 2. 快递员接单
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/addTask" -Method POST -Headers @{"Content-Type"="application/json"} -Body '{"courierId": 1}').Content

# 3. 测试用户取消（已出发场景）
(Invoke-WebRequest -Uri "http://127.0.0.1:8888/handleException" -Method POST -Headers @{"Content-Type"="application/json"} -Body "{`"order_id`": $orderId, `"exception_type`": 1, `"reason`": `"用户主动取消订单`", `"courier_id`": 1}").Content
```

## 预期响应示例

### 异常处理响应:
- **未出发场景**: `{"success":true, "message":"订单取消成功", "newStatus":"ORDER_STATUS_CANCELLED", "requiresNegotiation":false, "suggestedFee":0}`
- **已出发场景**: `{"success":true, "message":"快递员已出发，需要协商空跑费", "newStatus":"ORDER_STATUS_CANCELLED", "requiresNegotiation":true, "suggestedFee":5}`

## 注意事项
- 所有接口都使用POST方法
- Content-Type必须设置为"application/json"
- 服务器运行在本地8888端口
- 异常处理逻辑会根据任务状态自动判断是否需要协商费用