package services

import (
	"Lancio/database/entity"
	"Lancio/model"
	"Lancio/utils"
	"fmt"
	goruntime "runtime"

	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// StartTextSelectionHook 启动全局划词监听：按住 Ctrl 并用鼠标划词（松开左键）时，
// 自动获取前台应用中选中文本的原始内容，并通过 Wails 事件 "text:selected" 广播给前端。
// 该钩子需要独立消息循环，故在固定 OS 线程的 goroutine 中运行。
func StartTextSelectionHook() {
	app := model.AppInstance

	// 取词 worker：收到划词信号后异步执行复制取词，避免阻塞钩子线程
	go utils.SelectionWorker()

	// 翻译 worker：从 TextSelection channel 接收文本，判断长文本后调用对应 API
	go ReceiveTextSelection()

	// 全局低级鼠标钩子
	go func() {
		// 锁定 OS 线程
		goruntime.LockOSThread()
		// 记录钩子线程 ID，用于退出时发 WM_QUIT
		utils.InitHookThreadID()
		// 安装全局鼠标钩子（w32.SetWindowsHookEx 内部会 syscall.NewCallback）
		hHook := w32.SetWindowsHookEx(w32.WH_MOUSE_LL, utils.MouseHookCallback, 0, 0)
		if hHook == 0 {
			runtime.LogError(app.Ctx, "安装全局鼠标钩子失败")
			return
		}
		defer w32.UnhookWindowsHookEx(hHook)

		var message w32.MSG

		for {
			// 让线程活着，直到收到 WM_QUIT 或出错
			ret := w32.GetMessage(&message, 0, 0, 0)
			// ret == 0 表示收到 WM_QUIT；ret == -1 表示出错
			if ret == 0 || ret == -1 {
				break
			}
		}
	}()
}

// ReceiveTextSelection 接收划词文本：判断长短后走对应 API，统一广播 *entity.Vocabulary。
// - 短文本走 LookupWord（含三级缓存、回写），返回词典详情（音标、释义、例句）
// - 长文本走百度翻译，返回译文（音标、例句留空）
//
// 事件 "text:selected" 统一携带 Vocabulary 结构：
//
//	{ word, phonetic, translation, example_sentence }
func ReceiveTextSelection() {
	app := model.AppInstance
	for text := range model.TextSelection {
		if app == nil || app.Ctx == nil {
			continue
		}

		var vocab *entity.Vocabulary
		var err error

		if utils.IsLongText(text) {
			// 长文本：百度翻译，只填 word + translation
			translation, transErr := LookUpLongText(text)
			if transErr != nil {
				runtime.LogError(app.Ctx, fmt.Sprintf("百度翻译失败: %v", transErr))
				continue
			}
			vocab = &entity.Vocabulary{
				Word:        text,
				Translation: translation,
			}
		} else {
		
			// 判断是否单词
			if !utils.IsValidWord(text, "en") {
				continue
			}
			// 短文本：走完整查询（含内存/SQLite 缓存 + 在线查询）
			vocab, err = LookupWord(text, "en")
			if err != nil {
				runtime.LogError(app.Ctx, fmt.Sprintf("词典查询失败: %v", err))
				continue
			}
			// vocab 做个处理，它的翻译需要换行，每个词性只取2个意思
		}

		// 统一广播 Vocabulary 结构给前端
		runtime.EventsEmit(app.Ctx, "text:selected", vocab)
	}
}
// StopHook 停止全局鼠标钩子
func StopHook(){
	utils.StopHook()
}