package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/cleaner"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
	"github.com/dmitriy-dorofeev/tidysnap/internal/tui"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var (
		cleanupFlag    = flag.Bool("cleanup", false, "Запустить фоновую очистку (без TUI)")
		configPathFlag = flag.Bool("config-path", false, "Показать путь к конфигу")
		resetFlag      = flag.Bool("reset", false, "Сбросить настройки")
		uninstallFlag  = flag.Bool("uninstall", false, "Удалить plist и конфиг")
		versionFlag    = flag.Bool("version", false, "Показать версию")
	)
	flag.Parse()

	if *versionFlag {
		fmt.Printf("tidysnap %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	if *configPathFlag {
		path, err := config.ConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(path)
		return
	}

	if *resetFlag {
		if err := config.Reset(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка сброса: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Настройки сброшены.")
		return
	}

	if *uninstallFlag {
		if err := daemon.Uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка удаления: %v\n", err)
			os.Exit(1)
		}
		if err := config.Reset(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка удаления конфига: %v\n", err)
		}
		fmt.Println("TidySnap удалён.")
		return
	}

	if *cleanupFlag {
		if err := runCleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "Ошибка очистки: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// TUI mode
	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка запуска TUI: %v\n", err)
		os.Exit(1)
	}
}

func runCleanup() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфига: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0750); err != nil {
		return fmt.Errorf("ошибка создания папки логов: %w", err)
	}

	logFile, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("ошибка открытия лога: %w", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	ctx := context.Background()

	s := scanner.New(cfg.Extensions, cfg.RetentionDays)
	files, err := s.Scan(ctx, cfg.TargetDir)
	if err != nil {
		return fmt.Errorf("ошибка сканирования: %w", err)
	}

	c := cleaner.New(cfg.DryRun, logger)
	stats, err := c.Clean(ctx, files)
	if err != nil {
		return fmt.Errorf("ошибка очистки: %w", err)
	}

	logger.Printf("Cleanup complete: removed %d files, freed %d bytes", stats.FilesRemoved, stats.BytesFreed)
	return nil
}
