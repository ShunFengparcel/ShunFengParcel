<template>
  <view class="container">
    <!-- 顶部输入区域 -->
    <view class="input-section">
      <view class="input-header">
        <text class="input-title">粘贴识别</text>
        <text class="input-subtitle">或输入文本、智能拆分姓名、电话和地址</text>
      </view>

      <textarea v-model="inputText" class="input-textarea" placeholder="粘贴或输入收件人信息" :auto-height="true"
        :maxlength="500" />

      <view class="input-actions">
        <view class="action-btn" @click="handleImageRecognition">
          <text class="action-icon">📷</text>
          <text class="action-text">图片识别</text>
        </view>
        <view class="action-btn" @click="handleAddressCode">
          <text class="action-icon">📋</text>
          <text class="action-text">地址码</text>
        </view>
        <view class="action-btn primary" @click="handlePasteRecognition">
          <text class="action-text">粘贴并识别</text>
        </view>
      </view>
    </view>

    <!-- 手动输入表单 -->
    <view class="form-section">
      <view class="form-header">
        <view class="header-tabs">
          <view class="tab-item active">
            <text class="tab-icon">📍</text>
            <text class="tab-text">收件人</text>
          </view>
          <view class="tab-item">
            <text class="tab-text">地址簿</text>
          </view>
          <view class="tab-item">
            <text class="tab-text">顺丰ID</text>
          </view>
        </view>
      </view>

      <view class="form-content">
        <view class="form-row">
          <text class="form-label">姓名</text>
          <input v-model="formData.name" class="form-input" placeholder="请输入姓名" />
          <text class="form-divider">-</text>
          <text class="form-label">分机号</text>
          <input v-model="formData.extension" class="form-input small" placeholder="" />
          <view class="contacts-btn" @click="selectFromContacts">
            <text class="contacts-icon">📖</text>
            <text class="contacts-text">通讯录</text>
          </view>
        </view>

        <view class="form-row">
          <text class="form-label">电话</text>
          <input v-model="formData.phone" class="form-input" type="number" placeholder="请输入电话号码" />
        </view>

        <picker 
          mode="region" 
          :value="regionValue"
          @change="onRegionChange"
        >
          <view class="form-row clickable">
            <text class="form-label">省市区</text>
            <view class="form-value">
              <text v-if="formData.region" class="value-text">{{ formData.region }}</text>
              <text v-else class="placeholder-text">请选择</text>
              <text class="arrow">›</text>
            </view>
            <view class="location-btn" @click.stop="getCurrentLocation">
              <text class="location-icon">📍</text>
            </view>
          </view>
        </picker>

        <view class="form-row">
          <text class="form-label">详细地址</text>
          <input v-model="formData.detail" class="form-input" placeholder="例如**街**号**" />
        </view>

        <view class="form-row">
          <text class="form-label">公司名称</text>
          <input v-model="formData.company" class="form-input" placeholder="（选填）" />
        </view>

        <view class="form-row">
          <text class="form-label">地址标签</text>
          <input v-model="formData.tag" class="form-input" placeholder="（选填）" />
          <view class="add-tag-btn">
            <text>+</text>
          </view>
        </view>

        <view class="form-row checkbox-row">
          <checkbox :checked="formData.saveToAddressBook"
            @change="formData.saveToAddressBook = !formData.saveToAddressBook" />
          <text class="checkbox-label">保存到地址簿</text>
          <view class="clear-btn" @click="clearForm">
            <text>清空</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 历史地址 -->
    <view v-if="historyAddresses.length > 0" class="history-section">
      <view class="history-header">
        <text class="history-title">历史地址</text>
      </view>
      <view v-for="address in historyAddresses" :key="address.id" class="history-item"
        @click="selectHistoryAddress(address)">
        <view class="history-info">
          <text class="history-name">{{ address.name }}</text>
          <text class="history-phone">{{ formatPhone(address.phone) }}</text>
        </view>
        <text class="history-address">
          {{ address.province }}{{ address.city }}{{ address.district }}{{ address.detail }}
        </text>
      </view>
    </view>

    <!-- 底部按钮 -->
    <view class="footer">
      <view class="footer-btn cancel" @click="handleCancel">
        <text>邀请填写</text>
      </view>
      <view class="footer-btn confirm" @click="handleConfirm">
        <text>确定</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

// 接收参数
const addressType = ref('receiver') // sender 或 receiver

onLoad((options) => {
  if (options.type) {
    addressType.value = options.type
  }
})

