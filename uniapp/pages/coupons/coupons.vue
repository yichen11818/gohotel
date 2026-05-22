<template>
  <view class="page">
    <!-- 顶部状态栏占位 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <!-- 自定义导航 -->
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">优惠券中心</text>
      <view class="placeholder-view"></view>
    </view>

    <!-- 标签页 -->
    <view class="tabs-box">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item flex-center"
        :class="{ active: currentTab === tab.value }"
        @click="switchTab(tab.value)"
      >
        <text class="label">{{ tab.label }}</text>
        <view class="line" v-if="currentTab === tab.value"></view>
      </view>
    </view>

    <scroll-view class="list-container" scroll-y>
      <view v-if="loading && !couponList.length" class="state-box">
        <view class="loading-icon"></view>
        <text>正在加载...</text>
      </view>

      <view v-else-if="!couponList.length" class="state-box">
        <image class="empty-img" src="/static/icons/order-refund.png" mode="aspectFit" style="opacity: 0.1;" />
        <text>暂无优惠券</text>
      </view>

      <view class="coupon-list" v-else>
        <view
          v-for="coupon in couponList"
          :key="coupon.id"
          class="coupon-card"
          :class="[currentTab, { used: currentTab === 'used', expired: currentTab === 'expired' }]"
        >
          <view class="card-inner flex-between">
            <view class="coupon-left flex-center">
              <view class="price">
                <text class="symbol">¥</text>
                <text class="val">{{ coupon.amount }}</text>
              </view>
              <text class="condition">满{{ coupon.condition }}可用</text>
            </view>

            <view class="coupon-right">
              <view class="info-top">
                <text class="name text-ellipsis">{{ coupon.name }}</text>
                <text class="time">{{ coupon.validTime }}</text>
              </view>
              <view class="info-bottom flex-between">
                <text class="desc">全场通用</text>
                <button v-if="currentTab === 'unused'" class="use-btn" @click.stop="useCoupon(coupon)">立即使用</button>
                <image v-else-if="currentTab === 'used'" class="status-stamp" src="/static/icons/order-done.png" mode="aspectFit" />
                <image v-else class="status-stamp" src="/static/icons/order-refund.png" mode="aspectFit" />
              </view>
            </view>
          </view>
          <!-- 锯齿边缘视觉装饰 -->
          <view class="sawtooth-top"></view>
          <view class="sawtooth-bottom"></view>
        </view>
      </view>
      <view class="bottom-space"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { user } from '@/api/index.js'

const statusBarHeight = ref(44)
const tabs = [
  { label: '待使用', value: 'unused' },
  { label: '已使用', value: 'used' },
  { label: '已过期', value: 'expired' }
]

const currentTab = ref('unused')
const couponList = ref([])
const loading = ref(false)

const loadCoupons = async () => {
  loading.value = true
  try {
    const data = await user.getCoupons({ status: currentTab.value })
    couponList.value = data || []
  } catch (error) {
    console.error('Failed to load coupons:', error)
  } finally {
    loading.value = false
  }
}

const switchTab = (value) => {
  if (currentTab.value === value) return
  currentTab.value = value
  couponList.value = []
  loadCoupons()
}

const goBack = () => {
  uni.navigateBack()
}

const useCoupon = (coupon) => {
  uni.switchTab({ url: '/pages/index/index' })
}

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
  loadCoupons()
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

/* Tabs */
.tabs-box {
  display: flex;
  background: #fff;
  border-bottom: 1rpx solid #f8f8f8;

  .tab-item {
    flex: 1;
    padding: 24rpx 0;
    position: relative;
    flex-direction: column;

    .label {
      font-size: 28rpx;
      color: $text-sub;
      transition: all 0.2s;
    }

    .line {
      position: absolute;
      bottom: 0;
      width: 40rpx;
      height: 6rpx;
      background: $primary-color;
      border-radius: 10rpx;
    }

    &.active {
      .label {
        color: $text-main;
        font-weight: 700;
      }
    }
  }
}

.list-container {
  flex: 1;
  height: 0;
}

.state-box {
  padding: 200rpx 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: $text-sub;
  font-size: 26rpx;

  .empty-img {
    width: 200rpx;
    height: 200rpx;
    margin-bottom: 30rpx;
  }
}

.coupon-list {
  padding: 30rpx;
}

.coupon-card {
  height: 200rpx;
  background: #fff;
  margin-bottom: 30rpx;
  border-radius: 16rpx;
  position: relative;
  overflow: hidden;
  box-shadow: 0 8rpx 24rpx rgba(0,0,0,0.03);

  .card-inner {
    height: 100%;
    z-index: 2;
    position: relative;
  }

  .coupon-left {
    width: 220rpx;
    height: 100%;
    background: linear-gradient(135deg, #fdf8f2 0%, #faecd9 100%);
    flex-direction: column;
    border-right: 2rpx dashed #eee;

    .price {
      .symbol { font-size: 24rpx; color: $primary-color; font-weight: 700; }
      .val { font-size: 60rpx; color: $primary-color; font-weight: 700; }
    }

    .condition {
      font-size: 20rpx;
      color: #ad8551;
      margin-top: 4rpx;
    }
  }

  .coupon-right {
    flex: 1;
    height: 100%;
    padding: 24rpx 30rpx;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    box-sizing: border-box;

    .info-top {
      .name {
        font-size: 30rpx;
        font-weight: 700;
        color: $text-main;
        display: block;
      }
      .time {
        font-size: 20rpx;
        color: $text-sub;
        margin-top: 8rpx;
        display: block;
      }
    }

    .info-bottom {
      .desc { font-size: 22rpx; color: $text-second; }

      .use-btn {
        height: 52rpx;
        padding: 0 24rpx;
        background: $primary-color;
        color: #fff;
        font-size: 22rpx;
        border-radius: 100rpx;
        margin: 0;
        display: flex;
        align-items: center;
        border: 0;
      }

      .status-stamp {
        width: 80rpx;
        height: 80rpx;
        opacity: 0.15;
        position: absolute;
        right: 20rpx;
        bottom: 10rpx;
      }
    }
  }

  &.used, &.expired {
    filter: grayscale(1);
    opacity: 0.7;
    .coupon-left { background: #f5f5f5; }
    .price .symbol, .price .val, .condition { color: #999; }
  }
}

.bottom-space { height: 40rpx; }
</style>



