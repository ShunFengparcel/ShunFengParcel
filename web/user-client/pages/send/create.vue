<template>
  <view class="container">
    <!-- 实名认证提示 -->
    <view class="auth-banner">
      <text class="auth-text">根据相关法律法规要求，寄件须实名认证。推荐您提前线上实名，一次实名长期有效，寄件更便捷～</text>
      <text class="auth-link">立即实名</text>
    </view>

    <scroll-view class="content" scroll-y>
      <!-- 常规/模板切换 -->
      <view class="tab-bar">
        <view class="tab-item active">
          <text class="tab-text">常规</text>
        </view>
        <view class="tab-item">
          <text class="tab-text">+ 模板</text>
        </view>
      </view>

      <!-- 寄件人信息 -->
      <view class="address-card" @click="selectAddress('sender')">
        <view class="card-header">
          <view class="address-icon sender-icon">
            <text class="icon-text">寄</text>
          </view>
          <view class="card-content">
            <view v-if="senderAddress" class="address-info">
              <view class="name-phone">
                <text class="name">{{ senderAddress.name }}</text>
                <text class="phone">{{ formatPhone(senderAddress.phone) }}</text>
              </view>
              <text class="full-address">
                {{ senderAddress.province }}{{ senderAddress.city }}{{ senderAddress.district }}{{ senderAddress.detail }}
              </text>
            </view>
            <view v-else class="address-placeholder">
              <text class="placeholder-text">寄件人信息</text>
            </view>
          </view>
          <view class="card-action">
            <text class="action-text">地址簿</text>
          </view>
        </view>
      </view>

      <!-- 收件人信息 -->
      <view class="address-card receiver-card" @click="selectAddress('receiver')">
        <view class="card-header">
          <view class="address-icon receiver-icon">
            <text class="icon-text">收</text>
          </view>
          <view class="card-content">
            <view v-if="receiverAddress" class="address-info">
              <view class="name-phone">
                <text class="name">{{ receiverAddress.name }}</text>
                <text class="phone">{{ formatPhone(receiverAddress.phone) }}</text>
              </view>
              <text class="full-address">
                {{ receiverAddress.province }}{{ receiverAddress.city }}{{ receiverAddress.district }}{{ receiverAddress.detail }}
              </text>
            </view>
            <view v-else class="address-placeholder">
              <text class="placeholder-main">收件人信息</text>
              <text class="placeholder-sub">支持文本贴、图片识别、顺丰ID</text>
            </view>
          </view>
          <view class="card-action">
            <text class="action-text">地址簿</text>
          </view>
        </view>
      </view>

      <!-- 上门取件/服务点自寄 -->
      <view class="service-mode">
        <view 
          class="mode-item" 
          :class="{ active: pickupMode === 'pickup' }"
          @click="pickupMode = 'pickup'"
        >
          <view class="mode-icon">🚪</view>
          <text class="mode-text">上门取件</text>
        </view>
        <view 
          class="mode-item"
          :class="{ active: pickupMode === 'self' }"
          @click="pickupMode = 'self'"
        >
          <view class="mode-icon">📍</view>
          <text class="mode-text">服务点自寄</text>
        </view>
      </view>

      <!-- 期望上门时间 -->
      <view class="form-row" @click="showTimePicker = true">
        <text class="form-label">期望上门时间</text>
        <view class="form-value">
          <text class="time-badge">⏰ {{ pickupTime }}</text>
          <text class="form-text">{{ pickupDate || '今天 一小时内' }}</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 物品信息 -->
      <view class="form-row">
        <view class="form-label-wrapper">
          <text class="form-label">物品信息</text>
          <text class="required">必填</text>
        </view>
        <view class="form-value">
          <input 
            v-model="productType" 
            class="form-input-inline" 
            placeholder="如：文件、衣物等"
            placeholder-class="placeholder-style"
          />
        </view>
      </view>

      <!-- 重量 -->
      <view class="form-row">
        <view class="form-label-wrapper">
          <text class="form-label">重量</text>
          <text class="required">必填</text>
        </view>
        <view class="form-value">
          <input 
            v-model="weight" 
            type="digit"
            class="form-input-inline" 
            placeholder="请输入重量"
            placeholder-class="placeholder-style"
            @input="onWeightChange"
          />
          <text class="unit-text">kg</text>
        </view>
      </view>

      <!-- 服务类型 -->
      <view class="form-row" @click="showServicePicker = true">
        <text class="form-label">服务类型</text>
        <view class="form-value">
          <text class="form-text">{{ serviceType }}</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 付款方式 -->
      <view class="form-row" @click="showPaymentPicker = true">
        <text class="form-label">付款方式</text>
        <view class="form-value">
          <text class="payment-tip">可开通支付宝自动扣款</text>
          <text class="form-text">{{ paymentMethod }}</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 寄付现结 -->
      <view class="form-row">
        <text class="form-label">寄付现结</text>
        <view class="form-value">
          <text class="price-text">未保价物品最高赔7倍运费</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 保价 -->
      <view class="form-row" @click="showInsurancePicker = true">
        <text class="form-label">保价</text>
        <view class="form-value">
          <text class="form-text">{{ isInsured ? `¥${insuredValue}` : '未保价' }}</text>
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 增值服务 -->
      <view class="form-row">
        <text class="form-label">增值服务</text>
        <view class="form-value">
          <text class="arrow">›</text>
        </view>
      </view>

      <!-- 协议 -->
      <view class="agreement">
        <checkbox class="agreement-checkbox" :checked="agreedToTerms" @click="agreedToTerms = !agreedToTerms" />
        <text class="agreement-text">请阅读《电子运单契约条款》</text>
      </view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="footer">
      <view class="fee-section">
        <text class="fee-label">预估</text>
        <text class="fee-amount">¥ {{ totalFee > 0 ? totalFee.toFixed(2) : '--' }}</text>
        <text class="fee-detail" @click="showFeeDetail = !showFeeDetail">明细 {{ showFeeDetail ? '▲' : '▼' }}</text>
      </view>
      <view class="submit-btn" @click="submitOrder">
        <text class="submit-text">下单</text>
      </view>
    </view>

    <!-- 费用明细弹窗 -->
    <view v-if="showFeeDetail" class="fee-detail-modal" @click="showFeeDetail = false">
      <view class="fee-detail-content" @click.stop>
        <view class="fee-detail-header">
          <text class="fee-detail-title">费用明细</text>
          <text class="fee-detail-close" @click="showFeeDetail = false">✕</text>
        </view>
        <view class="fee-detail-body">
          <view class="fee-item">
            <text class="fee-item-label">基础运费（首重1kg）</text>
            <text class="fee-item-value">¥{{ pricingDetail.baseFee || 0 }}</text>
          </view>
          <view class="fee-item" v-if="pricingDetail.weightFee > 0">
            <text class="fee-item-label">续重费用</text>
            <text class="fee-item-value">¥{{ pricingDetail.weightFee || 0 }}</text>
          </view>
          <view class="fee-item" v-if="pricingDetail.insuranceFee > 0">
            <text class="fee-item-label">保价费用</text>
            <text class="fee-item-value">¥{{ pricingDetail.insuranceFee || 0 }}</text>
          </view>
          <view class="fee-divider"></view>
          <view class="fee-item fee-total">
            <text class="fee-item-label">合计</text>
            <text class="fee-item-value total">¥{{ totalFee.toFixed(2) }}</text>
          </view>
          <view class="fee-tip">
            <text class="fee-tip-text">预计送达：{{ estimatedTime || '--' }}</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 服务类型选择器 -->
    <picker 
      v-if="showServicePicker"
      mode="selector" 
      :range="serviceTypes" 
      :value="serviceTypeIndex"
      @change="onServiceTypeChange"
      @cancel="showServicePicker = false"
    >
      <view></view>
    </picker>

    <!-- 保价输入弹窗 -->
    <view v-if="showInsurancePicker" class="insurance-modal" @click="showInsurancePicker = false">
      <view class="insurance-content" @click.stop>
        <view class="insurance-header">
          <text class="insurance-title">设置保价</text>
          <text class="insurance-close" @click="showInsurancePicker = false">✕</text>
        </view>
        <view class="insurance-body">
          <view class="insurance-switch">
            <text class="insurance-label">是否保价</text>
            <switch :checked="isInsured" @change="onInsuranceSwitch" color="#FF4D4F" />
          </view>
          <view v-if="isInsured" class="insurance-input-wrapper">
            <text class="insurance-input-label">声明价值</text>
            <input 
              v-model="insuredValue" 
              type="digit"
              class="insurance-input" 
              placeholder="请输入物品价值"
              @input="onInsuredValueChange"
            />
            <text class="insurance-unit">元</text>
          </view>
          <view v-if="isInsured" class="insurance-tip">
            <text class="insurance-tip-text">保价费率：0.5%，最低1元</text>
          </view>
        </view>
        <view class="insurance-footer">
          <view class="insurance-btn" @click="confirmInsurance">
            <text class="insurance-btn-text">确定</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { useOrderStore } from '@/stores/order'

