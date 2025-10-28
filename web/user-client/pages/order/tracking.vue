<template>
  <view class="container">
    <!-- 头部信息 -->
    <view class="header">
      <view class="status-section">
        <text class="status-icon">{{ statusIcon }}</text>
        <view class="status-info">
          <text class="status-text">{{ statusText }}</text>
          <text class="status-desc">{{ statusDesc }}</text>
        </view>
      </view>
      
      <view class="tracking-no">
        <text>运单号：{{ trackingNo }}</text>
        <text class="copy-btn" @click="copyTrackingNo">复制</text>
      </view>
      
      <view v-if="estimatedTime" class="estimated-time">
        <text>预计送达：{{ estimatedTime }}</text>
      </view>
    </view>

    <!-- 地图轨迹 -->
    <view class="map-section">
      <map
        class="map"
        :latitude="mapCenter.latitude"
        :longitude="mapCenter.longitude"
        :markers="markers"
        :polyline="polyline"
        :show-location="true"
      ></map>
    </view>

    <!-- 物流时间轴 -->
    <view class="timeline-section">
      <view class="section-title">
        <text>物流信息</text>
      </view>
      
      <view class="timeline">
        <view 
          v-for="(log, index) in trackingLogs" 
          :key="index"
          class="timeline-item"
          :class="{ active: index === 0 }"
        >
          <view class="timeline-dot"></view>
          <view class="timeline-line" v-if="index < trackingLogs.length - 1"></view>
          
          <view class="timeline-content">
            <view class="timeline-time">{{ log.time }}</view>
            <view class="timeline-location">{{ log.location }}</view>
            <view class="timeline-desc">{{ log.description }}</view>
          </view>
        </view>
      </view>
    </view>

    <!-- 操作按钮 -->
    <view class="action-buttons">
      <view class="action-btn" @click="contactCourier">
        <text class="btn-icon">📞</text>
        <text class="btn-text">联系快递员</text>
      </view>
      <view class="action-btn" @click="shareTracking">
        <text class="btn-icon">📤</text>
        <text class="btn-text">分享物流</text>
      </view>
      <view class="action-btn" @click="complaint">
        <text class="btn-icon">💬</text>
        <text class="btn-text">投诉建议</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import { onLoad } from '@dcloudio/uni-app'

const orderStore = useOrderStore()

// 运单号
const trackingNo = ref('')

// 物流信息
const trackingInfo = ref(null)
const trackingLogs = ref([])

// 地图相关
const mapCenter = ref({
  latitude: 22.547,
  longitude: 114.085
})

const markers = ref([
  {
    id: 1,
    latitude: 22.547,
    longitude: 114.085,
    iconPath: '/static/images/marker-current.png',
    width: 30,
    height: 30
  }
])

const polyline = ref([
  {
    points: [
      { latitude: 22.547, longitude: 114.085 },
      { latitude: 22.557, longitude: 114.095 }
    ],
    color: '#D81E06',
    width: 4
  }
])

// 计算属性
const statusIcon = computed(() => {
  if (!trackingInfo.value) return '📦'
  
  const icons = {
    pending: '⏰',
    accepted: '✓',
    picked_up: '📦',
    in_transit: '🚚',
    out_for_delivery: '🚴',
    delivered: '✅',
    cancelled: '❌',
    exception: '⚠️'
  }
  return icons[trackingInfo.value.status] || '📦'
})

const statusText = computed(() => {
  if (!trackingInfo.value) return '查询中'
  
  const texts = {
    pending: '待接单',
    accepted: '已接单',
    picked_up: '已揽收',
    in_transit: '运输中',
    out_for_delivery: '派送中',
    delivered: '已签收',
    cancelled: '已取消',
    exception: '异常'
  }
  return texts[trackingInfo.value.status] || '未知状态'
})

const statusDesc = computed(() => {
  if (!trackingInfo.value || trackingLogs.value.length === 0) return ''
  return trackingLogs.value[0].description
})

const estimatedTime = computed(() => {
  return trackingInfo.value?.estimatedDeliveryTime || ''
})

// 方法
function loadTrackingInfo() {
  if (!trackingNo.value) return
  
  orderStore.getTrackingInfo(trackingNo.value).then(data => {
    trackingInfo.value = data
    trackingLogs.value = data.logs || []
  }).catch(err => {
    uni.showToast({
      title: '查询失败',
      icon: 'none'
    })
  })
}

function copyTrackingNo() {
  uni.setClipboardData({
    data: trackingNo.value,
    success: () => {
      uni.showToast({
        title: '已复制',
        icon: 'success'
      })
    }
  })
}

