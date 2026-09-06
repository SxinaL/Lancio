<template>
  <div class="bank-manager">

    <!-- ===== 页面标题 ===== -->
    <div class="page-header">
      <h2 class="page-title">
        单词管理
      </h2>

    </div>

    <!-- ===== 词库列表 ===== -->
    <div class="bank-list">
      <div v-if="banks.length === 0" class="empty-state">
        <span class="empty-icon">📭</span>
        <p>暂无词库，点击右下角 + 添加</p>
      </div>

      <div
        v-for="bank in banks"
        :key="bank.id"
        class="bank-item"
      >
        <div class="bank-info">
          <div class="bank-icon">📖</div>
          <div class="bank-detail">
            <span class="bank-name">{{ bank.name }}</span>
            <span class="bank-meta">{{ bank.wordCount }} 个单词</span>
          </div>
        </div>
        <div class="bank-actions">
          <button class="btn-icon" title="打开词库" @click="openBank(bank)">
            <svg width="14" height="14" viewBox="0 0 14 14"><path d="M3 2h5l3 3v7a1 1 0 01-1 1H3a1 1 0 01-1-1V3a1 1 0 011-1zm5 0v3h3" stroke="currentColor" stroke-width="1.2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
          <button class="btn-icon" title="编辑名称" @click="openEditDialog(bank)">
            <svg width="14" height="14" viewBox="0 0 14 14"><path d="M2 10.5V12h1.5l7.37-7.37-1.5-1.5L2 10.5zm9.96-6.96a.5.5 0 000-.7l-.8-.8a.5.5 0 00-.7 0l-.73.73 1.5 1.5.73-.73z" fill="currentColor"/></svg>
          </button>
          <button class="btn-icon btn-icon-danger" title="删除词库" @click="openDeleteConfirm(bank)">
            <svg width="14" height="14" viewBox="0 0 14 14"><path d="M2.5 4h9M5.5 2h3M3 4v8a1 1 0 001 1h6a1 1 0 001-1V4M5.5 6v4M8.5 6v4" stroke="currentColor" stroke-width="1.2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
        </div>
      </div>
    </div>

    <!-- ===== 添加 / 编辑 对话框 ===== -->
    <AddPannel
      :visible="dialogVisible"
      :isEditing="isEditing"
      :initialName="editingName"
      :submitting="submitting"
      @close="closeDialog"
      @confirm="confirmDialog"
    />

    <!-- ===== 删除确认弹窗 ===== -->
    <DeletePanel
      :visible="!!deleteTarget"
      :bankName="deleteTarget?.name"
      :wordCount="deleteTarget?.wordCount"
      :deleting="deleting"
      @close="deleteTarget = null"
      @confirm="confirmDelete"
    />

    <!-- ===== 浮动添加按钮 ===== -->
    <button class="fab" title="添加词库" @click="openAddDialog">
      <svg width="22" height="22" viewBox="0 0 22 22">
        <path d="M11 4v14M4 11h14" stroke="currentColor" stroke-width="2.4" stroke-linecap="round"/>
      </svg>
    </button>
  </div>

  
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { GetAllVocabularyBanks, CreateVocabularyBank, DeleteVocabularyBank, GetVocabularyCountByVocabularyBankID } from '../../wailsjs/go/handler/VocabularyBankHandler'
import AddPannel from '@/components/bankmanager/AddPannel.vue'
import DeletePanel from '@/components/bankmanager/DeletePanel.vue'

const router = useRouter()

function openBank(bank) {
  router.push({ name: 'VocManager', params: { id: bank.id }, query: { name: bank.name } })
}

// ===== 词库数据 =====
const banks = reactive([])
const loading = ref(false)

