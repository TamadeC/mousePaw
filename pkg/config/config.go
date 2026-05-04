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

type OperationType string

const (
	OpMove   OperationType = "move"
	OpClick  OperationType = "click"
	OpScroll OperationType = "scroll"
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

	OperationType  OperationType   `json:"operation_type"`
	MoveInterval   float64         `json:"move_interval"`
	MoveRandom     bool            `json:"move_random"`
	ClickInterval  float64         `json:"click_interval"`
	ClickType      ClickType       `json:"click_type"`
	ClickCount     int             `json:"click_count"`
	ScrollInterval float64         `json:"scroll_interval"`
	ScrollDir      ScrollDirection `json:"scroll_dir"`
	ScrollAmount   int             `json:"scroll_amount"`
	AutoStart      bool            `json:"auto_start"`
	MinimizeToTray bool            `json:"minimize_to_tray"`
}

var defaultConfig = Config{
	OperationType:  OpMove,
	MoveInterval:   5.0,
	MoveRandom:     true,
	ClickInterval:  3.0,
	ClickType:      ClickLeft,
	ClickCount:     1,
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

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return cfg
	}

	// 迁移旧配置：检查是否存在旧的 enabled 字段
	if moveEnabled, ok := raw["move_enabled"].(bool); ok && moveEnabled {
		cfg.OperationType = OpMove
	} else if clickEnabled, ok := raw["click_enabled"].(bool); ok && clickEnabled {
		cfg.OperationType = OpClick
	} else if scrollEnabled, ok := raw["scroll_enabled"].(bool); ok && scrollEnabled {
		cfg.OperationType = OpScroll
	}

	// 解析新字段
	if opType, ok := raw["operation_type"].(string); ok {
		cfg.OperationType = OperationType(opType)
	}
	if v, ok := raw["move_interval"].(float64); ok {
		cfg.MoveInterval = v
	}
	if v, ok := raw["move_random"].(bool); ok {
		cfg.MoveRandom = v
	}
	if v, ok := raw["click_interval"].(float64); ok {
		cfg.ClickInterval = v
	}
	if v, ok := raw["click_type"].(string); ok {
		cfg.ClickType = ClickType(v)
	}
	if v, ok := raw["click_count"].(float64); ok {
		cfg.ClickCount = int(v)
	}
	if v, ok := raw["scroll_interval"].(float64); ok {
		cfg.ScrollInterval = v
	}
	if v, ok := raw["scroll_dir"].(string); ok {
		cfg.ScrollDir = ScrollDirection(v)
	}
	if v, ok := raw["scroll_amount"].(float64); ok {
		cfg.ScrollAmount = int(v)
	}
	if v, ok := raw["auto_start"].(bool); ok {
		cfg.AutoStart = v
	}
	if v, ok := raw["minimize_to_tray"].(bool); ok {
		cfg.MinimizeToTray = v
	}

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
		OperationType:  c.OperationType,
		MoveInterval:   c.MoveInterval,
		MoveRandom:     c.MoveRandom,
		ClickInterval:  c.ClickInterval,
		ClickType:      c.ClickType,
		ClickCount:     c.ClickCount,
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
	c.OperationType = cfg.OperationType
	c.MoveInterval = cfg.MoveInterval
	c.MoveRandom = cfg.MoveRandom
	c.ClickInterval = cfg.ClickInterval
	c.ClickType = cfg.ClickType
	c.ClickCount = cfg.ClickCount
	c.ScrollInterval = cfg.ScrollInterval
	c.ScrollDir = cfg.ScrollDir
	c.ScrollAmount = cfg.ScrollAmount
	c.AutoStart = cfg.AutoStart
	c.MinimizeToTray = cfg.MinimizeToTray
}
