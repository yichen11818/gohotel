<template>
  <view class="container">
    <!-- 结果图标 -->
    <view class="result-section">
      <view class="icon-wrapper" :class="isSuccess ? 'success' : 'fail'">
        <TnIcon :name="isSuccess ? 'check' : 'close'" color="#fff" size="80" />
      </view>
      <text class="result-text">{{ isSuccess ? '支付成功' : '支付失败' }}</text>
      <text class="result-desc">{{ isSuccess ? '订单已确认，期待您的入住' : '支付遇到问题，请重试' }}</text>
    </view>

    <!-- 订单信息 -->
    <view class="order-card" v-if="isSuccess">
      <view class="order-row">
        <text class="label">订单编号</text>
        <text class="value">{{ orderId }}</text>
      </view>
      <view class="order-row">
        <text class="label">支付金额</text>
        <view class="price">
          <text class="currency">¥</text>
          <text class="amount">{{ amount }}</text>
        </view>
      </view>
      <view class="order-row">
        <text class="label">支付时间</text>
        <text class="value">{{ payTime }}</text>
      </view>
    </view>

    <!-- 温馨提示 -->
    <view class="tips-card" v-if="isSuccess">
      <view class="tips-title">
        <TnIcon name="info-circle" color="#C29D71" size="32" />
        <text>温馨提示</text>
      </view>
      <view class="tips-list">
        <text class="tip-item">• 入住时间为14:00后，请携带有效身份证件</text>
        <text class="tip-item">• 离店时间为12:00前，如需延迟请提前联系前台</text>
        <text class="tip-item">• 如有任何问题，请拨打酒店电话：027-12345678</text>
      </view>
    </view>

    <!-- 操作按钮 -->
    <view class="action-section">
      <view v-if="isSuccess" class="action-buttons">
        <view class="secondary-btn" @click="goToOrderDetail">
          <text>查看订单</text>
        </view>
        <view class="primary-btn" @click="goToHome">
          <text>返回首页</text>
        </view>
      </view>
      <view v-else class="action-buttons">
        <view class="secondary-btn" @click="goBack">
          <text>返回</text>
        </view>
        <view class="primary-btn" @click="retryPay">
          <text>重新支付</text>
        </view>
      </view>
    </view>

    <!-- 底部推荐 -->
    <view class="recommend-section" v-if="isSuccess">
      <view class="section-header">
        <text class="section-title">猜你喜欢</text>
      </view>
      <scroll-view class="recommend-scroll" scroll-x>
        <view class="recommend-list">
          <view 
            v-for="item in recommendList" 
            :key="item.id" 
            class="recommend-item"
            @click="goToRoomDetail(item.id)"
          >
            <image class="recommend-image" :src="item.image" mode="aspectFill"></image>
            <view class="recommend-info">
              <text class="recommend-name">{{ item.name }}</text>
              <view class="recommend-price">
                <text class="currency">¥</text>
                <text class="price">{{ item.price }}</text>
                <text class="unit">起</text>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>
  </view>
</template>

<script setup>
import { hotel } from '@/api/index.js'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'

const isSuccess = ref(true)
const orderId = ref('')
const amount = ref(0)
const payTime = ref('')
const recommendList = ref([])

// 加载推荐房型
const loadRecommend = async () => {
  try {
    const data = await hotel.getRoomTypes(1) // 暂时硬编码
    recommendList.value = data?.slice(0, 3).map(item => ({
      id: item.id,
      name: item.name,
      image: item.image_url || item.image,
      price: item.price
    })) || []
  } catch (error) {
    console.error('Failed to load recommends:', error)
  }
}

// 格式化当前时间
const formatCurrentTime = () => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  const seconds = String(now.getSeconds()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
}

// 返回首页
const goToHome = () => {
  uni.switchTab({
    url: '/pages/index/index'
  })
}

// 查看订单详情
const goToOrderDetail = () => {
  uni.navigateTo({
    url: `/pages/order-detail/order-detail?id=${orderId.value}`
  })
}

// 返回
const goBack = () => {
  uni.navigateBack()
}

// 重新支付
const retryPay = () => {
  uni.navigateBack()
}

// 跳转到房间详情
const goToRoomDetail = (id) => {
  uni.navigateTo({
    url: `/pages/room-detail/room-detail?id=${id}`
  })
}

onLoad((options) => {
  isSuccess.value = options?.success === 'true'
  orderId.value = options?.orderId || ''
  amount.value = options?.amount || 0
  payTime.value = formatCurrentTime()
  loadRecommend()
})
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #F5F5F5;
  padding: 48rpx 0;
}

