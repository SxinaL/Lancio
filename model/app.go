package model

import (
	"Lancio/database/method"
	"context"
	"os/exec"
)

type App struct {
	Ctx      context.Context
	Db       *method.DB
	FloatCmd map[string]*exec.Cmd // 子进程（浮窗）句柄
}

var AppInstance *App


func InitAppInstance(ctx context.Context, db *method.DB) {
	AppInstance = &App{
		Ctx: ctx,
		FloatCmd: make(map[string]*exec.Cmd),
		Db:  db,
	}
}
