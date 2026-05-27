# Конфигурация

## Расположение

| Ресурс | Путь |
|--------|------|
| Конфиг | `~/Library/Application Support/tidysnap/config.yaml` |
| Логи | `~/Library/Logs/tidysnap/cleanup.log` |
| Plist (демон) | `~/Library/LaunchAgents/com.tidysnap.plist` |
| Stdout демона | `~/Library/Logs/tidysnap/stdout.log` |
| Stderr демона | `~/Library/Logs/tidysnap/stderr.log` |

## Формат конфига

Файл `config.yaml` в формате YAML:

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

## Описание полей

| Поле | Тип | По умолчанию | Описание |
|------|-----|--------------|----------|
| `target_dir` | `string` | `~` (домашняя папка) | Целевая директория для сканирования |
| `extensions` | `[]string` | `.png, .jpg, .jpeg, .mov, .mp4, .gif` | Список расширений файлов для удаления |
| `retention_days` | `int` | `30` | Файлы старше этого количества дней будут удалены |
| `dry_run` | `bool` | `true` | Тестовый режим: показывает файлы, но не удаляет |
| `warning_ack` | `bool` | `false` | Флаг подтверждения предупреждения о системной папке |
| `log_path` | `string` | `~/Library/Logs/tidysnap/cleanup.log` | Путь к файлу логов |
| `check_interval_hours` | `int` | `24` | Интервал запуска фоновой очистки (в часах) |

## Примечания

- Конфиг создаётся автоматически при первой настройке через TUI.
- Если конфиг отсутствует, используются значения по умолчанию.
- Поле `warning_ack` сбрасывается при смене папки, чтобы предупреждение показывалось снова.
- Расширения регистронезависимы (`.PNG` и `.png` эквивалентны).
