<template>
  <div class="config-pannel">
    <div class="page-header">
      <h2 class="page-title">
        ⚙️ 设置
      </h2>
    </div>

    <div class="page-body">
      <ApiConfigForm ref="formRef" />
    </div>

    <div class="page-footer">
      <div v-if="statusMessage" :class="['status-bar', statusType]">
        <svg v-if="statusType === 'success'" width="14" height="14" viewBox="0 0 14 14">
          <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
          <path d="M4.5 7l2 2 3-4" stroke="currentColor" stroke-width="1.4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else width="14" height="14" viewBox="0 0 14 14">
          <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
          <path d="M5 5l4 4M9 5l-4 4" stroke="currentColor" stroke-width="1.4" fill="none" stroke-linecap="round"/>
        </svg>
        <span>{{ statusMessage }}</span>
      </div>
      <div class="footer-actions">
        <button class="btn btn-primary" :disabled="formLoading" @click="handleSave">
          <svg width="14" height="14" viewBox="0 0 14 14">
            <path d="M3 2h6l3 3v7a1 1 0 01-1 1H3a1 1 0 01-1-1V3a1 1 0 011-1zm6 0v3h3" stroke="currentColor" stroke-width="1.2" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
          {{ saving ? '保存中...' : '保存配置' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ApiConfigForm from '@/components/configpanel/ApiConfigForm.vue'

const formRef = ref(null)
const saving = ref(false)
const formLoading = ref(false)
const statusMessage = ref('')
const statusType = ref('success')

function showStatus(msg, type) {
  statusMessage.value = msg
  statusType.value = type
  setTimeout(() => {
    statusMessage.value = ''
  }, 4000)
}

async function handleSave() {
  if (!formRef.value) return
  saving.value = true
  formLoading.value = true
  const result = await formRef.value.saveConfig()
  showStatus(result.message, result.success ? 'success' : 'error')
  saving.value = false
  formLoading.value = false
}
</script>

<style scoped>
.config-pannel {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px 24px;
  overflow: hidden;
  gap: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.page-title {
  width: 100%;
  display: flex;
  align-items: baseline;
  gap: 10px;
  font-size: 22px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  margin: 0;
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  padding-bottom: 8px;
}

.page-body {
  flex: 1;
  overflow-y: auto;
}

.page-body::-webkit-scrollbar {
  width: 6px;
}

.page-body::-webkit-scrollbar-track {
  background: transparent;
}

.page-body::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 3px;
}

.page-body::-webkit-scrollbar-thumb:hover {
  background: rgba(79, 195, 247, 0.35);
}

.page-footer {
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px 0 4px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.status-bar {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 8px;
  font-size: 12px;
  line-height: 1.4;
}

.status-bar.success {
  background: rgba(34, 197, 94, 0.08);
  color: rgba(74, 222, 128, 0.85);
  border: 1px solid rgba(34, 197, 94, 0.12);
}

.status-bar.error {
  background: rgba(239, 68, 68, 0.08);
  color: rgba(248, 113, 113, 0.85);
  border: 1px solid rgba(239, 68, 68, 0.12);
}

.footer-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 34px;
  padding: 0 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  outline: none;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.04);
  color: rgba(255, 255, 255, 0.65);
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.85);
}

.btn-primary {
  background: rgba(59, 130, 246, 0.2);
  border-color: rgba(59, 130, 246, 0.25);
  color: rgba(96, 165, 250, 0.9);
}

.btn-primary:hover:not(:disabled) {
  background: rgba(59, 130, 246, 0.3);
  border-color: rgba(59, 130, 246, 0.4);
  color: rgba(147, 197, 253, 0.95);
}
</style>