<template>
  <view class="page">
    <view class="tabs">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @click="switchTab(tab.value)"
      >
        <text>{{ tab.label }}</text>
      </view>
    </view>

    <scroll-view class="list" scroll-y>
      <view v-if="loading" class="placeholder">正在加载订单...</view>
      <view v-else-if="!orderList.length" class="placeholder">暂无订单</view>

      <view v-for="order in orderList" :key="order.id" class="order-card" @click="goToOrderDetail(order.id)">
        <view class="order-header">
          <text class="order-no">订单号 {{ order.orderNo }}</text>
          <text class="status" :class="order.statusClass">{{ order.statusText }}</text>
        </view>
        <view class="order-body">
          <image class="order-image" :src="order.roomImage" mode="aspectFill" />
          <view class="order-info">
            <text class="room-name">{{ order.roomName }}</text>
            <text class="room-meta">{{ order.checkIn }} 至 {{ order.checkOut }}</text>
            <text class="room-meta">{{ order.totalDays }} 晚 · {{ order.roomNumber || '待分配房号' }}</text>
          </view>
        </view>
        <view class="order-footer">
          <text class="amount">¥{{ order.totalPrice }}</text>
          <button v-if="order.canCancel" class="ghost-btn" size="mini" @click.stop="handleCancel(order)">取消订单</button>
        </view>
      </view>
      <view class="bottom-space"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { booking } from '@/api/index.js'

const tabs = [
  { label: '全部', value: 'all' },
  { label: '待确认', value: 'pending' },
  { label: '待入住', value: 'confirmed' },
  { label: '入住中', value: 'checkin' },
  { label: '已完成', value: 'checkout' },
  { label: '已取消', value: 'cancelled' },
]

const currentTab = ref('all')
const orderList = ref([])
const loading = ref(false)

const loadOrders = async () => {
  loading.value = true
  try {
    orderList.value = await booking.getBookingList({
      status: currentTab.value,
      page: 1,
      pageSize: 50,
    })
  } catch (error) {
    console.error('load orders failed:', error)
  } finally {
    loading.value = false
  }
}

const switchTab = (value) => {
  currentTab.value = value
  loadOrders()
}

const goToOrderDetail = (id) => {
  uni.navigateTo({ url: `/pages/order-detail/order-detail?id=${id}` })
}

const handleCancel = (order) => {
  uni.showModal({
    title: '取消订单',
    content: '确定要取消这笔订单吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        uni.showLoading({ title: '处理中...' })
        await booking.cancelBooking(order.id, '用户主动取消')
        uni.hideLoading()
        uni.showToast({ title: '订单已取消', icon: 'success' })
        loadOrders()
      } catch (_error) {
        uni.hideLoading()
      }
    },
  })
}

onLoad((options) => {
  if (options?.type) {
    currentTab.value = options.type
  }
})

onShow(() => {
  loadOrders()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
}

.tabs {
  display: flex;
  gap: 16rpx;
  padding: 28rpx 24rpx 16rpx;
  overflow-x: auto;
  white-space: nowrap;
  background: #fff;
}

.tab-item {
  flex-shrink: 0;
  padding: 12rpx 24rpx;
  font-size: 24rpx;
  color: #6b7280;
  background: #f3f4f6;
  border-radius: 999rpx;
}

.tab-item.active {
  color: #fff;
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
}

.list {
  height: calc(100vh - 112rpx);
  padding: 24rpx;
  box-sizing: border-box;
}

.placeholder {
  padding: 80rpx 24rpx;
  text-align: center;
  color: #6b7280;
  font-size: 26rpx;
}

.order-card {
  padding: 24rpx;
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.order-header,
.order-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16rpx;
}

.order-no,
.room-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #111827;
}

.status {
  font-size: 24rpx;
}

.status.pending,
.status.confirmed,
.status.checkin {
  color: #b7791f;
}

.status.checkout {
  color: #059669;
}

.status.cancelled {
  color: #dc2626;
}

.order-body {
  display: flex;
  gap: 20rpx;
  margin-top: 20rpx;
}

.order-image {
  width: 180rpx;
  height: 140rpx;
  border-radius: 16rpx;
  background: #f3f4f6;
}

.order-info {
  flex: 1;
}

.room-meta,
.amount {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #4b5563;
}

.amount {
  font-size: 30rpx;
  color: #b7791f;
  font-weight: 700;
}

.ghost-btn {
  margin: 0;
  border-radius: 999rpx;
}

.bottom-space {
  height: 24rpx;
}
</style>
