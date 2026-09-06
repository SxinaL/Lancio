<template>
    <nav class="side-bar">
        <div class="side-nav">
            <RouterLink v-for="item in navItems" :key="item.id" class="nav-item"
                :class="{ active: activeNav === item.id }" @click="activeNav = item.id" :title="item.label" :to="item.link">
                <span class="nav-icon">{{ item.icon }}</span>
                {{ item.label }}
            </RouterLink>
        </div>

        <div class="side-bottom">
            <div class="nav-item" @click="onSettings" title="设置">
                <span class="nav-icon">⚙️</span>
                <span class="nav-label">设置</span>
            </div>
        </div>
    </nav>
</template>

<script setup>
import { ref } from 'vue'

const activeNav = ref('words')

const navItems = [
    { id: 'words', icon: '📝', label: '单词管理', link: '/main/bank-manager' },
    { id: 'review', icon: '🔄', label: '复习计划', link: '' },
    { id: 'stats', icon: '📊', label: '学习统计', link: '' },
    { id: 'float', icon: '🪟', label: '浮窗控制', link: '/main/float-manager' },
]

function onSettings() {
}
</script>

<style scoped>
.side-bar {
    width: 180px;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    background: rgba(0, 0, 0, 0.2);
    border-right: 1px solid rgba(255, 255, 255, 0.05);
    flex-shrink: 0;
    padding: 8px 0;
}

.side-nav {
    display: flex;
    flex-direction: column;
    gap: 2px;
    padding: 0 8px;
}

.side-bottom {
    padding: 0 8px;
    border-top: 1px solid rgba(255, 255, 255, 0.05);
    padding-top: 4px;
}

.nav-item {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 9px 12px;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.15s;
    color: rgba(255, 255, 255, 0.55);
    font-size: 16px;
    user-select: none;
    text-decoration: none;
}

.nav-item:hover {
    background: rgba(255, 255, 255, 0.05);
    color: rgba(255, 255, 255, 0.8);
}

.nav-item.active {
    background: rgba(79, 195, 247, 0.12);
    color: #4fc3f7;
}

.nav-icon {
    font-size: 16px;
    width: 22px;
    text-align: center;
    flex-shrink: 0;
}

.nav-label {
    white-space: nowrap;
}
</style>