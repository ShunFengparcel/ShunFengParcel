<template>
  <view class="container">
    <!-- 搜索栏 -->
    <view class="search-bar">
      <view class="search-input" @click="handleSearch">
        <text class="search-icon">🔍</text>
        <input 
          v-model="searchKeyword" 
          class="search-field" 
          placeholder="运单号查询/关键字索引"
          @confirm="handleSearchConfirm"
        />
      </view>
      <view class="scan-btn" @click="handleScan">
        <text class="scan-icon">📷</text>
      </view>
    </view>

    <!-- Tab 切换 -->
    <view class="tabs">
      <view 
        v-for="(tab, index) in tabs" 
        :key="index"
        class="tab-item"
        :class="{ active: currentTab === index }"
        @click="switchTab(index)"
      >
        <text class="tab-text">{{ tab.label }}</text>
        <text v-if="tab.count > 0" class="tab-count">{{ tab.count }}</text>
        <view v-if="currentTab === index" class="tab-line"></view>
      </view>
      <view class="filter-btn" @click="showFilter">
        <text>筛选</text>
        <text class="filter-icon">▼</text>
      </view>
    </view>

    <!-- 订单列表 -->
    <scroll-view 
      class="order-list" 
      scroll-y 
      @scrolltolower="loadMore"
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
    >
      <view v-if="orderList.length > 0">
        <view 
          v-for="order in orderList" 
          :key="order.id"
          class="order-card"
          @click="goToDetail(order.id)"
        >
          <view class="order-header">
            <view class="order-status">
              <text class="status-icon">{{ getStatusIcon(order.status) }}</text>
              <text class="status-text">{{ getStatusText(order.status) }}</text>
            </view>
            <text class="order-time">{{ formatTime(order.createdAt) }}</text>
          </view>

          <view class="order-content">
            <view class="address-info">
              <view class="address-item">
                <text class="address-label">寄</text>
                <text class="address-text">{{ order.senderAddress }}</text>
              </view>
              <view class="address-arrow">→</view>
              <view class="address-item">
                <text class="address-label">收</text>
                <text class="address-text">{{ order.receiverAddress }}</text>
              </view>
            </view>

            <view class="order-info">
              <text class="order-no">运单号：{{ order.orderNo }}</text>
              <text class="order-service">{{ order.serviceType }}</text>
            </view>
          </view>

          <view class="order-footer">
            <text class="order-fee">¥{{ order.estimatedFee }}</text>
            <view class="order-actions">
              <view v-if="order.status === 'pending'" class="action-btn cancel" @click.stop="cancelOrder(order.id)">
                取消订单
              </view>
              <view class="action-btn primary">查看详情</view>
            </view>
          </view>
        </view>
      </view>

      <!-- 空状态 -->
      <sf-empty 
        v-else
        text="您最近没有寄出快件，快去寄件吧~"
        :show-button="true"
        button-text="去寄件"
        @click="goToSend"
      >
        <view class="empty-image">📦</view>
      </sf-empty>

      <!-- 加载更多 -->
      <view v-if="orderList.length > 0 && hasMore" class="load-more">
        <text>加载中...</text>
      </view>
      <view v-if="orderList.length > 0 && !hasMore" class="load-more">
        <text>没有更多了</text>
      </view>
    </scroll-view>

    <!-- 客服浮动按钮 -->
    <view class="customer-service-btn">
      <image class="service-avatar" src="/static/images/customer-service.png" mode="aspectFill"></image>
      <text class="service-text">客服中心</text>
    </view>

    <!-- 历史查询提示 -->
    <view class="history-tip">
      <text>如需查询3个月以上运单</text>
      <text class="tip-link" @click="goToHistory">点击这里</text>
    </view>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import SfEmpty from '@/components/sf-empty/sf-empty.vue'

const orderStore = useOrderStore()

// 搜索关键词
const searchKeyword = ref('')

// Tab 数据
const tabs = ref([
  { label: '寄件', count: 0, type: 'send' },
  { label: '收件', count: 0, type: 'receive' },
  { label: '待支付', count: 0, type: 'unpaid' }
])

const currentTab = ref(0)

// 订单列表
const orderList = ref([])
const refreshing = ref(false)
const hasMore = ref(true)
const page = ref(1)

// 计算属性
const currentTabType = computed(() => tabs.value[currentTab.value].type)

