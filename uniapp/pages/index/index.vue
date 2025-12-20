<template>
  <view class="container">
    <!-- 自定义导航栏 -->
    <TnNavbar fixed :bottom-shadow="false" bg-color="rgba(255, 255, 255, 0.95)" :placeholder="true" class="navbar-glass">
      <template #back>
        <view class="navbar-left">
          <view class="location-icon-wrapper">
            <TnIcon name="location-fill" color="#C29D71" size="32" />
          </view>
          <text class="hotel-name">{{ hotelData?.name || '七天酒店' }}</text>
        </view>
      </template>
    </TnNavbar>

    <!-- 主内容区域 -->
    <scroll-view class="main-content" scroll-y>
      <!-- 加载骨架屏 -->
      <template v-if="loading || bannerLoading">
        <view class="skeleton-banner"></view>
        <view class="skeleton-address">
          <view class="skeleton-line"></view>
        </view>
        <view class="skeleton-booking">
          <view class="skeleton-block"></view>
          <view class="skeleton-block"></view>
          <view class="skeleton-block"></view>
        </view>
      </template>

      <template v-else>
        <!-- 酒店图片轮播 -->
        <view class="hotel-banner">
          <view class="banner-overlay"></view>
          <TnSwiper
            v-model="currentSwiperIndex"
            :data="hotelImages"
            width="100%"
            height="360"
            autoplay
            loop
            indicator
            indicator-type="dot"
            indicator-bg-color="rgba(255,255,255,0.4)"
            indicator-active-bg-color="#C29D71"
          >
            <template #default="{ data }">
              <view class="swiper-item">
                <image class="banner-img" :src="data" mode="aspectFill"></image>
                <view class="image-gradient"></view>
              </view>
            </template>
          </TnSwiper>
          <view class="image-count-badge">
            <TnIcon name="image" color="#fff" size="24" />
            <text>{{ currentSwiperIndex + 1 }}/{{ hotelImages.length }}</text>
          </view>
        </view>

        <!-- 酒店地址 -->
        <view class="address-section" @click="openMap">
          <view class="address-icon-wrapper">
            <TnIcon name="location-fill" color="#C29D71" size="36" class="address-icon" />
          </view>
          <text class="address-text">{{ address }}</text>
          <view class="address-arrow">
            <TnIcon name="right" color="#999" size="28" />
          </view>
        </view>

        <!-- 预订信息 -->
        <view class="booking-card" @click="openDatePicker">
          <view class="card-header">
            <text class="card-title">选择日期</text>
            <view class="edit-icon">
              <TnIcon name="edit" color="#C29D71" size="28" />
            </view>
          </view>

          <view class="booking-dates">
            <view class="date-column">
              <text class="label">{{ checkInLabel }}</text>
              <view class="date-row">
                <text class="date">{{ formatDisplayDate(checkInDate) }}</text>
                <text class="weekday">{{ getWeekDay(checkInDate) }}</text>
              </view>
            </view>

            <view class="nights-badge">
              <view class="arrow-icon">
                <TnIcon name="right" color="#C29D71" size="32" />
              </view>
              <text class="count">{{ nights }}晚</text>
            </view>

            <view class="date-column right">
              <text class="label">{{ checkOutLabel }}</text>
              <view class="date-row">
                <text class="date">{{ formatDisplayDate(checkOutDate) }}</text>
                <text class="weekday">{{ getWeekDay(checkOutDate) }}</text>
              </view>
            </view>
          </view>

          <!-- 搜索框 -->
          <view class="search-wrapper" @click.stop="goToSearch">
             <TnIcon name="search" color="#C29D71" size="32" class="search-icon" />
             <text class="search-placeholder">搜索房型、服务...</text>
             <view class="search-arrow">
               <TnIcon name="right" color="#ccc" size="28" />
             </view>
          </view>

          <!-- 立即预订按钮 -->
          <view class="action-area">
             <TnButton
               shape="round"
               size="xl"
               width="100%"
               height="100rpx"
               bg-color="linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%)"
               text-color="#FFFFFF"
               @click.stop="handleBooking"
             >
               <view class="btn-content">
                 <TnIcon name="calendar" color="#fff" size="32" />
                 <text class="btn-text">立即预订</text>
               </view>
             </TnButton>
          </view>

          <view class="guarantee-bar">
            <view class="guarantee-icon">
              <TnIcon name="check-circle-fill" color="#52C41A" size="28" />
            </view>
            <text>官方渠道 · 低价保证 · 安心入住</text>
          </view>
        </view>

        <!-- 功能入口 -->
        <view class="feature-grid">
          <view class="feature-card cinema" @click="navigateTo('cinema')">
            <view class="feature-bg-icon">
              <TnIcon name="video-fill" color="rgba(194, 157, 113, 0.1)" size="80" />
            </view>
            <view class="feature-content">
              <view class="feature-icon cinema-icon">
                <TnIcon name="video-fill" color="#C29D71" size="44" />
              </view>
              <text class="title">影院足道</text>
              <text class="desc">官方自营</text>
            </view>
            <view class="feature-arrow">
              <TnIcon name="right" color="#ddd" size="28" />
            </view>
          </view>

          <view class="feature-card points" @click="navigateTo('points')">
            <view class="feature-bg-icon">
              <TnIcon name="gift-fill" color="rgba(194, 157, 113, 0.1)" size="80" />
            </view>
            <view class="feature-content">
              <view class="feature-icon points-icon">
                <TnIcon name="gift-fill" color="#C29D71" size="44" />
              </view>
              <text class="title">积分商城</text>
              <text class="desc">体验兑换</text>
            </view>
            <view class="feature-arrow">
              <TnIcon name="right" color="#ddd" size="28" />
            </view>
          </view>

          <view class="feature-card vip" @click="navigateTo('vip')">
            <view class="feature-bg-icon">
              <TnIcon name="vip-fill" color="rgba(194, 157, 113, 0.1)" size="80" />
            </view>
            <view class="feature-content">
              <view class="feature-icon vip-icon">
                <TnIcon name="vip-fill" color="#C29D71" size="44" />
              </view>
              <text class="title">升级会员</text>
              <text class="desc">更多折扣</text>
            </view>
            <view class="feature-arrow">
              <TnIcon name="right" color="#ddd" size="28" />
            </view>
          </view>
        </view>

        <!-- 促销通告 -->
        <view class="notice-section" v-if="true">
          <view class="notice-wrapper">
            <TnNoticeBar
              :data="['🎉 双11特惠活动火热进行中！','🎁 新用户注册立享88元大礼包','✨ 会员专享每日折扣优惠']"
              direction="vertical"
              left-icon="sound"
              bg-color="linear-gradient(135deg, #FFF9F0 0%, #FFF4E5 100%)"
              color="#E67E22"
            />
          </view>
        </view>

        <!-- 促销横幅 -->
        <view class="promotion-banner" @click="goToPromotion">
          <image class="promotion-img" :src="promotionImage" mode="aspectFill"></image>
          <view class="promotion-overlay"></view>
          <view class="promotion-content">
            <view class="promotion-tag">
              <text class="tag-text">限时特惠</text>
            </view>
            <view class="promotion-info">
              <text class="promotion-title">周年庆大促</text>
              <text class="promotion-subtitle">立享5折优惠 ></text>
            </view>
          </view>
        </view>
      </template>

      <!-- 底部占位，避免被tabbar遮挡 -->
      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 日期选择弹窗 -->
    <TnPopup
      v-model="showCalendar"
      open-direction="bottom"
      :radius="32"
      height="80vh"
      :safe-area-inset-bottom="true"
      width="100%"
      @close="calendarVisible = false"
    >
      <view class="calendar-container">
        <!-- 头部 -->
        <view class="calendar-header">
          <text class="title">选择日期</text>
          <view class="close-btn" @click="closeCalendar">
            <TnIcon name="close" size="40" color="#999" />
          </view>
        </view>

        <!-- 移除了自定义的状态栏和星期栏，直接使用组件自带的以避免布局冲突 -->

        <!-- 日期选择区域 - 延迟渲染确保正确初始化 -->
        <view class="calendar-body">
          <TnCalendar
            v-if="calendarVisible"
            :key="calendarKey"
            v-model="selectedDateRange"
            mode="range"
            :min-date="minSelectDate"
            :max-date="maxSelectDate"
            active-bg-color="#C29D71"
            active-text-color="#fff"
            range-bg-color="rgba(194, 157, 113, 0.15)"
            range-text-color="#C29D71"
            range-start-desc="入住"
            range-end-desc="离店"
            @change="onDateRangeChange"
          />
          <!-- 加载提示 -->
          <view v-if="!calendarVisible && showCalendar" class="calendar-loading">
            <text>加载中...</text>
          </view>
        </view>

        <!-- 底部按钮 -->
        <view class="calendar-footer">
          <TnButton
            shape="round"
            size="xl"
            width="100%"
            height="90rpx"
            bg-color="#C29D71"
            text-color="#fff"
            :disabled="!tempCheckIn || !tempCheckOut"
            @click="confirmDateSelection"
          >
            确认选择
          </TnButton>
        </view>
      </view>
    </TnPopup>
  </view>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { hotel, booking, banner } from '@/api/index.js'
