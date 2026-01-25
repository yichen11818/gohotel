<template>
  <view class="container">
    <!-- 顶部动态渐变背景 -->
    <view class="header-bg">
      <view class="gradient-layer gradient-1"></view>
      <view class="gradient-layer gradient-2"></view>
      <view class="gradient-layer gradient-3"></view>
    </view>
    
    <!-- 顶部内容区 -->
    <view class="header-content">
      <!-- 右上角设置按钮 -->
      <view class="top-bar">
        <view class="menu-btn" @click="showMenu">
          <TnIcon name="more-vertical" color="#666" size="44" />
        </view>
        <view class="settings-btn" @click="goToSettings">
          <TnIcon name="setting" color="#666" size="44" />
        </view>
      </view>
      
      <!-- 个人信息 -->
      <view class="profile-section" @click="handleLogin">
        <image class="avatar" :src="userInfo.avatar || defaultAvatar" mode="aspectFill"></image>
        <view class="user-info">
          <text class="username">{{ userInfo.nickname || '点击登录' }}</text>
          <view class="member-info">
            <text class="member-label">{{ userInfo.memberLevel || 'N口27187358829' }}</text>
            <TnIcon v-if="userInfo.isVip" name="vip-fill" color="#C29D71" size="28" />
          </view>
        </view>
      </view>
      
      <!-- 数据统计 -->
      <view class="stats-row">
        <view class="stat-item" @click="goToWallet">
          <text class="stat-value">¥ {{ userStats.balance }}</text>
          <text class="stat-label">储值</text>
        </view>
        <view class="stat-divider"></view>
        <view class="stat-item" @click="goToPoints">
          <text class="stat-value">{{ userStats.points }}</text>
          <text class="stat-label">积分</text>
        </view>
        <view class="stat-divider"></view>
        <view class="stat-item" @click="goToCoupons">
          <text class="stat-value">{{ userStats.coupons }}</text>
          <text class="stat-label">代金券</text>
        </view>
      </view>
    </view>

    

    <!-- 功能卡片 -->
    <view class="function-card">
      <view class="function-row">
        <view class="function-item" @click="navigateTo('points')">
          <TnIcon name="medal" color="#C29D71" size="44" />
          <text class="function-text">房屋积分</text>
          <text class="function-desc">你的积分</text>
        </view>
        <view class="function-item" @click="navigateTo('alarm')">
          <TnIcon name="time-circle" color="#C29D71" size="44" />
          <text class="function-text">起床特性</text>
          <text class="function-desc">设置叫醒</text>
        </view>
        <view class="function-item" @click="navigateTo('shop')">
          <TnIcon name="cart" color="#C29D71" size="44" />
          <text class="function-text">会员购物</text>
          <text class="function-desc">浏览</text>
        </view>
        <view class="function-item" @click="navigateTo('license')">
          <TnIcon name="user-circle" color="#C29D71" size="44" />
          <text class="function-text">司机建档</text>
          <text class="function-desc">证件</text>
        </view>
      </view>
      
      <view class="function-row">
        <view class="function-item" @click="navigateTo('wifi')">
          <TnIcon name="wifi" color="#C29D71" size="44" />
          <text class="function-text">WiFi账单</text>
          <text class="function-desc">查看账单</text>
        </view>
        <view class="function-item" @click="navigateTo('recharge')">
          <TnIcon name="notification" color="#C29D71" size="44" />
          <text class="function-text">充值康帝卡</text>
          <text class="function-desc">优惠充值</text>
        </view>
        <view class="function-item" @click="navigateTo('rights')">
          <TnIcon name="close-circle" color="#C29D71" size="44" />
          <text class="function-text">更多权益</text>
          <text class="function-desc">可提前</text>
        </view>
      </view>
    </view>
 <!-- 预定订单 -->
 <view class="section">
      <view class="section-header">
        <text class="section-title">预定订单</text>
        <view class="section-more" @click="goToOrders('all')">
          <text class="more-text">全部订单</text>
          <TnIcon name="right" color="#999" size="28" />
        </view>
      </view>
      
      <view class="order-types">
        <view class="order-type" @click="goToOrders('pending')">
          <view class="order-icon-wrapper">
            <image class="order-icon-img" src="/static/icon/1.png" mode="aspectFit"></image>
          </view>
          <text class="order-type-text">待付款</text>
        </view>
        <view class="order-type" @click="goToOrders('checkin')">
          <view class="order-icon-wrapper">
            <image class="order-icon-img" src="/static/icon/2.png" mode="aspectFit"></image>
          </view>
          <text class="order-type-text">待入住</text>
        </view>
        <view class="order-type" @click="goToOrders('paid')">
          <view class="order-icon-wrapper">
            <image class="order-icon-img" src="/static/icon/3.png" mode="aspectFit"></image>
          </view>
          <text class="order-type-text">已支付</text>
        </view>
        <view class="order-type" @click="goToOrders('cancelled')">
          <view class="order-icon-wrapper">
            <image class="order-icon-img" src="/static/icon/4.png" mode="aspectFit"></image>
          </view>
          <text class="order-type-text">已取消</text>
        </view>
      </view>
    </view>
    <!-- 常用功能 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">常用功能</text>
        <view class="section-more" @click="showAllFunctions">
          <text class="more-text">全部功能</text>
          <TnIcon name="right" color="#999" size="28" />
        </view>
      </view>
      
      <view class="common-functions">
        <view class="common-item" @click="navigateTo('upgrade')">
          <view class="common-icon">
            <TnIcon name="vip" color="#C29D71" size="44" />
          </view>
          <text class="common-text">升级会员</text>
        </view>
        <view class="common-item" @click="navigateTo('strategy')">
          <view class="common-icon">
            <TnIcon name="message-circle" color="#C29D71" size="44" />
          </view>
          <text class="common-text">降费攻略</text>
        </view>
        <view class="common-item" @click="navigateTo('shopping')">
          <view class="common-icon">
            <TnIcon name="shopping-bag" color="#C29D71" size="44" />
          </view>
          <text class="common-text">购物订单</text>
        </view>
        <view class="common-item" @click="navigateTo('mall')">
          <view class="common-icon">
            <TnIcon name="list" color="#C29D71" size="44" />
          </view>
          <text class="common-text">商城订单</text>
        </view>
        <view class="common-item" @click="navigateTo('feedback')">
          <view class="common-icon">
            <TnIcon name="edit" color="#C29D71" size="44" />
          </view>
          <text class="common-text">反馈意见</text>
        </view>
        <view class="common-item" @click="navigateTo('invoice')">
          <view class="common-icon">
            <TnIcon name="file-text" color="#C29D71" size="44" />
          </view>
          <text class="common-text">切换发票</text>
        </view>
      </view>
    </view>

    <!-- 底部占位 -->
    <view class="bottom-placeholder"></view>
  </view>
