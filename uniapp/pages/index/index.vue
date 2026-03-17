<template>
  <view class="container">
    <!-- 沉浸式导航栏 -->
    <TnNavbar 
      fixed 
      :bottom-shadow="false" 
      :bg-color="navbarBgColor" 
      :placeholder="false" 
      class="custom-navbar"
    >
      <template #back>
        <view class="navbar-left" v-if="scrollTop > 100">
          <TnIcon name="location-fill" color="#333" size="32" />
          <text class="hotel-name">{{ hotelData?.name || '' }}</text>
        </view>
      </template>
    </TnNavbar>

    <!-- 顶部轮播区域 -->
    <view class="hero-section">
      <swiper
        class="banner-swiper"
        :indicator-dots="true"
        :autoplay="true"
        :interval="3000"
        :duration="500"
        indicator-color="rgba(255,255,255,0.3)"
        indicator-active-color="#C29D71"
        @change="handleBannerChange"
      >
        <swiper-item v-for="(item, index) in bannerImages" :key="index">
          <view class="swiper-item">
            <image class="banner-img" :src="item" mode="aspectFill" @error="handleImageError" @load="handleImageLoad"></image>
            <view class="image-gradient"></view>
          </view>
        </swiper-item>
      </swiper>
      
      <!-- 酒店信息浮层 -->
      <view class="hotel-info-overlay" v-if="hotelData">
        <text class="overlay-title">{{ hotelData.name }} · Premium</text>
        <view class="overlay-location" @click="openMap">
          <TnIcon name="location-fill" color="#fff" size="32" />
          <text class="location-text">{{ hotelData.address }}</text>
        </view>
      </view>
    </view>

    <!-- 主要内容区 -->
    <view class="main-content">
      <!-- 预订卡片 -->
      <view class="booking-card">
        <view class="date-selection" @click="openDatePicker">
          <view class="date-block">
            <text class="label">入住</text>
            <view class="date-value">
              <text class="date-text">{{ formatMonthDay(checkInDate) }}</text>
            </view>
            <text class="weekday">{{ getWeekDay(checkInDate) }}</text>
          </view>
          
          <view class="nights-divider">
            <view class="nights-badge">{{ nights }}晚</view>
          </view>
          
          <view class="date-block checkout">
            <text class="label">离店</text>
            <view class="date-value">
              <text class="date-text">{{ formatMonthDay(checkOutDate) }}</text>
            </view>
            <text class="weekday">{{ getWeekDay(checkOutDate) }}</text>
          </view>
        </view>

        <view class="search-btn-wrapper">
          <TnButton
            shape="round"
            size="lg"
            width="100%"
            height="104rpx"
            bg-color="#2A2A2A"
            text-color="#D4B184"
            @click.stop="handleBooking"
            :shadow="false"
          >
            <view class="btn-content">
              <text class="btn-text">立即预订</text>
            </view>
          </TnButton>
        </view>
      </view>

      <!-- 功能入口 -->
      <view class="feature-section">
        <view class="feature-item" @click="navigateTo('vip')">
          <view class="icon-box">
            <TnIcon name="vip-diamond" color="#C29D71" size="48" />
          </view>
          <text class="feature-name">会员权益</text>
        </view>
        <view class="feature-item" @click="navigateTo('points')">
          <view class="icon-box">
            <TnIcon name="gift" color="#C29D71" size="48" />
          </view>
          <text class="feature-name">积分商城</text>
        </view>
        <view class="feature-item" @click="navigateTo('cinema')">
          <view class="icon-box">
            <TnIcon name="video" color="#C29D71" size="48" />
          </view>
          <text class="feature-name">影音服务</text>
        </view>
        <view class="feature-item" @click="showMoreOptions">
          <view class="icon-box">
            <TnIcon name="menu-circle" color="#C29D71" size="48" />
          </view>
          <text class="feature-name">更多服务</text>
        </view>
      </view>

      <!-- 促销活动 -->
      <view class="promotion-section" v-if="promotionImage">
        <view class="promotion-card" @click="goToPromotion">
          <image :src="promotionImage" mode="aspectFill" class="promo-img"></image>
          <view class="promo-content">
            <view class="promo-tag">限时特惠</view>
            <text class="promo-title">尊享会员专享礼遇</text>
            <text class="promo-desc">预订立减 · 延迟退房</text>
          </view>
        </view>
      </view>

      <!-- 底部安全区 -->
      <view class="safe-area-bottom"></view>
    </view>

    <!-- 日期选择弹窗 -->
    <TnPopup
      v-model="showCalendar"
      open-direction="bottom"
      :radius="32"
      height="85vh"
      :safe-area-inset-bottom="true"
      bg-color="#fff"
    >
      <view class="calendar-wrapper">
        <view class="calendar-header">
          <text class="title">选择日期</text>
          <view class="close-icon" @click="closeCalendar">
            <TnIcon name="close" size="36" color="#999" />
          </view>
        </view>
        
        <view class="calendar-content">
          <TnCalendar
            v-if="calendarVisible"
            :key="calendarKey"
            v-model="selectedDateRange"
            mode="range"
            :min-date="minSelectDate"
            :max-date="maxSelectDate"
            active-bg-color="#C29D71"
            active-text-color="#fff"
            range-bg-color="rgba(194, 157, 113, 0.1)"
            range-text-color="#C29D71"
            @change="onDateRangeChange"
          />
        </view>

        <view class="calendar-footer">
          <TnButton
            shape="round"
            size="lg"
            width="100%"
            height="88rpx"
            bg-color="#2A2A2A"
            text-color="#D4B184"
            :disabled="!tempCheckIn || !tempCheckOut"
            @click="confirmDateSelection"
          >
            确认选择 ({{ tempNights }}晚)
          </TnButton>
        </view>
      </view>
    </TnPopup>
  </view>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { onLoad, onPageScroll } from '@dcloudio/uni-app'
