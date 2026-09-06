<template>
    <div class="word-content" @click="store.toggleTranslation()">
        <template v-if="store.currentWord">
            <div class="word-main">
                <div class="word-text">{{ store.currentWord.word }}</div>
                <div class="phonetic">{{ store.currentWord.phonetic || '' }}</div>
            </div>
            <div class="translation-section" :class="{ visible: store.showTranslation }">
                <div class="pos-tag">{{ store.currentWord.part_of_speech || '' }}</div>
                <div class="translation">{{ store.currentWord.translation }}</div>
                <div v-if="store.currentWord.example_sentence" class="example">
                    <span class="example-label">例句:</span>
                    {{ store.currentWord.example_sentence }}
                </div>
            </div>
            <div class="hint-text">
                {{ store.showTranslation ? '点击单词收起释义' : '点击单词查看释义' }}
            </div>
        </template>
        <div v-else class="empty-state">
            <div class="empty-icon">📚</div>
            <div>词库为空，请添加单词</div>
        </div>
    </div>
</template>

<script setup>
import { store } from '@/store.js';
</script>

<style scoped>
.word-content {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 16px 20px;
    cursor: pointer;
}
.word-main {
    text-align: center;
    transition: transform 0.2s;
}
.word-main:hover {
    transform: scale(1.02);
}
.word-text {
    font-size: 28px;
    font-weight: 700;
    color: #ffffff;
    letter-spacing: 1px;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
}
.phonetic {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.5);
    margin-top: 6px;
    font-style: italic;
    letter-spacing: 0.5px;
}
.translation-section {
    text-align: center;
    margin-top: 12px;
    opacity: 0;
    max-height: 0;
    overflow: hidden;
    transition: all 0.4s ease;
}
.translation-section.visible {
    opacity: 1;
    max-height: 200px;
}
.pos-tag {
    display: inline-block;
    background: rgba(79, 195, 247, 0.15);
    color: #4fc3f7;
    font-size: 11px;
    padding: 2px 8px;
    border-radius: 4px;
    margin-bottom: 6px;
    font-weight: 500;
}
.translation {
    font-size: 20px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
    margin: 4px 0;
}
.example {
    font-size: 12px;
    color: rgba(255, 255, 255, 0.45);
    margin-top: 8px;
    font-style: italic;
    line-height: 1.5;
    max-width: 280px;
}
.example-label {
    color: rgba(255, 255, 255, 0.3);
    font-style: normal;
}
.hint-text {
    font-size: 11px;
    color: rgba(255, 255, 255, 0.2);
    margin-top: 10px;
    transition: opacity 0.2s;
}
.empty-state {
    text-align: center;
    color: rgba(255, 255, 255, 0.4);
}
.empty-icon {
    font-size: 40px;
    margin-bottom: 8px;
}
</style>
