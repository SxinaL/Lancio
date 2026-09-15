<template>
    <div class="review-plan">
        <div class="page-header">
            <h2 class="page-title">
                🔄 复习计划
                <span class="title-meta">翻开卡片，检验记忆</span>
            </h2>
        </div>

        <div class="page-body">
            <ReviewCard ref="reviewCardRef" :wordList="wordList" @done="onReviewDone" />
            <button v-if="showStartBtn" class="btn-start-round" :disabled="loading" @click="startNewRound">
                <svg width="18" height="18" viewBox="0 0 14 14">
                    <path d="M2 2h4v3M12 12H8V9" stroke="currentColor" stroke-width="1.3" fill="none"
                        stroke-linecap="round" stroke-linejoin="round" />
                    <path d="M11 5a4 4 0 00-6.93-2M3 9a4 4 0 006.93 2" stroke="currentColor" stroke-width="1.3"
                        fill="none" stroke-linecap="round" />
                </svg>
                {{ loading ? '加载中...' : '开启新一轮复习' }}
            </button>
        </div>

        <div class="page-footer">
            <div class="footer-status" v-if="!reviewDone">
                <svg width="14" height="14" viewBox="0 0 14 14">
                    <circle cx="7" cy="7" r="6" stroke="currentColor" stroke-width="1.2" fill="none" />
                    <line x1="7" y1="4" x2="7" y2="8" stroke="currentColor" stroke-width="1.4" stroke-linecap="round" />
                    <circle cx="7" cy="10.5" r="0.7" fill="currentColor" />
                </svg>
                <span>点击卡片翻转查看释义，然后选择"认识"或"不认识"</span>

            </div>

        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import ReviewCard from '@/components/reviewplan/ReviewCard.vue'
import { GetReviewWords, GetActiveReviewWords, ResetReviewRound } from '../../wailsjs/go/handler/VocabularyBankHandler'
import { GetReviewWordsPerSession } from '../../wailsjs/go/handler/ConfigHandler'

const reviewCardRef = ref(null)
const wordList = ref([])       // 本轮复习单词列表（pending + reviewed），非空时按钮隐藏
const reviewDone = ref(false)  // 本轮是否已完成
const loading = ref(false)

// 按钮显隐：根据子组件当前的复习单词是否为空来判断
const showStartBtn = computed(() => !reviewCardRef.value?.currentWord)

// 点击"开启新一轮复习"：从配置读取数量 → 后端选词并标记pending → 加载卡片
async function startNewRound() {
    loading.value = true
    reviewDone.value = false
    try {
        const count = await GetReviewWordsPerSession()
        wordList.value = await GetReviewWords(count)
    } catch (err) {
        console.error('加载复习单词失败:', err)
        wordList.value = []
    } finally {
        loading.value = false
    }
}

// ReviewCard 全部完成时回调：清空列表，重置所有reviewed→unselected，为下一轮准备
function onReviewDone() {
  reviewDone.value = true
//   wordList.value = []
  ResetReviewRound().catch(err => console.error('重置复习状态失败:', err))
}

// 页面加载时检查是否有进行中的复习轮次，有则直接恢复
onMounted(async () => {
  try {
    const activeWords = await GetActiveReviewWords()
    if (activeWords && activeWords.length > 0) {
      wordList.value = activeWords
      reviewDone.value = false
    }
  } catch (err) {
    console.error('加载活跃复习单词失败:', err)
  }
})
</script>

<style scoped>
.review-plan {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    padding: 16px 24px;
    overflow: hidden;
}

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

.page-body {
    /* border:1px solid rgba(255, 255, 255, 0.06); */
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    overflow-y: auto;
}

.page-footer {
    flex-shrink: 0;
    padding-top: 16px;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
}

.footer-status {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.3);
}

.btn-start-round {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    height: 48px;
    padding: 0 32px;
    background: rgba(59, 130, 246, 0.15);
    border: 1px solid rgba(59, 130, 246, 0.25);
    border-radius: 12px;
    font-size: 16px;
    font-weight: 500;
    color: rgba(96, 165, 250, 0.9);
    cursor: pointer;
    transition: all 0.15s;
    outline: none;
}

.btn-start-round:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.btn-start-round:hover:not(:disabled) {
    background: rgba(59, 130, 246, 0.25);
    border-color: rgba(59, 130, 246, 0.4);
}
</style>