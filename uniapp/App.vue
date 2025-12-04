<script setup>
import { onLaunch, onShow, onHide } from '@dcloudio/uni-app'

onLaunch(() => {
  console.log('App Launch')
})

onShow(() => {
  console.log('App Show')
})

onHide(() => {
  console.log('App Hide')
})
</script>

<style lang="scss">
@import '@/uni.scss';
@import '@/uni_modules/tuniaoui-vue3/index.css';

page {
  background-color: #f8f8f8;
}
</style>

