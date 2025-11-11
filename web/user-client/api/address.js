/**
 * 地址相关 API
 */
import { get, post, put, del } from '@/utils/request'

/**
 * 获取地址列表
 * @param {number} userId - 用户ID
 */
export function getAddressList(userId) {
  return get('/api/v1/list/addresses', { user_id: userId })
}

/**
 * 创建地址
 * @param {object} data - 地址数据
 */
export function createAddress(data) {
  return post('/api/v1/create/addresses', data)
}

/**
 * 更新地址
 * @param {number} id - 地址ID
 * @param {object} data - 地址数据
 */
export function updateAddress(id, data) {
  return put(`/api/v1/addresses/${id}`, data)
}

/**
 * 删除地址
 * @param {number} id - 地址ID
 */
export function deleteAddress(id) {
  return del(`/api/v1/addresses/${id}`)
}

/**
 * 设置默认地址
 * @param {number} userId - 用户ID
 * @param {number} addressId - 地址ID
 */
export function setDefaultAddress(userId, addressId) {
  return post('/api/v1/addresses/default', { user_id: userId, address_id: addressId })
}
