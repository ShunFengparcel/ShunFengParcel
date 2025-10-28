import { createSSRApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import { useUserStore } from './stores/user'

export function createApp() {
  const app = createSSRApp(App)
  const pinia = createPinia()
  
  app.use(pinia)
  
  // 初始化用户状态
  const userStore = useUserStore()
  userStore.init()
  
  return {
    app,
    pinia
  }
}
