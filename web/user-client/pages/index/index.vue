<template>
  <view class="container">
    <!-- 测试登录按钮 -->
    <view class="test-login-btn" @click="goToLogin">
      <text>🔑 点击这里测试登录</text>
    </view>

    <!-- 顶部区域 -->
    <view class="header">
      <view class="location">
        <text class="location-text">上海</text>
        <text class="location-arrow">▼</text>
      </view>
      <view class="search-bar" @click="handleSearch">
        <text class="search-icon">🔍</text>
        <text class="search-placeholder">搜索运单或服务</text>
      </view>
      <view class="icons">
        <view class="icon-btn message-btn">
          <text class="icon">💬</text>
          <view class="badge" v-if="unreadCount > 0">{{ unreadCount }}</view>
        </view>
        <view class="icon-btn notice-btn">
          <text class="icon">📢</text>
        </view>
      </view>
    </view>

    <!-- 主要服务按钮 -->
    <view class="main-services">
      <view class="service-btn express-btn" hover-class="service-btn-hover" @click="navigateTo('/pages/send/create')">
        <view class="service-content">
          <text class="service-title">快速寄件</text>
          <text class="service-subtitle">一小时上门取件</text>
        </view>
        <text class="service-badge">SF</text>
      </view>
      <view class="service-btn logistics-btn" hover-class="service-btn-hover" @click="handleServiceClick">
        <view class="service-content">
          <text class="service-title">发物流</text>
        </view>
        <text class="service-icon">🚚</text>
      </view>
      <view class="service-btn scan-btn" hover-class="service-btn-hover" @click="handleScanClick">
        <view class="service-content">
          <text class="service-title">扫一扫</text>
        </view>
        <text class="service-icon">📱</text>
      </view>
    </view>

    <!-- 快捷入口 -->
    <view class="quick-entries">
      <view class="entry-item" hover-class="entry-item-hover">
        <text class="entry-icon">⏰</text>
        <text class="entry-text">运送时效</text>
      </view>
      <view class="entry-item" hover-class="entry-item-hover">
        <text class="entry-icon">💰</text>
        <text class="entry-text">收费标准</text>
      </view>
      <view class="entry-item" hover-class="entry-item-hover">
        <text class="entry-icon">🎧</text>
        <text class="entry-text">客服中心</text>
      </view>
    </view>

    <!-- 服务网格 -->
    <view class="service-grid">
      <view class="grid-item" hover-class="grid-item-hover" v-for="(service, index) in services" :key="index">
        <view class="grid-icon" :style="{ backgroundColor: service.color }">
          <text>{{ service.icon }}</text>
        </view>
        <text class="grid-label">{{ service.label }}</text>
        <view class="grid-badge" v-if="service.badge">{{ service.badge }}</view>
      </view>
    </view>

    <!-- 近期寄件和收件 -->
    <view class="recent-section">
      <view class="section-header">
        <text class="section-title">近期寄件</text>
        <text class="section-more" @click="navigateTo('/pages/order/list')">全部快递 ></text>
      </view>

      <!-- 订单列表 -->
      <view v-if="recentOrders.length > 0" class="order-list">
        <view v-for="order in recentOrders" :key="order.id" class="order-item" @click="goToOrderDetail(order.id)">
          <view class="order-header">
            <view class="order-status">
              <text class="status-icon">{{ getStatusIcon(order.status) }}</text>
              <text class="status-text">{{ getStatusText(order.status) }}</text>
            </view>
            <text class="order-time">{{ formatTime(order.created_at) }}</text>
          </view>

          <view class="order-route">
            <view class="route-item">
              <text class="route-label">寄</text>
              <text class="route-address">{{ order.sender_address || formatAddress(order) }}</text>
            </view>
            <text class="route-arrow">→</text>
            <view class="route-item">
              <text class="route-label">收</text>
              <text class="route-address">{{ order.receiver_address || formatReceiverAddress(order) }}</text>
            </view>
          </view>

          <view class="order-footer">
            <text class="order-no">运单号：{{ order.order_no }}</text>
            <text class="order-fee">¥{{ order.estimated_fee }}</text>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <view v-else class="empty-orders">
        <text class="empty-icon">📦</text>
        <text class="empty-text">暂无寄件记录</text>
        <view class="empty-btn" @click="navigateTo('/pages/send/create')">
          <text>立即寄件</text>
        </view>
      </view>
    </view>

    <!-- 优惠券横幅 -->
    <view class="coupon-banner">
      <view class="coupon-content">
        <text class="coupon-title">30元免单券</text>
        <text class="coupon-subtitle">天天送</text>
        <view class="coupon-image">🎁</view>
        <view class="coupon-btn">立即参与 ></view>
      </view>
      <view class="promo-items">
        <view class="promo-item">
          <text class="promo-icon">💰</text>
          <text class="promo-text">本月权益待领取 HOT</text>
          <text class="promo-desc">领积蓄权益</text>
        </view>
        <view class="promo-item">
          <text class="promo-icon">🎁</text>
          <text class="promo-text">兑换好礼</text>
        </view>
      </view>
      <view class="card-banner">
        <text class="card-title">新速运通卡</text>
        <text class="card-desc">8% 近利高达</text>
        <text class="card-discount">充值返利最高8%</text>
        <text class="card-value">去充值</text>
      </view>
    </view>

    <!-- 特色服务 -->
    <view class="special-services">
      <view class="special-header">
        <text class="special-title">特色服务</text>
        <text class="special-subtitle">顺丰服务 品质保证</text>
      </view>
      <view class="special-items">
        <view class="special-item">
          <view class="special-icon">📍</view>
          <view class="special-info">
            <text class="special-name">同城跑腿</text>
            <text class="special-desc">平均1小时送全城 | 新客3折起</text>
          </view>
        </view>
        <view class="special-item">
          <view class="special-icon">🚚</view>
          <view class="special-info">
            <text class="special-name">搬家搬厂</text>
            <text class="special-desc">省心安全放心 | 4万+人的已选</text>
            <text class="special-tag">去预约</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useOrderStore } from '@/stores/order'