// 方法
function switchTab(index) {
  currentTab.value = index
  page.value = 1
  orderList.value = []
  loadOrders()
}

function handleSearch() {
  console.log('搜索：', searchKeyword.value)
}

function handleSearchConfirm() {
  if (searchKeyword.value.trim()) {
    // 跳转到物流详情页
    uni.navigateTo({
      url: `/pages/order/tracking?trackingNo=${searchKeyword.value}`
    })
  }
}

function handleScan() {
  uni.scanCode({
    success: (res) => {
      searchKeyword.value = res.result
      handleSearchConfirm()
    }
  })
}

function showFilter() {
  uni.showActionSheet({
    itemList: ['全部', '待接单', '已接单', '运输中', '已完成'],
    success: (res) => {
      console.log('选择筛选：', res.tapIndex)
    }
  })
}

async function loadOrders() {
  try {
    console.log('🔵 [订单列表] 开始加载订单...')
    console.log('🔵 [订单列表] 当前页码:', page.value)
    console.log('🔵 [订单列表] 当前Tab:', currentTabType.value)
    
    uni.showLoading({
      title: '加载中...',
      mask: true
    })
    
    // 调用 API 获取当前用户的订单列表（自动根据登录用户ID过滤）
    console.log('🔵 [订单列表] 调用 orderStore.getOrderList...')
    const response = await orderStore.getOrderList('all', page.value, 10)
    
    console.log('🔵 [订单列表] API响应:', response)
    
    if (response && response.orders) {
      console.log('🔵 [订单列表] 收到订单数据:', response.orders.length, '条')
      
      // 处理订单数据，确保地址格式正确
      const processedOrders = response.orders.map(order => {
        console.log('🔵 [订单列表] 处理订单:', order.order_no || order.orderNo)
        return {
          ...order,
          senderAddress: order.sender_address || order.senderAddress || '未知地址',
          receiverAddress: order.receiver_address || order.receiverAddress || '未知地址',
          orderNo: order.order_no || order.orderNo || '',
          serviceType: order.service_type || order.serviceType || '顺丰标快',
          estimatedFee: order.estimated_fee || order.estimatedFee || 0,
          createdAt: order.created_at || order.createdAt || new Date().toISOString()
        }
      })
      
      if (page.value === 1) {
        orderList.value = processedOrders
      } else {
        orderList.value = [...orderList.value, ...processedOrders]
      }
      
      console.log('🔵 [订单列表] 最终订单列表:', orderList.value)
      
      // 更新 tab 计数（只显示当前用户的订单数量）
      tabs.value[0].count = response.total
      
      // 判断是否还有更多数据
      hasMore.value = orderList.value.length < response.total
    } else {
      console.log('🔴 [订单列表] 没有收到订单数据')
    }
    
    uni.hideLoading()
  } catch (error) {
    uni.hideLoading()
    console.error('🔴 [订单列表] 加载订单列表失败:', error)
    uni.showToast({
      title: '加载失败，请重试',
      icon: 'none'
    })
  }
}

function onRefresh() {
  refreshing.value = true
  page.value = 1
  orderList.value = []
  
  setTimeout(() => {
    loadOrders()
    refreshing.value = false
  }, 1000)
}

function loadMore() {
  if (!hasMore.value) return
  
  page.value++
  loadOrders()
}

function goToDetail(orderId) {
  // 找到对应的订单
  const order = orderList.value.find(o => o.id === orderId)
  if (order && order.orderNo) {
    // 跳转到物流轨迹页面
    uni.navigateTo({
      url: `/pages/order/tracking?orderId=${orderId}`
    })
  } else {
    uni.showToast({
      title: '订单信息不完整',
      icon: 'none'
    })
  }
}

function cancelOrder(orderId) {
  uni.showModal({
    title: '提示',
    content: '确定要取消订单吗？',
    success: (res) => {
      if (res.confirm) {
        orderStore.cancelOrder(orderId).then(() => {
          uni.showToast({
            title: '订单已取消',
            icon: 'success'
          })
          loadOrders()
        })
      }
    }
  })
}

function goToSend() {
  uni.switchTab({
    url: '/pages/send/index'
  })
}

function goToHistory() {
  uni.showToast({
    title: '历史订单查询功能开发中',
    icon: 'none'
  })
}

function getStatusIcon(status) {
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
  return icons[status] || '📦'
}

function getStatusText(status) {
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
  return texts[status] || '未知状态'
}

