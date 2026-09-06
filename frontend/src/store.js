import { reactive, readonly } from 'vue';
import { GetRandomVocabulary, GetAllVocabularyBanks, GetSelectedVocabularyBankID, SetSelectedVocabularyBankID } from '../wailsjs/go/handler/VocabularyBankHandler';
import { AddVocabulary,GetStarStatus, UpdateVocabulary,DeleteVocabulary } from '../wailsjs/go/handler/VocabularyHandler';
import { SetWindowOpacity, GetWindowOpacity, PinFloatWindow, UnpinFloatWindow, GetPinStatus, SetPinStatus } from '../wailsjs/go/handler/WindowHandler';



export const VocFloatWinStore = reactive({

    opacity: await GetWindowOpacity(),
    showSettings: false,
    showTranslation: false,
    // 词库列表与当前选中的加载词库
    vocabularyBanks: [],
    selectedBankId: 0,
    // 浮窗是否置顶（置顶后脱离桌面层跳到最顶层）
    isPinned: await GetPinStatus(),
    async init() {
        await this.loadVocabularyBanks();
    },
    // 切换置顶：置顶脱离桌面层到最顶层，非置顶回到桌面层
    async togglePin() {
        this.isPinned = !this.isPinned;
        await SetPinStatus(this.isPinned);
        try {
            if (this.isPinned) {
                await PinFloatWindow();
            } else {
                await UnpinFloatWindow();
            }
        } catch (err) {
            console.error('切换置顶失败:', err);
            this.isPinned = !this.isPinned;
        }
    },

    // 加载词库列表与当前选中的词库 ID
    async loadVocabularyBanks() {
        try {
            this.vocabularyBanks = await GetAllVocabularyBanks();
            this.selectedBankId = await GetSelectedVocabularyBankID();
        } catch (err) {
            console.error('获取词库列表失败:', err);
        }
    },
    // 切换当前加载的词库
    async selectVocabularyBank(id) {
        this.selectedBankId = id;
        try {
            await SetSelectedVocabularyBankID(id);
            // 切换后立即从新词库加载一个单词
            await WordStore.loadRandomWord();
        } catch (err) {
            console.error('切换词库失败:', err);
        }
    },

    // 停止滑动（松开滑块）时调用：将最终透明度同步到后端
    async commitOpacity() {
        await SetWindowOpacity(this.opacity);
    },
    // 打开设置面板
    openSettings() {
        this.showSettings = true;
    },
    // 关闭设置面板
    closeSettings(event) {
        if (event && event.target !== event.currentTarget) return;
        this.showSettings = false;
    },
    // 滑动过程中实时预览：仅更新本地透明度值
    setOpacity(value) {
        this.opacity = parseInt(value) / 100;
    },
    // 切换显示翻译
    toggleTranslation() {
        this.showTranslation = !this.showTranslation;
    },

});

export const WordStore = reactive({
    currentWord: null,
    starredWord: null,
    isStarred: false, // 当前单词是否已收藏（仅在划词翻译模式下使用）
    // starredBankId: 0, // 当前单词收藏在哪个词库（0 = 未收藏）
    isWord: false,    // 当前文本是否为单个英文单词（仅在划词翻译模式下控制星星按钮显示）
    showStarPopup: false, // 是否显示收藏弹窗（仅在划词翻译模式下使用）
       

    async loadRandomWord() {
        try {
            this.currentWord = await GetRandomVocabulary();
            // 换单词时重置翻译显示状态（showTranslation 属于 VocFloatWinStore）
            VocFloatWinStore.showTranslation = false;
            // 换词模式下单词来自词库，默认即已收藏，不需要重置 isStarred
        } catch (err) {
            console.error('获取单词失败:', err);
            this.currentWord = null;
        }
    },
    async nextWord() {
        await this.loadRandomWord();
    },
    // 划词翻译模式下：根据后端返回的 Vocabulary 对象更新当前单词
    // vocab 字段：{ word, phonetic, translation, example_sentence }
    async setWordFromText(vocab) {
        if (!vocab || !vocab.word) return;
        this.currentWord = {
            word: vocab.word,
            phonetic: vocab.phonetic || '',
            translation: vocab.translation || '',
            example_sentence: vocab.example_sentence || '',
            vocabulary_bank_id: 0,
        };
        
        VocFloatWinStore.showTranslation = false;
        // 仅当划取的是单个英文单词时显示收藏按钮；句子则隐藏
        this.isWord = this.checkIsWord(vocab.word);
        if (this.isWord) {
            let word,starStatus
            try {
                word = await GetStarStatus(this.currentWord.word);
                starStatus = true;
            }catch(err){
                starStatus = false;
            }
            this.isStarred = starStatus;
            if (this.isStarred){
                this.starredWord = word;
            }else{
                this.starredWord = this.currentWord;
            }

        } else {
            this.isStarred = false;
        }
    },
    clearWord(){
        this.currentWord = null;
    },
    // 判断文本是否为单个英文单词（无空白、仅由字母/连字符/撇号组成）
    checkIsWord(text) {
        if (!text) return false;
        const trimmed = text.trim();
        if (!trimmed) return false;
        // 含空白字符大于3为句子
        const whitespaceCount = (trimmed.match(/\s/g) || []).length;
        if (whitespaceCount > 3) return false;
        // 仅允许字母、连字符、撇号
        return /^[A-Za-z]+([-' ][A-Za-z]+)*$/.test(trimmed);
    },
    // 检查给定单词文本是否已存在于任意词库中，结果同步到 isStarred / starredBankId



    // 点击星星：弹出词库选择弹窗
    toggleStarPopup() {
        this.showStarPopup = !this.showStarPopup;
    },
    // 关闭弹窗
    closeStarPopup() {
        this.showStarPopup = false;
    },
    // 收藏到指定词库
    async collectToBank(bankId) {
       
        this.showStarPopup = false;
        this.starredWord.vocabulary_bank_id = bankId;
        if (this.isStarred){
            // 更新收藏词库ID
            await UpdateVocabulary(this.starredWord);
        }else{
            // 新增收藏记录
            await AddVocabulary(this.starredWord);
        }
     
        // 刷新单词信息
        this.starredWord,this.isStarred = await GetStarStatus(this.starredWord.word);

    },
    // 取消收藏：从词库中移除当前单词
    async uncollect() {
        this.isStarred = false;
        this.showStarPopup = false;
        // TODO: 后续调用后端接口从词库中移除单词
        await DeleteVocabulary(this.starredWord.id);
    },

});