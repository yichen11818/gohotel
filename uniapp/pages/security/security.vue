<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">修改密码</text>
      <view class="placeholder-view"></view>
    </view>

    <scroll-view class="content" scroll-y>
      <!-- 温馨提示 -->
      <view class="safe-header flex-center">
        <image class="icon-safe" src="/static/icons/menu-safe.png" mode="aspectFit" />
        <text class="safe-text">请定期更新密码以保障账号安全</text>
      </view>

      <!-- 表单卡片 -->
      <view class="form-card premium-card">
        <view class="field-item">
          <text class="label">当前密码</text>
          <input v-model="form.oldPassword" class="input" password placeholder="请输入当前密码" placeholder-class="placeholder" />
        </view>

        <view class="field-item">
          <text class="label">新密码</text>
          <input v-model="form.newPassword" class="input" password placeholder="请输入至少 6 位新密码" placeholder-class="placeholder" />
        </view>

        <view class="field-item">
          <text class="label">确认新密码</text>
          <input v-model="confirmPassword" class="input" password placeholder="请再次输入新密码" placeholder-class="placeholder" />
        </view>
      </view>

      <view class="action-box">
        <button class="save-btn premium-button" :disabled="submitting" @click="handleSubmit">
          <text>{{ submitting ? '正在提交...' : '确认更新密码' }}</text>
        </button>
      </view>
    </scroll-view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import { user } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'

const statusBarHeight = ref(44)
const submitting = ref(false)
const confirmPassword = ref('')
const form = ref({
  oldPassword: '',
  newPassword: '',
})

onShow(() => {
  if (!uni.getStorageSync(TOKEN_KEY)) {
    uni.redirectTo({ url: '/pages/login/login' })
  }
})

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
})

const goBack = () => {
  uni.navigateBack()
}

const handleSubmit = async () => {
  if (!form.value.oldPassword || !form.value.newPassword) {
    uni.showToast({ title: '请完整填写信息', icon: 'none' })
    return
  }

  if (form.value.newPassword.length < 6) {
    uni.showToast({ title: '新密码至少 6 位', icon: 'none' })
    return
  }

  if (form.value.newPassword !== confirmPassword.value) {
    uni.showToast({ title: '两次密码不一致', icon: 'none' })
    return
  }

  try {
    submitting.value = true
    uni.showLoading({ title: '处理中' })
    await user.changePassword({
      oldPassword: form.value.oldPassword,
      newPassword: form.value.newPassword,
    })
    uni.hideLoading()
    uni.showToast({ title: '修改成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (_error) {
    uni.hideLoading()
  } finally {
    submitting.value = false
  }
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

.safe-header {
  padding: 60rpx 0 40rpx;
  flex-direction: column;
  .icon-safe {
    width: 100rpx;
    height: 100rpx;
    margin-bottom: 24rpx;
    opacity: 0.8;
  }
  .safe-text {
    font-size: 24rpx;
    color: $text-sub;
  }
}

.form-card {
  padding: 0 40rpx;

  .field-item {
    padding: 36rpx 0;
    border-bottom: 1rpx solid #f8f8f8;
    &:last-child { border-bottom: 0; }

    .label {
      font-size: 26rpx;
      color: $text-sub;
      margin-bottom: 20rpx;
      display: block;
    }

    .input {
      font-size: 30rpx;
      color: $text-main;
      font-weight: 500;
    }

    .placeholder {
      color: #ccc;
      font-weight: 400;
    }
  }
}

.action-box {
  margin-top: 60rpx;
  .save-btn {
    height: 100rpx;
    font-size: 32rpx;
    font-weight: 700;
  }
}
</style>
