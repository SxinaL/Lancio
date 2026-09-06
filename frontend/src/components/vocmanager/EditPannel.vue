<template>
    <MainPopupPanel :visible="visible" title="✏️ 编辑单词" card-class="modal-card-edit" @close="$emit('close')">
        <div class="form-row">
            <label class="form-label">单词</label>
            <input class="form-input" v-model="form.word" placeholder="请输入单词" />
        </div>
        <div class="form-row">
            <label class="form-label">音标</label>
            <input class="form-input" v-model="form.phonetic" placeholder="如 /ˈæpl/" />
        </div>
        <div class="form-row">
            <label class="form-label">释义</label>
            <textarea class="form-textarea" v-model="form.translation" rows="3" placeholder="如 n. 苹果"></textarea>
        </div>
        <div class="form-row">
            <label class="form-label">例句</label>
            <textarea class="form-textarea" v-model="form.exampleSentence" rows="2" placeholder="可选"></textarea>
        </div>
        <div class="form-row">
            <label class="form-label">所属词库</label>
            <div ref="bankSelectRef" class="bank-select-wrap">
                <div class="bank-select-trigger" :class="{ open: bankListOpen }" @click="toggleBankList">
                    <span class="bank-select-value" :class="{ placeholder: !form.vocabularyId }">
                        {{ selectedBankName }}
                    </span>
                    <span class="bank-select-arrow">▾</span>
                </div>
                <Transition name="dropdown">
                    <ul v-if="bankListOpen" class="bank-list">
                        <li
                            v-for="bank in banks"
                            :key="bank.id"
                            class="bank-option"
                            :class="{ active: bank.id === form.vocabularyId }"
                            @click="chooseBank(bank.id)"
                        >
                            <span class="bank-option-name">{{ bank.name }}</span>
                            <span v-if="bank.id === form.vocabularyId" class="bank-option-check">✓</span>
                        </li>
                        <li v-if="!banks.length" class="bank-option empty">暂无词库</li>
                    </ul>
                </Transition>
            </div>
        </div>
        <template #footer>
            <button class="btn btn-secondary" @click="$emit('close')">取消</button>
            <button class="btn btn-primary" :disabled="saving" @click="$emit('save', form)">
                {{ saving ? '保存中...' : '保存' }}
            </button>
        </template>
    </MainPopupPanel>
</template>

<script setup>
import { reactive, ref, computed, watch, onMounted, onUnmounted } from 'vue'
import MainPopupPanel from '../../template/MainPopupPanel.vue'
import { GetAllVocabularyBanks } from '../../../wailsjs/go/handler/VocabularyBankHandler'

const props = defineProps({
    visible: { type: Boolean, default: false },
    word: { type: Object, default: () => ({}) },
    saving: { type: Boolean, default: false },
})

defineEmits(['close', 'save'])

const form = reactive({
    id: 0,
    word: '',
    phonetic: '',
    translation: '',
    exampleSentence: '',
    vocabularyId: 0,
})

// 全部词库列表，用于选择单词所属词库
const banks = ref([])
const bankListOpen = ref(false)
const bankSelectRef = ref(null)

const selectedBankName = computed(() => {
    if (!form.vocabularyId) return '请选择词库'
    const bank = banks.value.find((b) => b.id === form.vocabularyId)
    return bank ? bank.name : '请选择词库'
})

// 刷新词库列表
async function refreshBanks() {
    try {
        banks.value = await GetAllVocabularyBanks()
    } catch (err) {
        console.error('获取词库列表失败:', err)
    }
}

function toggleBankList() {
    bankListOpen.value = !bankListOpen.value
    if (bankListOpen.value) {
        refreshBanks()
    }
}

function chooseBank(id) {
    form.vocabularyId = id
    bankListOpen.value = false
}

// 点击外部关闭下拉
function onDocClick(e) {
    if (bankSelectRef.value && !bankSelectRef.value.contains(e.target)) {
        bankListOpen.value = false
    }
}

