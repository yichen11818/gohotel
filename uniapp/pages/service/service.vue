<template>
  <scroll-view class="page" scroll-y>
    <view class="content">
      <view class="hero-card">
        <text class="title">酒店服务</text>
        <text class="subtitle">当前版本提供真实可用的联系与入住说明</text>
        <view class="action-row">
          <button class="primary-btn" @click="callHotel" :disabled="!hotelInfo.phone">联系前台</button>
          <button class="ghost-btn" @click="copyAddress">复制地址</button>
        </view>
      </view>

      <view class="section">
        <text class="section-title">基础信息</text>
        <view class="info-row">
          <text>前台电话</text>
          <text>{{ hotelInfo.phone || '暂未配置' }}</text>
        </view>
        <view class="info-row">
          <text>服务时间</text>
          <text>{{ hotelInfo.serviceTime }}</text>
        </view>
        <view class="info-row">
          <text>入住 / 退房</text>
          <text>{{ hotelInfo.checkInTime }} / {{ hotelInfo.checkOutTime }}</text>
        </view>
        <view class="info-row">
          <text>酒店地址</text>
          <text class="multiline">{{ hotelInfo.address }}</text>
        </view>
      </view>

      <view class="section" v-if="notices.length">
        <text class="section-title">服务公告</text>
        <view v-for="notice in notices" :key="notice.id" class="notice-item">
          <text>{{ notice.title }}</text>
        </view>
      </view>

      <view class="section">
        <text class="section-title">入住提示</text>
        <text class="tip-text">- 提交预订后，订单状态会在“我的订单”中更新。</text>
        <text class="tip-text">- 如需帮助，请直接联系酒店前台。</text>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { hotel } from '@/api/index.js'

const hotelInfo = ref({
  phone: '',
  address: '欢迎入住七天酒店',
  serviceTime: '24小时服务',
  checkInTime: '14:00',
  checkOutTime: '12:00',
})
const notices = ref([])

const loadData = async () => {
  try {
    const detail = await hotel.getHotelDetail()
    hotelInfo.value = detail
    notices.value = detail.notices || []
  } catch (error) {
    console.error('load service data failed:', error)
  }
}

const callHotel = () => {
  if (!hotelInfo.value.phone) {
    uni.showToast({ title: '暂未配置联系电话', icon: 'none' })
    return
  }
  uni.makePhoneCall({ phoneNumber: hotelInfo.value.phone })
}

const copyAddress = () => {
  uni.setClipboardData({
    data: hotelInfo.value.address,
    success: () => {
      uni.showToast({ title: '地址已复制', icon: 'success' })
    },
  })
}

onShow(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: #f6f7fb;
}

.content {
  padding: 24rpx;
}

.hero-card,
.section {
  padding: 28rpx;
  margin-bottom: 24rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.hero-card {
  background: linear-gradient(135deg, #f6efe7 0%, #ffffff 100%);
}

.title,
.section-title {
  display: block;
  font-size: 32rpx;
  font-weight: 700;
  color: #111827;
}

.subtitle,
.tip-text {
  display: block;
  margin-top: 12rpx;
  font-size: 24rpx;
  color: #4b5563;
  line-height: 1.7;
}

.action-row {
  display: flex;
  gap: 20rpx;
  margin-top: 28rpx;
}

.primary-btn,
.ghost-btn {
  margin: 0;
  border-radius: 999rpx;
}

.primary-btn {
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 16rpx;
  padding: 18rpx 0;
  font-size: 26rpx;
  color: #374151;
  border-bottom: 1px solid #f3f4f6;
}

.info-row:last-child {
  border-bottom: 0;
}

.multiline {
  flex: 1;
  text-align: right;
  line-height: 1.6;
}

.notice-item {
  padding: 18rpx 0;
  font-size: 26rpx;
  color: #374151;
  border-bottom: 1px solid #f3f4f6;
}

.notice-item:last-child {
  border-bottom: 0;
}
</style>
