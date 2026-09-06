<template>
    <div class="word-table-wrap" :class="{ 'is-resizing': !!resizeState }">
        <div v-if="loading" class="state-tip">加载中...</div>
        <div v-else-if="words.length === 0" class="state-tip">该词库暂无单词</div>
        <template v-else>
            <table ref="headTableRef" class="word-table word-table-head">
                <colgroup>
                    <col v-for="col in columns" :key="col.key" :style="{ width: col.width + '%' }">
                </colgroup>
                <thead>
                    <tr>
                        <th
                            v-for="(col, i) in columns"
                            :key="col.key"
                            :class="'col-' + col.key"
                        >
                            <span v-if="col.key === 'action'" class="action-header">
                                <template v-if="!batchMode">
                                    <span>{{ col.label }}</span>
                                    <button
                                        class="btn-icon btn-icon-danger btn-batch-del"
                                        title="批量删除"
                                        @click="enterBatchMode"
                                    >
                                        <svg width="14" height="14" viewBox="0 0 14 14"><path d="M2.5 4h9M5.5 2h3M3 4v8a1 1 0 001 1h6a1 1 0 001-1V4M5.5 6v4M8.5 6v4" stroke="currentColor" stroke-width="1.2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
                                    </button>
                                </template>
                                <label v-else class="check-all" title="全选">
                                    <span style="margin-right: 5px;">全选</span>
                                    <input
                                        type="checkbox"
                                        :checked="allSelected"
                                        @change="toggleAll"
                                    />
                                </label>
                            </span>
                            <template v-else>
                                {{ col.label }}
                                <span
                                    v-if="i < columns.length - 1"
                                    class="resize-handle"
                                    @mousedown.prevent="startResize($event, i)"
                                ></span>
                            </template>
                        </th>
                    </tr>
                </thead>
            </table>
            <div class="word-table-body">
                <table class="word-table">
                    <colgroup>
                        <col v-for="col in columns" :key="col.key" :style="{ width: col.width + '%' }">
                    </colgroup>
                    <tbody>
                        <tr
                            v-for="w in words"
                            :key="w.id"
                            :class="{ 'row-selected': batchMode && isSelected(w) }"
                        >
                            <td class="col-word">{{ w.word }}</td>
                            <td class="col-phonetic">{{ w.phonetic }}</td>
                            <td class="col-translation">{{ w.translation }}</td>
                            <td class="col-action">
                                <template v-if="batchMode">
                                    <label class="row-check" title="选择">
                                        <input
                                            type="checkbox"
                                            :checked="isSelected(w)"
                                            @change="toggleOne(w)"
                                        />
                                    </label>
                                </template>
                                <button
                                    v-else
                                    class="btn-icon btn-icon-edit"
                                    title="编辑单词"
                                    @click="$emit('edit', w)"
                                >
                                    <svg width="14" height="14" viewBox="0 0 14 14"><path d="M9.5 1.5l3 3M2 11l8-8 2 2-8 8H2v-2z" stroke="currentColor" stroke-width="1.2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- ===== 批量删除操作栏 ===== -->
            <div v-if="batchMode" class="batch-bar">
                <span class="batch-info">已选 <strong>{{ selectedCount }}</strong> 项</span>
                <div class="batch-actions">
                    <button class="btn-batch btn-batch-cancel" @click="exitBatchMode">取消</button>
                    <button
                        class="btn-batch btn-batch-confirm"
                        :disabled="selectedCount === 0 || deleting"
                        @click="confirmBatchDelete"
                    >
                        {{ deleting ? '删除中...' : '删除选中' }}
                    </button>
                </div>
            </div>
        </template>
    </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { GetWordsByVocabularyBankID, DeleteVocabulary } from '../../../wailsjs/go/handler/VocabularyHandler'

const props = defineProps({
    bankId: { type: Number, required: true },
})

const emit = defineEmits(['edit', 'loaded'])

const columns = reactive([
    { key: 'word', label: '单词', width: 22 },
    { key: 'phonetic', label: '音标', width: 18 },
    { key: 'translation', label: '释义', width: 45 },
    { key: 'action', label: '操作', width: 15 },
])