// 输入文本
const inputText = ref('')

// 省市区选择器
const regionValue = ref([])

// 表单数据
const formData = ref({
  name: '',
  phone: '',
  extension: '',
  region: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  company: '',
  tag: '',
  saveToAddressBook: true
})

// 历史地址
const historyAddresses = ref([
  {
    id: 1,
    name: '康南南',
    phone: '18200004509',
    province: '上海市',
    city: '上海市',
    district: '浦东新区',
    detail: '盐大路2733-2735号海佳彩亮(上海运营中心)'
  },
  {
    id: 2,
    name: '康楠楠',
    phone: '18200004509',
    province: '上海市',
    city: '上海市',
    district: '浦东新区',
    detail: '盐大路2733号'
  }
])

// 方法
function handleImageRecognition() {
  uni.chooseImage({
    count: 1,
    success: (res) => {
      uni.showToast({
        title: '图片识别功能开发中',
        icon: 'none'
      })
    }
  })
}

function handleAddressCode() {
  uni.showToast({
    title: '地址码功能开发中',
    icon: 'none'
  })
}

function handlePasteRecognition() {
  if (!inputText.value.trim()) {
    uni.showToast({
      title: '请输入或粘贴地址信息',
      icon: 'none'
    })
    return
  }

  // 简单的地址识别逻辑
  const text = inputText.value

  // 提取手机号
  const phoneMatch = text.match(/1[3-9]\d{9}/)
  if (phoneMatch) {
    formData.value.phone = phoneMatch[0]
  }

  // 提取姓名（假设姓名在最前面，2-4个字）
  const nameMatch = text.match(/^[\u4e00-\u9fa5]{2,4}/)
  if (nameMatch) {
    formData.value.name = nameMatch[0]
  }

  uni.showToast({
    title: '识别成功',
    icon: 'success'
  })
}

function selectFromContacts() {
  uni.showToast({
    title: '通讯录功能开发中',
    icon: 'none'
  })
}

function onRegionChange(e) {
  const region = e.detail.value
  formData.value.province = region[0]
  formData.value.city = region[1]
  formData.value.district = region[2]
  formData.value.region = `${region[0]} ${region[1]} ${region[2]}`
  regionValue.value = region
}

function getCurrentLocation() {
  uni.getLocation({
    type: 'gcj02',
    success: (res) => {
      uni.showToast({
        title: '定位成功',
        icon: 'success'
      })
    },
    fail: () => {
      uni.showToast({
        title: '定位失败',
        icon: 'none'
      })
    }
  })
}

function clearForm() {
  formData.value = {
    name: '',
    phone: '',
    extension: '',
    region: '',
    province: '',
    city: '',
    district: '',
    detail: '',
    company: '',
    tag: '',
    saveToAddressBook: true
  }
  inputText.value = ''
}

function selectHistoryAddress(address) {
  formData.value = {
    name: address.name,
    phone: address.phone,
    extension: '',
    region: `${address.province} ${address.city} ${address.district}`,
    province: address.province,
    city: address.city,
    district: address.district,
    detail: address.detail,
    company: '',
    tag: '',
    saveToAddressBook: false
  }
  regionValue.value = [address.province, address.city, address.district]
}

function formatPhone(phone) {
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2')
}

function handleCancel() {
  uni.navigateBack()
}

function handleConfirm() {
  // 验证
  if (!formData.value.name) {
    uni.showToast({ title: '请输入姓名', icon: 'none' })
    return
  }
  if (!formData.value.phone) {
    uni.showToast({ title: '请输入电话', icon: 'none' })
    return
  }
  if (!formData.value.region) {
    uni.showToast({ title: '请选择省市区', icon: 'none' })
    return
  }
  if (!formData.value.detail) {
    uni.showToast({ title: '请输入详细地址', icon: 'none' })
    return
  }

  // 构建地址对象
  const address = {
    name: formData.value.name,
    phone: formData.value.phone,
    province: formData.value.province,
    city: formData.value.city,
    district: formData.value.district,
    detail: formData.value.detail,
    company: formData.value.company,
    tag: formData.value.tag
  }

  // 保存到全局数据
  const app = getApp()
  app.globalData = app.globalData || {}
  app.globalData.selectedAddress = address
  app.globalData.addressType = addressType.value

  // 返回上一页
  uni.navigateBack()
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F5F5F5;
  padding-bottom: 120rpx;
}

