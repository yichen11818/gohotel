<template>
  <scroll-view class="page" scroll-y>
    <view class="header">
      <view>
        <text class="title">{{ hotelInfo.name }}</text>
        <text class="subtitle">{{ hotelInfo.address }}</text>
      </view>
      <button class="mini-btn" size="mini" @click="goBack">返回</button>
    </view>

    <view class="content">
      <view class="summary-card">
        <view class="summary-row">
          <text>入住</text>
          <text>{{ checkInDate }}</text>
        </view>
        <view class="summary-row">
          <text>离店</text>
          <text>{{ checkOutDate }}</text>
        </view>
        <view class="summary-row">
          <text>酒店电话</text>
          <text>{{ hotelInfo.phone || '暂未配置' }}</text>
        </view>
        <view class="summary-row">
          <text>入住 / 退房</text>
          <text>{{ hotelInfo.checkInTime }} / {{ hotelInfo.checkOutTime }}</text>
        </view>
      </view>

      <view class="section-title">可订房型</view>
      <view v-if="loading" class="placeholder">正在加载房型...</view>
      <view v-else-if="!roomList.length" class="placeholder">当前暂无可订房型</view>

      <view v-for="room in roomList" :key="room.id" class="room-card" @click="handleRoomClick(room)">
        <image class="room-image" :src="room.image" mode="aspectFill" />
        <view class="room-body">
          <view class="room-top">
            <text class="room-name">{{ room.name }}</text>
            <text class="room-price">¥{{ room.price }}/晚</text>
          </view>
          <text class="room-desc">{{ room.area }}㎡ · {{ room.bedType }} · 可住 {{ room.capacity }} 人</text>
          <text class="room-text">{{ room.description }}</text>
          <view class="tag-list">
            <text v-for="item in room.facilities.slice(0, 4)" :key="item" class="tag">{{ item }}</text>
          </view>
          <button class="primary-btn" size="mini" @click.stop="handleRoomClick(room)">查看详情</button>
        </view>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { hotel } from '@/api/index.js'

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
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${room.id}&checkIn=${checkInDate.value}&checkOut=${checkOutDate.value}`,
  })
}

const goBack = () => {
  uni.navigateBack()
}

onLoad((options) => {
  if (options?.id) {
    hotelId.value = Number(options.id)
  }
  checkInDate.value = normalizeDateParam(options?.checkIn)
  checkOutDate.value = normalizeDateParam(options?.checkOut)
  ensureDates()
  loadData()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 24rpx;
  padding: 88rpx 24rpx 24rpx;
  background: linear-gradient(180deg, #e9d7bf 0%, #f6f7fb 100%);
}

.title {
  display: block;
  font-size: 38rpx;
  font-weight: 700;
  color: #1f2937;
}

.subtitle {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #4b5563;
  line-height: 1.6;
}

.mini-btn {
  margin: 0;
  background: #fff;
  border-radius: 999rpx;
}

.content {
  padding: 0 24rpx 32rpx;
}

.summary-card,
.room-card {
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.summary-card {
  padding: 28rpx;
}

.summary-row {
  display: flex;
  justify-content: space-between;
  gap: 24rpx;
  padding: 18rpx 0;
  font-size: 26rpx;
  color: #374151;
  border-bottom: 1px solid #f3f4f6;
}

.summary-row:last-child {
  border-bottom: 0;
}

.section-title {
  margin: 24rpx 0 16rpx;
  font-size: 30rpx;
  font-weight: 700;
  color: #1f2937;
}

.placeholder {
  padding: 48rpx 24rpx;
  text-align: center;
  color: #6b7280;
  font-size: 26rpx;
}

.room-image {
  width: 100%;
  height: 320rpx;
  border-radius: 24rpx 24rpx 0 0;
}

.room-body {
  padding: 24rpx;
}

.room-top {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  align-items: center;
}

.room-name {
  flex: 1;
  font-size: 30rpx;
  font-weight: 700;
  color: #111827;
}

.room-price {
  font-size: 30rpx;
  color: #b7791f;
  font-weight: 700;
}

.room-desc,
.room-text {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #4b5563;
  line-height: 1.6;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-top: 18rpx;
}

.tag {
  padding: 8rpx 18rpx;
  font-size: 22rpx;
  color: #8b5e34;
  background: #faf4ed;
  border-radius: 999rpx;
}

.primary-btn {
  margin-top: 22rpx;
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
  border-radius: 999rpx;
}
</style>
