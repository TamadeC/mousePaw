package autostart

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

const (
	regPath = `Software\Microsoft\Windows\CurrentVersion\Run`
	appName = "MousePaw"
)

func Enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %w", err)
	}

	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer key.Close()

	err = key.SetStringValue(appName, fmt.Sprintf(`"%s"`, exePath))
	if err != nil {
		return fmt.Errorf("设置注册表值失败: %w", err)
	}
	return nil
}

func Disable() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("打开注册表失败: %w", err)
	}
	defer key.Close()

	err = key.DeleteValue(appName)
	if err != nil && err != registry.ErrNotExist {
		return fmt.Errorf("删除注册表值失败: %w", err)
	}
	return nil
}

func IsEnabled() bool {
	key, err := registry.OpenKey(registry.CURRENT_USER, regPath, registry.QUERY_VALUE)
	if err != nil {
		return false
	}
	defer key.Close()

	val, _, err := key.GetStringValue(appName)
	if err != nil {
		return false
	}

	exePath, _ := os.Executable()
	exePath, _ = filepath.Abs(exePath)
	return val == fmt.Sprintf(`"%s"`, exePath)
}
