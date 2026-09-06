<template>
    <Teleport to="body">
        <Transition name="modal">
            <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
                <div class="modal-panel" :style="{ maxWidth: maxWidth }" @click.stop>
                    <div class="modal-header">
                        <span class="modal-title">{{ title }}</span>
                        <button class="modal-close-btn" @click="$emit('close')" title="关闭">✕</button>
                    </div>
                    <div class="modal-body">
                        <slot />
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup>
defineProps({
    visible: { type: Boolean, default: false },
    title: { type: String, default: '' },
    maxWidth: { type: String, default: '300px' },
});

defineEmits(['close']);
</script>

<style scoped>
.modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
    border-radius: 12px;
}

.modal-panel {
    position: relative;
    background: linear-gradient(135deg, #1e2a3a 0%, #1a1a2e 100%);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    width: 85%;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.modal-title {
    font-size: 14px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
}

.modal-close-btn {
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

.modal-close-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.9);
}

.modal-body {
    padding: 16px;
}

/* Vue Transition */
.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}

.modal-enter-active .modal-panel,
.modal-leave-active .modal-panel {
    transition: transform 0.25s ease;
}

.modal-enter-from .modal-panel {
    transform: translateY(20px);
}

.modal-leave-to .modal-panel {
    transform: translateY(20px);
}
</style>