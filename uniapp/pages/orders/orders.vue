<template>
  <view class="page">
    <!-- 顶部状态栏占位 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>

    <!-- 顶部导航 -->
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">我的订单</text>
      <view class="placeholder-view"></view>
    </view>

    <!-- 标签页 -->
    <view class="tabs-box">
      <scroll-view class="tabs-scroll" scroll-x :show-scrollbar="false">
        <view class="tabs-content">
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
      </scroll-view>
    </view>

    <scroll-view class="list-container" scroll-y @scrolltolower="onReachBottom">
      <view v-if="loading && !orderList.length" class="state-box">
        <view class="loading-icon"></view>
        <text>加载中...</text>
      </view>

      <view v-else-if="!orderList.length" class="state-box">
        <image class="empty-img" src="/static/icons/order.png" mode="aspectFit" style="opacity: 0.2; filter: grayscale(1);" />
        <text>暂无相关订单</text>
      </view>

      <view class="order-list" v-else>
        <view
          v-for="order in orderList"
          :key="order.id"
          class="order-card premium-card"
          @click="goToOrderDetail(order.id)"
        >
          <view class="card-header flex-between">
            <view class="order-no flex-center">
              <image class="icon-order" src="/static/icons/order.png" mode="aspectFit" />
              <text>订单号: {{ order.orderNo }}</text>
            </view>
            <text class="status-text" :class="order.status">{{ order.statusText }}</text>
          </view>

          <view class="card-body">
            <image class="room-image" :src="order.roomImage" mode="aspectFill" />
            <view class="room-info">
              <text class="room-name text-ellipsis">{{ order.roomName }}</text>
              <view class="date-info">
                <text>{{ formatDateDisplay(order.checkIn) }} - {{ formatDateDisplay(order.checkOut) }}</text>
                <text class="nights">共{{ order.totalDays }}晚</text>
              </view>
              <text class="room-num" v-if="order.roomNumber">房号: {{ order.roomNumber }}</text>
              <text class="room-num" v-else>待分配房号</text>
            </view>
          </view>

          <view class="card-footer flex-between">
            <view class="price-box flex-center">
              <text class="label">实付款</text>
              <text class="symbol">¥</text>
              <text class="val">{{ order.totalPrice }}</text>
            </view>
            <view class="actions flex-center">
              <button v-if="order.canCancel" class="btn ghost-btn" @click.stop="handleCancel(order)">取消订单</button>
              <button class="btn primary-btn">查看详情</button>
            </view>
          </view>
        </view>
        <view class="load-more" v-if="orderList.length > 5">
          <text>已经到底啦</text>
        </view>
      </view>
      <view class="bottom-space"></view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { booking } from '@/api/index.js'

const statusBarHeight = ref(44)
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
    const data = await booking.getBookingList({
      status: currentTab.value,
      page: 1,
      pageSize: 50,
    })
    orderList.value = data || []
  } catch (error) {
    console.error('load orders failed:', error)
  } finally {
    loading.value = false
  }
}

const switchTab = (value) => {
  if (currentTab.value === value) return
  currentTab.value = value
  orderList.value = []
  loadOrders()
}

const goToOrderDetail = (id) => {
  uni.navigateTo({ url: `/pages/order-detail/order-detail?id=${id}` })
}

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({ url: '/pages/index/index' })
  }
}

const handleCancel = (order) => {
  uni.showModal({
    title: '温馨提示',
    content: '确认要取消这笔订单吗？',
    confirmColor: '#C29D71',
    success: async (res) => {
      if (!res.confirm) return
      try {
        uni.showLoading({ title: '正在处理' })
        await booking.cancelBooking(order.id, '用户主动取消')
        uni.hideLoading()
        uni.showToast({ title: '已取消订单', icon: 'success' })
        loadOrders()
      } catch (_error) {
        uni.hideLoading()
      }
    },
  })
}

const formatDateDisplay = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getMonth() + 1}月${d.getDate()}日`
}

const onReachBottom = () => {
  // 分页逻辑
}

onLoad((options) => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
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
    font-size: 36rpx;
    font-weight: 700;
    color: $text-main;
  }

  .placeholder-view { width: 64rpx; }
}

/* Tabs */
.tabs-box {
  background: #fff;
  border-bottom: 1rpx solid #f8f8f8;

  .tabs-scroll {
    width: 100%;
    white-space: nowrap;
  }

  .tabs-content {
    display: inline-flex;
    padding: 0 20rpx;
  }

  .tab-item {
    padding: 24rpx 30rpx;
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
        transform: scale(1.05);
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
    width: 160rpx;
    height: 160rpx;
    margin-bottom: 30rpx;
  }
}

.order-list {
  padding: 30rpx;
}

.order-card {
  padding: 32rpx;
  margin-bottom: 30rpx;

  .card-header {
    padding-bottom: 24rpx;
    border-bottom: 1rpx solid #f8f8f8;

    .order-no {
      font-size: 24rpx;
      color: $text-sub;
      .icon-order {
        width: 28rpx;
        height: 28rpx;
        margin-right: 8rpx;
        opacity: 0.5;
      }
    }

    .status-text {
      font-size: 26rpx;
      font-weight: 600;

      &.pending, &.confirmed, &.checkin { color: $primary-color; }
      &.checkout { color: #07C160; }
      &.cancelled { color: $text-sub; text-decoration: line-through; }
    }
  }

  .card-body {
    display: flex;
    padding: 30rpx 0;

    .room-image {
      width: 160rpx;
      height: 160rpx;
      border-radius: 12rpx;
      background: #f5f5f5;
      margin-right: 24rpx;
    }

    .room-info {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: space-between;

      .room-name {
        font-size: 30rpx;
        font-weight: 700;
        color: $text-main;
      }

      .date-info {
        font-size: 24rpx;
        color: $text-second;
        .nights {
          margin-left: 16rpx;
          color: $primary-color;
          font-weight: 500;
        }
      }

      .room-num {
        font-size: 22rpx;
        color: $text-sub;
        background: $bg-color;
        padding: 4rpx 16rpx;
        border-radius: 4rpx;
        align-self: flex-start;
      }
    }
  }

  .card-footer {
    padding-top: 24rpx;
    border-top: 1rpx solid #f8f8f8;

    .price-box {
      .label { font-size: 22rpx; color: $text-sub; margin-right: 8rpx; }
      .symbol { font-size: 22rpx; color: #E64340; font-weight: 700; }
      .val { font-size: 36rpx; color: #E64340; font-weight: 700; }
    }

    .btn {
      height: 64rpx;
      padding: 0 24rpx;
      font-size: 24rpx;
      border-radius: 100rpx;
      margin-left: 16rpx;
      display: flex;
      align-items: center;
      justify-content: center;

      &.ghost-btn {
        background: #fff;
        border: 1rpx solid #eee;
        color: $text-second;
      }

      &.primary-btn {
        background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
        color: #fff;
        border: 0;
      }
    }
  }
}

.load-more {
  text-align: center;
  padding: 20rpx 0;
  font-size: 22rpx;
  color: $text-sub;
}

.bottom-space {
  height: 40rpx;
}
</style>
