<template>
  <view class="container">
    <!-- 地址列表 -->
    <view class="address-list">
      <view 
        v-for="(address, index) in addressList" 
        :key="address.id"
        class="address-card"
        @click="selectAddress(address)"
      >
        <view class="address-header">
          <view class="address-name">
            <text class="name">{{ address.name }}</text>
            <text class="phone">{{ address.phone }}</text>
          </view>
          <view v-if="address.isDefault" class="default-badge">默认</view>
        </view>

        <view class="address-content">
          <view v-if="address.tag" class="address-tag">{{ address.tag }}</view>
          <text class="address-text">
            {{ address.province }} {{ address.city }} {{ address.district }} {{ address.detail }}
          </text>
        </view>

        <view class="address-footer">
          <view class="address-actions">
            <view class="action-btn" @click.stop="editAddress(address)">
              <text class="action-icon">✏️</text>
              <text>编辑</text>
            </view>
            <view class="action-btn" @click.stop="deleteAddress(address.id)">
              <text class="action-icon">🗑️</text>
              <text>删除</text>
            </view>
            <view class="action-btn" @click.stop="setDefault(address.id)">
              <text class="action-icon">⭐</text>
              <text>{{ address.isDefault ? '已默认' : '设为默认' }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <sf-empty 
        v-if="addressList.length === 0"
        text="还没有保存地址"
        :show-button="true"
        button-text="添加地址"
        @click="addAddress"
      />
    </view>

    <!-- 添加地址按钮 -->
    <view class="add-btn" @click="addAddress">
      <text class="add-icon">+</text>
      <text class="add-text">添加新地址</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'
import SfEmpty from '@/components/sf-empty/sf-empty.vue'

const userStore = useUserStore()

// 地址列表
const addressList = ref([
  {
    id: 1,
    name: '张三',
    phone: '138****0000',
    province: '广东省',
    city: '深圳市',
    district: '南山区',
    detail: '科技园南区1001号',
    company: '顺丰科技',
    tag: '公司',
    isDefault: true
  },
  {
    id: 2,
    name: '李四',
    phone: '139****0000',
    province: '北京市',
    city: '北京市',
    district: '朝阳区',
    detail: '建国门外大街1号',
    tag: '家',
    isDefault: false
  }
])

// 方法
function selectAddress(address) {
  try {
    // 将选中的地址保存到全局数据
    const app = getApp()
    if (app) {
      app.globalData = app.globalData || {}
      app.globalData.selectedAddress = address
    }
  } catch (error) {
    console.error('保存地址失败:', error)
  }
  
  // 返回上一页
  uni.navigateBack()
}

function addAddress() {
  uni.navigateTo({
    url: '/pages/send/address-edit'
  })
}

function editAddress(address) {
  uni.navigateTo({
    url: `/pages/send/address-edit?id=${address.id}`
  })
}

function deleteAddress(id) {
  uni.showModal({
    title: '提示',
    content: '确定要删除这个地址吗？',
    success: (res) => {
      if (res.confirm) {
        const index = addressList.value.findIndex(a => a.id === id)
        if (index !== -1) {
          addressList.value.splice(index, 1)
          uni.showToast({
            title: '删除成功',
            icon: 'success'
          })
        }
      }
    }
  })
}

function setDefault(id) {
  addressList.value.forEach(address => {
    address.isDefault = address.id === id
  })
  uni.showToast({
    title: '设置成功',
    icon: 'success'
  })
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F8F8F8;
  padding-bottom: 120rpx;
}

/* 地址列表 */
.address-list {
  padding: 20rpx 30rpx;
}

.address-card {
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.address-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.address-name {
  display: flex;
  align-items: center;
  gap: 20rpx;
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

.default-badge {
  padding: 6rpx 16rpx;
  background-color: #D81E06;
  color: #FFFFFF;
  border-radius: 20rpx;
  font-size: 22rpx;
}

.address-content {
  margin-bottom: 20rpx;
}

.address-tag {
  display: inline-block;
  padding: 4rpx 12rpx;
  background-color: #F5F5F5;
  color: #666666;
  border-radius: 8rpx;
  font-size: 22rpx;
  margin-bottom: 10rpx;
}

.address-text {
  font-size: 28rpx;
  color: #666666;
  line-height: 1.6;
}

.address-footer {
  padding-top: 20rpx;
  border-top: 1rpx solid #F5F5F5;
}

.address-actions {
  display: flex;
  gap: 40rpx;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8rpx;
  font-size: 26rpx;
  color: #666666;
}

.action-icon {
  font-size: 28rpx;
}

/* 添加按钮 */
.add-btn {
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
  gap: 15rpx;
  font-size: 32rpx;
  font-weight: 600;
  box-shadow: 0 8rpx 30rpx rgba(216, 30, 6, 0.3);
}

.add-icon {
  font-size: 40rpx;
}
</style>
