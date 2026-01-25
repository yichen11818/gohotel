<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">我的积分</text>
        </view>
      </template>
    </TnNavbar>

    <scroll-view class="main-content" scroll-y>
      <!-- 积分卡片 -->
      <view class="points-card">
        <view class="points-header">
          <text class="points-label">当前积分</text>
          <view class="help-btn" @click="showPointsRule">
            <TnIcon name="question-circle" color="#fff" size="32" />
          </view>
        </view>
        <view class="points-amount">
          <text class="amount">{{ points }}</text>
          <text class="unit">分</text>
        </view>
        <view class="points-desc">
          <text>积分可用于兑换优惠券和礼品</text>
        </view>
      </view>

      <!-- 快捷入口 -->
      <view class="quick-actions">
        <view class="action-item" @click="goToMall">
          <view class="action-icon mall">
            <TnIcon name="gift" color="#FF9500" size="44" />
          </view>
          <text class="action-text">积分商城</text>
        </view>
        <view class="action-item" @click="goToTasks">
          <view class="action-icon task">
            <TnIcon name="list-check" color="#5AC8FA" size="44" />
          </view>
          <text class="action-text">赚积分</text>
        </view>
        <view class="action-item" @click="showPointsRule">
          <view class="action-icon rule">
            <TnIcon name="file-text" color="#34C759" size="44" />
          </view>
          <text class="action-text">积分规则</text>
        </view>
      </view>

      <!-- 积分记录 -->
      <view class="records-section">
        <view class="section-header">
          <text class="section-title">积分记录</text>
        </view>
        
        <view v-if="recordList.length > 0" class="record-list">
          <view 
            v-for="record in recordList" 
            :key="record.id"
            class="record-item"
          >
            <view class="record-left">
              <view class="record-icon" :class="record.type">
                <TnIcon :name="record.icon" :color="record.iconColor" size="36" />
              </view>
              <view class="record-info">
                <text class="record-title">{{ record.title }}</text>
                <text class="record-time">{{ record.time }}</text>
              </view>
            </view>
            <text class="record-points" :class="record.type">
              {{ record.type === 'add' ? '+' : '-' }}{{ record.points }}
            </text>
          </view>
        </view>
        
        <TnEmpty v-else description="暂无积分记录" :show-image="true" />
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

const points = ref(0)
const loading = ref(false)

const recordList = ref([])

onLoad(() => {
  loadPointsData()
})

const loadPointsData = async () => {
  loading.value = true
  try {
    const data = await user.getUserPoints()
    if (data) {
      points.value = data.points || 0
      // 记录列表也应从API获取
      // const records = await user.getPointsRecords()
      // recordList.value = records || []
    }
  } catch (error) {
    console.error('Failed to load points data:', error)
  } finally {
    loading.value = false
  }
}

const goBack = () => {
  uni.navigateBack()
}

const showPointsRule = () => {
  uni.showModal({
    title: '积分规则',
    content: '1. 完成订单可获得积分\n2. 每日签到可获得积分\n3. 积分可兑换优惠券和礼品\n4. 积分有效期为1年',
    showCancel: false
  })
}

const goToMall = () => {
  uni.showToast({
    title: '积分商城开发中',
    icon: 'none'
  })
}

const goToTasks = () => {
  uni.showToast({
    title: '任务中心开发中',
    icon: 'none'
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

/* 积分卡片 */
.points-card {
  padding: 48rpx 40rpx;
  background: linear-gradient(135deg, #5AC8FA 0%, #4A90E2 100%);
  border-radius: 32rpx;
  box-shadow: 0 16rpx 48rpx rgba(90, 200, 250, 0.3);
  margin-bottom: 24rpx;
  animation: fadeInUp 0.6s ease-out;
}

.points-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.points-label {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
}

.help-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 56rpx;
  height: 56rpx;
  
  &:active {
    opacity: 0.7;
  }
}

.points-amount {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
  margin-bottom: 16rpx;
}

.amount {
  font-size: 96rpx;
  font-weight: 700;
  color: #fff;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  line-height: 1;
}

.unit {
  font-size: 32rpx;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.points-desc {
  text {
    font-size: 24rpx;
    color: rgba(255, 255, 255, 0.8);
  }
}

/* 快捷入口 */
.quick-actions {
  display: flex;
  justify-content: space-around;
  padding: 40rpx 32rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.action-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
  
  &:active {
    opacity: 0.7;
  }
}

.action-icon {
  width: 88rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  
  &.mall {
    background: linear-gradient(135deg, rgba(255, 149, 0, 0.15) 0%, rgba(255, 149, 0, 0.08) 100%);
  }
  
  &.task {
    background: linear-gradient(135deg, rgba(90, 200, 250, 0.15) 0%, rgba(90, 200, 250, 0.08) 100%);
  }
  
  &.rule {
    background: linear-gradient(135deg, rgba(52, 199, 89, 0.15) 0%, rgba(52, 199, 89, 0.08) 100%);
  }
}

.action-text {
  font-size: 26rpx;
  color: #666;
  font-weight: 500;
}

/* 积分记录 */
.records-section {
  background: #fff;
  border-radius: 24rpx;
  padding: 32rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
}

.section-header {
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #333;
}

.record-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.record-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx;
  background: #FAFAFA;
  border-radius: 16rpx;
}

.record-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
  flex: 1;
}

.record-icon {
  width: 64rpx;
  height: 64rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  
  &.add {
    background: linear-gradient(135deg, rgba(52, 199, 89, 0.15) 0%, rgba(52, 199, 89, 0.08) 100%);
  }
  
  &.reduce {
    background: linear-gradient(135deg, rgba(255, 149, 0, 0.15) 0%, rgba(255, 149, 0, 0.08) 100%);
  }
}

.record-info {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.record-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.record-time {
  font-size: 24rpx;
  color: #999;
}

.record-points {
  font-size: 32rpx;
  font-weight: 700;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  
  &.add {
    color: #34C759;
  }
  
  &.reduce {
    color: #FF9500;
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



