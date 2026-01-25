<template>
  <view class="container">
    <!-- 导航栏 -->
    <TnNavbar fixed :bottom-shadow="true" bg-color="#ffffff" :placeholder="true">
      <template #back>
        <view class="navbar-content">
          <view class="back-btn" @click="goBack">
            <TnIcon name="left" color="#333" size="40" />
          </view>
          <text class="navbar-title">确认订单</text>
        </view>
      </template>
    </TnNavbar>

    <scroll-view class="main-content" scroll-y>
      <!-- 房间信息 -->
      <view class="room-card">
        <image class="room-image" :src="roomInfo.image" mode="aspectFill"></image>
        <view class="room-info">
          <text class="room-name">{{ roomInfo.name }}</text>
          <text class="room-specs">{{ roomInfo.area }}m² · {{ roomInfo.bedType }}</text>
          <view class="room-price">
            <text class="currency">¥</text>
            <text class="price">{{ roomInfo.price }}</text>
            <text class="unit">/晚</text>
          </view>
        </view>
      </view>

      <!-- 入住日期 -->
      <view class="date-card">
        <view class="date-row">
          <view class="date-item">
            <text class="label">入住</text>
            <text class="date">{{ formatDate(checkInDate) }}</text>
            <text class="weekday">{{ getWeekday(checkInDate) }}</text>
          </view>
          <view class="nights-badge">
            <text>共{{ nights }}晚</text>
          </view>
          <view class="date-item">
            <text class="label">离店</text>
            <text class="date">{{ formatDate(checkOutDate) }}</text>
            <text class="weekday">{{ getWeekday(checkOutDate) }}</text>
          </view>
        </view>
      </view>

      <!-- 入住人信息 -->
      <view class="guest-card">
        <view class="card-header">
          <text class="card-title">入住人信息</text>
        </view>
        
        <view class="input-group">
          <text class="input-label">姓名</text>
          <input 
            v-model="guestInfo.name" 
            type="text" 
            placeholder="请输入入住人姓名"
            class="input-field"
          />
        </view>
        
        <view class="input-group">
          <text class="input-label">手机号</text>
          <input 
            v-model="guestInfo.phone" 
            type="number" 
            placeholder="请输入手机号"
            maxlength="11"
            class="input-field"
          />
        </view>
        
        <view class="input-group">
          <text class="input-label">房间数</text>
          <view class="number-selector">
            <view class="num-btn" @click="decreaseRoom">
              <TnIcon name="minus" :color="roomCount > 1 ? '#333' : '#ccc'" size="28" />
            </view>
            <text class="num-value">{{ roomCount }}</text>
            <view class="num-btn" @click="increaseRoom">
              <TnIcon name="add" color="#333" size="28" />
            </view>
          </view>
        </view>
      </view>

      <!-- 优惠券 -->
      <view class="coupon-card" @click="selectCoupon">
        <view class="coupon-left">
          <TnIcon name="ticket" color="#C29D71" size="40" />
          <text class="coupon-label">优惠券</text>
        </view>
        <view class="coupon-right">
          <text v-if="selectedCoupon" class="coupon-discount">-¥{{ selectedCoupon.amount }}</text>
          <text v-else class="coupon-hint">{{ availableCoupons }}张可用</text>
          <TnIcon name="right" color="#ccc" size="28" />
        </view>
      </view>

      <!-- 备注 -->
      <view class="remark-card">
        <text class="remark-label">备注</text>
        <textarea 
          v-model="remark" 
          placeholder="如有特殊需求请备注（选填）"
          class="remark-input"
          :maxlength="200"
        />
      </view>

      <!-- 价格明细 -->
      <view class="price-detail-card">
        <view class="card-header">
          <text class="card-title">价格明细</text>
        </view>
        
        <view class="price-row">
          <text class="price-label">房费</text>
          <text class="price-value">¥{{ roomInfo.price }} × {{ nights }}晚 × {{ roomCount }}间</text>
        </view>
        
        <view v-if="selectedCoupon" class="price-row discount">
          <text class="price-label">优惠券</text>
          <text class="price-value">-¥{{ selectedCoupon.amount }}</text>
        </view>
        
        <view class="price-row total">
          <text class="price-label">合计</text>
          <view class="total-price">
            <text class="currency">¥</text>
            <text class="amount">{{ totalPrice }}</text>
          </view>
        </view>
      </view>

      <!-- 预订须知 -->
      <view class="notice-card">
        <view class="notice-item">
          <TnIcon name="check-circle" color="#52C41A" size="28" />
          <text>入住前1天18:00前可免费取消</text>
        </view>
        <view class="notice-item">
          <TnIcon name="check-circle" color="#52C41A" size="28" />
          <text>14:00后办理入住，12:00前退房</text>
        </view>
      </view>

      <view class="bottom-placeholder"></view>
    </scroll-view>

    <!-- 底部提交栏 -->
    <view class="submit-bar">
      <view class="price-info">
        <text class="label">应付金额</text>
        <view class="price">
          <text class="currency">¥</text>
          <text class="amount">{{ totalPrice }}</text>
        </view>
      </view>
      <view class="submit-btn" @click="handleSubmit">
        <text>提交订单</text>
      </view>
    </view>
  </view>