const userStore = useUserStore()
const orderStore = useOrderStore()
const unreadCount = ref(1)
const recentOrders = ref([])

// 页面加载时检查登录状态并加载订单
onMounted(async () => {
  if (!userStore.isLogin) {
    uni.reLaunch({
      url: '/pages/login/login'
    })
    return
  }

  // 加载最近的订单（最多3条）
  await loadRecentOrders()
})

// 加载最近的订单
async function loadRecentOrders() {
  try {
    const response = await orderStore.getOrderList('all', 1, 3)
    if (response && response.orders) {
      recentOrders.value = response.orders
    }
  } catch (error) {
    console.error('加载订单失败:', error)
  }
}

const services = ref([
  { icon: '📦', label: '同城跑腿', color: '#FF6B6B', badge: '低至3折' },
  { icon: '📱', label: '手机回收', color: '#4ECDC4', badge: '补贴20%' },
  { icon: '🍊', label: '生鲜水果', color: '#FFB347' },
  { icon: '👕', label: '旧衣回收', color: '#95E1D3' },
  { icon: '💳', label: '优惠充值卡', color: '#F38181' },
  { icon: '👔', label: '洗衣洗鞋', color: '#AA96DA' },
  { icon: '🎓', label: '至尊会员', color: '#FCBAD3' },
  { icon: '🎓', label: '学生会员', color: '#A8D8EA' },
  { icon: '📍', label: '亲情卡', color: '#FFD93D' },
  { icon: '🎁', label: '寄件返礼', color: '#6BCB77' }
])

const handleSearch = () => {
  // 弹出输入框让用户输入运单号
  uni.showModal({
    title: '查询快递',
    editable: true,
    placeholderText: '请输入运单号',
    success: (res) => {
      if (res.confirm && res.content) {
        // 跳转到物流轨迹页面
        uni.navigateTo({
          url: `/pages/order/tracking?trackingNo=${res.content}`
        })
      } else if (res.confirm && !res.content) {
        // 如果没有输入运单号，跳转到订单列表
        uni.navigateTo({
          url: '/pages/order/list'
        })
      }
    }
  })
}

const handleScan = () => {
  uni.scanCode({
    success: (res) => {
      console.log('扫码结果：', res)
    }
  })
}

const handleScanClick = () => {
  handleScan()
}

const handleServiceClick = () => {
  uni.showToast({
    title: '功能开发中',
    icon: 'none'
  })
}

const navigateTo = (url) => {
  uni.navigateTo({ url })
}

const goToLogin = () => {
  uni.reLaunch({
    url: '/pages/login/login'
  })
}

const goToOrderDetail = (orderId) => {
  uni.navigateTo({
    url: `/pages/order/detail?id=${orderId}`
  })
}

const getStatusIcon = (status) => {
  const icons = {
    pending: '⏰',
    confirmed: '✓',
    picked_up: '📦',
    in_transit: '🚚',
    out_for_delivery: '🚴',
    delivered: '✅',
    cancelled: '❌'
  }
  return icons[status] || '📦'
}