const orderStore = useOrderStore()

// 页面显示时检查是否有返回的地址数据
onShow(() => {
  try {
    const app = getApp()
    if (app && app.globalData && app.globalData.selectedAddress) {
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

      // 触发计价
      calculatePriceDebounced()
    }
  } catch (error) {
    console.error('获取全局数据失败:', error)
  }
})

// 地址信息
const senderAddress = ref(null)
const receiverAddress = ref(null)

// 取件方式
const pickupMode = ref('pickup') // pickup: 上门取件, self: 服务点自寄

// 期望上门时间
const pickupTime = ref('今天 一小时内')
const pickupDate = ref('')

// 物品信息
const productType = ref('')

// 重量
const weight = ref('')

// 服务类型
const serviceTypes = ['顺丰标快', '顺丰特快', '顺丰即日', '顺丰快递包裹']
const serviceType = ref('顺丰标快')
const serviceTypeIndex = ref(0)
const showServicePicker = ref(false)

// 付款方式
const paymentMethod = ref('寄付现结')

// 保价信息
const isInsured = ref(false)
const insuredValue = ref(0)
const showInsurancePicker = ref(false)

// 费用信息
const totalFee = ref(0)
const estimatedTime = ref('')
const pricingDetail = ref({
  baseFee: 0,
  weightFee: 0,
  insuranceFee: 0
})
const showFeeDetail = ref(false)

