<template>
  <div class="float-control">
    <!-- ===== 页面头部 ===== -->
    <div class="page-header">
      <h2 class="page-title">
        浮窗控制
      </h2>
    </div>

    <!-- ===== 控制项列表 ===== -->
    <div class="control-list">
      <FloatControlItem
        v-for="control in controls"
        :key="control.name"
        :name="control.name"
        :icon="control.icon"
        :running="floatRunning"
        @toggle="toggleFloat(control.windowName)"
      />
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

// 浮窗控制项列表
const controls = [
  {
    name: '单词显示',
    windowName: 'LancioFloat',
    icon: 'M3 5h8v6H3V5zm10 0h8v6h-8V5zM3 13h8v6H3v-6zm10 0h8v6h-8v-6z',
  },
]

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
  if (controls.length > 0) {
    for (let i = 0; i < controls.length; i++) {
      loadStatus(controls[i].windowName)
    }
  } 
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
  gap: 16px;
}

/* ===== 页面头部 ===== */
.page-header {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  flex-shrink: 0;
}

.page-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 28px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  margin-top: 5px;
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