const getStatusText = (status) => {
  const texts = {
    pending: '待接单',
    confirmed: '已接单',
    picked_up: '已揽收',
    in_transit: '运输中',
    out_for_delivery: '派送中',
    delivered: '已签收',
    cancelled: '已取消'
  }
  return texts[status] || '未知'
}

const formatTime = (time) => {
  if (!time) return ''
  const date = new Date(time)
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}-${day}`
}

const formatAddress = (order) => {
  return `${order.sender_city || ''}${order.sender_district || ''}`
}

const formatReceiverAddress = (order) => {
  return `${order.receiver_city || ''}${order.receiver_district || ''}`
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F5F5F5;
  padding-bottom: 20rpx;
}

/* 测试登录按钮 */
.test-login-btn {
  background: linear-gradient(135deg, #D81E06 0%, #FF4444 100%);
  color: white;
  padding: 30rpx;
  text-align: center;
  font-size: 32rpx;
  font-weight: bold;
  margin: 20rpx;
  border-radius: 15rpx;
  box-shadow: 0 8rpx 20rpx rgba(216, 30, 6, 0.3);
  animation: pulse 2s infinite;
}

@keyframes pulse {

  0%,
  100% {
    transform: scale(1);
  }

  50% {
    transform: scale(1.02);
  }
}

/* 顶部区域 */
.header {
  display: flex;
  align-items: center;
  padding: 20rpx 30rpx;
  background-color: #FFFFFF;
  gap: 20rpx;
}

.location {
  display: flex;
  align-items: center;
  gap: 5rpx;
}

.location-text {
  font-size: 28rpx;
  color: #333333;
  font-weight: 600;
}

.location-arrow {
  font-size: 20rpx;
  color: #666666;
}

.search-bar {
  flex: 1;
  display: flex;
  align-items: center;
  height: 60rpx;
  background-color: #F5F5F5;
  border-radius: 30rpx;
  padding: 0 20rpx;
  gap: 10rpx;
}

.search-icon {
  font-size: 28rpx;
}

.search-placeholder {
  font-size: 26rpx;
  color: #999999;
}

.icons {
  display: flex;
  gap: 20rpx;
}

.icon-btn {
  position: relative;
  width: 50rpx;
  height: 50rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon {
  font-size: 36rpx;
}

.badge {
  position: absolute;
  top: 0;
  right: 0;
  min-width: 32rpx;
  height: 32rpx;
  background-color: #FF4D4F;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
  color: #FFFFFF;
  padding: 0 8rpx;
}

/* 主要服务按钮 */
.main-services {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 30rpx;
}

.service-btn {
  flex: 1;
  height: 160rpx;
  border-radius: 16rpx;
  padding: 20rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.express-btn {
  background: linear-gradient(135deg, #FF4D4F 0%, #FF7875 100%);
}

.logistics-btn {
  background: linear-gradient(135deg, #434343 0%, #666666 100%);
}

.scan-btn {
  background: linear-gradient(135deg, #8B7355 0%, #A0826D 100%);
}

.service-content {
  display: flex;
  flex-direction: column;
  gap: 5rpx;
}

.service-title {
  font-size: 32rpx;
  color: #FFFFFF;
  font-weight: 600;
}

.service-subtitle {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.9);
}

.service-badge {
  position: absolute;
  bottom: 15rpx;
  right: 15rpx;
  font-size: 40rpx;
  color: rgba(255, 255, 255, 0.3);
  font-weight: bold;
}

.service-icon {
  position: absolute;
  bottom: 15rpx;
  right: 15rpx;
  font-size: 50rpx;
}

.service-btn-hover {
  opacity: 0.8;
  transform: scale(0.98);
}

/* 快捷入口 */
.quick-entries {
  display: flex;
  justify-content: space-around;
  padding: 30rpx;
  background-color: #FFFFFF;
  margin: 0 30rpx;
  border-radius: 16rpx;
}

.entry-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
}

.entry-icon {
  font-size: 40rpx;
}

.entry-text {
  font-size: 24rpx;
  color: #666666;
}

.entry-item-hover {
  opacity: 0.7;
}

/* 服务网格 */
.service-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 30rpx 20rpx;
  padding: 30rpx;
  background-color: #FFFFFF;
  margin: 20rpx 30rpx;
  border-radius: 16rpx;
}

.grid-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
  position: relative;
}

.grid-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 36rpx;
}

.grid-label {
  font-size: 22rpx;
  color: #333333;
  text-align: center;
}

.grid-badge {
  position: absolute;
  top: -5rpx;
  right: -10rpx;
  background-color: #FF4D4F;
  color: #FFFFFF;
  font-size: 18rpx;
  padding: 2rpx 8rpx;
  border-radius: 8rpx;
}

.grid-item-hover {
  opacity: 0.7;
  transform: scale(0.95);
}

/* 近期寄件 */
.recent-section {
  padding: 20rpx 30rpx;
  background-color: #FFFFFF;
  margin: 20rpx 30rpx;
  border-radius: 16rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 32rpx;
  color: #333333;
  font-weight: 600;
}

.section-more {
  font-size: 26rpx;
  color: #D81E06;
}

/* 订单列表 */
.order-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.order-item {
  background-color: #F8F8F8;
  border-radius: 12rpx;
  padding: 24rpx;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.order-status {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.status-icon {
  font-size: 28rpx;
}

.status-text {
  font-size: 26rpx;
  font-weight: 600;
  color: #333333;
}

.order-time {
  font-size: 22rpx;
  color: #999999;
}

.order-route {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 16rpx;
}

.route-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.route-label {
  width: 36rpx;
  height: 36rpx;
  background-color: #E5E5E5;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20rpx;
  color: #666666;
  flex-shrink: 0;
}

.route-address {
  flex: 1;
  font-size: 24rpx;
  color: #666666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.route-arrow {
  font-size: 20rpx;
  color: #CCCCCC;
  flex-shrink: 0;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 16rpx;
  border-top: 1rpx solid #E5E5E5;
}

.order-no {
  font-size: 22rpx;
  color: #999999;
}

.order-fee {
  font-size: 28rpx;
  font-weight: 600;
  color: #D81E06;
}

/* 空状态 */
.empty-orders {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 0;
}

.empty-icon {
  font-size: 120rpx;
  opacity: 0.3;
  margin-bottom: 20rpx;
}

.empty-text {
  font-size: 26rpx;
  color: #999999;
  margin-bottom: 30rpx;
}

.empty-btn {
  background-color: #D81E06;
  color: #FFFFFF;
  padding: 16rpx 48rpx;
  border-radius: 30rpx;
  font-size: 26rpx;
}

/* 优惠券横幅 */
.coupon-banner {
  margin: 20rpx 30rpx;
  background-color: #FFFFFF;
  border-radius: 16rpx;
  padding: 30rpx;
}

.coupon-content {
  background: linear-gradient(135deg, #FFF5E6 0%, #FFE7BA 100%);
  border-radius: 16rpx;
  padding: 30rpx;
  position: relative;
  margin-bottom: 20rpx;
}

.coupon-title {
  font-size: 48rpx;
  color: #FF6B00;
  font-weight: bold;
}

.coupon-subtitle {
  font-size: 24rpx;
  color: #FF6B00;
}

.coupon-btn {
  background-color: #FF6B00;
  color: #FFFFFF;
  padding: 10rpx 30rpx;
  border-radius: 30rpx;
  font-size: 24rpx;
  display: inline-block;
  margin-top: 20rpx;
}

.promo-items {
  display: flex;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.promo-item {
  flex: 1;
  background-color: #FFF9F0;
  border-radius: 12rpx;
  padding: 20rpx;
}

.card-banner {
  background: linear-gradient(135deg, #FFF9E6 0%, #FFE7BA 100%);
  border-radius: 12rpx;
  padding: 20rpx;
}

/* 特色服务 */
.special-services {
  margin: 20rpx 30rpx;
  background-color: #FFFFFF;
  border-radius: 16rpx;
  padding: 30rpx;
}

.special-header {
  margin-bottom: 20rpx;
}

.special-title {
  font-size: 32rpx;
  color: #333333;
  font-weight: 600;
}

.special-subtitle {
  font-size: 24rpx;
  color: #999999;
  margin-left: 10rpx;
}

.special-items {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.special-item {
  display: flex;
  gap: 20rpx;
  padding: 20rpx;
  background-color: #F8F8F8;
  border-radius: 12rpx;
}

.special-icon {
  font-size: 60rpx;
}

.special-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5rpx;
}

.special-name {
  font-size: 28rpx;
  color: #333333;
  font-weight: 600;
}

.special-desc {
  font-size: 22rpx;
  color: #999999;
}

.special-tag {
  color: #FF4D4F;
  font-size: 24rpx;
}
</style>
