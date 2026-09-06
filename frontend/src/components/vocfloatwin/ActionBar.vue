<template>
    <div v-if="!VocFloatWinStore.isPinned" class="action-bar">
        <button class="action-btn primary" @click="WordStore.nextWord()">
            <!-- <span class="btn-icon"></span> -->
            下一个
        </button>
    </div>
    <div v-else class="action-bar pinned">
        <span class="pinned-text">提示：可根据划词内容调整窗口大小</span>
    </div>

</template>

<script setup>
import { WordStore,VocFloatWinStore} from '@/store.js';
</script>

<style scoped>

.action-bar {
    display: flex;
    
    justify-content: center;
    padding: 10px 16px 14px;
    gap: 10px;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
    flex: 0 0 auto;
    /* border: 2px solid rgb(255, 0, 0); */
}
.action-btn {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 8px 24px;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
}
.action-btn.primary {
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    color: #ffffff;
    box-shadow: 0 4px 12px rgba(79, 195, 247, 0.3);
}
.action-btn.primary:hover {
    transform: translateY(-1px);
    box-shadow: 0 6px 16px rgba(79, 195, 247, 0.4);
}
.action-btn.primary:active {
    transform: translateY(0);
}
.btn-icon {
    font-size: 16px;
}
/* 置顶状态提示条：与整体深色玻璃风格保持一致 */
.action-bar.pinned {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 10px 16px 14px;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
    user-select: none;
}
.pinned-text {
    font-size: 13px;
    font-weight: 500;
    letter-spacing: 0.5px;
    color: rgba(255, 255, 255, 0.55);
}
</style>
