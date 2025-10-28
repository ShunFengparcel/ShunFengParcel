<template>
  <view class="container">
    <!-- 顶部搜索栏 -->
    <view class="search-bar">
      <view class="search-input" @click="handleSearch">
        <text class="search-icon">🔍</text>
        <text class="search-placeholder">搜索运单或服务</text>
      </view>
      <view class="scan-btn" @click="handleScan">
        <text class="scan-icon">📷</text>
      </view>
      <view class="message-btn">
        <text class="message-icon">💬</text>
        <view class="badge" v-if="unreadCount > 0">{{ unreadCount }}</view>
      </view>
    </view>

    <!-- 用户信息卡片 -->
    <view class="user-card">
      <view class="user-info">
        <image class="avatar" src="/static/images/default-avatar.png" mode="aspectFill"></image>
        <view class="user-detail">
          <view class="user-phone">182****4509</view>
          <view class="member-status">
            <text class="member-icon">✓</text>
            <text class="member-text">普通会员</text>
            <text class="upgrade-text">未实名 ></text>
          </view>
        </view>
        <view class="member-center-btn">
          <text>会员中心</text>
          <text class="arrow">></text>
        </view>
      </view>
      
      <view class="user-stats">
        <view class="stat-item" v-for="(item, index) in userStats" :key="index">
          <view class="stat-value">{{ item.value }}</view>
          <view class="stat-label">{{ item.label }}</view>
        </view>
      </view>
    </view>

    <!-- 会员推广卡片 -->
    <view class="promo-cards">
      <view class="promo-card promo-card-1">
        <view class="promo-content">
          <text class="promo-title">家庭账户</text>
          <text class="promo-subtitle">享10X8折互寄权益</text>
        </view>
        <view class="promo-action">去开通</view>
      </view>
      <view class="promo-card promo-card-2">
        <view class="promo-content">
          <text class="promo-title">至尊会员</text>
          <text class="promo-subtitle">全年最高省2292元</text>
        </view>
        <view class="promo-action">去升级</view>
      </view>
    </view>

    <!-- 主要服务入口 -->
    <view class="main-services">
      <view class="service-item service-express" @click="navigateTo('/pages/send/create')">
        <view class="service-icon">📦</view>
        <view class="service-info">
          <text class="service-title">快速寄件</text>
          <text class="service-subtitle">一小时上门取件</text>
        </view>
        <text class="service-badge">SF</text>
      </view>
      <view class="service-item service-logistics">
        <view class="service-icon">🚚</view>
        <view class="service-info">
          <text class="service-title">发物流</text>
          <text class="service-subtitle">大件/零担/整车</text>
        </view>
      </view>
      <view class="service-item service-scan">
        <view class="service-icon">📱</view>
        <view class="service-info">
          <text class="service-title">扫一扫</text>
        </view>
      </view>
    </view>

    <!-- 快捷入口 -->
    <view class="quick-entries">
      <view class="entry-item">
        <text class="entry-icon">⏰</text>
        <text class="entry-text">运送时效</text>
      </view>
      <view class="entry-item">
        <text class="entry-icon">💰</text>
        <text class="entry-text">收费标准</text>
      </view>
      <view class="entry-item">
        <text class="entry-icon">🎧</text>
        <text class="entry-text">客服中心</text>
      </view>
    </view>

    <!-- 增值服务网格 -->
    <view class="service-grid">
      <view class="grid-row">
        <view class="grid-item" v-for="(service, index) in services.slice(0, 5)" :key="index">
          <view class="grid-icon" :style="{ backgroundColor: service.color }">
            <text>{{ service.icon }}</text>
          </view>
          <text class="grid-label">{{ service.label }}</text>
          <view class="grid-badge" v-if="service.badge">{{ service.badge }}</view>
        </view>
      </view>
      <view class="grid-row">
        <view class="grid-item" v-for="(service, index) in services.slice(5, 10)" :key="index">
          <view class="grid-icon" :style="{ backgroundColor: service.color }">
            <text>{{ service.icon }}</text>
          </view>
          <text class="grid-label">{{ service.label }}</text>
          <view class="grid-badge" v-if="service.badge">{{ service.badge }}</view>
        </view>
      </view>
    </view>

    <!-- 在线客服和企业合作 -->
    <view class="bottom-services">
      <view class="bottom-service-item">
        <text class="bottom-icon">🎧</text>
        <text class="bottom-text">在线客服</text>
      </view>
      <view class="bottom-service-item">
        <text class="bottom-icon">🤝</text>
        <text class="bottom-text">企业合作</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'

const unreadCount = ref(1)

const userStats = ref([
  { value: '2', label: '积分' },
  { value: '0', label: '优惠券' },
  { value: '0', label: '卡券额' },
  { value: '0', label: '我的关爱包' },
  { value: '0', label: '碎能量' }
])

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
  uni.navigateTo({
    url: '/pages/order/list'
  })
}

const handleScan = () => {
  uni.scanCode({
    success: (res) => {
      console.log('扫码结果：', res)
    }
  })
}

