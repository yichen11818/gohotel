<template>
  <scroll-view class="page" scroll-y>
    <view class="content" v-if="orderDetail.id">
      <view class="status-card">
        <text class="status-title">{{ orderDetail.statusText }}</text>
        <text class="status-desc">订单号 {{ orderDetail.orderNo }}</text>
      </view>

      <view class="section room-card">
        <image class="room-image" :src="orderDetail.roomImage" mode="aspectFill" />
        <view class="room-body">
          <text class="room-name">{{ orderDetail.roomName }}</text>
          <text class="room-meta">入住 {{ orderDetail.checkIn }} · 离店 {{ orderDetail.checkOut }}</text>
          <text class="room-meta">{{ orderDetail.totalDays }} 晚 · {{ orderDetail.roomNumber || '待分配房号' }}</text>
        </view>
      </view>

      <view class="section">
        <text class="section-title">入住人信息</text>
        <view class="info-row">
          <text>姓名</text>
          <text>{{ orderDetail.guestName || '未填写' }}</text>
        </view>
        <view class="info-row">
          <text>手机号</text>
          <text>{{ orderDetail.guestPhone || '未填写' }}</text>
        </view>
        <view class="info-row">
          <text>身份证号</text>
          <text>{{ orderDetail.guestIdCard || '未填写' }}</text>
        </view>
        <view class="info-row">
          <text>特殊需求</text>
          <text class="multiline">{{ orderDetail.specialRequest || '无' }}</text>
        </view>
      </view>

      <view class="section">
        <text class="section-title">订单信息</text>
        <view class="info-row">
          <text>创建时间</text>
          <text>{{ orderDetail.createdAt || '-' }}</text>
        </view>
        <view class="info-row">
          <text>订单状态</text>
          <text>{{ orderDetail.statusText }}</text>
        </view>
        <view class="info-row">
          <text>支付状态</text>
          <text>{{ paymentStatusText }}</text>
        </view>
        <view class="info-row total-row">
          <text>订单金额</text>
          <text class="amount">¥{{ orderDetail.totalPrice }}</text>
        </view>
      </view>

      <view class="section" v-if="hotelInfo.phone || hotelInfo.address">
        <text class="section-title">酒店联系</text>
        <view class="info-row">
          <text>前台电话</text>
          <text>{{ hotelInfo.phone || '暂未配置' }}</text>
        </view>
        <view class="info-row">
          <text>酒店地址</text>
          <text class="multiline">{{ hotelInfo.address }}</text>
        </view>
      </view>
    </view>

    <view v-else class="placeholder">正在加载订单详情...</view>

    <view class="bottom-bar" v-if="orderDetail.id">
      <button v-if="orderDetail.canCancel" class="ghost-btn" @click="cancelOrder">取消订单</button>
      <button v-if="orderDetail.roomId" class="primary-btn" @click="viewRoom">查看房型</button>
    </view>
  </scroll-view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { booking, hotel } from '@/api/index.js'

const orderId = ref('')
const orderDetail = ref({})
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
    title: '取消订单',
    content: '确认取消这笔订单吗？',
    success: async (res) => {
      if (!res.confirm) return
      try {
        uni.showLoading({ title: '处理中...' })
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
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${orderDetail.value.roomId}&checkIn=${orderDetail.value.checkIn}&checkOut=${orderDetail.value.checkOut}`,
  })
}

onLoad((options) => {
  if (options?.id) {
    orderId.value = String(options.id)
    loadOrderDetail()
  }
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
  padding-bottom: 160rpx;
}

.content {
  padding: 24rpx;
}

.status-card,
.section {
  padding: 28rpx;
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.status-card {
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
}

.status-title,
.section-title,
.room-name {
  display: block;
  font-size: 32rpx;
  font-weight: 700;
}

.status-desc,
.room-meta,
.placeholder {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  line-height: 1.7;
}

.placeholder {
  padding: 120rpx 24rpx;
  text-align: center;
  color: #6b7280;
}

.room-card {
  display: flex;
  gap: 20rpx;
  align-items: center;
}

.room-image {
  width: 220rpx;
  height: 180rpx;
  border-radius: 20rpx;
  background: #f3f4f6;
}

.room-body {
  flex: 1;
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  padding: 18rpx 0;
  font-size: 26rpx;
  color: #374151;
  border-bottom: 1px solid #f3f4f6;
}

.info-row:last-child {
  border-bottom: 0;
}

.multiline {
  flex: 1;
  text-align: right;
  line-height: 1.6;
}

.total-row,
.amount {
  font-weight: 700;
}

.amount {
  color: #b7791f;
}

.bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: flex-end;
  gap: 20rpx;
  padding: 24rpx;
  background: rgba(255, 255, 255, 0.96);
  border-top: 1px solid #eef2f7;
}

.ghost-btn,
.primary-btn {
  margin: 0;
  border-radius: 999rpx;
}

.primary-btn {
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
}
</style>
