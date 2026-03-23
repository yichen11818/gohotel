<template>
  <scroll-view class="page" scroll-y>
    <swiper class="image-swiper" circular indicator-dots>
      <swiper-item v-for="(image, index) in roomImages" :key="index">
        <image class="room-image" :src="image" mode="aspectFill" />
      </swiper-item>
    </swiper>

    <view class="content">
      <view class="section">
        <view class="top-row">
          <view>
            <text class="room-name">{{ roomDetail.name }}</text>
            <text class="room-meta">{{ roomDetail.area }}㎡ · {{ roomDetail.bedType }} · 可住 {{ roomDetail.capacity }} 人</text>
          </view>
          <view class="price-box">
            <text class="price">¥{{ roomDetail.price }}</text>
            <text class="price-unit">/晚</text>
          </view>
        </view>
        <text class="room-desc">{{ roomDetail.description }}</text>
      </view>

      <view class="section">
        <text class="section-title">入住日期</text>
        <view class="date-grid">
          <picker mode="date" :value="checkInDate" @change="onCheckInChange">
            <view class="date-box">
              <text class="date-label">入住</text>
              <text class="date-value">{{ checkInDate }}</text>
            </view>
          </picker>
          <picker mode="date" :value="checkOutDate" @change="onCheckOutChange">
            <view class="date-box">
              <text class="date-label">离店</text>
              <text class="date-value">{{ checkOutDate }}</text>
            </view>
          </picker>
        </view>
        <view class="summary-row">
          <text>共 {{ nights }} 晚</text>
          <text>预计房费 ¥{{ totalPrice }}</text>
        </view>
      </view>

      <view class="section" v-if="roomDetail.facilities.length">
        <text class="section-title">房间设施</text>
        <view class="facility-list">
          <text v-for="item in roomDetail.facilities" :key="item" class="facility-item">{{ item }}</text>
        </view>
      </view>

      <view class="section rules">
        <text class="section-title">预订说明</text>
        <text class="rule-text">- 仅展示当前可预订房间，提交订单后等待酒店确认。</text>
        <text class="rule-text">- 入住日期需早于离店日期，取消规则以订单状态为准。</text>
      </view>
    </view>

    <view class="bottom-bar">
      <view>
        <text class="bottom-price">¥{{ totalPrice }}</text>
        <text class="bottom-note">共 {{ nights }} 晚</text>
      </view>
      <button class="primary-btn" @click="handleBook">立即预订</button>
    </view>
  </scroll-view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { hotel } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'

const roomId = ref('')
const roomImages = ref([])
const roomDetail = ref({
  id: '',
  name: '',
  price: 0,
  area: 0,
  bedType: '',
  capacity: 2,
  description: '',
  facilities: [],
  images: [],
})
const checkInDate = ref('')
const checkOutDate = ref('')

const nights = computed(() => {
  const start = new Date(checkInDate.value)
  const end = new Date(checkOutDate.value)
  const diff = end.getTime() - start.getTime()
  return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

const totalPrice = computed(() => Number(roomDetail.value.price || 0) * nights.value)

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

const ensureCheckoutAfterCheckin = () => {
  if (new Date(checkOutDate.value) <= new Date(checkInDate.value)) {
    const next = new Date(checkInDate.value)
    next.setDate(next.getDate() + 1)
    checkOutDate.value = formatDate(next)
  }
}

const onCheckInChange = (event) => {
  checkInDate.value = event.detail.value
  ensureCheckoutAfterCheckin()
}

const onCheckOutChange = (event) => {
  const nextValue = event.detail.value
  if (new Date(nextValue) <= new Date(checkInDate.value)) {
    uni.showToast({ title: '离店日期需晚于入住日期', icon: 'none' })
    return
  }
  checkOutDate.value = nextValue
}

const loadRoomDetail = async () => {
  try {
    const detail = await hotel.getRoomDetail(roomId.value)
    roomDetail.value = detail
    roomImages.value = detail.images?.length ? detail.images : [detail.image]
  } catch (error) {
    console.error('load room detail failed:', error)
  }
}

const handleBook = () => {
  const token = uni.getStorageSync(TOKEN_KEY)
  if (!token) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }

  uni.navigateTo({
    url: `/pages/booking/booking?roomId=${roomDetail.value.id}&checkIn=${checkInDate.value}&checkOut=${checkOutDate.value}`,
  })
}

onLoad((options) => {
  if (options?.id) {
    roomId.value = String(options.id)
  }
  checkInDate.value = normalizeDateParam(options?.checkIn)
  checkOutDate.value = normalizeDateParam(options?.checkOut)
  ensureDates()
  loadRoomDetail()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
  padding-bottom: 160rpx;
}

.image-swiper,
.room-image {
  width: 100%;
  height: 480rpx;
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

.top-row {
  display: flex;
  justify-content: space-between;
  gap: 20rpx;
}

.room-name,
.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: 700;
  color: #111827;
}

.room-meta,
.room-desc,
.rule-text,
.date-label,
.bottom-note {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #4b5563;
  line-height: 1.7;
}

.price-box {
  text-align: right;
}

.price,
.bottom-price {
  font-size: 40rpx;
  color: #b7791f;
  font-weight: 700;
}

.price-unit {
  font-size: 24rpx;
  color: #6b7280;
}

.date-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20rpx;
  margin-top: 20rpx;
}

.date-box {
  padding: 24rpx;
  background: #f9fafb;
  border-radius: 20rpx;
  border: 1px solid #eef2f7;
}

.date-value {
  display: block;
  margin-top: 8rpx;
  font-size: 28rpx;
  color: #111827;
  font-weight: 600;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  margin-top: 20rpx;
  font-size: 26rpx;
  color: #374151;
}

.facility-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-top: 20rpx;
}

.facility-item {
  padding: 10rpx 20rpx;
  background: #faf4ed;
  color: #8b5e34;
  font-size: 22rpx;
  border-radius: 999rpx;
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
