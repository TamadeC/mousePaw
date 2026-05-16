package recorder

import (
	"fmt"
	"math"
	"sync"
	"time"

	hook "github.com/robotn/gohook"
)

type RecorderStatus string

const (
	RecorderIdle      RecorderStatus = "idle"
	RecorderRecording RecorderStatus = "recording"
)

var mouseButtonNames = map[uint16]string{
	1: "left",
	2: "right",
	3: "center",
}

type RecorderStatusInfo struct {
	Status   RecorderStatus `json:"status"`
	Count    int            `json:"count"`
	Duration float64        `json:"duration"`
}

type Recorder struct {
	mu      sync.Mutex
	status  RecorderStatus
	actions []RecordedAction

	startTime    time.Time
	lastMoveTime time.Time
	lastMoveX    int
	lastMoveY    int

	onChange func(RecorderStatusInfo)
}

func NewRecorder() *Recorder {
	return &Recorder{
		status: RecorderIdle,
	}
}

func (r *Recorder) SetOnChange(fn func(RecorderStatusInfo)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.onChange = fn
}

func (r *Recorder) GetStatus() RecorderStatusInfo {
	r.mu.Lock()
	defer r.mu.Unlock()

	duration := float64(0)
	if r.status == RecorderRecording {
		duration = time.Since(r.startTime).Seconds()
	} else if len(r.actions) > 0 {
		duration = r.actions[len(r.actions)-1].TimeOffset
	}

	return RecorderStatusInfo{
		Status:   r.status,
		Count:    len(r.actions),
		Duration: duration,
	}
}

func (r *Recorder) Start() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.status = RecorderRecording
	r.actions = nil
	r.startTime = time.Now()
	r.lastMoveTime = time.Time{}
	r.lastMoveX = 0
	r.lastMoveY = 0

	r.registerHooks()
}

func (r *Recorder) Stop() *Recording {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != RecorderRecording {
		return nil
	}

	r.status = RecorderIdle

	duration := float64(0)
	if len(r.actions) > 0 {
		duration = r.actions[len(r.actions)-1].TimeOffset
	}

	rec := &Recording{
		Actions:   r.actions,
		Duration:  duration,
		CreatedAt: time.Now(),
	}

	r.emitChange()
	return rec
}

func (r *Recorder) GetRecording() *Recording {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.actions) == 0 {
		return nil
	}

	duration := float64(0)
	if len(r.actions) > 0 {
		duration = r.actions[len(r.actions)-1].TimeOffset
	}
	if r.status == RecorderRecording {
		duration = time.Since(r.startTime).Seconds()
	}

	return &Recording{
		Actions:   r.actions,
		Duration:  duration,
		CreatedAt: time.Now(),
	}
}

func (r *Recorder) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.actions = nil
	r.emitChange()
}

func (r *Recorder) registerHooks() {
	hook.Register(hook.MouseDown, []string{}, func(e hook.Event) {
		r.handleMouseDown(e)
	})
	hook.Register(hook.MouseMove, []string{}, func(e hook.Event) {
		r.handleMouseMove(e)
	})
	hook.Register(hook.MouseWheel, []string{}, func(e hook.Event) {
		r.handleMouseWheel(e)
	})
	hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
		r.handleKeyDown(e)
	})
}

func (r *Recorder) handleMouseDown(e hook.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != RecorderRecording {
		return
	}

	btnName, ok := mouseButtonNames[e.Button]
	if !ok {
		return
	}

	action := RecordedAction{
		TimeOffset: time.Since(r.startTime).Seconds(),
		Type:       ActionClick,
		X:          int(e.X),
		Y:          int(e.Y),
		Button:     btnName,
	}

	r.actions = append(r.actions, action)
	r.emitChangeUnsafe()
}

