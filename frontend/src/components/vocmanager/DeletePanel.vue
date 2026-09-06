<template>
    <MainPopupPanel :visible="visible" title="⚠️ 确认删除" card-class="modal-card-sm" @close="$emit('close')">
        <p class="confirm-text">
            确定要删除单词 <strong>{{ wordName }}</strong> 吗？
        </p>
        <p class="confirm-warn">此操作不可恢复。</p>
        <template #footer>
            <button class="btn btn-secondary" @click="$emit('close')">取消</button>
            <button class="btn btn-danger" :disabled="deleting" @click="$emit('confirm')">
                {{ deleting ? '删除中...' : '确认删除' }}
            </button>
        </template>
    </MainPopupPanel>
</template>

<script setup>
import MainPopupPanel from '../../template/MainPopupPanel.vue'

defineProps({
    visible: { type: Boolean, default: false },
    wordName: { type: String, default: '' },
    deleting: { type: Boolean, default: false },
})

defineEmits(['close', 'confirm'])
</script>

<style scoped>
.confirm-text {
    margin: 0 0 8px;
    font-size: 14px;
    color: rgba(255, 255, 255, 0.8);
    line-height: 1.6;
}

.confirm-text strong {
    color: rgba(255, 255, 255, 0.95);
}

.confirm-warn {
    margin: 0;
    font-size: 12px;
    color: rgba(255, 255, 255, 0.4);
    line-height: 1.5;
}

.btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 7px 16px;
    border: none;
    border-radius: 6px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.15s;
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.7);
}

.btn-secondary:hover {
    background: rgba(255, 255, 255, 0.14);
}

.btn-danger {
    background: #e81123;
    color: #fff;
}

.btn-danger:hover:not(:disabled) {
    background: #c90e1d;
}

.btn-danger:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}
</style>