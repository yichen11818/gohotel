<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">订单详情</text>
      <view class="placeholder-view"></view>
    </view>

    <view class="content" v-if="orderDetail.id">
      <!-- 状态卡片 -->
      <view class="status-banner">
        <view class="status-content">
          <view class="status-header flex-between">
            <text class="status-text">{{ orderDetail.statusText }}</text>
            <image class="status-icon" :src="getStatusIcon(orderDetail.status)" mode="aspectFit" />
          </view>
          <text class="order-no">订单号: {{ orderDetail.orderNo }}</text>
        </view>
      </view>

      <!-- 房型卡片 -->
      <view class="section room-card premium-card" @click="viewRoom">
        <image class="room-image" :src="orderDetail.roomImage" mode="aspectFill" />
        <view class="room-info">
          <text class="room-name">{{ orderDetail.roomName }}</text>
          <view class="date-info">
            <text>{{ formatDateDisplay(orderDetail.checkIn) }} - {{ formatDateDisplay(orderDetail.checkOut) }}</text>
            <text class="nights">共{{ orderDetail.totalDays }}晚</text>
          </view>
          <text class="room-num">{{ orderDetail.roomNumber || '待分配房号' }}</text>
        </view>
        <image class="icon-arrow" src="/static/icons/arrow-right.png" mode="aspectFit" />
      </view>

      <!-- 入住人信息 -->
      <view class="section info-section premium-card">
        <view class="section-header">
          <text class="title">入住人信息</text>
        </view>
        <view class="info-list">
          <view class="info-item flex-between">
            <text class="label">姓名</text>
            <text class="val">{{ orderDetail.guestName || '未填写' }}</text>
          </view>
          <view class="info-item flex-between">
            <text class="label">手机号</text>
            <text class="val">{{ orderDetail.guestPhone || '未填写' }}</text>
          </view>
          <view class="info-item flex-between" v-if="orderDetail.guestIdCard">
            <text class="label">身份证号</text>
            <text class="val">{{ formatIdCard(orderDetail.guestIdCard) }}</text>
          </view>
          <view class="info-item flex-between" v-if="orderDetail.specialRequest">
            <text class="label">特殊需求</text>
            <text class="val multiline">{{ orderDetail.specialRequest }}</text>
          </view>
        </view>
      </view>

      <!-- 订单信息 -->
      <view class="section info-section premium-card">
        <view class="section-header">
          <text class="title">订单信息</text>
        </view>
        <view class="info-list">
          <view class="info-item flex-between">
            <text class="label">下单时间</text>
            <text class="val">{{ orderDetail.createdAt || '-' }}</text>
          </view>
          <view class="info-item flex-between">
            <text class="label">支付状态</text>
            <text class="val">{{ paymentStatusText }}</text>
          </view>
          <view class="info-item flex-between amount-row">
            <text class="label">订单总额</text>
            <view class="price">
              <text class="symbol">¥</text>
              <text class="val">{{ orderDetail.totalPrice }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- 酒店信息 -->
      <view class="section hotel-section premium-card" v-if="hotelInfo.phone">
        <view class="section-header flex-between">
          <text class="title">酒店信息</text>
          <view class="call-btn flex-center" @click="makeCall(hotelInfo.phone)">
            <image class="icon" src="/static/icons/phone.png" mode="aspectFit" />
            <text>联系前台</text>
          </view>
        </view>
        <view class="hotel-info">
          <text class="address">{{ hotelInfo.address }}</text>
        </view>
      </view>
    </view>

    <view v-else class="state-box">
      <view class="loading-icon"></view>
      <text>正在获取订单详情...</text>
    </view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar safe-area-inset-bottom" v-if="orderDetail.id">
      <button v-if="orderDetail.canCancel" class="btn ghost-btn" @click="cancelOrder">取消订单</button>
      <button class="btn primary-btn" @click="viewRoom">查看房型</button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { booking, hotel } from '@/api/index.js'

const orderId = ref('')
const statusBarHeight = ref(44)
const orderDetail = ref({})

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({ url: '/pages/index/index' })
  }
}
const hotelInfo = ref({
  phone: '',
  address: '欢迎入住七天酒店',
})

const paymentStatusText = computed(() => {
  const status = orderDetail.value.paymentStatus
  if (status === 'paid') return '已支付'
  if (status === 'refunded') return '已退款'
  return '未支付'
})

const loadOrderDetail = async () => {
  try {
    const [detail, hotelDetail] = await Promise.all([
      booking.getBookingDetail(orderId.value),
      hotel.getHotelDetail(),
    ])
    orderDetail.value = detail
    hotelInfo.value = hotelDetail
  } catch (error) {
    console.error('load order detail failed:', error)
  }
}

