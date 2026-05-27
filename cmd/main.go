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
	"github.com/dmitriy-dorofeev/tidysnap/internal/i18n"
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
		cleanupFlag    = flag.Bool("cleanup", false, i18n.T("flag_cleanup"))
		configPathFlag = flag.Bool("config-path", false, i18n.T("flag_config_path"))
		resetFlag      = flag.Bool("reset", false, i18n.T("flag_reset"))
		uninstallFlag  = flag.Bool("uninstall", false, i18n.T("flag_uninstall"))
		versionFlag    = flag.Bool("version", false, i18n.T("flag_version"))
	)
	flag.Parse()

	if *versionFlag {
		fmt.Printf("tidysnap %s (commit: %s, built: %s)\n", version, commit, date)
		return
	}

	if *configPathFlag {
		path, err := config.ConfigPath()
		if err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("err_generic"), err)
			os.Exit(1)
		}
		fmt.Println(path)
		return
	}

	if *resetFlag {
		if err := config.Reset(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("err_reset"), err)
			os.Exit(1)
		}
		fmt.Println(i18n.T("settings_reset"))
		return
	}

	if *uninstallFlag {
		if err := daemon.Uninstall(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("err_uninstall"), err)
			os.Exit(1)
		}
		if err := config.Reset(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("err_uninstall_config"), err)
		}
		fmt.Println(i18n.T("tidysnap_removed"))
		return
	}

	if *cleanupFlag {
		if err := runCleanup(); err != nil {
			fmt.Fprintf(os.Stderr, i18n.T("err_cleanup"), err)
			os.Exit(1)
		}
		return
	}

	// TUI mode
	p := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, i18n.T("err_tui"), err)
		os.Exit(1)
	}
}

func runCleanup() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf(i18n.T("err_load_config"), err)
	}

	if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0750); err != nil {
		return fmt.Errorf(i18n.T("err_log_dir"), err)
	}

	logFile, err := os.OpenFile(cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf(i18n.T("err_open_log"), err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)

	ctx := context.Background()

	s := scanner.New(cfg.Extensions, cfg.RetentionDays)
	files, err := s.Scan(ctx, cfg.TargetDir)
	if err != nil {
		return fmt.Errorf(i18n.T("err_scan"), err)
	}

	c := cleaner.New(cfg.DryRun, logger)
	stats, err := c.Clean(ctx, files)
	if err != nil {
		return fmt.Errorf(i18n.T("err_clean"), err)
	}

	logger.Printf("Cleanup complete: removed %d files, freed %d bytes", stats.FilesRemoved, stats.BytesFreed)
	return nil
}
