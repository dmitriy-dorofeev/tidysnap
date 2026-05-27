# Development

## Dependencies

Main project dependencies:

| Package | Purpose |
|---------|---------|
| `github.com/charmbracelet/bubbletea` | TUI framework (Model-Update-View) |
| `github.com/charmbracelet/bubbles` | Ready-made components (spinner, file picker) |
| `github.com/charmbracelet/huh` | Forms and input in TUI |
| `github.com/charmbracelet/lipgloss` | Text styling |
| `github.com/dustin/go-humanize` | File size formatting |
| `gopkg.in/yaml.v3` | Config serialization |

## Makefile Commands

```bash
make build         # Build binary to bin/tidysnap
make run           # Run via go run
make install       # Build and copy to /usr/local/bin
make uninstall     # Remove binary, plist, and config
make clean         # Remove bin/ directory
make test          # Run all tests
make check         # Run formatting, vet, staticcheck, govulncheck, and gosec
make install-tools # Install linting and security tools
```

## Testing

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/scanner
go test ./internal/cleaner
go test ./internal/config
go test ./internal/daemon
go test ./internal/tui
go test ./cmd
```

## Test Structure

Tests are located next to the files they test (`*_test.go`):

- `cmd/main_test.go` — flag parsing tests
- `cmd/main_integration_test.go` — integration tests
- `internal/*/..._test.go` — unit tests

## Build with Version

To embed version, commit, and build date, use `ldflags`:

```bash
go build -ldflags "-X main.version=1.0.0 -X main.commit=abc123 -X main.date=2026-05-27" -o bin/tidysnap ./cmd/main.go
```

## Linting

Recommended to use `golangci-lint`:

```bash
golangci-lint run ./...
```

## Contributing

1. Fork the repository.
2. Create a branch: `git checkout -b feature/my-feature`.
3. Make changes and add tests.
4. Ensure `make test` passes.
5. Submit a Pull Request.
