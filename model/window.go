package model

import (
	"syscall"
)

type Window struct {
	WindowMode      string
	WindowName      string
	Hwnd            syscall.Handle // 缓存的主窗口句柄，避免每次 FindWindow 按名称查找
	OriginalWndProc uintptr
}

var WindowInstance *Window

func InitWindowInstance(WindowName string, WindowMode string) {
	WindowInstance = &Window{
		WindowName:      WindowName,
		WindowMode:      WindowMode,
		Hwnd:            0,
		OriginalWndProc: 0,
	}
}
