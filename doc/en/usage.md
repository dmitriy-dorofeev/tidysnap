# Usage

## Operating Modes

TidySnap supports two modes: **interactive TUI** (default) and **CLI** with flags.

## TUI Mode (Default)

Running without arguments opens the interactive interface:

```bash
tidysnap
```

### First Run (Welcome → Setup)

1. **Welcome screen** — press `Enter` or `s` to start setup.
2. **Folder selection** — navigate with `↑`/`↓`, go into folders with `→`/`Enter`, go back with `←`, and press `Space` to **select** the current folder.
3. **Settings form**:
   - **File extensions** — comma-separated, e.g., `.png, .mov, .mp4`
   - **Retention period (days)** — files older than this will be deleted
   - **Check interval (hours)** — how often to run background cleanup
   - **Dry Run mode** — show files but do not delete
4. **Warning** — if one of the system folders (`Desktop`, `Downloads`, `Documents`, `Movies`, `Music`, `Pictures`, `Public`, `Library`) is selected, a warning appears.
5. **Save** — settings are saved to `config.yaml`, and you proceed to the status screen.

### Status Screen

Displays current settings and daemon state:

| Key | Action |
|-----|--------|
| `r` | Run cleanup (scan + preview) |
| `l` | Open logs |
| `e` | Edit settings (folder and parameters) |
| `s` | Daemon action (depends on state): Install → Load → Start → Stop |
| `x` | Delete settings and plist (reset) |
| `q` / `Esc` | Quit |

### Preview Screen

After scanning (`r`), a list of files to be deleted is shown. If **Dry Run** is enabled, files are not actually deleted — only logged.

| Key | Action |
|-----|--------|
| `d` / `Enter` | Execute cleanup (delete or dry-run) |
| `q` / `Esc` | Back to status |

## CLI Flags

| Flag | Description |
|------|-------------|
| `--cleanup` | Run background cleanup (no TUI, for `launchd`) |
| `--config-path` | Show the full path to the configuration file |
| `--reset` | Reset settings (delete `config.yaml`) |
| `--uninstall` | Remove plist, config, and unload daemon |
| `--version` | Show version, commit, and build date |

### Examples

```bash
# Show config path
tidysnap --config-path

# Reset settings
tidysnap --reset

# Full uninstall
tidysnap --uninstall
```

## Daemon Management (launchd)

TidySnap uses `launchd` for background operation. Management is done via TUI (key `s`), but can also be done manually:

```bash
# Load daemon
launchctl load ~/Library/LaunchAgents/com.tidysnap.plist

# Unload daemon
launchctl unload ~/Library/LaunchAgents/com.tidysnap.plist

# Start immediately
launchctl start com.tidysnap

# Stop
launchctl stop com.tidysnap

# Check status
launchctl list | grep com.tidysnap
```

The daemon runs `tidysnap --cleanup` at the configured interval (`StartInterval` in seconds).
