<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">订单详情</text>
        </view>
      </template>
    </TnNavbar>

    <scroll-view class="main-content" scroll-y>
      <!-- 订单状态 -->
      <view class="status-card" :class="orderStatus">
        <view class="status-icon">
          <TnIcon :name="statusIcon" color="#fff" size="56" />
        </view>
        <view class="status-info">
          <text class="status-text">{{ statusText }}</text>
          <text class="status-desc">{{ statusDesc }}</text>
        </view>
      </view>

      <!-- 房间信息 -->
      <view class="room-card">
        <image class="room-image" :src="orderDetail.roomImage" mode="aspectFill"></image>
        <view class="room-info">
          <text class="room-name">{{ orderDetail.roomName }}</text>
          <text class="room-specs">{{ orderDetail.roomSpecs }}</text>
          <view class="room-tags">
            <text v-for="(tag, index) in orderDetail.tags" :key="index" class="tag">{{ tag }}</text>
          </view>
        </view>
      </view>

      <!-- 入住信息 -->
      <view class="info-card">
        <view class="card-title">入住信息</view>
        
        <view class="info-row">
          <text class="label">入住日期</text>
          <text class="value">{{ orderDetail.checkIn }} 至 {{ orderDetail.checkOut }}</text>
        </view>
        
        <view class="info-row">
          <text class="label">入住天数</text>
          <text class="value">{{ orderDetail.nights }}晚</text>
        </view>
        
        <view class="info-row">
          <text class="label">房间数量</text>
          <text class="value">{{ orderDetail.roomCount }}间</text>
        </view>
        
        <view class="info-row">
          <text class="label">入住人</text>
          <text class="value">{{ orderDetail.guestName }}</text>
        </view>
        
        <view class="info-row">
          <text class="label">联系电话</text>
          <text class="value">{{ orderDetail.guestPhone }}</text>
        </view>
      </view>

      <!-- 订单信息 -->
      <view class="info-card">
        <view class="card-title">订单信息</view>
        
        <view class="info-row">
          <text class="label">订单编号</text>
          <view class="value-with-action">
            <text class="value">{{ orderDetail.orderNo }}</text>
            <text class="copy-btn" @click="copyOrderNo">复制</text>
          </view>
        </view>
        
        <view class="info-row">
          <text class="label">下单时间</text>
          <text class="value">{{ orderDetail.createTime }}</text>
        </view>
        
        <view class="info-row">
          <text class="label">支付方式</text>
          <text class="value">{{ orderDetail.payMethod }}</text>
        </view>
        
        <view class="info-row" v-if="orderDetail.payTime">
          <text class="label">支付时间</text>
          <text class="value">{{ orderDetail.payTime }}</text>
        </view>
      </view>

      <!-- 价格明细 -->
      <view class="info-card">
        <view class="card-title">价格明细</view>
        
        <view class="info-row">
          <text class="label">房费</text>
          <text class="value">¥{{ orderDetail.roomPrice }} × {{ orderDetail.nights }}晚 × {{ orderDetail.roomCount }}间</text>
        </view>
        
        <view class="info-row" v-if="orderDetail.couponAmount">
          <text class="label">优惠券</text>
          <text class="value discount">-¥{{ orderDetail.couponAmount }}</text>
        </view>
        
        <view class="info-row total">
          <text class="label">实付金额</text>
          <view class="total-price">
            <text class="currency">¥</text>
            <text class="amount">{{ orderDetail.totalPrice }}</text>
          </view>
        </view>
      </view>

      <!-- 酒店信息 -->
      <view class="hotel-card">
        <view class="hotel-info">
          <text class="hotel-name">七天酒店</text>
          <text class="hotel-address">湖北省武汉市硚口区晴川街道沿河大道246号</text>
        </view>
        <view class="hotel-actions">
          <view class="action-btn" @click="callHotel">
            <TnIcon name="phone" color="#C29D71" size="40" />
            <text>电话</text>
          </view>
          <view class="action-btn" @click="navigateToHotel">
            <TnIcon name="location" color="#C29D71" size="40" />
            <text>导航</text>
          </view>
        </view>
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar" v-if="showBottomActions">
      <view v-if="orderStatus === 'pending'" class="action-row">
        <view class="cancel-btn" @click="cancelOrder">
          <text>取消订单</text>
        </view>
        <view class="pay-btn" @click="payOrder">
          <text>立即支付</text>
        </view>
      </view>
      
      <view v-else-if="orderStatus === 'completed'" class="action-row">
        <view class="review-btn" @click="writeReview">
          <text>评价订单</text>
        </view>
        <view class="rebook-btn" @click="rebookOrder">
          <text>再次预订</text>
        </view>
      </view>
      
      <view v-else-if="orderStatus === 'cancelled'" class="action-row">
        <view class="rebook-btn full" @click="rebookOrder">
          <text>重新预订</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { booking } from '@/api/index.js'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'