import TnSwiper from '@/uni_modules/tuniaoui-vue3/components/swiper/src/swiper.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnPopup from '@/uni_modules/tuniaoui-vue3/components/popup/src/popup.vue'
import TnCalendar from '@/uni_modules/tuniaoui-vue3/components/calendar/src/calendar.vue'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'
import TnNoticeBar from '@/uni_modules/tuniaoui-vue3/components/notice-bar/src/notice-bar.vue'

// 状态栏高度
const statusBarHeight = ref(0)
const navbarHeight = ref(44)

// 酒店信息
const hotelId = ref(1)
const hotelData = ref(null)
const currentSwiperIndex = ref(0)
const hotelImages = ref([
  'https://images.unsplash.com/photo-1631049307264-da0ec9d70304?w=800&q=80',
  'https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?w=800&q=80',
  'https://images.unsplash.com/photo-1590490360182-c33d57733427?w=800&q=80'
])
const promotionImage = ref('https://images.unsplash.com/photo-1571896349842-33c89424de2d?w=800&q=80')
const address = ref('湖北省武汉市硚口区晴川街道沿河大道246号')
const loading = ref(false)

// Banner数据
const banners = ref([])
const bannerLoading = ref(false)

// 日期相关
const checkInDate = ref(null)  // Date对象
const checkOutDate = ref(null) // Date对象
const showCalendar = ref(false)
const calendarVisible = ref(false) // 控制日历组件的实际渲染
const selectedDateRange = ref([])
const tempCheckIn = ref(null)
const tempCheckOut = ref(null)
const calendarKey = ref(0) // 用于强制重新渲染日历组件

