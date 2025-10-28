import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'

export function createApp() {
  const app = createSSRApp(App)
  const pinia = createPinia()
  
  app.use(pinia)
  
  // 延迟初始化用户状态，避免在 App 实例创建前访问
  app.onLaunch = function() {
    try {
      const { useUserStore } = require('./stores/user')
      const userStore = useUserStore()
      userStore.init()
    } catch (error) {
      console.error('初始化用户状态失败:', error)
    }
  }
  
  return {
    app,
    pinia
  }
}
