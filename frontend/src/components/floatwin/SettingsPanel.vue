<template>
    <Teleport to="body">
        <Transition name="modal">
            <div v-if="store.showSettings" class="settings-overlay" @click="store.closeSettings($event)">
                <div class="settings-panel" @click.stop>
                    <div class="settings-header">
                        <span class="settings-title">⚙️ 设置</span>
                        <button class="tool-btn" @click="store.closeSettings()" title="关闭">✕</button>
                    </div>
                    <div class="settings-body">
                        <div class="setting-item">
                            <div class="setting-label">
                                <span class="setting-icon">🔆</span>
                                <span>窗口透明度</span>
                            </div>
                            <div class="setting-control">
                                <input
                                    type="range"
                                    class="opacity-slider"
                                    min="10"
                                    max="100"
                                    :value="Math.round(settingStore.opacity * 100)"
                                    @input="settingStore.setOpacity($event.target.value)"
                                />
                                <span class="opacity-value">{{ Math.round(settingStore.opacity * 100) }}%</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </Transition>
    </Teleport>
</template>

<script setup>
import { settingStore,store } from '@/store.js';
</script>

<style scoped>
.settings-overlay {
    position: fixed;
    top: 0; left: 0; right: 0; bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    backdrop-filter: blur(4px);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 100;
    border-radius: 12px;
}
.settings-panel {
    background: linear-gradient(135deg, #1e2a3a 0%, #1a1a2e 100%);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    width: 85%;
    max-width: 300px;
    box-shadow: 0 12px 40px rgba(0, 0, 0, 0.5);
    overflow: hidden;
}
.settings-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}
.settings-title {
    font-size: 14px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
}
.tool-btn {
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
.tool-btn:hover {
    background: rgba(255, 255, 255, 0.1);
    color: rgba(255, 255, 255, 0.9);
}
.settings-body {
    padding: 16px;
}
.setting-item {
    margin-bottom: 4px;
}
.setting-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 13px;
    color: rgba(255, 255, 255, 0.7);
    margin-bottom: 10px;
}
.setting-icon {
    font-size: 16px;
}
.setting-control {
    display: flex;
    align-items: center;
    gap: 12px;
    padding: 0 4px;
}
.opacity-slider {
    -webkit-appearance: none;
    appearance: none;
    flex: 1;
    height: 4px;
    background: linear-gradient(90deg, rgba(79, 195, 247, 0.3), rgba(79, 195, 247, 0.8));
    border-radius: 2px;
    outline: none;
    cursor: pointer;
}
.opacity-slider::-webkit-slider-thumb {
    -webkit-appearance: none;
    appearance: none;
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    cursor: pointer;
    box-shadow: 0 2px 8px rgba(79, 195, 247, 0.4);
    transition: transform 0.2s;
}
.opacity-slider::-webkit-slider-thumb:hover {
    transform: scale(1.15);
}
.opacity-slider::-moz-range-thumb {
    width: 18px;
    height: 18px;
    border-radius: 50%;
    background: linear-gradient(135deg, #4fc3f7, #29b6f6);
    cursor: pointer;
    border: none;
    box-shadow: 0 2px 8px rgba(79, 195, 247, 0.4);
}
.opacity-value {
    font-size: 13px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.8);
    min-width: 36px;
    text-align: right;
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
.modal-enter-active .settings-panel,
.modal-leave-active .settings-panel {
    transition: transform 0.25s ease;
}
.modal-enter-from .settings-panel {
    transform: translateY(20px);
}
.modal-leave-to .settings-panel {
    transform: translateY(20px);
}
</style>
