# MousePaw - 银河麒麟系统构建脚本 (Windows PowerShell)

$ErrorActionPreference = "Stop"

$APP_NAME = "MousePaw"
$OUTPUT_DIR = ".\dist\kylin"
$DOCKER_IMAGE = "mousepaw-kylin-builder"

function Write-Info {
    param([string]$Message)
    Write-Host "[INFO] $Message" -ForegroundColor Green
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[WARN] $Message" -ForegroundColor Yellow
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "[ERROR] $Message" -ForegroundColor Red
}

function Check-Docker {
    try {
        docker version | Out-Null
        Write-Info "Docker 已就绪"
    }
    catch {
        Write-Error-Custom "Docker 未安装，请先安装 Docker Desktop"
        exit 1
    }
}

function Check-QEMU {
    try {
        docker run --rm --platform linux/arm64 arm64v8/debian:bookworm-slim echo "QEMU OK" 2>&1 | Out-Null
        Write-Info "QEMU 多架构支持已就绪"
    }
    catch {
        Write-Warn "QEMU 多架构支持未启用，尝试启用..."
        docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
    }
}

function Build-Frontend {
    Write-Info "构建前端资源..."
    Push-Location frontend
    if (Test-Path "package.json") {
        npm install
        npm run build
    }
    Pop-Location
    Write-Info "前端构建完成"
}

function Build-DockerImage {
    Write-Info "构建 Docker 镜像 ($DOCKER_IMAGE)..."
    docker build -f Dockerfile.kylin -t $DOCKER_IMAGE .
    Write-Info "Docker 镜像构建完成"
}

function Extract-Binary {
    Write-Info "提取二进制文件..."
    New-Item -ItemType Directory -Force -Path $OUTPUT_DIR | Out-Null

    $CONTAINER_ID = docker create $DOCKER_IMAGE
    docker cp "${CONTAINER_ID}:/app/MousePaw" "$OUTPUT_DIR/MousePaw"
    docker cp "${CONTAINER_ID}:/app/frontend/dist" "$OUTPUT_DIR/frontend/dist"
    docker rm $CONTAINER_ID

    Write-Info "二进制文件已提取到 $OUTPUT_DIR"
}

function Create-InstallScript {
    Write-Info "生成安装脚本..."
    
    $installScript = @'
#!/bin/bash

set -e

APP_NAME="MousePaw"
INSTALL_DIR="/opt/mousepaw"
DESKTOP_FILE="/usr/share/applications/mousepaw.desktop"

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

if [ "$EUID" -ne 0 ]; then
    log_error "请使用 sudo 运行此脚本"
    exit 1
fi

log_info "安装 ${APP_NAME}..."

apt-get update
apt-get install -y libgtk-3-0 libwebkit2gtk-4.0-37 libx11-6 libxtst6 libxkbcommon0

mkdir -p ${INSTALL_DIR}
cp -r * ${INSTALL_DIR}/
chmod +x ${INSTALL_DIR}/MousePaw

cat > ${DESKTOP_FILE} << DESKTOP
[Desktop Entry]
Type=Application
Name=MousePaw
Comment=鼠标防休眠工具
Exec=${INSTALL_DIR}/MousePaw
Icon=${INSTALL_DIR}/frontend/dist/favicon.ico
Terminal=false
Categories=Utility;
DESKTOP

log_info "安装完成！"
log_info "可以通过应用菜单或运行 ${INSTALL_DIR}/MousePaw 启动"
'@

    $installScript | Out-File -FilePath "$OUTPUT_DIR/install.sh" -Encoding UTF8
    Write-Info "安装脚本已生成"
}

function Create-Readme {
    Write-Info "生成说明文档..."
    
    $readme = @'
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
'@

    $readme | Out-File -FilePath "$OUTPUT_DIR/README.md" -Encoding UTF8
    Write-Info "说明文档已生成"
}

function Main {
    Write-Info "开始构建银河麒麟版本..."
    
    Check-Docker
    Check-QEMU
    Build-Frontend
    Build-DockerImage
    Extract-Binary
    Create-InstallScript
    Create-Readme
    
    Write-Info "构建完成！"
    Write-Info "输出目录: $OUTPUT_DIR"
    Write-Info ""
    Write-Info "将 $OUTPUT_DIR 目录复制到银河麒麟系统，运行 install.sh 即可安装"
}

Main
