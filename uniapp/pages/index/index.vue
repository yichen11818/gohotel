<template>
  <view class="container">
    <!-- 自定义导航栏 -->
    <TnNavbar fixed :bottom-shadow="false" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-left">
          <TnIcon name="location-fill" color="#E74C3C" size="36" />
          <text class="hotel-name">{{ hotelData?.name || '七天酒店' }}</text>
        </view>
      </template>
      <template #right>
        <view class="navbar-right">
          <view class="icon-btn" @click="showMoreOptions">
            <TnIcon name="more-vertical" color="#333" size="38" />
          </view>
          <view class="icon-btn" @click="scanQRCode">
            <TnIcon name="scan" color="#333" size="38" />
          </view>
        </view>
      </template>
    </TnNavbar>

    <!-- 主内容区域 -->
    <scroll-view class="main-content" scroll-y>
      <!-- 加载骨架屏 -->
      <template v-if="loading">
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
          <TnSwiper
            v-model="currentSwiperIndex"
            :data="hotelImages"
            width="100%"
            height="480"
            autoplay
            loop
            indicator
            indicator-type="dot"
            indicator-bg-color="rgba(255,255,255,0.4)"
            indicator-active-bg-color="#fff"
          >
            <template #default="{ data }">
              <view class="swiper-item">
                <image class="banner-img" :src="data" mode="aspectFill"></image>
              </view>
            </template>
          </TnSwiper>
        </view>

        <!-- 酒店地址 -->
        <view class="address-section" @click="openMap">
          <TnIcon name="location-fill" color="#333" size="36" class="address-icon" />
          <text class="address-text">{{ address }}</text>
        </view>

        <!-- 预订信息 -->
        <view class="booking-card" @click="openDatePicker">
          <view class="booking-dates">
            <view class="date-column">
              <text class="label">{{ checkInLabel }}</text>
              <view class="date-row">
                <text class="date">{{ formatDisplayDate(checkInDate) }}</text>
                <text class="week">{{ getWeekDay(checkInDate) }}</text>
              </view>
            </view>

            <view class="nights-badge-outline">
              <text class="count">共{{ nights }}晚</text>
              <TnIcon name="right" size="20" color="#999" />
            </view>

            <view class="date-column right">
              <text class="label">{{ checkOutLabel }}</text>
              <view class="date-row">
                <text class="week">{{ getWeekDay(checkOutDate) }}</text>
                <text class="date">{{ formatDisplayDate(checkOutDate) }}</text>
              </view>
            </view>
          </view>

          <!-- 搜索框 -->
          <view class="search-wrapper" @click.stop="goToSearch">
             <TnIcon name="search" color="#ccc" size="32" class="search-icon" />
             <text class="search-placeholder">输入关键词搜索酒店</text>
          </view>

          <!-- 立即预订按钮 -->
          <view class="action-area">
             <TnButton 
               shape="round" 
               size="xl" 
               width="100%" 
               height="100rpx"
               bg-color="#C29D71"
               text-color="#FFFFFF"
               @click.stop="handleBooking"
             >
               <text class="btn-text">立即预定</text>
             </TnButton>
          </view>
          
          <view class="guarantee-bar">
            <text>官方渠道预订享低价保证</text>
          </view>
        </view>

        <!-- 功能入口 -->
        <view class="feature-grid">
          <view class="feature-card" @click="navigateTo('cinema')">
            <view class="feature-text">
              <text class="title">影院足道</text>
              <text class="desc">官方自营 ></text>
            </view>
            <view class="feature-icon cinema-icon">
               <TnIcon name="video-fill" color="#E74C3C" size="44" />
            </view>
          </view>
          <view class="feature-card" @click="navigateTo('points')">
            <view class="feature-text">
              <text class="title">积分商城</text>
              <text class="desc">体验兑换 ></text>
            </view>
            <view class="feature-icon points-icon">
               <TnIcon name="gift-fill" color="#3498DB" size="44" />
            </view>
          </view>
          <view class="feature-card" @click="navigateTo('vip')">
            <view class="feature-text">
              <text class="title">升级会员</text>
              <text class="desc">更多折扣 ></text>
            </view>
            <view class="feature-icon vip-icon">
               <TnIcon name="vip-fill" color="#F39C12" size="44" />
            </view>
          </view>
        </view>

        <!-- 促销通告 -->
        <view class="notice-section" v-if="true">
           <TnNoticeBar 
             :data="['双11特惠活动火热进行中！','新用户注册立享88元大礼包']" 
             direction="vertical" 
             left-icon="sound"
             bg-color="#FFF4E5"
             color="#E67E22"
           />
        </view>

        <!-- 促销横幅 -->
        <view class="promotion-banner" @click="goToPromotion">
          <image class="promotion-img" :src="promotionImage" mode="aspectFill"></image>
          <view class="promotion-tag">限时特惠</view>
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
    >
      <view class="calendar-container">
        <!-- 头部 -->
        <view class="calendar-header">
          <text class="title">选择日期</text>
          <view class="close-btn" @click="showCalendar = false">
            <TnIcon name="close" size="40" color="#999" />
          </view>
        </view>
        
        <!-- 移除了自定义的状态栏和星期栏，直接使用组件自带的以避免布局冲突 -->

        <!-- 日期选择区域 - 移除scroll-view，防止高度塌陷 -->
        <view class="calendar-body">
          <TnCalendar 
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
import { hotel, booking } from '@/api/index.js'
import TnSwiper from '@/uni_modules/tuniaoui-vue3/components/swiper/src/swiper.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import TnPopup from '@/uni_modules/tuniaoui-vue3/components/popup/src/popup.vue'
import TnCalendar from '@/uni_modules/tuniaoui-vue3/components/calendar/src/calendar.vue'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'
import TnSearchBox from '@/uni_modules/tuniaoui-vue3/components/search-box/src/search-box.vue'
import TnTag from '@/uni_modules/tuniaoui-vue3/components/tag/src/tag.vue'
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

