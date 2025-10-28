<template>
  <view class="container">
    <scroll-view class="content" scroll-y>
      <!-- 寄件地址 -->
      <view class="section">
        <view class="section-title">寄件信息</view>
        <view class="address-card" @click="selectAddress('sender')">
          <view v-if="senderAddress" class="address-info">
            <view class="address-header">
              <text class="address-name">{{ senderAddress.name }}</text>
              <text class="address-phone">{{ senderAddress.phone }}</text>
            </view>
            <text class="address-detail">
              {{ senderAddress.province }} {{ senderAddress.city }} {{ senderAddress.district }} {{ senderAddress.detail }}
            </text>
          </view>
          <view v-else class="address-placeholder">
            <text class="placeholder-icon">+</text>
            <text class="placeholder-text">添加寄件地址</text>
          </view>
        </view>
      </view>

      <!-- 收件地址 -->
      <view class="section">
        <view class="section-title">收件信息</view>
        <view class="address-card" @click="selectAddress('receiver')">
          <view v-if="receiverAddress" class="address-info">
            <view class="address-header">
              <text class="address-name">{{ receiverAddress.name }}</text>
              <text class="address-phone">{{ receiverAddress.phone }}</text>
            </view>
            <text class="address-detail">
              {{ receiverAddress.province }} {{ receiverAddress.city }} {{ receiverAddress.district }} {{ receiverAddress.detail }}
            </text>
          </view>
          <view v-else class="address-placeholder">
            <text class="placeholder-icon">+</text>
            <text class="placeholder-text">添加收件地址</text>
          </view>
        </view>
      </view>

      <!-- 服务选择 -->
      <view class="section">
        <view class="section-title">选择服务</view>
        <view class="service-list">
          <view 
            v-for="service in services" 
            :key="service.type"
            class="service-item"
            :class="{ active: selectedService === service.type }"
            @click="selectService(service)"
          >
            <view class="service-info">
              <text class="service-name">{{ service.name }}</text>
              <text class="service-time">{{ service.timeLimit }}</text>
            </view>
            <text class="service-price">¥{{ calculateFee(service) }}</text>
          </view>
        </view>
      </view>

      <!-- 物品信息 -->
      <view class="section">
        <view class="section-title">物品信息</view>
        <view class="form-item">
          <text class="form-label">物品类型</text>
          <input v-model="productType" class="form-input" placeholder="如：文件、衣物等" />
        </view>
        <view class="form-item">
          <text class="form-label">重量（kg）</text>
          <input v-model.number="weight" type="digit" class="form-input" placeholder="请输入重量" />
        </view>
      </view>

      <!-- 增值服务 -->
      <view class="section">
        <view class="section-title">增值服务</view>
        <view class="checkbox-item" @click="isInsured = !isInsured">
          <checkbox :checked="isInsured" />
          <text class="checkbox-label">保价服务</text>
        </view>
        <view v-if="isInsured" class="form-item">
          <text class="form-label">声明价值（元）</text>
          <input v-model.number="insuredValue" type="digit" class="form-input" placeholder="最高20000元" />
          <text class="form-tip">保价费：¥{{ insuranceFee }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="footer">
      <view class="fee-info">
        <text class="fee-label">预估费用</text>
        <text class="fee-value">¥{{ totalFee }}</text>
      </view>
      <view class="submit-btn" @click="submitOrder">
        <text>提交订单</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useOrderStore } from '@/stores/order'

const orderStore = useOrderStore()

// 实时计价数据
const pricingData = ref(null)
const isCalculating = ref(false)

// 地址信息
const senderAddress = ref(null)
const receiverAddress = ref(null)

// 服务类型
const services = ref([
  { type: 'same_day', name: '顺丰即日', timeLimit: '当日达', priceRate: 2.0, firstWeight: 12 },
  { type: 'next_day', name: '顺丰特快', timeLimit: '次日达', priceRate: 1.5, firstWeight: 12 },
  { type: 'standard', name: '顺丰标快', timeLimit: '隔日达', priceRate: 1.0, firstWeight: 12 },
  { type: 'economy', name: '顺丰快递包裹', timeLimit: '3-5天', priceRate: 0.7, firstWeight: 12 }
])

const selectedService = ref('standard')

// 物品信息
const productType = ref('')
const weight = ref(1)

// 保价信息
const isInsured = ref(false)
const insuredValue = ref(0)

// 计算费用
const insuranceFee = computed(() => {
  if (!isInsured.value || !insuredValue.value) return 0
  return Math.max(1, insuredValue.value * 0.005).toFixed(2)
})

const totalFee = computed(() => {
  const service = services.value.find(s => s.type === selectedService.value)
  if (!service) return '0.00'
  
  const baseFee = parseFloat(calculateFee(service))
  const insurance = parseFloat(insuranceFee.value)
  
  return (baseFee + insurance).toFixed(2)
})

