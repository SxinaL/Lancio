<template>
  <div class="review-card-wrapper">
    <template v-if="!currentWord">
      <div class="empty-state">
        <div class="empty-icon">{{ totalInRound > 0 && reviewedCount >= totalInRound ? '🎉' : '⭐' }}</div>
        <div class="empty-text">
          {{ totalInRound > 0 && reviewedCount >= totalInRound ? '本轮复习已完成！' : '点击下方按钮开启新一轮复习' }}
        </div>
      </div>
    </template>
    <template v-else>
    <div class="progress-bar">
      <span class="progress-text">{{ queueDisplayIndex }} / {{ totalInRound }}</span>
      <div class="progress-track">
        <div class="progress-fill" :style="{ width: progressPercent + '%' }"></div>
      </div>
    </div>

    <div class="card-container" @click="flip">
      <div class="card-inner" :class="{ flipped: isFlipped }">
        <div class="card-front">
          <div class="card-content">
            <div class="word-main">{{ currentWord?.word || '—' }}</div>
            <div v-if="currentWord?.phonetic" class="word-phonetic">
              <span class="audio-btn" @click.stop="playAudio" title="播放发音">
                <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" width="14" height="14">
                  <path fill="currentColor" d="M5 15V9h4l5-5v16l-5-5zm11 1V7.95q1.125.525 1.813 1.625T18.5 12t-.687 2.4T16 16"/>
                </svg>
              </span>
              <span>{{ currentWord.phonetic }}</span>
            </div>
            <audio ref="audioRef" :src="audioSrc" style="display:none"></audio>
          </div>
          <div class="flip-hint">点击翻转查看释义</div>
        </div>
        <div class="card-back">
          <div class="card-content">
            <div class="word-main-sm">{{ currentWord?.word || '—' }}</div>
            <div class="translation-main">{{ currentWord?.translation || '暂无释义' }}</div>
            <div v-if="currentWord?.example_sentence" class="example-sentence">
              {{ currentWord.example_sentence }}
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- :class="{ visible: isFlipped } -->
    <div class="review-actions" >
      <button class="action-btn forget" :disabled="loading" @click.stop="markForget">
        <svg width="16" height="16" viewBox="0 0 16 16">
          <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
          <path d="M5 5l6 6M11 5l-6 6" stroke="currentColor" stroke-width="1.4" fill="none" stroke-linecap="round"/>
        </svg>
        不认识
      </button>
      <button class="action-btn know" :disabled="loading" @click.stop="markKnow">
        <svg width="16" height="16" viewBox="0 0 16 16">
          <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.2" fill="none"/>
          <path d="M5 8l2 2 4-4" stroke="currentColor" stroke-width="1.4" fill="none" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        认识
      </button>
    </div>
    </template>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { IncrementReviewCount, SetReviewStatus } from '../../../wailsjs/go/handler/VocabularyHandler'

const props = defineProps({
  wordList: { type: Array, default: () => [] }, // 本轮所有单词（pending + reviewed）
})

const emit = defineEmits(['done']) // 全部复习完成时触发

// ── 复习队列（核心数据结构）──
// 用队列替代线性索引：不认识→回到队尾循环，直到全部认识
const wordQueue = ref([])        // 当前待复习队列（仅 pending）
const totalInRound = ref(0)     // 本轮单词总数
const reviewedCount = ref(0)    // 已认识数量
const currentWord = ref(null)   // 当前展示的单词
const isFlipped = ref(false)    // 卡片是否翻转
const loading = ref(false)

// 进度百分比：已认识 / 总数
const progressPercent = computed(() => {
  if (totalInRound.value === 0) return 0
  return Math.round((reviewedCount.value / totalInRound.value) * 100)
})
// 显示序号：已认识 + 当前正在看
const queueDisplayIndex = computed(() => reviewedCount.value )

const audioRef = ref(null)
const audioSrc = computed(() => {
  if (!currentWord.value?.word) return ''
  return 'https://dict.youdao.com/dictvoice?audio=' + currentWord.value.word.replaceAll(' ', '%20')
})

function playAudio() {
  if (audioRef.value) {
    audioRef.value.currentTime = 0
    audioRef.value.play()
  }
}

// 监听外部传入的单词列表，拆分为 pending 队列和 reviewed 计数
watch(() => props.wordList, (list) => {
  if (list.length > 0) {
    const pending = list.filter(w => w.review_status === 'pending')
    const reviewed = list.filter(w => w.review_status === 'reviewed')
    wordQueue.value = [...pending]
    totalInRound.value = list.length
    reviewedCount.value = reviewed.length
    if (wordQueue.value.length > 0) {
      currentWord.value = wordQueue.value.shift()
    } else {
      currentWord.value = null
    }
    isFlipped.value = false
  } else {
    wordQueue.value = []
    totalInRound.value = 0
    reviewedCount.value = 0
    currentWord.value = null
  }
}, { immediate: true })