function formatTime(time) {
  // 简单的时间格式化
  return time.split(' ')[0]
}

onMounted(() => {
  console.log('🚀🚀🚀 订单列表页面已加载 🚀🚀🚀')
  uni.showToast({
    title: '正在加载订单...',
    icon: 'loading',
    duration: 1000
  })
  loadOrders()
})
</script>

<style scoped>
.container {
  min-height: 100vh;
  background-color: #F8F8F8;
  padding-bottom: 120rpx;
}

/* 搜索栏 */
.search-bar {
  display: flex;
  align-items: center;
  padding: 20rpx 30rpx;
  background-color: #FFFFFF;
  gap: 20rpx;
}

.search-input {
  flex: 1;
  display: flex;
  align-items: center;
  height: 70rpx;
  background-color: #F5F5F5;
  border-radius: 35rpx;
  padding: 0 30rpx;
  gap: 15rpx;
}

.search-icon {
  font-size: 32rpx;
  color: #999999;
}

.search-field {
  flex: 1;
  font-size: 28rpx;
  height: 100%;
}

.scan-btn {
  width: 70rpx;
  height: 70rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 40rpx;
}

/* Tab 切换 */
.tabs {
  display: flex;
  align-items: center;
  background-color: #FFFFFF;
  padding: 0 30rpx;
  border-bottom: 1rpx solid #E5E5E5;
}

.tab-item {
  position: relative;
  padding: 30rpx 40rpx;
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.tab-item.active .tab-text {
  color: #D81E06;
  font-weight: 600;
}

.tab-text {
  font-size: 30rpx;
  color: #666666;
}

.tab-count {
  font-size: 24rpx;
  color: #999999;
}

.tab-line {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 60rpx;
  height: 6rpx;
  background-color: #D81E06;
  border-radius: 3rpx;
}

.filter-btn {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10rpx;
  font-size: 28rpx;
  color: #666666;
}

.filter-icon {
  font-size: 20rpx;
}

/* 订单列表 */
.order-list {
  height: calc(100vh - 200rpx);
  padding: 20rpx 30rpx;
}

.order-card {
  background-color: #FFFFFF;
  border-radius: 20rpx;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.order-status {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.status-icon {
  font-size: 32rpx;
}

.status-text {
  font-size: 28rpx;
  font-weight: 600;
  color: #333333;
}

.order-time {
  font-size: 24rpx;
  color: #999999;
}

.order-content {
  margin-bottom: 20rpx;
}

.address-info {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.address-item {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.address-label {
  width: 40rpx;
  height: 40rpx;
  background-color: #F5F5F5;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 22rpx;
  color: #666666;
}

.address-text {
  flex: 1;
  font-size: 26rpx;
  color: #666666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.address-arrow {
  font-size: 24rpx;
  color: #CCCCCC;
}

.order-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-no {
  font-size: 24rpx;
  color: #999999;
}

.order-service {
  font-size: 24rpx;
  color: #D81E06;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20rpx;
  border-top: 1rpx solid #F5F5F5;
}

.order-fee {
  font-size: 32rpx;
  font-weight: 600;
  color: #D81E06;
}

.order-actions {
  display: flex;
  gap: 20rpx;
}

.action-btn {
  padding: 12rpx 30rpx;
  border-radius: 30rpx;
  font-size: 24rpx;
  border: 1rpx solid #E5E5E5;
  color: #666666;
}

.action-btn.cancel {
  color: #999999;
}

.action-btn.primary {
  background-color: #D81E06;
  color: #FFFFFF;
  border-color: #D81E06;
}

/* 空状态 */
.empty-image {
  font-size: 200rpx;
  opacity: 0.3;
}

/* 加载更多 */
.load-more {
  text-align: center;
  padding: 30rpx;
  font-size: 24rpx;
  color: #999999;
}

/* 客服按钮 */
.customer-service-btn {
  position: fixed;
  right: 30rpx;
  bottom: 200rpx;
  width: 120rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10rpx;
}

.service-avatar {
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background-color: #FFFFFF;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
}

.service-text {
  font-size: 22rpx;
  color: #666666;
}

/* 历史查询提示 */
.history-tip {
  position: fixed;
  bottom: 120rpx;
  left: 0;
  right: 0;
  text-align: center;
  font-size: 24rpx;
  color: #999999;
}

.tip-link {
  color: #D81E06;
  margin-left: 10rpx;
}
</style>