// 日期相关
const checkInDate = ref(null)  // Date对象
const checkOutDate = ref(null) // Date对象
const showCalendar = ref(false)
const selectedDateRange = ref([])
const tempCheckIn = ref(null)
const tempCheckOut = ref(null)

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

// 解析日期字符串
const parseDate = (dateStr) => {
  if (!dateStr) return null
  const [year, month, day] = dateStr.split('-').map(Number)
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
  
  showCalendar.value = true
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
    showCalendar.value = false
    
    uni.showToast({
      title: `已选择${nights.value}晚`,
      icon: 'none'
    })
  }
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
})

// 加载酒店详情
const loadHotelData = async () => {
  try {
    loading.value = true
    const data = await hotel.getHotelDetail(hotelId.value)
    hotelData.value = data
    
    // 更新页面数据
    if (data.images && data.images.length > 0) {
      hotelImages.value = data.images
    }
    if (data.address) {
      address.value = data.address
    }
  } catch (error) {
    console.error('加载酒店数据失败:', error)
    // 使用默认数据
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
  background-color: #f5f5f5;
}

/* 骨架屏 */
.skeleton-banner {
  width: 100%;
  height: 480rpx;
  background: linear-gradient(90deg, #f0f0f0 25%, #e8e8e8 50%, #f0f0f0 75%);
  background-size: 200% 100%;
  animation: skeleton-loading 1.5s infinite;
}

.skeleton-address {
  padding: 32rpx 40rpx;
  background-color: #fff;
  
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
  margin-top: 20rpx;
  
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
  width: 100%;
  height: 480rpx;
  background-color: #fff;
  position: relative;
  
  .swiper-item {
    width: 100%;
    height: 100%;
    position: relative;
    
    .banner-img {
      width: 100%;
      height: 100%;
    }
    
    .image-counter {
      position: absolute;
      right: 24rpx;
      bottom: 24rpx;
      background-color: rgba(0, 0, 0, 0.5);
      padding: 8rpx 20rpx;
      border-radius: 20rpx;
      
      text {
        color: #fff;
        font-size: 24rpx;
      }
    }
  }
}

/* 地址部分 */
.address-section {
  display: flex;
  align-items: center;
  padding: 32rpx 40rpx;
  background-color: #fff;
  gap: 12rpx;
  
  .address-text {
    flex: 1;
    font-size: 34rpx;
    color: #000;
    font-weight: 600;
    line-height: 1.4;
  }
  
  &:active {
    background-color: #f9f9f9;
  }
}

/* 预订卡片 */
.booking-card {
  margin: 24rpx;
  padding: 48rpx 40rpx;
  background-color: #fff;
  border-radius: 32rpx;
  box-shadow: 0 12rpx 40rpx rgba(0, 0, 0, 0.05);
  
  .booking-dates {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 48rpx;
    
    .date-column {
      display: flex;
      flex-direction: column;
      
      &.right {
        align-items: flex-end;
        
        .date-row {
          justify-content: flex-end;
        }
      }
      
      .label {
        font-size: 24rpx;
        color: #999;
        margin-bottom: 12rpx;
      }
      
      .date-row {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 8rpx;
        
        .date {
          font-size: 52rpx;
          font-weight: 600;
          color: #333;
          font-family: Helvetica, Arial, sans-serif;
          line-height: 1;
          letter-spacing: -1px;
        }
        
        .week {
          font-size: 24rpx;
          color: #666;
        }
      }
      
      &.right .date-row {
         align-items: flex-end;
      }
    }
    
    .nights-badge-outline {
      padding: 6rpx 20rpx;
      border: 1px solid #eee;
      border-radius: 100rpx;
      display: flex;
      align-items: center;
      gap: 4rpx;
      
      .count {
        font-size: 24rpx;
        color: #C29D71;
        font-weight: 500;
      }
    }
  }
  
  .search-wrapper {
    margin-bottom: 48rpx;
    padding: 24rpx 32rpx;
    background-color: #f9f9f9;
    border-radius: 100rpx;
    display: flex;
    align-items: center;
    gap: 16rpx;
    
    .search-placeholder {
      font-size: 28rpx;
      color: #bbb;
    }
  }
  
  .action-area {
    margin-bottom: 24rpx;
    box-shadow: 0 20rpx 40rpx -10rpx rgba(194, 157, 113, 0.4);
    border-radius: 100rpx;
    
    .btn-text {
      font-size: 34rpx;
      font-weight: 600;
      letter-spacing: 2rpx;
    }
  }
  
  .guarantee-bar {
    display: flex;
    align-items: center;
    justify-content: center;
    
    text {
      font-size: 24rpx;
      color: #999;
    }
  }
}

/* 功能网格 */
.feature-grid {
  display: flex;
  justify-content: space-between;
  padding: 0 24rpx;
  gap: 20rpx;
  margin-bottom: 24rpx;
  
  .feature-card {
    flex: 1;
    background-color: #fff;
    padding: 32rpx 24rpx;
    border-radius: 20rpx;
    display: flex;
    justify-content: space-between;
    align-items: center;
    box-shadow: 0 4rpx 16rpx rgba(0,0,0,0.02);
    
    .feature-text {
      display: flex;
      flex-direction: column;
      gap: 4rpx;
      
      .title {
        font-size: 28rpx;
        font-weight: 600;
        color: #333;
      }
      
      .desc {
        font-size: 20rpx;
        color: #ccc;
      }
    }
    
    .feature-icon {
      width: 80rpx;
      height: 80rpx;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      
      &.cinema-icon { background-color: #FFF5F5; }
      &.points-icon { background-color: #F0F9FF; }
      &.vip-icon { background-color: #FFF8E1; }
    }
    
    &:active {
      transform: scale(0.98);
      background-color: #fafafa;
    }
  }
}

/* 通告栏 */
.notice-section {
  margin: 0 24rpx 24rpx;
  border-radius: 16rpx;
  overflow: hidden;
}


/* 促销横幅 */
.promotion-banner {
  margin: 32rpx 24rpx;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.1);
  position: relative;
  
  &:active {
    transform: scale(0.98);
  }
  
  .promotion-img {
    width: 100%;
    height: 320rpx;
  }
  
  .promotion-tag {
    position: absolute;
    top: 20rpx;
    left: 20rpx;
    background: linear-gradient(135deg, #E74C3C 0%, #C0392B 100%);
    color: #fff;
    font-size: 22rpx;
    padding: 8rpx 20rpx;
    border-radius: 8rpx;
    font-weight: 500;
  }
}

/* 底部占位 */
.bottom-placeholder {
  height: 60rpx;
}

/* 日期弹窗样式重构 */
.calendar-container {
  background-color: #fff;
  height: 80vh;
  
  .calendar-header {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 32rpx;
    position: relative;
    border-bottom: 1px solid #f5f5f5;
    height: 100rpx;
    box-sizing: border-box;
    
    .title {
      font-size: 34rpx;
      font-weight: 600;
      color: #333;
    }
    
    .close-btn {
      position: absolute;
      right: 32rpx;
      top: 50%;
      transform: translateY(-50%);
      padding: 10rpx;
    }
  }
  
  .calendar-body {
    height: calc(80vh - 100rpx - 140rpx);
    overflow-y: auto;
  }
  
  .calendar-footer {
    padding: 24rpx 32rpx;
    padding-bottom: calc(24rpx + env(safe-area-inset-bottom));
    background-color: #fff;
    border-top: 1px solid #f5f5f5;
    height: 140rpx;
    box-sizing: border-box;
  }
}
</style>


