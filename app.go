package main

import (
	"context"
	"fmt"
	"strings"

	"mousePaw_new/pkg/autostart"
	"mousePaw_new/pkg/config"
	"mousePaw_new/pkg/engine"
	"mousePaw_new/pkg/log"
	"mousePaw_new/pkg/recorder"

	hook "github.com/robotn/gohook"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx      context.Context
	cfg      *config.Config
	engine   *engine.Engine
	logger   *log.Logger
	recorder *recorder.Recorder
}

func NewApp() *App {
	cfg := config.Load()
	logger := log.NewLogger()
	return &App{
		cfg:      cfg,
		engine:   engine.NewEngine(cfg, logger),
		logger:   logger,
		recorder: recorder.NewRecorder(),
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

	a.recorder.SetOnChange(func(info recorder.RecorderStatusInfo) {
		runtime.EventsEmit(a.ctx, "recorderStatus", info)
	})

	a.logger.Info("应用程序启动")
	go a.startHotkeyListener()
}

func (a *App) startHotkeyListener() {
	// 解析快捷键配置
	startKeys := parseHotkey(a.cfg.Hotkeys.Start)
	stopKeys := parseHotkey(a.cfg.Hotkeys.Stop)
	pauseKeys := parseHotkey(a.cfg.Hotkeys.Pause)

	// 注册开始快捷键
	hook.Register(hook.KeyDown, startKeys, func(e hook.Event) {
		status := a.engine.GetStatus()
		if status == engine.StatusStopped {
			a.engine.Start()
		}
	})

	// 注册停止快捷键
	hook.Register(hook.KeyDown, stopKeys, func(e hook.Event) {
		status := a.engine.GetStatus()
		if status == engine.StatusRunning || status == engine.StatusPaused {
			a.engine.Stop()
		}
	})

	// 注册暂停/恢复快捷键
	hook.Register(hook.KeyDown, pauseKeys, func(e hook.Event) {
		status := a.engine.GetStatus()
		if status == engine.StatusRunning {
			a.engine.Pause()
		} else if status == engine.StatusPaused {
			a.engine.Resume()
		}
	})

	s := hook.Start()
	<-hook.Process(s)
}

func parseHotkey(hotkey string) []string {
	keys := []string{}
	for _, key := range splitHotkey(hotkey) {
		keys = append(keys, strings.TrimSpace(key))
	}
	return keys
}

func splitHotkey(hotkey string) []string {
	var result []string
	current := ""
	for _, ch := range hotkey {
		if ch == '+' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
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

func (a *App) Pause() {
	a.engine.Pause()
}

func (a *App) Resume() {
	a.engine.Resume()
}

func (a *App) GetLogs() []log.LogEntry {
	return a.logger.GetEntries()
}

func (a *App) IsAutoStartEnabled() bool {
	return autostart.IsEnabled()
}

func (a *App) UpdateHotkeys(hotkeys config.HotkeyConfig) error {
	a.cfg.Hotkeys = hotkeys
	return a.cfg.Save()
}

func (a *App) StartRecording() {
	a.recorder.Start()
	a.logger.Info("开始录制鼠标和键盘操作")
}

func (a *App) StopRecording() *recorder.Recording {
	rec := a.recorder.Stop()
	if rec != nil {
		a.logger.Info(fmt.Sprintf("录制停止，共 %d 个动作，总时长 %.1f 秒", len(rec.Actions), rec.Duration))
	}
	return rec
}

func (a *App) GetRecordingStatus() recorder.RecorderStatusInfo {
	return a.recorder.GetStatus()
}

func (a *App) GetRecording() *recorder.Recording {
	return a.recorder.GetRecording()
}

func (a *App) ClearRecording() {
	a.recorder.Clear()
	a.logger.Info("录制数据已清除")
}

func (a *App) SaveRecording(name string) error {
	rec := a.recorder.GetRecording()
	if rec == nil {
		return fmt.Errorf("没有录制数据")
	}
	if err := recorder.SaveRecording(name, rec); err != nil {
		return fmt.Errorf("保存录制失败: %w", err)
	}
	a.logger.Info(fmt.Sprintf("录制已保存: %s", name))
	return nil
}

func (a *App) LoadRecording(name string) (*recorder.Recording, error) {
	rec, err := recorder.LoadRecording(name)
	if err != nil {
		return nil, fmt.Errorf("加载录制失败: %w", err)
	}
	a.logger.Info(fmt.Sprintf("录制已加载: %s，共 %d 个动作", name, len(rec.Actions)))
	return rec, nil
}

func (a *App) ListRecordings() []string {
	return recorder.ListRecordings()
}

func (a *App) DeleteRecording(name string) error {
	if err := recorder.DeleteRecording(name); err != nil {
		return fmt.Errorf("删除录制失败: %w", err)
	}
	a.logger.Info(fmt.Sprintf("录制已删除: %s", name))
	return nil
}
