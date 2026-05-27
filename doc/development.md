# Разработка

## Зависимости

Основные зависимости проекта:

| Пакет | Назначение |
|-------|------------|
| `github.com/charmbracelet/bubbletea` | TUI-фреймворк (Model-Update-View) |
| `github.com/charmbracelet/bubbles` | Готовые компоненты (spinner, file picker) |
| `github.com/charmbracelet/huh` | Формы и ввод данных в TUI |
| `github.com/charmbracelet/lipgloss` | Стилизация текста |
| `github.com/dustin/go-humanize` | Форматирование размеров файлов |
| `gopkg.in/yaml.v3` | Сериализация конфига |

## Команды Makefile

```bash
make build      # Собрать бинарник в bin/tidysnap
make run        # Запустить через go run
make install    # Собрать и скопировать в /usr/local/bin
make uninstall  # Удалить бинарник, plist и конфиг
make clean      # Удалить папку bin/
make test       # Запустить все тесты
```

## Тестирование

```bash
# Все тесты
go test ./...

# С покрытием
go test -cover ./...

# Конкретный пакет
go test ./internal/scanner
go test ./internal/cleaner
go test ./internal/config
go test ./internal/daemon
go test ./internal/tui
go test ./cmd
```

## Структура тестов

Тесты расположены рядом с тестируемыми файлами (`*_test.go`):

- `cmd/main_test.go` — тесты парсинга флагов
- `cmd/main_integration_test.go` — интеграционные тесты
- `internal/*/..._test.go` — модульные тесты

## Сборка с версией

Для встраивания версии, коммита и даты сборки используйте `ldflags`:

```bash
go build -ldflags "-X main.version=1.0.0 -X main.commit=abc123 -X main.date=2026-05-27" -o bin/tidysnap ./cmd/main.go
```

## Линтинг

Рекомендуется использовать `golangci-lint`:

```bash
golangci-lint run ./...
```

## Вклад в проект

1. Форкните репозиторий.
2. Создайте ветку: `git checkout -b feature/my-feature`.
3. Внесите изменения и добавьте тесты.
4. Убедитесь, что `make test` проходит.
5. Отправьте Pull Request.
