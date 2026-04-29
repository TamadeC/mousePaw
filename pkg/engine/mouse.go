package engine

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"mousePaw_new/pkg/config"
	"mousePaw_new/pkg/log"

	"github.com/go-vgo/robotgo"
)

type Status string

const (
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
)

type Engine struct {
	cfg      *config.Config
	status   Status
	mu       sync.Mutex
	stopChan chan struct{}
	logger   *log.Logger

	onStatus func(Status)
}

func NewEngine(cfg *config.Config, logger *log.Logger) *Engine {
	return &Engine{
		cfg:      cfg,
		status:   StatusStopped,
		stopChan: make(chan struct{}),
		logger:   logger,
	}
}

func (e *Engine) SetStatusCallback(fn func(Status)) {
	e.onStatus = fn
}

func (e *Engine) GetStatus() Status {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.status
}

func (e *Engine) setStatus(s Status) {
	e.mu.Lock()
	e.status = s
	e.mu.Unlock()
	if e.onStatus != nil {
		e.onStatus(s)
	}
}

func (e *Engine) Start() {
	e.mu.Lock()
	if e.status == StatusRunning {
		e.mu.Unlock()
		return
	}
	e.stopChan = make(chan struct{})
	e.mu.Unlock()

	e.setStatus(StatusRunning)
	e.logger.Info("引擎启动")
	go e.run()
}

func (e *Engine) Stop() {
	e.mu.Lock()
	if e.status == StatusStopped {
		e.mu.Unlock()
		return
	}
	e.mu.Unlock()

	close(e.stopChan)
	e.setStatus(StatusStopped)
	e.logger.Info("引擎停止")
}

func (e *Engine) ReloadConfig(cfg *config.Config) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cfg = cfg
	e.logger.Info("配置已更新")
}

func (e *Engine) run() {
	var moveTicker, clickTicker, scrollTicker *time.Ticker
	var moveDone, clickDone, scrollDone <-chan time.Time

	if e.cfg.MoveEnabled {
		interval := time.Duration(e.cfg.MoveInterval * float64(time.Second))
		moveTicker = time.NewTicker(interval)
		moveDone = moveTicker.C
		defer moveTicker.Stop()
		e.logger.Info(fmt.Sprintf("鼠标移动已启用，间隔 %.1f 秒", e.cfg.MoveInterval))
	}

	if e.cfg.ClickEnabled {
		interval := time.Duration(e.cfg.ClickInterval * float64(time.Second))
		clickTicker = time.NewTicker(interval)
		clickDone = clickTicker.C
		defer clickTicker.Stop()
		e.logger.Info(fmt.Sprintf("鼠标点击已启用，间隔 %.1f 秒", e.cfg.ClickInterval))
	}

	if e.cfg.ScrollEnabled {
		interval := time.Duration(e.cfg.ScrollInterval * float64(time.Second))
		scrollTicker = time.NewTicker(interval)
		scrollDone = scrollTicker.C
		defer scrollTicker.Stop()
		e.logger.Info(fmt.Sprintf("滚轮滚动已启用，间隔 %.1f 秒", e.cfg.ScrollInterval))
	}

	for {
		select {
		case <-e.stopChan:
			return
		case <-moveDone:
			e.doMove()
		case <-clickDone:
			e.doClick()
		case <-scrollDone:
			e.doScroll()
		}
	}
}

func (e *Engine) doMove() {
	w, h := robotgo.GetScreenSize()
	if w <= 0 || h <= 0 {
		return
	}

	targetX := rand.Intn(w)
	targetY := rand.Intn(h)

	curX, curY := robotgo.Location()
	steps := 30

	for i := 1; i <= steps; i++ {
		select {
		case <-e.stopChan:
			return
		default:
		}

		x := curX + (targetX-curX)*i/steps
		y := curY + (targetY-curY)*i/steps
		robotgo.Move(x, y)
		time.Sleep(15 * time.Millisecond)
	}

	e.logger.Info(fmt.Sprintf("鼠标移动到 (%d, %d)", targetX, targetY))
}

func (e *Engine) doClick() {
	btn := string(e.cfg.ClickType)
	count := e.cfg.ClickCount

	btnName := map[string]string{
		"left":   "左键",
		"right":  "右键",
		"middle": "中键",
	}[btn]

	for i := 0; i < count; i++ {
		robotgo.Click(btn)
		if i < count-1 {
			time.Sleep(50 * time.Millisecond)
		}
	}

	if count > 1 {
		e.logger.Info(fmt.Sprintf("执行点击: %s x%d", btnName, count))
	} else {
		e.logger.Info(fmt.Sprintf("执行点击: %s", btnName))
	}
}

func (e *Engine) doScroll() {
	dir := string(e.cfg.ScrollDir)
	amount := e.cfg.ScrollAmount

	dirName := map[string]string{
		"up":    "向上",
		"down":  "向下",
		"left":  "向左",
		"right": "向右",
	}[dir]

	robotgo.ScrollDir(amount, dir)
	e.logger.Info(fmt.Sprintf("执行滚动: %s %d格", dirName, amount))
}