.result-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 80rpx 0;
  
  .icon-wrapper {
    width: 160rpx;
    height: 160rpx;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 32rpx;
    animation: scaleIn 0.5s ease-out;
    
    &.success {
      background: linear-gradient(135deg, #34C759 0%, #30B350 100%);
      box-shadow: 0 16rpx 48rpx rgba(52, 199, 89, 0.35);
    }
    
    &.fail {
      background: linear-gradient(135deg, #FF3B30 0%, #E63227 100%);
      box-shadow: 0 16rpx 48rpx rgba(255, 59, 48, 0.35);
    }
  }
  
  .result-text {
    font-size: 44rpx;
    font-weight: 700;
    color: #333;
    margin-bottom: 16rpx;
  }
  
  .result-desc {
    font-size: 28rpx;
    color: #999;
  }
}

@keyframes scaleIn {
  0% {
    transform: scale(0);
    opacity: 0;
  }
  50% {
    transform: scale(1.1);
  }
  100% {
    transform: scale(1);
    opacity: 1;
  }
}

.order-card {
  margin: 0 32rpx 24rpx;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .order-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16rpx 0;
    
    &:not(:last-child) {
      border-bottom: 1rpx solid #F5F5F5;
    }
    
    .label {
      font-size: 28rpx;
      color: #999;
    }
    
    .value {
      font-size: 28rpx;
      color: #333;
    }
    
    .price {
      display: flex;
      align-items: baseline;
      
      .currency {
        font-size: 24rpx;
        color: #C29D71;
        font-weight: 600;
      }
      
      .amount {
        font-size: 36rpx;
        font-weight: 700;
        color: #C29D71;
      }
    }
  }
}

.tips-card {
  margin: 0 32rpx 24rpx;
  padding: 32rpx;
  background: linear-gradient(135deg, rgba(194, 157, 113, 0.1) 0%, rgba(194, 157, 113, 0.05) 100%);
  border-radius: 24rpx;
  
  .tips-title {
    display: flex;
    align-items: center;
    gap: 12rpx;
    margin-bottom: 20rpx;
    
    text {
      font-size: 30rpx;
      font-weight: 600;
      color: #333;
    }
  }
  
  .tips-list {
    display: flex;
    flex-direction: column;
    gap: 12rpx;
    
    .tip-item {
      font-size: 26rpx;
      color: #666;
      line-height: 1.6;
    }
  }
}

.action-section {
  padding: 0 32rpx;
  margin-bottom: 48rpx;
  
  .action-buttons {
    display: flex;
    gap: 24rpx;
  }
  
  .secondary-btn {
    flex: 1;
    height: 96rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: #fff;
    border-radius: 100rpx;
    border: 2rpx solid #E5E5E5;
    
    text {
      font-size: 32rpx;
      color: #666;
      font-weight: 500;
    }
    
    &:active {
      background: #F9F9F9;
    }
  }
  
  .primary-btn {
    flex: 1;
    height: 96rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #D4B184 0%, #C29D71 50%, #B88A5E 100%);
    border-radius: 100rpx;
    box-shadow: 0 12rpx 24rpx rgba(194, 157, 113, 0.35);
    
    text {
      font-size: 32rpx;
      color: #fff;
      font-weight: 600;
    }
    
    &:active {
      transform: scale(0.98);
    }
  }
}

.recommend-section {
  margin-top: 24rpx;
  
  .section-header {
    padding: 0 32rpx 24rpx;
    
    .section-title {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
    }
  }
  
  .recommend-scroll {
    width: 100%;
    white-space: nowrap;
  }
  
  .recommend-list {
    display: inline-flex;
    gap: 20rpx;
    padding: 0 32rpx 32rpx;
  }
  
  .recommend-item {
    width: 280rpx;
    background: #fff;
    border-radius: 20rpx;
    overflow: hidden;
    box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
    
    .recommend-image {
      width: 100%;
      height: 180rpx;
    }
    
    .recommend-info {
      padding: 20rpx;
      
      .recommend-name {
        display: block;
        font-size: 28rpx;
        font-weight: 600;
        color: #333;
        margin-bottom: 8rpx;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
      
      .recommend-price {
        display: flex;
        align-items: baseline;
        
        .currency {
          font-size: 22rpx;
          color: #C29D71;
        }
        
        .price {
          font-size: 32rpx;
          font-weight: 700;
          color: #C29D71;
        }
        
        .unit {
          font-size: 22rpx;
          color: #999;
          margin-left: 4rpx;
        }
      }
    }
    
    &:active {
      transform: scale(0.98);
    }
  }
}
</style>

