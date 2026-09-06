<template>
    <!-- @mousedown="onToolbarMouseDown" -->
    <div class="toolbar" >
       
    </div>
</template>

<script setup>
import { store } from '@/store.js';

let isDragging = false;
let startX, startY;

function minimizeWindow() {
    if (window.runtime?.WindowMinimise) {
        window.runtime.WindowMinimise();
    }
}
</script>

<style scoped>
.toolbar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    /* background: rgba(255, 255, 255, 0.03);
    border-bottom: 1px solid rgba(255, 255, 255, 0.06); */
    cursor: grab;
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
    /* background: rgba(255, 255, 255, 0.1); */
    color: rgba(255, 255, 255, 0.9);
}
.tool-btn.active {
    color: #4fc3f7;
}

</style>
