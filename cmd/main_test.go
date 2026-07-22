package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/logger"
)

func setTestHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	t.Cleanup(func() { os.Setenv("HOME", oldHome) })
	return tmp
}

func TestRunCleanup(t *testing.T) {
	tmp := setTestHome(t)
	targetDir := filepath.Join(tmp, "target")
	os.MkdirAll(targetDir, 0755)

	// Create an old file
	oldFile := filepath.Join(targetDir, "old.png")
	os.WriteFile(oldFile, []byte("data"), 0644)
	oldTime := time.Now().Add(-100 * 24 * time.Hour)
	os.Chtimes(oldFile, oldTime, oldTime)

	// Create config
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	cfg.TargetDir = targetDir
	cfg.Extensions = []string{".png"}
	cfg.RetentionDays = 30
	cfg.DryRun = false
	cfg.LogPath = filepath.Join(tmp, "cleanup.log")
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}

	// Run cleanup
	if err := runCleanup(); err != nil {
		t.Fatalf("runCleanup error: %v", err)
	}

	// Verify file was deleted
	if _, err := os.Stat(oldFile); !os.IsNotExist(err) {
		t.Error("old file should be deleted")
	}

	// Verify log was written
	logData, err := os.ReadFile(cfg.LogPath)
	if err != nil {
		t.Fatalf("log file not found: %v", err)
	}
	if !strings.Contains(string(logData), "Cleanup complete") {
		t.Errorf("log should contain 'Cleanup complete', got: %s", string(logData))
	}
}

func TestRunCleanup_DryRun(t *testing.T) {
	tmp := setTestHome(t)
	targetDir := filepath.Join(tmp, "target")
	os.MkdirAll(targetDir, 0755)

	oldFile := filepath.Join(targetDir, "old.png")
	os.WriteFile(oldFile, []byte("data"), 0644)
	oldTime := time.Now().Add(-100 * 24 * time.Hour)
	os.Chtimes(oldFile, oldTime, oldTime)

	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	cfg.TargetDir = targetDir
	cfg.Extensions = []string{".png"}
	cfg.RetentionDays = 30
	cfg.DryRun = true
	cfg.LogPath = filepath.Join(tmp, "cleanup.log")
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}

	if err := runCleanup(); err != nil {
		t.Fatalf("runCleanup error: %v", err)
	}

	if _, err := os.Stat(oldFile); os.IsNotExist(err) {
		t.Error("file should still exist in dry run")
	}
}

func TestRunCleanup_PrunesOldLogEntries(t *testing.T) {
	tmp := setTestHome(t)
	targetDir := filepath.Join(tmp, "target")
	os.MkdirAll(targetDir, 0755)

	oldFile := filepath.Join(targetDir, "old.png")
	os.WriteFile(oldFile, []byte("data"), 0644)
	oldTime := time.Now().Add(-100 * 24 * time.Hour)
	os.Chtimes(oldFile, oldTime, oldTime)

	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	cfg.TargetDir = targetDir
	cfg.Extensions = []string{".png"}
	cfg.RetentionDays = 30
	cfg.DryRun = false
	cfg.LogPath = filepath.Join(tmp, "cleanup.log")
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}

	// Pre-populate log with old and recent entries.
	recentEntry := time.Now().Format(logger.LogTimeFormat) + " Recent entry\n"
	oldEntry := time.Now().Add(-31*24*time.Hour).Format(logger.LogTimeFormat) + " Old entry\n"
	if err := os.WriteFile(cfg.LogPath, []byte(oldEntry+recentEntry), 0600); err != nil {
		t.Fatal(err)
	}

	if err := runCleanup(); err != nil {
		t.Fatalf("runCleanup error: %v", err)
	}

	logData, err := os.ReadFile(cfg.LogPath)
	if err != nil {
		t.Fatalf("log file not found: %v", err)
	}
	content := string(logData)
	if strings.Contains(content, "Old entry") {
		t.Error("old log entry should be pruned")
	}
	if !strings.Contains(content, "Recent entry") {
		t.Error("recent log entry should be kept")
	}
	if !strings.Contains(content, "Cleanup complete") {
		t.Errorf("log should contain 'Cleanup complete', got: %s", content)
	}
}
