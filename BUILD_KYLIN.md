# MousePaw - 银河麒麟系统打包指南

## 概述

本文档说明如何将 MousePaw 打包为银河麒麟系统（ARM64架构）可用的程序。

## 前置条件

### Windows 环境

1. **Docker Desktop**
   - 下载地址：https://www.docker.com/products/docker-desktop
   - 安装后确保启用多架构支持

2. **Node.js 和 npm**
   - 用于构建前端资源
   - 下载地址：https://nodejs.org/

### 银河麒麟环境

- 银河麒麟 V10 SP1 或更高版本
- ARM64 架构（飞腾/鲲鹏处理器）
- 图形桌面环境

## 构建步骤

### 方法一：使用 PowerShell 脚本（推荐）

```powershell
# 在项目根目录执行
.\build-kylin.ps1
```

### 方法二：使用 Bash 脚本

```bash
# 在项目根目录执行（需要 Git Bash 或 WSL）
bash build-kylin.sh
```

### 方法三：手动构建

1. **构建前端**
```bash
cd frontend
npm install
npm run build
cd ..
```

2. **构建 Docker 镜像**
```bash
docker build -f Dockerfile.kylin -t mousepaw-kylin-builder .
```

3. **提取二进制文件**
```bash
mkdir -p dist/kylin
CONTAINER_ID=$(docker create mousepaw-kylin-builder)
docker cp $CONTAINER_ID:/app/MousePaw dist/kylin/
docker cp $CONTAINER_ID:/app/frontend/dist dist/kylin/frontend/
docker rm $CONTAINER_ID
```

## 输出目录结构

```
dist/kylin/
├── MousePaw          # 主程序
├── install.sh        # 安装脚本
├── README.md         # 说明文档
└── frontend/
    └── dist/         # 前端资源
```

## 在银河麒麟上安装

### 方法一：使用安装脚本（推荐）

1. 将 `dist/kylin` 目录复制到银河麒麟系统

2. 执行安装脚本：
```bash
cd /path/to/kylin
sudo ./install.sh
```

3. 启动程序：
```bash
/opt/mousepaw/MousePaw
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

## 开机自启配置

### 方法一：程序内设置

在程序设置中启用"开机自启"功能。

### 方法二：手动配置

```bash
mkdir -p ~/.config/autostart
cat > ~/.config/autostart/mousepaw.desktop << 'EOF'
[Desktop Entry]
Type=Application
Name=MousePaw
Exec=/opt/mousepaw/MousePaw
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
EOF
```

## 常见问题

### Q: 程序无法启动

A: 检查是否已安装所有依赖：
```bash
sudo apt-get install -y libgtk-3-0 libwebkit2gtk-4.0-37 libx11-6 libxtst6 libxkbcommon0
```

### Q: 界面显示异常

A: 确保在图形桌面环境下运行，而不是 SSH 终端。

### Q: 快捷键不工作

A: 确保程序有权限监听全局键盘事件，可能需要以 root 权限运行或配置 udev 规则。

### Q: 如何查看日志

A: 日志文件位于：
```bash
~/.mousepaw/mousepaw.log
```

## 技术说明

### 依赖库说明

| 库 | 用途 |
|---|------|
| libgtk-3-0 | GTK 图形库 |
| libwebkit2gtk-4.0-37 | Web 渲染引擎（Wails 前端） |
| libx11-6 | X11 图形系统 |
| libxtst6 | X11 测试扩展（模拟输入） |
| libxkbcommon0 | 键盘处理 |

### 架构说明

- 程序使用 Wails v2 框架
- 前端使用 Web 技术（HTML/CSS/JS）
- 后端使用 Go 语言
- 使用 robotgo 实现鼠标控制
- 使用 gohook 实现全局快捷键

## 开发说明

如需修改代码并重新构建：

1. 修改代码后重新运行构建脚本
2. Docker 镜像会自动重新构建
3. 提取新的二进制文件到 `dist/kylin`
4. 重新部署到银河麒麟系统

## 许可证

本程序遵循项目原有的许可证协议。
