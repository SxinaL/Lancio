package main

import (
	"context"
	"fmt"
    "syscall"
	"unsafe"
	"DesktopVoc/database"
	"DesktopVoc/services"
	

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/gonutz/w32/v2"
)

// App 应用主结构体
type App struct {
	ctx  context.Context
	db   *database.DB
	dict services.DictionaryService
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// 初始化数据库
	db, err := database.New("vocabulary.db")
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("数据库初始化失败: %v", err))
		return
	}
	a.db = db

	// 初始化词典服务（默认占位实现）
	a.dict = services.NewDefaultDictionaryService()

	// 获取屏幕尺寸并定位窗口到右下角
	a.positionWindow(ctx)
}

// shutdown 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		a.db.Close()
	}
}

// positionWindow 将窗口定位到屏幕右下角
func (a *App) positionWindow(ctx context.Context) {
	screen, err := runtime.ScreenGetAll(ctx)
	if err != nil || len(screen) == 0 {
		return
	}

	primary := screen[0]
	winWidth := 360
	winHeight := 340
	x := primary.Width/2 + winWidth/2
	y := primary.Height/2 + winHeight/2

	runtime.WindowSetSize(ctx, winWidth, winHeight)
	runtime.WindowSetPosition(ctx, x, y)
}

func (a *App) SetOpacity(opacity float64) {
	if a.ctx == nil {
		return
	}

	// 限制透明度范围
	if opacity < 0.1 {
		opacity = 0.1
	}
	if opacity > 1.0 {
		opacity = 1.0
	}

	// 通过 Wails runtime 设置背景色透明度
	alpha := uint8(opacity * 255)
	runtime.WindowSetBackgroundColour(a.ctx, 27, 38, 54, alpha)
}

const (
	WS_CAPTION          = 0x00C00000
	WS_THICKFRAME       = 0x00040000
	WS_MINIMIZEBOX      = 0x00020000
	WS_MAXIMIZEBOX      = 0x00010000
	WS_SYSMENU          = 0x00080000
	WS_EX_DLGMODALFRAME = 0x00000001
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	findWindow = user32.NewProc("FindWindowW")
)

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
    r1, _, e1 := syscall.Syscall(findWindow.Addr(), 2, uintptr(unsafe.Pointer(cname)), uintptr(unsafe.Pointer(wname)), 0)
    if r1 == 0 {
        if e1 != 0 {
            err = error(e1)
        } else {
            err = syscall.EINVAL
        }
    }
    hwnd = syscall.Handle(r1)
    return
}

func (a *App) domReady(ctx context.Context) {
	hwnd,err := FindWindow("", "DesktopVoc")
	if err != nil {
        fmt.Println("error", err)
        return
    }

	go func() {
		// 用 w32 库一行代码移除边框（已封装好 Get/SetWindowLongPtr）
		style := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_STYLE)
		// 移除标题栏、系统菜单、最小化/最大化、可调整边框
		newStyle := style &^ (w32.WS_CAPTION | w32.WS_SYSMENU | w32.WS_MINIMIZEBOX | w32.WS_MAXIMIZEBOX | w32.WS_THICKFRAME)
		w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_STYLE, newStyle)

		// 移除对话框边框（扩展样式）
		exStyle := w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE)
		newExStyle := exStyle &^ 0x00000001 // WS_EX_DLGMODALFRAME
		w32.SetWindowLongPtr(w32.HWND(hwnd), w32.GWL_EXSTYLE, newExStyle)

		// 刷新边框
		w32.SetWindowPos(w32.HWND(hwnd), 0, 0, 0, 0, 0,
			w32.SWP_NOMOVE|w32.SWP_NOSIZE|w32.SWP_FRAMECHANGED)

		// 显示窗口
		runtime.WindowShow(ctx)
	}()
}
