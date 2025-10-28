<template>
  <view v-if="show" class="sf-loading">
    <view class="loading-mask" :style="{ backgroundColor: maskColor }"></view>
    <view class="loading-content">
      <view class="loading-spinner">
        <view class="spinner-item" v-for="i in 12" :key="i"></view>
      </view>
      <text v-if="text" class="loading-text">{{ text }}</text>
    </view>
  </view>
</template>

<script setup>
defineProps({
  show: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: ''
  },
  maskColor: {
    type: String,
    default: 'rgba(0, 0, 0, 0.5)'
  }
})
</script>

<style scoped>
.sf-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 9999;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.loading-content {
  position: relative;
  background-color: rgba(0, 0, 0, 0.7);
  border-radius: 20rpx;
  padding: 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20rpx;
}

.loading-spinner {
  width: 80rpx;
  height: 80rpx;
  position: relative;
  animation: rotate 1.2s linear infinite;
}

.spinner-item {
  position: absolute;
  width: 6rpx;
  height: 20rpx;
  background-color: #FFFFFF;
  border-radius: 3rpx;
  left: 50%;
  top: 0;
  margin-left: -3rpx;
  transform-origin: center 40rpx;
  opacity: 0.3;
}

.spinner-item:nth-child(1) { transform: rotate(0deg); animation-delay: 0s; }
.spinner-item:nth-child(2) { transform: rotate(30deg); animation-delay: 0.1s; }
.spinner-item:nth-child(3) { transform: rotate(60deg); animation-delay: 0.2s; }
.spinner-item:nth-child(4) { transform: rotate(90deg); animation-delay: 0.3s; }
.spinner-item:nth-child(5) { transform: rotate(120deg); animation-delay: 0.4s; }
.spinner-item:nth-child(6) { transform: rotate(150deg); animation-delay: 0.5s; }
.spinner-item:nth-child(7) { transform: rotate(180deg); animation-delay: 0.6s; }
.spinner-item:nth-child(8) { transform: rotate(210deg); animation-delay: 0.7s; }
.spinner-item:nth-child(9) { transform: rotate(240deg); animation-delay: 0.8s; }
.spinner-item:nth-child(10) { transform: rotate(270deg); animation-delay: 0.9s; }
.spinner-item:nth-child(11) { transform: rotate(300deg); animation-delay: 1s; }
.spinner-item:nth-child(12) { transform: rotate(330deg); animation-delay: 1.1s; }

@keyframes rotate {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.loading-text {
  color: #FFFFFF;
  font-size: 28rpx;
}
</style>
