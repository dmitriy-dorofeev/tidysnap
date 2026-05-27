# Configuration

## Locations

| Resource | Path |
|----------|------|
| Config | `~/Library/Application Support/tidysnap/config.yaml` |
| Logs | `~/Library/Logs/tidysnap/cleanup.log` |
| Plist (daemon) | `~/Library/LaunchAgents/com.tidysnap.plist` |
| Daemon stdout | `~/Library/Logs/tidysnap/stdout.log` |
| Daemon stderr | `~/Library/Logs/tidysnap/stderr.log` |

## Config Format

The `config.yaml` file in YAML format:

```yaml
target_dir: /Users/username/Desktop
extensions:
  - .png
  - .jpg
  - .jpeg
  - .mov
  - .mp4
  - .gif
retention_days: 30
dry_run: true
warning_ack: false
log_path: /Users/username/Library/Logs/tidysnap/cleanup.log
check_interval_hours: 24
```

## Field Descriptions

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `target_dir` | `string` | `~` (home folder) | Target directory for scanning |
| `extensions` | `[]string` | `.png, .jpg, .jpeg, .mov, .mp4, .gif` | List of file extensions to delete |
| `retention_days` | `int` | `30` | Files older than this many days will be deleted |
| `dry_run` | `bool` | `true` | Test mode: shows files but does not delete |
| `warning_ack` | `bool` | `false` | Flag confirming warning about system folder |
| `log_path` | `string` | `~/Library/Logs/tidysnap/cleanup.log` | Path to log file |
| `check_interval_hours` | `int` | `24` | Interval for background cleanup runs (in hours) |

## Notes

- Config is created automatically during the first TUI setup.
- If config is missing, default values are used.
- Extensions are case-insensitive (`.PNG` and `.png` are equivalent).
