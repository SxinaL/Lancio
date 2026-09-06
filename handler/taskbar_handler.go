package handler

import (

	"Lancio/utils"
)

type TaskbarHandler struct {}
func NewTaskbarHandler() *TaskbarHandler {
	return &TaskbarHandler{}
}
// MinimizeToTray 最小化到系统托盘：隐藏任务栏图标、隐藏窗口、确保托盘图标已启动
func (taskbarHandler *TaskbarHandler) MinimizeToTray() {
	utils.MinimizeToTray()
}

// RestoreFromTray 从托盘恢复主窗口并重新显示在任务栏
func (taskbarHandler *TaskbarHandler) RestoreFromTray() {
	utils.RestoreMainWindow()
}