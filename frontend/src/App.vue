<template>
    <div class="vocab-card" :style="{ opacity: store.opacity }">
        <!-- 窗口边缘拖拽调整大小手柄 -->
        <div class="resize-handle top" @mousedown.prevent="startResize($event, 'top')"></div>
        <div class="resize-handle bottom" @mousedown.prevent="startResize($event, 'bottom')"></div>
        <div class="resize-handle left" @mousedown.prevent="startResize($event, 'left')"></div>
        <div class="resize-handle right" @mousedown.prevent="startResize($event, 'right')"></div>
        <div class="resize-handle top-left" @mousedown.prevent="startResize($event, 'top-left')"></div>
        <div class="resize-handle top-right" @mousedown.prevent="startResize($event, 'top-right')"></div>
        <div class="resize-handle bottom-left" @mousedown.prevent="startResize($event, 'bottom-left')"></div>
        <div class="resize-handle bottom-right" @mousedown.prevent="startResize($event, 'bottom-right')"></div>

        <ToolBar />
        <WordCard />
        <ActionBar />
        <SettingsPanel />
    </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue';
import { store } from './store.js';
import ToolBar from './components/ToolBar.vue';
import WordCard from './components/WordCard.vue';
import ActionBar from './components/ActionBar.vue';
import SettingsPanel from './components/SettingsPanel.vue';

// ===== 窗口拖拽调整大小 =====
let resizeState = null;

async function startResize(e, direction) {
    e.preventDefault();
    if (!window.runtime) return;
    let WindowPosition = await window.runtime.WindowGetPosition();
    // 保存初始状态
    resizeState = {
        direction,
        startX: e.screenX,
        startY: e.screenY,
        startWidth: window.innerWidth,
        startHeight: window.innerHeight,
        startScreenX: WindowPosition.x,
        startScreenY: WindowPosition.y,

    };
    console.log(window.screen)
    console.log('resize', direction, resizeState);
    console.log(window.runtime.WindowGetPosition());
    document.addEventListener('mousemove', onResizeMove);
    document.addEventListener('mouseup', onResizeEnd);
}

async function onResizeMove(e) {
    if (!resizeState) return;
    let ScreenWidth, ScreenHeight, ScreenCurrent;
    let Screen = await window.runtime.ScreenGetAll();
    Screen.forEach((item, index) => {
        if (item.isCurrent) {
            ScreenWidth = item.width;
            ScreenHeight = item.height;
            ScreenCurrent = index;
        }
    })
    let WindowPosition = await window.runtime.WindowGetPosition();
    const {
        direction,
        startX,
        startY,
        startWidth,
        startHeight,
        startScreenX,
        startScreenY
    } = resizeState;

    const dx = e.screenX - startX;
    const dy = e.screenY - startY;
    let newWidth = startWidth;
    let newHeight = startHeight;
    let newX = startScreenX;
    let newY = startScreenY;

    // console.log(resizeState);
    // console.log(direction, dx, dy);



    // 计算新尺寸
    if (direction.includes('right')) {
        newWidth = startWidth + dx;
    }
    if (direction.includes('left')) {
        newWidth = startWidth - dx;
        newX = startScreenX + dx / window.screen.width * ScreenWidth;
    }
    if (direction.includes('bottom')) {
        newHeight = startHeight + dy;
    }
    if (direction.includes('top')) {
        newHeight = startHeight - dy;
        newY = startScreenY + dy / window.screen.height * ScreenHeight;
    }
    // 保持最小/最大尺寸限制（与 main.go 一致）
    const minW = 300, minH = 400;


    if (newWidth < minW) {
        newWidth = minW;
        if (direction.includes('left')) newX = startScreenX + (startWidth - minW)/ window.screen.width * ScreenWidth;
    }

    if (newHeight < minH) {
        newHeight = minH;
        if (direction.includes('top')) newY = startScreenY + (startHeight - minH)/ window.screen.height * ScreenHeight;;
    }
    // 适配横屏双屏幕
    if (newX > Screen[0].width) {
        newX = newX - Screen[0].width;
    }else if(newX < 0){
        newX = newX + Screen[ScreenCurrent].width;
    }


    console.log(e.screenX, e.screenY)
    console.log(startX, startY)
    console.log(newX, newY, newWidth, newHeight)
    // 先设置位置（如果有变化），再设置大小
    if (direction.includes('left') || direction.includes('top')) {
        window.runtime.WindowSetPosition(Math.round(newX), Math.round(newY));
    }
    window.runtime.WindowSetSize(Math.round(newWidth), Math.round(newHeight));



}

function onResizeEnd() {
    resizeState = null;
    document.removeEventListener('mousemove', onResizeMove);
    document.removeEventListener('mouseup', onResizeEnd);
}

onMounted(async () => {
    await store.loadRandomWord();
});

onUnmounted(() => {
    // 清理事件监听
    document.removeEventListener('mousemove', onResizeMove);
    document.removeEventListener('mouseup', onResizeEnd);
});
</script>

<style>
/* ===== 全局样式 ===== */
* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

html,
body {
    width: 100%;
    height: 100%;
    overflow: hidden;
    background: transparent;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
        "PingFang SC", "Microsoft YaHei", sans-serif;
    user-select: none;
    -webkit-user-select: none;

}

#app {
    width: 100%;
    height: 100%;
}


/* ===== 主卡片 ===== */
.vocab-card {
    width: 100%;
    height: 100%;
    /* background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%); */
    background-color: #0000003e;

    display: flex;
    flex-direction: column;
    overflow: hidden;
    position: relative;
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

/* ===== 窗口拖拽调整大小手柄 ===== */
.resize-handle {
    position: absolute;
    z-index: 200;
}

/* 四条边 */
.resize-handle.top {
    top: 0;
    left: 8px;
    right: 8px;
    height: 4px;
    cursor: n-resize;
}

.resize-handle.bottom {
    bottom: 0;
    left: 8px;
    right: 8px;
    height: 4px;
    cursor: s-resize;
}

.resize-handle.left {
    left: 0;
    top: 8px;
    bottom: 8px;
    width: 4px;
    cursor: w-resize;
}

.resize-handle.right {
    right: 0;
    top: 8px;
    bottom: 8px;
    width: 4px;
    cursor: e-resize;
}

/* 四个角 */
.resize-handle.top-left {
    top: 0;
    left: 0;
    width: 8px;
    height: 8px;
    cursor: nw-resize;
}

.resize-handle.top-right {
    top: 0;
    right: 0;
    width: 8px;
    height: 8px;
    cursor: ne-resize;
}

.resize-handle.bottom-left {
    bottom: 0;
    left: 0;
    width: 8px;
    height: 8px;
    cursor: sw-resize;
}

.resize-handle.bottom-right {
    bottom: 0;
    right: 0;
    width: 8px;
    height: 8px;
    cursor: se-resize;
}
</style>
