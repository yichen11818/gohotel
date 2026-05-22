<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">确认订单</text>
      <view class="placeholder-view"></view>
    </view>

    <scroll-view class="scroll-container" scroll-y>
      <view class="content">
        <!-- 房型卡片 -->
        <view class="section room-card premium-card">
          <image class="room-image" :src="roomInfo.image" mode="aspectFill" />
          <view class="room-body">
            <text class="room-name">{{ roomInfo.name }}</text>
            <text class="room-meta">{{ roomInfo.area }}㎡ · {{ roomInfo.bedType }}</text>
            <view class="room-price">
              <text class="symbol">¥</text>
              <text class="val">{{ roomInfo.price }}</text>
              <text class="unit">/晚</text>
            </view>
          </view>
        </view>

        <!-- 入住信息 -->
        <view class="section premium-card" @click="showCalendar = true">
          <view class="section-header flex-between">
            <text class="section-title">入住信息</text>
            <text class="color-primary" style="font-size: 24rpx;">修改</text>
          </view>
          <view class="date-list">
            <view class="date-row flex-between">
              <text class="label">入住日期</text>
              <text class="val">{{ checkInDate }}</text>
            </view>
            <view class="date-row flex-between">
              <text class="label">离店日期</text>
              <text class="val">{{ checkOutDate }}</text>
            </view>
            <view class="date-row flex-between">
              <text class="label">共计晚数</text>
              <text class="val color-primary">{{ nights }} 晚</text>
            </view>
          </view>
        </view>

        <!-- 入住人信息 -->
        <view class="section premium-card">
          <view class="section-header">
            <text class="section-title">入住人信息</text>
          </view>
          <view class="form-list">
            <view class="field">
              <text class="field-label">姓名</text>
              <input v-model="guestInfo.name" class="field-input" placeholder="请输入入住人姓名" placeholder-class="placeholder" />
            </view>
            <view class="field">
              <text class="field-label">手机号</text>
              <input v-model="guestInfo.phone" class="field-input" type="number" maxlength="11" placeholder="请输入手机号" placeholder-class="placeholder" />
            </view>
            <view class="field">
              <text class="field-label">身份证号</text>
              <input v-model="guestInfo.idCard" class="field-input" placeholder="选填，用于快捷入住" placeholder-class="placeholder" />
            </view>
            <view class="field">
              <text class="field-label">特殊需求</text>
              <textarea v-model="specialRequest" class="field-textarea" maxlength="200" placeholder="选填，如安静房、尽量高楼层等" placeholder-class="placeholder" />
            </view>
          </view>
        </view>

        <!-- 价格明细 -->
        <view class="section premium-card">
          <view class="section-header">
            <text class="section-title">价格明细</text>
          </view>
          <view class="price-list">
            <view class="price-row flex-between">
              <text class="label">在线支付房费</text>
              <text class="val">¥{{ roomInfo.price }} × {{ nights }}</text>
            </view>
            <view class="total-row flex-between">
              <text class="label">应付总额</text>
              <view class="price">
                <text class="symbol">¥</text>
                <text class="val">{{ totalPrice }}</text>
              </view>
            </view>
          </view>
        </view>

        <!-- 温馨提示 -->
        <view class="section tips">
          <text class="tip-text">· 提交后将生成预订订单，请在 15 分钟内完成确认。</text>
          <text class="tip-text">· 入住时请出示有效身份证件。</text>
        </view>
      </view>
      <view class="bottom-space"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar safe-area-inset-bottom">
      <view class="price-info">
        <text class="bottom-note">应付金额</text>
        <view class="bottom-price">
          <text class="symbol">¥</text>
          <text class="val">{{ totalPrice }}</text>
        </view>
      </view>
      <button class="primary-btn premium-button" :disabled="submitting" @click="handleSubmit">
        <text>{{ submitting ? '正在提交...' : '确认预订' }}</text>
      </button>
    </view>
    <HotelCalendar v-model:show="showCalendar" :check-in="checkInDate" :check-out="checkOutDate" @confirm="onCalendarConfirm" />
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { booking, hotel, user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'
import HotelCalendar from '@/components/hotel-calendar/hotel-calendar.vue'

const statusBarHeight = ref(44)
const showCalendar = ref(false)
const submitting = ref(false)
const roomInfo = ref({
  id: '',
  name: '',
  image: '',
  area: 0,
  bedType: '',
  price: 0,
})
const checkInDate = ref('')
const checkOutDate = ref('')
const specialRequest = ref('')
const guestInfo = ref({
  name: '',
  phone: '',
  idCard: '',
})

const nights = computed(() => {
  if (!checkInDate.value || !checkOutDate.value) return 1
  const start = new Date(checkInDate.value)
  const end = new Date(checkOutDate.value)
  const diff = end.getTime() - start.getTime()
  return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

const totalPrice = computed(() => Number(roomInfo.value.price || 0) * nights.value)

const onCalendarConfirm = (dates) => {
  checkInDate.value = dates.checkIn
  checkOutDate.value = dates.checkOut
}

const goBack = () => {
  uni.navigateBack()
}

const loadRoomInfo = async (roomId) => {
  try {
    const detail = await hotel.getRoomDetail(roomId)
    roomInfo.value = {
      id: detail.id,
      name: detail.name,
      image: detail.image,
      area: detail.area,
      bedType: detail.bedType,
      price: detail.price,
    }
  } catch (error) {
    console.error('load room info failed:', error)
  }
}

const fillGuestInfo = async () => {
  const token = uni.getStorageSync(TOKEN_KEY)
  if (!token) return

  // 1. 先尝试从缓存读
  const cache = uni.getStorageSync(USER_INFO_KEY)
  if (cache) {
    guestInfo.value.name = cache.real_name || cache.nickname || cache.username || ''
    guestInfo.value.phone = cache.phone || ''
    guestInfo.value.idCard = cache.id_card || ''
  }

  // 2. 异步请求最新数据（更准）
  try {
    const profile = await user.getUserInfo()
    guestInfo.value.name = profile.real_name || profile.nickname || profile.username || guestInfo.value.name
    guestInfo.value.phone = profile.phone || guestInfo.value.phone
    guestInfo.value.idCard = profile.id_card || guestInfo.value.idCard
  } catch (_error) {}
}

const handleSubmit = async () => {
  if (!guestInfo.value.name.trim()) {
    uni.showToast({ title: '请输入入住人姓名', icon: 'none' })
    return
  }
  if (!/^1\d{10}$/.test(guestInfo.value.phone)) {
    uni.showToast({ title: '手机号格式不正确', icon: 'none' })
    return
  }

  try {
    submitting.value = true
    uni.showLoading({ title: '提交预订' })

    const result = await booking.createBooking({
      room_id: roomInfo.value.id,
      check_in: checkInDate.value,
      check_out: checkOutDate.value,
      guest_name: guestInfo.value.name.trim(),
      guest_phone: guestInfo.value.phone,
      guest_id_card: guestInfo.value.idCard.trim(),
      special_request: specialRequest.value.trim(),
    })

    uni.hideLoading()
    uni.showToast({ title: '预订成功', icon: 'success' })

    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/order-detail/order-detail?id=${result.id}`,
      })
    }, 1000)
  } catch (_error) {
    uni.hideLoading()
  } finally {
    submitting.value = false
  }
}

onLoad((options) => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44

  if (options?.roomId) loadRoomInfo(options.roomId)
  checkInDate.value = options?.checkIn || ''
  checkOutDate.value = options?.checkOut || ''
  fillGuestInfo()
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

.scroll-container {
  flex: 1;
  height: 0;
}

.content {
  padding: 30rpx;
}

/* 房型卡片 */
.room-card {
  display: flex;
  padding: 30rpx;
  margin-bottom: 30rpx;

  .room-image {
    width: 160rpx;
    height: 160rpx;
    border-radius: 12rpx;
    background: #f5f5f5;
    margin-right: 24rpx;
  }

  .room-body {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;

    .room-name {
      font-size: 32rpx;
      font-weight: 700;
      color: $text-main;
    }
    .room-meta {
      font-size: 24rpx;
      color: $text-sub;
    }
    .room-price {
      .symbol { font-size: 24rpx; color: #E64340; font-weight: 700; }
      .val { font-size: 36rpx; color: #E64340; font-weight: 700; }
      .unit { font-size: 22rpx; color: $text-sub; margin-left: 4rpx; }
    }
  }
}

.section {
  padding: 40rpx;
  margin-bottom: 30rpx;

  .section-header {
    margin-bottom: 30rpx;
    .section-title {
      font-size: 32rpx;
      font-weight: 700;
      color: $text-main;
    }
  }
}

/* 日期/价格 列表 */
.date-list, .price-list {
  .date-row, .price-row {
    padding: 20rpx 0;
    font-size: 28rpx;
    .label { color: $text-sub; }
    .val { color: $text-main; font-weight: 500; }
    .color-primary { color: $primary-color; }
  }

  .total-row {
    margin-top: 10rpx;
    padding-top: 30rpx;
    border-top: 1rpx solid #f8f8f8;
    .label { font-size: 30rpx; font-weight: 700; color: $text-main; }
    .price {
      .symbol { font-size: 24rpx; color: #E64340; font-weight: 700; }
      .val { font-size: 40rpx; color: #E64340; font-weight: 700; }
    }
  }
}

/* 表单 */
.field {
  margin-bottom: 30rpx;
  &:last-child { margin-bottom: 0; }

  .field-label {
    font-size: 26rpx;
    color: $text-sub;
    margin-bottom: 16rpx;
    display: block;
  }

  .field-input, .field-textarea {
    width: 100%;
    background: #f9f9f9;
    border: 1rpx solid #eee;
    border-radius: 12rpx;
    padding: 20rpx 24rpx;
    box-sizing: border-box;
    font-size: 28rpx;
    color: $text-main;

    &:focus {
      background: #fff;
      border-color: $primary-color;
    }
  }

  .field-textarea { height: 160rpx; }
  .placeholder { color: #ccc; }
}

.tips {
  background: #FFF9F0;
  border: 0;
  .tip-text {
    display: block;
    font-size: 22rpx;
    color: #AD8551;
    line-height: 1.6;
    margin-bottom: 8rpx;
  }
}

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
  justify-content: space-between;
  align-items: center;

  .price-info {
    .bottom-note { font-size: 22rpx; color: $text-sub; display: block; }
    .bottom-price {
      .symbol { font-size: 24rpx; color: #E64340; font-weight: 700; }
      .val { font-size: 44rpx; color: #E64340; font-weight: 700; }
    }
  }

  .primary-btn {
    width: 280rpx;
    height: 100rpx;
    margin: 0;
  }
}

.bottom-space { height: 60rpx; }
</style>
