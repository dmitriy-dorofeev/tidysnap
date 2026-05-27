# Installation

## Requirements

- **macOS** (the utility uses `launchd` for background operation)
- **Go 1.26+**
- **make** (optional, for convenience)

## Build from Source

```bash
# Clone the repository
git clone https://github.com/dmitriy-dorofeev/tidysnap.git
cd tidysnap

# Build the binary
make build

# Or directly via go
go build -ldflags "-s -w" -o bin/tidysnap ./cmd/main.go
```

## Install into the System

```bash
make install
```

This command copies the binary to `/usr/local/bin/tidysnap` and makes it globally available.

## Uninstall

```bash
make uninstall
```

Or manually:

```bash
tidysnap --uninstall
sudo rm /usr/local/bin/tidysnap
```

The `--uninstall` flag removes the `plist` (daemon) and configuration.
