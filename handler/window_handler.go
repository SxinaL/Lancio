package handler

import (
	"Lancio/model"

	"Lancio/utils"
	"fmt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)
type MinSize struct {
	WindowSizeMinW int
	WindowSizeMinH int
}
type WindowHandler struct{}

func NewWindowHandler() *WindowHandler {
	return &WindowHandler{}
}

// SetWindowOpacity 设置窗口透明度
func (windowHandler *WindowHandler) SetWindowOpacity(opacity float64) {

	app := model.AppInstance
	if app.Ctx == nil {
		return
	}

	// 限制透明度范围
	if opacity < 0.1 {
		opacity = 0.1
	}
	if opacity > 1.0 {
		opacity = 1.0
	}
	utils.LoadConfig()
	model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName].WindowOpacity = opacity
	utils.SaveConfig()
}

// GetWindowOpacity 获取窗口透明度
func (windowHandler *WindowHandler) GetWindowOpacity() float64 {
	return utils.GetFloatConfig(model.WindowInstance.WindowName).WindowOpacity
}

// PinFloatWindow 置底
func (windowHandler *WindowHandler) PinFloatWindow() {
	fmt.Println("置于顶层")
	utils.SetTopMost()
}

// UnpinFloatWindow 取消置底
func (windowHandler *WindowHandler) UnpinFloatWindow() {
	fmt.Println("置于底层")
	utils.SetBottom()
}

// ToggleFloatWindow 切换浮窗子进程运行状态：开启或关闭，并把状态写入 config.json
func (windowHandler *WindowHandler) ToggleFloatWindow(windowName string) bool {
	if model.AppInstance == nil {
		return false
	}
	if windowHandler.IsFloatWindowRunning(windowName) {
		// 关闭浮窗
		if err := model.AppInstance.FloatCmd[windowName].Process.Kill(); err != nil {
			runtime.LogError(model.AppInstance.Ctx, fmt.Sprintf("关闭浮窗失败: %v", err))
			return true // 仍认为在运行
		}
		delete(model.AppInstance.FloatCmd, windowName)
		utils.LoadConfig()
		utils.GetFloatConfig(windowName).Launch = false
		utils.SaveConfig()
		return false
	}
	// 启动浮窗
	utils.SpawnFloatWindow(windowName)
	utils.LoadConfig()
	utils.GetFloatConfig(windowName).Launch = true
	utils.SaveConfig()
	return windowHandler.IsFloatWindowRunning(windowName)
}

// SaveFloatWindowPositionAndSize 保存浮窗位置和尺寸
func (windowHandler *WindowHandler) SaveFloatWindowPositionAndSize() {
	utils.SaveFloatWindowPositionAndSize()
}

// GetPinStatus 获取置底状态
func (windowHandler *WindowHandler) GetPinStatus() bool {
	isPinned := utils.GetFloatConfig(model.WindowInstance.WindowName).IsPinned
	fmt.Printf("GetPinStatus: %v\n", isPinned)
	return isPinned
}

// SetPinStatus 设置置底状态
func (windowHandler *WindowHandler) SetPinStatus(pinStatus bool) {
	utils.LoadConfig()
	utils.GetFloatConfig(model.WindowInstance.WindowName).IsPinned = pinStatus
	utils.SaveConfig()
}

// 功能类
// IsFloatWindowRunning 返回浮窗子进程是否正在运行
func (windowHandler *WindowHandler) IsFloatWindowRunning(windowName string) bool {
	return model.AppInstance != nil && model.AppInstance.FloatCmd != nil && model.AppInstance.FloatCmd[windowName] != nil
}

// GetWindowMode 获取窗口模式
func (windowHandler *WindowHandler) GetWindowMode() string {
	return model.WindowInstance.WindowMode
}

func (windowHandler *WindowHandler) GetMinSize() MinSize {
	
	return MinSize{
		WindowSizeMinW: model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName].WindowSizeMinW,
		WindowSizeMinH: model.AppConfig.WindowFloatConfigs[model.WindowInstance.WindowName].WindowSizeMinH, 	
	}
}
