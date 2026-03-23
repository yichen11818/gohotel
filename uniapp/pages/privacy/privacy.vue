<template>
  <scroll-view class="page" scroll-y>
    <view class="content">
      <view class="section">
        <text class="title">隐私政策</text>
        <text class="paragraph">我们仅在提供登录、预订、订单查询、个人资料维护等服务所必需的范围内处理你的信息。</text>
        <text class="paragraph">当前版本会保存登录令牌、用户基础资料与必要的订单信息，以便在移动端完成预订流程。</text>
        <text class="paragraph">如需清除本地数据，可在“设置”页面使用清除缓存或退出登录。</text>
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