import { hotel, banner } from '@/api/index.js'
import { API_BASE_URL } from '@/config/api.config.js'
import TnSwiper from '@/uni_modules/tuniaoui-vue3/components/swiper/src/swiper.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnPopup from '@/uni_modules/tuniaoui-vue3/components/popup/src/popup.vue'
import TnCalendar from '@/uni_modules/tuniaoui-vue3/components/calendar/src/calendar.vue'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'

// 状态管理
const loading = ref(false)
const scrollTop = ref(0)
const navbarBgColor = computed(() => {
  const opacity = Math.min(scrollTop.value / 100, 1)
  return `rgba(255, 255, 255, ${opacity})`
})

// 数据
const hotelId = ref(1)
const hotelData = ref(null)
const bannerImages = ref([])
const bannerLinks = ref([])
const promotionImage = ref('')
const currentSwiperIndex = ref(0)

// 日期相关
const checkInDate = ref(null)
const checkOutDate = ref(null)
const showCalendar = ref(false)
const calendarVisible = ref(false)
const calendarKey = ref(0)
const selectedDateRange = ref([])
const tempCheckIn = ref(null)
const tempCheckOut = ref(null)

// 计算属性
const nights = computed(() => {
  if (!checkInDate.value || !checkOutDate.value) return 1
  const diff = checkOutDate.value.getTime() - checkInDate.value.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

const tempNights = computed(() => {
  if (!tempCheckIn.value || !tempCheckOut.value) return 0
  const diff = tempCheckOut.value.getTime() - tempCheckIn.value.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

const minSelectDate = computed(() => formatDateString(new Date()))
const maxSelectDate = computed(() => {
  const date = new Date()
  date.setMonth(date.getMonth() + 6)
  return formatDateString(date)
})

const formatMonthDay = (date) => {
  if (!date) return '--'
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${month}月${day}日`
}

// 格式化工具
const formatDateString = (date) => {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

const formatDay = (date) => {
  return date ? String(date.getDate()) : '--'
}

const formatMonth = (date) => {
  return date ? `${date.getMonth() + 1}月` : '--'
}

const getWeekDay = (date) => {
  if (!date) return ''
  const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return weekDays[date.getDay()]
}

// 初始化
const initDates = () => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  checkInDate.value = today
  
  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  checkOutDate.value = tomorrow
}

const normalizeAssetUrl = (url) => {
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  if (url.startsWith('//')) return `https:${url}`
  const firstSegment = String(url).split('/')[0]
  if (firstSegment.includes('.') || firstSegment.includes(':')) return `http://${url}`
  if (url.startsWith('/')) return `${API_BASE_URL}${url}`
  return `${API_BASE_URL}/${url}`
}

const loadData = async () => {
  loading.value = true
  try {
    // 获取横幅数据
    const bannerRes = await banner.getActiveBanners()
    
    // 处理Banner数据 - 后端返回 { success: true, data: [...] }
    const bannerData = bannerRes?.data || (bannerRes?.success ? bannerRes.data : bannerRes)
    
    if (bannerData && Array.isArray(bannerData) && bannerData.length > 0) {
      bannerImages.value = bannerData
        .map(item => normalizeAssetUrl(item?.image_url))
        .filter(Boolean)
      bannerLinks.value = bannerData.map(item => item?.link_url || '')
      promotionImage.value = normalizeAssetUrl(bannerData?.[0]?.image_url)
    }
    
    // 酒店数据暂时使用固定值，等待后端提供酒店详情接口
    hotelData.value = {
      name: 'GoHotel精品酒店',
      address: '北京市朝阳区建国门外大街1号',
      latitude: 39.9042,
      longitude: 116.4074
    }
  } catch (error) {
    console.error('Data loading failed:', error)
    throw error
  } finally {
    loading.value = false
  }
}

const handleBannerChange = (e) => {
  currentSwiperIndex.value = e.detail.current
}

const handleImageLoad = (e) => {
  // 图片加载成功
}

const handleImageError = (e) => {
  console.error('Image load failed:', e)
}

const handleBannerClick = (index) => {
  const linkUrl = bannerLinks.value[index]
  if (linkUrl) {
    // 如果是小程序页面路径
    if (linkUrl.startsWith('/pages/')) {
      uni.navigateTo({ url: linkUrl })
    } else {
      // 如果是外部链接，使用 webview
      uni.navigateTo({
        url: `/pages/webview/webview?url=${encodeURIComponent(linkUrl)}`
      })
    }
  }
}

// 生命周期
onLoad((options) => {
  if (options?.hotelId) {
    hotelId.value = options.hotelId
  }
  initDates()
  loadData()
})

onPageScroll((e) => {
  scrollTop.value = e.scrollTop
})

// 交互方法
const openDatePicker = () => {
  tempCheckIn.value = checkInDate.value
  tempCheckOut.value = checkOutDate.value
  
  if (checkInDate.value && checkOutDate.value) {
    selectedDateRange.value = [
      formatDateString(checkInDate.value),
      formatDateString(checkOutDate.value)
    ]
  } else {
    selectedDateRange.value = []
  }
  
  showCalendar.value = true
  setTimeout(() => {
    calendarKey.value++
    calendarVisible.value = true
  }, 200)
}

const closeCalendar = () => {
  showCalendar.value = false
  calendarVisible.value = false
}

const onDateRangeChange = (value) => {
  if (Array.isArray(value) && value.length === 2) {
    tempCheckIn.value = parseDate(value[0])
    tempCheckOut.value = parseDate(value[1])
  } else {
    tempCheckIn.value = parseDate(value[0])
    tempCheckOut.value = null
  }
}

const parseDate = (dateStr) => {
  if (!dateStr) return null
  const normalizedStr = dateStr.replace(/\//g, '-')
  const [year, month, day] = normalizedStr.split('-').map(Number)
  return new Date(year, month - 1, day)
}

const confirmDateSelection = () => {
  if (tempCheckIn.value && tempCheckOut.value) {
    checkInDate.value = tempCheckIn.value
    checkOutDate.value = tempCheckOut.value
    closeCalendar()
  }
}

const handleBooking = () => {
  const token = uni.getStorageSync('gohotel_token')
  if (!token) {
    uni.navigateTo({ url: '/pages/login/login' })
    return
  }
  
  uni.navigateTo({
    url: `/pages/hotel/hotel?id=${hotelId.value}&checkIn=${formatDateString(checkInDate.value)}&checkOut=${formatDateString(checkOutDate.value)}`
  })
}

const openMap = () => {
  if (hotelData.value?.latitude && hotelData.value?.longitude) {
    uni.openLocation({
      latitude: Number(hotelData.value.latitude),
      longitude: Number(hotelData.value.longitude),
      name: hotelData.value.name,
      address: hotelData.value.address
    })
  }
}

const navigateTo = (path) => {
  uni.navigateTo({ url: `/pages/${path}/${path}` })
}

const goToPromotion = () => {
  uni.navigateTo({ url: '/pages/promotion/promotion' })
}

const showMoreOptions = () => {
  uni.showActionSheet({
    itemList: ['联系客服', '酒店介绍', '分享行程'],
    success: (res) => {
      if (res.tapIndex === 0) {
        uni.makePhoneCall({ phoneNumber: '400-888-8888' })
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background-color: #F9F9F9;
  box-sizing: border-box;
}

/* 导航栏 */
.navbar-left {
  display: flex;
  align-items: center;
  gap: 6rpx;
  width: 500rpx;
  
  .hotel-name {
    font-size: 32rpx;
    font-weight: 600;
    color: #333;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
  }
}

/* Hero Section */
.hero-section {
  position: relative;
  height: 560rpx;
  overflow: hidden;
  width: 100%;
  
  .banner-swiper {
    width: 100%;
    height: 100%;
  }
  
  .swiper-item {
    width: 100%;
    height: 100%;
    position: relative;
    
    .banner-img {
      width: 100%;
      height: 100%;
    }
  }
  
  .image-gradient {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(180deg, rgba(0,0,0,0.2) 0%, rgba(0,0,0,0) 45%, rgba(0,0,0,0.55) 100%);
  }
  
  .hotel-info-overlay {
    position: absolute;
    left: 32rpx;
    top: 380rpx;
    z-index: 2;
    
    .overlay-title {
      font-size: 48rpx;
      font-weight: bold;
      color: #fff;
      margin-bottom: 24rpx;
      display: block;
    }
    
    .overlay-location {
      display: inline-flex;
      align-items: center;
      gap: 12rpx;
      background: rgba(255,255,255,0.22);
      backdrop-filter: blur(10px);
      padding: 12rpx 20rpx;
      border-radius: 999rpx;
      
      .location-text {
        font-size: 24rpx;
        color: #fff;
        max-width: 440rpx;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

/* 主要内容区 */
.main-content {
  position: relative;
  z-index: 3;
  padding: 32rpx;
  display: flex;
  flex-direction: column;
  gap: 40rpx;
  margin-top: -80rpx;
  
  .booking-card {
    background: #fff;
    border-radius: 44rpx;
    padding: 32rpx;
    box-shadow: 0 20rpx 60rpx rgba(0,0,0,0.05);
    border: 1rpx solid #F0F0F0;
    height: 420rpx;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    box-sizing: border-box;
    
    .date-selection {
      display: flex;
      justify-content: space-between;
      align-items: center;
      height: 176rpx;
      
      .date-block {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 12rpx;
        
        &.checkout {
          align-items: flex-end;
        }
        
        .label {
          font-size: 22rpx;
          color: #999;
        }
        
        .date-value {
          .date-text {
            font-size: 36rpx;
            font-weight: bold;
            color: #333;
          }
        }
        
        .weekday {
          font-size: 22rpx;
          color: #999;
        }
      }
      
      .nights-divider {
        display: flex;
        align-items: center;
        justify-content: center;
        
        .nights-badge {
          font-size: 22rpx;
          font-weight: 600;
          color: #C29D71;
          background: rgba(194,157,113,0.12);
          border: 1rpx solid rgba(194,157,113,0.35);
          padding: 8rpx 24rpx;
          border-radius: 999rpx;
        }
      }
    }
    
    .search-btn-wrapper {
      width: 100%;
    }
    
    .btn-content {
      display: flex;
      align-items: center;
      justify-content: center;
      
      .btn-text {
        font-size: 32rpx;
        font-weight: bold;
        letter-spacing: 2rpx;
      }
    }
  }
}

/* 金刚区 */
.feature-section {
  display: flex;
  justify-content: space-between;
  padding: 24rpx 16rpx;
  background: #fff;
  border-radius: 44rpx;
  border: 1rpx solid #F0F0F0;
  height: 184rpx;
  align-items: center;
  box-sizing: border-box;
  
  .feature-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16rpx;
    width: 144rpx;
    
    .icon-box {
      width: 44rpx;
      height: 44rpx;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .feature-name {
      font-size: 22rpx;
      color: #666;
    }
  }
}

/* 促销区域 */
.promotion-section {
  .promotion-card {
    position: relative;
    height: 320rpx;
    border-radius: 44rpx;
    overflow: hidden;
    background: #111;
    
    .promo-img {
      width: 100%;
      height: 100%;
    }
    
    .promo-content {
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      padding: 36rpx 32rpx;
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      background: linear-gradient(90deg, rgba(0,0,0,0.6) 0%, rgba(0,0,0,0.05) 70%, transparent 100%);
      box-sizing: border-box;
      
      .promo-tag {
        align-self: flex-start;
        background: #C29D71;
        color: #fff;
        font-size: 22rpx;
        font-weight: 600;
        padding: 6rpx 24rpx;
        border-radius: 16rpx;
        margin-bottom: 16rpx;
      }
      
      .promo-title {
        font-size: 36rpx;
        font-weight: bold;
        color: #fff;
        margin-bottom: 12rpx;
      }
      
      .promo-desc {
        font-size: 24rpx;
        color: rgba(255,255,255,0.85);
      }
    }
  }
}

/* 日历弹窗 */
.calendar-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  
  .calendar-header {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    padding: 32rpx;
    border-bottom: 1rpx solid #f5f5f5;
    
    .title {
      font-size: 32rpx;
      font-weight: 600;
    }
    
    .close-icon {
      position: absolute;
      right: 32rpx;
      top: 50%;
      transform: translateY(-50%);
      padding: 10rpx;
    }
  }
  
  .calendar-content {
    flex: 1;
    overflow-y: auto;
  }
  
  .calendar-footer {
    padding: 24rpx 32rpx;
    padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
    border-top: 1rpx solid #f5f5f5;
  }
}

.safe-area-bottom {
  height: env(safe-area-inset-bottom);
}
</style>