// 从队列头部取下一个单词，队列为空则本轮完成
function loadNextWord() {
  if (wordQueue.value.length > 0) {
    currentWord.value = wordQueue.value.shift()
    isFlipped.value = false
  } else {
    currentWord.value = null
    // 本轮单词已全部认识，触发完成事件
    emit('done')
  }
}

function flip() {
  if (!currentWord.value || loading.value) return
  isFlipped.value = !isFlipped.value
}

// 认识：标记 reviewed + 复习次数+1 → 出队
async function markKnow() {
  if (!currentWord.value) return
  const wordId = currentWord.value.id
  loading.value = true
  try {
    await SetReviewStatus(wordId, 'reviewed') // 状态改为已复习
    await IncrementReviewCount(wordId)        // 累计复习次数 +1
    reviewedCount.value++
    loadNextWord()
  } catch (err) {
    console.error('标记认识失败:', err)
  } finally {
    loading.value = false
  }
}

// 不认识：复习次数+1，状态保持 pending → 回到队尾继续循环
async function markForget() {
  if (!currentWord.value) return
  const word = currentWord.value
  loading.value = true
  try {
    // 不认识也增加复习次数（已过目），但状态保持 pending，回到队尾
    await IncrementReviewCount(word.id)
    wordQueue.value.push(word)
    loadNextWord()
  } catch (err) {
    console.error('标记不认识失败:', err)
  } finally {
    loading.value = false
  }
}

defineExpose({ currentWord, startRound: () => {} })
</script>

<style scoped>
.review-card-wrapper {
  width: 100%;
  max-width: 420px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

.progress-bar {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 10px;
}

.progress-text {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.45);
  white-space: nowrap;
  min-width: 50px;
}

.progress-track {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 2px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: rgba(96, 165, 250, 0.6);
  border-radius: 2px;
  transition: width 0.4s ease;
}

.card-container {
  width: 100%;
  aspect-ratio: 1.6;
  perspective: 800px;
  cursor: pointer;
}

.card-inner {
  width: 100%;
  height: 100%;
  position: relative;
  transition: transform 0.5s ease;
  transform-style: preserve-3d;
}

.card-inner.flipped {
  transform: rotateY(180deg);
}

.card-front,
.card-back {
  position: absolute;
  inset: 0;
  backface-visibility: hidden;
  border-radius: 14px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.card-front {
  background: rgba(255, 255, 255, 0.03);
}

.card-back {
  background: rgba(59, 130, 246, 0.06);
  transform: rotateY(180deg);
}

.card-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  padding: 24px;
}

.word-main {
  font-size: 36px;
  font-weight: 700;
  color: rgba(255, 255, 255, 0.9);
  letter-spacing: 2px;
}

.word-main-sm {
  font-size: 30px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
}

.word-phonetic {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  display: flex;
  align-items: center;
  gap: 4px;
}

.audio-btn {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  color: rgba(255, 255, 255, 0.35);
  transition: color 0.15s;
}

.audio-btn:hover {
  color: #4fc3f7;
}

.translation-main {
  font-size: 22px;
  font-weight: 600;
  color: rgba(147, 197, 253, 0.9);
  text-align: center;
  line-height: 1.5;
}

.example-sentence {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
  font-style: italic;
  text-align: center;
  line-height: 1.5;
  margin-top: 8px;
  white-space: pre-wrap;
}

.flip-hint {
  position: absolute;
  bottom: 14px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.2);
}

.review-actions {
  display: flex;
  gap: 16px;
  opacity: 1;
  transform: translateY(8px);
  transition: opacity 0.3s, transform 0.3s;
}

/* .review-actions.visible {
  opacity: 1;
  transform: translateY(0);
} */

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  height: 40px;
  padding: 0 28px;
  border-radius: 10px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
  outline: none;
}

.action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.action-btn.forget {
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  color: rgba(248, 113, 113, 0.85);
}

.action-btn.forget:hover:not(:disabled) {
  background: rgba(239, 68, 68, 0.18);
  border-color: rgba(239, 68, 68, 0.35);
}

.action-btn.know {
  background: rgba(34, 197, 94, 0.1);
  border: 1px solid rgba(34, 197, 94, 0.2);
  color: rgba(74, 222, 128, 0.85);
}

.action-btn.know:hover:not(:disabled) {
  background: rgba(34, 197, 94, 0.18);
  border-color: rgba(34, 197, 94, 0.35);
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  gap: 12px;
}

.empty-icon {
  font-size: 40px;
}

.empty-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  text-align: center;
}
</style>