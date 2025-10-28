import { defineStore } from 'pinia'
import { ref } from 'vue'
import { post, get } from '@/utils/request'
import { useUserStore } from './user'

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

  async function createOrder(orderData) {
    try {
      const userStore = useUserStore()
      
      // 准备请求数据
      const requestData = {
        user_id: userStore.userInfo?.id || 1,
        sender_name: orderData.senderName,
        sender_phone: orderData.senderPhone,
        sender_address: orderData.senderAddress,
        receiver_name: orderData.receiverName,
        receiver_phone: orderData.receiverPhone,
        receiver_address: orderData.receiverAddress,
        pickup_mode: orderData.pickupMode || 'pickup',
        pickup_time: orderData.pickupTime || '今天 一小时内',
        service_type: orderData.serviceType || '顺丰标快',
        product_type: orderData.productType,
        weight: orderData.estimatedWeight || 0,
        payment_method: orderData.paymentMethod || '寄付现结',
        estimated_fee: orderData.estimatedFee || 0,
        is_insured: orderData.isInsured ? 1 : 0,
        insured_value: orderData.insuredValue || 0
      }
      
      // 调用创建订单 API
      const response = await post('/api/v1/orders', requestData)
      
      // 构建完整的订单对象
      const order = {
        id: response.id,
        orderNo: response.order_no,
        pickupCode: response.pickup_code,
        status: response.status,
        ...orderData,
        createdAt: new Date().toISOString()
      }
      
      orderList.value.unshift(order)
      setCurrentOrder(order)
      
      return order
    } catch (error) {
      console.error('创建订单失败:', error)
      throw error
    }
  }

  async function getOrderList(type = 'all', page = 1, pageSize = 10) {
    try {
      const userStore = useUserStore()
      
      // 确保用户已登录
      if (!userStore.isLogin || !userStore.userInfo?.id) {
        console.error('🔴 [OrderStore] 用户未登录')
        uni.showToast({
          title: '请先登录',
          icon: 'none'
        })
        // 跳转到登录页
        setTimeout(() => {
          uni.reLaunch({
            url: '/pages/login/login'
          })
        }, 1500)
        return { orders: [], total: 0 }
      }
      
      const userId = userStore.userInfo.id
      
      console.log('🟢 [OrderStore] 开始获取订单列表')
      console.log('🟢 [OrderStore] 用户ID:', userId)
      console.log('🟢 [OrderStore] 类型:', type, '页码:', page, '每页:', pageSize)
      
      // 构建查询参数
      const params = {
        user_id: userId,
        page,
        page_size: pageSize
      }
      
      // 根据类型设置状态过滤
      if (type !== 'all') {
        params.status = type
      }
      
      console.log('🟢 [OrderStore] 请求参数:', params)
      console.log('🟢 [OrderStore] 请求URL: /api/v1/orders')
      
      // 调用获取订单列表 API
      const response = await get('/api/v1/orders', params)
      
      console.log('🟢 [OrderStore] API响应:', response)
      
      if (response && response.orders) {
        console.log('🟢 [OrderStore] 成功获取', response.orders.length, '条订单')
        setOrderList(response.orders)
        return response
      }
      
      console.log('🟡 [OrderStore] 响应中没有订单数据')
      return { orders: [], total: 0 }
    } catch (error) {
      console.error('🔴 [OrderStore] 获取订单列表失败:', error)
      
      // 如果是401错误，跳转到登录页
      if (error.code === 401) {
        uni.showToast({
          title: '登录已过期，请重新登录',
          icon: 'none'
        })
        setTimeout(() => {
          uni.reLaunch({
            url: '/pages/login/login'
          })
        }, 1500)
      }
      
      return { orders: [], total: 0 }
    }
  }

  async function getOrderDetail(orderId, orderNo) {
    try {
      // 构建查询参数
      const params = {}
      if (orderId) {
        params.id = orderId
      } else if (orderNo) {
        params.order_no = orderNo
      } else {
        throw new Error('订单ID或订单号不能为空')
      }
      
      // 调用获取订单详情 API
      const order = await get('/api/v1/orders/detail', params)
      
      if (order) {
        setCurrentOrder(order)
      }
      
      return order
    } catch (error) {
      console.error('获取订单详情失败:', error)
      throw error
    }
  }

  async function getTrackingInfo(orderId) {
    try {
      // 调用获取物流信息 API
      const response = await get('/api/v1/tracking', { order_id: orderId })
      
      if (response && response.trackings) {
        const tracking = {
          orderNo: currentOrder.value?.orderNo || '',
          status: currentOrder.value?.status || '',
          estimatedDeliveryTime: '',
          logs: response.trackings.map(t => ({
            time: t.created_at,
            location: t.location,
            description: t.description,
            status: t.status
          }))
        }
        
        setTrackingInfo(tracking)
        return tracking
      }
      
      return null
    } catch (error) {
      console.error('获取物流信息失败:', error)
      throw error
    }
  }

  async function cancelOrder(orderId) {
    try {
      // TODO: 实现取消订单API
      const index = orderList.value.findIndex(o => o.id === orderId)
      if (index !== -1) {
        orderList.value[index].status = 'cancelled'
      }
      return true
    } catch (error) {
      console.error('取消订单失败:', error)
      throw error
    }
  }

  async function calculatePrice(pricingData) {
    try {
      console.log('💰 [OrderStore] 开始计价:', pricingData)
      
      const requestData = {
        sender_city: pricingData.senderCity,
        receiver_city: pricingData.receiverCity,
        weight: pricingData.weight,
        service_type: pricingData.serviceType || '顺丰标快',
        is_insured: pricingData.isInsured ? 1 : 0,
        insured_value: pricingData.insuredValue || 0
      }
      
      const response = await post('/api/v1/pricing/calculate', requestData)
      
      console.log('💰 [OrderStore] 计价结果:', response)
      
      return response
    } catch (error) {
      console.error('💰 [OrderStore] 计价失败:', error)
      throw error
    }
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
    cancelOrder,
    calculatePrice
  }
})
