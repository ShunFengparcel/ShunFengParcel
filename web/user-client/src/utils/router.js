/**
 * 路由拦截器
 * 处理页面跳转前的逻辑
 */

import { useUserStore } from '@/stores/user'

// 需要登录的页面列表
const authPages = [
  '/pages/send/create',
  '/pages/send/address',
  '/pages/order/detail',
  '/pages/profile/member',
  '/pages/profile/points',
  '/pages/profile/wallet'
]

/**
 * 检查是否需要登录
 */
function needAuth(url) {
  return authPages.some(page => url.includes(page))
}

/**
 * 跳转到登录页
 */
function toLogin(redirectUrl) {
  uni.navigateTo({
    url: `/pages/login/login?redirect=${encodeURIComponent(redirectUrl)}`
  })
}

/**
 * 路由跳转拦截
 */
export function navigateTo(options) {
  const { url } = options
  
  // 检查是否需要登录
  if (needAuth(url)) {
    const userStore = useUserStore()
    if (!userStore.isLogin) {
      toLogin(url)
      return
    }
  }
  
  // 执行跳转
  uni.navigateTo(options)
}

/**
 * 重定向拦截
 */
export function redirectTo(options) {
  const { url } = options
  
  // 检查是否需要登录
  if (needAuth(url)) {
    const userStore = useUserStore()
    if (!userStore.isLogin) {
      toLogin(url)
      return
    }
  }
  
  // 执行重定向
  uni.redirectTo(options)
}

/**
 * 切换 Tab 拦截
 */
export function switchTab(options) {
  // Tab 页面通常不需要登录验证
  uni.switchTab(options)
}

/**
 * 返回上一页
 */
export function navigateBack(options = {}) {
  uni.navigateBack(options)
}

/**
 * 重新加载
 */
export function reLaunch(options) {
  uni.reLaunch(options)
}

export default {
  navigateTo,
  redirectTo,
  switchTab,
  navigateBack,
  reLaunch
}
