/**
 * Toast 提示封装
 */

export function showToast(options) {
  const defaultOptions = {
    icon: 'none',
    duration: 2000,
    mask: false
  }
  
  if (typeof options === 'string') {
    uni.showToast({
      ...defaultOptions,
      title: options
    })
  } else {
    uni.showToast({
      ...defaultOptions,
      ...options
    })
  }
}

export function showSuccess(title, duration = 2000) {
  uni.showToast({
    title,
    icon: 'success',
    duration,
    mask: false
  })
}

export function showError(title, duration = 2000) {
  uni.showToast({
    title,
    icon: 'error',
    duration,
    mask: false
  })
}

export function showLoading(title = '加载中...') {
  uni.showLoading({
    title,
    mask: true
  })
}

export function hideLoading() {
  uni.hideLoading()
}

export default {
  showToast,
  showSuccess,
  showError,
  showLoading,
  hideLoading
}