const orderId = ref('')
const loading = ref(false)

// 订单详情
const orderDetail = ref({
  id: '',
  orderNo: '',
  status: '',
  roomName: '',
  roomImage: '',
  roomSpecs: '',
  tags: [],
  checkIn: '',
  checkOut: '',
  nights: 0,
  roomCount: 0,
  guestName: '',
  guestPhone: '',
  createTime: '',
  payTime: '',
  payMethod: '',
  roomPrice: 0,
  couponAmount: 0,
  totalPrice: 0
})

// 订单状态
const orderStatus = computed(() => orderDetail.value.status)

// 状态图标
const statusIcon = computed(() => {
  const icons = {
    pending: 'time-circle',
    checkin: 'check-circle',
    completed: 'check-circle-fill',
    cancelled: 'close-circle'
  }
  return icons[orderStatus.value] || 'info-circle'
})

// 状态文字
const statusText = computed(() => {
  const texts = {
    pending: '待付款',
    checkin: '待入住',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[orderStatus.value] || '未知状态'
})

// 状态描述
const statusDesc = computed(() => {
  const descs = {
    pending: '请在30分钟内完成支付',
    checkin: '预订成功，期待您的入住',
    completed: '感谢您的入住，期待再次光临',
    cancelled: '订单已取消'
  }
  return descs[orderStatus.value] || ''
})

// 是否显示底部操作
const showBottomActions = computed(() => {
  return ['pending', 'completed', 'cancelled'].includes(orderStatus.value)
})

// 返回
const goBack = () => {
  uni.navigateBack()
}

// 复制订单号
const copyOrderNo = () => {
  uni.setClipboardData({
    data: orderDetail.value.orderNo,
    success: () => {
      uni.showToast({ title: '已复制', icon: 'success' })
    }
  })
}

// 拨打酒店电话
const callHotel = () => {
  uni.makePhoneCall({
    phoneNumber: '027-12345678'
  })
}

// 导航到酒店
const navigateToHotel = () => {
  uni.openLocation({
    latitude: 30.56,
    longitude: 114.28,
    name: '七天酒店',
    address: '湖北省武汉市硚口区晴川街道沿河大道246号'
  })
}

// 取消订单
const cancelOrder = () => {
  uni.showModal({
    title: '取消订单',
    content: '确定要取消此订单吗？',
    success: (res) => {
      if (res.confirm) {
        uni.showToast({ title: '订单已取消', icon: 'success' })
        orderDetail.value.status = 'cancelled'
      }
    }
  })
}

// 支付订单
const payOrder = () => {
  uni.showLoading({ title: '支付中...' })
  setTimeout(() => {
    uni.hideLoading()
    uni.redirectTo({
      url: `/pages/pay-result/pay-result?success=true&orderId=${orderDetail.value.orderNo}&amount=${orderDetail.value.totalPrice}`
    })
  }, 1500)
}

// 评价订单
const writeReview = () => {
  uni.showToast({ title: '评价功能开发中', icon: 'none' })
}

// 再次预订
const rebookOrder = () => {
  uni.navigateTo({
    url: '/pages/room-detail/room-detail?id=1'
  })
}

onLoad((options) => {
  if (options?.id) {
    orderId.value = options.id
    loadOrderDetail()
  }
})

// 加载订单详情
const loadOrderDetail = async () => {
  loading.value = true
  try {
    const data = await booking.getBookingDetail(orderId.value)
    if (data) {
      orderDetail.value = {
        id: data.id,
        orderNo: data.order_no,
        status: data.status,
        roomName: data.room_type_name,
        roomImage: data.room_type_image,
        roomSpecs: `${data.room_area}m² · ${data.bed_type}`,
        tags: data.tags || [],
        checkIn: data.check_in_date,
        checkOut: data.check_out_date,
        nights: data.nights,
        roomCount: data.room_count,
        guestName: data.guest_name,
        guestPhone: data.guest_phone,
        createTime: data.created_at,
        payTime: data.paid_at,
        payMethod: data.payment_method === 'wechat' ? '微信支付' : '其他支付',
        roomPrice: data.price,
        couponAmount: data.coupon_amount || 0,
        totalPrice: data.total_amount
      }
    }
  } catch (error) {
    console.error('Failed to load order detail:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 140rpx;
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
  height: calc(100vh - 88rpx - 140rpx);
}

/* 状态卡片 */
.status-card {
  display: flex;
  align-items: center;
  gap: 24rpx;
  padding: 40rpx 32rpx;
  
  &.pending {
    background: linear-gradient(135deg, #FFB347 0%, #FF9500 100%);
  }
  
  &.checkin {
    background: linear-gradient(135deg, #5AC8FA 0%, #4A90E2 100%);
  }
  
  &.completed {
    background: linear-gradient(135deg, #34C759 0%, #30B350 100%);
  }
  
  &.cancelled {
    background: linear-gradient(135deg, #999 0%, #666 100%);
  }
  
  .status-icon {
    width: 80rpx;
    height: 80rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(255, 255, 255, 0.2);
    border-radius: 50%;
  }
  
  .status-info {
    display: flex;
    flex-direction: column;
    gap: 8rpx;
    
    .status-text {
      font-size: 36rpx;
      font-weight: 700;
      color: #fff;
    }
    
    .status-desc {
      font-size: 26rpx;
      color: rgba(255, 255, 255, 0.9);
    }
  }
}

/* 房间卡片 */
.room-card {
  display: flex;
  gap: 20rpx;
  padding: 24rpx;
  margin: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .room-image {
    width: 180rpx;
    height: 140rpx;
    border-radius: 16rpx;
  }
  
  .room-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 12rpx;
    
    .room-name {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
    }
    
    .room-specs {
      font-size: 26rpx;
      color: #999;
    }
    
    .room-tags {
      display: flex;
      gap: 12rpx;
      
      .tag {
        padding: 6rpx 12rpx;
        background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
        border-radius: 6rpx;
        font-size: 20rpx;
        color: #C29D71;
      }
    }
  }
}

/* 信息卡片 */
.info-card {
  margin: 0 24rpx 24rpx;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .card-title {
    font-size: 32rpx;
    font-weight: 600;
    color: #333;
    margin-bottom: 24rpx;
  }
  
  .info-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16rpx 0;
    
    .label {
      font-size: 28rpx;
      color: #999;
    }
    
    .value {
      font-size: 28rpx;
      color: #333;
      
      &.discount {
        color: #FF3B30;
      }
    }
    
    .value-with-action {
      display: flex;
      align-items: center;
      gap: 16rpx;
      
      .copy-btn {
        font-size: 24rpx;
        color: #C29D71;
        padding: 4rpx 12rpx;
        border: 1rpx solid #C29D71;
        border-radius: 8rpx;
        
        &:active {
          opacity: 0.7;
        }
      }
    }
    
    &.total {
      padding-top: 24rpx;
      margin-top: 16rpx;
      border-top: 1rpx solid #F0F0F0;
      
      .label {
        font-size: 30rpx;
        color: #333;
        font-weight: 600;
      }
      
      .total-price {
        display: flex;
        align-items: baseline;
        
        .currency {
          font-size: 28rpx;
          color: #C29D71;
          font-weight: 600;
        }
        
        .amount {
          font-size: 44rpx;
          font-weight: 700;
          color: #C29D71;
        }
      }
    }
  }
}

/* 酒店卡片 */
.hotel-card {
  margin: 0 24rpx 24rpx;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .hotel-info {
    margin-bottom: 24rpx;
    
    .hotel-name {
      display: block;
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
      margin-bottom: 12rpx;
    }
    
    .hotel-address {
      font-size: 26rpx;
      color: #999;
    }
  }
  
  .hotel-actions {
    display: flex;
    gap: 24rpx;
    
    .action-btn {
      flex: 1;
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 12rpx;
      padding: 24rpx;
      background: linear-gradient(135deg, rgba(194, 157, 113, 0.1) 0%, rgba(194, 157, 113, 0.05) 100%);
      border-radius: 16rpx;
      
      text {
        font-size: 26rpx;
        color: #666;
      }
      
      &:active {
        opacity: 0.7;
      }
    }
  }
}

.bottom-placeholder {
  height: 40rpx;
}

/* 底部操作栏 */
.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 24rpx 32rpx;
  padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
  background: #fff;
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.06);
  
  .action-row {
    display: flex;
    gap: 24rpx;
  }
  
  .cancel-btn,
  .review-btn {
    flex: 1;
    height: 88rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #F5F5F5;
    border-radius: 100rpx;
    
    text {
      font-size: 30rpx;
      color: #666;
      font-weight: 500;
    }
    
    &:active {
      background: #EEEEEE;
    }
  }
  
  .pay-btn,
  .rebook-btn {
    flex: 1;
    height: 88rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%);
    border-radius: 100rpx;
    box-shadow: 0 12rpx 24rpx rgba(194, 157, 113, 0.35);
    
    text {
      font-size: 30rpx;
      color: #fff;
      font-weight: 600;
    }
    
    &:active {
      transform: scale(0.98);
    }
    
    &.full {
      flex: none;
      width: 100%;
    }
  }
}
</style>

