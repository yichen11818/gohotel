<template>
  <scroll-view class="page" scroll-y>
    <view class="content">
      <view class="section room-card">
        <image class="room-image" :src="roomInfo.image" mode="aspectFill" />
        <view class="room-body">
          <text class="room-name">{{ roomInfo.name }}</text>
          <text class="room-meta">{{ roomInfo.area }}㎡ · {{ roomInfo.bedType }}</text>
          <text class="room-price">¥{{ roomInfo.price }}/晚</text>
        </view>
      </view>

      <view class="section">
        <text class="section-title">入住信息</text>
        <view class="date-row">
          <text>入住</text>
          <text>{{ checkInDate }}</text>
        </view>
        <view class="date-row">
          <text>离店</text>
          <text>{{ checkOutDate }}</text>
        </view>
        <view class="date-row">
          <text>共计</text>
          <text>{{ nights }} 晚</text>
        </view>
      </view>

      <view class="section">
        <text class="section-title">入住人信息</text>
        <view class="field">
          <text class="field-label">姓名</text>
          <input v-model="guestInfo.name" class="field-input" placeholder="请输入入住人姓名" />
        </view>
        <view class="field">
          <text class="field-label">手机号</text>
          <input v-model="guestInfo.phone" class="field-input" type="number" maxlength="11" placeholder="请输入手机号" />
        </view>
        <view class="field">
          <text class="field-label">身份证号</text>
          <input v-model="guestInfo.idCard" class="field-input" placeholder="选填，用于快捷入住" />
        </view>
        <view class="field">
          <text class="field-label">特殊需求</text>
          <textarea v-model="specialRequest" class="field-textarea" maxlength="200" placeholder="选填，如安静房、尽量高楼层等" />
        </view>
      </view>

      <view class="section">
        <text class="section-title">价格明细</text>
        <view class="date-row">
          <text>房费</text>
          <text>¥{{ roomInfo.price }} × {{ nights }} 晚</text>
        </view>
        <view class="date-row total-row">
          <text>合计</text>
          <text class="total-price">¥{{ totalPrice }}</text>
        </view>
      </view>

      <view class="section tips">
        <text class="tip-text">- 提交后将生成预订订单，等待酒店确认。</text>
        <text class="tip-text">- 当前版本不支持在线支付，订单请在“我的订单”中查看状态。</text>
      </view>
    </view>

    <view class="bottom-bar">
      <view>
        <text class="bottom-note">订单金额</text>
        <text class="bottom-price">¥{{ totalPrice }}</text>
      </view>
      <button class="primary-btn" :disabled="submitting" @click="handleSubmit">
        {{ submitting ? '提交中...' : '提交订单' }}
      </button>
    </view>
  </scroll-view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { booking, hotel, user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'

const submitting = ref(false)
const roomInfo = ref({
  id: '',
  name: '',
  image: 'https://dummyimage.com/720x420/f5f5f5/999999&text=Room',
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
  const start = new Date(checkInDate.value)
  const end = new Date(checkOutDate.value)
  const diff = end.getTime() - start.getTime()
  return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

const totalPrice = computed(() => Number(roomInfo.value.price || 0) * nights.value)

const formatDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const normalizeDateParam = (value) => {
  if (!value) return ''
  return String(value).slice(0, 10)
}

const ensureDates = () => {
  if (!checkInDate.value || !checkOutDate.value) {
    const today = new Date()
    const tomorrow = new Date()
    tomorrow.setDate(today.getDate() + 1)
    checkInDate.value = formatDate(today)
    checkOutDate.value = formatDate(tomorrow)
  }
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
  if (!token) {
    uni.redirectTo({ url: '/pages/login/login' })
    return
  }

  const cache = uni.getStorageSync(USER_INFO_KEY)
  if (cache) {
    guestInfo.value.name = cache.nickname || cache.real_name || cache.username || ''
    guestInfo.value.phone = cache.phone || ''
  }

  try {
    const profile = await user.getUserInfo()
    guestInfo.value.name = profile.nickname || profile.real_name || profile.username || guestInfo.value.name
    guestInfo.value.phone = profile.phone || guestInfo.value.phone
  } catch (_error) {}
}

const handleSubmit = async () => {
  if (!guestInfo.value.name.trim()) {
    uni.showToast({ title: '请输入入住人姓名', icon: 'none' })
    return
  }

  if (!/^1\d{10}$/.test(guestInfo.value.phone)) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }

  try {
    submitting.value = true
    uni.showLoading({ title: '提交中...' })

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
    uni.showToast({ title: '预订提交成功', icon: 'success' })

    setTimeout(() => {
      uni.redirectTo({
        url: `/pages/order-detail/order-detail?id=${result.id}`,
      })
    }, 800)
  } catch (_error) {
    uni.hideLoading()
  } finally {
    submitting.value = false
  }
}

onLoad((options) => {
  if (options?.roomId) {
    loadRoomInfo(options.roomId)
  }
  checkInDate.value = normalizeDateParam(options?.checkIn)
  checkOutDate.value = normalizeDateParam(options?.checkOut)
  ensureDates()
  fillGuestInfo()
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

.section {
  padding: 28rpx;
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
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

.room-name,
.section-title {
  display: block;
  font-size: 30rpx;
  font-weight: 700;
  color: #111827;
}

.room-meta,
.room-price,
.field-label,
.tip-text,
.bottom-note {
  display: block;
  margin-top: 10rpx;
  font-size: 24rpx;
  color: #4b5563;
}

.room-price,
.total-price,
.bottom-price {
  color: #b7791f;
  font-weight: 700;
}

.field {
  margin-top: 24rpx;
}

.field-input,
.field-textarea {
  width: 100%;
  margin-top: 12rpx;
  padding: 22rpx 24rpx;
  background: #f9fafb;
  border: 1px solid #eef2f7;
  border-radius: 18rpx;
  box-sizing: border-box;
}

.field-textarea {
  min-height: 160rpx;
}

.date-row {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  padding: 18rpx 0;
  font-size: 26rpx;
  color: #374151;
  border-bottom: 1px solid #f3f4f6;
}

.date-row:last-child {
  border-bottom: 0;
}

.total-row {
  font-weight: 700;
}

.tips {
  background: #fffbeb;
}

.bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 24rpx;
  padding: 24rpx;
  background: rgba(255, 255, 255, 0.96);
  border-top: 1px solid #eef2f7;
}

.primary-btn {
  min-width: 240rpx;
  margin: 0;
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
  border-radius: 999rpx;
}
</style>
