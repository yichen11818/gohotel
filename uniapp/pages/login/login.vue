<template>
  <view class="container">
    <view class="header-bg">
      <view class="gradient-layer"></view>
    </view>

    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
    </view>

    <view class="main-content">
      <view class="logo-section">
        <view class="logo-wrapper">
          <image class="logo" src="https://images.unsplash.com/photo-1566073771259-6a8506099945?w=200&h=200&fit=crop" mode="aspectFill"></image>
        </view>
        <text class="brand-name">七天酒店</text>
        <text class="brand-slogan">诚挚待客 · 舒适如家</text>
      </view>

      <view class="login-form">
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="people" color="#C29D71" size="40" />
          </view>
          <input
            v-model="username"
            type="text"
            placeholder="请输入用户名"
            class="input-field"
          />
          <view v-if="username" class="clear-btn" @click="username = ''">
            <TnIcon name="close-circle-fill" color="#ccc" size="32" />
          </view>
        </view>

        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="lock" color="#C29D71" size="40" />
          </view>
          <input
            v-model="password"
            password
            placeholder="请输入密码"
            class="input-field"
          />
        </view>

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

        <view class="register-link">
          <text>还没有账号？</text>
          <text class="link" @click="goToRegister">立即注册</text>
        </view>
      </view>

      <view class="other-login">
        <view class="divider">
          <view class="line"></view>
          <text class="text">其他登录方式</text>
          <view class="line"></view>
        </view>

        <view class="social-login">
          <view class="social-item wechat" @click="handleWechatLogin">
            <TnIcon name="logo-wechat" color="#07C160" size="56" />
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { login, wechatLogin } from '@/api/user.js'
import { TOKEN_KEY } from '@/config/api.config.js'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'

const username = ref('')
const password = ref('')
const agreed = ref(false)
const statusBarHeight = ref(44)

const canLogin = computed(() => {
  return username.value.trim().length >= 3 && password.value.length >= 6 && agreed.value
})

const handleLogin = async () => {
  if (!canLogin.value) return

  try {
    uni.showLoading({ title: '登录中...' })

    await login({
      username: username.value.trim(),
      password: password.value,
    })

    uni.hideLoading()
    uni.showToast({ title: '登录成功', icon: 'success' })

    setTimeout(() => {
      uni.navigateBack()
    }, 1200)
  } catch (error) {
    uni.hideLoading()
  }
}

const handleWechatLogin = () => {
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
        success: async (res) => {
          try {
            uni.showLoading({ title: '登录中...' })
            await wechatLogin({
              code: loginRes.code,
              nickname: res.userInfo.nickName,
              avatar: res.userInfo.avatarUrl,
            })
            uni.hideLoading()
            uni.showToast({ title: '登录成功', icon: 'success' })
            setTimeout(() => uni.navigateBack(), 1200)
          } catch (error) {
            uni.hideLoading()
          }
        }
      })
    }
  })
  // #endif

  // #ifndef MP-WEIXIN
  uni.showToast({ title: '请在微信小程序环境中使用微信登录', icon: 'none' })
  // #endif
}

const goBack = () => {
  uni.navigateBack()
}

const goToRegister = () => {
  uni.navigateTo({
    url: '/pages/register/register'
  })
}

const goToUserAgreement = () => {
  uni.navigateTo({
    url: '/pages/agreement/agreement'
  })
}

const goToPrivacy = () => {
  uni.navigateTo({
    url: '/pages/privacy/privacy'
  })
}

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
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

.nav-header {
  height: 88rpx;
  padding: 0 30rpx;
  position: relative;
  z-index: 100;

  .back-btn {
    width: 64rpx;
    height: 64rpx;
    background: rgba(255, 255, 255, 0.9);
    border-radius: 50%;
    box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
    .icon-back { width: 32rpx; height: 32rpx; }

    &:active {
      transform: scale(0.9);
    }
  }
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

.main-content {
  position: relative;
  z-index: 10;
  padding: 0 48rpx;
  padding-top: 100rpx;
}

.logo-section {
  text-align: center;
  margin-bottom: 80rpx;
}

.logo-wrapper {
  width: 160rpx;
  height: 160rpx;
  margin: 0 auto 24rpx;
  border-radius: 32rpx;
  overflow: hidden;
  box-shadow: 0 8rpx 32rpx rgba(194, 157, 113, 0.2);
}

.logo {
  width: 100%;
  height: 100%;
}

.brand-name {
  display: block;
  font-size: 52rpx;
  font-weight: 700;
  color: #333;
  margin-bottom: 12rpx;
}

.brand-slogan {
  font-size: 28rpx;
  color: #999;
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
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.04);
}

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
  padding: 8rpx;
}

.login-btn-wrapper {
  margin: 48rpx 0 32rpx;
}

.btn-text {
  font-size: 32rpx;
  font-weight: 600;
}

.agreement {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.checkbox {
  width: 32rpx;
  height: 32rpx;
  border-radius: 50%;
  border: 2rpx solid #d9d9d9;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-top: 6rpx;

  &.checked {
    background: #C29D71;
    border-color: #C29D71;
  }
}

.agreement-text {
  flex: 1;
  font-size: 24rpx;
  line-height: 1.8;
  color: #999;
}

.link {
  color: #C29D71;
}

.register-link {
  text-align: center;
  font-size: 26rpx;
  color: #999;
}

.other-login {
  margin-top: 80rpx;
}

.divider {
  display: flex;
  align-items: center;
  gap: 24rpx;
  margin-bottom: 32rpx;

  .line {
    flex: 1;
    height: 1rpx;
    background: #e8e8e8;
  }

  .text {
    font-size: 24rpx;
    color: #999;
  }
}

.social-login {
  display: flex;
  justify-content: center;
}

.social-item {
  width: 96rpx;
  height: 96rpx;
  border-radius: 50%;
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.06);
}
</style>
