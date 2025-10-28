<template>
  <view class="sf-empty">
    <image v-if="image" class="empty-image" :src="image" mode="aspectFit"></image>
    <view v-else class="empty-icon">📦</view>
    <text class="empty-text">{{ text }}</text>
    <view v-if="showButton" class="empty-button" @click="handleClick">
      {{ buttonText }}
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  image: {
    type: String,
    default: ''
  },
  text: {
    type: String,
    default: '暂无数据'
  },
  showButton: {
    type: Boolean,
    default: false
  },
  buttonText: {
    type: String,
    default: '去操作'
  }
})

const emit = defineEmits(['click'])

const handleClick = () => {
  emit('click')
}
</script>

<style scoped>
.sf-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 60rpx;
}

.empty-image {
  width: 300rpx;
  height: 300rpx;
  margin-bottom: 40rpx;
}

.empty-icon {
  font-size: 200rpx;
  margin-bottom: 40rpx;
  opacity: 0.3;
}

.empty-text {
  font-size: 28rpx;
  color: #999999;
  margin-bottom: 40rpx;
}

.empty-button {
  padding: 20rpx 60rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 50rpx;
  font-size: 28rpx;
}
</style>
