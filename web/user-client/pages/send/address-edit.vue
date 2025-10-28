<template>
  <view class="container">
    <view class="form">
      <!-- 地址类型 -->
      <view class="form-section">
        <view class="section-title">地址类型</view>
        <view class="type-tabs">
          <view 
            class="type-tab"
            :class="{ active: formData.type === 'sender' }"
            @click="formData.type = 'sender'"
          >
            寄件地址
          </view>
          <view 
            class="type-tab"
            :class="{ active: formData.type === 'receiver' }"
            @click="formData.type = 'receiver'"
          >
            收件地址
          </view>
        </view>
      </view>

      <!-- 联系人信息 -->
      <view class="form-section">
        <view class="form-item required">
          <text class="label">联系人</text>
          <input 
            v-model="formData.name" 
            class="input" 
            placeholder="请输入姓名（2-20字符）"
            maxlength="20"
          />
        </view>

        <view class="form-item required">
          <text class="label">手机号码</text>
          <input 
            v-model="formData.phone" 
            class="input" 
            type="number"
            placeholder="请输入11位手机号"
            maxlength="11"
          />
        </view>
      </view>

      <!-- 地址信息 -->
      <view class="form-section">
        <picker 
          mode="region"
          :value="regionValue"
          @change="handleRegionChange"
        >
          <view class="form-item required">
            <text class="label">所在地区</text>
            <view class="input-wrapper">
              <text v-if="regionText" class="input-text">{{ regionText }}</text>
              <text v-else class="placeholder">请选择省市区</text>
              <text class="arrow">></text>
            </view>
          </view>
        </picker>

        <view class="form-item required">
          <text class="label">详细地址</text>
          <textarea 
            v-model="formData.detail" 
            class="textarea" 
            placeholder="请输入详细地址（5-120字符）"
            maxlength="120"
            :show-confirm-bar="false"
          />
        </view>
      </view>

      <!-- 可选信息 -->
      <view class="form-section">
        <view class="form-item">
          <text class="label">公司名称</text>
          <input 
            v-model="formData.company" 
            class="input" 
            placeholder="选填"
          />
        </view>

        <view class="form-item" @click="showTagPicker">
          <text class="label">地址标签</text>
          <view class="input-wrapper">
            <text v-if="formData.tag" class="input-text">{{ formData.tag }}</text>
            <text v-else class="placeholder">选填</text>
            <text class="arrow">></text>
          </view>
        </view>

        <view class="form-item">
          <text class="label">设为默认地址</text>
          <switch 
            :checked="formData.isDefault" 
            @change="handleDefaultChange"
            color="#D81E06"
          />
        </view>
      </view>
    </view>

    <!-- 保存按钮 -->
    <view class="save-btn" @click="saveAddress">
      保存地址
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

// 表单数据
const formData = ref({
  id: null,
  type: 'sender',
  name: '',
  phone: '',
  province: '',
  city: '',
  district: '',
  detail: '',
  company: '',
  tag: '',
  isDefault: false
})

// 地区选择
const regionValue = ref([])

// 计算属性
const regionText = computed(() => {
  if (formData.value.province && formData.value.city && formData.value.district) {
    return `${formData.value.province} ${formData.value.city} ${formData.value.district}`
  }
  return ''
})

// 方法
function handleRegionChange(e) {
  const [province, city, district] = e.detail.value
  formData.value.province = province
  formData.value.city = city
  formData.value.district = district
  regionValue.value = [province, city, district]
}

function showTagPicker() {
  uni.showActionSheet({
    itemList: ['家', '公司', '学校', '其他'],
    success: (res) => {
      const tags = ['家', '公司', '学校', '其他']
      formData.value.tag = tags[res.tapIndex]
    }
  })
}

function handleDefaultChange(e) {
  formData.value.isDefault = e.detail.value
}

function validateForm() {
  if (!formData.value.name || formData.value.name.length < 2) {
    uni.showToast({
      title: '请输入正确的姓名',
      icon: 'none'
    })
    return false
  }

  if (!formData.value.phone || !/^1[3-9]\d{9}$/.test(formData.value.phone)) {
    uni.showToast({
      title: '请输入正确的手机号',
      icon: 'none'
    })
    return false
  }

  if (!formData.value.province || !formData.value.city || !formData.value.district) {
    uni.showToast({
      title: '请选择所在地区',
      icon: 'none'
    })
    return false
  }

  if (!formData.value.detail || formData.value.detail.length < 5) {
    uni.showToast({
      title: '请输入详细地址（至少5个字符）',
      icon: 'none'
    })
    return false
  }

  return true
}

function saveAddress() {
  if (!validateForm()) {
    return
  }

  // TODO: 调用地址验证 API（高德地图）
  // TODO: 调用保存地址 API

  uni.showToast({
    title: '保存成功',
    icon: 'success'
  })

  setTimeout(() => {
    uni.navigateBack()
  }, 1500)
}

onLoad((options) => {
  if (options.id) {
    // 编辑模式，加载地址数据
    formData.value.id = parseInt(options.id)
    // TODO: 加载地址详情
  }
  
  if (options.type) {
    formData.value.type = options.type
  }
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F8F8F8;
  padding-bottom: 120rpx;
}

/* 表单 */
.form {
  padding: 20rpx 0;
}

.form-section {
  background-color: #FFFFFF;
  margin-bottom: 20rpx;
  padding: 0 30rpx;
}

.section-title {
  padding: 30rpx 0 20rpx;
  font-size: 28rpx;
  color: #999999;
}

/* 类型选择 */
.type-tabs {
  display: flex;
  gap: 20rpx;
  padding-bottom: 30rpx;
}

.type-tab {
  flex: 1;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #F5F5F5;
  border-radius: 40rpx;
  font-size: 28rpx;
  color: #666666;
}

.type-tab.active {
  background-color: #FFE5E5;
  color: #D81E06;
  font-weight: 600;
}

/* 表单项 */
.form-item {
  display: flex;
  align-items: center;
  min-height: 100rpx;
  border-bottom: 1rpx solid #F5F5F5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-item.required .label::before {
  content: '*';
  color: #D81E06;
  margin-right: 8rpx;
}

.label {
  width: 180rpx;
  font-size: 28rpx;
  color: #333333;
}

.input {
  flex: 1;
  font-size: 28rpx;
  color: #333333;
}

.input-wrapper {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.input-text {
  flex: 1;
  font-size: 28rpx;
  color: #333333;
}

.placeholder {
  flex: 1;
  font-size: 28rpx;
  color: #CCCCCC;
}

.arrow {
  font-size: 24rpx;
  color: #CCCCCC;
}

.textarea {
  flex: 1;
  min-height: 120rpx;
  padding: 20rpx 0;
  font-size: 28rpx;
  color: #333333;
  line-height: 1.6;
}

/* 保存按钮 */
.save-btn {
  position: fixed;
  bottom: 30rpx;
  left: 30rpx;
  right: 30rpx;
  height: 90rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 45rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  font-weight: 600;
  box-shadow: 0 8rpx 30rpx rgba(216, 30, 6, 0.3);
}
</style>
