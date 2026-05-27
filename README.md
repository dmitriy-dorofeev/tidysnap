# TidySnap

Утилита для macOS с TUI-интерфейсом для автоматической очистки скриншотов и записей экрана.

## Быстрый старт

```bash
git clone https://github.com/dmitriy-dorofeev/tidysnap.git
cd tidysnap
make install
```

Запустите `tidysnap` и следуйте интерактивной настройке.

## Основные возможности

- 🔍 Поиск файлов **по расширениям** — работает с любыми языками системы
- ⚙️ Настройка папки, расширений и срока хранения через TUI
- 🧪 Тестовый режим (Dry Run) по умолчанию
- 🔄 Фоновая работа через `launchd`

## Команды

```bash
tidysnap              # TUI-режим
tidysnap --cleanup    # Фоновая очистка (для launchd)
tidysnap --config-path
tidysnap --reset
tidysnap --uninstall
tidysnap --version
```

## Документация

Подробная документация находится в папке [`doc/`](doc/):

- [Установка](doc/installation.md)
- [Использование](doc/usage.md)
- [Конфигурация](doc/configuration.md)
- [Архитектура](doc/architecture.md)
- [Безопасность](doc/security.md)
- [Разработка](doc/development.md)

## Лицензия

MIT
