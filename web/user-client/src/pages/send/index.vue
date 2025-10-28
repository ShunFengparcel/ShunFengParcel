<template>
  <view class="container">
    <!-- 会员优惠券横幅 -->
    <view class="promo-banner">
      <view class="promo-content">
        <text class="promo-title">20元无门槛寄件券</text>
        <text class="promo-subtitle">至尊会员专属</text>
      </view>
      <view class="promo-image">💰</view>
      <view class="promo-btn">立即查看 ></view>
    </view>

    <!-- 会员礼包提示 -->
    <view class="member-tip">
      <text class="tip-icon">✓</text>
      <text class="tip-text">普通会员</text>
      <text class="tip-desc">领本月惊喜礼包</text>
      <view class="tip-btn">去领取</view>
    </view>

    <!-- 主要服务 -->
    <view class="main-services">
      <view class="service-card express" @click="goToCreate('express')">
        <view class="service-icon">📦</view>
        <view class="service-info">
          <text class="service-title">寄快递</text>
          <text class="service-subtitle">一小时取件</text>
        </view>
      </view>
      <view class="service-card logistics" @click="goToCreate('logistics')">
        <view class="service-icon">🚚</view>
        <view class="service-info">
          <text class="service-title">发物流</text>
          <text class="service-subtitle">大件/零担/整车</text>
        </view>
      </view>
    </view>

    <!-- 服务类型网格 -->
    <view class="service-grid">
      <view class="grid-row">
        <view 
          v-for="(service, index) in services.slice(0, 4)" 
          :key="index"
          class="grid-item"
          @click="handleServiceClick(service)"
        >
          <view class="grid-icon" :style="{ backgroundColor: service.color }">
            <text>{{ service.icon }}</text>
            <view v-if="service.badge" class="grid-badge">{{ service.badge }}</view>
          </view>
          <text class="grid-label">{{ service.label }}</text>
        </view>
      </view>
      <view class="grid-row">
        <view 
          v-for="(service, index) in services.slice(4, 8)" 
          :key="index"
          class="grid-item"
          @click="handleServiceClick(service)"
        >
          <view class="grid-icon" :style="{ backgroundColor: service.color }">
            <text>{{ service.icon }}</text>
            <view v-if="service.badge" class="grid-badge">{{ service.badge }}</view>
          </view>
          <text class="grid-label">{{ service.label }}</text>
        </view>
      </view>
      <view class="grid-row">
        <view 
          v-for="(service, index) in services.slice(8, 12)" 
          :key="index"
          class="grid-item"
          @click="handleServiceClick(service)"
        >
          <view class="grid-icon" :style="{ backgroundColor: service.color }">
            <text>{{ service.icon }}</text>
            <view v-if="service.badge" class="grid-badge">{{ service.badge }}</view>
          </view>
          <text class="grid-label">{{ service.label }}</text>
        </view>
      </view>
    </view>

    <!-- 底部服务 -->
    <view class="bottom-services">
      <view class="bottom-service-item" @click="handleBottomService('customer')">
        <text class="bottom-icon">🎧</text>
        <text class="bottom-text">在线客服</text>
      </view>
      <view class="bottom-service-item" @click="handleBottomService('enterprise')">
        <text class="bottom-icon">🤝</text>
        <text class="bottom-text">企业合作</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'

const services = ref([
  { 
    icon: '📦', 
    label: '同城急送', 
    color: '#FF6B6B', 
    badge: '低至5折',
    type: 'same_city'
  },
  { 
    icon: '📷', 
    label: '扫码寄件', 
    color: '#4ECDC4',
    type: 'scan'
  },
  { 
    icon: '📋', 
    label: '批量寄', 
    color: '#FFB347',
    type: 'batch'
  },
  { 
    icon: '🏪', 
    label: '丰巢寄件', 
    color: '#95E1D3',
    type: 'locker'
  },
  { 
    icon: '⏰', 
    label: '港澳台寄送', 
    color: '#F38181',
    type: 'hk_tw'
  },
  { 
    icon: '🌍', 
    label: '国际寄送', 
    color: '#AA96DA',
    type: 'international'
  },
  { 
    icon: '🏪', 
    label: '服务点自寄', 
    color: '#FCBAD3',
    type: 'self_service'
  },
  { 
    icon: '🎓', 
    label: '校园专区', 
    color: '#A8D8EA',
    badge: '热卖百动',
    type: 'campus'
  },
  { 
    icon: '🛍️', 
    label: '网购退货', 
    color: '#FFD93D',
    type: 'return'
  },
  { 
    icon: '🍊', 
    label: '生鲜寄', 
    color: '#6BCB77',
    type: 'fresh'
  },
  { 
    icon: '💼', 
    label: '行李寄存', 
    color: '#87CEEB',
    type: 'luggage'
  },
  { 
    icon: '✈️', 
    label: '雪具寄', 
    color: '#DDA0DD',
    badge: '热卖百动',
    type: 'ski'
  }
])

