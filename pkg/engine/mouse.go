package engine

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"mousePaw_new/pkg/config"
	"mousePaw_new/pkg/log"
	"mousePaw_new/pkg/recorder"

	"github.com/go-vgo/robotgo"
)

type Status string

const (
	StatusStopped Status = "stopped"
	StatusRunning Status = "running"
	StatusPaused  Status = "paused"
)

type Engine struct {
	cfg       *config.Config
	status    Status
	mu        sync.Mutex
	stopChan  chan struct{}
	pauseChan chan struct{}
	logger    *log.Logger
	replay    *recorder.Replay

	onStatus func(Status)
}

func NewEngine(cfg *config.Config, logger *log.Logger) *Engine {
	rp := recorder.NewReplay()
	rp.SetOnLog(func(msg string) {
		logger.Info(msg)
	})
	return &Engine{
		cfg:       cfg,
		status:    StatusStopped,
		stopChan:  make(chan struct{}),
		pauseChan: make(chan struct{}),
		logger:    logger,
		replay:    rp,
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

	e.replay.Stop()
	close(e.stopChan)
	e.setStatus(StatusStopped)
	e.logger.Info("引擎停止")
}

func (e *Engine) Pause() {
	e.mu.Lock()
	if e.status != StatusRunning {
		e.mu.Unlock()
		return
	}
	e.mu.Unlock()

	close(e.pauseChan)
	e.replay.Pause()
	e.setStatus(StatusPaused)
	e.logger.Info("引擎暂停")
}

func (e *Engine) Resume() {
	e.mu.Lock()
	if e.status != StatusPaused {
		e.mu.Unlock()
		return
	}
	e.pauseChan = make(chan struct{})
	e.mu.Unlock()

	e.replay.Resume()
	e.setStatus(StatusRunning)
	e.logger.Info("引擎恢复")
}

func (e *Engine) ReloadConfig(cfg *config.Config) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.cfg = cfg
	e.logger.Info("配置已更新")
}

func (e *Engine) GetReplay() *recorder.Replay {
	return e.replay
}

func (e *Engine) run() {
	if e.cfg.OperationType == config.OpReplay {
		e.runReplayLoop()
		return
	}

	var ticker *time.Ticker
	var done <-chan time.Time

	switch e.cfg.OperationType {
	case config.OpMove:
		interval := time.Duration(e.cfg.MoveInterval * float64(time.Second))
		ticker = time.NewTicker(interval)
		done = ticker.C
		defer ticker.Stop()
		e.logger.Info(fmt.Sprintf("鼠标移动已启用，间隔 %.1f 秒", e.cfg.MoveInterval))
	case config.OpClick:
		interval := time.Duration(e.cfg.ClickInterval * float64(time.Second))
		ticker = time.NewTicker(interval)
		done = ticker.C
		defer ticker.Stop()
		e.logger.Info(fmt.Sprintf("鼠标点击已启用，间隔 %.1f 秒", e.cfg.ClickInterval))
	case config.OpScroll:
		interval := time.Duration(e.cfg.ScrollInterval * float64(time.Second))
		ticker = time.NewTicker(interval)
		done = ticker.C
		defer ticker.Stop()
		e.logger.Info(fmt.Sprintf("滚轮滚动已启用，间隔 %.1f 秒", e.cfg.ScrollInterval))
	case config.OpType:
		interval := time.Duration(e.cfg.TypeInterval * float64(time.Second))
		ticker = time.NewTicker(interval)
		done = ticker.C
		defer ticker.Stop()
		e.logger.Info(fmt.Sprintf("键盘输入已启用，间隔 %.1f 秒", e.cfg.TypeInterval))
	default:
		e.logger.Error("未知的操作类型: " + string(e.cfg.OperationType))
		return
	}

	for {
		select {
		case <-e.stopChan:
			return
		case <-e.pauseChan:
			// 暂停状态，等待恢复或停止
			e.logger.Info("引擎已暂停，等待恢复...")
			for {
				select {
				case <-e.stopChan:
					return
				default:
					// 检查是否已恢复（pauseChan 被重新创建）
					e.mu.Lock()
					paused := e.status == StatusPaused
					e.mu.Unlock()
					if !paused {
						// 已恢复，重新创建 ticker
						switch e.cfg.OperationType {
						case config.OpMove:
							interval := time.Duration(e.cfg.MoveInterval * float64(time.Second))
							ticker = time.NewTicker(interval)
							done = ticker.C
						case config.OpClick:
							interval := time.Duration(e.cfg.ClickInterval * float64(time.Second))
							ticker = time.NewTicker(interval)
							done = ticker.C
						case config.OpScroll:
							interval := time.Duration(e.cfg.ScrollInterval * float64(time.Second))
							ticker = time.NewTicker(interval)
							done = ticker.C
						case config.OpType:
							interval := time.Duration(e.cfg.TypeInterval * float64(time.Second))
							ticker = time.NewTicker(interval)
							done = ticker.C
						}
						break
					}
					time.Sleep(100 * time.Millisecond)
				}
				// 检查是否已恢复
				e.mu.Lock()
				paused := e.status == StatusPaused
				e.mu.Unlock()
				if !paused {
					break
				}
			}
		case <-done:
			switch e.cfg.OperationType {
			case config.OpMove:
				e.doMove()
			case config.OpClick:
				e.doClick()
			case config.OpScroll:
				e.doScroll()
			case config.OpType:
				e.doType()
			}
		}
	}
}

func (e *Engine) runReplayLoop() {
	file := e.cfg.ReplayFile
	if file == "" {
		e.logger.Error("未指定回放录制文件")
		e.Stop()
		return
	}

	rec, err := recorder.LoadRecording(file)
	if err != nil {
		e.logger.Error(fmt.Sprintf("加载录制文件失败: %v", err))
		e.Stop()
		return
	}

	e.replay.SetActions(rec.Actions)
	e.logger.Info(fmt.Sprintf("回放模式已启用，文件: %s，共 %d 个动作，总时长 %.1f 秒",
		file, len(rec.Actions), rec.Duration))

	interval := time.Duration(e.cfg.ReplayInterval * float64(time.Second))

	for {
		if !e.waitIfPaused() {
			return
		}

		e.replay.Start()
		e.replay.RunOnce()

		if e.replay.IsStopped() {
			e.Stop()
			return
		}

		if !e.cfg.ReplayRepeat {
			select {
			case <-e.stopChan:
				e.replay.Stop()
				return
			case <-e.pauseChan:
				e.logger.Info("引擎已暂停，等待恢复...")
				continue
			case <-time.After(interval):
			}
		}
	}
}

func (e *Engine) waitIfPaused() bool {
	for {
		e.mu.Lock()
		paused := e.status == StatusPaused
		e.mu.Unlock()

		if !paused {
			return true
		}

		select {
		case <-e.stopChan:
			e.replay.Stop()
			return false
		default:
			time.Sleep(100 * time.Millisecond)
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

func (e *Engine) doType() {
	text := e.cfg.TypeText
	if text == "" {
		return
	}

	robotgo.TypeStr(text)
	e.logger.Info(fmt.Sprintf("键盘输入: %s", text))
}
