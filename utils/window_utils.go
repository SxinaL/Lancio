package utils

import (
	"Lancio/model"
	"fmt"
	"log"

	// "log"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	user32       = syscall.NewLazyDLL("user32.dll")
	findWindow   = user32.NewProc("FindWindowW")
	findWindowEx = user32.NewProc("FindWindowExW")
	setParent    = user32.NewProc("SetParent")
	isWindow     = user32.NewProc("IsWindow")
)

// 该子类化功能子类化功能，需要监听 系统给特定窗口发送的特定消息并处理，并转发回原窗口的处理流程
func ModifyWndProc() {
	if model.WindowInstance.WindowMode != "float" {
		return
	}
	// 获取要子类化的浮窗句柄
	hwnd, err := FindWindow("", model.WindowInstance.WindowName)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	originalWndProc := model.WindowInstance.OriginalWndProc
	// 子类化：安装新的窗口过程，拦截特定系统消息
	w32.SetWindowLongPtr(
		w32.HWND(hwnd),
		w32.GWLP_WNDPROC,
		syscall.NewCallback(func(h syscall.Handle, msg uint32, wParam, lParam uintptr) uintptr {
			switch msg {
			case w32.WM_EXITSIZEMOVE:
				SaveFloatWindowPositionAndSize()
			}
			// 其余消息转发回原窗口过程，保证系统默认行为
			return w32.CallWindowProc(originalWndProc, w32.HWND(h), msg, wParam, lParam)

		}),
	)
}

// FindWindow 查找窗口句柄
func FindWindow(className, windowName string) (hwnd syscall.Handle, err error) {
	var cname, wname *uint16
	if className != "" {
		cname, err = syscall.UTF16PtrFromString(className)
		if err != nil {
			return 0, err
		}
	}
	if windowName != "" {
		wname, err = syscall.UTF16PtrFromString(windowName)
		if err != nil {
			return 0, err
		}
	}
	r1, _, e1 := findWindow.Call(uintptr(unsafe.Pointer(cname)), uintptr(unsafe.Pointer(wname)))
	if r1 == 0 {
		// e1 为 syscall.Errno(0) 时接口非 nil 但实际无错误（只是没找到窗口）
		if e1 == nil || e1 == syscall.Errno(0) {
			err = fmt.Errorf("未找到窗口 (className=%q, windowName=%q)", className, windowName)
		} else {
			err = e1
		}
	}
	hwnd = syscall.Handle(r1)
	return hwnd, err
}

// FindChildWindow 查找子窗口句柄
func FindChildWindow(parentHandle syscall.Handle, className, windowName string) (hwnd syscall.Handle, err error) {
	var cname, wname *uint16
	if className != "" {
		cname, err = syscall.UTF16PtrFromString(className)
		if err != nil {
			return 0, err
		}
	}
	if windowName != "" {
		wname, err = syscall.UTF16PtrFromString(windowName)
		if err != nil {
			return 0, err
		}
	}
	r1, _, e1 := findWindowEx.Call(uintptr(parentHandle), 0, uintptr(unsafe.Pointer(cname)), uintptr(unsafe.Pointer(wname)))
	if r1 == 0 {
		// e1 为 syscall.Errno(0) 时接口非 nil 但实际无错误（只是没找到窗口）
		if e1 == nil || e1 == syscall.Errno(0) {
			err = fmt.Errorf("未找到窗口 (className=%q, windowName=%q)", className, windowName)
		} else {
			err = e1
		}
	}
	hwnd = syscall.Handle(r1)
	return hwnd, err
}

// modifyWindowStyle 修改窗口样式
// 移除标题栏、系统菜单、最小化/最大化、可调整边框
// 移除对话框边框（扩展样式）
func ModifyWindowStyle() {
	ctx := model.AppInstance.Ctx
	hwnd, err := FindWindow("", model.WindowInstance.WindowName)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	// 移除标题栏、系统菜单、最小化/最大化、可调整边框
	style := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_STYLE)
	newStyle := style &^ (w32.WS_CAPTION | w32.WS_SYSMENU | w32.WS_MINIMIZEBOX | w32.WS_MAXIMIZEBOX | w32.WS_THICKFRAME)

	w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_STYLE, newStyle)

	// 移除对话框边框（扩展样式）
	exStyle := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE)
	newExStyle := exStyle &^ w32.WS_EX_DLGMODALFRAME // WS_EX_DLGMODALFRAME移除双边框

	if model.WindowInstance.WindowMode == "float" {
		newExStyle = (newExStyle | w32.WS_EX_TOOLWINDOW | w32.WS_EX_NOACTIVATE) &^ w32.WS_EX_APPWINDOW
	} // 浮窗额外移除任务栏, 并设置为工具窗口, 不点击窗口时不激活、不抢焦点

	w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE, newExStyle)

	// 刷新边框
	w32.SetWindowPos(w32.HWND(hwnd), 0, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_FRAMECHANGED)

	// 显示窗口

	if model.AppConfig.StartHidden && model.WindowInstance.WindowMode == "main" {
		StartTray(model.AppInstance)
	} else {
		runtime.WindowShow(ctx)
	}

}

