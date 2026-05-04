package main

import (
	"context"
	_ "embed"

	"github.com/getlantern/systray"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed build/appicon.ico
var iconData []byte

type SystemTray struct {
	ctx    context.Context
	app    *App
	onShow func()
	onQuit func()
}

func NewSystemTray(app *App) *SystemTray {
	return &SystemTray{
		app: app,
	}
}

func (st *SystemTray) SetContext(ctx context.Context) {
	st.ctx = ctx
}

func (st *SystemTray) SetOnShow(fn func()) {
	st.onShow = fn
}

func (st *SystemTray) SetOnQuit(fn func()) {
	st.onQuit = fn
}

func (st *SystemTray) Start() {
	go systray.Run(st.onReady, st.onExit)
}

func (st *SystemTray) onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("MousePaw")
	systray.SetTooltip("MousePaw - 鼠标自动化工具")

	mShow := systray.AddMenuItem("显示窗口", "显示主窗口")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("退出", "退出应用程序")

	go func() {
		for {
			select {
			case <-mShow.ClickedCh:
				if st.onShow != nil {
					st.onShow()
				}
			case <-mQuit.ClickedCh:
				if st.onQuit != nil {
					st.onQuit()
				}
				systray.Quit()
				return
			}
		}
	}()
}

func (st *SystemTray) onExit() {
	if st.onQuit != nil {
		st.onQuit()
	}
}

func (st *SystemTray) ShowWindow() {
	if st.ctx != nil {
		runtime.WindowShow(st.ctx)
	}
}
