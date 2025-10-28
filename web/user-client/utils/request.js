/**
 * API 请求封装
 * 统一处理请求和响应
 */

// API 基础配置
const BASE_URL = process.env.NODE_ENV === 'development' 
  ? 'http://localhost:18000' 
  : 'https://api.shunfeng.com'

const TIMEOUT = 10000

/**
 * 请求拦截器
 */
function requestInterceptor(config) {
  // 添加 token
  const token = uni.getStorageSync('token')
  if (token) {
    config.header = {
      ...config.header,
      'Authorization': `Bearer ${token}`
    }
  }
  
  // 添加通用 header
  config.header = {
    ...config.header,
    'Content-Type': 'application/json'
  }
  
  return config
}

/**
 * 响应拦截器
 */
function responseInterceptor(response) {
  const { data, statusCode } = response
  
  // HTTP 状态码检查
  if (statusCode !== 200) {
    handleError({ code: statusCode, message: '网络请求失败' })
    return Promise.reject(response)
  }
  
  // 业务状态码检查 (后端返回 code: 0 表示成功)
  if (data.code !== 0 && data.code !== 200) {
    handleError(data)
    return Promise.reject(data)
  }
  
  return data.data
}

/**
 * 错误处理
 */
function handleError(error) {
  const { code, message } = error
  
  let errorMessage = message || '请求失败'
  
  switch (code) {
    case 1002: // 未授权
      errorMessage = '请先登录'
      // 跳转到登录页
      uni.navigateTo({
        url: '/pages/login/login'
      })
      break
    case 1003: // 无权限
      errorMessage = '无权限访问'
      break
    case 1004: // 资源不存在
      errorMessage = '资源不存在'
      break
    case 3004: // 地址超出服务范围
      errorMessage = '该地址暂不支持服务'
      break
    case 5001: // 支付失败
      errorMessage = '支付失败，请重试'
      break
  }
  
  uni.showToast({
    title: errorMessage,
    icon: 'none',
    duration: 2000
  })
}

/**
 * 请求方法
 */
function request(options) {
  return new Promise((resolve, reject) => {
    // 请求拦截
    const config = requestInterceptor({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data || {},
      header: options.header || {},
      timeout: options.timeout || TIMEOUT
    })
    
    console.log('🌐 [Request] 完整URL:', config.url)
    console.log('🌐 [Request] 请求配置:', config)
    
    // 发起请求
    uni.request({
      ...config,
      success: (response) => {
        console.log('🌐 [Request] 响应状态码:', response.statusCode)
        console.log('🌐 [Request] 响应数据:', response.data)
        
        // 响应拦截
        responseInterceptor(response)
          .then(resolve)
          .catch(reject)
      },
      fail: (error) => {
        console.error('🌐 [Request] 请求失败:', error)
        handleError({ code: 0, message: '网络连接失败' })
        reject(error)
      }
    })
  })
}

/**
 * GET 请求
 */
export function get(url, params = {}, options = {}) {
  console.log('🌐 [Request] GET请求:', url)
  console.log('🌐 [Request] 参数:', params)
  return request({
    url,
    method: 'GET',
    data: params,
    ...options
  })
}

/**
 * POST 请求
 */
export function post(url, data = {}, options = {}) {
  return request({
    url,
    method: 'POST',
    data,
    ...options
  })
}

/**
 * PUT 请求
 */
export function put(url, data = {}, options = {}) {
  return request({
    url,
    method: 'PUT',
    data,
    ...options
  })
}

/**
 * DELETE 请求
 */
export function del(url, data = {}, options = {}) {
  return request({
    url,
    method: 'DELETE',
    data,
    ...options
  })
}

export default {
  get,
  post,
  put,
  del
}
