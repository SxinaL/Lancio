<template>
  <div class="voc-manager">
    <!-- ===== 页面头部 ===== -->
    <div class="page-header">
      <button class="btn-back" title="返回" @click="goBack">
        <svg width="16" height="16" viewBox="0 0 16 16"><path d="M10 3L5 8l5 5" stroke="currentColor" stroke-width="1.6" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
      </button>
      <h2 class="page-title">
        📖 {{ bankName }}
        <span class="title-meta">{{ wordCount }} 个单词</span>
      </h2>
    </div>

    <!-- ===== 单词表格 ===== -->
    <WordTable
      ref="wordTableRef"
      :bankId="Number(route.params.id)"
      @edit="openEditDialog"
      @loaded="onWordsLoaded"
    />

    <!-- ===== 编辑单词弹窗 ===== -->
    <EditPannel
      :visible="editVisible"
      :word="editWord"
      :saving="saving"
      @close="editVisible = false"
      @save="saveEdit"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { UpdateVocabulary } from '../../wailsjs/go/handler/VocabularyHandler'
import EditPannel from '@/components/vocmanager/EditPannel.vue'
import WordTable from '@/components/vocmanager/WordTable.vue'

const route = useRoute()
const router = useRouter()

const bankName = ref(route.query.name || '词库')
const wordCount = ref(0)
const wordTableRef = ref(null)

function onWordsLoaded({ count }) {
  wordCount.value = count
}

function goBack() {
  router.push({ name: 'BankManager' })
}

// ===== 编辑单词 =====
const editVisible = ref(false)
const editWord = ref(null)
const saving = ref(false)

function openEditDialog(w) {
  editWord.value = w
  editVisible.value = true
}

async function saveEdit(form) {
  if (!form.word.trim()) {
    alert('单词不能为空')
    return
  }
  saving.value = true
  try {
    await UpdateVocabulary({
      id: form.id,
      word: form.word,
      phonetic: form.phonetic,
      translation: form.translation,
      example_sentence: form.exampleSentence,
      vocabulary_bank_id: form.vocabularyId,
    })
    editVisible.value = false
    wordTableRef.value?.reload()
  } catch (err) {
    console.error('保存单词失败:', err)
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
/* ===== 容器 ===== */
.voc-manager {
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
  gap: 12px;
  flex-shrink: 0;
}

.btn-back {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 6px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.15s;
}
.btn-back:hover {
  background: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.95);
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

</style>