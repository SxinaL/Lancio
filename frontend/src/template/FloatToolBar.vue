<template>
    <!-- 窗口拖拽由父组件通过 @mousedown="startDrag" 监听实现（移除 --wails-draggable 避免冲突） -->
    <div>
        <slot name="menu-title"></slot>
        <slot name="menu-actions"></slot>
    </div>
</template>

<style scoped>




</style>
