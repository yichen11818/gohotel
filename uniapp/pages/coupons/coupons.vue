<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">优惠券</text>
        </view>
      </template>
    </TnNavbar>

    <!-- 标签页 -->
    <view class="tabs">
      <view 
        v-for="(tab, index) in tabs" 
        :key="index"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @click="switchTab(tab.value)"
      >
        <text class="tab-text">{{ tab.label }}</text>
        <view v-if="currentTab === tab.value" class="tab-indicator"></view>
      </view>
    </view>

    <scroll-view class="coupon-list" scroll-y>
      <template v-if="couponList.length > 0">
        <view 
          v-for="coupon in couponList" 
          :key="coupon.id"
          class="coupon-card"
          :class="{ expired: coupon.status === 'expired', used: coupon.status === 'used' }"
        >
          <view class="coupon-left">
            <view class="coupon-amount">
              <text class="amount-symbol">¥</text>
              <text class="amount-value">{{ coupon.amount }}</text>
            </view>
            <text class="coupon-condition">满{{ coupon.condition }}元可用</text>
          </view>
          
          <view class="coupon-divider">
            <view class="circle top"></view>
            <view class="dash-line"></view>
            <view class="circle bottom"></view>
          </view>
          
          <view class="coupon-right">
            <text class="coupon-name">{{ coupon.name }}</text>
            <text class="coupon-time">{{ coupon.validTime }}</text>
            <view v-if="currentTab === 'unused'" class="use-btn" @click.stop="useCoupon(coupon)">
              <text>立即使用</text>
            </view>
            <view v-else-if="currentTab === 'used'" class="status-tag used">
              <text>已使用</text>
            </view>
            <view v-else class="status-tag expired">
              <text>已过期</text>
            </view>
          </view>
        </view>
      </template>
      
      <!-- 空状态 -->
      <TnEmpty v-else description="暂无优惠券" :show-image="true" />
      
      <view class="bottom-placeholder"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { user } from '@/api/index.js'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnEmpty from '@/uni_modules/tuniaoui-vue3/components/empty/src/empty.vue'

const tabs = [
  { label: '可使用', value: 'unused' },
  { label: '已使用', value: 'used' },
  { label: '已过期', value: 'expired' }
]

const currentTab = ref('unused')
const couponList = ref([])
const loading = ref(false)

onLoad(() => {
  loadCoupons()
})

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
  currentTab.value = value
  loadCoupons()
}

const goBack = () => {
  uni.navigateBack()
}

const useCoupon = (coupon) => {
  uni.showModal({
    title: '使用优惠券',
    content: `确定使用${coupon.amount}元优惠券吗？`,
    success: (res) => {
      if (res.confirm) {
        // 跳转到预订页面
        uni.navigateTo({
          url: `/pages/hotel/hotel?couponId=${coupon.id}`
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

/* 标签页 */
.tabs {
  display: flex;
  background: #fff;
  padding: 0 24rpx;
  position: sticky;
  top: 88rpx;
  z-index: 10;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24rpx 0;
  position: relative;
  
  &:active {
    opacity: 0.7;
  }
}

.tab-text {
  font-size: 28rpx;
  color: #666;
  font-weight: 500;
  transition: all 0.3s ease;
}

.tab-item.active .tab-text {
  color: #C29D71;
  font-weight: 700;
  font-size: 32rpx;
}

.tab-indicator {
  position: absolute;
  bottom: 0;
  width: 40rpx;
  height: 6rpx;
  background: linear-gradient(90deg, #D4B184 0%, #C29D71 100%);
  border-radius: 3rpx;
}

/* 优惠券列表 */
.coupon-list {
  height: calc(100vh - 88rpx - 96rpx);
  padding: 24rpx;
}

.coupon-card {
  display: flex;
  height: 220rpx;
  margin-bottom: 24rpx;
  background: linear-gradient(135deg, #FFFFFF 0%, #FEFEFE 100%);
  border-radius: 20rpx;
  box-shadow: 0 8rpx 24rpx rgba(194, 157, 113, 0.15);
  overflow: hidden;
  position: relative;
  transition: all 0.3s ease;
  
  &:active {
    transform: scale(0.98);
  }
  
  &.expired,
  &.used {
    opacity: 0.6;
    filter: grayscale(0.8);
  }
}

.coupon-left {
  width: 220rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  background: linear-gradient(135deg, #D4B184 0%, #C29D71 100%);
  position: relative;
}

.coupon-amount {
  display: flex;
  align-items: baseline;
  gap: 4rpx;
}

.amount-symbol {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
}

.amount-value {
  font-size: 72rpx;
  font-weight: 800;
  color: #fff;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

.coupon-condition {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.9);
}

.coupon-divider {
  width: 1rpx;
  position: relative;
  background: transparent;
  
  .circle {
    position: absolute;
    width: 24rpx;
    height: 24rpx;
    background: #F5F5F5;
    border-radius: 50%;
    left: 50%;
    transform: translateX(-50%);
    
    &.top {
      top: -12rpx;
    }
    
    &.bottom {
      bottom: -12rpx;
    }
  }
  
  .dash-line {
    position: absolute;
    top: 24rpx;
    bottom: 24rpx;
    left: 50%;
    transform: translateX(-50%);
    width: 2rpx;
    background-image: linear-gradient(to bottom, #E5E5E5 0%, #E5E5E5 50%, transparent 50%, transparent 100%);
    background-size: 2rpx 12rpx;
    background-repeat: repeat-y;
  }
}

.coupon-right {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 32rpx;
  gap: 12rpx;
}

.coupon-name {
  font-size: 32rpx;
  font-weight: 700;
  color: #333;
}

.coupon-time {
  font-size: 24rpx;
  color: #999;
}

.use-btn {
  align-self: flex-start;
  margin-top: 12rpx;
  padding: 12rpx 32rpx;
  background: linear-gradient(135deg, #D4B184 0%, #C29D71 100%);
  border-radius: 100rpx;
  box-shadow: 0 8rpx 16rpx rgba(194, 157, 113, 0.3);
  
  text {
    font-size: 24rpx;
    color: #fff;
    font-weight: 600;
  }
  
  &:active {
    transform: scale(0.95);
  }
}

.status-tag {
  align-self: flex-start;
  margin-top: 12rpx;
  padding: 8rpx 24rpx;
  border-radius: 100rpx;
  
  text {
    font-size: 22rpx;
    font-weight: 500;
  }
  
  &.used {
    background: rgba(153, 153, 153, 0.1);
    
    text {
      color: #999;
    }
  }
  
  &.expired {
    background: rgba(255, 59, 48, 0.1);
    
    text {
      color: #FF3B30;
    }
  }
}

.bottom-placeholder {
  height: 40rpx;
}
</style>



