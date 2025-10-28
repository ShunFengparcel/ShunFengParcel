# 顺丰速运用户端小程序

## 项目说明

这是顺丰速运平台的用户端小程序，使用 uni-app + Vue 3 开发，支持微信小程序平台。

## 技术栈

- **框架**: uni-app 3.x
- **UI 框架**: Vue 3
- **状态管理**: Pinia
- **构建工具**: Vite

## 目录结构

```
web/user-client/
├── pages/              # 页面文件
│   ├── index/         # 首页
│   ├── order/         # 查快递
│   ├── send/          # 寄快递
│   └── profile/       # 我的
├── components/        # 公共组件
├── stores/           # 状态管理
├── utils/            # 工具函数
├── static/           # 静态资源
├── App.vue           # 应用入口
├── main.js           # 主入口文件
├── manifest.json     # 应用配置
├── pages.json        # 页面配置
└── vite.config.js    # Vite 配置
```

## 开发指南

### 安装依赖

```bash
cd web/user-client
npm install
```

### 开发模式

```bash
npm run dev:mp-weixin
```

### 构建生产版本

```bash
npm run build:mp-weixin
```

### 微信开发者工具

1. 打开微信开发者工具
2. 导入项目，选择 `web/user-client/dist/dev/mp-weixin` 目录
3. 配置 AppID（在 manifest.json 中）

## 功能模块

### 已实现

- [x] 项目基础结构
- [x] 首页布局和基本功能
- [x] TabBar 导航

### 待实现

- [ ] 查快递功能
- [ ] 寄快递功能
- [ ] 个人中心功能
- [ ] 地址管理
- [ ] 订单管理
- [ ] 物流轨迹
- [ ] 支付功能

## API 接口

所有 API 接口预留给后端 Kratos 框架实现，接口文档参见项目根目录的设计文档。

## 注意事项

1. 开发前请先配置微信小程序 AppID
2. 需要配置服务器域名白名单
3. 地图功能需要申请高德地图 API Key
4. 支付功能需要配置微信支付商户号

## 相关文档

- [uni-app 官方文档](https://uniapp.dcloud.net.cn/)
- [Vue 3 官方文档](https://cn.vuejs.org/)
- [Pinia 官方文档](https://pinia.vuejs.org/zh/)
