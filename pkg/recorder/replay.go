package recorder

import (
	"fmt"
	"sync"
	"time"

	"github.com/go-vgo/robotgo"
)

type ReplayStatus string

const (
	ReplayStopped ReplayStatus = "stopped"
	ReplayRunning ReplayStatus = "running"
	ReplayPaused  ReplayStatus = "paused"
)

type Replay struct {
	mu       sync.Mutex
	actions  []RecordedAction
	status   ReplayStatus
	stopChan chan struct{}

	onLog func(string)
}

func NewReplay() *Replay {
	return &Replay{
		status:   ReplayStopped,
		stopChan: make(chan struct{}),
	}
}

func (rp *Replay) SetOnLog(fn func(string)) {
	rp.onLog = fn
}

func (rp *Replay) log(msg string) {
	if rp.onLog != nil {
		rp.onLog(msg)
	}
}

func (rp *Replay) SetActions(actions []RecordedAction) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	rp.actions = actions
}

func (rp *Replay) GetStatus() ReplayStatus {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	return rp.status
}

func (rp *Replay) setStatus(s ReplayStatus) {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	rp.status = s
}

func (rp *Replay) IsStopped() bool {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	return rp.status == ReplayStopped
}

func (rp *Replay) IsPaused() bool {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	return rp.status == ReplayPaused
}

func (rp *Replay) Start() {
	rp.setStatus(ReplayRunning)
}

func (rp *Replay) Pause() {
	rp.setStatus(ReplayPaused)
}

func (rp *Replay) Resume() {
	rp.setStatus(ReplayRunning)
}

func (rp *Replay) Stop() {
	rp.mu.Lock()
	if rp.status == ReplayStopped {
		rp.mu.Unlock()
		return
	}
	rp.status = ReplayStopped
	rp.stopChan = make(chan struct{})
	rp.mu.Unlock()
}

func (rp *Replay) RunOnce() {
	actions := rp.getActions()
	if len(actions) == 0 {
		return
	}

	var lastOffset float64

	for _, a := range actions {
		if rp.IsStopped() {
			return
		}

		rp.waitForResume()

		if rp.IsStopped() {
			return
		}

		delta := a.TimeOffset - lastOffset
		if delta > 0 {
			rp.sleepUntil(delta, lastOffset)
		}
		lastOffset = a.TimeOffset

		if rp.IsStopped() {
			return
		}

		rp.executeAction(a)
	}

	rp.log(fmt.Sprintf("回放完成，共 %d 个动作", len(actions)))
}

func (rp *Replay) getActions() []RecordedAction {
	rp.mu.Lock()
	defer rp.mu.Unlock()
	result := make([]RecordedAction, len(rp.actions))
	copy(result, rp.actions)
	return result
}

func (rp *Replay) waitForResume() {
	for rp.IsPaused() && !rp.IsStopped() {
		time.Sleep(100 * time.Millisecond)
	}
}

func (rp *Replay) sleepUntil(delta float64, lastOffset float64) {
	elapsed := time.Duration(0)
	target := time.Duration(delta * float64(time.Second))
	step := 50 * time.Millisecond

	for elapsed < target {
		if rp.IsStopped() {
			return
		}
		if rp.IsPaused() {
			rp.waitForResume()
			if rp.IsStopped() {
				return
			}
		}
		time.Sleep(step)
		elapsed += step
	}
}

func (rp *Replay) executeAction(a RecordedAction) {
	switch a.Type {
	case ActionMove:
		robotgo.Move(a.X, a.Y)
		rp.log(fmt.Sprintf("回放 → 移动到 (%d, %d)", a.X, a.Y))
	case ActionClick:
		robotgo.Move(a.X, a.Y)
		time.Sleep(30 * time.Millisecond)
		robotgo.Click(a.Button)
		rp.log(fmt.Sprintf("回放 → %s点击 (%d, %d)", buttonLabel(a.Button), a.X, a.Y))
	case ActionScroll:
		robotgo.ScrollDir(a.Amount, a.Direction)
		rp.log(fmt.Sprintf("回放 → 滚动 %s %d格", directionLabel(a.Direction), a.Amount))
	case ActionKey:
		robotgo.TypeStr(a.KeyChar)
		rp.log(fmt.Sprintf("回放 → 按键: %s", a.KeyChar))
	}
}

func buttonLabel(b string) string {
	switch b {
	case "left":
		return "左键"
	case "right":
		return "右键"
	case "center":
		return "中键"
	default:
		return b
	}
}

func directionLabel(d string) string {
	switch d {
	case "up":
		return "向上"
	case "down":
		return "向下"
	case "left":
		return "向左"
	case "right":
		return "向右"
	default:
		return d
	}
}
