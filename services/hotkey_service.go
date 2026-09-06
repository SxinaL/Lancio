package services

import (
	"Lancio/model"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func ToggleFloatWindow() {

	app := model.AppInstance
	// 通知前端状态变了
    runtime.EventsEmit(app.Ctx, "pin:changed","")
}