const words = reactive([])
const loading = ref(false)

async function loadWords() {
    loading.value = true
    try {
        const list = await GetWordsByVocabularyBankID(props.bankId)
        words.splice(0, words.length, ...list)
        emit('loaded', { count: list.length })
        // 数据刷新后重置批量选择状态
        selectedIds.clear()
        batchMode.value = false
    } catch (err) {
        console.error('加载单词失败:', err)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    loadWords()
})

defineExpose({ reload: loadWords })

const headTableRef = ref(null)
const resizeState = ref(null)

function startResize(e, index) {
    const tableEl = headTableRef.value
    if (!tableEl) return
    const tableWidth = tableEl.getBoundingClientRect().width
    resizeState.value = {
        index,
        startX: e.clientX,
        startWidths: [columns[index].width, columns[index + 1].width],
        tableWidth,
    }
    document.addEventListener('mousemove', onResizeMove)
    document.addEventListener('mouseup', onResizeEnd)
}

function onResizeMove(e) {
    const state = resizeState.value
    if (!state) return

    const dx = e.clientX - state.startX
    const deltaPercent = (dx / state.tableWidth) * 100
    const minWidth = 8

    let newWidth = state.startWidths[0] + deltaPercent
    let nextWidth = state.startWidths[1] - deltaPercent

    if (newWidth < minWidth) {
        newWidth = minWidth
        nextWidth = state.startWidths[0] + state.startWidths[1] - minWidth
    }
    if (nextWidth < minWidth) {
        nextWidth = minWidth
        newWidth = state.startWidths[0] + state.startWidths[1] - minWidth
    }

    columns[state.index].width = Math.round(newWidth * 10) / 10
    columns[state.index + 1].width = Math.round(nextWidth * 10) / 10
}

function onResizeEnd() {
    resizeState.value = null
    document.removeEventListener('mousemove', onResizeMove)
    document.removeEventListener('mouseup', onResizeEnd)
}

// ===== 批量删除 =====
const batchMode = ref(false)
const selectedIds = reactive(new Set())
const deleting = ref(false)

const selectedCount = computed(() => selectedIds.size)
const allSelected = computed(() => words.length > 0 && selectedIds.size === words.length)

function enterBatchMode() {
    batchMode.value = true
}

function exitBatchMode() {
    batchMode.value = false
    selectedIds.clear()
}

function toggleAll(e) {
    if (e.target.checked) {
        words.forEach((w) => selectedIds.add(w.id))
    } else {
        selectedIds.clear()
    }
}

function toggleOne(w) {
    if (selectedIds.has(w.id)) {
        selectedIds.delete(w.id)
    } else {
        selectedIds.add(w.id)
    }
}

function isSelected(w) {
    return selectedIds.has(w.id)
}

async function confirmBatchDelete() {
    const ids = Array.from(selectedIds)
    if (ids.length === 0) return
    deleting.value = true
    try {
        for (const id of ids) {
            await DeleteVocabulary(id)
        }
        await loadWords()
    } catch (err) {
        console.error('批量删除失败:', err)
    } finally {
        deleting.value = false
    }
}

onUnmounted(() => {
    document.removeEventListener('mousemove', onResizeMove)
    document.removeEventListener('mouseup', onResizeEnd)
})
</script>

<style scoped>
.word-table-wrap {
    flex: 1;
    overflow: hidden;
    display: flex;
    flex-direction: column;
    border: 1px solid rgba(255, 255, 255, 0.05);
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.02);
}

.state-tip {
    padding: 40px;
    text-align: center;
    color: rgba(255, 255, 255, 0.4);
    font-size: 14px;
}

.word-table {
    width: 100%;
    border-collapse: collapse;
    font-size: 13px;
    table-layout: fixed;
}

.is-resizing {
    user-select: none;
    cursor: col-resize;
}
.is-resizing * {
    cursor: col-resize;
}

.word-table-head {
    flex-shrink: 0;
}

.word-table th {
    position: relative;
    padding: 10px 14px;
    text-align: left;
    font-weight: 600;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.5);
    background: rgba(30, 42, 58, 0.98);
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    white-space: nowrap;
}

