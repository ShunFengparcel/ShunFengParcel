/**
 * 订单相关 API
 */
import { get, post } from '@/utils/request'

/**
 * 创建订单
 * @param {object} data - 订单数据
 */
export function createOrder(data) {
  return post('/api/v1/orders', data)
}

/**
 * 获取订单列表
 * @param {object} params - 查询参数
 */
export function getOrderList(params) {
  return get('/api/v1/orders', params)
}

/**
 * 获取订单详情
 * @param {number} id - 订单ID
 * @param {string} orderNo - 订单号
 */
export function getOrderDetail(id, orderNo) {
  const params = {}
  if (id) params.id = id
  if (orderNo) params.order_no = orderNo
  return get('/api/v1/orders/detail', params)
}

/**
 * 获取物流轨迹
 * @param {number} orderId - 订单ID
 */
export function getTracking(orderId) {
  return get('/api/v1/tracking', { order_id: orderId })
}

/**
 * 取消订单
 * @param {number} orderId - 订单ID
 */
export function cancelOrder(orderId) {
  return post(`/api/v1/orders/${orderId}/cancel`)
}

/**
 * 计算运费
 * @param {object} data - 计算参数
 */
export function calculateFee(data) {
  return post('/api/v1/orders/calculate-fee', data)
}
