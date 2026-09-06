<template>
    <FloatPopupPanel
        :visible="VocFloatWinStore.showSettings"
        title="⚙️ 设置"
        @close="VocFloatWinStore.closeSettings()"
    >
        <div class="setting-item">
            <div class="setting-label">
                <span class="setting-icon">⭐</span>
                <span>加载词库</span>
            </div>
            <div class="setting-control">
                <div class="bank-select-wrap">
                    <div class="bank-select-trigger" :class="{ open: bankListOpen }"
                        @click="toggleBankList">
                        <span class="bank-select-value">{{ selectedBankName }}</span>
                        <span class="bank-select-arrow">▾</span>
                    </div>
                    <Transition name="dropdown">
                        <ul v-if="bankListOpen" class="bank-list">
                            <li v-for="bank in VocFloatWinStore.vocabularyBanks" :key="bank.id"
                                class="bank-option"
                                :class="{ active: bank.id === VocFloatWinStore.selectedBankId }"
                                @click="chooseBank(bank.id)">
                                <span class="bank-option-name">{{ bank.name }}</span>
                                <span v-if="bank.id === VocFloatWinStore.selectedBankId"
                                    class="bank-option-check">✓</span>
                            </li>
                            <li v-if="!VocFloatWinStore.vocabularyBanks.length"
                                class="bank-option empty">
                                暂无词库
                            </li>
                        </ul>
                    </Transition>
                </div>
            </div>
        </div>
        <div class="setting-item">
            <div class="setting-label">
                <span class="setting-icon">🔆</span>
                <span>窗口透明度</span>
            </div>
            <div class="setting-control">
                <input type="range" class="opacity-slider" min="10" max="100"
                    :value="Math.round(VocFloatWinStore.opacity * 100)"
                    @input="VocFloatWinStore.setOpacity($event.target.value)"
                    @change="VocFloatWinStore.commitOpacity()" />
                <span class="opacity-value">{{ Math.round(VocFloatWinStore.opacity * 100) }}%</span>
            </div>
        </div>
    </FloatPopupPanel>
</template>

<script setup>
import { ref, computed } from 'vue';
import { VocFloatWinStore } from '@/store.js';
import FloatPopupPanel from '../../template/FloatPopupPanel.vue';

const bankListOpen = ref(false);

const selectedBankName = computed(() => {
    const bank = VocFloatWinStore.vocabularyBanks.find(b => b.id === VocFloatWinStore.selectedBankId);
    return bank ? bank.name : '未选择';
});

function toggleBankList() {
    bankListOpen.value = !bankListOpen.value;
    if (bankListOpen.value) {
        VocFloatWinStore.loadVocabularyBanks();
    }
}

async function chooseBank(id) {
    bankListOpen.value = false;
    await VocFloatWinStore.selectVocabularyBank(id);
}
</script>

<style scoped>
.setting-item {
    margin-bottom: 4px;
}

.setting-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.7);
    margin-bottom: 10px;
}

.setting-icon {
    font-size: 16px;
}

.setting-control {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 0 4px;
}

.bank-select-wrap {
    position: relative;
    flex: 1;
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
    padding: 6px 10px;
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

.bank-select-value {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
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

.opacity-slider {
    -webkit-appearance: none;
    appearance: none;
    flex: 1;
    height: 4px;
    background: linear-gradient(90deg, rgba(79, 195, 247, 0.3), rgba(79, 195, 247, 0.8));
    border-radius: 2px;
    outline: none;
    cursor: pointer;
}

.opacity-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    cursor: pointer;
    box-shadow: 0 2px 8px rgba(79, 195, 247, 0.4);
    transition: transform 0.2s;
}

.opacity-slider::-webkit-slider-thumb:hover {
    transform: scale(1.15);
}

.opacity-slider::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    cursor: pointer;
    border: none;
    box-shadow: 0 2px 8px rgba(79, 195, 247, 0.4);
}

.opacity-value {
    font-size: 13px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.8);
    min-width: 36px;
    text-align: right;
}
</style>