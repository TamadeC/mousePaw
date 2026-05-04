#!/bin/bash

set -e

GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

log_info "安装编译依赖..."
sudo apt-get update
sudo apt-get install -y \
    golang \
    gcc \
    pkg-config \
    libgtk-3-dev \
    libwebkit2gtk-4.0-dev \
    libx11-dev \
    libxtst-dev \
    libxkbcommon-dev \
    libxkbcommon-x11-dev \
    libx11-xcb-dev \
    libayatana-appindicator3-dev \
    npm

log_info "安装前端依赖..."
cd frontend
npm install
npm run build
cd ..

log_info "下载 Go 依赖..."
export GOPROXY=https://goproxy.cn,https://goproxy.io,direct
go mod download

log_info "编译程序..."
go build -tags webkit2_4.0 -o MousePaw .

log_info "编译完成！可执行文件: ./MousePaw"
log_info "运行: ./MousePaw"
