<template>
  <view class="page">
    <!-- 头部：背景与基本信息 -->
    <view class="user-header">
      <view class="header-content" :style="{ paddingTop: statusBarHeight + 'px' }">
        <view class="top-actions flex-between">
          <view class="member-tag flex-center" v-if="isLoggedIn">
            <image class="icon-crown" src="/static/icons/crown.png" mode="aspectFit" />
            <text>{{ profile.memberLevel }}</text>
          </view>
          <view class="action-icons flex-center">
            <image class="icon-action" src="/static/icons/settings-white.png" @click="goToSettings" mode="aspectFit" />
          </view>
        </view>

        <view class="user-profile flex-center" @click="handleProfileClick">
          <image class="avatar" :src="profile.avatar || defaultAvatar" mode="aspectFill" />
          <view class="info">
            <text class="nickname">{{ isLoggedIn ? profile.nickname : '立即登录' }}</text>
            <text class="desc">{{ isLoggedIn ? `用户名: ${profile.username}` : '登录享受更多会员权益' }}</text>
          </view>
        </view>
      </view>
      <view class="header-bg"></view>
    </view>

    <view class="content">
      <!-- 资产面板 -->
      <view class="asset-card premium-card flex-between">
        <view class="asset-item flex-center" @click="goToWallet">
          <text class="val">{{ isLoggedIn ? profile.balance : '--' }}</text>
          <text class="label">余额</text>
        </view>
        <view class="divider"></view>
        <view class="asset-item flex-center" @click="goToPoints">
          <text class="val">{{ isLoggedIn ? profile.points : '--' }}</text>
          <text class="label">积分</text>
        </view>
        <view class="divider"></view>
        <view class="asset-item flex-center" @click="goToCoupons">
          <text class="val">{{ isLoggedIn ? '3' : '--' }}</text>
          <text class="label">优惠券</text>
        </view>
      </view>

      <!-- 订单面板 -->
      <view class="order-panel premium-card">
        <view class="panel-header flex-between" @click="goToOrders('all')">
          <text class="title">我的订单</text>
          <view class="all flex-center">
            <text>全部订单</text>
            <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
          </view>
        </view>
        <view class="order-types flex-between">
          <view class="type-item flex-center" @click="goToOrders('pending')">
            <image class="icon" src="/static/icons/order-pay.png" mode="aspectFit" />
            <text>待确认</text>
          </view>
          <view class="type-item flex-center" @click="goToOrders('confirmed')">
            <image class="icon" src="/static/icons/order-stay.png" mode="aspectFit" />
            <text>待入住</text>
          </view>
          <view class="type-item flex-center" @click="goToOrders('checkout')">
            <image class="icon" src="/static/icons/order-done.png" mode="aspectFit" />
            <text>已完成</text>
          </view>
          <view class="type-item flex-center" @click="goToOrders('cancelled')">
            <image class="icon" src="/static/icons/order-refund.png" mode="aspectFit" />
            <text>已取消</text>
          </view>
        </view>
      </view>

      <!-- 常用功能 -->
      <view class="menu-list premium-card">
        <view class="menu-item flex-between" @click="goToProfile">
          <view class="left flex-center">
            <image class="icon" src="/static/icons/menu-profile.png" mode="aspectFit" />
            <text>个人资料</text>
          </view>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
        <view class="menu-item flex-between" @click="goToSecurity">
          <view class="left flex-center">
            <image class="icon" src="/static/icons/menu-safe.png" mode="aspectFit" />
            <text>账号安全</text>
          </view>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
        <view class="menu-item flex-between" @click="goToService">
          <view class="left flex-center">
            <image class="icon" src="/static/icons/menu-service.png" mode="aspectFit" />
            <text>酒店服务</text>
          </view>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
        <view class="menu-item flex-between" @click="goToAbout">
          <view class="left flex-center">
            <image class="icon" src="/static/icons/menu-about.png" mode="aspectFit" />
            <text>关于我们</text>
          </view>
          <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { user } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'

const statusBarHeight = ref(44)
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
    profile.value = { avatar: '', nickname: '用户', username: '', memberLevel: '普通会员', points: 0, balance: '0.00' }
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

