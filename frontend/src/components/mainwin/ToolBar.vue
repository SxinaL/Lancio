<template>
    <MainToolBar>
        <template #title-left>
            <div class="title-left">
                <span class="app-logo"><svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" width="16"
                        height="16">
                        <path fill="#ffcc33"
                            d="M7.212 2c-.474 0-.891.314-1.021.77l-2.25 7.874a1.062 1.062 0 0 0 1.021 1.354H6.23l-1.17 4.68c-.264 1.055 1.04 1.777 1.795.995l3.343-3.413a1 1 0 0 0-.157-.066l-1.091-.355a1.424 1.424 0 0 1 0-2.684l.006-.002l.018-.006l1.061-.344a1.2 1.2 0 0 0 .746-.761l.355-1.093A1.42 1.42 0 0 1 12.48 8h.021a1.42 1.42 0 0 1 1.342.95l.002.006l.006.017l.348 1.068q.015.046.034.09l.003.008l1.297-1.324l.003-.004c.641-.667.18-1.811-.766-1.811h-2.564l1.261-3.594l.003-.008A1.062 1.062 0 0 0 12.461 2zm5.666 7.282l.348 1.071a2.2 2.2 0 0 0 1.398 1.397l1.072.348l.021.006a.423.423 0 0 1 0 .798l-1.071.348a2.2 2.2 0 0 0-1.399 1.397l-.348 1.07a.423.423 0 0 1-.798 0l-.349-1.07a2.2 2.2 0 0 0-.65-.977a2.2 2.2 0 0 0-.748-.426l-1.072-.348a.423.423 0 0 1 0-.798l1.072-.348a2.2 2.2 0 0 0 1.377-1.397l.348-1.07a.423.423 0 0 1 .799 0m4.905 7.931l-.766-.248a1.58 1.58 0 0 1-.998-.998l-.25-.765a.302.302 0 0 0-.57 0l-.248.765a1.58 1.58 0 0 1-.984.998l-.765.248a.302.302 0 0 0 0 .57l.765.249a1.58 1.58 0 0 1 1 1.002l.248.764a.302.302 0 0 0 .57 0l.249-.764a1.58 1.58 0 0 1 .999-.999l.765-.248a.302.302 0 0 0 0-.57z" />
                    </svg></span>
                <span class="app-name">Lancio</span>
                <span class="app-version">v1.0</span>
            </div>
        </template>
        <template #title-center>
            <div class="title-center">
            </div>
        </template>
        <template #title-right>
            <div class="title-right">
                <button class="title-btn" title="最小化" @click="onMinimize">
                    <svg width="12" height="12" viewBox="0 0 12 12">
                        <rect y="5" width="12" height="2" fill="currentColor" />
                    </svg>
                </button>

                <button v-if="!isMaximized" class="title-btn" title="最大化" @click="onMaximize">
                    <svg width="12" height="12" viewBox="0 0 12 12">
                        <rect x="1" y="1" width="10" height="10" stroke="currentColor" stroke-width="1.5" fill="none" />
                    </svg>
                </button>
                <button v-else class="title-btn" title="恢复窗口" @click="onMaximize">
                    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="12" height="12">
                        <path fill="currentColor" d="M4 20h12v-9H4zm14-5v-2h2V4H8v5H6V2h16v13zM2 22V9h16v13zm8-6.5" />
                    </svg>
                </button>
                <button class="title-btn title-btn-close" title="关闭" @click="onClose">
                    <svg width="12" height="12" viewBox="0 0 12 12">
                        <path d="M1 1l10 10M11 1L1 11" stroke="currentColor" stroke-width="1.5" />
                    </svg>
                </button>
            </div>
        </template>
    </MainToolBar>
</template>

<script setup>
let isMaximized = ref(false)
import { MinimizeToTray } from '../../../wailsjs/go/handler/TaskbarHandler'
import MainToolBar from '../../template/MainToolBar.vue'
import { ref } from 'vue'

function onMinimize() {
    // try {
    //     MinimizeToTray()
    // } catch (err) {
    // if (window.runtime?.WindowMinimise) {
    window.runtime.WindowMinimise()
    //     }
    // }
}
function onMaximize() {
    if (isMaximized.value) {
        window.runtime.WindowUnmaximise()
    } else {
        window.runtime.WindowMaximise()
    }
    isMaximized.value = !isMaximized.value
}
function onClose() {
    MinimizeToTray()
    // if (window.runtime?.Quit) {
    //     window.runtime.Quit()
    // }
}
</script>

<style scoped>
.title-left {
    display: flex;
    align-items: center;
    gap: 8px;
}

.app-logo {
    font-size: 18px;
}

.app-name {
    font-size: 14px;
    font-weight: 600;
    color: rgba(255, 255, 255, 0.9);
    letter-spacing: 0.5px;
}

.app-version {
    font-size: 10px;
    color: rgba(255, 255, 255, 0.3);
    padding: 1px 6px;
    background: rgba(255, 255, 255, 0.06);
    border-radius: 8px;
}

.title-center {
    flex: 1;
}

.title-right {
    display: flex;
    gap: 2px;
    --wails-draggable: none;
}

.title-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 28px;
    background: transparent;
    border: none;
    border-radius: 4px;
    color: rgba(255, 255, 255, 0.5);
    cursor: pointer;
    transition: all 0.15s;
}

.title-btn:hover {
    background: rgba(255, 255, 255, 0.08);
    color: rgba(255, 255, 255, 0.85);
}

.title-btn-close:hover {
    background: #e81123;
    color: #fff;
}
</style>