.resize-handle {
    position: absolute;
    right: 0;
    top: 0;
    bottom: 0;
    width: 4px;
    cursor: col-resize;
    z-index: 1;
}
.resize-handle:hover,
.resize-handle:active {
    background: rgba(79, 195, 247, 0.5);
}

.word-table-body {
    flex: 1;
    overflow: auto;
}

.word-table-body::-webkit-scrollbar {
    width: 8px;
    height: 8px;
}
.word-table-body::-webkit-scrollbar-track {
    background: transparent;
}
.word-table-body::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.1);
    border-radius: 4px;
    transition: background 0.2s;
}
.word-table-body::-webkit-scrollbar-thumb:hover {
    background: rgba(79, 195, 247, 0.4);
}
.word-table-body::-webkit-scrollbar-corner {
    background: transparent;
}

.word-table td {
    padding: 8px 14px;
    color: rgba(255, 255, 255, 0.75);
    border-bottom: 1px solid rgba(255, 255, 255, 0.03);
    word-break: break-word;
}

.word-table tbody tr:hover {
    background: rgba(255, 255, 255, 0.03);
}

.col-word {
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.col-phonetic {
    color: rgba(255, 255, 255, 0.45);
    font-size: 12px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.col-translation {
    color: rgba(255, 255, 255, 0.75);
    white-space: pre-line;
    word-break: break-word;
    overflow-wrap: break-word;
}

.col-action {
    text-align: center;
    white-space: nowrap;
}
.col-action .btn-icon {
    display: inline-flex;
    vertical-align: middle;
}
.col-action .btn-icon + .btn-icon {
    margin-left: 4px;
}

.btn-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
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

.btn-icon-edit:hover {
    background: rgba(79, 195, 247, 0.15);
    color: #4fc3f7;
}

/* ===== 表头操作列（批量删除入口） ===== */
.action-header {
    display: flex;
    align-items: center;
    gap: 8px;
}

.btn-batch-del {
    width: 22px;
    height: 22px;
}
.btn-batch-del:hover {
    background: rgba(232, 17, 35, 0.15);
    color: #e81123;
}

/* ===== 勾选框 ===== */
.row-check,
.check-all {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;

}

.row-check input,
.check-all input {
    appearance: none;
    -webkit-appearance: none;
    width: 15px;
    height: 15px;
    margin: 0;
    border: 1.5px solid rgba(255, 255, 255, 0.35);
    border-radius: 4px;
    background: transparent;
    cursor: pointer;
    position: relative;
    transition: all 0.15s;
}

.row-check input:hover,
.check-all input:hover {
    border-color: rgba(255, 255, 255, 0.75);
}

.row-check input:checked,
.check-all input:checked {
    background: #e81123;
    border-color: #e81123;
}

.row-check input:checked::after,
.check-all input:checked::after {
    content: '';
    position: absolute;
    left: 4px;
    top: 1px;
    width: 4px;
    height: 8px;
    border: solid #fff;
    border-width: 0 2px 2px 0;
    transform: rotate(45deg);
}

/* ===== 选中行高亮 ===== */
.word-table tbody tr.row-selected {
    background: rgba(232, 17, 35, 0.12);
}

/* ===== 批量删除操作栏 ===== */
.batch-bar {
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 12px;
    padding: 10px 14px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);
    background: rgba(30, 42, 58, 0.98);
    border-radius: 0 0 8px 8px;
}

.batch-info {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.6);
}
.batch-info strong {
    color: #e81123;
}

.batch-actions {
    display: flex;
    gap: 8px;
}

.btn-batch {
    display: inline-flex;
    align-items: center;
    padding: 6px 16px;
    border: none;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
}
.btn-batch-cancel {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.7);
}
.btn-batch-cancel:hover {
    background: rgba(255, 255, 255, 0.14);
}
.btn-batch-confirm {
    background: #e81123;
    color: #fff;
}
.btn-batch-confirm:hover:not(:disabled) {
    background: #c90e1d;
}
.btn-batch-confirm:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}
</style>