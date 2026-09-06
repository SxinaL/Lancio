<template>
    <div class="resize-handle top" @mousedown.prevent="startResize($event, 'top')"></div>
    <div class="resize-handle bottom" @mousedown.prevent="startResize($event, 'bottom')"></div>
    <div class="resize-handle left" @mousedown.prevent="startResize($event, 'left')"></div>
    <div class="resize-handle right" @mousedown.prevent="startResize($event, 'right')"></div>
    <div class="resize-handle top-left" @mousedown.prevent="startResize($event, 'top-left')"></div>
    <div class="resize-handle top-right" @mousedown.prevent="startResize($event, 'top-right')"></div>
    <div class="resize-handle bottom-left" @mousedown.prevent="startResize($event, 'bottom-left')"></div>
    <div class="resize-handle bottom-right" @mousedown.prevent="startResize($event, 'bottom-right')"></div>
</template>

<script setup>
import { onUnmounted } from 'vue'
import { SaveFloatWindowPositionAndSize, GetMinSize } from '../../../wailsjs/go/handler/WindowHandler.js'

let resizeState = null

async function startResize(e, direction) {
    e.preventDefault()
    if (!window.runtime) return
    const WindowPosition = await window.runtime.WindowGetPosition()

    resizeState = {
        direction,
        startX: e.screenX,
        startY: e.screenY,
        startWidth: window.innerWidth,
        startHeight: window.innerHeight,
        startScreenX: WindowPosition.x,
        startScreenY: WindowPosition.y,
    }

    document.addEventListener('mousemove', onResizeMove)
    document.addEventListener('mouseup', onResizeEnd)
}

async function onResizeMove(e) {
    // 移动窗口许可
    let MoveFlag = true;
    if (!resizeState) return
    let ScreenWidth, ScreenHeight, ScreenCurrent
    const Screen = await window.runtime.ScreenGetAll()
    Screen.forEach((item, index) => {
        if (item.isCurrent) {
            ScreenWidth = item.width
            ScreenHeight = item.height
            ScreenCurrent = index
        }
    })

    const {
        direction,
        startX,
        startY,
        startWidth,
        startHeight,
        startScreenX,
        startScreenY,
    } = resizeState

    const dx = e.screenX - startX
    const dy = e.screenY - startY
    let newWidth = startWidth
    let newHeight = startHeight
    let newX = startScreenX
    let newY = startScreenY

    if (direction.includes('right')) {
        newWidth = startWidth + dx
    }
    if (direction.includes('left')) {
        newWidth = startWidth - dx
        newX = startScreenX + (dx / window.screen.width) * ScreenWidth
    }
    if (direction.includes('bottom')) {
        newHeight = startHeight + dy
    }
    if (direction.includes('top')) {
        newHeight = startHeight - dy
        newY = startScreenY + (dy / window.screen.height) * ScreenHeight
    }
    const minSize = await GetMinSize()
    const minW = minSize.WindowSizeMinW;
    const minH = minSize.WindowSizeMinH;


    if (newWidth < minW) {
        newWidth = minW
        MoveFlag = false;
        if (direction.includes('left'))
            newX = startScreenX + ((startWidth - minW) / window.screen.width) * ScreenWidth
    }

    if (newHeight < minH) {
        newHeight = minH

        if (direction.includes('top'))
            newY = startScreenY + ((startHeight - minH) / window.screen.height) * ScreenHeight
    }

    if (newX > Screen[0].width) {
        newX = newX - Screen[0].width
    } else if (newX < 0) {
        newX = newX + Screen[ScreenCurrent].width
    }

    if (direction.includes('left') || direction.includes('top')) {
        window.runtime.WindowSetPosition(Math.round(newX), Math.round(newY))
    }
    if (MoveFlag) {
        window.runtime.WindowSetSize(Math.round(newWidth), Math.round(newHeight))
    }
}

async function onResizeEnd() {
    resizeState = null
    document.removeEventListener('mousemove', onResizeMove)
    document.removeEventListener('mouseup', onResizeEnd)
    await SaveFloatWindowPositionAndSize()
}

onUnmounted(() => {
    document.removeEventListener('mousemove', onResizeMove)
    document.removeEventListener('mouseup', onResizeEnd)
})
</script>

<style scoped>
.resize-handle {
    position: absolute;
    z-index: 200;
}

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