</template>

<script setup>
import { hotel, booking, user } from '@/api/index.js'
import { TOKEN_KEY, USER_INFO_KEY } from '@/config/api.config.js'
import TnNavbar from '@/uni_modules/tuniaoui-vue3/components/navbar/src/navbar.vue'
import TnIcon from '@/uni_modules/tuniaoui-vue3/components/icon/src/icon.vue'
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

// 房间信息
const roomInfo = ref({
  id: '',
  name: '',
  image: '',
  area: 0,
  bedType: '',
  price: 0
})

// 加载房间信息
const loadRoomInfo = async (roomId) => {
  try {
    const data = await hotel.getRoomTypes(1) // 暂时硬编码酒店ID为1
    const detail = data.find(item => item.id == roomId)
    if (detail) {
      roomInfo.value = {
        id: detail.id,
        name: detail.name,
        image: detail.image_url || detail.image,
        area: detail.area,
        bedType: detail.bed_type,
        price: detail.price
      }
    }
  } catch (error) {
    console.error('Failed to load room info:', error)
  }
}

// 优惠券相关
const coupons = ref([])
const availableCoupons = computed(() => coupons.value.length)

const loadCoupons = async () => {
  try {
    const data = await user.getCoupons({ status: 'unused' })
    coupons.value = data || []
  } catch (error) {
    console.error('Failed to load coupons:', error)
  }
}

// 备注
const remark = ref('')

// 计算入住天数
const nights = computed(() => {
  const diff = checkOutDate.value.getTime() - checkInDate.value.getTime()
  return Math.ceil(diff / (1000 * 60 * 60 * 24))
})

// 计算总价
const totalPrice = computed(() => {
  let price = roomInfo.value.price * nights.value * roomCount.value
  if (selectedCoupon.value) {
    price -= selectedCoupon.value.amount
  }
  return Math.max(0, price)
})

