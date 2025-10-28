# 小程序 getApp() 错误完整解决方案

## 错误信息
```
TypeError: Cannot read property '$vm' of undefined
at getAppVm (vendor.js? [sm]:7424)
```

## 根本原因

这个错误是因为在某些组件初始化时，App 实例还没有完全创建好，导致 `getApp()` 返回 undefined。

## 已实施的修复

### 1. 更新 main.js
移除了在 createApp 时立即初始化 userStore 的代码，改为在 App.onLaunch 中初始化。

### 2. 更新 App.vue
在 onLaunch 生命周期中：
- 初始化 globalData
- 初始化 userStore
- 添加 try-catch 保护

### 3. 更新页面组件
在所有使用 `getApp()` 的地方添加 try-catch 保护：
- `pages/send/create.vue`
- `pages/send/address.vue`

## 完整的重新编译步骤

### 步骤1：停止所有进程

1. 停止 npm 编译进程（Ctrl+C）
2. 关闭微信开发者工具

### 步骤2：清理缓存

```bash
cd web/user-client

# 删除编译输出
rm -rf dist

# 删除 node_modules 缓存
rm -rf node_modules/.cache

# 如果问题严重，可以删除 node_modules 重新安装
# rm -rf node_modules
# npm install
```

### 步骤3：重新编译

```bash
npm run dev:mp-weixin
```

等待编译完成，看到类似输出：
```
DONE  Build complete. Watching for changes...
```

### 步骤4：重新打开微信开发者工具

1. 打开微信开发者工具
2. 导入项目：选择 `web/user-client/dist/dev/mp-weixin` 目录
3. 填写 AppID：`wx7f11be428011deb9`

### 步骤5：清除微信开发者工具缓存

1. 点击菜单：工具 → 清除缓存
2. 勾选所有选项：
   - 清除文件缓存
   - 清除网络缓存
   - 清除授权数据
3. 点击"清除"

### 步骤6：重新编译小程序

在微信开发者工具中：
1. 点击"编译"按钮
2. 或按快捷键 `Ctrl+B` (Windows) / `Cmd+B` (Mac)

### 步骤7：测试

1. 打开登录页面
2. 点击"微信一键登录"
3. 查看控制台是否还有错误

## 如果问题仍然存在

### 方案1：检查编译输出

查看 `web/user-client/dist/dev/mp-weixin/App.js` 文件，确认修改已经编译进去。

### 方案2：使用 HBuilderX

如果使用命令行编译有问题，可以尝试使用 HBuilderX：

1. 下载并安装 HBuilderX
2. 打开项目：`web/user-client`
3. 运行 → 运行到小程序模拟器 → 微信开发者工具

### 方案3：简化 getApp() 使用

如果上述方法都不行，可以完全移除 `getApp()` 的使用，改用 Pinia store 来传递数据。

#### 修改 create.vue

**原来的方式（使用 getApp）：**
```javascript
const app = getApp()
app.globalData.selectedAddress = address
```

**新的方式（使用 Pinia）：**
```javascript
// 在 stores/order.js 中添加
const tempAddress = ref(null)
const tempAddressType = ref('')

function setTempAddress(address, type) {
  tempAddress.value = address
  tempAddressType.value = type
}

// 在 create.vue 中使用
import { useOrderStore } from '@/stores/order'
const orderStore = useOrderStore()

// 保存地址
orderStore.setTempAddress(address, 'sender')

// 获取地址
const address = orderStore.tempAddress
```

## 验证修复

### 1. 检查控制台日志

应该看到：
```
App Launch
用户状态初始化完成
```

### 2. 测试功能

- ✅ 登录功能正常
- ✅ 创建订单页面正常打开
- ✅ 选择地址功能正常
- ✅ 查看订单列表正常

### 3. 没有错误

控制台不应该再出现：
```
TypeError: Cannot read property '$vm' of undefined
```

## 预防措施

### 1. 始终使用 try-catch

```javascript
try {
  const app = getApp()
  if (app && app.globalData) {
    // 使用 app.globalData
  }
} catch (error) {
  console.error('获取App实例失败:', error)
}
```

### 2. 优先使用 Pinia Store

对于需要跨页面传递的数据，优先使用 Pinia store 而不是 globalData。

### 3. 检查生命周期

确保在正确的生命周期中访问 App 实例：
- ✅ onLaunch
- ✅ onShow
- ✅ onReady
- ❌ setup 函数顶层（太早）

## 联系支持

如果按照以上步骤操作后问题仍然存在，请提供：

1. 完整的错误堆栈
2. 微信开发者工具版本
3. uni-app 版本
4. 操作系统版本

这将帮助我们更好地诊断问题。
