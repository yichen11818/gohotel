<template>
  <view class="container">
    <!-- 图片轮播 -->
    <view class="image-section">
      <TnSwiper
        v-model="currentImageIndex"
        :data="roomImages"
        width="100%"
        height="500"
        autoplay
        loop
        indicator
        indicator-type="dot"
        indicator-bg-color="rgba(255,255,255,0.4)"
        indicator-active-bg-color="#C29D71"
      >
        <template #default="{ data }">
          <view class="swiper-item">
            <image class="room-image" :src="data" mode="aspectFill"></image>
          </view>
        </template>
      </TnSwiper>
      
      <!-- 返回按钮 -->
      <view class="back-btn" @click="goBack">
        <TnIcon name="left" color="#fff" size="40" />
      </view>
      
      <!-- 图片计数 -->
      <view class="image-count">
        <text>{{ currentImageIndex + 1 }}/{{ roomImages.length }}</text>
      </view>
    </view>

    <!-- 房间信息 -->
    <scroll-view class="content-scroll" scroll-y>
      <!-- 基本信息 -->
      <view class="room-info-card">
        <view class="room-header">
          <text class="room-name">{{ roomDetail.name }}</text>
          <view class="room-tags">
            <text v-for="(tag, index) in roomDetail.tags" :key="index" class="tag">{{ tag }}</text>
          </view>
        </view>
        
        <view class="room-price">
          <text class="currency">¥</text>
          <text class="price">{{ roomDetail.price }}</text>
          <text class="unit">/晚</text>
          <text class="original-price" v-if="roomDetail.originalPrice">¥{{ roomDetail.originalPrice }}</text>
        </view>
        
        <view class="room-specs">
          <view class="spec-item">
            <TnIcon name="cube" color="#C29D71" size="32" />
            <text>{{ roomDetail.area }}m²</text>
          </view>
          <view class="spec-item">
            <TnIcon name="bed" color="#C29D71" size="32" />
            <text>{{ roomDetail.bedType }}</text>
          </view>
          <view class="spec-item">
            <TnIcon name="users" color="#C29D71" size="32" />
            <text>{{ roomDetail.maxGuests }}人入住</text>
          </view>
          <view class="spec-item">
            <TnIcon name="window" color="#C29D71" size="32" />
            <text>{{ roomDetail.window }}</text>
          </view>
        </view>
      </view>

      <!-- 入住日期 -->
      <view class="date-card" @click="openDatePicker">
        <view class="date-header">
          <text class="title">入住日期</text>
          <view class="edit-icon">
            <TnIcon name="edit" color="#C29D71" size="28" />
          </view>
        </view>
        <view class="date-content">
          <view class="date-item">
            <text class="label">入住</text>
            <text class="date">{{ formatDate(checkInDate) }}</text>
          </view>
          <view class="nights">
            <text>{{ nights }}晚</text>
          </view>
          <view class="date-item">
            <text class="label">离店</text>
            <text class="date">{{ formatDate(checkOutDate) }}</text>
          </view>
        </view>
      </view>

      <!-- 房间设施 -->
      <view class="facility-card">
        <view class="card-title">房间设施</view>
        <view class="facility-grid">
          <view v-for="(facility, index) in roomDetail.facilities" :key="index" class="facility-item">
            <TnIcon :name="facility.icon" color="#C29D71" size="36" />
            <text>{{ facility.name }}</text>
          </view>
        </view>
      </view>

      <!-- 预订须知 -->
      <view class="notice-card">
        <view class="card-title">预订须知</view>
        <view class="notice-list">
          <view class="notice-item">
            <text class="label">入住时间</text>
            <text class="value">14:00后</text>
          </view>
          <view class="notice-item">
            <text class="label">离店时间</text>
            <text class="value">12:00前</text>
          </view>
          <view class="notice-item">
            <text class="label">取消政策</text>
            <text class="value">入住前1天18:00前免费取消</text>
          </view>
          <view class="notice-item">
            <text class="label">早餐</text>
            <text class="value">{{ roomDetail.breakfast ? '含双早' : '不含早' }}</text>
          </view>
        </view>
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 底部操作栏 -->
    <view class="bottom-bar">
      <view class="price-info">
        <view class="total-price">
          <text class="currency">¥</text>
          <text class="amount">{{ totalPrice }}</text>
        </view>
        <text class="price-desc">共{{ nights }}晚</text>
      </view>
      <view class="book-btn" @click="handleBook">
        <text>立即预订</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { hotel, booking } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'