// 格式化日期
const formatDate = (date) => {
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日`
}

// 获取星期
const getWeekday = (date) => {
  const weekdays = ['周日', '周一', '周二', '周三', '周四', '周五', '周六']
  return weekdays[date.getDay()]
}

// 增加房间
const increaseRoom = () => {
  if (roomCount.value < 5) {
    roomCount.value++
  }
}

// 减少房间
const decreaseRoom = () => {
  if (roomCount.value > 1) {
    roomCount.value--
  }
}

// 选择优惠券
const selectCoupon = () => {
  uni.navigateTo({
    url: '/pages/coupons/coupons?select=true'
  })
}

// 返回
const goBack = () => {
  uni.navigateBack()
}

// 提交订单
const handleSubmit = async () => {
  // 验证
  if (!guestInfo.value.name.trim()) {
    uni.showToast({ title: '请输入入住人姓名', icon: 'none' })
    return
  }
  
  if (!guestInfo.value.phone || guestInfo.value.phone.length !== 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  
  try {
    uni.showLoading({ title: '提交中...' })
    
    const orderData = {
      hotelId: 1, // 暂时硬编码
      roomTypeId: roomInfo.value.id,
      checkInDate: checkInDate.value.toISOString().split('T')[0],
      checkOutDate: checkOutDate.value.toISOString().split('T')[0],
      roomCount: roomCount.value,
      guestName: guestInfo.value.name,
      guestPhone: guestInfo.value.phone,
      remark: remark.value
    }
    
    const result = await booking.createBooking(orderData)
    uni.hideLoading()
    
    if (result && result.id) {
      // 模拟支付
      uni.showModal({
        title: '确认支付',
        content: `订单已创建，需支付 ¥${totalPrice.value}`,
        success: async (res) => {
          if (res.confirm) {
            try {
              uni.showLoading({ title: '支付中...' })
              await booking.payBooking(result.id, 'wechat')
              uni.hideLoading()
              
              uni.redirectTo({
                url: `/pages/pay-result/pay-result?success=true&orderId=${result.id}&amount=${totalPrice.value}`
              })
            } catch (payError) {
              uni.hideLoading()
              uni.redirectTo({
                url: `/pages/pay-result/pay-result?success=false&orderId=${result.id}&amount=${totalPrice.value}`
              })
            }
          } else {
            // 取消支付跳转到订单列表
            uni.redirectTo({ url: '/pages/orders/orders' })
          }
        }
      })
    }
  } catch (error) {
    uni.hideLoading()
  }
}

onLoad((options) => {
  // 解析参数
  if (options?.roomId) {
    loadRoomInfo(options.roomId)
  }
  if (options?.checkIn) {
    checkInDate.value = new Date(options.checkIn)
  }
  if (options?.checkOut) {
    checkOutDate.value = new Date(options.checkOut)
  }
  
  // 预填用户信息
  const userInfo = uni.getStorageSync(USER_INFO_KEY)
  if (userInfo) {
    guestInfo.value.name = userInfo.nickname || ''
    guestInfo.value.phone = userInfo.phone || ''
  }
  
  loadCoupons()
})
</script>

<style lang="scss" scoped>
.container {
  min-height: 100vh;
  background: #F5F5F5;
  padding-bottom: 140rpx;
}

.navbar-content {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60rpx;
  height: 60rpx;
  
  &:active {
    opacity: 0.6;
  }
}

.navbar-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #333;
}

.main-content {
  height: calc(100vh - 88rpx - 140rpx);
  padding: 24rpx;
}

/* 房间卡片 */
.room-card {
  display: flex;
  gap: 20rpx;
  padding: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .room-image {
    width: 180rpx;
    height: 140rpx;
    border-radius: 16rpx;
  }
  
  .room-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    
    .room-name {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
    }
    
    .room-specs {
      font-size: 26rpx;
      color: #999;
    }
    
    .room-price {
      display: flex;
      align-items: baseline;
      
      .currency {
        font-size: 24rpx;
        color: #C29D71;
        font-weight: 600;
      }
      
      .price {
        font-size: 36rpx;
        font-weight: 700;
        color: #C29D71;
      }
      
      .unit {
        font-size: 22rpx;
        color: #999;
      }
    }
  }
}

/* 日期卡片 */
.date-card {
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .date-row {
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
      
      .weekday {
        font-size: 24rpx;
        color: #666;
      }
    }
    
    .nights-badge {
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

/* 入住人卡片 */
.guest-card {
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .card-header {
    margin-bottom: 24rpx;
    
    .card-title {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
    }
  }
}

.input-group {
  display: flex;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #F5F5F5;
  
  &:last-child {
    border-bottom: none;
  }
  
  .input-label {
    width: 140rpx;
    font-size: 28rpx;
    color: #666;
  }
  
  .input-field {
    flex: 1;
    font-size: 28rpx;
    color: #333;
  }
  
  .number-selector {
    display: flex;
    align-items: center;
    gap: 32rpx;
    
    .num-btn {
      width: 56rpx;
      height: 56rpx;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #F5F5F5;
      border-radius: 50%;
      
      &:active {
        background: #EEEEEE;
      }
    }
    
    .num-value {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
      min-width: 40rpx;
      text-align: center;
    }
  }
}

/* 优惠券卡片 */
.coupon-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .coupon-left {
    display: flex;
    align-items: center;
    gap: 16rpx;
    
    .coupon-label {
      font-size: 30rpx;
      color: #333;
      font-weight: 500;
    }
  }
  
  .coupon-right {
    display: flex;
    align-items: center;
    gap: 12rpx;
    
    .coupon-discount {
      font-size: 28rpx;
      color: #FF3B30;
      font-weight: 600;
    }
    
    .coupon-hint {
      font-size: 28rpx;
      color: #C29D71;
    }
  }
  
  &:active {
    opacity: 0.7;
  }
}

/* 备注卡片 */
.remark-card {
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .remark-label {
    display: block;
    font-size: 28rpx;
    color: #666;
    margin-bottom: 16rpx;
  }
  
  .remark-input {
    width: 100%;
    height: 160rpx;
    font-size: 28rpx;
    color: #333;
    background: #F9F9F9;
    border-radius: 12rpx;
    padding: 20rpx;
    box-sizing: border-box;
  }
}

/* 价格明细卡片 */
.price-detail-card {
  padding: 32rpx;
  background: #fff;
  border-radius: 24rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.04);
  
  .card-header {
    margin-bottom: 24rpx;
    
    .card-title {
      font-size: 32rpx;
      font-weight: 600;
      color: #333;
    }
  }
  
  .price-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 16rpx 0;
    
    .price-label {
      font-size: 28rpx;
      color: #666;
    }
    
    .price-value {
      font-size: 28rpx;
      color: #333;
    }
    
    &.discount .price-value {
      color: #FF3B30;
    }
    
    &.total {
      padding-top: 24rpx;
      margin-top: 16rpx;
      border-top: 1rpx solid #F0F0F0;
      
      .price-label {
        font-size: 30rpx;
        color: #333;
        font-weight: 600;
      }
      
      .total-price {
        display: flex;
        align-items: baseline;
        
        .currency {
          font-size: 28rpx;
          color: #C29D71;
          font-weight: 600;
        }
        
        .amount {
          font-size: 44rpx;
          font-weight: 700;
          color: #C29D71;
        }
      }
    }
  }
}

/* 预订须知 */
.notice-card {
  padding: 24rpx 32rpx;
  background: rgba(82, 196, 26, 0.08);
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  
  .notice-item {
    display: flex;
    align-items: center;
    gap: 12rpx;
    padding: 8rpx 0;
    
    text {
      font-size: 26rpx;
      color: #52C41A;
    }
  }
}

.bottom-placeholder {
  height: 40rpx;
}

/* 底部提交栏 */
.submit-bar {
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
    
    .label {
      font-size: 24rpx;
      color: #999;
    }
    
    .price {
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
      }
    }
  }
  
  .submit-btn {
    padding: 24rpx 80rpx;
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

