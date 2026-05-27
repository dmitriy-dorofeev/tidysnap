package main

import (
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
		fmt.Println(config.ConfigPath())
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
		runCleanup()
		return
	}

	// TUI mode
	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка запуска TUI: %v\n", err)
		os.Exit(1)
	}
}

func runCleanup() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфига: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0750); err != nil {
		log.Fatalf("Ошибка создания папки логов: %v", err)
	}

	logFile, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatalf("Ошибка открытия лога: %v", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	s := scanner.New(cfg.Extensions, cfg.RetentionDays)
	files, err := s.Scan(cfg.TargetDir)
	if err != nil {
		logger.Fatalf("Ошибка сканирования: %v", err)
	}

	c := cleaner.New(cfg.DryRun, logger)
	stats, err := c.Clean(files)
	if err != nil {
		logger.Fatalf("Ошибка очистки: %v", err)
	}

	logger.Printf("Cleanup complete: removed %d files, freed %d bytes", stats.FilesRemoved, stats.BytesFreed)
}
