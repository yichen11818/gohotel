<template>
  <view class="container">
    <!-- 顶部装饰背景 -->
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

    <!-- 主内容区 -->
    <view class="main-content">
      <!-- 标题区域 -->
      <view class="title-section">
        <text class="title">注册账号</text>
        <text class="subtitle">欢迎加入七天酒店</text>
      </view>

      <!-- 注册表单 -->
      <view class="register-form">
        <!-- 用户名输入 -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="my" color="#C29D71" size="40" />
          </view>
          <input
            v-model="formData.username"
            type="text"
            placeholder="请输入用户名 (3-20位)"
            maxlength="20"
            class="input-field"
          />
        </view>

        <!-- 邮箱输入 -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="email" color="#C29D71" size="40" />
          </view>
          <input
            v-model="formData.email"
            type="text"
            placeholder="请输入邮箱"
            class="input-field"
          />
        </view>

        <!-- 密码输入 -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="lock" color="#C29D71" size="40" />
          </view>
          <input
            v-model="formData.password"
            type="password"
            placeholder="请输入密码 (至少6位)"
            maxlength="20"
            class="input-field"
          />
        </view>

        <!-- 确认密码输入 -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="lock-fill" color="#C29D71" size="40" />
          </view>
          <input
            v-model="formData.confirmPassword"
            type="password"
            placeholder="请确认密码"
            maxlength="20"
            class="input-field"
          />
        </view>

        <!-- 手机号输入 (可选) -->
        <view class="input-group">
          <view class="input-icon">
            <TnIcon name="phone" color="#C29D71" size="40" />
          </view>
          <input
            v-model="formData.phone"
            type="number"
            placeholder="请输入手机号 (可选)"
            maxlength="11"
            class="input-field"
          />
        </view>

        <!-- 注册按钮 -->
        <view class="register-btn-wrapper">
          <TnButton
            shape="round"
            size="xl"
            width="100%"
            height="100rpx"
            bg-color="linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%)"
            text-color="#FFFFFF"
            :disabled="!canRegister"
            @click="handleRegister"
          >
            <text class="btn-text">立即注册</text>
          </TnButton>
        </view>

        <!-- 登录链接 -->
        <view class="login-link">
            <text>已有账号？</text>
            <text class="link" @click="goBack">立即登录</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { register } from '@/api/user.js'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'

const statusBarHeight = ref(44)
const formData = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  phone: ''
})

// 是否可以注册
const canRegister = computed(() => {
  return formData.username.length >= 3 &&
         formData.email &&
         formData.password.length >= 6 &&
         formData.confirmPassword === formData.password
})

// 注册
const handleRegister = async () => {
  if (!canRegister.value) return

  // 简单的邮箱验证
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(formData.email)) {
      uni.showToast({ title: '请输入有效的邮箱格式', icon: 'none' })
      return
  }

  try {
    uni.showLoading({ title: '注册中...' })

    await register({
      username: formData.username,
      email: formData.email,
      password: formData.password,
      phone: formData.phone
    })

    uni.hideLoading()
    uni.showToast({ title: '注册成功', icon: 'success' })

    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (error) {
    uni.hideLoading()
    // 错误处理通常在 request.js 中统一处理，或者在这里显示
    if (error.msg || error.message) {
         uni.showToast({ title: error.msg || error.message, icon: 'none' })
    }
  }
}

// 返回
const goBack = () => {
  uni.navigateBack()
}

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
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

.title-section {
  margin-bottom: 60rpx;

  .title {
    display: block;
    font-size: 48rpx;
    font-weight: 700;
    color: #333;
    margin-bottom: 16rpx;
  }

  .subtitle {
    font-size: 28rpx;
    color: #999;
  }
}

.register-form {
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
}

.register-btn-wrapper {
  margin-top: 60rpx;
  box-shadow: 0 16rpx 32rpx rgba(194, 157, 113, 0.35);
  border-radius: 100rpx;

  .btn-text {
    font-size: 36rpx;
    font-weight: 700;
    letter-spacing: 4rpx;
  }
}

.login-link {
    margin-top: 30rpx;
    text-align: center;
    font-size: 28rpx;
    color: #666;

    .link {
        color: #C29D71;
        margin-left: 10rpx;
        font-weight: 500;
    }
}
</style>
