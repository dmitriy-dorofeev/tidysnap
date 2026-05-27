# Установка

## Требования

- **macOS** (утилита использует `launchd` для фоновой работы)
- **Go 1.26+**
- **make** (опционально, для удобства)

## Сборка из исходников

```bash
# Клонировать репозиторий
git clone https://github.com/dmitriy-dorofeev/tidysnap.git
cd tidysnap

# Собрать бинарник
make build

# Или напрямую через go
go build -ldflags "-s -w" -o bin/tidysnap ./cmd/main.go
```

## Установка в систему

```bash
make install
```

Команда копирует бинарник в `/usr/local/bin/tidysnap` и делает его доступным глобально.

## Удаление

```bash
make uninstall
```

Или вручную:

```bash
tidysnap --uninstall
sudo rm /usr/local/bin/tidysnap
```

Флаг `--uninstall` удаляет `plist` (демон) и конфигурацию.