// 方法
function calculateFee(service) {
  if (!weight.value) return service.firstWeight.toFixed(2)
  
  const firstWeight = service.firstWeight
  const additionalWeight = Math.ceil(weight.value - 1)
  const additionalFee = additionalWeight * 2 * service.priceRate
  
  return (firstWeight * service.priceRate + additionalFee).toFixed(2)
}

function selectService(service) {
  selectedService.value = service.type
}

function selectAddress(type) {
  // 保存当前选择的地址类型到全局
  getApp().globalData = getApp().globalData || {}
  getApp().globalData.addressType = type
  
  uni.navigateTo({
    url: `/pages/send/address-input?type=${type}`
  })
}

// 页面显示时检查是否有返回的地址数据
onShow(() => {
  const app = getApp()
  if (app.globalData && app.globalData.selectedAddress) {
    const address = app.globalData.selectedAddress
    const type = app.globalData.addressType
    
    if (type === 'sender') {
      senderAddress.value = address
    } else {
      receiverAddress.value = address
    }
    
    // 清除全局数据
    delete app.globalData.selectedAddress
    delete app.globalData.addressType
  }
})

function submitOrder() {
  // 验证
  if (!senderAddress.value) {
    uni.showToast({ title: '请选择寄件地址', icon: 'none' })
    return
  }
  if (!receiverAddress.value) {
    uni.showToast({ title: '请选择收件地址', icon: 'none' })
    return
  }
  if (!weight.value) {
    uni.showToast({ title: '请输入重量', icon: 'none' })
    return
  }

  // 创建订单
  const orderData = {
    senderName: senderAddress.value.name,
    senderPhone: senderAddress.value.phone,
    senderAddress: `${senderAddress.value.province}${senderAddress.value.city}${senderAddress.value.district}${senderAddress.value.detail}`,
    receiverName: receiverAddress.value.name,
    receiverPhone: receiverAddress.value.phone,
    receiverAddress: `${receiverAddress.value.province}${receiverAddress.value.city}${receiverAddress.value.district}${receiverAddress.value.detail}`,
    serviceType: services.value.find(s => s.type === selectedService.value).name,
    productType: productType.value,
    estimatedWeight: weight.value,
    estimatedFee: totalFee.value,
    isInsured: isInsured.value,
    insuredValue: insuredValue.value
  }

  orderStore.createOrder(orderData).then((order) => {
    uni.showToast({
      title: '下单成功',
      icon: 'success'
    })
    
    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/order/detail?id=${order.id}`
      })
    }, 1500)
  })
}
</script>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #F8F8F8;
}

.content {
  flex: 1;
  padding-bottom: 20rpx;
}

.section {
  margin: 20rpx 30rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 20rpx;
}

/* 地址卡片 */
.address-card {
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
}

.address-info {
  display: flex;
  flex-direction: column;
  gap: 15rpx;
}

.address-header {
  display: flex;
  gap: 20rpx;
}

.address-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.address-phone {
  font-size: 28rpx;
  color: #666666;
}

.address-detail {
  font-size: 28rpx;
  color: #666666;
  line-height: 1.6;
}

.address-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 15rpx;
  padding: 40rpx 0;
  color: #999999;
}

.placeholder-icon {
  font-size: 48rpx;
}

.placeholder-text {
  font-size: 28rpx;
}

/* 服务列表 */
.service-list {
  background-color: #FFFFFF;
  border-radius: 20rpx;
  overflow: hidden;
}

.service-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #F5F5F5;
}

.service-item:last-child {
  border-bottom: none;
}

.service-item.active {
  background-color: #FFF5F5;
}

.service-info {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.service-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
}

.service-time {
  font-size: 24rpx;
  color: #999999;
}

.service-price {
  font-size: 32rpx;
  font-weight: 600;
  color: #D81E06;
}

/* 表单 */
.form-item {
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.form-label {
  display: block;
  font-size: 28rpx;
  color: #333333;
  margin-bottom: 15rpx;
}

.form-input {
  width: 100%;
  padding: 20rpx;
  background-color: #F5F5F5;
  border-radius: 12rpx;
  font-size: 28rpx;
}

.form-tip {
  display: block;
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #D81E06;
}

.checkbox-item {
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  gap: 15rpx;
}

.checkbox-label {
  font-size: 28rpx;
  color: #333333;
}

/* 底部操作栏 */
.footer {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 20rpx 30rpx;
  background-color: #FFFFFF;
  box-shadow: 0 -4rpx 20rpx rgba(0, 0, 0, 0.05);
}

.fee-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 5rpx;
}

.fee-label {
  font-size: 24rpx;
  color: #999999;
}

.fee-value {
  font-size: 36rpx;
  font-weight: 600;
  color: #D81E06;
}

.submit-btn {
  flex: 1;
  height: 80rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  font-weight: 600;
}
</style>
