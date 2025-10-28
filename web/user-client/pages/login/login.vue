<template>
  <view class="container">
    <view class="logo-section">
      <view class="logo">
        <text class="logo-text">SF</text>
      </view>
      <text class="app-name">顺丰速运</text>
    </view>

    <view class="login-form">
      <view class="agreement">
        <checkbox-group @change="handleAgreementChange">
          <label class="agreement-label">
            <checkbox :checked="agreed" color="#D81E06" />
            <text class="agreement-text">
              阅读并同意《顺丰速运用户协议》和《顺丰速运隐私政策》
            </text>
          </label>
        </checkbox-group>
      </view>

      <button 
        class="login-btn" 
        :class="{ disabled: !agreed }"
        @click="handleWechatLogin"
        open-type="getUserInfo"
        @getuserinfo="onGetUserInfo"
      >
        <text class="login-btn-text">微信一键登录</text>
      </button>
    </view>

    <view class="customer-service">
      <text class="service-text">遇到问题？拨打客服电话 </text>
      <text class="service-phone">95338</text>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const agreed = ref(false)

const handleAgreementChange = (e) => {
  agreed.value = e.detail.value.length > 0
}

const handleWechatLogin = () => {
  if (!agreed.value) {
    uni.showToast({ title: '请先阅读并同意用户协议', icon: 'none' })
    return
  }

  uni.showLoading({ title: '登录中...' })
  
  // 调用微信登录
  uni.login({
    provider: 'weixin',
    success: (loginRes) => {
      console.log('wx.login success:', loginRes)
      const code = loginRes.code
      
      // 获取用户信息
      uni.getUserInfo({
        provider: 'weixin',
        success: (infoRes) => {
          console.log('getUserInfo success:', infoRes)
          
          // 调用后端登录接口
          userStore.loginWithWechat(code, infoRes.userInfo).then(() => {
            uni.hideLoading()
            uni.showToast({ title: '登录成功', icon: 'success' })
            setTimeout(() => {
              uni.switchTab({ url: '/pages/index/index' })
            }, 1500)
          }).catch((err) => {
            uni.hideLoading()
            uni.showToast({ title: '登录失败', icon: 'none' })
            console.error('Login failed:', err)
          })
        },
        fail: (err) => {
          uni.hideLoading()
          uni.showToast({ title: '获取用户信息失败', icon: 'none' })
          console.error('getUserInfo failed:', err)
        }
      })
    },
    fail: (err) => {
      uni.hideLoading()
      uni.showToast({ title: '微信登录失败', icon: 'none' })
      console.error('wx.login failed:', err)
    }
  })
}

const onGetUserInfo = (e) => {
  console.log('onGetUserInfo:', e)
  if (e.detail.userInfo) {
    handleWechatLogin()
  }
}
</script>

<style scoped>
.container {
  min-height: 100vh;
  background: linear-gradient(180deg, #F5F5F5 0%, #FFFFFF 50%);
  padding: 0 60rpx;
  position: relative;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 200rpx;
  margin-bottom: 120rpx;
}

.logo {
  width: 180rpx;
  height: 180rpx;
  background-color: #000000;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 8rpx solid #FFFFFF;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.1);
}

.logo-text {
  font-size: 72rpx;
  font-weight: bold;
  color: #FFFFFF;
}

.app-name {
  font-size: 36rpx;
  font-weight: 600;
  color: #333333;
  margin-top: 30rpx;
}

.login-form {
  margin-bottom: 80rpx;
}

.agreement {
  margin-bottom: 60rpx;
  padding: 0 10rpx;
}

.agreement-label {
  display: flex;
  align-items: flex-start;
}

.agreement-text {
  font-size: 24rpx;
  color: #666666;
  line-height: 36rpx;
  margin-left: 10rpx;
  flex: 1;
}

.agreement-tip {
  font-size: 24rpx;
  color: #666666;
  line-height: 36rpx;
  margin-left: 50rpx;
}

.login-btn {
  height: 100rpx;
  background: linear-gradient(135deg, #07C160 0%, #38D97F 100%);
  border-radius: 50rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 20rpx rgba(7, 193, 96, 0.3);
  border: none;
  line-height: 100rpx;
}

.login-btn.disabled {
  background: #E0E0E0;
  box-shadow: none;
}

.login-btn-text {
  font-size: 32rpx;
  color: #FFFFFF;
  font-weight: 600;
}

.customer-service {
  position: absolute;
  bottom: 80rpx;
  left: 0;
  right: 0;
  text-align: center;
}

.service-text {
  font-size: 24rpx;
  color: #999999;
}

.service-phone {
  font-size: 24rpx;
  color: #333333;
  font-weight: 600;
}
</style>