const handleProfileClick = () => isLoggedIn.value ? goToProfile() : uni.navigateTo({ url: '/pages/login/login' })
const goToOrders = (type) => requireLogin(() => uni.navigateTo({ url: `/pages/orders/orders?type=${type}` }))
const goToProfile = () => requireLogin(() => uni.navigateTo({ url: '/pages/profile/profile' }))
const goToSecurity = () => requireLogin(() => uni.navigateTo({ url: '/pages/security/security' }))
const goToSettings = () => uni.navigateTo({ url: '/pages/settings/settings' })
const goToService = () => uni.switchTab({ url: '/pages/service/service' })
const goToAbout = () => uni.navigateTo({ url: '/pages/about/about' })
const goToWallet = () => requireLogin(() => uni.navigateTo({ url: '/pages/wallet/wallet' }))
const goToPoints = () => requireLogin(() => uni.navigateTo({ url: '/pages/points/points' }))
const goToCoupons = () => requireLogin(() => uni.navigateTo({ url: '/pages/coupons/coupons' }))

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
})

onShow(() => loadProfile())
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: $bg-color;
  padding-bottom: 60rpx;
}

/* 头部 */
.user-header {
  position: relative;
  height: 480rpx;

  .header-content {
    position: relative;
    z-index: 10;
    padding: 0 40rpx;
  }

  .top-actions {
    height: 88rpx;
    margin-bottom: 20rpx;

    .member-tag {
      padding: 6rpx 20rpx;
      background: rgba(0,0,0,0.2);
      border-radius: 100rpx;
      color: #fff;
      font-size: 20rpx;
      backdrop-filter: blur(4px);

      .icon-crown {
        width: 24rpx;
        height: 24rpx;
        margin-right: 8rpx;
      }
    }

    .icon-action {
      width: 44rpx;
      height: 44rpx;
    }
  }

  .user-profile {
    justify-content: flex-start;
    padding-top: 20rpx;

    .avatar {
      width: 128rpx;
      height: 128rpx;
      border-radius: 50%;
      border: 4rpx solid rgba(255,255,255,0.3);
      background: #fff;
      margin-right: 24rpx;
    }

    .nickname {
      font-size: 40rpx;
      font-weight: 700;
      color: #fff;
      display: block;
      margin-bottom: 8rpx;
    }

    .desc {
      font-size: 24rpx;
      color: rgba(255,255,255,0.8);
    }
  }

  .header-bg {
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  }
}

.content {
  padding: 0 30rpx;
  margin-top: -60rpx;
  position: relative;
  z-index: 20;
}

/* 资产卡片 */
.asset-card {
  padding: 40rpx 0;
  margin-bottom: 30rpx;

  .asset-item {
    flex: 1;
    flex-direction: column;

    .val {
      font-size: 36rpx;
      font-weight: 700;
      color: $text-main;
      margin-bottom: 8rpx;
    }
    .label {
      font-size: 24rpx;
      color: $text-sub;
    }
  }

  .divider {
    width: 1rpx;
    height: 40rpx;
    background: #f0f0f0;
  }
}

/* 订单面板 */
.order-panel {
  padding: 0;
  margin-bottom: 30rpx;

  .panel-header {
    padding: 30rpx 32rpx;
    border-bottom: 1rpx solid #f8f8f8;

    .title {
      font-size: 30rpx;
      font-weight: 600;
      color: $text-main;
    }
    .all {
      font-size: 24rpx;
      color: $text-sub;
      .icon-arrow { width: 24rpx; height: 24rpx; margin-left: 4rpx; }
    }
  }

  .order-types {
    padding: 40rpx 20rpx;

    .type-item {
      flex: 1;
      flex-direction: column;

      .icon {
        width: 56rpx;
        height: 56rpx;
        margin-bottom: 16rpx;
      }
      text {
        font-size: 24rpx;
        color: $text-second;
      }
    }
  }
}

/* 菜单列表 */
.menu-list {
  padding: 0 32rpx;

  .menu-item {
    height: 110rpx;
    border-bottom: 1rpx solid #f8f8f8;
    &:last-child { border-bottom: 0; }

    .left {
      .icon {
        width: 40rpx;
        height: 40rpx;
        margin-right: 24rpx;
      }
      text {
        font-size: 28rpx;
        color: $text-main;
      }
    }

    .icon-arrow {
      width: 28rpx;
      height: 28rpx;
      opacity: 0.3;
    }
  }
}
</style>