function goToCreate(type) {
  uni.navigateTo({
    url: `/pages/send/create?type=${type}`
  })
}

function handleServiceClick(service) {
  if (service.type === 'scan') {
    // 扫码寄件
    uni.scanCode({
      success: (res) => {
        console.log('扫码结果：', res)
        uni.showToast({
          title: '扫码寄件功能开发中',
          icon: 'none'
        })
      }
    })
  } else {
    uni.showToast({
      title: `${service.label}功能开发中`,
      icon: 'none'
    })
  }
}

function handleBottomService(type) {
  if (type === 'customer') {
    uni.showToast({
      title: '在线客服功能开发中',
      icon: 'none'
    })
  } else {
    uni.showToast({
      title: '企业合作功能开发中',
      icon: 'none'
    })
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: linear-gradient(180deg, #FFF4E5 0%, #F8F8F8 30%);
  padding-bottom: 20rpx;
}

/* 优惠券横幅 */
.promo-banner {
  margin: 20rpx 30rpx;
  background: linear-gradient(135deg, #FFE5E5 0%, #FFD1D1 100%);
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  position: relative;
  overflow: hidden;
}

.promo-content {
  flex: 1;
}

.promo-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #D81E06;
  display: block;
  margin-bottom: 10rpx;
}

.promo-subtitle {
  font-size: 24rpx;
  color: #666666;
  display: block;
}

.promo-image {
  font-size: 80rpx;
  margin: 0 20rpx;
}

.promo-btn {
  padding: 12rpx 24rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 30rpx;
  font-size: 24rpx;
}

/* 会员提示 */
.member-tip {
  margin: 0 30rpx 20rpx;
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 20rpx 30rpx;
  display: flex;
  align-items: center;
  gap: 15rpx;
}

.tip-icon {
  width: 40rpx;
  height: 40rpx;
  background-color: #52C41A;
  color: #FFFFFF;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
}

.tip-text {
  font-size: 28rpx;
  color: #333333;
}

.tip-desc {
  flex: 1;
  font-size: 24rpx;
  color: #999999;
}

.tip-btn {
  padding: 8rpx 20rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 20rpx;
  font-size: 22rpx;
}

/* 主要服务 */
.main-services {
  display: flex;
  gap: 20rpx;
  padding: 0 30rpx;
  margin-bottom: 20rpx;
}

.service-card {
  flex: 1;
  height: 200rpx;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: #FFFFFF;
}

.service-card.express {
  background: linear-gradient(135deg, #FF3B30 0%, #D81E06 100%);
}

.service-card.logistics {
  background: linear-gradient(135deg, #555555 0%, #333333 100%);
}

.service-icon {
  font-size: 80rpx;
  margin-bottom: 15rpx;
}

.service-info {
  text-align: center;
}

.service-title {
  font-size: 32rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 8rpx;
}

.service-subtitle {
  font-size: 24rpx;
  opacity: 0.9;
  display: block;
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
  width: 25%;
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
}

.grid-icon {
  width: 100rpx;
  height: 100rpx;
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48rpx;
  margin-bottom: 10rpx;
  position: relative;
}

.grid-badge {
  position: absolute;
  top: -10rpx;
  right: -10rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  font-size: 18rpx;
  padding: 4rpx 8rpx;
  border-radius: 10rpx;
  white-space: nowrap;
}

.grid-label {
  font-size: 24rpx;
  color: #666666;
  text-align: center;
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