</template>

<script setup>
import { user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'
import { ref, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'

// 默认头像
const defaultAvatar = 'https://images.unsplash.com/photo-1535713875002-d1d0cf377fde?w=200&h=200&fit=crop'

// 用户信息
const userInfo = ref({
  avatar: '',
  nickname: '',
  memberLevel: '',
  isVip: false
})

// 用户统计数据
const userStats = ref({
  balance: '0.00',
  points: '0',
  coupons: '0'
})

// 页面加载
onLoad(() => {
  refreshData()
})

// 页面显示时刷新数据
onShow(() => {
  refreshData()
})

const refreshData = () => {
  const token = uni.getStorageSync(TOKEN_KEY)
  if (token) {
    loadUserInfo()
    loadUserStats()
  } else {
    // 未登录状态重置数据
    userInfo.value = {
      avatar: '',
      nickname: '',
      memberLevel: '',
      isVip: false
    }
    userStats.value = {
      balance: '0.00',
      points: '0',
      coupons: '0'
    }
  }
}

// 加载用户信息
const loadUserInfo = async () => {
  try {
    const data = await user.getUserInfo()
    if (data) {
      userInfo.value = {
        avatar: data.avatar || defaultAvatar,
        nickname: data.nickname || '用户',
        memberLevel: data.member_level_name || '普通会员',
        isVip: data.is_vip || false
      }
      uni.setStorageSync(USER_INFO_KEY, data)
    }
  } catch (error) {
    console.error('Failed to load user info:', error)
  }
}

// 加载用户统计数据
const loadUserStats = async () => {
  try {
    const [pointsRes, couponsRes] = await Promise.all([
      user.getUserPoints(),
      user.getCoupons({ status: 'unused' })
    ])
    
    userStats.value = {
      balance: userInfo.value.balance || '0.00',
      points: String(pointsRes?.points || 0),
      coupons: String(couponsRes?.length || 0)
    }
  } catch (error) {
    console.error('Failed to load user stats:', error)
  }
}

// 处理登录
const handleLogin = () => {
  const token = uni.getStorageSync(TOKEN_KEY)
  if (!token) {
    uni.navigateTo({
      url: '/pages/login/login'
    })
  } else {
    uni.navigateTo({
      url: '/pages/settings/settings'
    })
  }
}

// 前往设置
const goToSettings = () => {
  uni.navigateTo({
    url: '/pages/settings/settings'
  })
}

// 显示菜单
const showMenu = () => {
  uni.showActionSheet({
    itemList: ['分享给好友', '意见反馈', '关于我们'],
    success: (res) => {
      if (res.tapIndex === 0) {
        // 分享
        uni.showToast({ title: '分享功能开发中', icon: 'none' })
      } else if (res.tapIndex === 1) {
        // 意见反馈
        uni.showToast({ title: '反馈功能开发中', icon: 'none' })
      } else if (res.tapIndex === 2) {
        // 关于我们
        uni.navigateTo({ url: '/pages/about/about' })
      }
    }
  })
}

// 前往钱包
const goToWallet = () => {
  checkLoginAndNavigate('/pages/wallet/wallet')
}

// 前往积分
const goToPoints = () => {
  checkLoginAndNavigate('/pages/points/points')
}

// 前往代金券
const goToCoupons = () => {
  checkLoginAndNavigate('/pages/coupons/coupons')
}

// 前往订单
const goToOrders = (type) => {
  checkLoginAndNavigate(`/pages/orders/orders?type=${type}`)
}

// 显示全部功能
const showAllFunctions = () => {
  uni.showToast({
    title: '功能开发中',
    icon: 'none'
  })
}

// 导航到指定页面
const navigateTo = (type) => {
  uni.showToast({
    title: `${type}功能开发中`,
    icon: 'none'
  })
}

// 检查登录并导航
const checkLoginAndNavigate = (url) => {
  const token = uni.getStorageSync(TOKEN_KEY)
  if (!token) {
    uni.showModal({
      title: '提示',
      content: '请先登录',
      success: (res) => {
        if (res.confirm) {
          uni.navigateTo({
            url: '/pages/login/login'
          })
        }
      }
    })
    return
  }
  
  uni.navigateTo({ url })
}
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #FAFAFA;
  padding-bottom: 120rpx;
  position: relative;
  overflow: hidden;
}

/* 顶部动态渐变背景 */
.header-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 500rpx;
  z-index: 0;
  overflow: hidden;
}

