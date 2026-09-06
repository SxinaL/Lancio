package utils

import (
	"Lancio/model"
	"syscall"
	"github.com/gonutz/w32/v2"
	"strings"
	"time"
)

var (
	kernel32               = syscall.NewLazyDLL("kernel32.dll")
	procGetCurrentThreadId = kernel32.NewProc("GetCurrentThreadId")
	procPostThreadMessage  = user32.NewProc("PostThreadMessageW")

	// hookThreadID 钩子线程 ID，用于退出时发 WM_QUIT
	hookThreadID uint32
)

// getCurrentThreadId 获取当前 Windows 线程 ID
func getCurrentThreadId() uint32 {
	ret, _, _ := procGetCurrentThreadId.Call()
	return uint32(ret)
}

// postThreadMessage 向指定线程投递消息（不阻塞）
func postThreadMessage(threadID uint32, msg uint32, wParam, lParam uintptr) bool {
	ret, _, _ := procPostThreadMessage.Call(uintptr(threadID), uintptr(msg), wParam, lParam)
	return ret != 0
}

// StopHook 通知钩子线程退出消息循环（发送 WM_QUIT）
// 调用后钩子线程会从 GetMessage 返回 0，触发 UnhookWindowsHookEx
func StopHook() {
	if hookThreadID != 0 {
		postThreadMessage(hookThreadID, w32.WM_QUIT, 0, 0)
	}
}

// InitHookThreadID 由钩子 goroutine 在锁定 OS 线程后调用，
// 记录当前线程 ID 供 StopHook 使用
func InitHookThreadID() {
	hookThreadID = getCurrentThreadId()
}

// mouseDownAt 记录按住左键的时间，用于区分「长按划词」与「普通点击」
var mouseDownAt time.Time

// mouseHookCallback 鼠标钩子回调：Ctrl + 鼠标左键划词。
// 在「左键按下」时检测 Ctrl 并记录，在「左键释放」时触发取词，
// 避免 Ctrl（菜单键）激活菜单模式导致释放消息延迟、触发变慢。
// 为避免普通点击被误判为划词，仅当按住时间超过 1s 才触发取词。
func MouseHookCallback(code int, wParam w32.WPARAM, lParam w32.LPARAM) w32.LRESULT {
	if code >= 0 {
		switch wParam {
		case w32.WM_LBUTTONDOWN:
			// 按下左键时 Ctrl 已按住 → 标记本次为 Ctrl 划词，并记录按下时间
			
			if w32.GetAsyncKeyState(w32.VK_CONTROL)&0x8000 != 0 {
				model.CtrlSelecting = true
				
				mouseDownAt = time.Now()
			}
		case w32.WM_LBUTTONUP:
			if model.CtrlSelecting {
				model.CtrlSelecting = false
				// 仅当按住时间超过 1s 才判定为划词，避免快速点击误触发
				if time.Since(mouseDownAt) >= 500*time.Millisecond {
					// 通道异步消息传递使用
				
					select {
					case model.SelectionChan <- struct{}{}:
					default:
					}
				}
			}
		}
	}
	return w32.CallNextHookEx(0, code, wParam, lParam)
}

// selectionWorker 处理划词信号：保存剪贴板 -> 模拟 Ctrl+C -> 读取选中文本 -> 恢复剪贴板 -> 投递给翻译队列
func SelectionWorker() {
	for range model.SelectionChan {
		// 稍作延迟，确保鼠标已释放、选区稳定
		time.Sleep(50 * time.Millisecond)

		// 保存原剪贴板文本（仅文本，便于恢复）
		// OldClip := ReadClipboardText()
		
		// 模拟 Ctrl+C，让前台应用把选中内容复制到剪贴板
		SimulateCtrlC()

		// 等待目标应用完成复制
		time.Sleep(200 * time.Millisecond)
		// NewClip := ReadClipboardText()
	
		
		// 读取划词结果
		text := ReadClipboardText()

		// 恢复原剪贴板文本（原剪贴板是文本时才恢复）
		// if OldClip != "" {
		// 	_ = writeClipboardText(OldClip)
		// }

		text = strings.TrimSpace(text)
	
		if text == "" {continue}
		// 非阻塞投递给翻译队列，满了就丢弃（说明上一次还在处理）
		select {
		case model.TextSelection <- text:
		default:
		}
	}
}