// 协议
const agreedToTerms = ref(false)

// 计价防抖定时器
let pricingTimer = null

// 方法
function formatPhone(phone) {
  if (!phone) return ''
  // 格式化手机号为 182****4509
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

function selectAddress(type) {
  try {
    // 保存当前选择的地址类型到全局
    const app = getApp()
    if (app) {
      app.globalData = app.globalData || {}
      app.globalData.addressType = type
    }
  } catch (error) {
    console.error('保存全局数据失败:', error)
  }
  
  uni.navigateTo({
    url: '/pages/send/address'
  })
}

// 重量变化时触发计价
function onWeightChange() {
  calculatePriceDebounced()
}

// 服务类型变化
function onServiceTypeChange(e) {
  serviceTypeIndex.value = e.detail.value
  serviceType.value = serviceTypes[e.detail.value]
  showServicePicker.value = false
  calculatePriceDebounced()
}

// 保价开关
function onInsuranceSwitch(e) {
  isInsured.value = e.detail.value
  if (!isInsured.value) {
    insuredValue.value = 0
  }
}

// 保价金额变化
function onInsuredValueChange() {
  // 实时计价会在确认时触发
}

// 确认保价
function confirmInsurance() {
  showInsurancePicker.value = false
  calculatePriceDebounced()
}

// 防抖计价
function calculatePriceDebounced() {
  if (pricingTimer) {
    clearTimeout(pricingTimer)
  }
  
  pricingTimer = setTimeout(() => {
    calculatePrice()
  }, 500)
}

// 实时计价
async function calculatePrice() {
  // 检查必要条件
  if (!senderAddress.value || !receiverAddress.value) {
    return
  }
  
  const weightNum = parseFloat(weight.value)
  if (!weightNum || weightNum <= 0) {
    return
  }

  try {
    console.log('💰 开始计价...')
    
    const pricingData = {
      senderCity: senderAddress.value.city,
      receiverCity: receiverAddress.value.city,
      weight: weightNum,
      serviceType: serviceType.value,
      isInsured: isInsured.value,
      insuredValue: parseFloat(insuredValue.value) || 0
    }

    const result = await orderStore.calculatePrice(pricingData)
    
    if (result && result.data) {
      pricingDetail.value = {
        baseFee: result.data.base_fee || 0,
        weightFee: result.data.weight_fee || 0,
        insuranceFee: result.data.insurance_fee || 0
      }
      totalFee.value = result.data.total_fee || 0
      estimatedTime.value = result.data.estimated_time || ''
      
      console.log('💰 计价成功:', totalFee.value)
    }
  } catch (error) {
    console.error('💰 计价失败:', error)
    uni.showToast({
      title: '计价失败',
      icon: 'none'
    })
  }
}

async function submitOrder() {
  // 验证
  if (!senderAddress.value) {
    uni.showToast({ title: '请选择寄件地址', icon: 'none' })
    return
  }
  if (!receiverAddress.value) {
    uni.showToast({ title: '请选择收件地址', icon: 'none' })
    return
  }
  if (!productType.value) {
    uni.showToast({ title: '请填写物品信息', icon: 'none' })
    return
  }
  if (!weight.value || parseFloat(weight.value) <= 0) {
    uni.showToast({ title: '请输入重量', icon: 'none' })
    return
  }
  if (!agreedToTerms.value) {
    uni.showToast({ title: '请阅读并同意电子运单契约条款', icon: 'none' })
    return
  }

  // 显示加载提示
  uni.showLoading({
    title: '提交中...',
    mask: true
  })

  try {
    // 创建订单
    const orderData = {
      senderName: senderAddress.value.name,
      senderPhone: senderAddress.value.phone,
      senderAddress: `${senderAddress.value.province} ${senderAddress.value.city} ${senderAddress.value.district} ${senderAddress.value.detail}`,
      receiverName: receiverAddress.value.name,
      receiverPhone: receiverAddress.value.phone,
      receiverAddress: `${receiverAddress.value.province} ${receiverAddress.value.city} ${receiverAddress.value.district} ${receiverAddress.value.detail}`,
      pickupMode: pickupMode.value,
      pickupTime: pickupTime.value,
      serviceType: serviceType.value,
      productType: productType.value,
      paymentMethod: paymentMethod.value,
      isInsured: isInsured.value,
      insuredValue: parseFloat(insuredValue.value) || 0,
      estimatedFee: totalFee.value,
      estimatedWeight: parseFloat(weight.value)
    }

    const order = await orderStore.createOrder(orderData)
    
    uni.hideLoading()
    
    uni.showToast({
      title: '下单成功',
      icon: 'success',
      duration: 1500
    })
    
    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/order/detail?id=${order.id}`
      })
    }, 1500)
  } catch (error) {
    uni.hideLoading()
    console.error('提交订单失败:', error)
    uni.showToast({
      title: error.message || '下单失败，请重试',
      icon: 'none',
      duration: 2000
    })
  }
}
</script>

<style scoped>
.container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: #F5F5F5;
}

/* 实名认证横幅 */
.auth-banner {
  background: linear-gradient(135deg, #FFE8D6 0%, #FFD9B8 100%);
  padding: 20rpx 30rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20rpx;
}

.auth-text {
  flex: 1;
  font-size: 24rpx;
  color: #8B4513;
  line-height: 1.5;
}

.auth-link {
  font-size: 26rpx;
  color: #D81E06;
  font-weight: 600;
  white-space: nowrap;
}

.content {
  flex: 1;
}

/* 标签栏 */
.tab-bar {
  display: flex;
  background-color: #FFFFFF;
  padding: 20rpx 30rpx;
  gap: 30rpx;
}

.tab-item {
  padding: 10rpx 30rpx;
  border-radius: 30rpx;
}

.tab-item.active {
  background-color: #333333;
}

.tab-item.active .tab-text {
  color: #FFFFFF;
}

.tab-text {
  font-size: 28rpx;
  color: #666666;
}

/* 地址卡片 */
.address-card {
  background-color: #FFFFFF;
  padding: 30rpx;
  margin-bottom: 2rpx;
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: 20rpx;
}

.address-icon {
  width: 56rpx;
  height: 56rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background-color: #333333;
}

.sender-icon {
  background-color: #333333;
}

.receiver-icon {
  background-color: #FF4D4F;
}

.icon-text {
  font-size: 24rpx;
  font-weight: 600;
  color: #FFFFFF;
}

.card-content {
  flex: 1;
  min-height: 56rpx;
  display: flex;
  align-items: center;
}

.address-info {
  width: 100%;
}

.name-phone {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 12rpx;
}

.name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.phone {
  font-size: 28rpx;
  color: #666666;
}

.full-address {
  font-size: 26rpx;
  color: #999999;
  line-height: 1.6;
  display: block;
}

.address-placeholder {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.placeholder-main {
  font-size: 30rpx;
  color: #333333;
}

.placeholder-text {
  font-size: 30rpx;
  color: #333333;
}

.placeholder-sub {
  font-size: 24rpx;
  color: #999999;
}

.card-action {
  flex-shrink: 0;
  padding-left: 20rpx;
}

.action-text {
  font-size: 26rpx;
  color: #666666;
}

/* 服务模式选择 */
.service-mode {
  background-color: #FFFFFF;
  padding: 20rpx 30rpx 30rpx;
  margin-bottom: 2rpx;
  display: flex;
  gap: 30rpx;
}

.mode-item {
  flex: 1;
  padding: 40rpx 20rpx;
  border: 2rpx solid #E5E5E5;
  border-radius: 12rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  background-color: #FAFAFA;
  position: relative;
}

.mode-item.active {
  border-color: #FF4D4F;
  background-color: #FFFFFF;
}

.mode-item.active::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 6rpx;
  background-color: #FF4D4F;
  border-radius: 0 0 10rpx 10rpx;
}

.mode-icon {
  font-size: 48rpx;
}

.mode-text {
  font-size: 28rpx;
  color: #333333;
  font-weight: 500;
}

/* 表单行 */
.form-row {
  background-color: #FFFFFF;
  padding: 32rpx 30rpx;
  margin-bottom: 2rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-height: 100rpx;
}

.form-label-wrapper {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.form-label {
  font-size: 30rpx;
  color: #333333;
  font-weight: 500;
}

.required {
  font-size: 20rpx;
  color: #FFFFFF;
  background-color: #FF4D4F;
  padding: 3rpx 10rpx;
  border-radius: 4rpx;
}

.form-value {
  display: flex;
  align-items: center;
  gap: 12rpx;
  flex: 1;
  justify-content: flex-end;
}

.form-text {
  font-size: 28rpx;
  color: #666666;
}

.time-badge {
  background-color: #F5F5F5;
  color: #333333;
  font-size: 24rpx;
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
}

.payment-tip {
  font-size: 22rpx;
  color: #999999;
}

.price-text {
  font-size: 24rpx;
  color: #999999;
}

.arrow {
  font-size: 28rpx;
  color: #CCCCCC;
  font-weight: 300;
}

.form-input-inline {
  flex: 1;
  font-size: 28rpx;
  color: #333333;
  text-align: right;
}

.placeholder-style {
  color: #CCCCCC;
}

/* 协议 */
.agreement {
  background-color: transparent;
  padding: 30rpx;
  margin-top: 20rpx;
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.agreement-checkbox {
  transform: scale(0.85);
}

.agreement-text {
  font-size: 24rpx;
  color: #666666;
}

/* 底部操作栏 */
.footer {
  display: flex;
  align-items: center;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #333333;
  gap: 30rpx;
}

.fee-section {
  flex: 1;
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.fee-label {
  font-size: 26rpx;
  color: #FFFFFF;
}

.fee-amount {
  font-size: 32rpx;
  color: #FFFFFF;
  font-weight: 600;
}

.fee-detail {
  font-size: 24rpx;
  color: #CCCCCC;
  margin-left: 10rpx;
}

.submit-btn {
  width: 240rpx;
  height: 80rpx;
  background-color: #FFFFFF;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.unit-text {
  font-size: 28rpx;
  color: #666666;
  margin-left: 8rpx;
}

/* 费用明细弹窗 */
.fee-detail-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: flex-end;
  z-index: 1000;
}

.fee-detail-content {
  width: 100%;
  background-color: #FFFFFF;
  border-radius: 32rpx 32rpx 0 0;
  padding-bottom: env(safe-area-inset-bottom);
}

.fee-detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 40rpx 30rpx 20rpx;
  border-bottom: 1rpx solid #F0F0F0;
}

.fee-detail-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.fee-detail-close {
  font-size: 40rpx;
  color: #999999;
  padding: 0 10rpx;
}

.fee-detail-body {
  padding: 30rpx;
}

.fee-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 0;
}

.fee-item-label {
  font-size: 28rpx;
  color: #666666;
}

.fee-item-value {
  font-size: 28rpx;
  color: #333333;
  font-weight: 500;
}

.fee-divider {
  height: 1rpx;
  background-color: #F0F0F0;
  margin: 20rpx 0;
}

.fee-total {
  padding: 30rpx 0 20rpx;
}

.fee-total .fee-item-label {
  font-size: 32rpx;
  color: #333333;
  font-weight: 600;
}

.fee-total .fee-item-value.total {
  font-size: 36rpx;
  color: #FF4D4F;
  font-weight: 600;
}

.fee-tip {
  background-color: #F5F5F5;
  padding: 20rpx;
  border-radius: 8rpx;
  margin-top: 20rpx;
}

.fee-tip-text {
  font-size: 24rpx;
  color: #666666;
}

/* 保价弹窗 */
.insurance-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.insurance-content {
  width: 600rpx;
  background-color: #FFFFFF;
  border-radius: 24rpx;
  overflow: hidden;
}

.insurance-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 40rpx 30rpx 20rpx;
  border-bottom: 1rpx solid #F0F0F0;
}

.insurance-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
}

.insurance-close {
  font-size: 40rpx;
  color: #999999;
  padding: 0 10rpx;
}

.insurance-body {
  padding: 30rpx;
}

.insurance-switch {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 0;
}

.insurance-label {
  font-size: 28rpx;
  color: #333333;
}

.insurance-input-wrapper {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 30rpx 0;
}

.insurance-input-label {
  font-size: 28rpx;
  color: #333333;
  white-space: nowrap;
}

.insurance-input {
  flex: 1;
  height: 80rpx;
  border: 1rpx solid #E5E5E5;
  border-radius: 8rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
}

.insurance-unit {
  font-size: 28rpx;
  color: #666666;
}

.insurance-tip {
  background-color: #FFF7E6;
  padding: 20rpx;
  border-radius: 8rpx;
}

.insurance-tip-text {
  font-size: 24rpx;
  color: #FF8C00;
}

.insurance-footer {
  padding: 20rpx 30rpx 30rpx;
}

.insurance-btn {
  height: 80rpx;
  background-color: #FF4D4F;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.insurance-btn-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #FFFFFF;
}
</style>
