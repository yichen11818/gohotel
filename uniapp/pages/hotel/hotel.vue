<template>
  <view class="page">
    <!-- 头部：酒店基础信息 -->
    <view class="hotel-header">
      <view class="header-content" :style="{ paddingTop: statusBarHeight + 'px' }">
        <view class="nav-bar">
          <view class="back-btn flex-center" @click="goBack">
            <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" />
          </view>
        </view>
        <view class="hotel-main">
          <text class="hotel-name">{{ hotelInfo.name }}</text>
          <view class="hotel-sub-info">
            <text class="hotel-address">{{ hotelInfo.address }}</text>
            <view class="info-tags">
              <text class="info-tag">4.8分</text>
              <text class="info-tag">高档型</text>
              <text class="info-tag">2023年装修</text>
            </view>
          </view>
        </view>
      </view>
      <view class="header-bg"></view>
    </view>

    <view class="content">
      <!-- 日期汇总条 -->
      <view class="date-summary premium-card flex-between" @click="handleDateChange">
        <view class="date-group flex-center">
          <view class="date-item">
            <text class="label">入住</text>
            <text class="val">{{ checkInDisplay }}</text>
          </view>
          <view class="duration">
            <text class="duration-val">{{ nights }}晚</text>
          </view>
          <view class="date-item">
            <text class="label">离店</text>
            <text class="val">{{ checkOutDisplay }}</text>
          </view>
        </view>
        <text class="change-btn">修改日期</text>
      </view>

      <!-- 房型列表 -->
      <view class="room-section">
        <view class="section-title">可订房型</view>

        <view v-if="loading" class="loading-state flex-center">
          <view class="loading-spinner"></view>
          <text>加载中...</text>
        </view>

        <view v-else-if="!roomList.length" class="empty-state">
          <image src="/static/images/empty.png" mode="aspectFit" class="empty-img" />
          <text>暂无可预订房型</text>
        </view>

        <view v-else class="room-list">
          <view v-for="room in roomList" :key="room.id" class="room-item premium-card" @click="handleRoomClick(room)">
            <view class="room-content">
              <image class="room-image" :src="room.image" mode="aspectFill" />
              <view class="room-info">
                <view class="info-header">
                  <text class="room-name">{{ room.name }}</text>
                  <view class="room-tags">
                    <text class="tag">{{ room.area }}㎡</text>
                    <text class="tag">{{ room.bedType }}</text>
                  </view>
                </view>
                <view class="info-footer flex-between">
                  <view class="price-box">
                    <text class="symbol">¥</text>
                    <text class="price">{{ room.price }}</text>
                    <text class="unit">起</text>
                  </view>
                  <button class="book-btn premium-button" @click.stop="handleRoomClick(room)">查看</button>
                </view>
              </view>
            </view>
          </view>
        </view>
      </view>
    </view>

    <HotelCalendar v-model:show="showCalendar" :check-in="checkInDate" :check-out="checkOutDate" @confirm="onCalendarConfirm" />
  </view>
</template>

<script setup>
import { ref, computed } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { hotel } from '@/api/index.js'
import HotelCalendar from '@/components/hotel-calendar/hotel-calendar.vue'

const statusBarHeight = ref(44)
const showCalendar = ref(false)
const hotelInfo = ref({
  name: '七天酒店',
  address: '欢迎入住七天酒店',
  phone: '',
  checkInTime: '14:00',
  checkOutTime: '12:00',
})
const roomList = ref([])
const loading = ref(false)
const hotelId = ref(hotel.DEFAULT_HOTEL_ID || 1)
const checkInDate = ref('')
const checkOutDate = ref('')
const HOTEL_TAB_CONTEXT_KEY = 'gohotel_hotel_tab_context'

const checkInDisplay = computed(() => {
  if (!checkInDate.value) return '--'
  const d = new Date(checkInDate.value)
  return `${d.getMonth() + 1}月${d.getDate()}日`
})

const checkOutDisplay = computed(() => {
  if (!checkOutDate.value) return '--'
  const d = new Date(checkOutDate.value)
  return `${d.getMonth() + 1}月${d.getDate()}日`
})

const nights = computed(() => {
  if (!checkInDate.value || !checkOutDate.value) return 1
  const start = new Date(checkInDate.value)
  const end = new Date(checkOutDate.value)
  const diff = end.getTime() - start.getTime()
  return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)))
})

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

const applyPageContext = (context = {}) => {
  const nextHotelId = Number(context.hotelId || context.id)
  if (Number.isFinite(nextHotelId) && nextHotelId > 0) hotelId.value = nextHotelId
  const nextCheckIn = normalizeDateParam(context.checkIn)
  const nextCheckOut = normalizeDateParam(context.checkOut)
  if (nextCheckIn) checkInDate.value = nextCheckIn
  if (nextCheckOut) checkOutDate.value = nextCheckOut
  ensureDates()
}

const syncStoredContext = () => {
  try {
    applyPageContext(uni.getStorageSync(HOTEL_TAB_CONTEXT_KEY) || {})
  } catch (_error) {
    ensureDates()
  }
}

const loadData = async () => {
  loading.value = true
  try {
    const [detail, rooms] = await Promise.all([
      hotel.getHotelDetail(hotelId.value),
      hotel.getRoomTypes(hotelId.value),
    ])
    hotelInfo.value = detail
    roomList.value = rooms
  } catch (error) {
    console.error('load hotel data failed:', error)
  } finally {
    loading.value = false
  }
}

