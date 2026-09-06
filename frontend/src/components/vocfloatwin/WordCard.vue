<template>
    <!-- 词典模式（划词翻译）：单词、音标、释义、例句都直接展示，点击切换全显/全隐 -->
    <div v-if="VocFloatWinStore.isPinned" class="word-content dict-mode" >
        <template v-if="WordStore.currentWord">
            <div class="word-main">
                <div class="word-text omit" :title="WordStore.currentWord.word">{{ WordStore.currentWord.word }}</div>
               
                <div v-if="WordStore.currentWord.phonetic" class="phonetic">
                    <span class="phonetic-label audio-btn" @click.stop="playAudio" title="播放发音">
                       <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" width="12" height="12" ><path fill="currentColor" d="M5 15V9h4l5-5v16l-5-5zm11 1V7.95q1.125.525 1.813 1.625T18.5 12t-.687 2.4T16 16"/></svg>
                       
                    </span>
                    <span class="phonetic-label">{{ WordStore.currentWord.phonetic }}</span>
                    
                    <audio ref="audioRef" :src="audioSrc" style="display:none"></audio>
                </div>
            </div>
             
            <div class="translation-section" :class="{ visible: true }">
                <div class="translation">{{ WordStore.currentWord.translation }}</div>
            </div>
      
        </template>
        <div v-else class="empty-state">
            <div class="empty-icon">🔍</div>
            <div>划词翻译</div>
        </div>
    </div>

    <!-- 词库模式（默认/换词模式）：点击切换释义显隐 -->
    <div v-else class="word-content lib-mode" @click="VocFloatWinStore.toggleTranslation()">
        <template v-if="WordStore.currentWord">
            <div class="word-main">
                <div class="word-text" :title="WordStore.currentWord.word">{{ WordStore.currentWord.word }}</div>
                <div class="phonetic">
                      <span class="phonetic-label audio-btn" @click.stop="playAudio" title="播放发音">
                       <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" width="12" height="12" ><path fill="currentColor" d="M5 15V9h4l5-5v16l-5-5zm11 1V7.95q1.125.525 1.813 1.625T18.5 12t-.687 2.4T16 16"/></svg>
                    </span>
                    <audio ref="audioRef" :src="audioSrc" style="display:none"></audio>
                    <span class="phonetic-label">{{ WordStore.currentWord.phonetic || '' }}</span>
                </div>
            </div>
            <div class="translation-section" :class="{ visible: VocFloatWinStore.showTranslation }" >
               
                <div class="translation"> {{ WordStore.currentWord.translation }}</div>
                <div v-if="WordStore.currentWord.example_sentence" class="example">
                    <span class="example-label">例句:</span>
                    {{ WordStore.currentWord.example_sentence }}
                </div>
            </div>
            <div class="hint-text">
                {{ VocFloatWinStore.showTranslation ? '点击单词收起释义' : '点击单词查看释义' }}
            </div>
        </template>
        <div v-else class="empty-state">
            <div class="empty-icon">📚</div>
            <div>词库为空，请添加单词</div>
        </div>
    </div>
</template>

<script setup>
import { computed, ref,  onUnmounted } from 'vue';
import { VocFloatWinStore, WordStore } from '@/store.js';
import { useResizeObserver } from '@vueuse/core';

// const translationLines = computed(() => {
//     if (!WordStore.currentWord?.translation) return ['暂无释义'];
//     return WordStore.currentWord.translation.split('\n');
// });
const audiourl = computed(() => {
    if (!WordStore.currentWord?.word) return '';
    let url = 'https://dict.youdao.com/dictvoice?audio='+WordStore.currentWord.word
    return url.replaceAll(' ', '%20');
});
const audioRef = ref(null);
const audioSrc = computed(() => audiourl.value);

function playAudio() {
    // console.log(audioSrc.value);
    if (audioRef.value) {
        audioRef.value.currentTime = 0;
        audioRef.value.play();
    }
}

// const transSectionRef = ref(null);
// const lineClamp = ref(3);

// function calcLineClamp() {
//     const translationHeight = transSectionRef.value.clientHeight-10;
//     if (!translationHeight) return;
//     lineClamp.value = Math.max(1, Math.floor(translationHeight / 20));
// }

// let resizeObserver = null;
// useResizeObserver(transSectionRef, calcLineClamp);

onUnmounted(() => {
    if (resizeObserver) {
        resizeObserver.disconnect();
        resizeObserver = null;
    }
});
</script>

<style scoped>
.word-content {
    display: flex;
    flex: 1;
    flex-direction: column;
    justify-content: flex-start;
    padding: 16px 20px;
    cursor: pointer;
    overflow: hidden;
    /* border: 1px solid rgb(255, 4, 4); */
}
.word-main {
    text-align: center;
    transition: transform 0.2s;
    /* border: 1px solid rgb(255, 0, 0); */
}
.word-text {
    text-align: center;
    font-size: 28px;
    font-weight: 700;
    color: #ffffff;
    letter-spacing: 1px;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.167);
    overflow: hidden;
}
.phonetic {
    font-size: 13px;
    color: rgba(255, 255, 255, 0.5);
    margin-top: 6px;
    font-style: italic;
    letter-spacing: 0.5px;
}
.phonetic-label{
    text-align: center;
    margin-right: 5px;
    font-size: 12px;
    /* border: 1px solid rgba(255, 255, 255, 0.2); */
}
.audio-btn:hover {
    color: #4fc3f7;
}
.translation-section {
    
    text-align: center;
    color: rgba(255, 255, 255, 0.9);
    margin-top: 12px;
    opacity: 0.4;
    overflow-x: hidden;
    overflow-y: scroll;
    transition: all 0.4s ease;
    scrollbar-width: none;
    /* border: 1px solid rgb(255, 1, 1); */
    /* padding: 5px; */
}

.translation-section::-webkit-scrollbar {
    width: 0;      /* 隐藏纵向滚动条宽度 */
    height: 0;     /* 隐藏横向滚动条高度 */
    display: none; /* 部分旧版浏览器兼容 */
}


.translation {
    font-size: 15px;
    font-weight: 600;
    white-space: pre-line;   /* ← 这一行是保留\n换行符的关键 */
    color: rgba(255, 255, 255, 0.9);
    overflow: hidden;
    /* border: 1px solid rgb(255, 0, 0); */
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
    font-size: 10px;
    color: rgba(255, 255, 255, 0.2);
    margin-top: 10px;
    transition: opacity 0.2s;
    text-align: center;
}
.empty-state {
    
    display: flex;
    flex: 1;
    flex-direction: column;
    justify-content: center;
    /* padding: 16px 20px; */
    /* border: 1px solid rgba(255, 0, 0, 0.2); */
    text-align: center;
    color: rgba(255, 255, 255, 0.4);
}
.empty-icon {
    font-size: 40px;
    margin-bottom: 8px;
}
.empty-translation {
    color: rgba(255, 255, 255, 0.4);
    font-size: 15px;
    font-style: italic;
}
/* ===== 词库模式专用 ===== */
.lib-mode .translation-section {
    max-height: 0;
}
.lib-mode .translation-section.visible {
    opacity: 1;
    max-height: 200px;
}
/* ===== 划词模式专用 ===== */

.dict-mode .translation-section {
    flex: 1 1 0;
    max-height: none;  
}

/* 超出部分省略号 */
.omit{
    display: -webkit-box;
    -webkit-box-orient: vertical;
    -webkit-line-clamp: var(--clamp, 1);
    line-clamp: var(--clamp, 1);
}

</style>