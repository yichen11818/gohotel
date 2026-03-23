<template>
  <view class="page">
    <view class="content">
      <view class="section">
        <text class="section-title">个人资料</text>
        <view class="field">
          <text class="field-label">真实姓名</text>
          <input v-model="form.real_name" class="field-input" placeholder="请输入真实姓名" />
        </view>
        <view class="field">
          <text class="field-label">手机号</text>
          <input v-model="form.phone" class="field-input" type="number" maxlength="11" placeholder="请输入手机号" />
        </view>
        <view class="field">
          <text class="field-label">头像 URL</text>
          <input v-model="form.avatar" class="field-input" placeholder="请输入头像链接（选填）" />
        </view>
        <button class="primary-btn" :disabled="saving" @click="handleSave">
          {{ saving ? '保存中...' : '保存资料' }}
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

const saving = ref(false)
const form = ref({
  real_name: '',
  phone: '',
  avatar: '',
})

const loadProfile = async () => {
  if (!uni.getStorageSync(TOKEN_KEY)) {
    uni.redirectTo({ url: '/pages/login/login' })
    return
  }

  try {
    const profile = await user.getUserInfo()
    form.value = {
      real_name: profile.real_name || profile.nickname || '',
      phone: profile.phone || '',
      avatar: profile.avatar || '',
    }
  } catch (error) {
    console.error('load profile failed:', error)
  }
}

const handleSave = async () => {
  if (form.value.phone && !/^1\d{10}$/.test(form.value.phone)) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }

  try {
    saving.value = true
    uni.showLoading({ title: '保存中...' })
    await user.updateUserInfo(form.value)
    uni.hideLoading()
    uni.showToast({ title: '保存成功', icon: 'success' })
  } catch (_error) {
    uni.hideLoading()
  } finally {
    saving.value = false
  }
}

onShow(() => {
  loadProfile()
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
