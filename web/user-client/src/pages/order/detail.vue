<template>
  <view class="container">
    <!-- 订单状态 -->
    <view class="status-card">
      <view class="status-icon">{{ statusIcon }}</view>
      <view class="status-info">
        <text class="status-text">{{ statusText }}</text>
        <text class="status-desc">{{ statusDesc }}</text>
      </view>
    </view>

    <!-- 物流进度 -->
    <view class="progress-card" @click="goToTracking">
      <view class="progress-header">
        <text class="progress-title">物流进度</text>
        <text class="progress-link">查看详情 ></text>
      </view>
      <view v-if="latestLog" class="progress-content">
        <text class="progress-time">{{ latestLog.time }}</text>
        <text class="progress-desc">{{ latestLog.description }}</text>
      </view>
    </view>

    <!-- 收寄件信息 -->
    <view class="info-card">
      <view class="info-section">
        <view class="info-header">
          <text class="info-label">寄件信息</text>
        </view>
        <view class="info-row">
          <text class="info-name">{{ order.senderName }}</text>
          <text class="info-phone">{{ order.senderPhone }}</text>
        </view>
        <view class="info-address">{{ order.senderAddress }}</view>
      </view>

      <view class="info-divider"></view>

      <view class="info-section">
        <view class="info-header">
          <text class="info-label">收件信息</text>
        </view>
        <view class="info-row">
          <text class="info-name">{{ order.receiverName }}</text>
          <text class="info-phone">{{ order.receiverPhone }}</text>
        </view>
        <view class="info-address">{{ order.receiverAddress }}</view>
      </view>
    </view>

    <!-- 订单信息 -->
    <view class="detail-card">
      <view class="detail-title">订单信息</view>
      <view class="detail-row">
        <text class="detail-label">运单号</text>
        <view class="detail-value">
          <text>{{ order.orderNo }}</text>
          <text class="copy-btn" @click="copyOrderNo">复制</text>
        </view>
      </view>
      <view class="detail-row">
        <text class="detail-label">取件码</text>
        <text class="detail-value">{{ order.pickupCode }}</text>
      </view>
      <view class="detail-row">
        <text class="detail-label">服务类型</text>
        <text class="detail-value">{{ order.serviceType }}</text>
      </view>
      <view class="detail-row">
        <text class="detail-label">下单时间</text>
        <text class="detail-value">{{ order.createdAt }}</text>
      </view>
      <view v-if="order.actualPickupTime" class="detail-row">
        <text class="detail-label">取件时间</text>
        <text class="detail-value">{{ order.actualPickupTime }}</text>
      </view>
    </view>

    <!-- 费用信息 -->
    <view class="fee-card">
      <view class="fee-title">费用信息</view>
      <view class="fee-row">
        <text class="fee-label">运费</text>
        <text class="fee-value">¥{{ order.estimatedFee }}</text>
      </view>
      <view v-if="order.isInsured" class="fee-row">
        <text class="fee-label">保价费</text>
        <text class="fee-value">¥{{ calculateInsuranceFee(order.insuredValue) }}</text>
      </view>
      <view class="fee-divider"></view>
      <view class="fee-row total">
        <text class="fee-label">实付金额</text>
        <text class="fee-value">¥{{ order.actualFee || order.estimatedFee }}</text>
      </view>
    </view>

    <!-- 操作按钮 -->
    <view class="action-bar">
      <view v-if="order.status === 'pending'" class="action-btn cancel" @click="cancelOrder">
        取消订单
      </view>
      <view class="action-btn primary" @click="contactService">
        联系客服
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import { onLoad } from '@dcloudio/uni-app'

const orderStore = useOrderStore()

// 订单ID
const orderId = ref(null)

// 订单信息
const order = ref({
  id: 1,
  orderNo: 'SF1234567890',
  pickupCode: '123456',
  status: 'in_transit',
  senderName: '张三',
  senderPhone: '138****0000',
  senderAddress: '广东省深圳市南山区科技园南区1001号',
  receiverName: '李四',
  receiverPhone: '139****0000',
  receiverAddress: '北京市朝阳区建国门外大街1号',
  serviceType: '顺丰标快',
  estimatedFee: 23.00,
  actualFee: 23.00,
  isInsured: false,
  insuredValue: 0,
  createdAt: '2024-01-19 10:30:00',
  actualPickupTime: '2024-01-19 14:30:00'
})

const latestLog = ref({
  time: '2024-01-19 16:00',
  description: '快件已到达深圳转运中心'
})

// 计算属性
const statusIcon = computed(() => {
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
  return icons[order.value.status] || '📦'
})

const statusText = computed(() => {
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
  return texts[order.value.status] || '未知状态'
})

