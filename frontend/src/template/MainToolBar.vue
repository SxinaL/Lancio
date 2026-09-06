<template>
    <header class="title-bar">
        <slot name="title-left"></slot>
        <slot name="title-center"></slot>
        <slot name="title-right"></slot>
    </header>
</template>

<style scoped>
.title-bar {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 40px;
    padding: 0 12px;
    background: rgba(0, 0, 0, 0.35);
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    --wails-draggable: drag;
    flex-shrink: 0;
}
</style>