package main

import (
	"context"
	"fmt"

	"mousePaw_new/pkg/autostart"
	"mousePaw_new/pkg/config"
	"mousePaw_new/pkg/engine"
	"mousePaw_new/pkg/log"

	"github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	cfg    *config.Config
	engine *engine.Engine
	logger *log.Logger
}

func NewApp() *App {
	cfg := config.Load()
	logger := log.NewLogger()
	return &App{
		cfg:    cfg,
		engine: engine.NewEngine(cfg, logger),
		logger: logger,
	}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	a.engine.SetStatusCallback(func(status engine.Status) {
		runtime.EventsEmit(a.ctx, "statusChanged", string(status))
	})

	a.logger.SetOnAppend(func(entry log.LogEntry) {
		runtime.EventsEmit(a.ctx, "newLog", entry)
	})

	a.logger.Info("应用程序启动")
	go a.startHotkeyListener()
}

func (a *App) startHotkeyListener() {
	hook.Register(hook.KeyDown, []string{"f6"}, func(e hook.Event) {
		status := a.engine.GetStatus()
		if status == engine.StatusStopped {
			a.engine.Start()
		}
	})

	hook.Register(hook.KeyDown, []string{"f7"}, func(e hook.Event) {
		status := a.engine.GetStatus()
		if status == engine.StatusRunning {
			a.engine.Stop()
		}
	})

	s := hook.Start()
	<-hook.Process(s)
}

func (a *App) domready(ctx context.Context) {
}

func (a *App) shutdown(ctx context.Context) {
	a.logger.Info("应用程序关闭")
	hook.End()
	a.engine.Stop()
	a.logger.Close()
}

func (a *App) GetConfig() config.Config {
	return a.cfg.GetAll()
}

func (a *App) UpdateConfig(cfg config.Config) error {
	a.cfg.Update(cfg)
	a.engine.ReloadConfig(a.cfg)

	if cfg.AutoStart {
		if err := autostart.Enable(); err != nil {
			return fmt.Errorf("设置开机自启失败: %w", err)
		}
	} else {
		if err := autostart.Disable(); err != nil {
			return fmt.Errorf("取消开机自启失败: %w", err)
		}
	}

	return a.cfg.Save()
}

func (a *App) GetStatus() string {
	return string(a.engine.GetStatus())
}

func (a *App) Start() {
	a.engine.Start()
}

func (a *App) Stop() {
	a.engine.Stop()
}

func (a *App) GetLogs() []log.LogEntry {
	return a.logger.GetEntries()
}

func (a *App) IsAutoStartEnabled() bool {
	return autostart.IsEnabled()
}
