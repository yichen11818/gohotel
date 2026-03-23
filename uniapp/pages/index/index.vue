<template>
  <scroll-view class="page" scroll-y>
    <view class="hero">
      <swiper class="hero-swiper" circular autoplay indicator-dots>
        <swiper-item v-for="(image, index) in bannerImages" :key="index">
          <image class="hero-image" :src="image" mode="aspectFill" />
        </swiper-item>
      </swiper>
      <view class="hero-mask"></view>
      <view class="hero-content">
        <text class="hotel-name">{{ hotelInfo.name }}</text>
        <text class="hotel-address">{{ hotelInfo.address }}</text>
        <view class="hero-actions">
          <button class="ghost-btn" size="mini" @click="openMap">地图导航</button>
          <button class="ghost-btn" size="mini" @click="callHotel" :disabled="!hotelInfo.phone">联系前台</button>
        </view>
      </view>
    </view>

    <view class="content">
      <view class="card">
        <view class="card-header">
          <text class="card-title">入住日期</text>
          <text class="card-subtitle">选择日期后查看可订房型</text>
        </view>
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
        <view class="booking-summary">
          <text>共 {{ nights }} 晚</text>
          <text v-if="hotelInfo.lowestPrice">低至 ¥{{ hotelInfo.lowestPrice }}/晚</text>
        </view>
        <button class="primary-btn" @click="goToHotel">查看可订房型</button>
      </view>

      <view v-if="notices.length" class="card">
        <view class="card-header">
          <text class="card-title">酒店公告</text>
        </view>
        <view v-for="notice in notices" :key="notice.id" class="notice-item">
          <text class="notice-title">{{ notice.title }}</text>
        </view>
      </view>

      <view class="card">
        <view class="card-header">
          <text class="card-title">酒店信息</text>
        </view>
        <view class="info-row">
          <text class="info-label">前台电话</text>
          <text class="info-value">{{ hotelInfo.phone || '暂未配置' }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">服务时间</text>
          <text class="info-value">{{ hotelInfo.serviceTime }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">入住 / 退房</text>
          <text class="info-value">{{ hotelInfo.checkInTime }} / {{ hotelInfo.checkOutTime }}</text>
        </view>
        <text class="hotel-intro">{{ hotelIntroText }}</text>
      </view>

      <view class="quick-grid">
        <view class="quick-item" @click="goToHotel">
          <text class="quick-title">房型预订</text>
          <text class="quick-desc">浏览可订房型</text>
        </view>
        <view class="quick-item" @click="goToService">
          <text class="quick-title">酒店服务</text>
          <text class="quick-desc">查看服务信息</text>
        </view>
        <view class="quick-item" @click="goToOrders">
          <text class="quick-title">我的订单</text>
          <text class="quick-desc">查看预订状态</text>
        </view>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { banner, hotel } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'

const hotelId = ref(hotel.DEFAULT_HOTEL_ID || 1)
const hotelInfo = ref({
  id: hotelId.value,
  name: '七天酒店',
  address: '欢迎入住七天酒店',
  phone: '',
  lowestPrice: 0,
  latitude: 30.56,
  longitude: 114.28,
  serviceTime: '24小时服务',
  checkInTime: '14:00',
  checkOutTime: '12:00',
  intro: '',
})
const bannerImages = ref(['https://dummyimage.com/1200x720/f3ede5/8b6b47&text=GoHotel'])
const notices = ref([])
const checkInDate = ref('')
const checkOutDate = ref('')

const nights = computed(() => {
  const start = new Date(checkInDate.value)
  const end = new Date(checkOutDate.value)
  const diff = end.getTime() - start.getTime()
  return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

const hotelIntroText = computed(() => {
  return String(hotelInfo.value.intro || '').replace(/<[^>]+>/g, '').trim() || '欢迎入住七天酒店，享受舒适便捷的入住体验。'
})

const formatDate = (date) => {
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const ensureDates = () => {
  if (checkInDate.value && checkOutDate.value) {
    return
  }

  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  checkInDate.value = formatDate(today)
  checkOutDate.value = formatDate(tomorrow)
}

const loadData = async () => {
  try {
    const [hotelRes, bannerRes] = await Promise.all([
      hotel.getHotelDetail(hotelId.value),
      banner.getActiveBanners().catch(() => []),
    ])

    hotelInfo.value = hotelRes
    notices.value = hotelRes.notices || []

    if (Array.isArray(bannerRes) && bannerRes.length > 0) {
      bannerImages.value = bannerRes.map((item) => item.image_url).filter(Boolean)
    }
  } catch (error) {
    console.error('load home data failed:', error)
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

const goToHotel = () => {
  uni.navigateTo({
    url: `/pages/hotel/hotel?checkIn=${checkInDate.value}&checkOut=${checkOutDate.value}`,
  })
}

const goToService = () => {
  uni.switchTab({ url: '/pages/service/service' })
}

const goToOrders = () => {
  const token = uni.getStorageSync(TOKEN_KEY)
  if (!token) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }

  uni.navigateTo({ url: '/pages/orders/orders' })
}

const callHotel = () => {
  if (!hotelInfo.value.phone) {
    uni.showToast({ title: '暂未配置联系电话', icon: 'none' })
    return
  }

  uni.makePhoneCall({ phoneNumber: hotelInfo.value.phone })
}

const openMap = () => {
  uni.openLocation({
    latitude: Number(hotelInfo.value.latitude),
    longitude: Number(hotelInfo.value.longitude),
    name: hotelInfo.value.name,
    address: hotelInfo.value.address,
  })
}

onLoad((options) => {
  if (options?.hotelId) {
    hotelId.value = Number(options.hotelId)
  }
  ensureDates()
})

onShow(() => {
  ensureDates()
  loadData()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
}

.hero {
  position: relative;
  height: 480rpx;
}

.hero-swiper,
.hero-image {
  width: 100%;
  height: 100%;
}

.hero-mask {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(17, 24, 39, 0.1) 0%, rgba(17, 24, 39, 0.58) 100%);
}

.hero-content {
  position: absolute;
  left: 32rpx;
  right: 32rpx;
  bottom: 32rpx;
  color: #fff;
}

.hotel-name {
  display: block;
  font-size: 40rpx;
  font-weight: 700;
}

.hotel-address {
  display: block;
  margin-top: 12rpx;
  font-size: 26rpx;
  opacity: 0.92;
}

.hero-actions {
  display: flex;
  gap: 20rpx;
  margin-top: 24rpx;
}

.ghost-btn {
  margin: 0;
  padding: 0 28rpx;
  background: rgba(255, 255, 255, 0.18);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.32);
  border-radius: 999rpx;
}

.content {
  padding: 24rpx;
  margin-top: -24rpx;
}

.card {
  padding: 28rpx;
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.06);
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
  margin-bottom: 24rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1f2937;
}

.card-subtitle {
  font-size: 24rpx;
  color: #6b7280;
}

.date-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20rpx;
}

.date-box {
  padding: 24rpx;
  background: #f9fafb;
  border: 1px solid #eef2f7;
  border-radius: 20rpx;
}

.date-label {
  display: block;
  font-size: 24rpx;
  color: #6b7280;
}

.date-value {
  display: block;
  margin-top: 10rpx;
  font-size: 30rpx;
  font-weight: 600;
  color: #111827;
}

.booking-summary,
.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 16rpx;
  margin-top: 20rpx;
  font-size: 26rpx;
  color: #4b5563;
}

.primary-btn {
  margin-top: 28rpx;
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
  border-radius: 999rpx;
}

.notice-item {
  padding: 18rpx 0;
  border-bottom: 1px solid #f1f5f9;
}

.notice-item:last-child {
  border-bottom: 0;
  padding-bottom: 0;
}

.notice-title {
  font-size: 26rpx;
  color: #374151;
  line-height: 1.6;
}

.info-label {
  color: #6b7280;
}

.info-value {
  flex: 1;
  text-align: right;
  color: #111827;
}

.hotel-intro {
  display: block;
  margin-top: 24rpx;
  font-size: 26rpx;
  color: #4b5563;
  line-height: 1.7;
}

.quick-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 20rpx;
  padding-bottom: 36rpx;
}

.quick-item {
  padding: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.quick-title {
  display: block;
  font-size: 28rpx;
  font-weight: 600;
  color: #1f2937;
}

.quick-desc {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #6b7280;
  line-height: 1.5;
}
</style>