const handleRoomClick = (room) => {
  hotel.trackRoomBehavior(room.id, 'view_room', 'hotel_list', {
    hotel_id: hotelId.value,
    check_in: checkInDate.value,
    check_out: checkOutDate.value,
  })
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${room.id}&checkIn=${checkInDate.value}&checkOut=${checkOutDate.value}`,
  })
}

const handleDateChange = () => {
  showCalendar.value = true
}

const onCalendarConfirm = (dates) => {
  checkInDate.value = dates.checkIn
  checkOutDate.value = dates.checkOut
  // 如果需要，可以在这里触发重新加载房型列表数据
  // loadData()
}

const goBack = () => {
  const pages = getCurrentPages()
  pages.length > 1 ? uni.navigateBack() : uni.switchTab({ url: '/pages/index/index' })
}

onLoad((options) => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
  applyPageContext(options || {})
  syncStoredContext()
})

onShow(() => {
  syncStoredContext()
  loadData()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: $bg-color;
  padding-bottom: 60rpx;
}

/* 头部设计 */
.hotel-header {
  position: relative;
  height: 460rpx;
  overflow: hidden;

  .header-content {
    position: relative;
    z-index: 10;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 0 32rpx;
  }

  .nav-bar {
    height: 88rpx;
    display: flex;
    align-items: center;

    .back-btn {
      width: 64rpx;
      height: 64rpx;
      background: rgba(0,0,0,0.2);
      border-radius: 50%;

      .icon-back {
        width: 32rpx;
        height: 32rpx;
      }
    }
  }

  .hotel-main {
    margin-top: auto;
    padding-bottom: 60rpx;
    color: #fff;

    .hotel-name {
      font-size: 44rpx;
      font-weight: 700;
      margin-bottom: 16rpx;
      display: block;
    }

    .hotel-sub-info {
      .hotel-address {
        font-size: 24rpx;
        opacity: 0.9;
        display: block;
        margin-bottom: 16rpx;
      }

      .info-tags {
        display: flex;
        gap: 12rpx;

        .info-tag {
          font-size: 20rpx;
          padding: 4rpx 16rpx;
          background: rgba(255,255,255,0.2);
          border-radius: 4rpx;
          backdrop-filter: blur(4px);
        }
      }
    }
  }

  .header-bg {
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
    &::after {
      content: '';
      position: absolute;
      inset: 0;
      background: rgba(0,0,0,0.15);
    }
  }
}

.content {
  padding: 0 30rpx;
  margin-top: -40rpx;
  position: relative;
  z-index: 20;
}

/* 日期汇总 */
.date-summary {
  padding: 24rpx 32rpx;
  margin-bottom: 30rpx;

  .date-group {
    .date-item {
      .label {
        font-size: 20rpx;
        color: $text-sub;
        display: block;
      }
      .val {
        font-size: 28rpx;
        font-weight: 600;
        color: $text-main;
      }
    }

    .duration {
      margin: 0 30rpx;
      padding: 0 16rpx;
      height: 32rpx;
      line-height: 32rpx;
      border: 1rpx solid #eee;
      border-radius: 16rpx;

      .duration-val {
        font-size: 20rpx;
        color: $text-second;
      }
    }
  }

  .change-btn {
    font-size: 24rpx;
    color: $primary-color;
    font-weight: 500;
  }
}

/* 房型列表 */
.room-section {
  .section-title {
    font-size: 32rpx;
    font-weight: 700;
    color: $text-main;
    margin-bottom: 24rpx;
    padding-left: 10rpx;
  }

  .room-item {
    padding: 0;
    margin-bottom: 24rpx;

    .room-content {
      display: flex;
      padding: 20rpx;
      gap: 20rpx;
    }

    .room-image {
      width: 220rpx;
      height: 220rpx;
      border-radius: $radius-md;
      flex-shrink: 0;
    }

    .room-info {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      padding: 4rpx 0;

      .info-header {
        .room-name {
          font-size: 30rpx;
          font-weight: 600;
          color: $text-main;
          display: block;
          margin-bottom: 12rpx;
        }

        .room-tags {
          display: flex;
          gap: 12rpx;

          .tag {
            font-size: 22rpx;
            color: $text-sub;
            background: $bg-color;
            padding: 4rpx 12rpx;
            border-radius: 4rpx;
          }
        }
      }

      .info-footer {
        .price-box {
          display: flex;
          align-items: baseline;

          .symbol {
            font-size: 24rpx;
            color: #E64340;
            font-weight: 700;
          }
          .price {
            font-size: 40rpx;
            color: #E64340;
            font-weight: 700;
            margin: 0 4rpx;
          }
          .unit {
            font-size: 20rpx;
            color: $text-sub;
          }
        }

        .book-btn {
          width: 120rpx;
          height: 60rpx;
          font-size: 26rpx;
          margin: 0;
        }
      }
    }
  }
}

.loading-state, .empty-state {
  padding: 100rpx 0;
  flex-direction: column;
  color: $text-sub;
  font-size: 26rpx;

  .empty-img {
    width: 200rpx;
    height: 200rpx;
    margin-bottom: 20rpx;
  }
}
</style>