async function loadBanks() {
  loading.value = true
  try {
    const list = await GetAllVocabularyBanks()
    banks.splice(0, banks.length)
    for (const v of list) {
      let wordCount = 0
      try {
        wordCount = await GetVocabularyCountByVocabularyBankID(v.id)
      } catch (_) {}
      banks.push({ id: v.id, name: v.name, wordCount })
    }
  } catch (err) {
    console.error('加载词库失败:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadBanks()
})

// ===== 添加 / 编辑 对话框 =====
const dialogVisible = ref(false)
const isEditing = ref(false)
const editingId = ref(null)
const editingName = ref('')
const submitting = ref(false)

function openAddDialog() {
  isEditing.value = false
  editingId.value = null
  editingName.value = ''
  dialogVisible.value = true
}

function openEditDialog(bank) {
  isEditing.value = true
  editingId.value = bank.id
  editingName.value = bank.name
  dialogVisible.value = true
}

function closeDialog() {
  dialogVisible.value = false
  editingName.value = ''
}

async function confirmDialog(name) {
  if (!name) return

  if (isEditing.value) {
    const bank = banks.find(b => b.id === editingId.value)
    if (bank) bank.name = name
    closeDialog()
  } else {
    submitting.value = true
    try {
      const created = await CreateVocabularyBank(name, '')
      banks.push({ id: created.id, name: created.name, wordCount: 0 })
      closeDialog()
    } catch (err) {
      console.error('创建词库失败:', err)
    } finally {
      submitting.value = false
    }
  }
}

// ===== 删除确认 =====
const deleteTarget = ref(null)
const deleting = ref(false)

function openDeleteConfirm(bank) {
  deleteTarget.value = bank
}

async function confirmDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await DeleteVocabularyBank(deleteTarget.value.id)
    const idx = banks.findIndex(b => b.id === deleteTarget.value.id)
    if (idx !== -1) banks.splice(idx, 1)
    deleteTarget.value = null
  } catch (err) {
    console.error('删除词库失败:', err)
  } finally {
    deleting.value = false
  }
}
</script>

<style scoped>
/* ===== 容器 ===== */
.bank-manager {
  /* flex: 1;
  height: 100%; */
  width: 100%;
  display: flex;
  flex-direction: column;
  padding: 16px 24px;
  overflow: hidden;
  gap: 16px;
  /* border: 1px solid rgb(255, 0, 0); */
}

/* ===== 页面头部 ===== */
.page-header {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16px;
  flex-shrink: 0;
  /* border: 1px solid rgb(255, 0, 0); */
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

.title-icon {
  font-size: 22px;
}

/* ===== 操作按钮 ===== */
.btn-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  background: transparent;
  border: none;
  border-radius: 5px;
  color: rgba(255, 255, 255, 0.35);
  cursor: pointer;
  transition: all 0.15s;
}
.btn-icon:hover {
  background: rgba(255, 255, 255, 0.08);
  color: rgba(255, 255, 255, 0.8);
}
.btn-icon-danger:hover {
  background: rgba(232, 17, 35, 0.15);
  color: #e81123;
}

/* ===== 词库列表 ===== */
.bank-list {
  flex: 1;
  overflow-y: auto;
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 10px;
  align-content: start;
  align-items: start;
  /* border: 1px solid rgb(255, 0, 0); */
}

.bank-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.025);
  border: 1px solid rgba(255, 255, 255, 0.04);
  transition: background 0.15s, border-color 0.15s;
}
.bank-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.08);
}

.bank-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.bank-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.bank-detail {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.bank-name {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.85);
}

.bank-meta {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
}

.bank-actions {
  display: flex;
  gap: 4px;
  opacity: 0;
  transition: opacity 0.15s;
}
.bank-item:hover .bank-actions {
  opacity: 1;
}

/* ===== 空状态 ===== */
.empty-state {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.2);
}
.empty-icon {
  font-size: 40px;
}
.empty-state p {
  margin: 0;
  font-size: 13px;
}

/* ===== 浮动添加按钮 (FAB) ===== */
.fab {
  position: fixed;
  right: 24px;
  bottom: 24px;
  width: 48px;
  height: 48px;
  border: none;
  border-radius: 50%;
  background: linear-gradient(135deg, #4fc3f7, #29b6f6);
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 16px rgba(79, 195, 247, 0.4);
  transition: all 0.2s ease;
  z-index: 500;
}
.fab:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(79, 195, 247, 0.55);
}
.fab:active {
  transform: translateY(0);
}
</style>