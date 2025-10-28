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

<style>
/* 全局样式 */
page {
  background-color: #F8F8F8;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB',
    'Microsoft YaHei', 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

/* 顺丰红色主题 */
:root {
  --sf-primary: #D81E06;
  --sf-primary-light: #FF3B30;
  --sf-primary-dark: #C01A05;
  --sf-text-primary: #333333;
  --sf-text-secondary: #666666;
  --sf-text-tertiary: #999999;
  --sf-border: #E5E5E5;
  --sf-bg: #F8F8F8;
  --sf-white: #FFFFFF;
}
</style>
