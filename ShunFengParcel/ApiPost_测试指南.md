# ApiPost 测试指南 - 顺丰快递系统

## 🚀 服务器信息
- **基础URL**: `http://127.0.0.1:8888`
- **服务状态**: 正在运行 ✅

## 📋 ApiPost 接口配置

### 1. 创建订单
**接口名称**: 创建订单  
**请求方式**: POST  
**URL**: `http://127.0.0.1:8888/createOrder`  
**Headers**:
```
Content-Type: application/json
```
**Body (raw JSON)**:
```json
{
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
}
```

### 2. 查看订单列表
**接口名称**: 订单列表  
**请求方式**: POST  
**URL**: `http://127.0.0.1:8888/orders`  
**Headers**:
```
Content-Type: application/json
```
**Body (raw JSON)**:
```json
{}
```

### 3. 快递员接单
**接口名称**: 快递员接单  
**请求方式**: POST  
**URL**: `http://127.0.0.1:8888/addTask`  
**Headers**:
```
Content-Type: application/json
```
**Body (raw JSON)**:
```json
{
  "courierId": 1
}
```

### 4. 查看任务列表
**接口名称**: 任务列表  
**请求方式**: POST  
**URL**: `http://127.0.0.1:8888/taskList`  
**Headers**:
```
Content-Type: application/json
```
**Body (raw JSON)**:
```json
{}
```

### 5. 异常处理（核心功能）
**接口名称**: 用户取消订单  
**请求方式**: POST  
**URL**: `http://127.0.0.1:8888/handleException`  
**Headers**:
```
Content-Type: application/json
```
**Body (raw JSON)**:
```json
{
  "order_id": 1,
  "exception_type": 1,
  "reason": "用户主动取消订单",
  "courier_id": 1
}
```

**异常类型说明**:
- `exception_type: 1` - 用户取消订单
- `exception_type: 2` - 快递员取消订单
- `exception_type: 3` - 其他异常

## 🎯 ApiPost 测试步骤

### 步骤1: 创建接口集合
1. 打开ApiPost
2. 新建项目："顺丰快递系统测试"
3. 设置环境变量：
   - 变量名: `baseUrl`
   - 变量值: `http://127.0.0.1:8888`

### 步骤2: 添加接口
按照上面的配置，依次添加5个接口到你的项目中。

### 步骤3: 完整测试流程
1. **创建订单** → 复制返回的 `orderId`
2. **查看订单列表** → 确认订单存在且状态为"pending"
3. **快递员接单** → 任务状态变为"accepted"
4. **查看任务列表** → 确认任务被接取
5. **异常处理测试**:
   - 将步骤1返回的 `orderId` 填入异常处理接口的 `order_id` 字段
   - 发送请求测试用户取消订单

## 📊 预期响应结果

### 创建订单响应:
```json
{
  "orderId": 18
}
```

### 快递员接单响应:
```json
{
  "task": {
    "id": "17",
    "courierId": "1",
    "courierName": "",
    "taskType": "Pickup",
    "priority": 1,
    "taskStatus": "Accepted",
    "createdAt": ""
  }
}
```

### 异常处理响应:
**未出发场景**:
```json
{
  "success": true,
  "message": "订单取消成功",
  "newStatus": "ORDER_STATUS_CANCELLED",
  "requiresNegotiation": false,
  "suggestedFee": 0
}
```

**已出发场景**:
```json
{
  "success": true,
  "message": "快递员已出发，需要协商空跑费",
  "newStatus": "ORDER_STATUS_CANCELLED",
  "requiresNegotiation": true,
  "suggestedFee": 5
}
```

## 🔧 ApiPost 使用技巧

### 1. 使用变量
在ApiPost中设置变量来简化测试：
- `{{baseUrl}}` 替代 `http://127.0.0.1:8888`
- `{{orderId}}` 保存创建订单返回的ID

### 2. 后置脚本
在"创建订单"接口的后置脚本中添加：
```javascript
// 保存订单ID到环境变量
const response = JSON.parse(responseBody);
if (response.orderId) {
    apt.setEnvironmentVariable("orderId", response.orderId);
}
```

### 3. 前置脚本
在"异常处理"接口的前置脚本中添加：
```javascript
// 自动使用保存的订单ID
const orderId = apt.getEnvironmentVariable("orderId");
if (orderId) {
    const body = JSON.parse(requestBody);
    body.order_id = parseInt(orderId);
    requestBody = JSON.stringify(body);
}
```

## ⚠️ 注意事项
1. 确保服务器在 `http://127.0.0.1:8888` 上运行
2. 所有接口都使用 POST 方法
3. 必须设置 `Content-Type: application/json`
4. `order_id` 必须是实际存在的订单ID
5. 测试异常处理前，确保订单已被快递员接取（状态为accepted）

## 🚀 快速开始
1. 复制上述接口配置到ApiPost
2. 按顺序执行：创建订单 → 快递员接单 → 异常处理
3. 观察不同场景下的响应结果