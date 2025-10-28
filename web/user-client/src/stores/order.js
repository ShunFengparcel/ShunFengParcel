import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useOrderStore = defineStore('order', () => {
  // 状态
  const orderList = ref([])
  const currentOrder = ref(null)
  const trackingInfo = ref(null)

  // 方法
  function setOrderList(list) {
    orderList.value = list
  }

  function setCurrentOrder(order) {
    currentOrder.value = order
  }

  function setTrackingInfo(info) {
    trackingInfo.value = info
  }

  function createOrder(orderData) {
    // TODO: 调用创建订单 API
    return new Promise((resolve) => {
      const mockOrder = {
        id: Date.now(),
        orderNo: 'SF' + Date.now(),
        ...orderData,
        status: 'pending',
        pickupCode: Math.random().toString().slice(2, 8),
        createdAt: new Date().toISOString()
      }
      
      orderList.value.unshift(mockOrder)
      setCurrentOrder(mockOrder)
      
      resolve(mockOrder)
    })
  }

  function getOrderList(type = 'all') {
    // TODO: 调用获取订单列表 API
    return Promise.resolve(orderList.value)
  }

  function getOrderDetail(orderId) {
    // TODO: 调用获取订单详情 API
    const order = orderList.value.find(o => o.id === orderId)
    if (order) {
      setCurrentOrder(order)
    }
    return Promise.resolve(order)
  }

  function getTrackingInfo(trackingNo) {
    // TODO: 调用获取物流信息 API
    const mockTracking = {
      orderNo: trackingNo,
      status: 'in_transit',
      estimatedDeliveryTime: '2024-01-20 18:00',
      logs: [
        {
          time: '2024-01-19 14:30',
          location: '深圳市南山区',
          description: '快件已揽收'
        },
        {
          time: '2024-01-19 16:00',
          location: '深圳转运中心',
          description: '快件已到达转运中心'
        }
      ]
    }
    
    setTrackingInfo(mockTracking)
    return Promise.resolve(mockTracking)
  }

  function cancelOrder(orderId) {
    // TODO: 调用取消订单 API
    return new Promise((resolve) => {
      const index = orderList.value.findIndex(o => o.id === orderId)
      if (index !== -1) {
        orderList.value[index].status = 'cancelled'
      }
      resolve(true)
    })
  }

  return {
    // 状态
    orderList,
    currentOrder,
    trackingInfo,
    // 方法
    setOrderList,
    setCurrentOrder,
    setTrackingInfo,
    createOrder,
    getOrderList,
    getOrderDetail,
    getTrackingInfo,
    cancelOrder
  }
})
