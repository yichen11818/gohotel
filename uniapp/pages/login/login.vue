<template>
  <view class="container">
    <!-- 顶部装饰背景 -->
    <view class="header-bg">
      <view class="gradient-layer"></view>
    </view>

    <!-- 返回按钮 -->
    <view class="back-btn" @click="goBack">
      <TnIcon name="left" color="#333" size="40" />
    </view>

    <!-- 主内容区 -->
    <view class="main-content">
      <!-- Logo区域 -->
      <view class="logo-section">
        <view class="logo-wrapper">
          <image class="logo" src="https://images.unsplash.com/photo-1566073771259-6a8506099945?w=200&h=200&fit=crop" mode="aspectFill"></image>
        </view>
        <text class="brand-name">七天酒店</text>
        <text class="brand-slogan">品质生活，从这里开始</text>
      </view>

      <!-- 登录表单 -->
      <view class="login-form">
        <!-- 手机号输入 -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="phone" color="#C29D71" size="40" />
          </view>
          <input 
            v-model="phone" 
            type="number" 
            placeholder="请输入手机号" 
            maxlength="11"
            class="input-field"
          />
          <view v-if="phone" class="clear-btn" @click="phone = ''">
            <TnIcon name="close-circle-fill" color="#ccc" size="32" />
          </view>
        </view>

        <!-- 验证码输入 -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="shield" color="#C29D71" size="40" />
          </view>
          <input 
            v-model="code" 
            type="number" 
            placeholder="请输入验证码" 
            maxlength="6"
            class="input-field"
          />
          <view 
            class="code-btn" 
            :class="{ disabled: countdown > 0 }"
            @click="sendCode"
          >
            <text>{{ countdown > 0 ? `${countdown}s` : '获取验证码' }}</text>
          </view>
        </view>

        <!-- 登录按钮 -->
        <view class="login-btn-wrapper">
          <TnButton 
            shape="round" 
            size="xl" 
            width="100%" 
            height="100rpx"
            bg-color="linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%)"
            text-color="#FFFFFF"
            :disabled="!canLogin"
            @click="handleLogin"
          >
            <text class="btn-text">登录</text>
          </TnButton>
        </view>

        <!-- 协议 -->
        <view class="agreement">
          <view class="checkbox" :class="{ checked: agreed }" @click="agreed = !agreed">
            <TnIcon v-if="agreed" name="check" color="#fff" size="20" />
          </view>
          <text class="agreement-text">
            登录即表示同意
            <text class="link" @click.stop="goToUserAgreement">《用户协议》</text>
            和
            <text class="link" @click.stop="goToPrivacy">《隐私政策》</text>
          </text>
        </view>
        
        <!-- 注册链接 -->
        <view class="register-link">
            <text>还没有账号？</text>
            <text class="link" @click="goToRegister">立即注册</text>
        </view>
      </view>

      <!-- 其他登录方式 -->
      <view class="other-login">
        <view class="divider">
          <view class="line"></view>
          <text class="text">其他登录方式</text>
          <view class="line"></view>
        </view>

        <view class="social-login">
          <view class="social-item wechat" @click="wechatLogin">
            <TnIcon name="logo-wechat" color="#07C160" size="56" />
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { login, sendVerifyCode } from '@/api/user.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'

const phone = ref('')
const code = ref('')
const countdown = ref(0)
const agreed = ref(false)

// 是否可以登录
const canLogin = computed(() => {
  return phone.value.length === 11 && code.value.length >= 4 && agreed.value
})

// 发送验证码
const sendCode = async () => {
  if (countdown.value > 0) return
  
  if (phone.value.length !== 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  
  try {
    uni.showLoading({ title: '发送中...' })
    await sendVerifyCode(phone.value, 'login')
    uni.hideLoading()
    uni.showToast({ title: '验证码已发送', icon: 'success' })
    countdown.value = 60
    
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    uni.hideLoading()
  }
}

// 登录
const handleLogin = async () => {
  if (!canLogin.value) return
  
  try {
    uni.showLoading({ title: '登录中...' })
    
    const result = await login({
      phone: phone.value,
      code: code.value
    })
    
    uni.hideLoading()
    uni.showToast({ title: '登录成功', icon: 'success' })
    
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (error) {
    uni.hideLoading()
  }
}

// 微信登录
const wechatLogin = () => {
  if (!agreed.value) {
    uni.showToast({ title: '请先同意用户协议', icon: 'none' })
    return
  }
  
  // #ifdef MP-WEIXIN
  uni.login({
    provider: 'weixin',
    success: (loginRes) => {
      uni.getUserProfile({
        desc: '用于完善用户资料',
        success: (res) => {
          // 调用后端接口进行登录
          handleRealWechatLogin(loginRes.code, res.userInfo)
        }
      })
    }
  })
  // #endif
}

const handleRealWechatLogin = async (code, userInfo) => {
  try {
    uni.showLoading({ title: '登录中...' })
    const result = await post('/auth/wechat-login', {
      code,
      nickname: userInfo.nickName,
      avatar: userInfo.avatarUrl
    })
    
    if (result.token) {
      uni.setStorageSync(TOKEN_KEY, result.token)
    }
    if (result.userInfo) {
      uni.setStorageSync(USER_INFO_KEY, result.userInfo)
    }
    
    uni.hideLoading()
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1500)
  } catch (error) {
    uni.hideLoading()
  }
}

// 返回
const goBack = () => {
  uni.navigateBack()
}

// 用户协议
const goToUserAgreement = () => {
  uni.navigateTo({
    url: '/pages/agreement/agreement'
  })
}

// 隐私政策
const goToPrivacy = () => {
  uni.navigateTo({
    url: '/pages/privacy/privacy'
  })
}

onLoad(() => {
  // 检查是否已登录
  const token = uni.getStorageSync(TOKEN_KEY)
  if (token) {
    uni.navigateBack()
  }
})
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #FAFAFA;
  position: relative;
}