const navigateTo = (url) => {
  uni.navigateTo({ url })
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F8F8F8;
  padding-bottom: 20rpx;
}

/* 搜索栏 */
.search-bar {
  display: flex;
  align-items: center;
  padding: 20rpx 30rpx;
  background-color: #FFFFFF;
  gap: 20rpx;
}

.search-input {
  flex: 1;
  display: flex;
  align-items: center;
  height: 70rpx;
  background-color: #F5F5F5;
  border-radius: 35rpx;
  padding: 0 30rpx;
  gap: 15rpx;
}

.search-icon {
  font-size: 32rpx;
}

.search-placeholder {
  font-size: 28rpx;
  color: #999999;
}

.scan-btn, .message-btn {
  position: relative;
  width: 70rpx;
  height: 70rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
}

.badge {
  position: absolute;
  top: 0;
  right: 0;
  min-width: 32rpx;
  height: 32rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  font-size: 20rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0 8rpx;
}

/* 用户卡片 */
.user-card {
  margin: 20rpx 30rpx;
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
}

.user-info {
  display: flex;
  align-items: center;
  margin-bottom: 30rpx;
}

.avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  margin-right: 20rpx;
}

.user-detail {
  flex: 1;
}

.user-phone {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 10rpx;
}

.member-status {
  display: flex;
  align-items: center;
  gap: 10rpx;
  font-size: 24rpx;
}

.member-icon {
  color: #52C41A;
}

.member-text {
  color: #666666;
}

.upgrade-text {
  color: #D81E06;
}

.member-center-btn {
  display: flex;
  align-items: center;
  gap: 5rpx;
  padding: 10rpx 20rpx;
  background-color: #FFF5F5;
  color: #D81E06;
  border-radius: 30rpx;
  font-size: 24rpx;
}

.user-stats {
  display: flex;
  justify-content: space-around;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 36rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 10rpx;
}

.stat-label {
  font-size: 24rpx;
  color: #999999;
}

/* 推广卡片 */
.promo-cards {
  display: flex;
  gap: 20rpx;
  padding: 0 30rpx;
  margin-bottom: 20rpx;
}

.promo-card {
  flex: 1;
  height: 160rpx;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  overflow: hidden;
}

.promo-card-1 {
  background: linear-gradient(135deg, #FFE5E5 0%, #FFD1D1 100%);
}

.promo-card-2 {
  background: linear-gradient(135deg, #FFF4E5 0%, #FFE8CC 100%);
}

.promo-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333333;
  display: block;
  margin-bottom: 10rpx;
}

.promo-subtitle {
  font-size: 22rpx;
  color: #666666;
  display: block;
}

.promo-action {
  align-self: flex-start;
  padding: 8rpx 20rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 30rpx;
  font-size: 22rpx;
}

/* 主要服务 */
.main-services {
  display: flex;
  gap: 20rpx;
  padding: 0 30rpx;
  margin-bottom: 20rpx;
}

.service-item {
  flex: 1;
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}

.service-express {
  background: linear-gradient(135deg, #FF3B30 0%, #D81E06 100%);
  color: #FFFFFF;
}

.service-logistics {
  background-color: #333333;
  color: #FFFFFF;
}

.service-scan {
  background-color: #8B6F47;
  color: #FFFFFF;
}

.service-icon {
  font-size: 60rpx;
  margin-bottom: 15rpx;
}

.service-title {
  font-size: 28rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 5rpx;
}

.service-subtitle {
  font-size: 22rpx;
  opacity: 0.8;
  display: block;
}

.service-badge {
  position: absolute;
  top: 20rpx;
  right: 20rpx;
  font-size: 20rpx;
  font-weight: 600;
}

/* 快捷入口 */
.quick-entries {
  display: flex;
  justify-content: space-around;
  background-color: #FFFFFF;
  padding: 30rpx;
  margin: 0 30rpx 20rpx;
  border-radius: 20rpx;
}

.entry-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
}

.entry-icon {
  font-size: 48rpx;
}

.entry-text {
  font-size: 24rpx;
  color: #666666;
}

/* 服务网格 */
.service-grid {
  background-color: #FFFFFF;
  padding: 30rpx;
  margin: 0 30rpx 20rpx;
  border-radius: 20rpx;
}

.grid-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 40rpx;
}

.grid-row:last-child {
  margin-bottom: 0;
}

.grid-item {
  width: 20%;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}

.grid-icon {
  width: 80rpx;
  height: 80rpx;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
  margin-bottom: 10rpx;
}

.grid-label {
  font-size: 22rpx;
  color: #666666;
  text-align: center;
}

.grid-badge {
  position: absolute;
  top: -5rpx;
  right: 10rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  font-size: 18rpx;
  padding: 2rpx 8rpx;
  border-radius: 10rpx;
}

/* 底部服务 */
.bottom-services {
  display: flex;
  gap: 20rpx;
  padding: 0 30rpx;
}

.bottom-service-item {
  flex: 1;
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15rpx;
}

.bottom-icon {
  font-size: 40rpx;
}

.bottom-text {
  font-size: 28rpx;
  color: #333333;
}
</style>
