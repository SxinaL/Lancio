package utils

import (
	"log"
	"syscall"

	"Lancio/party/systray"

	"Lancio/model"

	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	trayIcon    []byte // 托盘图标（ICO 格式），由 main 注入
	trayStarted bool
)

// SetTrayIcon 注入托盘图标字节（ICO 格式）
func SetTrayIcon(icon []byte) {
	trayIcon = icon
}

// StartTray 启动系统托盘（幂等，可在最小化到托盘时调用）
func StartTray(app *model.App) {
	if trayStarted {
		return
	}
	trayStarted = true

	go systray.Run(func() {
		systray.SetIcon(trayIcon)
		systray.SetTitle("Lancio")
		systray.SetTooltip("Lancio 桌面单词")

		// 左键点击托盘图标 → 直接恢复并显示主窗口（右键仍弹出菜单）
		systray.SetOnLeftClick(func() {
			go func() {
				defer func() {
					if r := recover(); r != nil {
						// 这行会直接写进 exe 旁边的 lancio.log 里
						log.Printf("恢复窗口时发生 panic: %v", r)
						// 如果还想看堆栈，可以加上 debug.Stack()
					}
				}()
				RestoreMainWindow()
			}()
		})

		mShow := systray.AddMenuItem("显示主窗口", "恢复并显示主窗口")
		mQuit := systray.AddMenuItem("退出", "退出 Lancio")

		go func() {
			for {
				select {
				case <-mShow.ClickedCh:
					go func() {
						defer func() {
							if r := recover(); r != nil {
								// 这行会直接写进 exe 旁边的 lancio.log 里
								log.Printf("恢复窗口时发生 panic: %v", r)
								// 如果还想看堆栈，可以加上 debug.Stack()
							}
						}()
						RestoreMainWindow()
					}()
				case <-mQuit.ClickedCh:
					systray.Quit()
					runtime.Quit(app.Ctx)
				}
			}
		}()
	}, func() {})
}

// HideFromTaskbar 通过扩展样式将窗口从任务栏隐藏（工具窗口样式）
func HideFromTaskbar(hwnd syscall.Handle) {
	exStyle := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE)
	newExStyle := (exStyle | w32.WS_EX_TOOLWINDOW) &^ w32.WS_EX_APPWINDOW
	w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE, newExStyle)
	w32.SetWindowPos(w32.HWND(hwnd), 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_FRAMECHANGED)
}

// ShowInTaskbar 恢复窗口在任务栏显示
func ShowInTaskbar(hwnd syscall.Handle) {
	exStyle := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE)
	// 移除工具窗口样式，恢复应用窗口样式，确保窗口重新出现在任务栏
	newExStyle := (exStyle &^ w32.WS_EX_TOOLWINDOW) | w32.WS_EX_APPWINDOW
	w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE, newExStyle)
	w32.SetWindowPos(w32.HWND(hwnd), 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_FRAMECHANGED)
}
