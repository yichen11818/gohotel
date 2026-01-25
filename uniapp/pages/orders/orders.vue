<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">我的订单</text>
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

    <!-- 订单列表 -->
    <scroll-view class="order-list" scroll-y>
      <template v-if="orderList.length > 0">
        <view 
          v-for="order in orderList" 
          :key="order.id"
          class="order-card"
          @click="goToOrderDetail(order.id)"
        >
          <view class="order-header">
            <text class="order-number">订单号: {{ order.orderNo }}</text>
            <text class="order-status" :class="order.statusClass">{{ order.statusText }}</text>
          </view>
          
          <view class="order-content">
            <image class="hotel-image" :src="order.hotelImage" mode="aspectFill"></image>
            <view class="order-info">
              <text class="hotel-name">{{ order.hotelName }}</text>
              <text class="room-type">{{ order.roomType }}</text>
              <text class="date-info">{{ order.checkIn }} 至 {{ order.checkOut }}</text>
            </view>
          </view>
          
          <view class="order-footer">
            <view class="price-info">
              <text class="price-label">实付金额:</text>
              <text class="price-value">¥{{ order.totalPrice }}</text>
            </view>
            <view class="order-actions">
              <view 
                v-for="action in order.actions" 
                :key="action.type"
                class="action-btn"
                :class="action.type"
                @click.stop="handleAction(action.type, order)"
              >
                <text>{{ action.text }}</text>
              </view>
            </view>
          </view>
        </view>
      </template>
      
      <!-- 空状态 -->
      <TnEmpty v-else description="暂无订单" :show-image="true" />
      
      <view class="bottom-placeholder"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { booking } from '@/api/index.js'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnEmpty from '@/uni_modules/tuniaoui-vue3/components/empty/src/empty.vue'

// 标签页配置
const tabs = [
  { label: '全部', value: 'all' },
  { label: '待付款', value: 'pending' },
  { label: '待入住', value: 'confirmed' },
  { label: '已完成', value: 'completed' },
  { label: '已取消', value: 'cancelled' }
]

const currentTab = ref('all')
const orderList = ref([])
const loading = ref(false)

onLoad((options) => {
  if (options?.type) {
    currentTab.value = options.type
  }
  loadOrders()
})

// 切换标签
const switchTab = (value) => {
  currentTab.value = value
  loadOrders()
}

// 加载订单
const loadOrders = async () => {
  loading.value = true
  try {
    const params = {
      status: currentTab.value === 'all' ? undefined : currentTab.value,
      page: 1,
      pageSize: 20
    }
    const data = await booking.getBookingList(params)
    orderList.value = data || []
  } catch (error) {
    console.error('Failed to load orders:', error)
  } finally {
    loading.value = false
  }
}

// 返回
const goBack = () => {
  uni.navigateBack()
}

// 前往订单详情
const goToOrderDetail = (orderId) => {
  uni.navigateTo({
    url: `/pages/order-detail/order-detail?id=${orderId}`
  })
}

// 处理操作
const handleAction = (type, order) => {
  switch (type) {
    case 'pay':
      // 支付
      uni.showToast({ title: '前往支付', icon: 'none' })
      break
    case 'cancel':
      // 取消订单
      uni.showModal({
        title: '提示',
        content: '确定要取消订单吗？',
        success: (res) => {
          if (res.confirm) {
            uni.showToast({ title: '订单已取消', icon: 'success' })
          }
        }
      })
      break
    case 'delete':
      // 删除订单
      uni.showModal({
        title: '提示',
        content: '确定要删除订单吗？',
        success: (res) => {
          if (res.confirm) {
            uni.showToast({ title: '订单已删除', icon: 'success' })
          }
        }
      })
      break
  }
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
  animation: slideIn 0.3s ease;
}

@keyframes slideIn {
  from {
    width: 0;
  }
  to {
    width: 40rpx;
  }
}

/* 订单列表 */
.order-list {
  height: calc(100vh - 88rpx - 96rpx);
  padding: 24rpx;
}

.order-card {
  background: #fff;
  border-radius: 24rpx;
  padding: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  &:active {
    transform: scale(0.98);
  }
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
  padding-bottom: 20rpx;
  border-bottom: 1rpx solid #F0F0F0;
}

.order-number {
  font-size: 24rpx;
  color: #999;
}

.order-status {
  font-size: 26rpx;
  font-weight: 600;
  
  &.pending {
    color: #FF9500;
  }
  
  &.checkin {
    color: #5AC8FA;
  }
  
  &.completed {
    color: #34C759;
  }
  
  &.cancelled {
    color: #FF3B30;
  }
}

.order-content {
  display: flex;
  gap: 20rpx;
  margin-bottom: 20rpx;
}

.hotel-image {
  width: 160rpx;
  height: 120rpx;
  border-radius: 12rpx;
}

.order-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.hotel-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.room-type {
  font-size: 26rpx;
  color: #666;
}

.date-info {
  font-size: 24rpx;
  color: #999;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20rpx;
  border-top: 1rpx solid #F0F0F0;
}

.price-info {
  display: flex;
  align-items: baseline;
  gap: 8rpx;
}

.price-label {
  font-size: 24rpx;
  color: #666;
}

.price-value {
  font-size: 32rpx;
  font-weight: 700;
  color: #C29D71;
}

.order-actions {
  display: flex;
  gap: 16rpx;
}

.action-btn {
  padding: 12rpx 32rpx;
  border-radius: 100rpx;
  font-size: 26rpx;
  transition: all 0.3s ease;
  
  &.pay {
    background: linear-gradient(135deg, #D4B184 0%, #C29D71 100%);
    color: #fff;
  }
  
  &.cancel,
  &.delete {
    background: #F5F5F5;
    color: #666;
  }
  
  &:active {
    transform: scale(0.95);
  }
}

.bottom-placeholder {
  height: 40rpx;
}
</style>



