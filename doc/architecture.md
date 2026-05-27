# Архитектура

## Общая схема

```
cmd/main.go
    ├── config      (загрузка/сохранение YAML)
    ├── scanner     (рекурсивный обход и фильтрация файлов)
    ├── cleaner     (удаление или dry-run с логированием)
    ├── daemon      (launchd: plist, load/unload/start/stop)
    └── tui         (Bubble Tea: экраны, формы, навигация)
```

## Структура проекта

```
tidysnap/
├── cmd/
│   └── main.go                  # Точка входа, парсинг флагов, оркестрация
├── internal/
│   ├── config/
│   │   ├── config.go            # Load, Save, Reset, Exists
│   │   └── defaults.go          # DefaultConfig, пути к ресурсам macOS
│   ├── scanner/
│   │   └── scanner.go           # ScanResult, CleanupStats, Scanner.Scan()
│   ├── cleaner/
│   │   └── cleaner.go           # Cleaner.Clean() — удаление или dry-run
│   ├── daemon/
│   │   ├── install.go           # Install, Uninstall, IsRunning, NextRunTime
│   │   └── plist.go             # GeneratePlist, WritePlist, RemovePlist
│   └── tui/
│       ├── model.go             # Главная модель Bubble Tea, экраны, сообщения
│       ├── welcome.go           # Экран приветствия
│       ├── folderpicker.go      # Выбор папки через file picker
│       ├── setup.go             # Форма настроек (huh)
│       ├── warning.go           # Предупреждение о системных папках
│       ├── status.go            # Экран статуса и управления
│       ├── preview.go           # Предпросмотр файлов к удалению
│       ├── logview.go           # Просмотр логов
│       └── reset.go             # Подтверждение сброса
├── bin/                         # Сборочный артефакт
├── Makefile
├── go.mod
└── README.md
```

## Модули

### `config`

Отвечает за сериализацию/десериализацию конфигурации в YAML. Использует стандартные пути macOS (`~/Library/Application Support/`, `~/Library/Logs/`, `~/Library/LaunchAgents/`).

### `scanner`

Рекурсивно обходит `target_dir` через `filepath.Walk`. Фильтрует файлы по:
- расширению (регистронезависимо)
- времени модификации (файлы старше `retention_days`)

Возвращает слайс `ScanResult` с путём, размером, временем и расширением.

### `cleaner`

Принимает слайс `ScanResult` и выполняет удаление (или имитацию в dry-run). Ведёт статистику (`CleanupStats`) и логирует каждую операцию через `log.Logger`.

### `daemon`

Абстракция над `launchd`:
- `Install` — генерирует plist, записывает в `LaunchAgents`, загружает и стартует
- `Uninstall` — останавливает, выгружает, удаляет plist
- `IsInstalled` / `IsLoaded` / `IsRunning` — проверка состояния
- `NextRunTime` — эвристика: время последней модификации лога + `check_interval_hours`

### `tui`

Построен на [Bubble Tea](https://github.com/charmbracelet/bubbletea) и [huh](https://github.com/charmbracelet/huh).

Экраны (конечный автомат в `model.go`):

| Экран | Описание |
|-------|----------|
| `screenWelcome` | Приветствие и начало настройки |
| `screenFolderPicker` | Интерактивный выбор папки |
| `screenSetup` | Форма параметров (расширения, срок, интервал, dry-run) |
| `screenWarning` | Предупреждение при выборе Desktop/Downloads/Documents |
| `screenStatus` | Главный экран: статус, запуск очистки, управление демоном |
| `screenPreview` | Список файлов к удалению |
| `screenLogView` | Просмотр файла логов |
| `screenResetConfirm` | Подтверждение удаления настроек |

## Поток данных

### Интерактивный запуск

```
main.go
  → tea.NewProgram(tui.InitialModel())
    → model.Init(): config.Load()
      → если конфига нет → screenWelcome
      → если конфиг есть → screenStatus
        → [r] → scanner.Scan() → screenPreview
          → подтверждение → cleaner.Clean() → screenStatus
        → [s] → daemon.Install() / Start() / Stop() / Load()
```

### Фоновый запуск (launchd)

```
launchd → tidysnap --cleanup
  → config.Load()
  → scanner.Scan()
  → cleaner.Clean()
  → запись в log_file
```