/* 输入区域 */
.input-section {
  background-color: #FFFFFF;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.input-header {
  margin-bottom: 20rpx;
}

.input-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
  display: block;
  margin-bottom: 8rpx;
}

.input-subtitle {
  font-size: 24rpx;
  color: #999999;
  display: block;
}

.input-textarea {
  width: 100%;
  min-height: 120rpx;
  padding: 20rpx;
  background-color: #F8F8F8;
  border-radius: 12rpx;
  font-size: 28rpx;
  margin-bottom: 20rpx;
}

.input-actions {
  display: flex;
  gap: 20rpx;
}

.action-btn {
  flex: 1;
  height: 70rpx;
  border: 1rpx solid #E5E5E5;
  border-radius: 35rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  background-color: #FFFFFF;
}

.action-btn.primary {
  background-color: #FF4D4F;
  border-color: #FF4D4F;
}

.action-btn.primary .action-text {
  color: #FFFFFF;
}

.action-icon {
  font-size: 32rpx;
}

.action-text {
  font-size: 26rpx;
  color: #333333;
}

/* 表单区域 */
.form-section {
  background-color: #FFFFFF;
  margin-bottom: 20rpx;
}

.form-header {
  padding: 20rpx 30rpx;
  border-bottom: 1rpx solid #F0F0F0;
}

.header-tabs {
  display: flex;
  gap: 40rpx;
}

.tab-item {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 10rpx 0;
}

.tab-item.active {
  border-bottom: 4rpx solid #FF4D4F;
}

.tab-icon {
  font-size: 28rpx;
}

.tab-text {
  font-size: 28rpx;
  color: #666666;
}

.tab-item.active .tab-text {
  color: #333333;
  font-weight: 600;
}

.form-content {
  padding: 20rpx 30rpx;
}

.form-row {
  display: flex;
  align-items: center;
  padding: 25rpx 0;
  border-bottom: 1rpx solid #F0F0F0;
  gap: 20rpx;
}

.form-row.clickable {
  cursor: pointer;
}

.form-row.checkbox-row {
  border-bottom: none;
  padding-top: 30rpx;
}

.form-label {
  font-size: 28rpx;
  color: #333333;
  min-width: 120rpx;
}

.form-input {
  flex: 1;
  font-size: 28rpx;
  color: #333333;
}

.form-input.small {
  flex: 0.5;
}

.form-divider {
  font-size: 28rpx;
  color: #CCCCCC;
}

.contacts-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  padding: 8rpx 20rpx;
  background-color: #F5F5F5;
  border-radius: 20rpx;
}

.contacts-icon {
  font-size: 24rpx;
}

.contacts-text {
  font-size: 24rpx;
  color: #666666;
}

.form-value {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.value-text {
  font-size: 28rpx;
  color: #333333;
}

.placeholder-text {
  font-size: 28rpx;
  color: #CCCCCC;
}

.arrow {
  font-size: 32rpx;
  color: #CCCCCC;
}

.location-btn {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #F5F5F5;
  border-radius: 50%;
}

.location-icon {
  font-size: 32rpx;
}

.add-tag-btn {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #F5F5F5;
  border-radius: 50%;
  font-size: 32rpx;
  color: #666666;
}

.checkbox-label {
  flex: 1;
  font-size: 26rpx;
  color: #666666;
  margin-left: 10rpx;
}

.clear-btn {
  padding: 8rpx 24rpx;
  background-color: #F5F5F5;
  border-radius: 20rpx;
  font-size: 24rpx;
  color: #666666;
}

/* 历史地址 */
.history-section {
  background-color: #FFFFFF;
  padding: 30rpx;
}

.history-header {
  margin-bottom: 20rpx;
}

.history-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333333;
}

.history-item {
  padding: 25rpx 0;
  border-bottom: 1rpx solid #F0F0F0;
}

.history-info {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 12rpx;
}

.history-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
}

.history-phone {
  font-size: 26rpx;
  color: #666666;
}

.history-address {
  font-size: 26rpx;
  color: #999999;
  line-height: 1.6;
}

/* 底部按钮 */
.footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #FFFFFF;
  box-shadow: 0 -2rpx 20rpx rgba(0, 0, 0, 0.05);
  gap: 20rpx;
}

.footer-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 30rpx;
  font-weight: 600;
}

.footer-btn.cancel {
  background-color: #F5F5F5;
  color: #666666;
}

.footer-btn.confirm {
  background-color: #D81E06;
  color: #FFFFFF;
}
</style>
