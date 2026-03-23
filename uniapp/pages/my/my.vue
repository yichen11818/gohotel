<template>
  <scroll-view class="page" scroll-y>
    <view class="header">
      <view class="profile-card" @click="handleProfileClick">
        <image class="avatar" :src="profile.avatar || defaultAvatar" mode="aspectFill" />
        <view class="profile-info">
          <text class="nickname">{{ isLoggedIn ? profile.nickname : '点击登录' }}</text>
          <text class="meta">{{ isLoggedIn ? `${profile.memberLevel} · 用户名 ${profile.username}` : '登录后查看订单与个人资料' }}</text>
        </view>
      </view>
      <button class="settings-btn" size="mini" @click="goToSettings">设置</button>
    </view>

    <view class="content">
      <view class="section stats-section" v-if="isLoggedIn">
        <view class="stat-item">
          <text class="stat-value">{{ profile.points }}</text>
          <text class="stat-label">积分</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">¥{{ profile.balance }}</text>
          <text class="stat-label">余额</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ profile.memberLevel }}</text>
          <text class="stat-label">会员等级</text>
        </view>
      </view>

      <view class="section">
        <view class="menu-item" @click="goToOrders('all')">
          <text>我的订单</text>
          <text class="menu-desc">查看全部预订</text>
        </view>
        <view class="menu-item" @click="goToOrders('pending')">
          <text>待确认订单</text>
          <text class="menu-desc">跟进最新预订状态</text>
        </view>
        <view class="menu-item" @click="goToOrders('confirmed')">
          <text>待入住订单</text>
          <text class="menu-desc">查看待入住安排</text>
        </view>
      </view>

      <view class="section" v-if="isLoggedIn">
        <view class="menu-item" @click="goToProfile">
          <text>个人资料</text>
          <text class="menu-desc">编辑姓名、手机号与头像</text>
        </view>
        <view class="menu-item" @click="goToSecurity">
          <text>修改密码</text>
          <text class="menu-desc">更新登录密码</text>
        </view>
      </view>

      <view class="section">
        <view class="menu-item" @click="goToService">
          <text>酒店服务</text>
          <text class="menu-desc">联系前台、查看服务说明</text>
        </view>
        <view class="menu-item" @click="goToAbout">
          <text>关于我们</text>
          <text class="menu-desc">查看酒店介绍</text>
        </view>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { user } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'

const defaultAvatar = 'https://dummyimage.com/160x160/f3ede5/8b6b47&text=User'
const isLoggedIn = ref(false)
const profile = ref({
  avatar: '',
  nickname: '用户',
  username: '',
  memberLevel: '普通会员',
  points: 0,
  balance: '0.00',
})

const loadProfile = async () => {
  isLoggedIn.value = Boolean(uni.getStorageSync(TOKEN_KEY))

  if (!isLoggedIn.value) {
    profile.value = {
      avatar: '',
      nickname: '用户',
      username: '',
      memberLevel: '普通会员',
      points: 0,
      balance: '0.00',
    }
    return
  }

  try {
    const data = await user.getUserInfo()
    profile.value = {
      avatar: data.avatar || '',
      nickname: data.nickname || data.username || '用户',
      username: data.username || '',
      memberLevel: data.member_level_name || '普通会员',
      points: Number(data.points || 0),
      balance: data.balance || '0.00',
    }
  } catch (error) {
    console.error('load profile failed:', error)
  }
}

const requireLogin = (callback) => {
  if (!isLoggedIn.value) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }
  callback()
}

const handleProfileClick = () => {
  if (!isLoggedIn.value) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }
  goToProfile()
}

const goToOrders = (type) => {
  requireLogin(() => {
    uni.navigateTo({ url: `/pages/orders/orders?type=${type}` })
  })
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

const goToSettings = () => {
  uni.navigateTo({ url: '/pages/settings/settings' })
}

const goToService = () => {
  uni.switchTab({ url: '/pages/service/service' })
}

const goToAbout = () => {
  uni.navigateTo({ url: '/pages/about/about' })
}

onShow(() => {
  loadProfile()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: linear-gradient(180deg, #f4eadb 0%, #f6f7fb 220rpx);
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 20rpx;
  padding: 88rpx 24rpx 24rpx;
}

.profile-card {
  display: flex;
  align-items: center;
  gap: 20rpx;
  flex: 1;
}

.avatar {
  width: 112rpx;
  height: 112rpx;
  border-radius: 50%;
  background: #fff;
}

.nickname {
  display: block;
  font-size: 34rpx;
  font-weight: 700;
  color: #1f2937;
}

.meta {
  display: block;
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #4b5563;
  line-height: 1.6;
}

.settings-btn {
  margin: 0;
  background: rgba(255, 255, 255, 0.92);
  border-radius: 999rpx;
}

.content {
  padding: 0 24rpx 32rpx;
}

.section {
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.stats-section {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12rpx;
  padding: 28rpx 0;
}

.stat-item {
  text-align: center;
}

.stat-value {
  display: block;
  font-size: 28rpx;
  font-weight: 700;
  color: #b7791f;
}

.stat-label,
.menu-desc {
  display: block;
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #6b7280;
}

.menu-item {
  padding: 28rpx 24rpx;
  border-bottom: 1px solid #f3f4f6;
}

.menu-item:last-child {
  border-bottom: 0;
}
</style>
