<template>
  <view class="page">
    <scroll-view class="scroll-container" scroll-y>
      <!-- 头部：管家 Banner -->
      <view class="hero-section">
        <view class="hero-bg"></view>
        <view class="hero-content">
          <view class="hero-text">
            <text class="title">您好，您的随身管家</text>
            <text class="subtitle">24小时竭诚为您提供周到服务</text>
          </view>
          <view class="hero-badge flex-center">
            <image class="icon" src="/static/icons/service.png" mode="aspectFit" />
            <text>24h 在线</text>
          </view>
        </view>
      </view>

      <view class="content">
        <!-- 核心快捷服务 -->
        <view class="grid-card premium-card">
          <view class="grid-item flex-center" @click="callHotel">
            <view class="icon-box call"><image src="/static/icons/phone.png" mode="aspectFit" /></view>
            <text class="label">联系前台</text>
          </view>
          <view class="grid-item flex-center" @click="copyAddress">
            <view class="icon-box map"><image src="/static/icons/location.png" mode="aspectFit" /></view>
            <text class="label">酒店位置</text>
          </view>
          <view class="grid-item flex-center" @click="showWifi">
            <view class="icon-box wifi"><image src="/static/icons/menu-service.png" mode="aspectFit" /></view>
            <text class="label">WiFi信息</text>
          </view>
          <view class="grid-item flex-center" @click="handleService('luggage')">
            <view class="icon-box luggage"><image src="/static/icons/menu-safe.png" mode="aspectFit" /></view>
            <text class="label">行李寄送</text>
          </view>
        </view>

        <!-- 入住指南 -->
        <view class="section premium-card">
          <view class="section-header">
            <text class="section-title">服务指南</text>
          </view>
          <view class="guide-list">
            <view class="guide-item flex-between">
              <view class="item-left flex-center">
                <view class="dot"></view>
                <text class="label">入住时间</text>
              </view>
              <text class="val">{{ hotelInfo.checkInTime || '14:00' }} 以后</text>
            </view>
            <view class="guide-item flex-between">
              <view class="item-left flex-center">
                <view class="dot"></view>
                <text class="label">退房时间</text>
              </view>
              <text class="val">{{ hotelInfo.checkOutTime || '12:00' }} 以前</text>
            </view>
            <view class="guide-item flex-between">
              <view class="item-left flex-center">
                <view class="dot"></view>
                <text class="label">酒店地址</text>
              </view>
              <text class="val address-text">{{ hotelInfo.address }}</text>
            </view>
          </view>
        </view>

        <!-- 酒店公告 -->
        <view class="section notice-section premium-card" v-if="notices.length">
          <view class="section-header flex-between">
            <view class="flex-center">
              <text class="section-title">最新公告</text>
              <view class="badge-new">NEW</view>
            </view>
          </view>
          <swiper class="notice-swiper" vertical autoplay circular interval="4000">
            <swiper-item v-for="notice in notices" :key="notice.id">
              <view class="notice-item flex-center">
                <image class="icon-notice" src="/static/icons/order.png" mode="aspectFit" />
                <text class="notice-title text-ellipsis">{{ notice.title }}</text>
              </view>
            </swiper-item>
          </swiper>
        </view>

        <!-- 酒店设施看板 -->
        <view class="section premium-card">
          <view class="section-header">
            <text class="section-title">设施开放时间</text>
          </view>
          <view class="facility-grid">
            <view class="f-item flex-between">
              <text class="f-name">早餐厅 (2F)</text>
              <text class="f-time">07:00 - 10:00</text>
            </view>
            <view class="f-item flex-between">
              <text class="f-name">健身房 (3F)</text>
              <text class="f-time">09:00 - 22:00</text>
            </view>
            <view class="f-item flex-between">
              <text class="f-name">商务中心 (1F)</text>
              <text class="f-time">24小时开放</text>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>

    <!-- 底部一键客服 -->
    <view class="bottom-action safe-area-inset-bottom">
      <button class="service-btn premium-button flex-center" open-type="contact" @click="goServiceChat">
        <image class="icon" src="/static/icons/service.png" mode="aspectFit" style="filter: brightness(100);" />
        <text>咨询在线管家</text>
      </button>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { hotel } from '@/api/index.js'