// SpawnFloatWindow 启动浮窗子进程
func SpawnFloatWindow(windowName string) {
	app := model.AppInstance
	exe, err := os.Executable()
	if err != nil {
		runtime.LogError(app.Ctx, fmt.Sprintf("获取可执行文件路径失败: %v", err))
		return
	}

	cmd := exec.Command(exe, "--float")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// 浮窗子进程使用独立的热重载端口，避免与主进程 -devserver 端口冲突
	// （devserver 环境变量优先级最高，覆盖继承自主进程的端口值）
	if devServer := pickFreePort(); devServer != "" {
		cmd.Env = append(os.Environ(), "devserver="+devServer)
	}
	if err := cmd.Start(); err != nil {
		runtime.LogError(app.Ctx, fmt.Sprintf("启动浮窗子进程失败: %v", err))
		return
	}

	app.FloatCmd[windowName] = cmd

}

// SetBottom 将窗口置于底层并从任务栏隐藏
// 适用于桌面小工具类窗口（如浮窗），使其像桌面挂件一样固定在最底层且不抢焦点
func SetBottom() {

	// 获取当前窗口句柄并缓存，供后续置顶/取消置顶等操作使用
	hwnd, err := FindWindow("", model.WindowInstance.WindowName)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	ctx := model.AppInstance.Ctx
	// 1. 移除设置重叠窗口样式
	exStyle := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE)
	newExStyle := exStyle &^ w32.WS_EX_TOPMOST
	w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE, newExStyle)
	// 2. 找到桌面窗口（Progman）。仅靠 SetWindowPos(HWND_BOTTOM) 置底是一次性的，
	// 窗口被激活或切换窗口后会被顶上来；把窗口挂载为桌面子窗口才能真正固定在最底层。
	progman := w32.FindWindow("Progman", "")
	if progman == 0 {
		// 找不到桌面窗口时降级为一次性置底
		runtime.LogWarning(ctx, "未找到桌面窗口(Progman)，降级为 HWND_BOTTOM 置底")
		w32.SetWindowPos(w32.HWND(hwnd), w32.HWND_BOTTOM, 0, 0, 0, 0,
			w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE)
		return
	}
	const wmSpawnWorker = 0x052C
	// 通知 Progman 先创建桌面图标列表（SHELLDLL_DefView），保证桌面层级完整
	w32.SendMessage(progman, wmSpawnWorker, 0, 0)
	// 等待 SHELLDLL_DefView 出现（最多约 300ms）
	var defView w32.HWND
	for i := 0; i < 30; i++ {
		defView = w32.FindWindowEx(progman, 0, "SHELLDLL_DefView", "")
		if defView != 0 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	// 挂载目标：默认挂到 Progman（桌面图标之上、其他窗口之下）
	// 若希望“壁纸级”效果（桌面图标之下），可将 target 改为挂在图标列表上：
	//   if listView := w32.FindWindowEx(defView, 0, "SysListView32", ""); listView != 0 {
	//       target = listView
	//   }
	target := progman

	// 3. 将窗口挂载为桌面的子窗口，实现永久置底
	if ret, _, _ := setParent.Call(uintptr(hwnd), uintptr(target)); ret == 0 {
		runtime.LogError(ctx, "挂载到桌面失败(SetParent 返回 0)")
	}

	// 4. 在桌面窗口内将窗口置于 Z 序底部并刷新样式
	w32.SetWindowPos(w32.HWND(hwnd), w32.HWND_TOP, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE|w32.SWP_FRAMECHANGED)
	// runtime.WindowSetAlwaysOnTop(ctx, false)
}

