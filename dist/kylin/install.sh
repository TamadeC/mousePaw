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
