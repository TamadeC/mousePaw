package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

type ClickType string

const (
	ClickLeft   ClickType = "left"
	ClickRight  ClickType = "right"
	ClickMiddle ClickType = "middle"
)

type ScrollDirection string

const (
	ScrollUp    ScrollDirection = "up"
	ScrollDown  ScrollDirection = "down"
	ScrollLeft  ScrollDirection = "left"
	ScrollRight ScrollDirection = "right"
)

type Config struct {
	mu sync.RWMutex `json:"-"`

	MoveEnabled    bool    `json:"move_enabled"`
	MoveInterval   float64 `json:"move_interval"`
	MoveRandom     bool    `json:"move_random"`
	ClickEnabled   bool    `json:"click_enabled"`
	ClickInterval  float64 `json:"click_interval"`
	ClickType      ClickType `json:"click_type"`
	ClickCount     int     `json:"click_count"`
	ScrollEnabled  bool    `json:"scroll_enabled"`
	ScrollInterval float64 `json:"scroll_interval"`
	ScrollDir      ScrollDirection `json:"scroll_dir"`
	ScrollAmount   int     `json:"scroll_amount"`
	AutoStart      bool    `json:"auto_start"`
	MinimizeToTray bool    `json:"minimize_to_tray"`
}

var defaultConfig = Config{
	MoveEnabled:    false,
	MoveInterval:   5.0,
	MoveRandom:     true,
	ClickEnabled:   false,
	ClickInterval:  3.0,
	ClickType:      ClickLeft,
	ClickCount:     1,
	ScrollEnabled:  false,
	ScrollInterval: 5.0,
	ScrollDir:      ScrollDown,
	ScrollAmount:   3,
	AutoStart:      false,
	MinimizeToTray: true,
}

func configPath() string {
	exePath, _ := os.Executable()
	return filepath.Join(filepath.Dir(exePath), "mousepaw_config.json")
}

func Load() *Config {
	cfg := &Config{}
	*cfg = defaultConfig

	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	_ = json.Unmarshal(data, cfg)
	return cfg
}

func (c *Config) Save() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(), data, 0644)
}

func (c *Config) GetAll() Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Config{
		MoveEnabled:    c.MoveEnabled,
		MoveInterval:   c.MoveInterval,
		MoveRandom:     c.MoveRandom,
		ClickEnabled:   c.ClickEnabled,
		ClickInterval:  c.ClickInterval,
		ClickType:      c.ClickType,
		ClickCount:     c.ClickCount,
		ScrollEnabled:  c.ScrollEnabled,
		ScrollInterval: c.ScrollInterval,
		ScrollDir:      c.ScrollDir,
		ScrollAmount:   c.ScrollAmount,
		AutoStart:      c.AutoStart,
		MinimizeToTray: c.MinimizeToTray,
	}
}

func (c *Config) Update(cfg Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.MoveEnabled = cfg.MoveEnabled
	c.MoveInterval = cfg.MoveInterval
	c.MoveRandom = cfg.MoveRandom
	c.ClickEnabled = cfg.ClickEnabled
	c.ClickInterval = cfg.ClickInterval
	c.ClickType = cfg.ClickType
	c.ClickCount = cfg.ClickCount
	c.ScrollEnabled = cfg.ScrollEnabled
	c.ScrollInterval = cfg.ScrollInterval
	c.ScrollDir = cfg.ScrollDir
	c.ScrollAmount = cfg.ScrollAmount
	c.AutoStart = cfg.AutoStart
	c.MinimizeToTray = cfg.MinimizeToTray
}
