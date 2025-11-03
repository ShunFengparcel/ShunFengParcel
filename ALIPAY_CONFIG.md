# 支付宝配置说明

## 公钥错误修复

本项目遇到的错误：`ERROR msg=签名验证失败:alipay: alipay public key not found`

## 问题分析

支付宝公钥用于验证支付宝返回的回调通知。如果公钥未正确配置，将导致无法验证签名。

## 解决方案

### 1. 获取支付宝公钥

请按照以下步骤从支付宝开放平台获取公钥：

1. 登录 [支付宝开放平台](https://open.alipay.com)
2. 进入 **开发者应用** -> 您的应用
3. 找到 **应用信息** 或 **开发设置**
4. 在 **密钥管理** 中，下载 **支付宝公钥** (不是应用公钥)
5. 将公钥内容复制出来

### 2. 配置公钥

#### 方式一：在配置文件中配置（推荐）

编辑 `configs/config.yaml`，在 `Alipay` 部分添加 `PublicKey` 字段：

```yaml
Alipay:
  AppId: "您的应用ID"
  Akey: "您的应用私钥"
  PublicKey: "您从支付宝获取的公钥"
```

#### 方式二：直接在代码中配置

在 `internal/service/payment.go` 和 `utils/alipay.go` 中，使用以下代码：

```go
publicKey := "您的支付宝公钥"
if err := client.LoadAliPayPublicKey(publicKey); err != nil {
    log.Errorf("加载支付宝公钥失败: %v", err)
    return nil
}
```

### 3. 公钥格式

支付宝公钥应该是 Base64 编码的字符串，不需要 PEM 格式的头尾。

正确格式示例：
```
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...（省略）
```

### 4. 相关文件

本项目已修改以下文件来支持公钥配置：

- `configs/config.yaml` - 配置文件中添加了 `PublicKey` 字段
- `internal/service/payment.go` - 在 `NewPaymentService` 中加载公钥
- `utils/alipay.go` - 在 `Alipayment` 函数中加载公钥

### 5. 验证配置

启动应用后，如果公钥配置正确，支付宝回调验证应该能够成功。

错误日志：`签名验证失败:alipay: alipay public key not found` 说明公钥未正确加载。

## 常见问题

### Q: 公钥和私钥有什么区别？
- **私钥**：用于对请求进行签名（保密）
- **公钥**：用于验证支付宝回调通知的签名（可公开）

### Q: 从哪里获取应用私钥？
在支付宝开放平台的 **密钥管理** 中，生成或查看您的 **应用私钥**。

### Q: 为什么还是收不到签名验证成功的消息？
1. 确保公钥已正确加载
2. 检查应用 ID 是否正确
3. 确保网络连接正常
4. 查看应用日志获取详细错误信息

## 参考资源

- [支付宝开放平台](https://open.alipay.com)
- [支付宝SDK文档](https://github.com/smartwalle/alipay)
- [RSA签名验签工具](https://docs.open.alipay.com/291/105971)
