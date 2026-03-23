<template>
  <view class="page">
    <view class="content">
      <view class="section">
        <text class="section-title">修改密码</text>
        <view class="field">
          <text class="field-label">当前密码</text>
          <input v-model="form.oldPassword" class="field-input" password placeholder="请输入当前密码" />
        </view>
        <view class="field">
          <text class="field-label">新密码</text>
          <input v-model="form.newPassword" class="field-input" password placeholder="请输入至少 6 位新密码" />
        </view>
        <view class="field">
          <text class="field-label">确认新密码</text>
          <input v-model="confirmPassword" class="field-input" password placeholder="请再次输入新密码" />
        </view>
        <button class="primary-btn" :disabled="submitting" @click="handleSubmit">
          {{ submitting ? '提交中...' : '更新密码' }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { user } from '@/api/index.js'
import { TOKEN_KEY } from '@/config/api.config.js'

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

const handleSubmit = async () => {
  if (!form.value.oldPassword || !form.value.newPassword) {
    uni.showToast({ title: '请完整填写密码信息', icon: 'none' })
    return
  }

  if (form.value.newPassword.length < 6) {
    uni.showToast({ title: '新密码至少 6 位', icon: 'none' })
    return
  }

  if (form.value.newPassword !== confirmPassword.value) {
    uni.showToast({ title: '两次输入的新密码不一致', icon: 'none' })
    return
  }

  try {
    submitting.value = true
    uni.showLoading({ title: '提交中...' })
    await user.changePassword({
      oldPassword: form.value.oldPassword,
      newPassword: form.value.newPassword,
    })
    uni.hideLoading()
    uni.showToast({ title: '密码修改成功', icon: 'success' })
    form.value.oldPassword = ''
    form.value.newPassword = ''
    confirmPassword.value = ''
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
  background: #f6f7fb;
}

.content {
  padding: 24rpx;
}

.section {
  padding: 28rpx;
  background: #fff;
  border-radius: 24rpx;
  box-shadow: 0 12rpx 30rpx rgba(15, 23, 42, 0.05);
}

.section-title,
.field-label {
  display: block;
  font-size: 28rpx;
  color: #111827;
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
}

.field {
  margin-top: 24rpx;
}

.field-input {
  width: 100%;
  margin-top: 12rpx;
  padding: 22rpx 24rpx;
  box-sizing: border-box;
  background: #f9fafb;
  border: 1px solid #eef2f7;
  border-radius: 18rpx;
}

.primary-btn {
  margin-top: 32rpx;
  background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
  color: #fff;
  border-radius: 999rpx;
}
</style>