// SetTopMost 将窗口置于顶层
func SetTopMost() {
	ctx := model.AppInstance.Ctx
	parenthwnd, err := FindWindow("Progman", "")
	if err != nil {
		fmt.Println("error", err)
		return
	}
	hwnd, err := FindChildWindow(parenthwnd, "", model.WindowInstance.WindowName)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	// style := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_STYLE)
	// newStyle := (style &^ w32.WS_CHILD) | w32.WS_POPUP
	// w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_STYLE, newStyle)
	exStyle := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE)
	newExStyle := exStyle | w32.WS_EX_TOPMOST
	w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE, newExStyle)
	if ret, _, _ := setParent.Call(uintptr(hwnd), 0); ret == 0 {
		runtime.LogError(ctx, "还原桌面工具失败(SetParent 返回 0)")
		log.Println("还原桌面工具失败(SetParent 返回 0) w32.GetLastError() =", w32.GetLastError())
		return
	}
	// w32.ShowWindow(w32.HWND(hwnd), w32.SW_SHOW)
	w32.SetWindowPos(w32.HWND(hwnd), w32.HWND_TOPMOST, 0, 0, 0, 0,
		w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_NOACTIVATE|w32.SWP_FRAMECHANGED)
}

// 设置浮窗位置和尺寸
func SetFloatWindowPositionAndSize() {
	app := model.AppInstance
	window := model.WindowInstance
	if window.WindowMode != "float" {
		return
	}
	// SetWindowAsDesktopTool(ctx, window)
	runtime.WindowSetPosition(app.Ctx, model.AppConfig.WindowFloatConfigs[window.WindowName].WindowPositionX, model.AppConfig.WindowFloatConfigs[window.WindowName].WindowPositionY)
	runtime.WindowSetSize(app.Ctx, model.AppConfig.WindowFloatConfigs[window.WindowName].WindowSizeW, model.AppConfig.WindowFloatConfigs[window.WindowName].WindowSizeH)

}

// getMainWindowHandle 安全获取主窗口句柄
// 优先使用缓存的 HWND（通过 IsWindow 验证有效性），无效时降级为 FindWindow 按名称查找
func getWindowHandle() (syscall.Handle, error) {
	// 先尝试缓存的句柄
	if hwnd := model.WindowInstance.Hwnd; hwnd != 0 {
		ret, _, _ := isWindow.Call(uintptr(hwnd))
		if ret != 0 {
			return hwnd, nil
		}
		// 缓存的 HWND 已失效（窗口被销毁重建），清零并降级查找
		log.Println("getMainWindowHandle: 缓存的 HWND 已失效，降级为 FindWindow 查找")
		model.WindowInstance.Hwnd = 0
	}

	// 降级：通过窗口名称查找
	hwnd, err := FindWindow("", model.WindowInstance.WindowName)
	if err != nil {
		return 0, err
	}
	model.WindowInstance.Hwnd = hwnd
	return hwnd, nil
}

// restoreMainWindow 从托盘恢复主窗口到任务栏
// 使用原生 Windows API 而非 Wails runtime 调用，避免 WebView2 不响应时阻塞
func RestoreMainWindow() {
	app := model.AppInstance
	if app == nil || app.Ctx == nil {
		log.Println("RestoreMainWindow: app 或 ctx 为空，无法恢复窗口")
		return
	}

	hwnd, err := getWindowHandle()
	if err != nil {
		log.Printf("RestoreMainWindow: 未找到窗口 %q: %v", model.WindowInstance.WindowName, err)
		return
	}

	// 先检查窗口是否已可见，避免重复操作
	if w32.IsWindowVisible(w32.HWND(hwnd)) {
		ShowInTaskbar(hwnd)
		return
	}

	// 恢复任务栏显示（可能触发窗口样式变更，操作后刷新 HWND 缓存）
	ShowInTaskbar(hwnd)

	// 使用原生 API 显示窗口，不依赖 Wails runtime（WebView2 可能卡死导致 runtime 调用阻塞）
	w32.ShowWindow(w32.HWND(hwnd), w32.SW_SHOW)
	w32.SetForegroundWindow(w32.HWND(hwnd))
}

func MinimizeToTray() {
	app := model.AppInstance
	if app == nil || app.Ctx == nil {
		log.Println("MinimizeToTray: app 或 ctx 为空")
		return
	}

	// 确保托盘图标已启动
	StartTray(model.AppInstance)

	hwnd, err := getWindowHandle()
	if err != nil {
		log.Printf("MinimizeToTray: 未找到窗口 %q: %v", model.WindowInstance.WindowName, err)
		return
	}

	// 隐藏任务栏图标（可能触发窗口样式变更，操作后刷新 HWND 缓存）
	HideFromTaskbar(hwnd)

	// 使用原生 API 隐藏窗口
	w32.ShowWindow(w32.HWND(hwnd), w32.SW_HIDE)
}
