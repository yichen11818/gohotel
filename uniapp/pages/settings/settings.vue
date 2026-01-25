<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">设置</text>
        </view>
      </template>
    </TnNavbar>

    <scroll-view class="main-content" scroll-y>
      <!-- 账号设置 -->
      <view class="section">
        <view class="section-title">账号设置</view>
        <view class="setting-item" @click="goToProfile">
          <view class="item-left">
            <TnIcon name="user" color="#C29D71" size="40" />
            <text class="item-text">个人信息</text>
          </view>
          <view class="item-right">
            <TnIcon name="right" color="#ccc" size="28" />
          </view>
        </view>
        
        <view class="setting-item" @click="goToSecurity">
          <view class="item-left">
            <TnIcon name="lock" color="#C29D71" size="40" />
            <text class="item-text">账号安全</text>
          </view>
          <view class="item-right">
            <TnIcon name="right" color="#ccc" size="28" />
          </view>
        </view>
        
        <view class="setting-item" @click="goToAddress">
          <view class="item-left">
            <TnIcon name="location" color="#C29D71" size="40" />
            <text class="item-text">收货地址</text>
          </view>
          <view class="item-right">
            <TnIcon name="right" color="#ccc" size="28" />
          </view>
        </view>
      </view>

      <!-- 通知设置 -->
      <view class="section">
        <view class="section-title">通知设置</view>
        <view class="setting-item">
          <view class="item-left">
            <TnIcon name="notification" color="#C29D71" size="40" />
            <text class="item-text">消息通知</text>
          </view>
          <view class="item-right">
            <TnSwitch v-model="settings.notification" bg-color="#C29D71" />
          </view>
        </view>
        
        <view class="setting-item">
          <view class="item-left">
            <TnIcon name="bell" color="#C29D71" size="40" />
            <text class="item-text">订单提醒</text>
          </view>
          <view class="item-right">
            <TnSwitch v-model="settings.orderNotice" bg-color="#C29D71" />
          </view>
        </view>
      </view>

      <!-- 其他设置 -->
      <view class="section">
        <view class="section-title">其他</view>
        <view class="setting-item" @click="goToAbout">
          <view class="item-left">
            <TnIcon name="info-circle" color="#C29D71" size="40" />
            <text class="item-text">关于我们</text>
          </view>
          <view class="item-right">
            <text class="version-text">v1.0.0</text>
            <TnIcon name="right" color="#ccc" size="28" />
          </view>
        </view>
        
        <view class="setting-item" @click="goToPrivacy">
          <view class="item-left">
            <TnIcon name="file-text" color="#C29D71" size="40" />
            <text class="item-text">隐私政策</text>
          </view>
          <view class="item-right">
            <TnIcon name="right" color="#ccc" size="28" />
          </view>
        </view>
        
        <view class="setting-item" @click="clearCache">
          <view class="item-left">
            <TnIcon name="delete" color="#C29D71" size="40" />
            <text class="item-text">清除缓存</text>
          </view>
          <view class="item-right">
            <text class="cache-text">{{ cacheSize }}</text>
            <TnIcon name="right" color="#ccc" size="28" />
          </view>
        </view>
      </view>

      <!-- 退出登录 -->
      <view class="logout-section">
        <TnButton 
          shape="round" 
          size="xl" 
          width="100%" 
          height="100rpx"
          bg-color="#fff"
          text-color="#FF3B30"
          border
          @click="handleLogout"
        >
          退出登录
        </TnButton>
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnSwitch from '@/uni_modules/tuniaoui-vue3/components/switch/src/switch.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'

const settings = ref({
  notification: true,
  orderNotice: true
})

const cacheSize = ref('12.5MB')

onLoad(() => {
  loadSettings()
})

const loadSettings = () => {
  // TODO: 从本地存储加载设置
}

const goBack = () => {
  uni.navigateBack()
}

const goToProfile = () => {
  uni.navigateTo({
    url: '/pages/profile/profile'
  })
}

const goToSecurity = () => {
  uni.navigateTo({
    url: '/pages/security/security'
  })
}

const goToAddress = () => {
  uni.navigateTo({
    url: '/pages/address/address'
  })
}

const goToAbout = () => {
  uni.navigateTo({
    url: '/pages/about/about'
  })
}

const goToPrivacy = () => {
  uni.navigateTo({
    url: '/pages/privacy/privacy'
  })
}

const clearCache = () => {
  uni.showModal({
    title: '提示',
    content: '确定要清除缓存吗？',
    success: (res) => {
      if (res.confirm) {
        uni.showLoading({ title: '清除中...' })
        setTimeout(() => {
          uni.hideLoading()
          cacheSize.value = '0KB'
          uni.showToast({
            title: '清除成功',
            icon: 'success'
          })
        }, 1000)
      }
    }
  })
}

const handleLogout = () => {
  uni.showModal({
    title: '提示',
    content: '确定要退出登录吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await user.logout()
        } catch (e) {}
        uni.reLaunch({
          url: '/pages/index/index'
        })
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #F5F5F5;
}

.navbar-content {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60rpx;
  height: 60rpx;
  
  &:active {
    opacity: 0.6;
  }
}

.navbar-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #333;
}

.main-content {
  height: calc(100vh - 88rpx);
  padding: 24rpx;
}

.section {
  background: #fff;
  border-radius: 24rpx;
  padding: 0 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.section-title {
  padding: 32rpx 0 16rpx;
  font-size: 28rpx;
  font-weight: 600;
  color: #999;
}

.setting-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 0;
  border-bottom: 1rpx solid #F0F0F0;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:active {
    opacity: 0.7;
  }
}

.item-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.item-text {
  font-size: 30rpx;
  color: #333;
  font-weight: 500;
}

.item-right {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.version-text,
.cache-text {
  font-size: 26rpx;
  color: #999;
}

.logout-section {
  margin-top: 48rpx;
  padding: 0 8rpx;
}

.bottom-placeholder {
  height: 40rpx;
}
</style>



