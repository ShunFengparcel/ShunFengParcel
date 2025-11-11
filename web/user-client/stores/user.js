import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref('')
  const userInfo = ref(null)
  const addresses = ref([])
  const memberInfo = ref(null)

  // 计算属性
  const isLogin = computed(() => !!token.value)
  const phone = computed(() => userInfo.value?.phone || '')
  const memberLevel = computed(() => memberInfo.value?.level || 'normal')

  // 方法
  function setToken(newToken) {
    token.value = newToken
    uni.setStorageSync('token', newToken)
  }

  function setUserInfo(info) {
    userInfo.value = info
    uni.setStorageSync('userInfo', info)
  }

  function setAddresses(list) {
    addresses.value = list
  }

  function setMemberInfo(info) {
    memberInfo.value = info
  }

  function login(phone, code) {
    // TODO: 调用登录 API
    return new Promise((resolve) => {
      // Mock 数据 - 使用 user_id = 9 来匹配测试数据
      const mockToken = 'mock_token_' + Date.now()
      const mockUserInfo = {
        id: 9,
        phone: phone,
        nickname: '顺丰用户',
        avatar: '/static/images/default-avatar.png',
        memberLevel: 'normal',
        points: 2,
        coupons: 0,
        balance: 0
      }
      
      setToken(mockToken)
      setUserInfo(mockUserInfo)
      
      resolve(mockUserInfo)
    })
  }

  function loginWithWechat(code, userInfo) {
    // 调用后端微信登录接口
    return new Promise((resolve, reject) => {
      uni.request({
        url: 'http://localhost:18000/user/wxlogin',
        method: 'POST',
        data: { 
          code: code,
          nickname: userInfo.nickName || '',
          avatar_url: userInfo.avatarUrl || ''
        },
        header: {
          'Content-Type': 'application/json'
        },
        success: (res) => {
          console.log('Backend response:', res)
          if (res.statusCode === 200 && res.data.code === 200) {
            const data = res.data
            setToken(data.token)
            setUserInfo({
              id: data.user_id,
              nickname: data.nickname,
              avatar: userInfo.avatarUrl || '/static/images/default-avatar.png',
              openid: data.openid,
              memberLevel: 'normal',
              points: 0,
              coupons: 0,
              balance: 0
            })
            resolve(data)
          } else {
            reject(res.data.message || '登录失败')
          }
        },
        fail: (err) => {
          console.error('Request failed:', err)
          reject('网络请求失败')
        }
      })
    })
  }

  function logout() {
    token.value = ''
    userInfo.value = null
    addresses.value = []
    memberInfo.value = null
    
    uni.removeStorageSync('token')
    uni.removeStorageSync('userInfo')
    
    uni.reLaunch({
      url: '/pages/index/index'
    })
  }

  function getUserInfo() {
    // TODO: 调用获取用户信息 API
    return Promise.resolve(userInfo.value)
  }

  function getAddresses() {
    // TODO: 调用获取地址列表 API
    return Promise.resolve(addresses.value)
  }

  function getMemberInfo() {
    // TODO: 调用获取会员信息 API
    return Promise.resolve(memberInfo.value)
  }

  // 初始化：从本地存储恢复数据
  function init() {
    const savedToken = uni.getStorageSync('token')
    const savedUserInfo = uni.getStorageSync('userInfo')
    
    if (savedToken) {
      token.value = savedToken
    }
    if (savedUserInfo) {
      userInfo.value = savedUserInfo
    }
  }

  return {
    // 状态
    token,
    userInfo,
    addresses,
    memberInfo,
    // 计算属性
    isLogin,
    phone,
    memberLevel,
    // 方法
    setToken,
    setUserInfo,
    setAddresses,
    setMemberInfo,
    login,
    loginWithWechat,
    logout,
    getUserInfo,
    getAddresses,
    getMemberInfo,
    init
  }
})
