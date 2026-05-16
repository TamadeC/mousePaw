package recorder

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ActionType string

const (
	ActionMove   ActionType = "move"
	ActionClick  ActionType = "click"
	ActionScroll ActionType = "scroll"
	ActionKey    ActionType = "key"
)

type RecordedAction struct {
	TimeOffset float64    `json:"time_offset"`
	Type       ActionType `json:"type"`

	X int `json:"x,omitempty"`
	Y int `json:"y,omitempty"`

	Button    string `json:"button,omitempty"`
	Direction string `json:"direction,omitempty"`
	Amount    int    `json:"amount,omitempty"`

	KeyChar string `json:"keychar,omitempty"`
}

type Recording struct {
	Actions   []RecordedAction `json:"actions"`
	Duration  float64          `json:"duration"`
	CreatedAt time.Time        `json:"created_at"`
}

func recordingDir() string {
	exePath, _ := os.Executable()
	return filepath.Dir(exePath)
}

func fileName(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		name = "recording"
	}
	if !strings.HasSuffix(name, "_recording.json") {
		name = name + "_recording.json"
	}
	return filepath.Join(recordingDir(), name)
}

func LoadRecording(name string) (*Recording, error) {
	data, err := os.ReadFile(fileName(name))
	if err != nil {
		return nil, err
	}
	var r Recording
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

func SaveRecording(name string, r *Recording) error {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName(name), data, 0644)
}

func ListRecordings() []string {
	dir := recordingDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), "_recording.json") {
			name := strings.TrimSuffix(e.Name(), "_recording.json")
			names = append(names, name)
		}
	}
	return names
}

func DeleteRecording(name string) error {
	return os.Remove(fileName(name))
}
