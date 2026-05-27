# TidySnap

A macOS utility with a TUI interface for automatic cleanup of screenshots and screen recordings.

## Quick Start

```bash
git clone https://github.com/dmitriy-dorofeev/tidysnap.git
cd tidysnap
make install
```

Run `tidysnap` and follow the interactive setup.

## Features

- 🔍 **Extension-based file search** — works with any system language
- ⚙️ Configure folder, extensions, and retention period via TUI
- 🧪 Dry Run mode enabled by default
- 🔄 Background operation via `launchd`
- 🌐 Automatic English / Russian interface based on system locale

## Commands

```bash
tidysnap              # TUI mode
tidysnap --cleanup    # Background cleanup (for launchd)
tidysnap --config-path
tidysnap --reset
tidysnap --uninstall
tidysnap --version
```

## Documentation

Detailed documentation is available in the [`doc/`](doc/) folder:

- [Installation](doc/en/installation.md)
- [Usage](doc/en/usage.md)
- [Configuration](doc/en/configuration.md)
- [Architecture](doc/en/architecture.md)
- [Security](doc/en/security.md)
- [Development](doc/en/development.md)

Documentation is also available in [Russian](doc/ru/).

## License

MIT