const statusDesc = computed(() => {
  const descs = {
    pending: '快递员正在赶来的路上',
    accepted: '快递员已接单，即将上门取件',
    picked_up: '快件已揽收，正在运输中',
    in_transit: '快件正在运输途中',
    out_for_delivery: '快递员正在派送中',
    delivered: '快件已签收，感谢使用顺丰速运',
    cancelled: '订单已取消',
    exception: '快件出现异常，请联系客服'
  }
  return descs[order.value.status] || ''
})

// 方法
function loadOrderDetail() {
  if (!orderId.value) return
  
  orderStore.getOrderDetail(orderId.value).then(data => {
    if (data) {
      order.value = data
    }
  })
}

function goToTracking() {
  uni.navigateTo({
    url: `/pages/order/tracking?trackingNo=${order.value.orderNo}`
  })
}

function copyOrderNo() {
  uni.setClipboardData({
    data: order.value.orderNo,
    success: () => {
      uni.showToast({
        title: '已复制',
        icon: 'success'
      })
    }
  })
}

function calculateInsuranceFee(value) {
  return Math.max(1, value * 0.005).toFixed(2)
}

function cancelOrder() {
  uni.showModal({
    title: '提示',
    content: '确定要取消订单吗？',
    success: (res) => {
      if (res.confirm) {
        orderStore.cancelOrder(orderId.value).then(() => {
          uni.showToast({
            title: '订单已取消',
            icon: 'success'
          })
          setTimeout(() => {
            uni.navigateBack()
          }, 1500)
        })
      }
    }
  })
}

function contactService() {
  uni.showToast({
    title: '客服功能开发中',
    icon: 'none'
  })
}

onLoad((options) => {
  if (options.id) {
    orderId.value = parseInt(options.id)
    loadOrderDetail()
  }
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F8F8F8;
  padding-bottom: 120rpx;
}

/* 状态卡片 */
.status-card {
  background: linear-gradient(135deg, #FF3B30 0%, #D81E06 100%);
  padding: 40rpx 30rpx;
  display: flex;
  align-items: center;
  gap: 20rpx;
  color: #FFFFFF;
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

/* 物流进度 */
.progress-card {
  background-color: #FFFFFF;
  margin: 20rpx 30rpx;
  border-radius: 20rpx;
  padding: 30rpx;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.progress-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
}

.progress-link {
  font-size: 26rpx;
  color: #D81E06;
}

.progress-content {
  display: flex;
  flex-direction: column;
  gap: 10rpx;
}

.progress-time {
  font-size: 24rpx;
  color: #999999;
}

.progress-desc {
  font-size: 28rpx;
  color: #666666;
}

/* 信息卡片 */
.info-card {
  background-color: #FFFFFF;
  margin: 0 30rpx 20rpx;
  border-radius: 20rpx;
  padding: 30rpx;
}

.info-section {
  padding: 20rpx 0;
}

.info-header {
  margin-bottom: 20rpx;
}

.info-label {
  font-size: 28rpx;
  color: #999999;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15rpx;
}

.info-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.info-phone {
  font-size: 28rpx;
  color: #666666;
}

.info-address {
  font-size: 28rpx;
  color: #666666;
  line-height: 1.6;
}

.info-divider {
  height: 1rpx;
  background-color: #F5F5F5;
  margin: 20rpx 0;
}

/* 详情卡片 */
.detail-card {
  background-color: #FFFFFF;
  margin: 0 30rpx 20rpx;
  border-radius: 20rpx;
  padding: 30rpx;
}

.detail-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 30rpx;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25rpx;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  font-size: 28rpx;
  color: #999999;
}

.detail-value {
  font-size: 28rpx;
  color: #333333;
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.copy-btn {
  padding: 8rpx 20rpx;
  background-color: #F5F5F5;
  border-radius: 20rpx;
  font-size: 24rpx;
  color: #D81E06;
}

/* 费用卡片 */
.fee-card {
  background-color: #FFFFFF;
  margin: 0 30rpx 20rpx;
  border-radius: 20rpx;
  padding: 30rpx;
}

.fee-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 30rpx;
}

.fee-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 25rpx;
}

.fee-row:last-child {
  margin-bottom: 0;
}

.fee-label {
  font-size: 28rpx;
  color: #666666;
}

.fee-value {
  font-size: 28rpx;
  color: #333333;
}

.fee-divider {
  height: 1rpx;
  background-color: #F5F5F5;
  margin: 20rpx 0;
}

.fee-row.total .fee-label {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
}

.fee-row.total .fee-value {
  font-size: 36rpx;
  font-weight: 600;
  color: #D81E06;
}

/* 操作栏 */
.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #FFFFFF;
  padding: 20rpx 30rpx;
  display: flex;
  gap: 20rpx;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
  flex: 1;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 40rpx;
  font-size: 28rpx;
  font-weight: 600;
}

.action-btn.cancel {
  background-color: #F5F5F5;
  color: #666666;
}

.action-btn.primary {
  background-color: #D81E06;
  color: #FFFFFF;
}
</style>
