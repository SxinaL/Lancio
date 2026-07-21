<template>
    <div class="action-bar">
        <button class="action-btn primary" @click="store.nextWord()">
            <span class="btn-icon">🎲</span>
            下一个
        </button>
    </div>
</template>

<script setup>
import { store } from '../store.js';
</script>

<style scoped>
.action-bar {
    display: flex;
    justify-content: center;
    padding: 10px 16px 14px;
    gap: 10px;
    border-top: 1px solid rgba(255, 255, 255, 0.06);
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
</style>