const hotelInfo = ref({
  phone: '',
  address: '欢迎入住七天酒店',
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

const showWifi = () => {
  uni.showModal({
    title: '酒店WiFi信息',
    content: '账号：7Days_Guest\n密码：7days8888',
    confirmText: '去连接',
    confirmColor: '#C29D71',
  })
}

const handleService = (type) => {
  uni.showToast({ title: '服务申请已提交，请等待处理', icon: 'none' })
}

const goServiceChat = () => {
  // 模拟跳转到客服聊天
}

onShow(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: $bg-color;
  display: flex;
  flex-direction: column;
}

.scroll-container {
  flex: 1;
  height: 0;
}

/* 头部沉浸式 */
.hero-section {
  height: 400rpx;
  position: relative;
  overflow: hidden;

  .hero-bg {
    position: absolute;
    inset: 0;
    background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
    clip-path: ellipse(110% 70% at 50% 20%);
  }

  .hero-content {
    position: relative;
    z-index: 10;
    padding: 80rpx 40rpx;
    display: flex;
    justify-content: space-between;
    align-items: center;

    .hero-text {
      .title {
        font-size: 44rpx;
        font-weight: 700;
        color: #fff;
        display: block;
      }
      .subtitle {
        font-size: 24rpx;
        color: rgba(255,255,255,0.8);
        margin-top: 12rpx;
        display: block;
      }
    }

    .hero-badge {
      background: rgba(255,255,255,0.2);
      backdrop-filter: blur(8rpx);
      padding: 12rpx 24rpx;
      border-radius: 100rpx;
      .icon { width: 32rpx; height: 32rpx; margin-right: 8rpx; }
      text { font-size: 20rpx; color: #fff; font-weight: 500; }
    }
  }
}

.content {
  padding: 0 30rpx;
  margin-top: -140rpx;
  position: relative;
  z-index: 20;
}

/* 快捷服务 Grid */
.grid-card {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  padding: 40rpx 20rpx;
  margin-bottom: 30rpx;

  .grid-item {
    flex-direction: column;
    .icon-box {
      width: 88rpx;
      height: 88rpx;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      margin-bottom: 16rpx;
      image { width: 44rpx; height: 44rpx; }

      &.call { background: #EBF7FF; image { filter: invert(42%) sepia(82%) saturate(1352%) hue-rotate(185deg) brightness(101%) contrast(101%); } }
      &.map { background: #E7F9F1; image { filter: invert(53%) sepia(76%) saturate(415%) hue-rotate(104deg) brightness(92%) contrast(89%); } }
      &.wifi { background: #FFF4ED; image { filter: invert(65%) sepia(45%) saturate(542%) hue-rotate(345deg) brightness(101%) contrast(92%); } }
      &.luggage { background: #F3F1FF; image { filter: invert(62%) sepia(91%) saturate(301%) hue-rotate(212deg) brightness(91%) contrast(95%); } }
    }
    .label { font-size: 24rpx; color: $text-main; font-weight: 500; }
  }
}

.section {
  padding: 40rpx;
  margin-bottom: 30rpx;

  .section-header {
    margin-bottom: 30rpx;
    display: flex;
    align-items: center;
    .section-title {
      font-size: 32rpx;
      font-weight: 700;
      color: $text-main;
    }
    .badge-new {
      background: #E64340;
      color: #fff;
      font-size: 18rpx;
      padding: 2rpx 10rpx;
      border-radius: 4rpx;
      margin-left: 12rpx;
    }
  }
}

/* 服务指南列表 */
.guide-list {
  .guide-item {
    padding: 24rpx 0;
    border-bottom: 1rpx solid #f8f8f8;
    &:last-child { border-bottom: 0; }

    .item-left {
      .dot { width: 8rpx; height: 8rpx; border-radius: 50%; background: $primary-color; margin-right: 16rpx; }
      .label { font-size: 28rpx; color: $text-sub; }
    }
    .val {
      font-size: 28rpx; color: $text-main; font-weight: 500;
      &.address-text { max-width: 360rpx; text-align: right; line-height: 1.6; }
    }
  }
}

/* 公告 Swiper */
.notice-swiper {
  height: 60rpx;
  .notice-item {
    height: 100%;
    justify-content: flex-start;
    .icon-notice { width: 32rpx; height: 32rpx; margin-right: 16rpx; opacity: 0.5; }
    .notice-title { font-size: 26rpx; color: $text-second; }
  }
}

/* 设施 Grid */
.facility-grid {
  .f-item {
    padding: 20rpx 24rpx;
    background: #f9f9f9;
    border-radius: 12rpx;
    margin-bottom: 16rpx;
    &:last-child { margin-bottom: 0; }
    .f-name { font-size: 26rpx; color: $text-main; font-weight: 500; }
    .f-time { font-size: 24rpx; color: $primary-color; font-weight: 600; }
  }
}

/* 底部操作 */
.bottom-action {
  padding: 24rpx 40rpx calc(24rpx + env(safe-area-inset-bottom));
  background: #fff;
  box-shadow: 0 -4rpx 20rpx rgba(0,0,0,0.05);

  .service-btn {
    height: 100rpx;
    font-size: 32rpx;
    .icon { width: 44rpx; height: 44rpx; margin-right: 16rpx; }
  }
}
</style>
