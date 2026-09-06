<template>
    <!--  :style="{ backgroundColor: rabg(0,0,0,1)}" -->
    <div class="vocab-card" :style="{ backgroundColor: `rgba(0,0,0,${VocFloatWinStore.opacity})` }">
        <ResizeHandles />

        <ToolBar @togglePin="togglePin" />
        <WordCard />
        <ActionBar />
        <SettingsPanel />
        <StarPannel />
    </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue';
import { VocFloatWinStore, WordStore } from '@/store.js';
import ToolBar from '@/components/vocfloatwin/ToolBar.vue';
import WordCard from '@/components/vocfloatwin/WordCard.vue';
import ActionBar from '@/components/vocfloatwin/ActionBar.vue';
import SettingsPanel from '@/components/vocfloatwin/SettingsPanel.vue';
import StarPannel from '@/components/vocfloatwin/StarPannel.vue';
import ResizeHandles from '@/components/vocfloatwin/ResizeHandles.vue';

// ===== 自动轮换单词定时器 =====
let wordTimer = null
const WORD_INTERVAL = 30000 // 30秒换一个单词

// 处理划词结果：接收后端广播的 Vocabulary 对象
// { word, phonetic, translation, example_sentence }
function onTextSelected(vocab) {
    WordStore.setWordFromText(vocab);
}
// 切换固定模式
async function togglePin() {
    VocFloatWinStore.togglePin()
    if (!VocFloatWinStore.isPinned) {
        await WordStore.loadRandomWord();
        wordTimer = setInterval(() => {
            WordStore.loadRandomWord()
        }, WORD_INTERVAL)
    } else {
        console.log('停止定时器')
        clearInterval(wordTimer)
        wordTimer = null
        WordStore.clearWord()
    }
}
onMounted(async () => {

    // 启动定时器，每 30 秒自动换一个单词
    if (!VocFloatWinStore.isPinned) {
        await WordStore.loadRandomWord();
        wordTimer = setInterval(() => {
            WordStore.loadRandomWord()
        }, WORD_INTERVAL)
    }
    await VocFloatWinStore.init()

    window.runtime?.EventsOn?.('text:selected', onTextSelected)
    window.runtime?.EventsOn?.('pin:changed', togglePin)

});

onUnmounted(() => {
    // 清理单词定时器
    if (wordTimer) {
        clearInterval(wordTimer)
        wordTimer = null
    }
    // 注销划词事件监听
    window.runtime?.EventsOff?.('text:selected')
});
</script>

<style>
.vocab-card {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    /* overflow: hidden; */
    position: relative;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
    /* border: 1px solid rgb(255, 0, 0); */
}
</style>