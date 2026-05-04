package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	appName    = "MousePaw"
	desktopTpl = `[Desktop Entry]
Type=Application
Name=MousePaw
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`
)

func getDesktopFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("获取用户目录失败: %w", err)
	}
	return filepath.Join(homeDir, ".config", "autostart", appName+".desktop"), nil
}

func Enable() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	exePath, err = filepath.Abs(exePath)
	if err != nil {
		return fmt.Errorf("获取绝对路径失败: %w", err)
	}

	desktopPath, err := getDesktopFilePath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(desktopPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	content := fmt.Sprintf(desktopTpl, exePath)
	if err := os.WriteFile(desktopPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("写入桌面文件失败: %w", err)
	}

	return nil
}

func Disable() error {
	desktopPath, err := getDesktopFilePath()
	if err != nil {
		return err
	}

	if err := os.Remove(desktopPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除桌面文件失败: %w", err)
	}

	return nil
}

func IsEnabled() bool {
	desktopPath, err := getDesktopFilePath()
	if err != nil {
		return false
	}

	_, err = os.Stat(desktopPath)
	return err == nil
}
