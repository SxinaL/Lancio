<template>
    <!-- @mousedown="onToolbarMouseDown" -->
    <div class="toolbar" >
        <span class="app-title">DesktopVoc </span>
        <div class="toolbar-actions">
            <button class="tool-btn" @click="store.openSettings()" title="设置">
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24"><path fill="#9094a2" d="m9.25 22l-.4-3.2q-.325-.125-.612-.3t-.563-.375L4.7 19.375l-2.75-4.75l2.575-1.95Q4.5 12.5 4.5 12.338v-.675q0-.163.025-.338L1.95 9.375l2.75-4.75l2.975 1.25q.275-.2.575-.375t.6-.3l.4-3.2h5.5l.4 3.2q.325.125.613.3t.562.375l2.975-1.25l2.75 4.75l-2.575 1.95q.025.175.025.338v.674q0 .163-.05.338l2.575 1.95l-2.75 4.75l-2.95-1.25q-.275.2-.575.375t-.6.3l-.4 3.2zm2.8-6.5q1.45 0 2.475-1.025T15.55 12t-1.025-2.475T12.05 8.5q-1.475 0-2.488 1.025T8.55 12t1.013 2.475T12.05 15.5"/></svg>
            </button>
            <button
                class="tool-btn"
                :class="{ active: store.isPinned }"
                @click="store.togglePin()"
                title="窗口置顶"
            >
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" v-if="store.isPinned"><path fill="#9094a2" d="M18 3v2h-1v6l2 3v2h-6v7h-2v-7H5v-2l2-3V5H6V3z"/></svg>
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" v-else><path fill="#9094a2" d="m22.313 10.175l-1.415 1.414l-.707-.707l-4.242 4.243l-.707 3.536l-1.415 1.414l-4.242-4.243l-4.95 4.95l-1.414-1.414l4.95-4.95l-4.243-4.243l1.414-1.414l3.536-.707l4.242-4.243l-.707-.707l1.414-1.414z"/></svg>
            </button>
            <button class="tool-btn" @click="minimizeWindow()" title="最小化">
                ─
            </button>
            <button class="tool-btn" @click="store.quitApp()" title="退出">
                ✕
            </button>
        </div>
    </div>
</template>

<script setup>
import { store } from '../store.js';

let isDragging = false;
let startX, startY;

// async function onToolbarMouseDown(e) {
//     // 按钮点击不触发拖拽
//     if (e.target.closest('.tool-btn')) return;

//     isDragging = true;
//     startX = e.clientX;
//     startY = e.clientY;
    
//     document.addEventListener('mousemove', onMouseMove);
//     document.addEventListener('mouseup', onMouseUp);
// }

// async function onMouseMove(e) {
//     if (!isDragging) return;
//     let windowPosition = await window.runtime.WindowGetPosition();
//     let PositionX = windowPosition.x;
//     let PositionY = windowPosition.y;
    
//     let dx = e.clientX - startX ;
//     let dy = e.clientY - startY ;

//     if (window.runtime?.WindowSetPosition) {
//         window.runtime.WindowSetPosition(
//             Math.round(PositionX + dx),
//             Math.round(PositionY + dy)
//         );
//     }
// }

// function onMouseUp() {
//     isDragging = false;
//     document.removeEventListener('mousemove', onMouseMove);
//     document.removeEventListener('mouseup', onMouseUp);
// }

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