.gradient-layer {
  position: absolute;
  width: 200%;
  height: 100%;
  opacity: 0.6;
  filter: blur(60rpx);
}

.gradient-1 {
  background: linear-gradient(135deg, 
    rgba(255, 182, 193, 0.5) 0%,
    rgba(221, 160, 221, 0.5) 25%,
    rgba(176, 224, 230, 0.5) 50%,
    rgba(255, 218, 185, 0.5) 75%,
    rgba(255, 182, 193, 0.5) 100%
  );
  animation: gradientShift1 15s ease-in-out infinite;
}

.gradient-2 {
  background: linear-gradient(45deg,
    rgba(230, 230, 250, 0.4) 0%,
    rgba(255, 240, 245, 0.4) 33%,
    rgba(240, 248, 255, 0.4) 66%,
    rgba(230, 230, 250, 0.4) 100%
  );
  animation: gradientShift2 20s ease-in-out infinite;
  animation-delay: -5s;
}

.gradient-3 {
  background: linear-gradient(225deg,
    rgba(255, 228, 225, 0.3) 0%,
    rgba(255, 250, 240, 0.3) 50%,
    rgba(245, 245, 220, 0.3) 100%
  );
  animation: gradientShift3 25s ease-in-out infinite;
  animation-delay: -10s;
}

/* 渐变动画 */
@keyframes gradientShift1 {
  0%, 100% {
    transform: translate(-25%, -10%) rotate(0deg);
  }
  25% {
    transform: translate(-15%, -5%) rotate(5deg);
  }
  50% {
    transform: translate(-5%, 0%) rotate(10deg);
  }
  75% {
    transform: translate(-15%, -5%) rotate(5deg);
  }
}

