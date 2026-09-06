package main

import (
	// "Lancio/model"
	"Lancio/handler"
	"Lancio/model"
	"Lancio/services"
	"Lancio/utils"
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	// "github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/windows/Lancio.ico
var trayIcon []byte
func initHandlers() (*handler.WindowHandler, *handler.VocabularyHandler, *handler.VocabularyBankHandler, *handler.TaskbarHandler) {
	window_handler := handler.NewWindowHandler()
	vocabulary_handler := handler.NewVocabularyHandler()
	vocabularybank_handler := handler.NewBankHandler()
	taskbar_handler := handler.NewTaskbarHandler()
	return window_handler, vocabulary_handler, vocabularybank_handler, taskbar_handler
}
// RunMainWindow 运行主窗口 (800x600)
func RunMainWindow() error {
	window_handler, vocabulary_handler, vocabularybank_handler, taskbar_handler := initHandlers()
	
	window := model.WindowInstance
	err := wails.Run(&options.App{
		Title:     window.WindowName,
		Width:     800,
		Height:    600,
		MinWidth:  100,
		MinHeight: 100,
		StartHidden: model.AppConfig.StartHidden,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 26, G: 26, B: 46, A: 255},
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			// WebviewUserDataPath:  "./webview-data-main",
		},
		Mac: &mac.Options{
			Appearance: mac.NSAppearanceNameDarkAqua,
		},
		Linux:      &linux.Options{},
		OnStartup:  services.Startup,
		OnDomReady: services.DomReady,
		OnShutdown: services.Shutdown,
		Bind: []interface{}{
			window_handler,
			vocabulary_handler,
			vocabularybank_handler,
			taskbar_handler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
		return err
	}
	return nil
}

// RunFloatWindow 运行浮窗 (275x300, 无边框透明)
func RunFloatWindow()  error {
	window := model.WindowInstance
	cfg := model.AppConfig
	window_handler, vocabulary_handler, vocabularybank_handler, taskbar_handler := initHandlers()
	
	err := wails.Run(&options.App{
		Title:     window.WindowName,
		Width:     cfg.WindowFloatConfigs[window.WindowName].WindowSizeW,
		Height:    cfg.WindowFloatConfigs[window.WindowName].WindowSizeH,
		MinWidth:  cfg.WindowFloatConfigs[window.WindowName].WindowSizeMinW,
		MinHeight: cfg.WindowFloatConfigs[window.WindowName].WindowSizeMinH,

		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 0},
		// Frameless: true,

		Windows: &windows.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			// WebviewUserDataPath:  "./webview-data-float",
		},
		Mac: &mac.Options{
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		Linux:      &linux.Options{},
		OnStartup:  services.Startup,
		OnDomReady: services.DomReady,
		OnShutdown: services.Shutdown,

		Bind: []interface{}{
			window_handler,
			vocabulary_handler,
			vocabularybank_handler,
			taskbar_handler,
		},
	})

	if err != nil {
		println("Error:", err.Error())
		return err
	}
	return nil
}
func LogInit(){
	exePath, err := os.Executable()
	if err != nil {
        // 如果获取失败，就退而求其次存到当前工作目录
        exePath = "."
    }
	logDir := filepath.Dir(exePath)
    logFile := filepath.Join(logDir, "Lancio.log")
	  f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        // 如果连 exe 目录都写不了（比如装在 Program Files 里权限不够），就写到用户临时目录
        tmpDir := os.TempDir()
        logFile = filepath.Join(tmpDir, "lancio.log")
        f, _ = os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    }

    // 3. 关键命令：让标准库 log 包的所有输出都指向这个文件
    log.SetOutput(f)

    // 4. 给每条日志加上日期时间和文件名行号，否则根本不知道是哪一秒卡死的
    log.SetFlags(log.LstdFlags | log.Lshortfile)
    
    log.Println("======= 日志启动 =======")

}
func HotKeyInit(){
	utils.StartHotkeyListener(
		utils.HotkeyEntry{
			Modifiers: utils.ModAlt | utils.ModCtrl,
			VK:        utils.VK_L,
			Callback:  services.ToggleFloatWindow,
		},
	)
	// 注册热键


}
func main() {
	utils.LoadConfig()
	// 注入托盘图标（.ico）
	utils.SetTrayIcon(trayIcon)

	if len(os.Args) > 1 && os.Args[1] == "--float" {
		model.InitWindowInstance("LancioFloat", "float")
		os.Args = []string{os.Args[0]}
		HotKeyInit()
		RunFloatWindow()
		return
	}
	LogInit()
	model.InitWindowInstance("LancioMain", "main")
	RunMainWindow()
}
