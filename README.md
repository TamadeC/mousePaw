# MousePaw - Mouse Automation Tool

A cross-platform desktop application for preventing system idle by automating mouse movements, clicks, scrolling, and keyboard input. Built with [Wails v2](https://wails.io/) (Go + React/TypeScript).

[中文文档](README_zh.md)

## Features

- **Mouse Movement**: Automatically move the cursor to random screen positions at configurable intervals (1--60 seconds), with smooth interpolation
- **Mouse Click**: Perform left/right/middle clicks at adjustable intervals (0.5--30 seconds) with repeat counts (1--10)
- **Mouse Scroll**: Scroll in any direction (up/down/left/right) with configurable intervals (1--30 seconds) and amounts (1--20)
- **Keyboard Input**: Type configurable text repeatedly at adjustable intervals (0.5--10 seconds)
- **Global Hotkeys**: Configurable system-wide shortcuts for Start / Stop / Pause-Resume (default: Ctrl+F6 / Ctrl+F7 / Ctrl+F8)
- **Operation Recording**: Record mouse movements, clicks, scrolls, and keyboard input with timestamps; save/load/manage multiple recording files
- **Recording Replay**: Replay recorded operations at configurable intervals (1--300 seconds) with optional continuous looping
- **System Tray**: Minimize to system tray on close, with "Show Window" and "Quit" context menu
- **Auto-start**: Optional OS-level startup integration (Windows Registry / Linux `.desktop` autostart)
- **Real-time Logging**: In-memory ring buffer + file logger with live UI updates via events
- **Configuration Persistence**: Settings auto-saved to `mousepaw_config.json` alongside the executable, with backward-compatible migration

## Screenshots

*Add screenshots of the application UI here*

## Installation

### Prerequisites

- Windows 10/11 (primary), or Linux with GTK/WebKit2
- [Go 1.25+](https://go.dev/dl/)
- [Node.js 16+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation)

### Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone and Run

```bash
git clone https://github.com/yourusername/mousePaw.git
cd mousePaw
wails dev
```

### Build Executable

```bash
wails build
```

The executable will be created in `build/bin/`.

### Kylin Linux (ARM64)

For building on Kylin Linux (ARM64 architecture), see [BUILD_KYLIN.md](BUILD_KYLIN.md) for Docker-based cross-compilation and deployment instructions.

## Usage

1. **Launch the application** - Run `MousePaw.exe` or use `wails dev` for development
2. **Configure automation** - Select an operation mode and adjust its parameters:
   - **Mouse Move**: Interval (1--60s), random positioning toggle
   - **Mouse Click**: Button type (left/right/middle), interval (0.5--30s), repeat count (1--10)
   - **Scroll**: Direction (up/down/left/right), interval (1--30s), amount (1--20)
   - **Keyboard Input**: Interval (0.5--10s), text content to type
   - **Recording Replay**: Interval (1--300s), loop toggle, select a recording file
3. **Start automation** - Click "Start" button or press the configured hotkey (default: **Ctrl+F6**)
4. **Pause/Resume** - Click "Pause" button or press the configured hotkey (default: **Ctrl+F8**)
5. **Stop automation** - Click "Stop" button or press the configured hotkey (default: **Ctrl+F7**)
6. **View logs** - Switch to "Execution Logs" tab to monitor real-time activity
7. **System tray** - Close the window to minimize to tray; right-click the tray icon for Show/Quit options
8. **Recording** - Switch to "Operation Recording" tab to record, save, and manage operation recordings for replay

## Configuration

Configuration is stored in `mousepaw_config.json` next to the executable:

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

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `operation_type` | string | `"move"` | Current operation mode: `move`, `click`, `scroll`, `type`, `replay` |
| `move_interval` | number | `5.0` | Seconds between mouse movements (1--60) |
| `move_random` | boolean | `true` | Move to random screen positions |
| `click_interval` | number | `3.0` | Seconds between clicks (0.5--30) |
| `click_type` | string | `"left"` | Mouse button: `left`, `right`, `middle` |
| `click_count` | integer | `1` | Consecutive clicks per interval (1--10) |
| `scroll_interval` | number | `5.0` | Seconds between scroll actions (1--30) |
| `scroll_dir` | string | `"down"` | Scroll direction: `up`, `down`, `left`, `right` |
| `scroll_amount` | integer | `3` | Scroll lines per action (1--20) |
| `type_interval` | number | `1.0` | Seconds between keyboard inputs (0.5--10) |
| `type_text` | string | `""` | Text typed at each interval |
| `replay_interval` | number | `30.0` | Seconds between replay cycles (1--300) |
| `replay_repeat` | boolean | `false` | Loop replay continuously when enabled |
| `replay_file` | string | `""` | Recording file name to replay |
| `auto_start` | boolean | `false` | Launch on system startup |
| `minimize_to_tray` | boolean | `true` | Minimize to tray instead of closing |
| `hotkeys.start` | string | `"ctrl+f6"` | Start automation hotkey |
| `hotkeys.stop` | string | `"ctrl+f7"` | Stop automation hotkey |
| `hotkeys.pause` | string | `"ctrl+f8"` | Pause/Resume automation hotkey |

## Hotkeys

Default global hotkeys (configurable in System Settings tab):

| Hotkey | Action |
|--------|--------|
| **Ctrl+F6** | Start automation |
| **Ctrl+F7** | Stop automation |
| **Ctrl+F8** | Pause / Resume automation |

Hotkeys work globally even when the MousePaw window is not focused. Changes to hotkeys take effect after restarting the application.

## Project Structure

```
mousePaw/
├── main.go                     # Application entry point, Wails setup, window config
├── app.go                      # Core App struct, Wails bindings, hotkey listener
├── systray.go                  # System tray integration
├── icon.go                     # Programmatic PNG icon generation
├── wails.json                  # Wails framework configuration
│
├── pkg/
│   ├── autostart/
│   │   ├── autostart_windows.go    # Windows Registry auto-start (HKCU Run key)
│   │   └── autostart_linux.go      # Linux .config/autostart .desktop file
│   ├── config/
│   │   └── config.go               # Config struct, JSON load/save, backward compat
│   ├── engine/
│   │   └── mouse.go                # Automation engine (Start/Stop/Pause/Resume)
│   ├── log/
│   │   └── logger.go               # In-memory ring buffer + file logger
│   └── recorder/
│       ├── recording.go            # Action data structures + JSON I/O
│       ├── recorder.go             # Global event recording (gohook-based)
│       └── replay.go               # Recording replay engine
│
├── frontend/
│   ├── package.json                # React 18, Vite, Tailwind CSS, TypeScript
│   ├── vite.config.ts              # Vite build configuration
│   ├── tailwind.config.js          # Tailwind CSS dark theme
│   ├── src/
│   │   ├── main.tsx                # React entry point
│   │   ├── App.tsx                 # Main UI: tabs, controls, settings, logs
│   │   ├── style.css               # Tailwind directives + custom styles
│   │   └── assets/                 # Fonts and images
│   └── wailsjs/                    # Auto-generated Wails JS bindings
│
├── build/                          # Build assets (icons, manifests, NSIS installer)
├── dist/                           # Kylin Linux distribution output
├── Dockerfile.kylin                # Docker build for Kylin ARM64 cross-compilation
├── build-kylin.sh                  # Linux/macOS Kylin build script
├── build-kylin.ps1                 # Windows Kylin build script
├── build-kylin-local.sh            # Native Kylin build script
├── BUILD_KYLIN.md                  # Kylin build and deployment guide
└── LICENSE                         # MIT License
```

## Development

### Live Development Mode

```bash
wails dev
```

This starts:
- Vite dev server with hot-reload for frontend changes
- Dev server at `http://localhost:34115` for browser-based debugging with Go method access

### Building for Production

```bash
wails build               # Windows (default platform)
```

For Kylin Linux ARM64:

```bash
.\build-kylin.ps1         # Windows
bash build-kylin.sh       # Linux/macOS
bash build-kylin-local.sh # Native on Kylin
```

## Architecture

The application follows a layered architecture:

```
┌────────────────────────────────────────────────┐
│            Frontend (React 18 / TypeScript)     │
│  App.tsx → Tabs: Settings, Recording, System,   │
│            Logs                                 │
│  Tailwind CSS dark theme                        │
└──────────────────┬─────────────────────────────┘
                   │ Wails IPC (JSON bindings + events)
┌──────────────────▼─────────────────────────────┐
│               App (app.go)                      │
│  Config, Engine, Recorder, Logger management    │
│  Methods: GetConfig, Start/Stop, StartRecording │
│           StopRecording, LoadRecording, etc.    │
└──┬──────┬──────────┬──────────┬───────────────┘
   │      │          │          │
   ▼      ▼          ▼          ▼
┌──────┐ ┌────────┐ ┌──────┐ ┌──────┐
│Config│ │ Engine │ │Recorder│ │Logger│
│JSON  │ │ mouse  │ │gohook │ │ring+│
│I/O   │ │        │ │event  │ │file │
└──┬───┘ └───┬────┘ └───┬───┘ └──────┘
   │         │          │
   ▼         ▼          ▼
┌──────┐ ┌───────┐ ┌──────────┐
│Auto- │ │robotgo│ │Recording │
│start │ │+Replay│ │JSON I/O  │
└──────┘ └───────┘ └──────────┘
```

## Dependencies

### Go Backend

| Package | Purpose |
|---------|---------|
| [Wails v2](https://wails.io/) | Desktop app framework (Go + WebView) |
| [robotgo](https://github.com/go-vgo/robotgo) | Mouse/keyboard automation |
| [gohook](https://github.com/robotn/gohook) | Global keyboard event hooking |
| [systray](https://github.com/getlantern/systray) | System tray icon and menu |
| [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) | Windows registry access |

### Frontend

- React 18
- TypeScript
- Tailwind CSS 3
- Vite 3

## License

MIT License - see [LICENSE](LICENSE) file for details.