@keyframes gradientShift2 {
  0%, 100% {
    transform: translate(0%, 0%) rotate(0deg);
  }
  33% {
    transform: translate(-10%, 5%) rotate(-5deg);
  }
  66% {
    transform: translate(-20%, -5%) rotate(-10deg);
  }
}

@keyframes gradientShift3 {
  0%, 100% {
    transform: translate(-10%, 5%) rotate(0deg);
  }
  50% {
    transform: translate(-30%, -10%) rotate(15deg);
  }
}

/* 头部内容区 */
.header-content {
  position: relative;
  z-index: 1;
  padding: 0 24rpx;
  animation: fadeInDown 0.6s ease-out;
}

/* 顶部操作栏 */
.top-bar {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 16rpx;
  padding: 32rpx 0 24rpx;
}

.menu-btn,
.settings-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64rpx;
  height: 64rpx;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10rpx);
  -webkit-backdrop-filter: blur(10rpx);
  border-radius: 50%;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
  
  &:active {
    transform: scale(0.9);
    background: rgba(255, 255, 255, 1);
  }
}

/* 个人信息区域 */
.profile-section {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 16rpx 0 32rpx;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  border: 4rpx solid rgba(255, 255, 255, 0.9);
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.1);
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.username {
  font-size: 40rpx;
  font-weight: 700;
  color: #333;
  letter-spacing: 0.5rpx;
}

.member-info {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.member-label {
  font-size: 22rpx;
  color: #666;
}

/* 数据统计 */
.stats-row {
  display: flex;
  align-items: center;
  justify-content: space-around;
  padding: 32rpx;
  margin: 0 0 24rpx;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20rpx);
  -webkit-backdrop-filter: blur(20rpx);
  border-radius: 24rpx;
  box-shadow: 0 8rpx 24rpx rgba(0, 0, 0, 0.06);
  border: 1rpx solid rgba(255, 255, 255, 0.8);
}

.stat-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
  
  &:active {
    opacity: 0.7;
  }
}

