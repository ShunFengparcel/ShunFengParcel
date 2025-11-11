/**
 * API 请求封装
 * 统一处理请求和响应
 */

// API 基础配置
// 开发环境使用本机IP（微信小程序无法直接访问localhost）
// 生产环境使用域名
// 注意：微信小程序只能访问https或指定的域名，本地开发时需要在微信开发者工具中设置不校验合法域名
const BASE_URL = process.env.NODE_ENV === 'development'
  ? 'http://localhost:18000'  // 开发环境使用localhost
  : 'https://jinhuijuan.fun'  // 生产环境使用你的域名

// 确保API路径前缀正确
// 根据后端实际API路径结构设置
const API_PREFIX = '/api/v1'

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
  console.log('🌐 [ResponseInterceptor] ========== 开始处理响应 ==========')
  
  const { data, statusCode } = response

  console.log('🌐 [ResponseInterceptor] statusCode:', statusCode)
  console.log('🌐 [ResponseInterceptor] data类型:', typeof data)
  console.log('🌐 [ResponseInterceptor] data是否为null:', data === null)
  console.log('🌐 [ResponseInterceptor] data是否为undefined:', data === undefined)
  console.log('🌐 [ResponseInterceptor] data:', data)
  
  // 如果data不存在，直接返回错误
  if (!data) {
    console.error('🌐 [ResponseInterceptor] ❌ 响应数据为空！')
    handleError({ code: 500, message: '响应数据为空' })
    return Promise.reject({ code: 500, message: '响应数据为空' })
  }
  
  console.log('🌐 [ResponseInterceptor] data.code:', data.code)
  console.log('🌐 [ResponseInterceptor] data.message:', data.message)
  console.log('🌐 [ResponseInterceptor] data.data:', data.data)
  
  // 尝试序列化原始数据
  try {
    const rawDataStr = JSON.stringify(data, null, 2)
    console.log('🌐 [ResponseInterceptor] 完整响应JSON:')
    console.log(rawDataStr)
  } catch (e) {
    console.error('🌐 [ResponseInterceptor] 无法序列化原始data:', e)
  }

  // HTTP 状态码检查
  if (statusCode !== 200) {
    console.error('🌐 [ResponseInterceptor] ❌ HTTP状态码错误:', statusCode)
    handleError({ code: statusCode, message: '网络请求失败' })
    return Promise.reject(response)
  }

  // 业务状态码检查 (后端返回 code: 0 或 "0" 表示成功)
  const code = parseInt(data.code)
  console.log('🌐 [ResponseInterceptor] 业务状态码:', code)
  
  if (code !== 0 && code !== 200) {
    console.error('🌐 [ResponseInterceptor] ❌ 业务状态码错误:', code)
    handleError(data)
    return Promise.reject(data)
  }

  // 检查data.data是否存在
  if (!data.data) {
    console.warn('🌐 [ResponseInterceptor] ⚠️ data.data为空')
    return null
  }
  
  // 如果是订单列表，打印详细信息
  if (data.data.orders) {
    console.log('🌐 [ResponseInterceptor] ✅ 这是订单列表响应')
    console.log('🌐 [ResponseInterceptor] orders数组长度:', data.data.orders.length)
    console.log('🌐 [ResponseInterceptor] total:', data.data.total)
    if (data.data.orders.length > 0) {
      console.log('🌐 [ResponseInterceptor] 第一条订单:', data.data.orders[0])
    }
  }
  
  console.log('🌐 [ResponseInterceptor] ========== 返回 data.data ==========')
  return data.data
}

/**
 * 将蛇形命名转换为驼峰命名
 */
