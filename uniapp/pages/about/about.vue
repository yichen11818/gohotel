<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">关于我们</text>
      <view class="placeholder-view"></view>
    </view>

    <scroll-view class="content" scroll-y>
      <!-- 品牌 Logo 形象区 -->
      <view class="brand-section flex-center">
        <view class="logo-wrapper">
          <image class="logo" src="https://images.unsplash.com/photo-1566073771259-6a8506099945?w=200&h=200&fit=crop" mode="aspectFill" />
        </view>
        <text class="brand-name">七天酒店</text>
        <text class="brand-version">Version 1.2.0</text>
      </view>

      <!-- 酒店介绍 -->
      <view class="section premium-card">
        <view class="section-header">
          <text class="title">品牌故事</text>
        </view>
        <text class="paragraph">
          七天酒店致力于为每一位住客提供高品质、人性化的住宿体验。我们拥有舒适的客房设计、先进的服务设施以及专业的服务团队，确保您的每一次下榻都能感受到家一般的温暖与尊贵。
        </text>
        <text class="paragraph">
          无论是商务出差还是休闲旅游，七天酒店都是您的理想之选。
        </text>
      </view>

      <!-- 联系我们 -->
      <view class="section premium-card">
        <view class="section-header">
          <text class="title">联系我们</text>
        </view>
        <view class="contact-list">
          <view class="contact-item flex-between" @click="makeCall">
            <text class="label">客服电话</text>
            <text class="val color-primary">400-888-7777</text>
          </view>
          <view class="contact-item flex-between">
            <text class="label">官方邮箱</text>
            <text class="val">service@7dayshotel.com</text>
          </view>
          <view class="contact-item flex-between">
            <text class="label">官方网址</text>
            <text class="val">www.7dayshotel.com</text>
          </view>
        </view>
      </view>

      <view class="footer-info">
        <text>© 2026 七天酒店管理有限公司 版权所有</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'

const statusBarHeight = ref(44)

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
})

const goBack = () => {
  uni.navigateBack()
}

const makeCall = () => {
  uni.makePhoneCall({ phoneNumber: '400-888-7777' })
}
</script>

<style scoped lang="scss">
.page {
  min-height: 100vh;
  background: $bg-color;
  display: flex;
  flex-direction: column;
}

.nav-header {
  height: 88rpx;
  padding: 0 30rpx;
  background: #fff;

  .back-btn {
    width: 64rpx;
    height: 64rpx;
    .icon-back { width: 40rpx; height: 40rpx; }
  }

  .page-title {
    font-size: 34rpx;
    font-weight: 700;
    color: $text-main;
  }

  .placeholder-view { width: 64rpx; }
}

.content {
  flex: 1;
  height: 0;
  padding: 0 30rpx;
}

/* 品牌区 */
.brand-section {
  padding: 80rpx 0;
  flex-direction: column;

  .logo-wrapper {
    width: 160rpx;
    height: 160rpx;
    border-radius: 40rpx;
    overflow: hidden;
    box-shadow: 0 12rpx 30rpx rgba(194, 157, 113, 0.2);
    margin-bottom: 30rpx;
    background: #fff;
    .logo { width: 100%; height: 100%; }
  }

  .brand-name {
    font-size: 40rpx;
    font-weight: 700;
    color: $text-main;
    margin-bottom: 12rpx;
  }

  .brand-version {
    font-size: 24rpx;
    color: $text-sub;
  }
}

.section {
  padding: 40rpx;
  margin-bottom: 30rpx;

  .section-header {
    margin-bottom: 24rpx;
    .title {
      font-size: 32rpx;
      font-weight: 700;
      color: $text-main;
    }
  }

  .paragraph {
    font-size: 28rpx;
    color: $text-second;
    line-height: 1.8;
    margin-bottom: 20rpx;
    display: block;
    &:last-child { margin-bottom: 0; }
  }
}

/* 联系列表 */
.contact-list {
  .contact-item {
    padding: 24rpx 0;
    border-bottom: 1rpx solid #f8f8f8;
    &:last-child { border-bottom: 0; }

    .label { font-size: 28rpx; color: $text-sub; }
    .val { font-size: 28rpx; color: $text-main; font-weight: 500; }
    .color-primary { color: $primary-color; }
  }
}

.footer-info {
  margin-top: 40rpx;
  padding-bottom: 60rpx;
  text-align: center;
  font-size: 22rpx;
  color: #ccc;
}
</style>