.stat-value {
  font-size: 44rpx;
  font-weight: 700;
  color: #333;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

.stat-label {
  font-size: 24rpx;
  color: #999;
}

.stat-divider {
  width: 1rpx;
  height: 60rpx;
  background: linear-gradient(to bottom, transparent, #E5E5E5, transparent);
}

/* 功能卡片 */
.function-card {
  margin: 0 24rpx 24rpx;
  padding: 40rpx 32rpx;
  /* 背景色调整：加入微弱的紫红暖色调，模拟奢华感 */
  background: linear-gradient(135deg, #2E2528 0%, #1B1B1E 100%);
  border-radius: 32rpx;
  box-shadow: 0 16rpx 48rpx rgba(0, 0, 0, 0.2);
  border: 1rpx solid rgba(255, 255, 255, 0.05);
  animation: fadeInUp 0.6s ease-out 0.1s backwards;
  position: relative;
  overflow: hidden;
  
  /* 底部纹饰 */
  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-image: 
      /* 组合纹饰：左紧右疏，中间点缀 */
      url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 700 350' preserveAspectRatio='none'%3E%3Cdefs%3E%3ClinearGradient id='g_left' x1='0%25' y1='100%25' x2='100%25' y2='0%25'%3E%3Cstop offset='0%25' stop-color='%23E5CFA0' stop-opacity='0.25'/%3E%3Cstop offset='100%25' stop-color='%23E5CFA0' stop-opacity='0'/%3E%3C/linearGradient%3E%3ClinearGradient id='g_right' x1='100%25' y1='100%25' x2='0%25' y2='0%25'%3E%3Cstop offset='0%25' stop-color='%23E5CFA0' stop-opacity='0.1'/%3E%3Cstop offset='100%25' stop-color='%23E5CFA0' stop-opacity='0'/%3E%3C/linearGradient%3E%3C/defs%3E%3C!-- 左侧：紧致的束状线条 (密集起始) --%3E%3Cpath d='M0 350 Q 150 150 400 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3Cpath d='M15 350 Q 165 150 415 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3Cpath d='M30 350 Q 180 150 430 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3Cpath d='M45 350 Q 195 150 445 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3Cpath d='M60 350 Q 210 150 460 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3Cpath d='M75 350 Q 225 150 475 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3Cpath d='M90 350 Q 240 150 490 0' stroke='url(%23g_left)' fill='none' stroke-width='0.8'/%3E%3C!-- 中间底部：少量延伸 --%3E%3Cpath d='M150 350 Q 300 250 550 100' stroke='url(%23g_left)' fill='none' stroke-width='0.6'/%3E%3Cpath d='M200 350 Q 350 250 600 100' stroke='url(%23g_left)' fill='none' stroke-width='0.6'/%3E%3C!-- 右侧：稀疏的同心圆弧 --%3E%3Ccircle cx='700' cy='350' r='150' stroke='url(%23g_right)' fill='none' stroke-width='0.6'/%3E%3Ccircle cx='700' cy='350' r='250' stroke='url(%23g_right)' fill='none' stroke-width='0.6'/%3E%3Ccircle cx='700' cy='350' r='350' stroke='url(%23g_right)' fill='none' stroke-width='0.6'/%3E%3Ccircle cx='700' cy='350' r='450' stroke='url(%23g_right)' fill='none' stroke-width='0.6'/%3E%3Ccircle cx='700' cy='350' r='550' stroke='url(%23g_right)' fill='none' stroke-width='0.6'/%3E%3C/svg%3E");
    background-size: 100% 100%;
    background-position: center bottom;
    background-repeat: no-repeat;
    opacity: 1;
    pointer-events: none;
    z-index: 0;
    animation: waveFlowComplex 8s ease-in-out infinite;
  }
  
  /* 确保内容在纹饰之上 */
  .function-row {
    position: relative;
    z-index: 1;
  }
}

.function-row {
  display: flex;
  justify-content: space-between;
  margin-bottom: 48rpx;
  
  &:last-child {
    margin-bottom: 0;
  }
}

.function-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  transition: all 0.3s ease;
  
  &:active {
    transform: scale(0.95);
    opacity: 0.8;
  }
}

.function-text {
  font-size: 26rpx;
  color: #FFFFFF;
  font-weight: 500;
  text-align: center;
}

.function-desc {
  font-size: 20rpx;
  color: rgba(255, 255, 255, 0.6);
  text-align: center;
}

/* 分区 */
.section {
  margin: 0 24rpx 24rpx;
  padding: 32rpx;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20rpx);
  -webkit-backdrop-filter: blur(20rpx);
  border-radius: 32rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.04);
  border: 1rpx solid rgba(255, 255, 255, 0.6);
  animation: fadeInUp 0.6s ease-out 0.2s backwards;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 32rpx;
}

.section-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #333;
  letter-spacing: 0.5rpx;
}

.section-more {
  display: flex;
  align-items: center;
  gap: 8rpx;
  
  &:active {
    opacity: 0.6;
  }
}

.more-text {
  font-size: 26rpx;
  color: #999;
}

/* 订单类型 */
.order-types {
  display: flex;
  justify-content: space-around;
  padding: 20rpx 0;
}

.order-type {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  transition: all 0.3s ease;
  
  &:active {
    opacity: 0.7;
    transform: scale(0.95);
  }
}

.order-icon-wrapper {
  width: 100rpx;
  height: 100rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
}

.order-icon-img {
  width: 100%;
  height: 100%;
}

.order-type-text {
  font-size: 26rpx;
  color: #666;
  font-weight: 500;
}

/* 常用功能 */
.common-functions {
  display: flex;
  flex-wrap: wrap;
  gap: 32rpx;
}

.common-item {
  width: calc((100% - 64rpx) / 3);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  padding: 20rpx 0;
  transition: all 0.3s ease;
  
  &:active {
    transform: scale(0.95);
    opacity: 0.7;
  }
}

.common-icon {
  width: 80rpx;
  height: 80rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, rgba(194, 157, 113, 0.12) 0%, rgba(194, 157, 113, 0.06) 100%);
  border-radius: 50%;
}

.common-text {
  font-size: 26rpx;
  color: #666;
  font-weight: 500;
  text-align: center;
}

/* 底部占位 */
.bottom-placeholder {
  height: 40rpx;
}

/* 动画 */
@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes waveFlowComplex {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}
</style>
