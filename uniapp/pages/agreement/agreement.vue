<template>
  <scroll-view class="page" scroll-y>
    <view class="content">
      <view class="section">
        <text class="title">用户协议</text>
        <text class="paragraph">注册与登录即表示你同意在本平台内使用七天酒店提供的真实在线服务，包括账号登录、房型浏览、提交预订、订单查询与个人资料维护。</text>
        <text class="paragraph">请确保提交的入住人姓名、手机号等信息真实有效；若因信息错误导致入住失败，酒店有权要求重新核验。</text>
        <text class="paragraph">当前版本暂不支持在线支付、优惠券、收藏、评价等能力，相关入口不会计入本协议约定范围。</text>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
</script>

<style scoped lang="scss">
.page { min-height: 100vh; background: #f6f7fb; }
.content { padding: 24rpx; }
.section { padding: 28rpx; background: #fff; border-radius: 24rpx; box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05); }
.title { display: block; font-size: 32rpx; font-weight: 700; color: #111827; }
.paragraph { display: block; margin-top: 16rpx; font-size: 26rpx; color: #4b5563; line-height: 1.8; }
</style>
