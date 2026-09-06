<template>
    <!-- 收藏弹窗：选择目标词库 -->
    <FloatPopupPanel
        :visible="WordStore.showStarPopup"
        title="⭐ 收藏到词库"
        max-width="280px"
        @close="onClose"
    >
        <div class="star-popup-word">{{ WordStore.currentWord?.word }}</div>
        <div v-if="VocFloatWinStore.vocabularyBanks.length === 0" class="star-popup-empty">
            暂无词库，请先在设置中创建词库
        </div>
        <template v-else>
            <!-- 下拉选择目标词库（检索所有词库，初始未选择） -->
            <div class="bank-select-wrap">
                <div class="bank-select-trigger" :class="{ open: bankListOpen }" @click="toggleBankList">
                    <span class="bank-select-value" :class="{ placeholder: !selectedBankId }">
                        {{ selectedBankName }}
                    </span>
                    <span class="bank-select-arrow">▾</span>
                </div>
                <Transition name="dropdown">
                    <ul v-if="bankListOpen" class="bank-list">
                        <li v-for="bank in VocFloatWinStore.vocabularyBanks" :key="bank.id"
                            class="bank-option"
                            :class="{ active: bank.id === selectedBankId }"
                            @click="chooseBank(bank.id)">
                            <span class="bank-option-name">{{ bank.name }}</span>
                            <span v-if="bank.id === selectedBankId" class="bank-option-check">✓</span>
                        </li>
                    </ul>
                </Transition>
            </div>
            <!-- 底部按钮：
                已收藏且选中同一词库 → 取消收藏
                未收藏或切换词库 → 收藏
            -->
            <button
                class="star-confirm-btn"
                :class="{ cancel: isSameBank, disabled: !isSameBank && !selectedBankId }"
                :disabled="!isSameBank && !selectedBankId"
                @click="onConfirm"
            >
                {{ isSameBank ? '取消收藏' : '收藏' }}
            </button>
        </template>
    </FloatPopupPanel>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { WordStore, VocFloatWinStore } from '@/store.js';
import FloatPopupPanel from '../../template/FloatPopupPanel.vue';

// 下拉列表是否展开
const bankListOpen = ref(false);
// 当前选中的目标词库 ID（0 = 未选择）
const selectedBankId = ref(0);

// 打开弹窗时：如果已收藏则默认选中已收藏的词库，否则重置为「未选择」；刷新词库列表
watch(() => WordStore.showStarPopup, (open) => {
    if (open) {
        bankListOpen.value = false;
        if (WordStore.isStarred && WordStore.starredWord.vocabulary_bank_id) {
            selectedBankId.value = WordStore.starredWord.vocabulary_bank_id;
        } else {
            selectedBankId.value = 0;
        }
        VocFloatWinStore.loadVocabularyBanks();
    }
});

const selectedBankName = computed(() => {
    if (!selectedBankId.value) return '未选择';
    const bank = VocFloatWinStore.vocabularyBanks.find(b => b.id === selectedBankId.value);
    return bank ? bank.name : '未选择';
});

// 当前选中的词库是否就是已收藏的词库（用于判断按钮是"取消收藏"还是"收藏"）
const isSameBank = computed(() => {
    return WordStore.isStarred && selectedBankId.value === WordStore.starredBankId;
});

function toggleBankList() {
    bankListOpen.value = !bankListOpen.value;
}

function chooseBank(id) {
    selectedBankId.value = id;
    bankListOpen.value = false;
}

function onClose() {
    bankListOpen.value = false;
    WordStore.closeStarPopup();
}

// 底部按钮：已收藏且同一词库 → 取消收藏；否则 → 收藏到所选词库
function onConfirm() {
    if (isSameBank.value) {
        WordStore.uncollect();
    } else if (selectedBankId.value) {
        WordStore.collectToBank(selectedBankId.value);
    }
}
</script>

<style scoped>
/* ===== 收藏弹窗内容 ===== */
.star-popup-word {
    font-size: 18px;
    font-weight: 700;
    color: #4fc3f7;
    margin-bottom: 12px;
    text-align: center;
}
.star-popup-empty {
    text-align: center;
    color: rgba(255, 255, 255, 0.4);
    font-size: 13px;
    padding: 16px 0;
}

/* ===== 下拉选择词库（与 SettingsPanel 保持一致） ===== */
.bank-select-wrap {
    position: relative;
    margin-bottom: 12px;
}

.bank-select-trigger {
    display: flex;
    align-items: center;
    justify-content: space-between;
    background: rgba(255, 255, 255, 0.08);
    border: 1px solid rgba(255, 255, 255, 0.15);
    border-radius: 6px;
    color: rgba(255, 255, 255, 0.85);
    font-size: 13px;
    padding: 8px 12px;
    cursor: pointer;
    user-select: none;
    transition: border-color 0.2s, background 0.2s;
}

.bank-select-trigger:hover {
    background: rgba(255, 255, 255, 0.12);
    border-color: rgba(255, 255, 255, 0.25);
}

.bank-select-trigger.open {
    border-color: #4fc3f7;
}

.bank-select-value.placeholder {
    color: rgba(255, 255, 255, 0.4);
}

.bank-select-arrow {
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

/* ===== 底部收藏/取消收藏按钮 ===== */
.star-confirm-btn {
    width: 100%;
    padding: 9px 0;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    color: #fff;
    box-shadow: 0 4px 12px rgba(79, 195, 247, 0.3);
}

.star-confirm-btn:hover {
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(79, 195, 247, 0.4);
}

/* 已收藏：按钮切换为取消收藏样式 */
.star-confirm-btn.cancel {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.8);
    box-shadow: none;
}

.star-confirm-btn.cancel:hover {
    background: rgba(255, 82, 82, 0.25);
    color: #ff6b6b;
}

/* 未选择词库且未收藏时不可点击 */
.star-confirm-btn.disabled {
    opacity: 0.4;
    cursor: not-allowed;
    transform: none;
    box-shadow: none;
}

.star-confirm-btn.disabled:hover {
    transform: none;
    box-shadow: none;
}
</style>