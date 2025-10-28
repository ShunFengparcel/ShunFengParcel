# TabBar 图标修复说明

## 问题

底部 TabBar 的图标显示为占位符（图片图标），而不是实际的图标。

## 原因

当前的图标文件可能是占位符或者格式不正确。

## 解决方案

### 方案1：使用 iconfont 图标（推荐）

微信小程序支持使用 iconfont 字体图标，这样更清晰且不需要图片文件。

#### 步骤：

1. **访问 iconfont.cn**
   - 网址：https://www.iconfont.cn/
   - 搜索需要的图标：首页、搜索、发送、用户

2. **下载图标**
   - 选择合适的图标加入购物车
   - 下载为 PNG 格式
   - 尺寸：81x81 像素（微信小程序推荐）

3. **替换图标文件**
   将下载的图标重命名并放到 `web/user-client/static/tabbar/` 目录：
   - `home.png` - 首页（未选中）
   - `home-active.png` - 首页（选中，红色）
   - `search.png` - 查快递（未选中）
   - `search-active.png` - 查快递（选中，红色）
   - `send.png` - 寄快递（未选中）
   - `send-active.png` - 寄快递（选中，红色）
   - `user.png` - 我的（未选中）
   - `user-active.png` - 我的（选中，红色）

### 方案2：使用 emoji 临时替代

如果暂时无法获取图标，可以先用 emoji 表情：

修改 `pages.json`：

```json
{
  "tabBar": {
    "color": "#7A7E83",
    "selectedColor": "#D81E06",
    "borderStyle": "black",
    "backgroundColor": "#FFFFFF",
    "list": [
      {
        "pagePath": "pages/index/index",
        "text": "🏠 首页"
      },
      {
        "pagePath": "pages/order/list",
        "text": "🔍 查快递"
      },
      {
        "pagePath": "pages/send/index",
        "text": "📦 寄快递"
      },
      {
        "pagePath": "pages/profile/index",
        "text": "👤 我的"
      }
    ]
  }
}
```

注意：这种方式会移除 iconPath，只显示文字和 emoji。

### 方案3：创建简单的纯色图标

我可以帮你创建简单的纯色图标。

## 推荐的图标规格

- **尺寸**：81x81 像素
- **格式**：PNG
- **背景**：透明
- **颜色**：
  - 未选中：灰色 (#7A7E83)
  - 选中：红色 (#D81E06)

## 快速测试

### 临时方案：移除图标，只显示文字

修改 `web/user-client/pages.json`，移除所有 `iconPath` 和 `selectedIconPath`：

```json
{
  "tabBar": {
    "color": "#7A7E83",
    "selectedColor": "#D81E06",
    "borderStyle": "black",
    "backgroundColor": "#FFFFFF",
    "list": [
      {
        "pagePath": "pages/index/index",
        "text": "首页"
      },
      {
        "pagePath": "pages/order/list",
        "text": "查快递"
      },
      {
        "pagePath": "pages/send/index",
        "text": "寄快递"
      },
      {
        "pagePath": "pages/profile/index",
        "text": "我的"
      }
    ]
  }
}
```

这样至少可以先让功能正常使用，图标可以后续再添加。

## 重新编译

修改后需要重新编译：

```bash
cd web/user-client
npm run dev:mp-weixin
```

然后在微信开发者工具中点击"编译"。

## 在线图标资源

- **iconfont**: https://www.iconfont.cn/
- **iconpark**: https://iconpark.oceanengine.com/
- **flaticon**: https://www.flaticon.com/

搜索关键词：
- home / 首页
- search / 搜索
- send / 发送 / 快递
- user / 用户 / 个人中心
