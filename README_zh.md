# MousePaw - 鼠标自动化工具

一款跨平台桌面应用程序，通过自动化鼠标移动、点击、滚轮和键盘输入来防止系统空闲休眠。基于 [Wails v2](https://wails.io/) 框架构建（Go 后端 + React/TypeScript 前端）。

[English Documentation](README.md)

## 功能特性

- **鼠标移动**：按可配置的时间间隔（1--60秒）将光标移动到屏幕随机位置，移动轨迹平滑过渡
- **鼠标点击**：按可调间隔（0.5--30秒）执行左键/右键/中键点击，支持重复次数（1--10次）
- **鼠标滚轮**：按可配置间隔（1--30秒）在任意方向（上/下/左/右）滚动，支持自定义滚动量（1--20格）
- **键盘输入**：按可调间隔（0.5--10秒）重复输入配置的文本内容
- **全局热键**：可自定义的系统级快捷键，分别控制启动 / 停止 / 暂停-恢复（默认：Ctrl+F6 / Ctrl+F7 / Ctrl+F8）
- **操作录制**：录制鼠标移动、点击、滚轮和键盘输入，记录时间戳；支持保存/加载/管理多个录制文件
- **录制回放**：按可配置的时间间隔（1--300秒）回放录制的操作序列，支持循环回放
- **系统托盘**：关闭窗口时最小化到托盘，右键菜单提供"显示窗口"和"退出"选项
- **开机自启**：可选的系统启动集成（Windows 注册表 / Linux `.desktop` 自启动）
- **实时日志**：内存环形缓冲区 + 文件日志记录器，通过事件机制实时推送到界面
- **配置持久化**：设置自动保存到可执行文件旁的 `mousepaw_config.json`，支持旧配置向后兼容迁移

## 截图

*在此处添加应用程序界面截图*

## 安装说明

### 系统要求

- Windows 10/11（主要），或带 GTK/WebKit2 的 Linux 系统
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

### 银河麒麟 Linux（ARM64）

银河麒麟系统（ARM64 架构）的构建方法请参阅 [BUILD_KYLIN.md](BUILD_KYLIN.md)，其中包含基于 Docker 的交叉编译和部署说明。

## 使用方法

1. **启动应用程序** - 运行 `MousePaw.exe` 或使用 `wails dev` 进行开发调试
2. **配置自动化** - 选择操作模式并调整参数：
   - **鼠标移动**：间隔（1--60秒）、随机定位开关
   - **鼠标点击**：按键类型（左/右/中）、间隔（0.5--30秒）、重复次数（1--10）
   - **滚轮滚动**：方向（上/下/左/右）、间隔（1--30秒）、滚动量（1--20）
   - **键盘输入**：间隔（0.5--10秒）、要输入的文本内容
   - **录制回放**：间隔（1--300秒）、循环开关、选择录制文件
3. **启动自动化** - 点击"启动"按钮或按下配置的热键（默认：**Ctrl+F6**）
4. **暂停/恢复** - 点击"暂停"按钮或按下配置的热键（默认：**Ctrl+F8**）
5. **停止自动化** - 点击"停止"按钮或按下配置的热键（默认：**Ctrl+F7**）
6. **查看日志** - 切换到"执行日志"选项卡查看实时活动记录
7. **系统托盘** - 关闭窗口后最小化到托盘，右键托盘图标可选择显示窗口或退出程序
8. **录制** - 切换到"操作录制"选项卡，可录制、保存和管理操作序列用于回放

## 配置说明

配置存储在可执行文件旁的 `mousepaw_config.json` 中：

```json
{
  "operation_type": "move",
  "move_interval": 5.0,
  "move_random": true,
  "click_interval": 3.0,
  "click_type": "left",
  "click_count": 1,
  "scroll_interval": 5.0,
  "scroll_dir": "down",
  "scroll_amount": 3,
  "type_interval": 1.0,
  "type_text": "",
  "replay_interval": 30.0,
  "replay_repeat": false,
  "replay_file": "",
  "auto_start": false,
  "minimize_to_tray": true,
  "hotkeys": {
    "start": "ctrl+f6",
    "stop": "ctrl+f7",
    "pause": "ctrl+f8"
  }
}
```

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `operation_type` | string | `"move"` | 当前操作模式：`move`（移动）、`click`（点击）、`scroll`（滚动）、`type`（输入）、`replay`（回放） |
| `move_interval` | number | `5.0` | 鼠标移动间隔秒数（1--60） |
| `move_random` | boolean | `true` | 是否移动到随机位置 |
| `click_interval` | number | `3.0` | 点击间隔秒数（0.5--30） |
| `click_type` | string | `"left"` | 鼠标按键：`left`（左键）、`right`（右键）、`middle`（中键） |
| `click_count` | integer | `1` | 每次触发时的连击次数（1--10） |
| `scroll_interval` | number | `5.0` | 滚动间隔秒数（1--30） |
| `scroll_dir` | string | `"down"` | 滚动方向：`up`、`down`、`left`、`right` |
| `scroll_amount` | integer | `3` | 每次滚动的格数（1--20） |
| `type_interval` | number | `1.0` | 键盘输入间隔秒数（0.5--10） |
| `type_text` | string | `""` | 每次输入时键入的文本 |
| `replay_interval` | number | `30.0` | 回放周期间隔秒数（1--300） |
| `replay_repeat` | boolean | `false` | 启用后循环连续回放 |
| `replay_file` | string | `""` | 要回放的录制文件名 |
| `auto_start` | boolean | `false` | 是否开机自启 |
| `minimize_to_tray` | boolean | `true` | 关闭窗口时是否最小化到托盘 |
| `hotkeys.start` | string | `"ctrl+f6"` | 启动自动化热键 |
| `hotkeys.stop` | string | `"ctrl+f7"` | 停止自动化热键 |
| `hotkeys.pause` | string | `"ctrl+f8"` | 暂停/恢复自动化热键 |

## 快捷键

默认全局热键（可在"系统设置"选项卡中自定义）：

| 快捷键 | 功能 |
|--------|------|
| **Ctrl+F6** | 启动自动化 |
| **Ctrl+F7** | 停止自动化 |
| **Ctrl+F8** | 暂停 / 恢复自动化 |

热键为全局生效，即使 MousePaw 窗口未聚焦也能使用。修改热键后需重启程序生效。

## 项目结构

```
mousePaw/
├── main.go                     # 应用程序入口，Wails 窗口配置
├── app.go                      # 核心 App 结构体，Wails 绑定，热键监听
├── systray.go                  # 系统托盘集成
├── icon.go                     # 程序化 PNG 图标生成
├── wails.json                  # Wails 框架配置
│
├── pkg/
│   ├── autostart/
│   │   ├── autostart_windows.go    # Windows 注册表自启动（HKCU Run 键）
│   │   └── autostart_linux.go      # Linux .config/autostart .desktop 文件
│   ├── config/
│   │   └── config.go               # 配置结构体，JSON 加载/保存，向后兼容
│   ├── engine/
│   │   └── mouse.go                # 自动化引擎（启动/停止/暂停/恢复）
│   ├── log/
│   │   └── logger.go               # 内存环形缓冲区 + 文件日志记录器
│   └── recorder/
│       ├── recording.go            # 操作数据结构 + JSON I/O
│       ├── recorder.go             # 全局事件录制（基于 gohook）
│       └── replay.go               # 录制回放引擎
│
├── frontend/
│   ├── package.json                # React 18、Vite、Tailwind CSS、TypeScript
│   ├── vite.config.ts              # Vite 构建配置
│   ├── tailwind.config.js          # Tailwind CSS 暗色主题
│   ├── src/
│   │   ├── main.tsx                # React 入口
│   │   ├── App.tsx                 # 主界面：选项卡、控件、设置、日志
│   │   ├── style.css               # Tailwind 指令 + 自定义样式
│   │   └── assets/                 # 字体和图片资源
│   └── wailsjs/                    # 自动生成的 Wails JS 绑定
│
├── build/                          # 构建资源（图标、清单、NSIS 安装程序）
├── dist/                           # 银河麒麟 Linux 分发输出
├── Dockerfile.kylin                # 麒麟 ARM64 交叉编译 Docker 镜像
├── build-kylin.sh                  # Linux/macOS 麒麟构建脚本
├── build-kylin.ps1                 # Windows 麒麟构建脚本
├── build-kylin-local.sh            # 麒麟原生构建脚本
├── BUILD_KYLIN.md                  # 麒麟构建和部署指南
└── LICENSE                         # MIT 许可证
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
wails build               # Windows（默认平台）
```

银河麒麟 Linux ARM64 构建：

```bash
.\build-kylin.ps1         # Windows
bash build-kylin.sh       # Linux/macOS
bash build-kylin-local.sh # 麒麟原生构建
```

## 架构说明

应用程序采用分层架构设计：

```
┌────────────────────────────────────────────────┐
│            前端 (React 18 / TypeScript)         │
│  App.tsx → 选项卡：操作设置、操作录制、系统设置、  │
│            执行日志                              │
│  Tailwind CSS 暗色主题                          │
└──────────────────┬─────────────────────────────┘
                   │ Wails IPC（JSON 绑定 + 事件）
┌──────────────────▼─────────────────────────────┐
│               App (app.go)                      │
│  配置、引擎、录制器、日志管理                      │
│  方法：GetConfig, Start/Stop, StartRecording    │
│        StopRecording, LoadRecording 等          │
└──┬──────┬──────────┬──────────┬───────────────┘
   │      │          │          │
   ▼      ▼          ▼          ▼
┌──────┐ ┌────────┐ ┌──────┐ ┌──────┐
│ 配置 │ │  引擎  │ │ 录制器│ │ 日志 │
│JSON  │ │ mouse  │ │gohook│ │环形+ │
│I/O   │ │        │ │ 录制 │ │文件  │
└──┬───┘ └───┬────┘ └──┬───┘ └──────┘
   │         │         │
   ▼         ▼         ▼
┌──────┐ ┌───────┐ ┌──────────┐
│自启动│ │robotgo│ │ 录制文件  │
│      │ │+ 回放  │ │ JSON I/O │
└──────┘ └───────┘ └──────────┘
```

## 依赖项

### Go 后端

| 包 | 用途 |
|----|------|
| [Wails v2](https://wails.io/) | 桌面应用框架（Go + WebView） |
| [robotgo](https://github.com/go-vgo/robotgo) | 鼠标/键盘自动化 |
| [gohook](https://github.com/robotn/gohook) | 全局键盘事件钩子 |
| [systray](https://github.com/getlantern/systray) | 系统托盘图标和菜单 |
| [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) | Windows 注册表访问 |

### 前端

- React 18
- TypeScript
- Tailwind CSS 3
- Vite 3

## 许可证

MIT 许可证 - 详情请参阅 [LICENSE](LICENSE) 文件。
