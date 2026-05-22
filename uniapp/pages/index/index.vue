<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="custom-nav" :style="{ paddingTop: statusBarHeight + 'px' }">
      <view class="nav-content">
        <text class="nav-title">{{ hotelInfo.name }}</text>
      </view>
    </view>

    <!-- Banner 区域 -->
    <view class="hero">
      <swiper class="hero-swiper" circular autoplay interval="5000" duration="1000">
        <swiper-item v-for="(image, index) in bannerImages" :key="index">
          <image class="hero-image" :src="image" mode="aspectFill" />
        </swiper-item>
      </swiper>
      <view class="hero-mask"></view>
      <view class="hero-content" :style="{ bottom: '80rpx' }">
        <view class="hotel-info-tag">
          <text class="tag-text">优选酒店</text>
        </view>
        <text class="hotel-name">{{ hotelInfo.name }}</text>
        <view class="hotel-location" @click="goToMap">
          <image class="icon-location" src="/static/icons/location-white.png" mode="aspectFit" />
          <text class="address-text">{{ hotelInfo.address }}</text>
        </view>
      </view>
    </view>

    <view class="content">
      <!-- 核心操作区：日期选择 -->
      <view class="card date-card premium-card">
        <view class="date-selection" @click="showCalendar = true">
          <view class="date-item">
            <text class="date-label">入住</text>
            <view class="date-val-wrap">
              <text class="date-day">{{ checkInDateDisplay.day }}</text>
              <text class="date-month">{{ checkInDateDisplay.month }}月</text>
            </view>
          </view>
          <view class="date-duration">
            <view class="duration-line"></view>
            <text class="duration-text">{{ nights }}晚</text>
            <view class="duration-line"></view>
          </view>
          <view class="date-item">
            <text class="date-label">离店</text>
            <view class="date-val-wrap">
              <text class="date-day">{{ checkOutDateDisplay.day }}</text>
              <text class="date-month">{{ checkOutDateDisplay.month }}月</text>
            </view>
          </view>
        </view>
        <button class="primary-btn premium-button" @click="goToHotel">立即预订</button>
      </view>

      <!-- 金刚区 (Quick Actions) -->
      <view class="quick-grid">
        <view class="quick-item" @click="goToService">
          <view class="quick-icon-wrap service">
            <image class="quick-icon" src="/static/icons/service.png" mode="aspectFit" />
          </view>
          <text class="quick-title">酒店服务</text>
        </view>
        <view class="quick-item" @click="goToOrders">
          <view class="quick-icon-wrap orders">
            <image class="quick-icon" src="/static/icons/order.png" mode="aspectFit" />
          </view>
          <text class="quick-title">我的订单</text>
        </view>
        <view class="quick-item" @click="callHotel">
          <view class="quick-icon-wrap contact">
            <image class="quick-icon" src="/static/icons/phone.png" mode="aspectFit" />
          </view>
          <text class="quick-title">联系前台</text>
        </view>
        <view class="quick-item" @click="goToMap">
          <view class="quick-icon-wrap map">
            <image class="quick-icon" src="/static/icons/map.png" mode="aspectFit" />
          </view>
          <text class="quick-title">地图导航</text>
        </view>
      </view>

      <!-- 智能推荐 -->
      <view v-if="recommendedRooms.length" class="section">
        <view class="section-header">
          <view class="header-left">
            <text class="section-title">智能推荐</text>
            <text class="section-subtitle">为你精选最适合的房型</text>
          </view>
          <text class="more-link" @click="goToHotel">查看全部</text>
        </view>
        <scroll-view scroll-x class="recommend-scroll">
          <view class="recommend-list">
            <view
              v-for="room in recommendedRooms"
              :key="room.id"
              class="recommend-card premium-card"
              @click="goToRecommendedRoom(room)"
            >
              <image class="recommend-image" :src="room.image" mode="aspectFill" />
              <view class="recommend-info">
                <text class="room-name">{{ room.name }}</text>
                <view class="price-wrap">
                  <text class="price-symbol">¥</text>
                  <text class="price-value">{{ room.price }}</text>
                  <text class="price-unit">/晚起</text>
                </view>
                <view class="recommend-tags">
                  <text v-for="tag in room.recommendationTags.slice(0, 2)" :key="tag" class="tag">
                    {{ tag }}
                  </text>
                </view>
              </view>
            </view>
          </view>
        </scroll-view>
      </view>

      <!-- 酒店信息 -->
      <view class="section">
        <view class="section-header">
          <text class="section-title">酒店详情</text>
        </view>
        <view class="card info-card premium-card">
          <view class="info-row">
            <text class="info-label">服务时间</text>
            <text class="info-value">{{ hotelInfo.serviceTime }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">入住/退房</text>
            <text class="info-value">{{ hotelInfo.checkInTime }} / {{ hotelInfo.checkOutTime }}</text>
          </view>
          <view class="hotel-intro-box">
            <text class="hotel-intro-text">{{ hotelIntroText }}</text>
          </view>
        </view>
      </view>
    </view>
    <HotelCalendar v-model:show="showCalendar" :check-in="checkInDate" :check-out="checkOutDate" @confirm="onCalendarConfirm" />
  </view>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { banner, hotel } from '@/api/index.js'
import { API_BASE_URL, TOKEN_KEY } from '@/config/api.config.js'
import HotelCalendar from '@/components/hotel-calendar/hotel-calendar.vue'

const statusBarHeight = ref(44)
const showCalendar = ref(false)
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
  coverImages: [],
})
const bannerImages = ref(['https://dummyimage.com/1200x720/f3ede5/8b6b47&text=GoHotel'])
const notices = ref([])
const recommendedRooms = ref([])
const checkInDate = ref('')
const checkOutDate = ref('')
const HOTEL_TAB_CONTEXT_KEY = 'gohotel_hotel_tab_context'