func (r *Recorder) handleMouseMove(e hook.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != RecorderRecording {
		return
	}

	x := int(e.X)
	y := int(e.Y)
	now := time.Now()

	if !r.lastMoveTime.IsZero() {
		dt := now.Sub(r.lastMoveTime).Milliseconds()
		dx := math.Abs(float64(x - r.lastMoveX))
		dy := math.Abs(float64(y - r.lastMoveY))
		if dt < 200 && dx < 20 && dy < 20 {
			return
		}
	}

	r.lastMoveTime = now
	r.lastMoveX = x
	r.lastMoveY = y

	action := RecordedAction{
		TimeOffset: time.Since(r.startTime).Seconds(),
		Type:       ActionMove,
		X:          x,
		Y:          y,
	}

	r.actions = append(r.actions, action)
	r.emitChangeUnsafe()
}

func (r *Recorder) handleMouseWheel(e hook.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != RecorderRecording {
		return
	}

	dir := "up"
	if e.Rotation < 0 {
		dir = "down"
	}

	amount := int(e.Amount)
	if amount == 0 {
		amount = 3
	}

	action := RecordedAction{
		TimeOffset: time.Since(r.startTime).Seconds(),
		Type:       ActionScroll,
		Direction:  dir,
		Amount:     amount,
	}

	r.actions = append(r.actions, action)
	r.emitChangeUnsafe()
}

func (r *Recorder) handleKeyDown(e hook.Event) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.status != RecorderRecording {
		return
	}

	kc := string(e.Keychar)
	if kc == "" || kc == "\x00" || kc == "\uffff" {
		return
	}

	action := RecordedAction{
		TimeOffset: time.Since(r.startTime).Seconds(),
		Type:       ActionKey,
		KeyChar:    kc,
	}

	r.actions = append(r.actions, action)
	r.emitChangeUnsafe()
}

func (r *Recorder) emitChange() {
	if r.onChange == nil {
		return
	}

	duration := float64(0)
	if r.status == RecorderRecording {
		duration = time.Since(r.startTime).Seconds()
	} else if len(r.actions) > 0 {
		duration = r.actions[len(r.actions)-1].TimeOffset
	}

	r.onChange(RecorderStatusInfo{
		Status:   r.status,
		Count:    len(r.actions),
		Duration: duration,
	})
}

func (r *Recorder) emitChangeUnsafe() {
	if r.onChange == nil {
		return
	}

	duration := float64(0)
	if r.status == RecorderRecording {
		duration = time.Since(r.startTime).Seconds()
	} else if len(r.actions) > 0 {
		duration = r.actions[len(r.actions)-1].TimeOffset
	}

	r.onChange(RecorderStatusInfo{
		Status:   r.status,
		Count:    len(r.actions),
		Duration: duration,
	})
}

func (r *Recorder) Actions() []RecordedAction {
	r.mu.Lock()
	defer r.mu.Unlock()
	result := make([]RecordedAction, len(r.actions))
	copy(result, r.actions)
	return result
}

func (r *Recorder) ActionCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.actions)
}

func ActionTypeLabel(t ActionType) string {
	switch t {
	case ActionMove:
		return "鼠标移动"
	case ActionClick:
		return "鼠标点击"
	case ActionScroll:
		return "滚轮滚动"
	case ActionKey:
		return "键盘输入"
	default:
		return string(t)
	}
}

func (a RecordedAction) Detail() string {
	switch a.Type {
	case ActionMove:
		return fmt.Sprintf("(%d, %d)", a.X, a.Y)
	case ActionClick:
		btnName := map[string]string{
			"left": "左键", "right": "右键", "center": "中键",
		}[a.Button]
		if btnName == "" {
			btnName = a.Button
		}
		return fmt.Sprintf("%s at (%d, %d)", btnName, a.X, a.Y)
	case ActionScroll:
		dirName := map[string]string{
			"up": "向上", "down": "向下", "left": "向左", "right": "向右",
		}[a.Direction]
		if dirName == "" {
			dirName = a.Direction
		}
		return fmt.Sprintf("%s %d格", dirName, a.Amount)
	case ActionKey:
		return fmt.Sprintf("按键: %s", a.KeyChar)
	default:
		return ""
	}
}
