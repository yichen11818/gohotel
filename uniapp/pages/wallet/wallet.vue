<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">我的钱包</text>
        </view>
      </template>
    </TnNavbar>

    <scroll-view class="main-content" scroll-y>
      <!-- 余额卡片 -->
      <view class="balance-card">
        <view class="balance-header">
          <text class="balance-label">账户余额（元）</text>
          <view class="eye-btn" @click="toggleBalance">
            <TnIcon :name="showBalance ? 'eye' : 'eye-off'" color="#fff" size="32" />
          </view>
        </view>
        <view class="balance-amount">
          <text v-if="showBalance" class="amount">{{ balance }}</text>
          <text v-else class="amount">****</text>
        </view>
        <view class="balance-actions">
          <view class="balance-action" @click="handleRecharge">
            <TnIcon name="add-circle" color="#fff" size="40" />
            <text>充值</text>
          </view>
          <view class="balance-action" @click="handleWithdraw">
            <TnIcon name="minus-circle" color="#fff" size="40" />
            <text>提现</text>
          </view>
        </view>
      </view>

      <!-- 功能入口 -->
      <view class="function-section">
        <view class="function-item" @click="goToCoupons">
          <view class="function-icon coupon">
            <TnIcon name="ticket" color="#FF9500" size="44" />
          </view>
          <view class="function-info">
            <text class="function-name">优惠券</text>
            <text class="function-desc">{{ couponsCount }}张可用</text>
          </view>
          <TnIcon name="right" color="#ccc" size="28" />
        </view>
        
        <view class="function-item" @click="goToPoints">
          <view class="function-icon points">
            <TnIcon name="gift" color="#5AC8FA" size="44" />
          </view>
          <view class="function-info">
            <text class="function-name">积分</text>
            <text class="function-desc">{{ pointsCount }}积分</text>
          </view>
          <TnIcon name="right" color="#ccc" size="28" />
        </view>
        
        <view class="function-item" @click="goToBank">
          <view class="function-icon bank">
            <TnIcon name="credit-card" color="#34C759" size="44" />
          </view>
          <view class="function-info">
            <text class="function-name">银行卡</text>
            <text class="function-desc">管理银行卡</text>
          </view>
          <TnIcon name="right" color="#ccc" size="28" />
        </view>
      </view>

      <!-- 交易记录 -->
      <view class="transaction-section">
        <view class="section-header">
          <text class="section-title">交易记录</text>
          <view class="section-more" @click="goToAllRecords">
            <text class="more-text">查看全部</text>
            <TnIcon name="right" color="#999" size="28" />
          </view>
        </view>
        
        <view v-if="transactionList.length > 0" class="transaction-list">
          <view 
            v-for="item in transactionList" 
            :key="item.id"
            class="transaction-item"
          >
            <view class="transaction-icon" :class="item.type">
              <TnIcon :name="item.icon" :color="item.iconColor" size="40" />
            </view>
            <view class="transaction-info">
              <text class="transaction-title">{{ item.title }}</text>
              <text class="transaction-time">{{ item.time }}</text>
            </view>
            <text class="transaction-amount" :class="item.type">
              {{ item.type === 'income' ? '+' : '-' }}¥{{ item.amount }}
            </text>
          </view>
        </view>
        
        <TnEmpty v-else description="暂无交易记录" :show-image="true" />
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { user } from '@/api/index.js'
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnEmpty from '@/uni_modules/tuniaoui-vue3/components/empty/src/empty.vue'

const showBalance = ref(true)
const balance = ref('0.00')
const couponsCount = ref(0)
const pointsCount = ref(0)
const loading = ref(false)

const transactionList = ref([])

onLoad(() => {
  loadWalletData()
})

const loadWalletData = async () => {
  loading.value = true
  try {
    const [infoRes, pointsRes, couponsRes] = await Promise.all([
      user.getUserInfo(),
      user.getUserPoints(),
      user.getCoupons({ status: 'unused' })
    ])
    
    if (infoRes) {
      balance.value = infoRes.balance || '0.00'
    }
    if (pointsRes) {
      pointsCount.value = pointsRes.points || 0
    }
    if (couponsRes) {
      couponsCount.value = couponsRes.length || 0
    }
    
    // 假设有交易记录API，如果没有则保持空
    // const records = await user.getTransactionRecords()
    // transactionList.value = records || []
  } catch (error) {
    console.error('Failed to load wallet data:', error)
  } finally {
    loading.value = false
  }
}