// 计算日期显示格式
const checkInDateDisplay = computed(() => {
  if (!checkInDate.value) return { month: '--', day: '--' }
  const d = new Date(checkInDate.value)
  return { month: d.getMonth() + 1, day: d.getDate() }
})

const checkOutDateDisplay = computed(() => {
  if (!checkOutDate.value) return { month: '--', day: '--' }
  const d = new Date(checkOutDate.value)
  return { month: d.getMonth() + 1, day: d.getDate() }
})

const nights = computed(() => {
  if (!checkInDate.value || !checkOutDate.value) return 1
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
  if (checkInDate.value && checkOutDate.value) return

  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  checkInDate.value = formatDate(today)
  checkOutDate.value = formatDate(tomorrow)
}

const loadData = async () => {
  try {
    const [hotelRes, bannerRes, recommendationRes] = await Promise.all([
      hotel.getHotelDetail(hotelId.value),
      banner.getActiveBanners().catch(() => []),
      hotel.getRecommendedRooms({ limit: 5 }).catch(() => []),
    ])

    hotelInfo.value = hotelRes
    notices.value = hotelRes.notices || []
    recommendedRooms.value = recommendationRes

    const normalizedBanners = (Array.isArray(bannerRes) ? bannerRes : [])
      .map((item) => item?.image_url || item?.temp_url || item?.image)
      .filter(Boolean)

    if (normalizedBanners.length > 0) {
      bannerImages.value = normalizedBanners
    }
  } catch (error) {
    console.error('load home data failed:', error)
  }
}

const onCalendarConfirm = (dates) => {
  checkInDate.value = dates.checkIn
  checkOutDate.value = dates.checkOut
}

const goToHotel = () => {
  uni.setStorageSync(HOTEL_TAB_CONTEXT_KEY, {
    hotelId: hotelId.value,
    checkIn: checkInDate.value,
    checkOut: checkOutDate.value,
  })
  uni.switchTab({ url: '/pages/hotel/hotel' })
}

const goToRecommendedRoom = (room) => {
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${room.id}&checkIn=${checkInDate.value}&checkOut=${checkOutDate.value}`,
  })
}

const goToService = () => uni.switchTab({ url: '/pages/service/service' })
const goToOrders = () => {
  const token = uni.getStorageSync(TOKEN_KEY)
  token ? uni.navigateTo({ url: '/pages/orders/orders' }) : uni.navigateTo({ url: '/pages/login/login' })
}
const callHotel = () => {
  const phone = String(hotelInfo.value.phone || '').trim()
  phone ? uni.makePhoneCall({ phoneNumber: phone }) : uni.showToast({ title: '暂未配置联系电话', icon: 'none' })
}
const goToMap = () => {
  const { latitude, longitude, name, address } = hotelInfo.value
  uni.openLocation({ latitude: Number(latitude), longitude: Number(longitude), name, address })
}

onLoad((options) => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
  if (options?.hotelId) hotelId.value = Number(options.hotelId)
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
  padding-bottom: 60rpx;
}

/* 自定义导航栏 */
.custom-nav {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  background: rgba(255, 255, 255, 0);
  transition: background 0.3s;

  .nav-content {
    height: 88rpx;
    display: flex;
    align-items: center;
    padding: 0 32rpx;
  }

  .nav-title {
    font-size: 34rpx;
    font-weight: 600;
    color: #fff;
    opacity: 0; /* 滚动时根据需求显示 */
  }
}

.hero {
  position: relative;
  height: 600rpx;

  .hero-swiper, .hero-image {
    width: 100%;
    height: 100%;
  }

  .hero-mask {
    position: absolute;
    inset: 0;
    background: linear-gradient(180deg, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0.6) 100%);
  }

  .hero-content {
    position: absolute;
    left: 40rpx;
    right: 40rpx;
    color: #fff;

    .hotel-info-tag {
      display: inline-block;
      padding: 4rpx 16rpx;
      background: #c9a977;
      border-radius: 8rpx;
      margin-bottom: 16rpx;

      .tag-text {
        font-size: 20rpx;
        font-weight: 600;
      }
    }

    .hotel-name {
      display: block;
      font-size: 48rpx;
      font-weight: 700;
      margin-bottom: 12rpx;
      text-shadow: 0 2rpx 4rpx rgba(0,0,0,0.3);
    }

    .hotel-location {
      display: flex;
      align-items: center;
      opacity: 0.9;

      .icon-location {
        width: 24rpx;
        height: 24rpx;
        margin-right: 8rpx;
      }

      .address-text {
        font-size: 24rpx;
      }
    }
  }
}

.content {
  padding: 0 30rpx;
  margin-top: -60rpx;
  position: relative;
  z-index: 10;
}

.card {
  padding: 40rpx;
  margin-bottom: 30rpx;
}

/* 日期选择卡片 */
.date-card {
  .date-selection {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 40rpx;

    .date-item {
      .date-label {
        font-size: 24rpx;
        color: #6b7280;
        margin-bottom: 10rpx;
        display: block;
      }

      .date-val-wrap {
        display: flex;
        align-items: baseline;

        .date-day {
          font-size: 44rpx;
          font-weight: 700;
          color: #111827;
          margin-right: 4rpx;
        }

        .date-month {
          font-size: 26rpx;
          color: #111827;
        }
      }
    }

    .date-duration {
      display: flex;
      flex-direction: column;
      align-items: center;

      .duration-line {
        width: 60rpx;
        height: 1rpx;
        background: #eee;
      }

      .duration-text {
        font-size: 22rpx;
        color: #c9a977;
        margin: 4rpx 0;
        padding: 2rpx 12rpx;
        border: 1rpx solid #c9a977;
        border-radius: 20rpx;
      }
    }
  }

  .primary-btn {
    height: 90rpx;
    font-size: 32rpx;
  }
}

/* 金刚区 */
.quick-grid {
  display: flex;
  justify-content: space-between;
  padding: 10rpx 10rpx 40rpx;

  .quick-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 25%;

    .quick-icon-wrap {
      width: 96rpx;
      height: 96rpx;
      border-radius: 32rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 16rpx;
      box-shadow: 0 8rpx 20rpx rgba(0,0,0,0.05);

      &.service { background: #EBF3FF; }
      &.orders { background: #FFF4E8; }
      &.contact { background: #E8F8F2; }
      &.map { background: #F3EBFF; }

      .quick-icon {
        width: 48rpx;
        height: 48rpx;
      }
    }

    .quick-title {
      font-size: 24rpx;
      color: #4b5563;
      font-weight: 500;
    }
  }
}

/* 章节标题 */
.section {
  margin-bottom: 40rpx;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-end;
    margin-bottom: 24rpx;
    padding: 0 10rpx;

    .header-left {
      .section-title {
        font-size: 36rpx;
        font-weight: 700;
        color: #1f2937;
        display: block;
      }

      .section-subtitle {
        font-size: 22rpx;
        color: #6b7280;
        margin-top: 4rpx;
      }
    }

    .more-link {
      font-size: 24rpx;
      color: #c9a977;
    }
  }
}

/* 推荐滚动列表 */
.recommend-scroll {
  width: 100%;
  white-space: nowrap;

  .recommend-list {
    display: inline-flex;
    padding: 10rpx;
    gap: 24rpx;
  }

  .recommend-card {
    width: 420rpx;
    display: inline-block;
    padding: 0;

    .recommend-image {
      width: 100%;
      height: 260rpx;
    }

    .recommend-info {
      padding: 20rpx;

      .room-name {
        font-size: 28rpx;
        font-weight: 600;
        color: #1f2937;
        display: block;
        margin-bottom: 12rpx;
      }

      .price-wrap {
        margin-bottom: 12rpx;

        .price-symbol {
          font-size: 24rpx;
          color: #E64340;
          font-weight: 700;
        }

        .price-value {
          font-size: 36rpx;
          color: #E64340;
          font-weight: 700;
        }

        .price-unit {
          font-size: 20rpx;
          color: #6b7280;
          margin-left: 4rpx;
        }
      }

      .recommend-tags {
        display: flex;
        gap: 8rpx;

        .tag {
          font-size: 20rpx;
          padding: 4rpx 12rpx;
          background: #F5F5F7;
          color: #6b7280;
          border-radius: 4rpx;
        }
      }
    }
  }
}

/* 详情卡片 */
.info-card {
  .info-row {
    display: flex;
    justify-content: space-between;
    margin-bottom: 24rpx;

    &:last-of-type { margin-bottom: 0; }

    .info-label {
      font-size: 28rpx;
      color: #6b7280;
    }

    .info-value {
      font-size: 28rpx;
      color: #1f2937;
      font-weight: 500;
    }
  }

  .hotel-intro-box {
    margin-top: 30rpx;
    padding-top: 30rpx;
    border-top: 1rpx solid #F2F2F7;

    .hotel-intro-text {
      font-size: 26rpx;
      color: #4b5563;
      line-height: 1.6;
    }
  }
}
</style>
