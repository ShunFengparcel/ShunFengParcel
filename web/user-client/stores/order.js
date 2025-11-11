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
      const response = await post('/api/v1/create/orders', requestData)
      
      console.log('📦 [OrderStore] 创建订单响应:', response)
      
      // 构建完整的订单对象（后端返回驼峰命名）
      const order = {
        id: response.id,
        orderNo: response.orderNo,
        pickupCode: response.pickupCode,
        status: response.status,
        ...orderData,
        createdAt: new Date().toISOString()
      }
      
      console.log('📦 [OrderStore] 构建的订单对象:', order)
      
      orderList.value.unshift(order)
      setCurrentOrder(order)
      
      return order
    } catch (error) {
      console.error('创建订单失败:', error)
      throw error
    }
  }

  async function getOrderList(type = 'all', page = 1, pageSize = 10) {
    console.log('🟢🟢🟢 [OrderStore] ========== getOrderList 函数被调用 ==========')
    console.log('🟢 [OrderStore] 参数: type=', type, 'page=', page, 'pageSize=', pageSize)
    
    try {
      console.log('🟢 [OrderStore] ===== 开始获取订单列表 =====')
      
      // 临时方案：直接使用硬编码的 user_id = 9，跳过登录检查
      const userId = 9
      
      console.log('🟢 [OrderStore] 🔧 使用硬编码 user_id:', userId)
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
      console.log('🟢 [OrderStore] 请求URL: /api/v1/list/orders')
      
      // 调用获取订单列表 API
      const response = await get('/api/v1/list/orders', params)
      
      console.log('🟢 [OrderStore] ===== API响应分析 =====')
      console.log('🟢 [OrderStore] response类型:', typeof response)
      console.log('🟢 [OrderStore] response是否为null:', response === null)
      console.log('🟢 [OrderStore] response是否为undefined:', response === undefined)
      
      // 尝试序列化
      try {
        const responseStr = JSON.stringify(response)
        console.log('🟢 [OrderStore] response序列化成功，长度:', responseStr.length)
        console.log('🟢 [OrderStore] response内容:', responseStr.substring(0, 200))
      } catch (e) {
        console.error('🟢 [OrderStore] ❌ 无法序列化response:', e.message)
      }
      
      console.log('🟢 [OrderStore] response:', response)
      console.log('🟢 [OrderStore] response.orders:', response.orders)
      console.log('🟢 [OrderStore] response.total:', response.total)
      
      if (response && response.orders && Array.isArray(response.orders)) {
        console.log('🟢 [OrderStore] ✅ 获取到', response.orders.length, '条订单')
        
        if (response.orders.length > 0) {
          console.log('🟢 [OrderStore] 第一条订单:', response.orders[0])
          setOrderList(response.orders)
          return { 
            orders: response.orders, 
            total: response.total || response.orders.length 
          }
        }
      }
      
      console.log('🟡 [OrderStore] 没有订单数据')
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
      let order
      
      if (orderId) {
        // 使用 RESTful 风格的路径
        order = await get(`/api/v1/details/orders/${orderId}`)
      } else if (orderNo) {
        // 使用查询参数
        order = await get('/api/v1/details/orders', { order_no: orderNo })
      } else {
        throw new Error('订单ID或订单号不能为空')
      }
      
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
      // 使用 RESTful 风格的路径
      const response = await get(`/api/v1/tracking/${orderId}`)
      
      if (response && response.trackings) {
        const tracking = {
          orderNo: currentOrder.value?.orderNo || '',
          status: currentOrder.value?.status || '',
          estimatedDeliveryTime: '',
          logs: response.trackings.map(t => ({
            time: t.createdAt || t.created_at,
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
