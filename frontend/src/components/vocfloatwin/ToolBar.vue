<template>
    <div @mousedown="startDrag">
        <FloatToolBar class="toolbar">
        <template #menu-title>
            <span class="app-title">LancioFloat</span>
        </template>
        <template #menu-actions>
            <div class="toolbar-actions">
                <!-- 收藏 -->
                <button
                    v-if="VocFloatWinStore.isPinned && WordStore.isWord"
                    class="tool-btn"
                    :class="{ active: WordStore.isStarred }"
                    @click="WordStore.toggleStarPopup()"
                    title="收藏"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="16" height="16">
                        <path fill="currentColor" d="m5.825 21l1.625-7.025L2 9.25l7.2-.625L12 2l2.8 6.625l7.2.625l-5.45 4.725L18.175 21L12 17.275z" />
                    </svg>
                </button>

                <!-- 切换划词模式 -->
                <button
                    class="tool-btn"
                    :class="{ active: VocFloatWinStore.isPinned }"
                    @click="$emit('togglePin')"
                    title="切换划词模式"
                >
                    <svg v-if="VocFloatWinStore.isPinned" xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                        <path fill="currentColor" d="M18 3v2h-1v6l2 3v2h-6v7h-2v-7H5v-2l2-3V5H6V3z" />
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                        <path fill="currentColor" d="m22.313 10.175l-1.415 1.414l-.707-.707l-4.242 4.243l-.707 3.536l-1.415 1.414l-4.242-4.243l-4.95 4.95l-1.414-1.414l4.95-4.95l-4.243-4.243l1.414-1.414l3.536-.707l4.242-4.243l-.707-.707l1.414-1.414z" />
                    </svg>
                </button>

                <!-- 设置 -->
                <button class="tool-btn" @click="VocFloatWinStore.openSettings()" title="设置">
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24">
                        <path fill="currentColor" d="m9.25 22l-.4-3.2q-.325-.125-.612-.3t-.563-.375L4.7 19.375l-2.75-4.75l2.575-1.95Q4.5 12.5 4.5 12.338v-.675q0-.163.025-.338L1.95 9.375l2.75-4.75l2.975 1.25q.275-.2.575-.375t.6-.3l.4-3.2h5.5l.4 3.2q.325.125.613.3t.562.375l2.975-1.25l2.75 4.75l-2.575 1.95q.025.175.025.338v.674q0 .163-.05.338l2.575 1.95l-2.75 4.75l-2.95-1.25q-.275.2-.575.375t-.6.3l-.4 3.2zm2.8-6.5q1.45 0 2.475-1.025T15.55 12t-1.025-2.475T12.05 8.5q-1.475 0-2.488 1.025T8.55 12t1.013 2.475T12.05 15.5" />
                    </svg>
                </button>
            </div>
        </template>
    </FloatToolBar>
    </div>
</template>

<script setup>
import { VocFloatWinStore, WordStore } from '@/store.js';
import FloatToolBar from '../../template/FloatToolBar.vue';

defineEmits(['togglePin']);

function startDrag(e) {
    if (e.target.closest('.tool-btn')) {
        e.preventDefault();
        e.stopPropagation();
        return;
    }
}
</script>

<style scoped>
.toolbar {
    display: flex;
    flex: 0 0 auto;
    /* border: 1px solid rgb(255, 0, 0); */
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    cursor: grab;
    user-select: none;
    --wails-draggable: drag;
}
.toolbar:active {
    cursor: grabbing;
}
.app-title {
    font-size: 13px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.8);
    letter-spacing: 0.5px;
}

.toolbar-actions {
    display: flex;
    gap: 2px;
}

.tool-btn {
    background: transparent;
    border: none;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    font-size: 14px;
    padding: 2px 6px;
    border-radius: 4px;
    transition: all 0.2s;
    line-height: 1;
}

.tool-btn:hover {
    color: rgba(255, 255, 255, 0.9);
}

.tool-btn.active {
    color: #4fc3f7;
}
</style>