# MousePaw - Mouse Automation Tool

A Windows desktop application for automating mouse movements, clicks, and scrolling. Built with [Wails v2](https://wails.io/) (Go + React/TypeScript).

[中文文档](README_zh.md)

## Features

- **Mouse Movement**: Automatically move cursor to random positions at configurable intervals (1-60 seconds)
- **Mouse Click**: Perform left/right/middle clicks with adjustable intervals (0.5-30 seconds) and repeat counts (1-10)
- **Mouse Scroll**: Scroll in any direction (up/down/left/right) with configurable intervals (1-30 seconds) and amounts (1-20)
- **Global Hotkeys**: F6 to start, F7 to stop (works even when window is not focused)
- **Auto-start**: Optional Windows startup integration
- **Real-time Logging**: View automation activity in the built-in log viewer
- **Configuration Persistence**: Settings are saved automatically and restored on next launch

## Screenshots

*Add screenshots of the application UI here*

## Installation

### Prerequisites

- Windows 10/11
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

## Usage

1. **Launch the application** - Run `mousepaw.exe` or use `wails dev` for development
2. **Configure automation** - Use the tabs to enable/disable different mouse actions:
   - **Mouse Move**: Set interval and enable random cursor movement
   - **Mouse Click**: Configure click type, interval, and repeat count
   - **Scroll**: Set scroll direction, interval, and amount
3. **Start automation** - Click "Start" button or press **F6**
4. **Stop automation** - Click "Stop" button or press **F7**
5. **View logs** - Switch to "Execution Logs" tab to monitor activity

## Configuration

Configuration is stored in `mousepaw_config.json` next to the executable:

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

## Hotkeys

| Key | Action |
|-----|--------|
| **F6** | Start automation (global) |
| **F7** | Stop automation (global) |

## Project Structure

```
mousePaw/
├── main.go                 # Application entry point
├── app.go                  # Core App struct with Wails bindings
├── icon.go                 # Programmatic icon generation
├── pkg/
│   ├── autostart/          # Windows registry auto-start
│   ├── config/             # Configuration management
│   ├── engine/             # Mouse automation engine
│   └── log/                # In-memory + file logger
├── frontend/               # React/TypeScript UI
│   ├── src/
│   │   ├── App.tsx         # Main UI component
│   │   └── ...
│   └── wailsjs/            # Auto-generated Wails bindings
└── build/                  # Build assets and output
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
wails build
```

## Dependencies

### Go Backend
- [Wails v2](https://wails.io/) - Desktop app framework
- [robotgo](https://github.com/go-vgo/robotgo) - Mouse/keyboard automation
- [gohook](https://github.com/robotn/gohook) - Global keyboard/mouse event hooking
- [golang.org/x/sys](https://pkg.go.dev/golang.org/x/sys) - Windows registry access

### Frontend
- React 18
- TypeScript
- Tailwind CSS
- Vite

## License

MIT License - see [LICENSE](LICENSE) file for details.