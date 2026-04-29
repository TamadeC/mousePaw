# MousePaw - 鼠标自动化工具

一款 Windows 桌面鼠标自动化工具，可以自动执行鼠标移动、点击和滚轮操作。基于 [Wails v2](https://wails.io/) 框架构建（Go 后端 + React/TypeScript 前端）。

[English Documentation](README.md)

## 功能特性

- **鼠标移动**：按可配置的时间间隔（1-60秒）自动将鼠标移动到随机位置
- **鼠标点击**：按可调间隔（0.5-30秒）执行左键/右键/中键点击，支持重复次数（1-10次）
- **鼠标滚轮**：按可配置间隔（1-30秒）在任意方向（上/下/左/右）滚动，滚动量（1-20）
- **全局热键**：F6 开始，F7 停止（即使窗口未聚焦也能使用）
- **开机自启**：可选的 Windows 启动集成
- **实时日志**：在内置日志查看器中查看自动化活动
- **配置持久化**：设置自动保存，下次启动时恢复

## 截图

*在此处添加应用程序界面截图*

## 安装说明

### 系统要求

- Windows 10/11
- [Go 1.25+](https://go.dev/dl/)
- [Node.js 16+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 克隆并运行

```bash
git clone https://github.com/yourusername/mousePaw.git
cd mousePaw
wails dev
```

### 构建可执行文件

```bash
wails build
```

可执行文件将生成在 `build/bin/` 目录中。

## 使用方法

1. **启动应用程序** - 运行 `mousepaw.exe` 或使用 `wails dev` 进行开发
2. **配置自动化** - 使用选项卡启用/禁用不同的鼠标操作：
   - **鼠标移动**：设置时间间隔并启用随机光标移动
   - **鼠标点击**：配置点击类型、时间间隔和重复次数
   - **鼠标滚轮**：设置滚动方向、时间间隔和滚动量
3. **开始自动化** - 点击"开始"按钮或按 **F6**
4. **停止自动化** - 点击"停止"按钮或按 **F7**
5. **查看日志** - 切换到"执行日志"选项卡查看活动记录

## 配置说明

配置存储在可执行文件旁边的 `mousepaw_config.json` 中：

```json
{
  "move": {
    "enabled": true,
    "interval": 5,
    "random": true
  },
  "click": {
    "enabled": false,
    "interval": 2,
    "button": "left",
    "count": 1
  },
  "scroll": {
    "enabled": false,
    "interval": 3,
    "direction": "up",
    "amount": 3
  },
  "auto_start": false,
  "minimize_to_tray": true
}
```

## 快捷键

| 按键 | 功能 |
|------|------|
| **F6** | 开始自动化（全局） |
| **F7** | 停止自动化（全局） |

## 项目结构

```
mousePaw/
├── main.go                 # 应用程序入口点
├── app.go                  # 核心 App 结构体与 Wails 绑定
├── icon.go                 # 程序化图标生成
├── pkg/
│   ├── autostart/          # Windows 注册表开机自启
│   ├── config/             # 配置管理
│   ├── engine/             # 鼠标自动化引擎
│   └── log/                # 内存 + 文件日志记录器
├── frontend/               # React/TypeScript 用户界面
│   ├── src/
│   │   ├── App.tsx         # 主 UI 组件
│   │   └── ...
│   └── wailsjs/            # 自动生成的 Wails 绑定
└── build/                  # 构建资源和输出
```

## 开发指南

### 实时开发模式

```bash
wails dev
```

此命令启动：
- Vite 开发服务器，支持前端热重载
- 开发服务器位于 `http://localhost:34115`，可在浏览器中调试并访问 Go 方法

### 生产环境构建

```bash
wails build
```

## 依赖项

### Go 后端
- [Wails v2](https://wails.io/) - 桌面应用程序框架
- [robotgo](https://github.com/go-vgo/robotgo) - 鼠标/键盘自动化
- [gohook](https://github.com/robotn/gohook) - 全局键盘/鼠标事件钩子
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) - Windows 注册表访问

### 前端
- React 18
- TypeScript
- Tailwind CSS
- Vite

## 许可证

MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。