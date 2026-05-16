# MousePaw - 银河麒麟系统版本

## 系统要求

- 银河麒麟 V10 SP1 或更高版本
- ARM64 架构（飞腾/鲲鹏处理器）

## 安装方法

### 方法一：使用安装脚本（推荐）

```bash
sudo ./install.sh
```

### 方法二：手动安装

1. 安装依赖：
```bash
sudo apt-get update
sudo apt-get install -y libgtk-3-0 libwebkit2gtk-4.0-37 libx11-6 libxtst6 libxkbcommon0
```

2. 复制文件：
```bash
sudo mkdir -p /opt/mousepaw
sudo cp -r * /opt/mousepaw/
sudo chmod +x /opt/mousepaw/MousePaw
```

3. 运行程序：
```bash
/opt/mousepaw/MousePaw
```

## 开机自启

在程序设置中启用"开机自启"功能，或手动创建桌面文件：

```bash
mkdir -p ~/.config/autostart
cat > ~/.config/autostart/mousepaw.desktop << DESKTOP
[Desktop Entry]
Type=Application
Name=MousePaw
Exec=/opt/mousepaw/MousePaw
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
DESKTOP
```

## 快捷键

- F6: 启动引擎
- F7: 停止引擎

## 问题反馈

如遇到问题，请检查：
1. 是否已安装所有依赖
2. 是否在图形桌面环境下运行
3. 查看日志文件：~/.mousepaw/mousepaw.log
