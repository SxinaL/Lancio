import { reactive } from 'vue';
import { GetRandomWord, SetOpacity } from '../wailsjs/go/main/App';

export const store = reactive({
    currentWord: null,
    showTranslation: false,
    isPinned: true,
    showSettings: false,
    opacity: 1.0,

    async loadRandomWord() {
        try {
            this.currentWord = await GetRandomWord();
            this.showTranslation = false;
        } catch (err) {
            console.error('获取单词失败:', err);
            this.currentWord = null;
        }
    },

    toggleTranslation() {
        this.showTranslation = !this.showTranslation;
    },

    async nextWord() {
        await this.loadRandomWord();
    },

    togglePin() {
        this.isPinned = !this.isPinned;
        if (window.runtime?.WindowSetAlwaysOnTop) {
            window.runtime.WindowSetAlwaysOnTop(this.isPinned);
        }
    },

    quitApp() {
        if (window.runtime?.Quit) {
            window.runtime.Quit();
        }
    },

    openSettings() {
        this.showSettings = true;
    },

    closeSettings(event) {
        if (event && event.target !== event.currentTarget) return;
        this.showSettings = false;
    },

    setOpacity(value) {
        const opacity = parseInt(value) / 100;
        this.opacity = opacity;
        try {
            SetOpacity(opacity);
        } catch (e) {
            // 静默处理
        }
    },
});