function convertSnakeToCamel(obj) {
  if (obj === null || obj === undefined) {
    return obj
  }

  // 先序列化再反序列化，确保返回纯对象（避免Proxy问题）
  const plainObj = JSON.parse(JSON.stringify(obj))

  if (Array.isArray(plainObj)) {
    return plainObj.map(item => convertSnakeToCamel(item))
  }

  if (typeof plainObj === 'object') {
    const newObj = {}
    for (const key in plainObj) {
      if (plainObj.hasOwnProperty(key)) {
        // 将蛇形命名转换为驼峰命名
        const camelKey = key.replace(/_([a-z])/g, (match, letter) => letter.toUpperCase())
        newObj[camelKey] = convertSnakeToCamel(plainObj[key])
      }
    }
    return newObj
  }

  return plainObj
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
  console.log('🌐 [Request] ========== 开始请求 ==========');
  console.log('🌐 [Request] 原始 options:', options);

  return new Promise((resolve, reject) => {
    // 构建完整URL - 确保API前缀正确
    let finalUrl = options.url || '';
    
    // 如果URL不是以http开头，则添加BASE_URL和API_PREFIX
    if (!finalUrl.startsWith('http')) {
      // 确保URL路径正确
      if (finalUrl.startsWith('/')) {
        finalUrl = finalUrl.substring(1); // 移除开头的斜杠
      }
      
      // 确保API_PREFIX没有多余的斜杠
      let prefix = API_PREFIX;
      if (prefix.endsWith('/')) {
        prefix = prefix.slice(0, -1);
      }
      
      // 构建完整URL
      finalUrl = `${BASE_URL}${prefix}/${finalUrl}`;
      
      // 打印完整URL用于调试
      console.log('🌐 [Request] 完整URL构建过程:');
      console.log('🌐 [Request] BASE_URL:', BASE_URL);
      console.log('🌐 [Request] API_PREFIX:', API_PREFIX);
      console.log('🌐 [Request] 路径部分:', finalUrl);
    }
    
    console.log('🌐 [Request] 最终URL:', finalUrl);
    console.log('🌐 [Request] 请求方法:', options.method || 'GET');
    console.log('🌐 [Request] 请求参数:', options.data || '无参数');
    
    // 请求拦截
    const config = requestInterceptor({
      url: finalUrl,
      method: options.method || 'GET',
      data: options.data || {},
      header: options.header || {},
      timeout: options.timeout || TIMEOUT
    })

    console.log('🌐 [Request] 完整URL:', config.url)
    console.log('🌐 [Request] 请求方法:', config.method)
    console.log('🌐 [Request] 请求数据:', config.data)
    console.log('🌐 [Request] 请求头:', config.header)

    // 发起请求
    uni.request({
      ...config,
      success: (response) => {
        console.log('🌐 [Request] ========== 收到响应 ==========')
        console.log('🌐 [Request] 响应状态码:', response.statusCode)
        console.log('🌐 [Request] 响应数据:', response.data)
        console.log('🌐 [Request] 响应头:', response.header)

        // 响应拦截
        responseInterceptor(response)
          .then((data) => {
            console.log('🌐 [Request] 拦截器处理后的数据:', data)
            resolve(data)
          })
          .catch((err) => {
            console.error('🌐 [Request] 拦截器拒绝:', err)
            reject(err)
          })
      },
      fail: (error) => {
        console.error('🌐 [Request] ========== 请求失败 ==========')
        console.error('🌐 [Request] 错误信息:', error)
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
  console.log('🌐🌐🌐 [GET] ========== GET函数被调用 ==========')
  console.log('🌐 [GET] URL:', url)
  console.log('🌐 [GET] 参数:', params)
  console.log('🌐 [GET] 选项:', options)

  // 将参数拼接到URL中
  let finalUrl = url;
  if (params && Object.keys(params).length > 0) {
    const queryString = Object.keys(params)
      .filter(key => params[key] !== undefined && params[key] !== null)
      .map(key => `${encodeURIComponent(key)}=${encodeURIComponent(params[key])}`)
      .join('&');
    
    if (queryString) {
      finalUrl += (url.includes('?') ? '&' : '?') + queryString;
    }
  }

  console.log('🌐 [GET] 最终URL:', finalUrl)

  return request({
    url: finalUrl,
    method: 'GET',
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