import TnSwiper from '@/uni_modules/tuniaoui-vue3/components/swiper/src/swiper.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const currentImageIndex = ref(0)
const roomId = ref('')
const loading = ref(false)

// 房间图片
const roomImages = ref([])

// 房间详情
const roomDetail = ref({
  id: '',
  name: '',
  price: 0,
  original_price: 0,
  area: 0,
  bed_type: '',
  max_guests: 0,
  window: '',
  breakfast: false,
  tags: [],
  facilities: []
})

// 日期
const checkInDate = ref(new Date())
const checkOutDate = ref(new Date(Date.now() + 24 * 60 * 60 * 1000))

// 计算入住天数
const nights = computed(() => {
  const diff = checkOutDate.value.getTime() - checkInDate.value.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

// 计算总价
const totalPrice = computed(() => {
  return roomDetail.value.price * nights.value
})

// 格式化日期
const formatDate = (date) => {
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日`
}

// 返回
const goBack = () => {
  uni.navigateBack()
}

// 打开日期选择器
const openDatePicker = () => {
  uni.showToast({ title: '日期选择功能', icon: 'none' })
}

// 立即预订
const handleBook = () => {
  const token = uni.getStorageSync(TOKEN_KEY)
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
  
  // 跳转到预订确认页
  uni.navigateTo({
    url: `/pages/booking/booking?roomId=${roomDetail.value.id}&checkIn=${checkInDate.value.toISOString()}&checkOut=${checkOutDate.value.toISOString()}`
  })
}

onLoad((options) => {
  if (options?.id) {
    roomId.value = options.id
    loadRoomDetail()
  }
  
  // 解析日期参数
  if (options?.checkIn) {
    checkInDate.value = new Date(options.checkIn)
  }
  if (options?.checkOut) {
    checkOutDate.value = new Date(options.checkOut)
  }
})

// 加载房间详情
const loadRoomDetail = async () => {
  loading.value = true
  try {
    // 假设 hotel.js 中有 getRoomDetail 方法，如果没有则需要添加
    // 目前 api/hotel.js 只有 getRoomList, getHotelDetail, getRoomTypes 等
    // 我们可以根据业务逻辑获取房型详情
    const data = await hotel.getRoomTypes(1) // 这里暂时用 hotelId=1，实际应从 options 获取
    const detail = data.find(item => item.id == roomId.value)
    if (detail) {
      roomDetail.value = detail
      roomImages.value = detail.images || [detail.image_url || detail.image]
    }
  } catch (error) {
    console.error('Failed to load room detail:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 140rpx;
}

.image-section {
  position: relative;
  width: 100%;
  height: 500rpx;
  
  .swiper-item {
    width: 100%;
    height: 100%;
    
    .room-image {
      width: 100%;
      height: 100%;
    }
  }
  
  .back-btn {
    position: absolute;
    top: 88rpx;
    left: 32rpx;
    width: 72rpx;
    height: 72rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: rgba(0, 0, 0, 0.3);
    border-radius: 50%;
    z-index: 10;
    backdrop-filter: blur(10px);
    
    &:active {
      transform: scale(0.9);
    }
  }
  
  .image-count {
    position: absolute;
    bottom: 24rpx;
    right: 24rpx;
    padding: 8rpx 20rpx;
    background: rgba(0, 0, 0, 0.5);
    border-radius: 100rpx;
    
    text {
      font-size: 24rpx;
      color: #fff;
    }
  }
}

.content-scroll {
  height: calc(100vh - 500rpx - 140rpx);
}

.room-info-card {
  margin: -40rpx 24rpx 24rpx;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 10;
  
  .room-header {
    margin-bottom: 24rpx;
    
    .room-name {
      font-size: 40rpx;
      font-weight: 700;
      color: #333;
      display: block;
      margin-bottom: 16rpx;
    }
    
    .room-tags {
      display: flex;
      gap: 12rpx;
      flex-wrap: wrap;
      
      .tag {
        padding: 8rpx 16rpx;
        background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
        border-radius: 8rpx;
        font-size: 22rpx;
        color: #C29D71;
      }
    }
  }
  
  .room-price {
    display: flex;
    align-items: baseline;
    margin-bottom: 24rpx;
    
    .currency {
      font-size: 28rpx;
      color: #C29D71;
      font-weight: 600;
    }
    
    .price {
      font-size: 56rpx;
      font-weight: 700;
      color: #C29D71;
      font-family: -apple-system, BlinkMacSystemFont, sans-serif;
    }
    
    .unit {
      font-size: 24rpx;
      color: #999;
      margin-left: 4rpx;
    }
    
    .original-price {
      font-size: 26rpx;
      color: #ccc;
      text-decoration: line-through;
      margin-left: 16rpx;
    }
  }
  
  .room-specs {
    display: flex;
    justify-content: space-between;
    padding-top: 24rpx;
    border-top: 1rpx solid #F0F0F0;
    
    .spec-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8rpx;
      
      text {
        font-size: 24rpx;
        color: #666;
      }
    }
  }
}

.date-card {
  margin: 0 24rpx 24rpx;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .date-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 24rpx;
    
    .title {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
    }
  }
  
  .date-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    
    .date-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 8rpx;
      
      .label {
        font-size: 24rpx;
        color: #999;
      }
      
      .date {
        font-size: 32rpx;
        font-weight: 600;
        color: #333;
      }
    }
    
    .nights {
      padding: 12rpx 24rpx;
      background: linear-gradient(135deg, rgba(194, 157, 113, 0.15) 0%, rgba(194, 157, 113, 0.08) 100%);
      border-radius: 100rpx;
      
      text {
        font-size: 26rpx;
        color: #C29D71;
        font-weight: 500;
      }
    }
  }
}

.facility-card,
.notice-card {
  margin: 0 24rpx 24rpx;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .card-title {
    font-size: 32rpx;
    font-weight: 600;
    color: #333;
    margin-bottom: 24rpx;
  }
}

.facility-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 24rpx;
  
  .facility-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12rpx;
    
    text {
      font-size: 24rpx;
      color: #666;
    }
  }
}

.notice-list {
  .notice-item {
    display: flex;
    justify-content: space-between;
    padding: 20rpx 0;
    border-bottom: 1rpx solid #F5F5F5;
    
    &:last-child {
      border-bottom: none;
    }
    
    .label {
      font-size: 28rpx;
      color: #666;
    }
    
    .value {
      font-size: 28rpx;
      color: #333;
      font-weight: 500;
    }
  }
}

.bottom-placeholder {
  height: 40rpx;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  height: 140rpx;
  padding: 0 32rpx;
  padding-bottom: env(safe-area-inset-bottom);
  background: #fff;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 -4rpx 16rpx rgba(0, 0, 0, 0.06);
  
  .price-info {
    display: flex;
    flex-direction: column;
    gap: 4rpx;
    
    .total-price {
      display: flex;
      align-items: baseline;
      
      .currency {
        font-size: 28rpx;
        color: #C29D71;
        font-weight: 600;
      }
      
      .amount {
        font-size: 48rpx;
        font-weight: 700;
        color: #C29D71;
        font-family: -apple-system, BlinkMacSystemFont, sans-serif;
      }
    }
    
    .price-desc {
      font-size: 24rpx;
      color: #999;
    }
  }
  
  .book-btn {
    padding: 24rpx 64rpx;
    background: linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%);
    border-radius: 100rpx;
    box-shadow: 0 12rpx 24rpx rgba(194, 157, 113, 0.35);
    
    text {
      font-size: 32rpx;
      font-weight: 600;
      color: #fff;
    }
    
    &:active {
      transform: scale(0.95);
    }
  }
}
</style>

