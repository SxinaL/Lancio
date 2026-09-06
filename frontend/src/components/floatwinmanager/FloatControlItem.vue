<template>
    <div class="control-item">
        <div class="control-info">
            <div class="control-icon">
                <svg width="22" height="22" viewBox="0 0 24 24">
                    <path fill="currentColor" :d="icon" />
                </svg>
            </div>
            <div class="control-detail">
                <span class="control-name">{{ name }}</span>
                <span class="control-meta">{{ running ? '运行中' : '已关闭' }}</span>
            </div>
        </div>
        <div class="control-actions">
            <button
                class="toggle-switch"
                :class="{ active: running }"
                role="switch"
                :aria-checked="running"
                :title="running ? '点击关闭浮窗' : '点击开启浮窗'"
                @click="$emit('toggle')"
            >
                <span class="toggle-knob"></span>
            </button>
        </div>
    </div>
</template>

<script setup>
defineProps({
    name: { type: String, required: true },
    icon: { type: String, required: true },
    running: { type: Boolean, default: false },
})

defineEmits(['toggle'])
</script>

<style scoped>
.control-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px;
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.025);
    border: 1px solid rgba(255, 255, 255, 0.04);
    transition: background 0.15s, border-color 0.15s;
}
.control-item:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(255, 255, 255, 0.08);
}

.control-info {
    display: flex;
    align-items: center;
    gap: 12px;
}

.control-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 38px;
    height: 38px;
    background: rgba(79, 195, 247, 0.08);
    border-radius: 8px;
    color: #4fc3f7;
    flex-shrink: 0;
}

.control-detail {
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.control-name {
    font-size: 14px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.85);
}

.control-meta {
    font-size: 11px;
    color: rgba(255, 255, 255, 0.35);
}

.control-actions {
    display: flex;
    gap: 4px;
}

.toggle-switch {
    position: relative;
    width: 44px;
    height: 24px;
    background: rgba(255, 255, 255, 0.1);
    border: none;
    border-radius: 12px;
    cursor: pointer;
    transition: background 0.2s ease;
    padding: 0;
}
.toggle-switch.active {
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    box-shadow: 0 0 8px rgba(79, 195, 247, 0.4);
}

.toggle-knob {
    position: absolute;
    top: 2px;
    left: 2px;
    width: 20px;
    height: 20px;
    background: #fff;
    border-radius: 50%;
    transition: transform 0.2s ease;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
}
.toggle-switch.active .toggle-knob {
    transform: translateX(20px);
}
</style>