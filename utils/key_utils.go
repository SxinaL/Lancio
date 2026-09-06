package utils

import (
	"github.com/gonutz/w32/v2"
)

// procKeybdEvent 键盘事件系统调用
var procKeybdEvent = user32.NewProc("keybd_event")

// SimulateCtrlC 模拟按下 Ctrl+C（复制）
func SimulateCtrlC() {
	// Ctrl 按下
	procKeybdEvent.Call(w32.VK_CONTROL, 0, 0, 0)
	// C 按下
	procKeybdEvent.Call(uintptr('C'), 0, 0, 0)
	// C 松开
	procKeybdEvent.Call(uintptr('C'), 0, w32.KEYEVENTF_KEYUP, 0)
	// Ctrl 松开
	procKeybdEvent.Call(w32.VK_CONTROL, 0, w32.KEYEVENTF_KEYUP, 0)
}
