<template>
  <div class="float-control">
    <!-- ===== 页面头部 ===== -->
    <div class="page-header">
      <h2 class="page-title">
       🪟 浮窗控制
        <span class="title-meta">控制桌面浮窗的显示与隐藏</span>
      </h2>
    </div>

    <!-- ===== 控制项列表 ===== -->
    <div class="control-list">
      <FloatControlItem
        name="单词显示"
        :running="floatRunning"
        @toggle="toggleFloat('LancioFloat')"
      >
        <template #icon>
          <svg width="22" height="22" viewBox="0 0 22 22">
            <rect x="2" y="4" width="18" height="14" rx="2" stroke="currentColor" stroke-width="1.3" fill="none"/>
            <line x1="2" y1="8" x2="20" y2="8" stroke="currentColor" stroke-width="1.3"/>
            <circle cx="5.5" cy="6" r="0.8" fill="#ef4444"/>
            <circle cx="8" cy="6" r="0.8" fill="#f59e0b"/>
            <circle cx="10.5" cy="6" r="0.8" fill="#22c55e"/>
          </svg>
        </template>
      </FloatControlItem>
    </div>

    <!-- ===== 状态提示 ===== -->
    <ToastMessage :message="message" :type="messageType" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { IsFloatWindowRunning, ToggleFloatWindow } from '../../wailsjs/go/handler/WindowHandler.js'
import FloatControlItem from '@/components/floatwinmanager/FloatControlItem.vue'
import ToastMessage from '@/components/floatwinmanager/ToastMessage.vue'

const floatRunning = ref(false)
const message = ref('')
const messageType = ref('info')

async function loadStatus(windowName) {
  try {
    floatRunning.value = await IsFloatWindowRunning(windowName)
  } catch (err) {
    console.error('获取浮窗状态失败:', err)
  }
}

async function toggleFloat(windowName) {
  try {
    const running = await ToggleFloatWindow(windowName)
    floatRunning.value = running
    message.value = running ? '浮窗已开启' : '浮窗已关闭'
    messageType.value = running ? 'success' : 'info'
    setTimeout(() => { message.value = '' }, 2000)
  } catch (err) {
    console.error('切换浮窗状态失败:', err)
    message.value = '操作失败'
    messageType.value = 'error'
    setTimeout(() => { message.value = '' }, 2000)
  }
}

onMounted(() => {
  loadStatus('LancioFloat')
})
</script>

<style scoped>
/* ===== 容器 ===== */
.float-control {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px 24px;
  overflow: hidden;
  gap: 28px;
}

/* ===== 页面头部 ===== */
.page-header {
  flex-shrink: 0;
}

.page-title {
  display: flex;
  align-items: baseline;
  gap: 10px;
  font-size: 22px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
}

.title-meta {
  font-size: 12px;
  font-weight: 400;
  color: rgba(255, 255, 255, 0.4);
}

.title-icon-svg {
  flex-shrink: 0;
  color: rgba(255, 255, 255, 0.55);
}

/* ===== 控制项列表 ===== */
.control-list {
  flex: 1;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  align-content: start;
  align-items: start;
}
</style>