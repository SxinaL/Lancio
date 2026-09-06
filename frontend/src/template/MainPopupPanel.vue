<template>
    <Teleport to="body">
        <Transition name="modal">
            <div v-if="visible" class="modal-overlay" @click.self="$emit('close')">
                <div class="modal-card" :class="cardClass">
                    <div class="modal-header">
                        <span>{{ title }}</span>
                        <button class="modal-close" @click="$emit('close')">
                            <svg width="14" height="14" viewBox="0 0 14 14"><path d="M1 1l12 12M13 1L1 13" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/></svg>
                        </button>
                    </div>
                    <div class="modal-body">
                        <slot />
                    </div>
                    <div class="modal-footer">
                        <slot name="footer" />
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
    cardClass: { type: String, default: '' },
})

defineEmits(['close'])
</script>

<style scoped>
.modal-overlay {
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(3px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
}

.modal-card {
    background: linear-gradient(135deg, #1e2a3a, #1a1a2e);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
    width: 420px;
    max-width: 90vw;
    box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
    overflow: hidden;
}

.modal-card-sm {
    width: 380px;
}

.modal-card-edit {
    width: 460px;
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px 20px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
    font-size: 15px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
}

.modal-close {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 26px;
    height: 26px;
    background: transparent;
    border: none;
    border-radius: 4px;
    color: rgba(255, 255, 255, 0.4);
    cursor: pointer;
    transition: all 0.15s;
}

.modal-close:hover {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.8);
}

.modal-body {
    padding: 20px;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    padding: 12px 20px 16px;
}

.modal-enter-active,
.modal-leave-active {
    transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
    opacity: 0;
}

.modal-enter-active .modal-card,
.modal-leave-active .modal-card {
    transition: transform 0.2s ease, opacity 0.2s ease;
}

.modal-enter-from .modal-card {
    transform: scale(0.95) translateY(8px);
    opacity: 0;
}

.modal-leave-to .modal-card {
    transform: scale(0.95) translateY(8px);
    opacity: 0;
}
</style>