.header-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 400rpx;
  overflow: hidden;
  z-index: 0;
  pointer-events: none;
  
  .gradient-layer {
    width: 100%;
    height: 100%;
    background: linear-gradient(180deg, rgba(194, 157, 113, 0.15) 0%, transparent 100%);
  }
}

.back-btn {
  position: absolute;
  top: 88rpx;
  left: 32rpx;
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
  z-index: 100;
  
  &:active {
    transform: scale(0.9);
  }
}

.main-content {
  position: relative;
  z-index: 10;
  padding: 0 48rpx;
  padding-top: 180rpx;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 80rpx;
  
  .logo-wrapper {
    width: 160rpx;
    height: 160rpx;
    border-radius: 32rpx;
    overflow: hidden;
    box-shadow: 0 16rpx 48rpx rgba(194, 157, 113, 0.3);
    margin-bottom: 32rpx;
    
    .logo {
      width: 100%;
      height: 100%;
    }
  }
  
  .brand-name {
    font-size: 48rpx;
    font-weight: 700;
    color: #333;
    margin-bottom: 12rpx;
    letter-spacing: 2rpx;
  }
  
  .brand-slogan {
    font-size: 26rpx;
    color: #999;
  }
}

.login-form {
  margin-bottom: 60rpx;
}

.input-group {
  display: flex;
  align-items: center;
  height: 110rpx;
  background: #fff;
  border-radius: 20rpx;
  padding: 0 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .input-icon {
    margin-right: 20rpx;
  }
  
  .input-field {
    flex: 1;
    height: 100%;
    font-size: 30rpx;
    color: #333;
  }
  
  .clear-btn {
    padding: 10rpx;
    
    &:active {
      opacity: 0.6;
    }
  }
  
  .code-btn {
    padding: 16rpx 24rpx;
    background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
    border-radius: 100rpx;
    
    text {
      font-size: 26rpx;
      color: #C29D71;
      font-weight: 500;
    }
    
    &.disabled {
      opacity: 0.6;
    }
    
    &:active:not(.disabled) {
      transform: scale(0.95);
    }
  }
}

.login-btn-wrapper {
  margin-top: 48rpx;
  box-shadow: 0 16rpx 32rpx rgba(194, 157, 113, 0.35);
  border-radius: 100rpx;
  
  .btn-text {
    font-size: 36rpx;
    font-weight: 700;
    letter-spacing: 4rpx;
  }
}

.agreement {
  display: flex;
  align-items: flex-start;
  margin-top: 32rpx;
  padding: 0 8rpx;
  
  .checkbox {
    width: 36rpx;
    height: 36rpx;
    border: 2rpx solid #ddd;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-right: 12rpx;
    margin-top: 4rpx;
    transition: all 0.3s ease;
    
    &.checked {
      background: linear-gradient(135deg, #D4B184 0%, #C29D71 100%);
      border-color: #C29D71;
    }
  }
  
  .agreement-text {
    flex: 1;
    font-size: 24rpx;
    color: #999;
    line-height: 1.6;
    
    .link {
      color: #C29D71;
    }
  }
}

.register-link {
    margin-top: 20rpx;
    text-align: center;
    font-size: 26rpx;
    color: #999;
    
    .link {
        color: #C29D71;
        margin-left: 10rpx;
        font-weight: 500;
    }
}

.other-login {
  .divider {
    display: flex;
    align-items: center;
    margin-bottom: 48rpx;
    
    .line {
      flex: 1;
      height: 1rpx;
      background: linear-gradient(to right, transparent, #E5E5E5, transparent);
    }
    
    .text {
      padding: 0 24rpx;
      font-size: 24rpx;
      color: #999;
    }
  }
  
  .social-login {
    display: flex;
    justify-content: center;
    gap: 48rpx;
  }
  
  .social-item {
    width: 100rpx;
    height: 100rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fff;
    border-radius: 50%;
    box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.08);
    transition: all 0.3s ease;
    
    &:active {
      transform: scale(0.9);
      box-shadow: 0 4rpx 12rpx rgba(0, 0, 0, 0.08);
    }
  }
}
</style>

