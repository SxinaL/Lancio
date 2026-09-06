<template>
    <Transition name="toast">
        <div v-if="message" class="toast" :class="type">
            {{ message }}
        </div>
    </Transition>
</template>

<script setup>
defineProps({
    message: { type: String, default: '' },
    type: { type: String, default: 'info' },
})
</script>

<style scoped>
.toast {
    position: fixed;
    bottom: 24px;
    left: 50%;
    transform: translateX(-50%);
    padding: 8px 16px;
    border-radius: 6px;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.9);
    background: rgba(30, 42, 58, 0.95);
    border: 1px solid rgba(255, 255, 255, 0.08);
    z-index: 1000;
}
.toast.success {
    border-color: rgba(79, 195, 247, 0.5);
}
.toast.error {
    border-color: rgba(232, 17, 35, 0.5);
}

.toast-enter-active {
    animation: toast-in 0.2s ease;
}
.toast-leave-active {
    animation: toast-in 0.2s ease reverse;
}
@keyframes toast-in {
    from { opacity: 0; transform: translate(-50%, 8px); }
    to   { opacity: 1; transform: translate(-50%, 0); }
}
</style>