onMounted(() => document.addEventListener('mousedown', onDocClick))
onUnmounted(() => document.removeEventListener('mousedown', onDocClick))

watch(() => props.visible, async (val) => {
    if (val && props.word) {
        form.id = props.word.id || 0
        form.word = props.word.word || ''
        form.phonetic = props.word.phonetic || ''
        form.translation = props.word.translation || ''
        form.exampleSentence = props.word.example_sentence || ''
        form.vocabularyId = props.word.vocabulary_bank_id || 0
    }
    if (val) {
        bankListOpen.value = false
        refreshBanks()
    }
})
</script>

<style scoped>
.form-row {
    margin-bottom: 14px;
}

.form-row:last-child {
    margin-bottom: 0;
}

.form-label {
    display: block;
    margin-bottom: 6px;
    font-size: 12px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.6);
    letter-spacing: 0.3px;
}

.form-input,
.form-textarea {
    width: 100%;
    padding: 8px 12px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 6px;
    color: rgba(255, 255, 255, 0.9);
    font-size: 13px;
    font-family: inherit;
    outline: none;
    transition: all 0.15s;
    box-sizing: border-box;
}

.form-input:focus,
.form-textarea:focus {
    border-color: rgba(79, 195, 247, 0.5);
    background: rgba(79, 195, 247, 0.04);
}

.form-textarea {
    resize: vertical;
    min-height: 60px;
    line-height: 1.5;
}

.bank-select-wrap {
    position: relative;
    width: 100%;
}

.bank-select-trigger {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 6px;
    color: rgba(255, 255, 255, 0.9);
    font-size: 13px;
    padding: 8px 12px;
    cursor: pointer;
    user-select: none;
    box-sizing: border-box;
    transition: border-color 0.2s, background 0.2s;
}

.bank-select-trigger:hover {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.2);
}

.bank-select-trigger.open {
    border-color: rgba(79, 195, 247, 0.5);
    background: rgba(79, 195, 247, 0.04);
}

.bank-select-value {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.bank-select-value.placeholder {
    color: rgba(255, 255, 255, 0.35);
}

.bank-select-arrow {
    margin-left: 8px;
    font-size: 11px;
    color: rgba(255, 255, 255, 0.5);
    transition: transform 0.2s;
}

.bank-select-trigger.open .bank-select-arrow {
    transform: rotate(180deg);
}

.bank-list {
    position: absolute;
    top: calc(100% + 4px);
    left: 0;
    right: 0;
    z-index: 1000;
    max-height: 180px;
    overflow-y: auto;
    list-style: none;
    margin: 0;
    padding: 4px;
    background: linear-gradient(135deg, #243447 0%, #1c1c31 100%);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 8px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.5);
}

.bank-option {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 7px 10px;
    border-radius: 5px;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.8);
    cursor: pointer;
    transition: background 0.15s, color 0.15s;
}

.bank-option:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #fff;
}

.bank-option.active {
    background: rgba(79, 195, 247, 0.18);
    color: #4fc3f7;
}

.bank-option.empty {
    justify-content: center;
    color: rgba(255, 255, 255, 0.4);
    cursor: default;
}

.bank-option-check {
    color: #4fc3f7;
    font-size: 13px;
}

.dropdown-enter-active,
.dropdown-leave-active {
    transition: opacity 0.18s, transform 0.18s;
    transform-origin: top center;
}

.dropdown-enter-from,
.dropdown-leave-to {
    opacity: 0;
    transform: translateY(-6px);
}

.btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 7px 16px;
    border: none;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
}

.btn-primary {
    background: #4fc3f7;
    color: #fff;
}

.btn-primary:hover:not(:disabled) {
    background: #29b6f6;
}

.btn-primary:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.7);
}

.btn-secondary:hover {
    background: rgba(255, 255, 255, 0.14);
}
</style>