const toggleBalance = () => {
  showBalance.value = !showBalance.value
}

const goBack = () => {
  uni.navigateBack()
}

const handleRecharge = () => {
  uni.navigateTo({
    url: '/pages/recharge/recharge'
  })
}

const handleWithdraw = () => {
  uni.navigateTo({
    url: '/pages/withdraw/withdraw'
  })
}

const goToCoupons = () => {
  uni.navigateTo({
    url: '/pages/coupons/coupons'
  })
}

const goToPoints = () => {
  uni.navigateTo({
    url: '/pages/points/points'
  })
}

const goToBank = () => {
  uni.navigateTo({
    url: '/pages/bank-card/bank-card'
  })
}

const goToAllRecords = () => {
  uni.navigateTo({
    url: '/pages/transaction-records/transaction-records'
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

/* 余额卡片 */
.balance-card {
  padding: 48rpx 40rpx;
  background: linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%);
  border-radius: 32rpx;
  box-shadow: 0 16rpx 48rpx rgba(194, 157, 113, 0.3);
  margin-bottom: 24rpx;
  animation: fadeInUp 0.6s ease-out;
}

.balance-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.balance-label {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
}

.eye-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56rpx;
  height: 56rpx;
  
  &:active {
    opacity: 0.7;
  }
}

.balance-amount {
  margin-bottom: 40rpx;
}

.amount {
  font-size: 80rpx;
  font-weight: 700;
  color: #fff;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

.balance-actions {
  display: flex;
  gap: 32rpx;
}

.balance-action {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12rpx;
  padding: 24rpx;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 20rpx;
  backdrop-filter: blur(10px);
  
  text {
    font-size: 26rpx;
    color: #fff;
    font-weight: 500;
  }
  
  &:active {
    background: rgba(255, 255, 255, 0.3);
  }
}

/* 功能区 */
.function-section {
  background: #fff;
  border-radius: 24rpx;
  padding: 0 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.function-item {
  display: flex;
  align-items: center;
  padding: 32rpx 0;
  border-bottom: 1rpx solid #F0F0F0;
  
  &:last-child {
    border-bottom: none;
  }
  
  &:active {
    opacity: 0.7;
  }
}

.function-icon {
  width: 72rpx;
  height: 72rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  margin-right: 24rpx;
  
  &.coupon {
    background: linear-gradient(135deg, rgba(255, 149, 0, 0.15) 0%, rgba(255, 149, 0, 0.08) 100%);
  }
  
  &.points {
    background: linear-gradient(135deg, rgba(90, 200, 250, 0.15) 0%, rgba(90, 200, 250, 0.08) 100%);
  }
  
  &.bank {
    background: linear-gradient(135deg, rgba(52, 199, 89, 0.15) 0%, rgba(52, 199, 89, 0.08) 100%);
  }
}

.function-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.function-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.function-desc {
  font-size: 24rpx;
  color: #999;
}

/* 交易记录 */
.transaction-section {
  background: #fff;
  border-radius: 24rpx;
  padding: 32rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #333;
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

.transaction-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.transaction-item {
  display: flex;
  align-items: center;
  padding: 24rpx;
  background: #FAFAFA;
  border-radius: 16rpx;
  
  &:active {
    background: #F5F5F5;
  }
}

.transaction-icon {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  margin-right: 20rpx;
  
  &.income {
    background: linear-gradient(135deg, rgba(52, 199, 89, 0.15) 0%, rgba(52, 199, 89, 0.08) 100%);
  }
  
  &.expense {
    background: linear-gradient(135deg, rgba(255, 59, 48, 0.15) 0%, rgba(255, 59, 48, 0.08) 100%);
  }
}

.transaction-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.transaction-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.transaction-time {
  font-size: 24rpx;
  color: #999;
}

.transaction-amount {
  font-size: 32rpx;
  font-weight: 700;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  
  &.income {
    color: #34C759;
  }
  
  &.expense {
    color: #FF3B30;
  }
}

.bottom-placeholder {
  height: 40rpx;
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
</style>



