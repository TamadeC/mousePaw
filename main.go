package main

import (
	"context"
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	tray := NewSystemTray(app)

	err := wails.Run(&options.App{
		Title:     "MousePaw",
		Width:     480,
		Height:    600,
		MinWidth:  400,
		MinHeight: 500,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 15, G: 23, B: 42, A: 255},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			tray.SetContext(ctx)
			tray.SetOnShow(func() {
				runtime.WindowShow(ctx)
			})
			tray.SetOnQuit(func() {
				runtime.Quit(ctx)
			})
			tray.Start()
		},
		OnDomReady: app.domready,
		OnShutdown: app.shutdown,
		OnBeforeClose: func(ctx context.Context) bool {
			if app.cfg.MinimizeToTray {
				runtime.WindowHide(ctx)
				return true
			}
			return false
		},
		HideWindowOnClose: true,
		Windows: &windows.Options{
			Theme: windows.Dark,
		},
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
