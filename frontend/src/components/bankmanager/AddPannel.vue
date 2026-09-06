<template>
    <MainPopupPanel :visible="visible" :title="isEditing ? '修改词库名称' : '添加新词库'" @close="$emit('close')">
        <label class="input-label">词库名称</label>
        <input
            ref="inputRef"
            v-model="inputName"
            class="text-input"
            placeholder="输入词库名称，如「四级词汇」「GRE核心」"
            maxlength="30"
            @keyup.enter="$emit('confirm', inputName.trim())"
        />
        <template #footer>
            <button class="btn btn-secondary" @click="$emit('close')">取消</button>
            <button class="btn btn-primary" :disabled="!inputName.trim() || submitting" @click="$emit('confirm', inputName.trim())">
                {{ submitting ? '创建中...' : (isEditing ? '保存' : '添加') }}
            </button>
        </template>
    </MainPopupPanel>
</template>

<script setup>
import { ref, nextTick, watch } from 'vue'
import MainPopupPanel from '../../template/MainPopupPanel.vue'

const props = defineProps({
    visible: { type: Boolean, default: false },
    isEditing: { type: Boolean, default: false },
    initialName: { type: String, default: '' },
    submitting: { type: Boolean, default: false },
})

defineEmits(['close', 'confirm'])

const inputName = ref('')
const inputRef = ref(null)

watch(() => props.visible, (val) => {
    if (val) {
        inputName.value = props.initialName
        nextTick(() => inputRef.value?.focus())
    }
})
</script>

<style scoped>
.input-label {
    display: block;
    font-size: 12px;
    font-weight: 500;
    color: rgba(255, 255, 255, 0.5);
    margin-bottom: 8px;
}

.text-input {
    width: 100%;
    padding: 10px 14px;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 6px;
    color: rgba(255, 255, 255, 0.9);
    font-size: 14px;
    outline: none;
    transition: border-color 0.15s;
    box-sizing: border-box;
}

.text-input:focus {
    border-color: rgba(79, 195, 247, 0.5);
    box-shadow: 0 0 0 3px rgba(79, 195, 247, 0.1);
}

.text-input::placeholder {
    color: rgba(255, 255, 255, 0.25);
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

.btn-primary {
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    color: #fff;
    box-shadow: 0 2px 8px rgba(79, 195, 247, 0.3);
}

.btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 4px 12px rgba(79, 195, 247, 0.4);
}

.btn-primary:disabled {
    opacity: 0.4;
    cursor: not-allowed;
}

.btn-secondary {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.7);
}

.btn-secondary:hover {
    background: rgba(255, 255, 255, 0.14);
}
</style>