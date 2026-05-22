<template>
  <view class="calendar-mask" :class="{ show: show }" @click="close">
    <view class="calendar-wrapper" @click.stop>
      <view class="header">
        <text class="title">选择日期</text>
        <view class="close-btn flex-center" @click="close">
          <text class="icon-close">×</text>
        </view>
      </view>

      <view class="week-bar">
        <text class="week-item weekend">日</text>
        <text class="week-item">一</text>
        <text class="week-item">二</text>
        <text class="week-item">三</text>
        <text class="week-item">四</text>
        <text class="week-item">五</text>
        <text class="week-item weekend">六</text>
      </view>

      <scroll-view scroll-y class="month-scroll">
        <view class="month-block" v-for="month in months" :key="month.title">
          <view class="month-title">{{ month.title }}</view>
          <view class="day-grid">
            <view
              class="day-item"
              v-for="(day, index) in month.days"
              :key="index"
              :class="[
                day.empty ? 'empty' : '',
                day.disabled ? 'disabled' : '',
                isCheckIn(day) ? 'check-in active' : '',
                isCheckOut(day) ? 'check-out active' : '',
                isInRange(day) ? 'in-range' : ''
              ]"
              @click="onDayClick(day)"
            >
              <template v-if="!day.empty">
                <text class="day-num">{{ day.day }}</text>
                <text class="day-text">{{ getDayText(day) }}</text>
              </template>
            </view>
          </view>
        </view>
      </scroll-view>

      <view class="footer safe-area-inset-bottom">
        <button class="confirm-btn premium-button" :disabled="!tempCheckOut" @click="confirm">
          {{ tempCheckOut ? '完成' : '请选择离店日期' }}
        </button>
      </view>
    </view>
  </view>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  show: { type: Boolean, default: false },
  checkIn: { type: String, default: '' },
  checkOut: { type: String, default: '' }
})

const emit = defineEmits(['update:show', 'confirm'])

const months = ref([])
const tempCheckIn = ref('')
const tempCheckOut = ref('')

const formatDateStr = (date) => {
  const y = date.getFullYear()
  const m = String(date.getMonth() + 1).padStart(2, '0')
  const d = String(date.getDate()).padStart(2, '0')
  return `${y}-${m}-${d}`
}

const initCalendar = () => {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  const tempMonths = []
  for (let i = 0; i < 6; i++) {
    const d = new Date(today.getFullYear(), today.getMonth() + i, 1)
    const year = d.getFullYear()
    const month = d.getMonth()
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    const firstDay = d.getDay()

    const days = []
    for (let j = 0; j < firstDay; j++) {
      days.push({ empty: true })
    }
    for (let j = 1; j <= daysInMonth; j++) {
      const dateObj = new Date(year, month, j)
      const dateStr = formatDateStr(dateObj)
      days.push({
        empty: false,
        day: j,
        dateStr: dateStr,
        disabled: dateObj.getTime() < today.getTime(),
        isToday: dateObj.getTime() === today.getTime()
      })
    }
    tempMonths.push({
      title: `${year}年${month + 1}月`,
      days
    })
  }
  months.value = tempMonths
}

watch(() => props.show, (newVal) => {
  if (newVal) {
    if (!months.value.length) initCalendar()
    tempCheckIn.value = props.checkIn
    tempCheckOut.value = props.checkOut
  }
})

const isCheckIn = (day) => !day.empty && day.dateStr === tempCheckIn.value
const isCheckOut = (day) => !day.empty && day.dateStr === tempCheckOut.value
const isInRange = (day) => {
  if (day.empty || !tempCheckIn.value || !tempCheckOut.value) return false
  return day.dateStr > tempCheckIn.value && day.dateStr < tempCheckOut.value
}

const getDayText = (day) => {
  if (isCheckIn(day)) return '入住'
  if (isCheckOut(day)) return '离店'
  if (day.isToday) return '今天'
  return ''
}

const onDayClick = (day) => {
  if (day.empty || day.disabled) return

  if (!tempCheckIn.value || (tempCheckIn.value && tempCheckOut.value)) {
    tempCheckIn.value = day.dateStr
    tempCheckOut.value = ''
  } else {
    if (day.dateStr < tempCheckIn.value) {
      tempCheckIn.value = day.dateStr
    } else if (day.dateStr > tempCheckIn.value) {
      tempCheckOut.value = day.dateStr
    } else {
      // Clicked same day again
      tempCheckIn.value = day.dateStr
      tempCheckOut.value = ''
    }
  }
}

const close = () => {
  emit('update:show', false)
}

const confirm = () => {
  if (tempCheckIn.value && tempCheckOut.value) {
    emit('confirm', { checkIn: tempCheckIn.value, checkOut: tempCheckOut.value })
    close()
  }
}
</script>

<style scoped lang="scss">
.calendar-mask {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 999;
  visibility: hidden;
  opacity: 0;
  transition: all 0.3s;

  &.show {
    visibility: visible;
    opacity: 1;
    .calendar-wrapper {
      transform: translateY(0);
    }
  }
}

.calendar-wrapper {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  background: #fff;
  border-radius: 32rpx 32rpx 0 0;
  transform: translateY(100%);
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex;
  flex-direction: column;
  height: 80vh;
}

.header {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100rpx;
  position: relative;
  border-bottom: 1rpx solid #f8f8f8;

  .title {
    font-size: 32rpx;
    font-weight: 700;
    color: #333;
  }

  .close-btn {
    position: absolute;
    right: 0;
    top: 0;
    width: 100rpx;
    height: 100rpx;
    .icon-close { font-size: 44rpx; color: #999; }
  }
}

.week-bar {
  display: flex;
  height: 80rpx;
  align-items: center;
  background: #fafafa;

  .week-item {
    flex: 1;
    text-align: center;
    font-size: 24rpx;
    color: #666;
    &.weekend { color: #E64340; }
  }
}

.month-scroll {
  flex: 1;
  height: 0;
}

.month-block {
  padding-bottom: 40rpx;

  .month-title {
    text-align: center;
    padding: 30rpx 0;
    font-size: 30rpx;
    font-weight: 600;
    color: #333;
  }

  .day-grid {
    display: flex;
    flex-wrap: wrap;

    .day-item {
      width: 14.285%;
      height: 100rpx;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      position: relative;
      margin-bottom: 10rpx;

      .day-num { font-size: 32rpx; color: #333; font-weight: 500; z-index: 2; }
      .day-text { font-size: 20rpx; color: #c9a977; margin-top: 4rpx; z-index: 2; height: 28rpx; }

      &.disabled {
        .day-num { color: #ccc; }
        .day-text { color: transparent; }
      }

      &.in-range {
        background: rgba(201, 169, 119, 0.1);
      }

      &.active {
        background: #c9a977;
        border-radius: 8rpx;
        .day-num, .day-text { color: #fff; }

        &.check-in {
          border-top-right-radius: 0;
          border-bottom-right-radius: 0;
        }
        &.check-out {
          border-top-left-radius: 0;
          border-bottom-left-radius: 0;
        }
      }
    }
  }
}

.footer {
  padding: 24rpx 40rpx;
  background: #fff;
  box-shadow: 0 -4rpx 16rpx rgba(0,0,0,0.05);

  .confirm-btn {
    height: 90rpx;
    border-radius: 100rpx;
    font-size: 32rpx;
    margin: 0;
    background: linear-gradient(135deg, #c9a977 0%, #ad8551 100%);
    color: #fff;

    &[disabled] {
      background: #ccc;
      box-shadow: none;
    }
  }
}
</style>
