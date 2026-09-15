<template>
  <div class="api-config-form">
    <div class="section-title">
      <h3>百度翻译 API</h3>
      <p>配置百度翻译开放平台的 API 密钥</p>
    </div>

    <div class="form-group">
      <label class="form-label" for="appid">APP ID</label>
      <input
        id="appid"
        v-model="form.appid"
        type="text"
        class="form-input"
        placeholder="请输入百度翻译 APP ID"
        spellcheck="false"
      />
      <span class="form-hint">百度翻译开放平台注册后获取的应用 ID</span>
    </div>

    <div class="form-group">
      <label class="form-label" for="apikey">API Key（密钥）</label>
      <div class="input-wrapper">
        <input
          id="apikey"
          v-model="form.key"
          :type="showKey ? 'text' : 'password'"
          class="form-input"
          placeholder="请输入百度翻译 API Key"
          spellcheck="false"
        />
        <button class="toggle-key-btn" title="切换密钥可见性" @click="showKey = !showKey">
          <svg v-if="!showKey" width="16" height="16" viewBox="0 0 16 16">
            <path d="M8 3C4.5 3 1.5 5.5 0.5 8c1 2.5 4 5 7.5 5s6.5-2.5 7.5-5c-1-2.5-4-5-7.5-5z" stroke="currentColor" stroke-width="1.2" fill="none"/>
            <circle cx="8" cy="8" r="2" stroke="currentColor" stroke-width="1.2" fill="none"/>
          </svg>
          <svg v-else width="16" height="16" viewBox="0 0 16 16">
            <path d="M8 3C4.5 3 1.5 5.5 0.5 8c1 2.5 4 5 7.5 5s6.5-2.5 7.5-5c-1-2.5-4-5-7.5-5z" stroke="currentColor" stroke-width="1.2" fill="none"/>
            <line x1="2" y1="2" x2="14" y2="14" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
          </svg>
        </button>
      </div>
      <span class="form-hint">百度翻译开放平台注册后获取的密钥，请妥善保管</span>
    </div>

    <div class="section-title">
      <h3>复习计划</h3>
      <p>每次复习从数据库中抽取的单词数量</p>
    </div>

    <div class="form-group">
      <label class="form-label" for="reviewCount">每次复习单词数</label>
      <div class="number-input-wrap">
        <input
          id="reviewCount"
          v-model.number="form.reviewCount"
          type="number"
          class="form-input no-spinner"
          placeholder="10"
          min="1"
          max="100"
        />
        <div class="spinner-btns">
          <button class="spinner-btn spinner-up" @click="incrementReviewCount" :disabled="form.reviewCount >= 100" title="增加">
            <svg width="10" height="6" viewBox="0 0 10 6">
              <path d="M1 5l4-4 4 4" stroke="currentColor" stroke-width="1.3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
          <button class="spinner-btn spinner-down" @click="decrementReviewCount" :disabled="form.reviewCount <= 1" title="减少">
            <svg width="10" height="6" viewBox="0 0 10 6">
              <path d="M1 1l4 4 4-4" stroke="currentColor" stroke-width="1.3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </button>
        </div>
      </div>
      <span class="form-hint">优先抽取复习次数最少的单词，不足时自动补足</span>
    </div>

    <div class="form-actions">
      <button
        :class="['btn-test', { 'btn-test-ok': testOk, 'btn-test-fail': testResult && !testOk }]"
        :disabled="formLoading"
        @click="testConnection"
      >
        <svg v-if="!testResult" width="14" height="14" viewBox="0 0 14 14">
          <path d="M2 2h4v3M12 12H8V9" stroke="currentColor" stroke-width="1.3" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
          <path d="M11 5a4 4 0 00-6.93-2M3 9a4 4 0 006.93 2" stroke="currentColor" stroke-width="1.3" fill="none" stroke-linecap="round"/>
        </svg>
        <svg v-else-if="testOk" width="14" height="14" viewBox="0 0 14 14">
          <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
          <path d="M4.5 7l2 2 3-4" stroke="currentColor" stroke-width="1.4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <svg v-else width="14" height="14" viewBox="0 0 14 14">
          <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
          <path d="M5 5l4 4M9 5l-4 4" stroke="currentColor" stroke-width="1.4" fill="none" stroke-linecap="round"/>
        </svg>
        {{ testing ? '测试中...' : testResult || '测试连接' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { GetAPIConfig, SaveAPIConfig, TestAPIConnection, GetReviewWordsPerSession, SetReviewWordsPerSession } from '../../../wailsjs/go/handler/ConfigHandler'

const form = reactive({
  appid: '',
  key: '',
  reviewCount: 10,
})

const showKey = ref(false)
const loading = ref(false)
const testing = ref(false)
const formLoading = ref(false)
const testResult = ref('')
const testOk = ref(false)

onMounted(async () => {
  try {
    const config = await GetAPIConfig()
    form.appid = config.appid || ''
    form.key = config.key || ''
    form.reviewCount = await GetReviewWordsPerSession()
  } catch (err) {
    console.error('加载配置失败:', err)
  }
})

async function saveConfig() {
  if (!form.appid.trim() || !form.key.trim()) {
    return { success: false, message: 'APP ID 和 API Key 不能为空' }
  }
  loading.value = true
  try {
    await SaveAPIConfig(form.appid.trim(), form.key.trim())
    await SetReviewWordsPerSession(form.reviewCount || 10)
    return { success: true, message: '配置保存成功' }
  } catch (err) {
    console.error('保存配置失败:', err)
    return { success: false, message: '保存配置失败: ' + (err.message || err) }
  } finally {
    loading.value = false
  }
}

async function testConnection() {
  if (testResult.value) {
    testResult.value = ''
    testOk.value = false
    return
  }
  if (!form.appid.trim() || !form.key.trim()) {
    testResult.value = '请先填写 APP ID 和 API Key'
    return
  }
  testing.value = true
  formLoading.value = true
  try {
    await TestAPIConnection(form.appid.trim(), form.key.trim())
    testOk.value = true
    testResult.value = '连接成功'
  } catch (err) {
    testOk.value = false
    testResult.value = '连接失败: ' + (err.message || err)
  } finally {
    testing.value = false
    formLoading.value = false
  }
}

defineExpose({ saveConfig, loading })

// 自定义数字输入框的 +/- 按钮逻辑
function incrementReviewCount() {
  if (form.reviewCount < 100) form.reviewCount++
}
function decrementReviewCount() {
  if (form.reviewCount > 1) form.reviewCount--
}
</script>

<style scoped>
.api-config-form {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.section-title {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.section-title h3 {
  font-size: 15px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.8);
  margin: 0;
}

.section-title p {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.35);
  margin: 0;
  line-height: 1.4;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.7);
}

.form-input {
  width: 100%;
  height: 38px;
  padding: 0 12px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.85);
  outline: none;
  transition: border-color 0.2s, background 0.2s;
}