function contactCourier() {
  uni.showModal({
    title: '联系快递员',
    content: '是否拨打快递员电话？',
    success: (res) => {
      if (res.confirm) {
        uni.makePhoneCall({
          phoneNumber: '13800138000'
        })
      }
    }
  })
}

function shareTracking() {
  uni.showShareMenu({
    withShareTicket: true,
    success: () => {
      console.log('分享成功')
    }
  })
}

function complaint() {
  uni.showToast({
    title: '投诉建议功能开发中',
    icon: 'none'
  })
}

onLoad((options) => {
  if (options.trackingNo) {
    trackingNo.value = options.trackingNo
    loadTrackingInfo()
  } else if (options.orderId) {
    // 通过订单ID加载物流信息
    loadTrackingInfoByOrderId(options.orderId)
  }
})

onMounted(() => {
  if (!trackingNo.value) {
    uni.showToast({
      title: '请输入运单号',
      icon: 'none'
    })
  }
})

// 通过订单ID加载物流信息
async function loadTrackingInfoByOrderId(orderId) {
  try {
    uni.showLoading({
      title: '加载中...',
      mask: true
    })
    
    const data = await orderStore.getTrackingInfo(orderId)
    
    uni.hideLoading()
    
    if (data) {
      trackingNo.value = data.orderNo || ''
      trackingInfo.value = data
      trackingLogs.value = data.logs || []
    } else {
      uni.showToast({
        title: '暂无物流信息',
        icon: 'none'
      })
    }
  } catch (error) {
    uni.hideLoading()
    console.error('加载物流信息失败:', error)
    uni.showToast({
      title: '加载失败',
      icon: 'none'
    })
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F8F8F8;
}

/* 头部信息 */
.header {
  background: linear-gradient(135deg, #FF3B30 0%, #D81E06 100%);
  padding: 40rpx 30rpx;
  color: #FFFFFF;
}

.status-section {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 30rpx;
}

.status-icon {
  font-size: 80rpx;
}

.status-info {
  flex: 1;
}

.status-text {
  font-size: 36rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 10rpx;
}

.status-desc {
  font-size: 26rpx;
  opacity: 0.9;
  display: block;
}

.tracking-no {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: 28rpx;
  margin-bottom: 15rpx;
}

.copy-btn {
  padding: 8rpx 20rpx;
  background-color: rgba(255, 255, 255, 0.2);
  border-radius: 20rpx;
  font-size: 24rpx;
}

.estimated-time {
  font-size: 26rpx;
  opacity: 0.9;
}

/* 地图 */
.map-section {
  height: 400rpx;
  margin: 20rpx 30rpx;
  border-radius: 20rpx;
  overflow: hidden;
}

.map {
  width: 100%;
  height: 100%;
}

/* 时间轴 */
.timeline-section {
  background-color: #FFFFFF;
  margin: 20rpx 30rpx;
  border-radius: 20rpx;
  padding: 30rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 30rpx;
}

.timeline {
  position: relative;
}

.timeline-item {
  position: relative;
  padding-left: 60rpx;
  padding-bottom: 40rpx;
}

.timeline-item:last-child {
  padding-bottom: 0;
}

.timeline-dot {
  position: absolute;
  left: 0;
  top: 0;
  width: 24rpx;
  height: 24rpx;
  background-color: #E5E5E5;
  border-radius: 50%;
  border: 4rpx solid #FFFFFF;
  box-shadow: 0 0 0 2rpx #E5E5E5;
}

.timeline-item.active .timeline-dot {
  background-color: #D81E06;
  box-shadow: 0 0 0 2rpx #D81E06;
}

.timeline-line {
  position: absolute;
  left: 12rpx;
  top: 24rpx;
  bottom: -40rpx;
  width: 2rpx;
  background-color: #E5E5E5;
}

.timeline-content {
  padding-top: 0;
}

.timeline-time {
  font-size: 24rpx;
  color: #999999;
  margin-bottom: 10rpx;
}

.timeline-item.active .timeline-time {
  color: #D81E06;
}

.timeline-location {
  font-size: 28rpx;
  color: #666666;
  margin-bottom: 5rpx;
}

.timeline-item.active .timeline-location {
  color: #333333;
  font-weight: 600;
}

.timeline-desc {
  font-size: 26rpx;
  color: #999999;
  line-height: 1.6;
}

.timeline-item.active .timeline-desc {
  color: #666666;
}

/* 操作按钮 */
.action-buttons {
  display: flex;
  gap: 20rpx;
  padding: 30rpx;
}

.action-btn {
  flex: 1;
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
}

.btn-icon {
  font-size: 48rpx;
}

.btn-text {
  font-size: 24rpx;
  color: #666666;
}
</style>
