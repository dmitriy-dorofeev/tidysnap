# Architecture

## Overview

```
cmd/main.go
    ├── config      (load/save YAML)
    ├── scanner     (recursive walk and file filtering)
    ├── cleaner     (deletion or dry-run with logging)
    ├── daemon      (launchd: plist, load/unload/start/stop)
    └── tui         (Bubble Tea: screens, forms, navigation)
```

## Project Structure

```
tidysnap/
├── cmd/
│   └── main.go                  # Entry point, flag parsing, orchestration
├── internal/
│   ├── config/
│   │   ├── config.go            # Load, Save, Reset, Exists
│   │   └── defaults.go          # DefaultConfig, macOS resource paths
│   ├── scanner/
│   │   └── scanner.go           # ScanResult, CleanupStats, Scanner.Scan()
│   ├── cleaner/
│   │   └── cleaner.go           # Cleaner.Clean() — deletion or dry-run
│   ├── daemon/
│   │   ├── install.go           # Install, Uninstall, IsRunning, NextRunTime
│   │   └── plist.go             # GeneratePlist, WritePlist, RemovePlist
│   └── tui/
│       ├── model.go             # Main Bubble Tea model, screens, messages
│       ├── keys.go              # Keyboard layout normalization (QWERTY / ЙЦУКЕН)
│       ├── welcome.go           # Welcome screen
│       ├── folderpicker.go      # Folder selection via file picker
│       ├── setup.go             # Settings form (huh)
│       ├── warning.go           # Warning about system folders
│       ├── status.go            # Status and management screen
│       ├── preview.go           # Preview of files to delete
│       ├── logview.go           # Log viewer
│       └── reset.go             # Reset confirmation
├── bin/                         # Build artifact
├── Makefile
├── go.mod
└── README.md
```

## Modules

### `config`

Handles serialization/deserialization of configuration to YAML. Uses standard macOS paths (`~/Library/Application Support/`, `~/Library/Logs/`, `~/Library/LaunchAgents/`).

### `scanner`

Recursively walks `target_dir` via `filepath.Walk`. Filters files by:
- extension (case-insensitive)
- modification time (files older than `retention_days`)

Returns a slice of `ScanResult` with path, size, time, and extension.

### `cleaner`

Accepts a slice of `ScanResult` and performs deletion (or dry-run simulation). Tracks statistics (`CleanupStats`) and logs every operation via `log.Logger`.

### `daemon`

Abstraction over `launchd`:
- `Install` — generates plist, writes to `LaunchAgents`, loads and starts
- `Uninstall` — stops, unloads, removes plist
- `IsInstalled` / `IsLoaded` / `IsRunning` — state checks
- `NextRunTime` — heuristic: last log modification time + `check_interval_hours`

### `tui`

Built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [huh](https://github.com/charmbracelet/huh).

Screens (finite state machine in `model.go`):

| Screen | Description |
|--------|-------------|
| `screenWelcome` | Welcome and start setup |
| `screenFolderPicker` | Interactive folder selection |
| `screenSetup` | Parameter form (extensions, retention, interval, dry-run) |
| `screenWarning` | Warning when selecting system folders (Desktop, Downloads, Documents, Movies, Music, Pictures, Public, Library) |
| `screenStatus` | Main screen: status, run cleanup, manage daemon |
| `screenPreview` | List of files to be deleted |
| `screenLogView` | Log file viewer |
| `screenResetConfirm` | Confirmation for deleting settings |

## Data Flow

### Interactive Run

```
main.go
  → tea.NewProgram(tui.InitialModel())
    → model.Init(): config.Load()
      → if no config → screenWelcome
      → if config exists → screenStatus
        → [r] → scanner.Scan() → screenPreview
          → confirmation → cleaner.Clean() → screenStatus
        → [s] → daemon.Install() / Start() / Stop() / Load()
```

### Background Run (launchd)

```
launchd → tidysnap --cleanup
  → config.Load()
  → scanner.Scan()
  → cleaner.Clean()
  → write to log_file
```