.form-input::placeholder {
  color: rgba(255, 255, 255, 0.2);
}

.form-input:focus {
  border-color: rgba(96, 165, 250, 0.5);
  background: rgba(255, 255, 255, 0.06);
}

.form-hint {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.25);
  line-height: 1.4;
}

.number-input-wrap {
  position: relative;
  display: flex;
  align-items: center;
}

.no-spinner {
  -moz-appearance: textfield;
}

.no-spinner::-webkit-outer-spin-button,
.no-spinner::-webkit-inner-spin-button {
  -webkit-appearance: none;
  margin: 0;
}

.spinner-btns {
  position: absolute;
  right: 2px;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.spinner-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 17px;
  background: transparent;
  border: none;
  border-radius: 3px;
  color: rgba(255, 255, 255, 0.3);
  cursor: pointer;
  transition: all 0.15s;
  outline: none;
}

.spinner-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.06);
  color: rgba(255, 255, 255, 0.6);
}

.spinner-btn:active:not(:disabled) {
  background: rgba(79, 195, 247, 0.15);
  color: rgba(79, 195, 247, 0.7);
}

.spinner-btn:disabled {
  opacity: 0.2;
  cursor: not-allowed;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-wrapper .form-input {
  padding-right: 40px;
}

.toggle-key-btn {
  position: absolute;
  right: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: none;
  border: none;
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.3);
  cursor: pointer;
  transition: color 0.15s, background 0.15s;
}

.toggle-key-btn:hover {
  color: rgba(255, 255, 255, 0.6);
  background: rgba(255, 255, 255, 0.06);
}

.form-actions {
  display: flex;
  gap: 10px;
}

.btn-test {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 34px;
  padding: 0 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.65);
  cursor: pointer;
  transition: all 0.15s;
  outline: none;
}

.btn-test:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-test:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.85);
}

.btn-test-ok {
  border-color: rgba(34, 197, 94, 0.5) !important;
  color: rgba(74, 222, 128, 0.85) !important;
}

.btn-test-fail {
  border-color: rgba(239, 68, 68, 0.5) !important;
  color: rgba(248, 113, 113, 0.85) !important;
}
</style>