<template>
  <view class="container">
    <TnNavbar :fixed="true" :placeholder="true" title="客房列表" bg-color="transparent" :bottom-shadow="false" />
    
    <view class="room-list">
      <template v-if="loading">
        <view v-for="i in 3" :key="i" class="skeleton-card"></view>
      </template>
      <template v-else-if="roomList.length > 0">
        <view v-for="(room, index) in roomList" :key="index" class="room-card" @click="handleRoomClick(room)">
          <!-- 图片区域 -->
          <view class="image-wrapper">
            <image :src="room.image_url || room.image" mode="aspectFill" class="room-image"></image>
            <view class="image-mask"></view>
            <text class="room-name">{{ room.name }}</text>
          </view>
          
          <!-- 内容区域 -->
          <view class="content-wrapper">
            <!-- 描述/地址行 -->
            <view class="info-row location">
              <text class="tn-icon-location-fill icon"></text>
              <text class="text">{{ room.description || room.desc }}</text>
            </view>
            
            <!-- 评分/标签行 -->
            <view class="info-row scores">
              <text class="score">{{ room.score || '5.0' }}分</text>
              <text class="sub-text">{{ room.booked_count || '0' }}+消费</text>
              <text class="sub-text">{{ room.comment_count || '0' }}评论</text>
            </view>
            
            <!-- 价格行 -->
            <view class="price-row">
              <text class="currency">¥</text>
              <text class="amount">{{ room.price }}</text>
              <text class="suffix">起</text>
            </view>
          </view>
        </view>
      </template>
      <view v-else class="empty-state">
        <text>暂无可用房型</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { hotel } from '@/api/index.js'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnButton from '@/uni_modules/tuniaoui-vue3/components/button/src/button.vue'
import { onLoad } from '@dcloudio/uni-app'

const roomList = ref([])
const loading = ref(false)
const hotelId = ref(1)

onLoad((options) => {
  if (options?.id) {
    hotelId.value = options.id
  }
  loadRoomTypes()
})

const loadRoomTypes = async () => {
  loading.value = true
  try {
    // 使用新的后端API获取可用房间列表
    const data = await hotel.getAvailableRooms({
      page: 1,
      page_size: 20
    })
    
    // 处理返回的数据结构
    if (data && Array.isArray(data)) {
      roomList.value = data.map(room => ({
        id: room.id,
        name: room.room_number || `${room.room_type}房间`,
        description: room.room_type || room.description,
        price: room.price,
        image_url: room.image_url,
        score: room.score || '5.0',
        booked_count: room.booked_count || 0,
        comment_count: room.comment_count || 0
      }))
    } else if (data && data.data) {
      // 如果返回的是分页数据结构
      roomList.value = data.data.map(room => ({
        id: room.id,
        name: room.room_number || `${room.room_type}房间`,
        description: room.room_type || room.description,
        price: room.price,
        image_url: room.image_url,
        score: room.score || '5.0',
        booked_count: room.booked_count || 0,
        comment_count: room.comment_count || 0
      }))
    } else {
      roomList.value = []
    }
  } catch (error) {
    console.error('Failed to load room types:', error)
    throw error
  } finally {
    loading.value = false
  }
}

const handleRoomClick = (room) => {
  // 跳转到房间详情页
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${room.id}`
  })
}
</script>

<style lang="scss" scoped>
// 文本省略混入
@mixin text-ellipsis-1 {
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.container {
  min-height: 100vh;
  background-color: #F9F9F9;
}

.room-list {
  padding: 32rpx;
}

.room-card {
  background: #ffffff;
  border-radius: 32rpx;
  box-shadow: 0 8rpx 32rpx rgba(0, 0, 0, 0.04);
  border: 1rpx solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
  margin-bottom: 32rpx;
  transition: all 0.3s ease;
  
  &:active {
    transform: scale(0.98);
  }

  .image-wrapper {
    position: relative;
    width: 100%;
    height: 360rpx;
    
    .room-image {
      width: 100%;
      height: 100%;
    }
    
    .image-mask {
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 70%;
      background: linear-gradient(to top, rgba(0,0,0,0.7) 0%, transparent 100%);
      z-index: 1;
    }
    
    .room-name {
      position: absolute;
      bottom: 24rpx;
      left: 32rpx;
      color: #fff;
      font-size: 36rpx;
      font-weight: 600;
      z-index: 2;
      letter-spacing: 1rpx;
    }
  }
  
  .content-wrapper {
    padding: 32rpx;
    position: relative;
    
    .info-row {
      display: flex;
      align-items: center;
      margin-bottom: 12rpx;
      
      &.location {
        color: #666;
        font-size: 26rpx;
        
        .icon {
          font-size: 28rpx;
          margin-right: 8rpx;
          color: #C29D71;
        }
      }
      
      &.scores {
        .score {
          color: #C29D71;
          font-weight: 600;
          font-size: 30rpx;
          margin-right: 16rpx;
        }
        
        .sub-text {
          color: #999;
          font-size: 24rpx;
          margin-right: 16rpx;
          background: #f5f5f5;
          padding: 4rpx 12rpx;
          border-radius: 4rpx;
        }
      }
    }
    
    .price-row {
      position: absolute;
      right: 32rpx;
      bottom: 32rpx;
      display: flex;
      align-items: baseline;
      
      .currency {
        color: #333;
        font-size: 24rpx;
        margin-right: 4rpx;
        font-weight: 600;
      }
      
      .amount {
        color: #333;
        font-size: 48rpx;
        font-weight: 700;
        font-family: 'Din', sans-serif;
      }
      
      .suffix {
        color: #999;
        font-size: 22rpx;
        margin-left: 4rpx;
      }
    }
  }
}

.skeleton-card {
  height: 500rpx;
  background: #eee;
  border-radius: 24rpx;
  margin-bottom: 32rpx;
  animation: skeleton-blink 1.5s infinite;
}

@keyframes skeleton-blink {
  0% { opacity: 0.6; }
  50% { opacity: 0.3; }
  100% { opacity: 0.6; }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
  color: #999;
  font-size: 28rpx;
}
</style>
