<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">系统设置</text>
      <view class="placeholder-view"></view>
    </view>

    <scroll-view class="content" scroll-y>
      <view class="section premium-card">
        <view class="menu-item flex-between" @click="goToProfile">
          <view class="item-left">
            <text class="label">个人资料</text>
            <text class="menu-desc">编辑姓名、手机号与头像</text>
          </view>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
        <view class="menu-item flex-between" @click="goToSecurity">
          <view class="item-left">
            <text class="label">修改密码</text>
            <text class="menu-desc">保障您的账号安全</text>
          </view>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
      </view>

      <view class="section premium-card">
        <view class="menu-item flex-between" @click="goToAbout">
          <text class="label">关于我们</text>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
        <view class="menu-item flex-between" @click="goToPrivacy">
          <text class="label">隐私政策</text>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
        <view class="menu-item flex-between" @click="goToAgreement">
          <text class="label">用户协议</text>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
      </view>

      <view class="section premium-card">
        <view class="menu-item flex-between" @click="clearCache">
          <view class="item-left">
            <text class="label">清除缓存</text>
            <text class="menu-desc">释放手机存储空间</text>
          </view>
          <text class="val">{{ cacheSize }}</text>
        </view>
      </view>

      <view class="action-box">
        <button class="logout-btn flex-center" @click="handleLogout">
          <text>退出登录</text>
        </button>
      </view>

      <view class="version-info">
        <text>Version 1.2.0</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'

const statusBarHeight = ref(44)
const cacheSize = ref('12.4MB')

const goBack = () => {
  uni.navigateBack()
}

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
    confirmColor: '#C29D71',
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
    confirmColor: '#E64340',
    success: async (res) => {
      if (!res.confirm) return
      await user.logout()
      uni.showToast({ title: '已安全退出', icon: 'success' })
      setTimeout(() => {
        uni.reLaunch({ url: '/pages/index/index' })
      }, 800)
    },
  })
}

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: $bg-color;
  display: flex;
  flex-direction: column;
}

.nav-header {
  height: 88rpx;
  padding: 0 30rpx;
  background: #fff;

  .back-btn {
    width: 64rpx;
    height: 64rpx;
    .icon-back { width: 40rpx; height: 40rpx; }
  }

  .page-title {
    font-size: 34rpx;
    font-weight: 700;
    color: $text-main;
  }

  .placeholder-view { width: 64rpx; }
}

.content {
  flex: 1;
  height: 0;
  padding: 30rpx;
}

.section {
  margin-bottom: 30rpx;
  padding: 0 40rpx;
}

.menu-item {
  padding: 36rpx 0;
  border-bottom: 1rpx solid #f8f8f8;

  &:last-child { border-bottom: 0; }

  .item-left {
    .label {
      font-size: 30rpx;
      color: $text-main;
      font-weight: 500;
    }
    .menu-desc {
      display: block;
      margin-top: 8rpx;
      font-size: 22rpx;
      color: $text-sub;
    }
  }

  .icon-arrow {
    width: 32rpx;
    height: 32rpx;
    opacity: 0.2;
  }

  .val {
    font-size: 26rpx;
    color: $text-sub;
  }
}

.action-box {
  margin-top: 60rpx;

  .logout-btn {
    height: 100rpx;
    background: #fff;
    color: #E64340;
    font-size: 32rpx;
    font-weight: 700;
    border-radius: 20rpx;
    border: 1rpx solid rgba(230, 67, 64, 0.1);
    box-shadow: 0 8rpx 20rpx rgba(230, 67, 64, 0.05);

    &:active {
      background: #fff5f5;
      opacity: 0.8;
    }
  }
}

.version-info {
  margin-top: 40rpx;
  text-align: center;
  font-size: 22rpx;
  color: #ccc;
  padding-bottom: 40rpx;
}
</style>
