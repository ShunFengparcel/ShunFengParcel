# Jaeger 连接验证指南

## 如何判断是否连接成功？

### 方法 1：查看启动日志（最直接）

启动服务后，查看控制台输出：

#### ✅ 连接成功的日志：
```
========================================
正在初始化 Jaeger 链路追踪...
Jaeger 服务器: 14.103.153.242:4318
协议: OTLP HTTP
服务名称: ShunFengParcel
========================================
✅ OTLP Exporter 创建成功
✅ 资源信息创建成功
✅ TracerProvider 已设置
正在发送测试 span 验证连接...
✅ 测试 span 已发送
========================================
🎉 Jaeger 链路追踪初始化成功！
   - 采样率: 100% (所有请求都会被追踪)
   - 查看链路: http://14.103.153.242:16686
   - 服务名称: ShunFengParcel
   - 测试 span: jaeger-connection-test
========================================
💡 提示: 等待 5-10 秒后在 Jaeger UI 中查看测试链路
========================================
```

#### ❌ 连接失败的日志：
```
========================================
正在初始化 Jaeger 链路追踪...
Jaeger 服务器: 14.103.153.242:4318
协议: OTLP HTTP
服务名称: ShunFengParcel
========================================
❌ Jaeger 连接失败: connection refused
   请检查:
   1. Jaeger 服务器是否运行
   2. 网络连接是否正常
   3. 端口 4318 是否开放
   服务将继续运行，但不会上报链路数据
```

---

### 方法 2：查看 Jaeger UI 中的测试链路

1. **等待 5-10 秒**（让数据上报）

2. **打开 Jaeger UI**：
   ```
   http://14.103.153.242:16686
   ```

3. **查找测试链路**：
   - Service 下拉框选择：`ShunFengParcel`
   - Operation 下拉框选择：`jaeger-connection-test`
   - 点击 `Find Traces`

4. **验证结果**：
   - ✅ 如果看到 `jaeger-connection-test` 链路，说明连接成功
   - ❌ 如果没有任何数据，说明连接失败

---

### 方法 3：发送真实请求验证

#### 步骤 1：发送测试请求
```bash
curl -X POST http://localhost:8000/admin/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test_user","password":"test_pass"}'
```

#### 步骤 2：查看日志
日志中应该包含 trace.id 和 span.id：
```
ts=2025-10-23T... trace.id=abc123... span.id=def456... msg=...
```

#### 步骤 3：在 Jaeger UI 中查找
- 使用日志中的 `trace.id` 在 Jaeger UI 搜索
- 或者在 Service 列表中查看最新的 traces

---

### 方法 4：使用 PowerShell 脚本自动验证

运行测试脚本：
```powershell
.\scripts\test-jaeger.ps1
```

脚本会自动：
1. ✅ 检查端口连通性
2. ✅ 发送测试请求
3. ✅ 等待数据上报
4. ✅ 打开 Jaeger UI

---

## 常见问题排查

### Q1: 日志显示连接成功，但 UI 中看不到数据？

**可能原因**：
1. 数据还在上报中（等待 5-10 秒）
2. Jaeger UI 时间范围不对
3. Service 名称选择错误

**解决方案**：
```bash
# 1. 等待 10 秒
sleep 10

# 2. 刷新 Jaeger UI 页面

# 3. 确认 Service 选择的是 "ShunFengParcel"

# 4. 调整时间范围为 "Last 15 minutes"
```

---

### Q2: 日志显示连接失败？

**检查清单**：

1. **检查网络连通性**：
   ```powershell
   Test-NetConnection -ComputerName 14.103.153.242 -Port 4318
   ```
   
2. **检查 Jaeger 服务器状态**：
   ```bash
   # 访问 Jaeger UI
   http://14.103.153.242:16686
   ```
   
3. **检查防火墙**：
   - 确保 4318 端口未被防火墙阻止
   
4. **检查 Jaeger 配置**：
   - 确认 Jaeger 启用了 OTLP HTTP receiver

---

### Q3: 如何确认数据正在实时上报？

**实时监控方法**：

1. **查看服务日志**：
   每个请求都会生成 trace.id 和 span.id

2. **使用 Jaeger UI 实时刷新**：
   - 设置时间范围为 "Last 5 minutes"
   - 每隔几秒刷新页面
   - 观察 trace 数量是否增加

3. **查看 Jaeger 服务器日志**（如果有权限）：
   ```bash
   # 查看 Jaeger collector 日志
   docker logs -f jaeger-collector
   ```

---

## 验证步骤总结

### 快速验证（1 分钟）

```bash
# 1. 启动服务
go run ./cmd/ShunFengParcel

# 2. 查看启动日志，确认看到：
#    ✅ OTLP Exporter 创建成功
#    ✅ 测试 span 已发送

# 3. 等待 10 秒

# 4. 打开 Jaeger UI
#    http://14.103.153.242:16686

# 5. 查找 "jaeger-connection-test" 链路
```

### 完整验证（3 分钟）

```bash
# 1. 启动服务
go run ./cmd/ShunFengParcel

# 2. 发送测试请求
curl -X POST http://localhost:8000/admin/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"123456"}'

# 3. 查看响应日志中的 trace.id

# 4. 在 Jaeger UI 中搜索该 trace.id

# 5. 查看完整的调用链路
```

---

## 成功标志

当你看到以下所有标志时，说明 Jaeger 集成完全成功：

- ✅ 启动日志显示 "🎉 Jaeger 链路追踪初始化成功！"
- ✅ Jaeger UI 中能看到 "jaeger-connection-test" 链路
- ✅ 发送请求后，日志包含 trace.id 和 span.id
- ✅ Jaeger UI 中能搜索到对应的 trace
- ✅ 链路数据包含完整的调用信息（耗时、状态码等）

---

## 监控建议

### 开发环境
- 保持 100% 采样率
- 实时查看 Jaeger UI
- 关注错误链路

### 生产环境
- 降低采样率到 10-20%
- 设置告警规则
- 定期检查慢请求

---

## 相关命令

```powershell
# 检查端口
Test-NetConnection -ComputerName 14.103.153.242 -Port 4318

# 查看服务进程
Get-Process | Where-Object {$_.ProcessName -like "*ShunFeng*"}

# 停止服务
Stop-Process -Name ShunFengParcel

# 运行测试脚本
.\scripts\test-jaeger.ps1
```

---

## 技术支持

如果按照以上步骤仍然无法验证连接，请检查：

1. Jaeger 服务器版本（需要 v1.35+ 支持 OTLP）
2. 网络策略和安全组配置
3. Jaeger 的 OTLP receiver 配置
4. 服务器防火墙规则

更多信息请参考：
- [Jaeger 官方文档](https://www.jaegertracing.io/docs/)
- [OTLP 协议文档](https://opentelemetry.io/docs/specs/otlp/)
