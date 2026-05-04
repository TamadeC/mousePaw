#!/bin/bash

set -e

APP_NAME="MousePaw"
OUTPUT_DIR="./dist/kylin"
DOCKER_IMAGE="mousepaw-kylin-builder"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker 未安装，请先安装 Docker Desktop"
        exit 1
    fi
    log_info "Docker 已就绪"
}

check_qemu() {
    if ! docker run --rm --platform linux/arm64 arm64v8/debian:bookworm-slim echo "QEMU OK" &> /dev/null; then
        log_warn "QEMU 多架构支持未启用，尝试启用..."
        docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
    fi
    log_info "QEMU 多架构支持已就绪"
}

build_frontend() {
    log_info "构建前端资源..."
    cd frontend
    if [ -f "package.json" ]; then
        npm install
        npm run build
    fi
    cd ..
    log_info "前端构建完成"
}

build_docker_image() {
    log_info "构建 Docker 镜像 (${DOCKER_IMAGE})..."
    docker build -f Dockerfile.kylin -t ${DOCKER_IMAGE} .
    log_info "Docker 镜像构建完成"
}

extract_binary() {
    log_info "提取二进制文件..."
    mkdir -p ${OUTPUT_DIR}

    CONTAINER_ID=$(docker create ${DOCKER_IMAGE})
    docker cp ${CONTAINER_ID}:/app/MousePaw ${OUTPUT_DIR}/MousePaw
    docker cp ${CONTAINER_ID}:/app/frontend/dist ${OUTPUT_DIR}/frontend/dist
    docker rm ${CONTAINER_ID}

    log_info "二进制文件已提取到 ${OUTPUT_DIR}"
}

create_install_script() {
    log_info "生成安装脚本..."
    cat > ${OUTPUT_DIR}/install.sh << 'EOF'
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
EOF

    chmod +x ${OUTPUT_DIR}/install.sh
    log_info "安装脚本已生成"
}

create_readme() {
    log_info "生成说明文档..."
    cat > ${OUTPUT_DIR}/README.md << 'EOF'
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
EOF

    log_info "说明文档已生成"
}

main() {
    log_info "开始构建银河麒麟版本..."
    
    check_docker
    check_qemu
    build_frontend
    build_docker_image
    extract_binary
    create_install_script
    create_readme
    
    log_info "构建完成！"
    log_info "输出目录: ${OUTPUT_DIR}"
    log_info ""
    log_info "将 ${OUTPUT_DIR} 目录复制到银河麒麟系统，运行 install.sh 即可安装"
}

main "$@"
