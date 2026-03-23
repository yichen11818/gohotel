<template>
  <view class="page">
    <view class="content">
      <view class="section">
        <view class="menu-item" @click="goToProfile">
          <text>个人资料</text>
          <text class="menu-desc">编辑姓名、手机号与头像</text>
        </view>
        <view class="menu-item" @click="goToSecurity">
          <text>修改密码</text>
          <text class="menu-desc">保障账号安全</text>
        </view>
      </view>

      <view class="section">
        <view class="menu-item" @click="goToAbout">
          <text>关于我们</text>
          <text class="menu-desc">查看酒店介绍与联系方式</text>
        </view>
        <view class="menu-item" @click="goToPrivacy">
          <text>隐私政策</text>
          <text class="menu-desc">查看平台隐私说明</text>
        </view>
        <view class="menu-item" @click="goToAgreement">
          <text>用户协议</text>
          <text class="menu-desc">查看注册与入住相关协议</text>
        </view>
      </view>

      <view class="section">
        <view class="menu-item" @click="clearCache">
          <text>清除缓存</text>
          <text class="menu-desc">当前缓存 {{ cacheSize }}</text>
        </view>
      </view>

      <button class="logout-btn" @click="handleLogout">退出登录</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'

const cacheSize = ref('0KB')

const requireLogin = (callback) => {
  if (!uni.getStorageSync(TOKEN_KEY)) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }
  callback()
}

const goToProfile = () => {
  requireLogin(() => {
    uni.navigateTo({ url: '/pages/profile/profile' })
  })
}

const goToSecurity = () => {
  requireLogin(() => {
    uni.navigateTo({ url: '/pages/security/security' })
  })
}

const goToAbout = () => {
  uni.navigateTo({ url: '/pages/about/about' })
}

const goToPrivacy = () => {
  uni.navigateTo({ url: '/pages/privacy/privacy' })
}

const goToAgreement = () => {
  uni.navigateTo({ url: '/pages/agreement/agreement' })
}

const clearCache = () => {
  uni.showModal({
    title: '清除缓存',
    content: '确定清除本地缓存吗？不会影响当前登录状态。',
    success: (res) => {
      if (!res.confirm) return
      const token = uni.getStorageSync(TOKEN_KEY)
      const userInfo = uni.getStorageSync(USER_INFO_KEY)
      uni.clearStorageSync()
      if (token) {
        uni.setStorageSync(TOKEN_KEY, token)
      }
      if (userInfo) {
        uni.setStorageSync(USER_INFO_KEY, userInfo)
      }
      cacheSize.value = '0KB'
      uni.showToast({ title: '缓存已清除', icon: 'success' })
    },
  })
}

const handleLogout = () => {
  uni.showModal({
    title: '退出登录',
    content: '确认退出当前账号吗？',
    success: async (res) => {
      if (!res.confirm) return
      await user.logout()
      uni.showToast({ title: '已退出登录', icon: 'success' })
      setTimeout(() => {
        uni.reLaunch({ url: '/pages/index/index' })
      }, 800)
    },
  })
}
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
}

.content {
  padding: 32rpx 24rpx;
}

.section {
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.menu-item {
  padding: 28rpx 24rpx;
  border-bottom: 1px solid #f3f4f6;
}

.menu-item:last-child {
  border-bottom: 0;
}

.menu-desc {
  display: block;
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #6b7280;
}

.logout-btn {
  margin-top: 36rpx;
  background: #fff;
  color: #dc2626;
  border-radius: 999rpx;
}
</style>