// 计算入住天数
const nights = computed(() => {
  if (!checkInDate.value || !checkOutDate.value) return 1
  const diff = checkOutDate.value.getTime() - checkInDate.value.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

// 临时天数（日历选择时）
const tempNights = computed(() => {
  if (!tempCheckIn.value || !tempCheckOut.value) return 0
  const diff = tempCheckOut.value.getTime() - tempCheckIn.value.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

// 计算日期标签（今天入住/明天入住等）
const checkInLabel = computed(() => getDateLabel(checkInDate.value, '入住'))
const checkOutLabel = computed(() => getDateLabel(checkOutDate.value, '离店'))

// 日期选择范围
const minSelectDate = computed(() => formatDateString(new Date()))
const maxSelectDate = computed(() => {
  const date = new Date()
  date.setMonth(date.getMonth() + 3) // 最多预订3个月后
  return formatDateString(date)
})

// 初始化日期
const initDates = () => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  checkInDate.value = today

  const tomorrow = new Date(today)
  tomorrow.setDate(tomorrow.getDate() + 1)
  checkOutDate.value = tomorrow
}

// 格式化日期为字符串 YYYY-MM-DD
const formatDateString = (date) => {
  if (!date) return ''
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

// 格式化显示日期 MM月DD日
const formatDisplayDate = (date) => {
  if (!date) return ''
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  return `${month}月${day}日`
}

// 获取星期几
const getWeekDay = (date) => {
  if (!date) return ''
  const weekDays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return weekDays[date.getDay()]
}

// 获取日期标签
const getDateLabel = (date, suffix) => {
  if (!date) return suffix
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const targetDate = new Date(date)
  targetDate.setHours(0, 0, 0, 0)

  const diff = Math.floor((targetDate - today) / (1000 * 60 * 60 * 24))

  if (diff === 0) return `今天${suffix}`
  if (diff === 1) return `明天${suffix}`
  if (diff === 2) return `后天${suffix}`
  return suffix
}

// 解析日期字符串（支持 YYYY-MM-DD 和 YYYY/MM/DD 两种格式）
const parseDate = (dateStr) => {
  if (!dateStr) return null
  // 将 / 替换为 - 统一处理
  const normalizedStr = dateStr.replace(/\//g, '-')
  const [year, month, day] = normalizedStr.split('-').map(Number)
  if (isNaN(year) || isNaN(month) || isNaN(day)) return null
  return new Date(year, month - 1, day)
}

// 打开日期选择器
const openDatePicker = () => {
  // 初始化临时日期
  tempCheckIn.value = checkInDate.value
  tempCheckOut.value = checkOutDate.value

  // 设置已选日期范围
  if (checkInDate.value && checkOutDate.value) {
    selectedDateRange.value = [
      formatDateString(checkInDate.value),
      formatDateString(checkOutDate.value)
    ]
  } else {
    selectedDateRange.value = []
  }

  // 先打开弹窗
  showCalendar.value = true

  // 延迟渲染日历组件，确保 Popup 动画完成后再初始化
  setTimeout(() => {
    calendarKey.value++
    calendarVisible.value = true
  }, 300)
}

// 日期范围变化
const onDateRangeChange = (value) => {
  if (Array.isArray(value) && value.length === 2) {
    tempCheckIn.value = parseDate(value[0])
    tempCheckOut.value = parseDate(value[1])
  } else if (Array.isArray(value) && value.length === 1) {
    tempCheckIn.value = parseDate(value[0])
    tempCheckOut.value = null
  } else {
    tempCheckIn.value = null
    tempCheckOut.value = null
  }
}

// 确认日期选择
const confirmDateSelection = () => {
  if (tempCheckIn.value && tempCheckOut.value) {
    checkInDate.value = tempCheckIn.value
    checkOutDate.value = tempCheckOut.value
    closeCalendar()

    uni.showToast({
      title: `已选择${nights.value}晚`,
      icon: 'none'
    })
  }
}

// 关闭日历
const closeCalendar = () => {
  showCalendar.value = false
  calendarVisible.value = false
}

// 打开地图
const openMap = () => {
  if (hotelData.value?.latitude && hotelData.value?.longitude) {
    uni.openLocation({
      latitude: hotelData.value.latitude,
      longitude: hotelData.value.longitude,
      name: hotelData.value.name || '七天酒店',
      address: address.value
    })
  } else {
    uni.showToast({
      title: address.value,
      icon: 'none',
      duration: 2000
    })
  }
}

// 更多选项
const showMoreOptions = () => {
  uni.showActionSheet({
    itemList: ['分享给好友', '收藏酒店', '联系客服'],
    success: (res) => {
      if (res.tapIndex === 0) {
        // 分享
      } else if (res.tapIndex === 1) {
        // 收藏
        handleFavorite()
      } else if (res.tapIndex === 2) {
        // 联系客服
        uni.makePhoneCall({
          phoneNumber: '400-123-4567'
        })
      }
    }
  })
}

// 收藏酒店
const handleFavorite = async () => {
  const token = uni.getStorageSync('gohotel_token')
  if (!token) {
    uni.showToast({ title: '请先登录', icon: 'none' })
    return
  }

  try {
    await hotel.favoriteHotel(hotelId.value)
    uni.showToast({ title: '收藏成功', icon: 'success' })
  } catch (error) {
    console.error('收藏失败:', error)
  }
}

// 扫码
const scanQRCode = () => {
  uni.scanCode({
    success: (res) => {
      console.log('扫码结果:', res)
    }
  })
}

// 加载Banner数据
const loadBanners = async () => {
  try {
    bannerLoading.value = true
    const data = await banner.getActiveBanners()
    banners.value = data || []
    // 如果有banner数据，更新轮播图和促销横幅
    if (banners.value.length > 0) {
      // 按sort字段排序
      banners.value.sort((a, b) => a.sort - b.sort)

      // 更新酒店轮播图
      const bannerImages = banners.value.map(item => item.image_url)
      if (bannerImages.length > 0) {
        hotelImages.value = bannerImages
      }
      // 更新促销横幅（使用第一个banner）
      promotionImage.value = banners.value[0].image_url
    }
  } catch (error) {
    console.error('加载Banner数据失败:', error)
  } finally {
    bannerLoading.value = false
  }
}

onLoad((options) => {
  // 获取系统信息
  const systemInfo = uni.getSystemInfoSync()
  statusBarHeight.value = systemInfo.statusBarHeight || 0
  navbarHeight.value = statusBarHeight.value + 44

  // 初始化日期
  initDates()

  // 如果有传入酒店ID
  if (options?.hotelId) {
    hotelId.value = options.hotelId
  }

  // 加载酒店数据
  loadHotelData()

  // 加载Banner数据
  loadBanners()
})

// 加载酒店详情
const loadHotelData = async () => {
  try {
    loading.value = true
    // 尝试调用API，如果失败则使用模拟数据
    try {
      const data = await hotel.getHotelDetail(hotelId.value)
      hotelData.value = data

      if (data.images && data.images.length > 0) {
        hotelImages.value = data.images
      }
      if (data.address) {
        address.value = data.address
      }
    } catch (e) {
      console.log('API调用失败，使用模拟数据')
      // 模拟数据
      hotelData.value = {
        name: '七天酒店',
        address: '湖北省武汉市硚口区晴川街道沿河大道246号',
        latitude: 30.56,
        longitude: 114.28,
        images: [
          'https://images.unsplash.com/photo-1631049307264-da0ec9d70304?w=800&q=80',
          'https://images.unsplash.com/photo-1582719478250-c89cae4dc85b?w=800&q=80',
          'https://images.unsplash.com/photo-1590490360182-c33d57733427?w=800&q=80'
        ]
      }
    }
  } catch (error) {
    console.error('加载酒店数据失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const goToSearch = () => {
  uni.navigateTo({
    url: '/pages/search/search'
  })
}

// 处理预订
const handleBooking = async () => {
  // 检查是否登录
  const token = uni.getStorageSync('gohotel_token')
  if (!token) {
    uni.showModal({
      title: '提示',
      content: '请先登录',
      success: (res) => {
        if (res.confirm) {
          uni.navigateTo({
            url: '/pages/login/login'
          })
        }
      }
    })
    return
  }

  // 跳转到房型选择页面，传递日期参数
  uni.navigateTo({
    url: `/pages/hotel/hotel?id=${hotelId.value}&checkIn=${formatDateString(checkInDate.value)}&checkOut=${formatDateString(checkOutDate.value)}`
  })
}

// 导航到其他页面
const navigateTo = (type) => {
  const routes = {
    cinema: '/pages/cinema/cinema',
    points: '/pages/points/points',
    vip: '/pages/vip/vip'
  }

  if (routes[type]) {
    uni.navigateTo({
      url: routes[type]
    })
  }
}

// 查看促销详情
const goToPromotion = () => {
  uni.navigateTo({
    url: '/pages/promotion/promotion'
  })
}
</script>

<style lang="scss" scoped>
.container {
  width: 100%;
  height: 100vh;
  background: linear-gradient(180deg, #F8F8F8 0%, #FAFAFA 100%);
}

/* 导航栏玻璃效果 */
.navbar-glass {
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
}

/* 导航栏左侧 */
.navbar-left {
  display: flex;
  align-items: center;
  gap: 12rpx;
  white-space: nowrap;

  .location-icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    animation: pulse 2s ease-in-out infinite;
  }
}

@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.1);
  }
}

.hotel-name {
  font-size: 34rpx;
  font-weight: 700;
  color: #333;
  white-space: nowrap;
  letter-spacing: 0.5rpx;
}
/* 骨架屏 */
.skeleton-banner {
  width: calc(100% - 48rpx);
  margin: 24rpx;
  height: 360rpx;
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
  border-radius: 24rpx;
}

.skeleton-address {
  padding: 32rpx 40rpx;
  background-color: #fff;
  margin: 0 24rpx;
  border-radius: 20rpx;
  margin-bottom: 24rpx;

  .skeleton-line {
    height: 32rpx;
    width: 70%;
    background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
    background-size: 200% 100%;
    animation: skeleton-loading 1.5s infinite;
    border-radius: 8rpx;
  }
}

.skeleton-booking {
  display: flex;
  justify-content: space-between;
  padding: 40rpx;
  background-color: #fff;
  margin: 0 24rpx;
  border-radius: 32rpx;

  .skeleton-block {
    width: 180rpx;
    height: 80rpx;
    background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
    background-size: 200% 100%;
    animation: skeleton-loading 1.5s infinite;
    border-radius: 12rpx;
  }
}

@keyframes skeleton-loading {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

/* 主内容 */
.main-content {
  width: 100%;
  height: 100vh;
  overflow-y: scroll;
}

/* 酒店横幅 */
.hotel-banner {
  width: calc(100% - 48rpx);
  margin: 24rpx;
  height: 360rpx;
  background-color: #fff;
  position: relative;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 12rpx 40rpx rgba(194, 157, 113, 0.15);
  transition: all 0.3s ease;

  &:active {
    transform: scale(0.98);
  }

  .banner-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(180deg, transparent 60%, rgba(0, 0, 0, 0.3) 100%);
    z-index: 1;
    pointer-events: none;
  }

  // 确保 swiper 组件也有圆角
  :deep(.tn-swiper) {
    border-radius: 24rpx;
    overflow: hidden;
    height: 100% !important;
  }

  .swiper-item {
    width: 100%;
    height: 100%;
    position: relative;

    .banner-img {
      width: 100%;
      height: 100%;
      transition: transform 0.3s ease;
    }

    .image-gradient {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 100rpx;
      background: linear-gradient(to top, rgba(0, 0, 0, 0.3), transparent);
      z-index: 1;
    }
  }

  .image-count-badge {
    position: absolute;
    right: 24rpx;
    bottom: 24rpx;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(10px);
    padding: 8rpx 20rpx;
    border-radius: 30rpx;
    display: flex;
    align-items: center;
    gap: 8rpx;
    z-index: 2;

    text {
      color: #fff;
      font-size: 24rpx;
      font-weight: 500;
    }
  }
}

/* 地址部分 */
.address-section {
  display: flex;
  align-items: center;
  padding: 32rpx 40rpx;
  margin: 0 24rpx 24rpx;
  background: linear-gradient(135deg, #FFFFFF 0%, #FAFAFA 100%);
  border-radius: 20rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.03);
  gap: 12rpx;
  transition: all 0.3s ease;

  .address-icon-wrapper {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 56rpx;
    height: 56rpx;
    background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.05) 100%);
    border-radius: 50%;
  }

  .address-text {
    flex: 1;
    font-size: 28rpx;
    color: #333;
    font-weight: 500;
    line-height: 1.5;
  }

  .address-arrow {
    display: flex;
    align-items: center;
    opacity: 0.6;
    transition: all 0.3s ease;
  }

  &:active {
    background: linear-gradient(135deg, #FAFAFA 0%, #F5F5F5 100%);
    transform: scale(0.98);

    .address-arrow {
      opacity: 1;
      transform: translateX(4rpx);
    }
  }
}

/* 预订卡片 */
.booking-card {
  margin: 0 24rpx 24rpx;
  padding: 40rpx;
  background: linear-gradient(135deg, #FFFFFF 0%, #FEFEFE 100%);
  border-radius: 32rpx;
  box-shadow: 0 16rpx 48rpx rgba(194, 157, 113, 0.12);
  border: 1rpx solid rgba(194, 157, 113, 0.08);
  transition: all 0.3s ease;

  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 32rpx;

    .card-title {
      font-size: 32rpx;
      font-weight: 700;
      color: #333;
      letter-spacing: 0.5rpx;
    }

    .edit-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 48rpx;
      height: 48rpx;
      background: linear-gradient(135deg, rgba(194, 157, 113, 0.1) 0%, rgba(194, 157, 113, 0.05) 100%);
      border-radius: 50%;
      transition: all 0.3s ease;

      &:active {
        transform: scale(0.9);
      }
    }
  }

  .booking-dates {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 40rpx;
    padding: 32rpx;
    background: linear-gradient(135deg, rgba(194, 157, 113, 0.05) 0%, rgba(194, 157, 113, 0.02) 100%);
    border-radius: 24rpx;

    .date-column {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8rpx;

      .label {
        font-size: 24rpx;
        color: #999;
        font-weight: 500;
      }

      .date-row {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 4rpx;

        .date {
          font-size: 48rpx;
          font-weight: 700;
          color: #333;
          font-family: -apple-system, BlinkMacSystemFont, sans-serif;
          line-height: 1;
        }

        .weekday {
          font-size: 22rpx;
          color: #999;
        }
      }
    }

    .nights-badge {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8rpx;

      .arrow-icon {
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .count {
        font-size: 22rpx;
        color: #C29D71;
        font-weight: 600;
      }
    }
  }

  .search-wrapper {
    margin-bottom: 40rpx;
    padding: 28rpx 32rpx;
    background: linear-gradient(135deg, #F8F8F8 0%, #FAFAFA 100%);
    border-radius: 100rpx;
    display: flex;
    align-items: center;
    gap: 16rpx;
    transition: all 0.3s ease;

    .search-placeholder {
      flex: 1;
      font-size: 28rpx;
      color: #999;
    }

    .search-arrow {
      display: flex;
      align-items: center;
      opacity: 0.6;
    }

    &:active {
      background: linear-gradient(135deg, #F5F5F5 0%, #F8F8F8 100%);
      transform: scale(0.98);
    }
  }

  .action-area {
    margin-bottom: 24rpx;
    box-shadow: 0 16rpx 32rpx rgba(194, 157, 113, 0.35);
    border-radius: 100rpx;
    transition: all 0.3s ease;

    &:active {
      transform: scale(0.98);
      box-shadow: 0 12rpx 24rpx rgba(194, 157, 113, 0.25);
    }

    .btn-content {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 16rpx;
    }

    .btn-text {
      font-size: 36rpx;
      font-weight: 700;
      letter-spacing: 2rpx;
    }
  }

  .guarantee-bar {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8rpx;

    .guarantee-icon {
      display: flex;
      align-items: center;
    }

    text {
      font-size: 22rpx;
      color: #999;
    }
  }
}

/* 功能网格 */
.feature-grid {
  display: flex;
  flex-direction: column;
  padding: 0 24rpx;
  gap: 16rpx;
  margin-bottom: 24rpx;

  .feature-card {
    position: relative;
    background: linear-gradient(135deg, #FFFFFF 0%, #FEFEFE 100%);
    padding: 32rpx;
    border-radius: 24rpx;
    display: flex;
    align-items: center;
    box-shadow: 0 8rpx 24rpx rgba(194, 157, 113, 0.08);
    border: 1rpx solid rgba(194, 157, 113, 0.06);
    overflow: hidden;
    transition: all 0.3s ease;

    .feature-bg-icon {
      position: absolute;
      right: -10rpx;
      top: 50%;
      transform: translateY(-50%);
      opacity: 0.5;
      z-index: 0;
    }

    .feature-content {
      flex: 1;
      display: flex;
      align-items: center;
      gap: 20rpx;
      z-index: 1;

      .feature-icon {
        width: 72rpx;
        height: 72rpx;
        display: flex;
        align-items: center;
        justify-content: center;
        border-radius: 50%;

        &.cinema-icon {
          background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
        }

        &.points-icon {
          background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
        }

        &.vip-icon {
          background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
        }
      }

      .title {
        font-size: 32rpx;
        font-weight: 700;
        color: #333;
        margin-bottom: 4rpx;
      }

      .desc {
        font-size: 22rpx;
        color: #999;
      }
    }

    .feature-arrow {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 44rpx;
      height: 44rpx;
      background: rgba(0, 0, 0, 0.02);
      border-radius: 50%;
      z-index: 1;
      transition: all 0.3s ease;
    }

    &:active {
      transform: scale(0.98);
      background: linear-gradient(135deg, #FAFAFA 0%, #F8F8F8 100%);

      .feature-arrow {
        background: rgba(0, 0, 0, 0.04);
        transform: translateX(4rpx);
      }
    }
  }
}

/* 通告栏 */
.notice-section {
  margin: 0 24rpx 24rpx;

  .notice-wrapper {
    border-radius: 20rpx;
    overflow: hidden;
    box-shadow: 0 4rpx 16rpx rgba(230, 126, 34, 0.08);

    :deep(.tn-notice-bar) {
      border-radius: 20rpx;
    }
  }
}

/* 促销横幅 */
.promotion-banner {
  margin: 0 24rpx 32rpx;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 16rpx 48rpx rgba(194, 157, 113, 0.2);
  position: relative;
  transition: all 0.3s ease;

  &:active {
    transform: scale(0.98);
    box-shadow: 0 12rpx 36rpx rgba(194, 157, 113, 0.15);
  }

  .promotion-img {
    width: 100%;
    height: 320rpx;
  }

  .promotion-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: linear-gradient(180deg, transparent 0%, rgba(0, 0, 0, 0.5) 100%);
    z-index: 1;
  }

  .promotion-content {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    z-index: 2;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    padding: 28rpx;

    .promotion-tag {
      align-self: flex-start;
      background: linear-gradient(135deg, #FF6B6B 0%, #EE5A6F 100%);
      color: #fff;
      padding: 10rpx 24rpx;
      border-radius: 30rpx;
      box-shadow: 0 8rpx 16rpx rgba(255, 107, 107, 0.3);
      backdrop-filter: blur(10px);

      .tag-text {
        font-size: 22rpx;
        font-weight: 600;
        letter-spacing: 1rpx;
      }
    }

    .promotion-info {
      display: flex;
      flex-direction: column;
      gap: 8rpx;

      .promotion-title {
        font-size: 40rpx;
        font-weight: 800;
        color: #fff;
        text-shadow: 0 4rpx 8rpx rgba(0, 0, 0, 0.3);
        letter-spacing: 1rpx;
      }

      .promotion-subtitle {
        font-size: 26rpx;
        color: rgba(255, 255, 255, 0.9);
        text-shadow: 0 2rpx 4rpx rgba(0, 0, 0, 0.2);
      }
    }
  }
}

/* 底部占位 */
.bottom-placeholder {
  height: 80rpx;
}

/* 全局动画 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30rpx);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slideInRight {
  from {
    opacity: 0;
    transform: translateX(30rpx);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

/* 为各个模块添加入场动画 */
.hotel-banner {
  animation: fadeInUp 0.6s ease-out;
}

.address-section {
  animation: fadeInUp 0.6s ease-out 0.1s backwards;
}

.booking-card {
  animation: fadeInUp 0.6s ease-out 0.2s backwards;
}

.feature-grid {
  .feature-card {
    &:nth-child(1) {
      animation: fadeInUp 0.6s ease-out 0.3s backwards;
    }
    &:nth-child(2) {
      animation: fadeInUp 0.6s ease-out 0.4s backwards;
    }
    &:nth-child(3) {
      animation: fadeInUp 0.6s ease-out 0.5s backwards;
    }
  }
}

.notice-section {
  animation: fadeInUp 0.6s ease-out 0.6s backwards;
}

.promotion-banner {
  animation: fadeInUp 0.6s ease-out 0.7s backwards;
}

/* 悬浮效果增强 */
@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-8rpx);
  }
}

/* 为按钮添加渐变闪烁效果 */
.action-area {
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: linear-gradient(
      45deg,
      transparent 30%,
      rgba(255, 255, 255, 0.3) 50%,
      transparent 70%
    );
    animation: shimmer 3s infinite;
  }
}

@keyframes shimmer {
  0% {
    transform: translateX(-100%) translateY(-100%) rotate(45deg);
  }
  100% {
    transform: translateX(100%) translateY(100%) rotate(45deg);
  }
}

/* 日期弹窗样式重构 */
.calendar-container {
  background: linear-gradient(180deg, #FFFFFF 0%, #FAFAFA 100%);
  height: 80vh;

  .calendar-header {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32rpx;
    position: relative;
    border-bottom: 1rpx solid rgba(194, 157, 113, 0.1);
    height: 100rpx;
    box-sizing: border-box;
    background: #fff;

    .title {
      font-size: 36rpx;
      font-weight: 700;
      color: #333;
      letter-spacing: 1rpx;
    }

    .close-btn {
      position: absolute;
      right: 32rpx;
      top: 50%;
      transform: translateY(-50%);
      padding: 10rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      width: 56rpx;
      height: 56rpx;
      background: rgba(0, 0, 0, 0.03);
      border-radius: 50%;
      transition: all 0.3s ease;

      &:active {
        background: rgba(0, 0, 0, 0.06);
        transform: translateY(-50%) scale(0.9);
      }
    }
  }

  .calendar-body {
    height: calc(80vh - 100rpx - 140rpx);
    overflow-y: auto;
    background: #fff;

    .calendar-loading {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 400rpx;
      color: #999;
      font-size: 28rpx;
    }

    // 覆盖日历组件样式，确保在小程序中正确显示
    :deep(.tn-calendar) {
      width: 100%;
    }

    :deep(.tn-calendar__data) {
      height: 650rpx !important;
      min-height: 650rpx !important;
    }

    // swiper 在小程序中必须有明确的固定高度
    :deep(.tn-calendar__data__swiper) {
      height: 650rpx !important;
      min-height: 650rpx !important;
    }

    :deep(.tn-calendar__data__swiper-item) {
      height: 650rpx !important;
      overflow: visible !important;
    }

    :deep(.tn-calendar__data__dates) {
      display: flex !important;
      flex-wrap: wrap !important;
      min-height: 600rpx;
    }

    :deep(.tn-calendar__data__date) {
      display: flex !important;
      align-items: center !important;
      justify-content: center !important;
    }
  }

  .calendar-footer {
    padding: 28rpx 32rpx;
    padding-bottom: calc(28rpx + env(safe-area-inset-bottom));
    background: linear-gradient(180deg, #FFFFFF 0%, #FAFAFA 100%);
    border-top: 1rpx solid rgba(194, 157, 113, 0.1);
    height: 140rpx;
    box-sizing: border-box;

    :deep(.tn-button) {
      box-shadow: 0 12rpx 24rpx rgba(194, 157, 113, 0.3);

      &:active {
        transform: scale(0.98);
      }
    }
  }
}
</style>




