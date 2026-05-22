<template>
  <view class="page">
    <!-- 自定义导航栏 -->
    <view class="status-bar" :style="{ height: statusBarHeight + 'px' }"></view>
    <view class="nav-header flex-between">
      <view class="back-btn flex-center" @click="goBack">
        <image class="icon-back" src="/static/icons/back-white.png" mode="aspectFit" style="filter: brightness(0);" />
      </view>
      <text class="page-title">个人资料</text>
      <view class="placeholder-view"></view>
    </view>

    <scroll-view class="content" scroll-y>
      <!-- 头像修改区 -->
      <view class="avatar-section flex-center">
        <view class="avatar-wrapper" @click="handleChooseAvatar">
          <image class="avatar" :src="form.avatar || '/static/icons/people.png'" mode="aspectFill" />
          <view class="camera-tag flex-center">
            <image src="/static/icons/menu-service.png" mode="aspectFit" style="width: 24rpx; height: 24rpx; filter: brightness(100);" />
          </view>
        </view>
        <text class="avatar-tip">点击修改头像</text>
      </view>

      <!-- 表单卡片 -->
      <view class="form-card premium-card">
        <view class="field-item">
          <text class="label">真实姓名</text>
          <input v-model="form.real_name" class="input" placeholder="请输入真实姓名" placeholder-class="placeholder" />
        </view>

        <view class="field-item">
          <text class="label">手机号</text>
          <input v-model="form.phone" class="input" type="number" maxlength="11" placeholder="请输入手机号" placeholder-class="placeholder" />
        </view>

        <view class="field-item">
          <text class="label">头像链接</text>
          <input v-model="form.avatar" class="input" placeholder="可粘贴图片链接" placeholder-class="placeholder" />
        </view>
      </view>

      <view class="action-box">
        <button class="save-btn premium-button" :disabled="saving" @click="handleSave">
          <text>{{ saving ? '正在保存...' : '确认保存' }}</text>
        </button>
      </view>

      <view class="safe-tips">
        <text>· 您的个人信息将严格保密，仅用于订单确认。</text>
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
    uni.showLoading({ title: '保存中' })
    await user.updateUserInfo(form.value)
    uni.hideLoading()
    uni.showToast({ title: '资料已更新', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (_error) {
    uni.hideLoading()
  } finally {
    saving.value = false
  }
}

const handleChooseAvatar = () => {
  uni.showToast({ title: '请在输入框粘贴链接或直接输入', icon: 'none' })
}

const goBack = () => {
  uni.navigateBack()
}

onLoad(() => {
  const sysInfo = uni.getSystemInfoSync()
  statusBarHeight.value = sysInfo.statusBarHeight || 44
})

onShow(() => {
  loadProfile()
})
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

/* 头像区 */
.avatar-section {
  padding: 60rpx 0;
  flex-direction: column;

  .avatar-wrapper {
    position: relative;
    width: 160rpx;
    height: 160rpx;
    margin-bottom: 24rpx;

    .avatar {
      width: 100%;
      height: 100%;
      border-radius: 50%;
      background: #eee;
      border: 4rpx solid #fff;
      box-shadow: 0 8rpx 20rpx rgba(0,0,0,0.1);
    }

    .camera-tag {
      position: absolute;
      right: 0;
      bottom: 0;
      width: 48rpx;
      height: 48rpx;
      background: $primary-color;
      border-radius: 50%;
      border: 4rpx solid #fff;
    }
  }

  .avatar-tip {
    font-size: 24rpx;
    color: $text-sub;
  }
}

/* 表单卡片 */
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

.safe-tips {
  margin-top: 40rpx;
  padding: 0 20rpx;
  text-align: center;
  font-size: 22rpx;
  color: $text-sub;
}
</style>
