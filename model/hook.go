package model


var (
	SelectionChan = make(chan struct{}, 1)
	TextSelection = make(chan string,1)
)

// CtrlSelecting 记录本次按下左键时是否按住 Ctrl（避免在释放瞬间检测 Ctrl 时被系统菜单模式延迟）
var CtrlSelecting bool
