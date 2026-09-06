package services

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"Lancio/database/method"
	"Lancio/model"
	"Lancio/utils"

	"github.com/gonutz/w32/v2"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Startup 应用启动时调用
func Startup(ctx context.Context) {

	// 初始化数据库
	db, err := method.New(utils.GetFilePath("vocabulary.db"))
	if err != nil {
		runtime.LogError(ctx, fmt.Sprintf("数据库初始化失败: %v", err))
		return
	}
	model.AppInstance = &model.App{
		Ctx:      ctx,
		FloatCmd: make(map[string]*exec.Cmd),
		Db:       db,
	}

	// 启动内存缓存定期清理（每 10 分钟清掉过期条目，避免长尾泄漏）
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			utils.CleanupResultCache()
		}
	}()

	// 主窗口进程启动时，根据配置决定是否拉起浮窗子进程
	if model.WindowInstance.WindowMode == "main" {
		for windowName, config := range model.AppConfig.WindowFloatConfigs {
			if config.Launch {
				utils.SpawnFloatWindow(windowName)
			}
		}
	}

	// 浮窗进程启动全局划词监听：按住 Ctrl + 鼠标划词即可获取选中文本
	if model.WindowInstance.WindowMode == "float" {
		StartTextSelectionHook()
	}
}

// DomReady 应用启动时调用，移除浮窗边框
func DomReady(ctx context.Context) {

	// 获取窗口句柄并缓存，后续所有窗口操作都通过缓存的句柄进行，
	// 避免每次 FindWindow 按名称查找（窗口标题可能被前端修改导致查找失败）
	hwnd, err := utils.FindWindow("", model.WindowInstance.WindowName)
	if err != nil {
		fmt.Println("error", err)
		return
	}
	
	model.WindowInstance.OriginalWndProc = w32.GetWindowLongPtr(w32.HWND(hwnd), w32.GWLP_WNDPROC)

	utils.ModifyWndProc()
	utils.ModifyWindowStyle()
	utils.SetFloatWindowPositionAndSize()

	if model.WindowInstance.WindowMode == "float" {

		isPinned := model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName].IsPinned

		if isPinned {
			utils.SetTopMost()
		} else {
			utils.SetBottom()
		}
	}
	model.WindowInstance.Hwnd = hwnd
}

// Shutdown 应用关闭时调用
func Shutdown(ctx context.Context) {
	if model.AppInstance == nil {
		return
	}

	// 浮窗进程退出时，通知钩子线程优雅退出
	if model.WindowInstance.WindowMode == "float" {
		StopHook()
		utils.UnregisterAllHotkeys()
		utils.StopHotkeyListener()
	}

	// 主窗口关闭时，终止浮窗子进程
	if model.AppInstance.FloatCmd != nil {
		for _, cmd := range model.AppInstance.FloatCmd {
			if cmd.Process != nil {
				killCmd := exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid))
				killCmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
				if err := killCmd.Run(); err != nil {
					runtime.LogError(ctx, fmt.Sprintf("终止浮窗子进程失败: %v", err))
				}
			}
		}
	}
	if model.AppInstance.Db != nil {
		model.AppInstance.Db.Close()
	}
}