const cancelOrder = () => {
  uni.showModal({
    title: '温馨提示',
    content: '确认要取消这笔订单吗？',
    confirmColor: '#C29D71',
    success: async (res) => {
      if (!res.confirm) return
      try {
        uni.showLoading({ title: '正在处理' })
        await booking.cancelBooking(orderDetail.value.id, '用户主动取消')
        uni.hideLoading()
        uni.showToast({ title: '订单已取消', icon: 'success' })
        loadOrderDetail()
      } catch (_error) {
        uni.hideLoading()
      }
    },
  })
}

const viewRoom = () => {
  if (!orderDetail.value.roomId) return
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${orderDetail.value.roomId}&checkIn=${orderDetail.value.checkIn}&checkOut=${orderDetail.value.checkOut}`,
  })
}

const getStatusIcon = (status) => {
  const icons = {
    pending: '/static/icons/order-pay.png',
    confirmed: '/static/icons/order-stay.png',
    checkout: '/static/icons/order-done.png',
    cancelled: '/static/icons/order-refund.png',
  }
  return icons[status] || '/static/icons/order.png'
}

const formatDateDisplay = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getMonth() + 1}月${d.getDate()}日`
}

const formatIdCard = (id) => {
  if (!id) return ''
  return id.replace(/^(.{4}).+(.{4})$/, '$1**********$2')
}

const makeCall = (phone) => {
  uni.makePhoneCall({ phoneNumber: phone })
}

onLoad((options) => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
  if (options?.id) {
    orderId.value = String(options.id)
    loadOrderDetail()
  }
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: $bg-color;
  padding-bottom: 180rpx;
}

/* 顶部状态条 */
.status-banner {
  height: 320rpx;
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  padding: 60rpx 40rpx;
  box-sizing: border-box;
  color: #fff;

  .status-header {
    margin-bottom: 16rpx;
    .status-text {
      font-size: 44rpx;
      font-weight: 700;
    }
    .status-icon {
      width: 80rpx;
      height: 80rpx;
      opacity: 0.9;
    }
  }

  .order-no {
    font-size: 24rpx;
    opacity: 0.8;
  }
}

.content {
  padding: 0 30rpx;
  margin-top: -60rpx;
}

/* 房型卡片 */
.room-card {
  display: flex;
  align-items: center;
  padding: 30rpx;
  margin-bottom: 30rpx;

  .room-image {
    width: 160rpx;
    height: 160rpx;
    border-radius: 12rpx;
    background: #f5f5f5;
    margin-right: 24rpx;
  }

  .room-info {
    flex: 1;
    .room-name {
      font-size: 32rpx;
      font-weight: 700;
      color: $text-main;
      display: block;
      margin-bottom: 12rpx;
    }
    .date-info {
      font-size: 24rpx;
      color: $text-second;
      .nights { margin-left: 12rpx; color: $primary-color; }
    }
    .room-num {
      margin-top: 12rpx;
      font-size: 22rpx;
      color: $text-sub;
      background: $bg-color;
      padding: 2rpx 12rpx;
      border-radius: 4rpx;
      display: inline-block;
    }
  }

  .icon-arrow {
    width: 32rpx;
    height: 32rpx;
    opacity: 0.3;
  }
}

/* 信息列表 */
.info-section {
  padding: 40rpx;
  margin-bottom: 30rpx;

  .section-header {
    margin-bottom: 30rpx;
    .title {
      font-size: 32rpx;
      font-weight: 700;
      color: $text-main;
    }
  }

  .info-item {
    padding: 20rpx 0;
    font-size: 28rpx;

    .label { color: $text-sub; }
    .val {
      color: $text-main;
      font-weight: 500;
      &.multiline { max-width: 400rpx; text-align: right; }
    }

    &.amount-row {
      margin-top: 10rpx;
      padding-top: 30rpx;
      border-top: 1rpx solid #f8f8f8;
      .price {
        .symbol { font-size: 24rpx; color: #E64340; font-weight: 700; }
        .val { font-size: 40rpx; color: #E64340; font-weight: 700; }
      }
    }
  }
}

/* 酒店信息 */
.hotel-section {
  padding: 40rpx;

  .call-btn {
    background: $bg-color;
    padding: 8rpx 20rpx;
    border-radius: 100rpx;
    .icon { width: 28rpx; height: 28rpx; margin-right: 8rpx; }
    text { font-size: 22rpx; color: $primary-color; font-weight: 500; }
  }

  .hotel-info {
    margin-top: 20rpx;
    .address {
      font-size: 24rpx;
      color: $text-second;
      line-height: 1.6;
    }
  }
}

.state-box {
  padding: 200rpx 0;
  text-align: center;
  color: $text-sub;
  font-size: 26rpx;
}

/* 底部操作栏 */
.bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 100;
  background: #fff;
  padding: 24rpx 40rpx calc(24rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -4rpx 20rpx rgba(0,0,0,0.05);
  display: flex;
  justify-content: flex-end;
  gap: 24rpx;

  .btn {
    height: 80rpx;
    padding: 0 40rpx;
    font-size: 28rpx;
    border-radius: 100rpx;
    margin: 0;
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
</style>
