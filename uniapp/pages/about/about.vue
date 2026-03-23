<template>
  <scroll-view class="page" scroll-y>
    <view class="content">
      <view class="section">
        <text class="title">关于七天酒店</text>
        <text class="paragraph">七天酒店面向住客提供房型浏览、在线预订、订单查询、个人资料维护等基础能力。</text>
        <text class="paragraph">当前移动端以真实后端能力为准，所有页面与按钮均围绕可实际使用的预订流程收口。</text>
      </view>
      <view class="section">
        <text class="title">联系我们</text>
        <text class="paragraph">如需帮助，请在“酒店服务”页面联系前台，或前往首页查看酒店地址与服务时间。</text>
      </view>
    </view>
  </scroll-view>
</template>

<script setup>
</script>

<style scoped lang="scss">
.page { min-height: 100vh; background: #f6f7fb; }
.content { padding: 24rpx; }
.section { padding: 28rpx; margin-bottom: 24rpx; background: #fff; border-radius: 24rpx; box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05); }
.title { display: block; font-size: 32rpx; font-weight: 700; color: #111827; }
.paragraph { display: block; margin-top: 16rpx; font-size: 26rpx; color: #4b5563; line-height: 1.8; }
</style>
