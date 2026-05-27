package config

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TargetDir == "" {
		t.Error("TargetDir should not be empty")
	}
	if len(cfg.Extensions) == 0 {
		t.Error("Extensions should not be empty")
	}
	if cfg.RetentionDays != 30 {
		t.Errorf("RetentionDays = %d, want 30", cfg.RetentionDays)
	}
	if cfg.DryRun != true {
		t.Errorf("DryRun = %v, want true", cfg.DryRun)
	}
	if cfg.CheckIntervalHours != 24 {
		t.Errorf("CheckIntervalHours = %d, want 24", cfg.CheckIntervalHours)
	}
}

func TestConfigPath(t *testing.T) {
	path := ConfigPath()
	if path == "" {
		t.Error("ConfigPath should not be empty")
	}
	if filepath.Base(path) != "config.yaml" {
		t.Errorf("basename = %q, want config.yaml", filepath.Base(path))
	}
}

func TestConfigDir(t *testing.T) {
	dir := ConfigDir()
	if dir == "" {
		t.Error("ConfigDir should not be empty")
	}
}

func TestLogPath(t *testing.T) {
	path := LogPath()
	if path == "" {
		t.Error("LogPath should not be empty")
	}
	if filepath.Base(path) != "cleanup.log" {
		t.Errorf("basename = %q, want cleanup.log", filepath.Base(path))
	}
}

func TestPlistPath(t *testing.T) {
	path := PlistPath()
	if path == "" {
		t.Error("PlistPath should not be empty")
	}
	if filepath.Ext(path) != ".plist" {
		t.Errorf("ext = %q, want .plist", filepath.Ext(path))
	}
}

func TestLoad_DefaultWhenMissing(t *testing.T) {
	// Temporarily override config path by using a temp dir
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg == nil {
		t.Fatal("cfg is nil")
	}
	if cfg.RetentionDays != 30 {
		t.Errorf("RetentionDays = %d, want 30", cfg.RetentionDays)
	}
}

func TestLoad_ReadError(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	// Create a directory instead of a file to cause a read error
	if err := os.RemoveAll(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(path, 0755); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("expected error when config path is a directory")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	cfg := &Config{
		TargetDir:          "/tmp/test",
		Extensions:         []string{".png", ".jpg"},
		RetentionDays:      7,
		DryRun:             false,
		WarningAck:         true,
		LogPath:            "/tmp/test.log",
		CheckIntervalHours: 12,
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	if !Exists() {
		t.Error("Exists should return true after Save")
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}
	if loaded.TargetDir != cfg.TargetDir {
		t.Errorf("TargetDir = %q, want %q", loaded.TargetDir, cfg.TargetDir)
	}
	if loaded.RetentionDays != cfg.RetentionDays {
		t.Errorf("RetentionDays = %d, want %d", loaded.RetentionDays, cfg.RetentionDays)
	}
	if loaded.DryRun != cfg.DryRun {
		t.Errorf("DryRun = %v, want %v", loaded.DryRun, cfg.DryRun)
	}
}

func TestSave_MkdirError(t *testing.T) {
	// Make parent dir a file so MkdirAll fails
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	path := ConfigPath()
	parent := filepath.Dir(path)
	// Create intermediate dirs, then create parent as a file
	if err := os.MkdirAll(filepath.Dir(parent), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(parent, []byte("block"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := DefaultConfig()
	err := Save(cfg)
	if err == nil {
		t.Fatal("expected error when parent path is a file")
	}
}

func TestReset(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	cfg := DefaultConfig()
	if err := Save(cfg); err != nil {
		t.Fatalf("Save error: %v", err)
	}
	if !Exists() {
		t.Fatal("Exists should be true")
	}

	if err := Reset(); err != nil {
		t.Fatalf("Reset error: %v", err)
	}

	if Exists() {
		t.Error("Exists should be false after Reset")
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not: yaml: ["), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestCleanup(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	cfg := DefaultConfig()
	Save(cfg)

	if err := Cleanup(); err != nil {
		t.Fatalf("Cleanup error: %v", err)
	}

	if Exists() {
		t.Error("Exists should be false after Cleanup")
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	cfg := Config{
		TargetDir:          "/tmp/target",
		Extensions:         []string{".png", ".mov"},
		RetentionDays:      14,
		DryRun:             true,
		WarningAck:         false,
		LogPath:            "/tmp/log",
		CheckIntervalHours: 6,
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var loaded Config
	if err := yaml.Unmarshal(data, &loaded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if loaded.TargetDir != cfg.TargetDir {
		t.Errorf("TargetDir = %q, want %q", loaded.TargetDir, cfg.TargetDir)
	}
	if loaded.RetentionDays != cfg.RetentionDays {
		t.Errorf("RetentionDays = %d, want %d", loaded.RetentionDays, cfg.RetentionDays)
	}
	if loaded.CheckIntervalHours != cfg.CheckIntervalHours {
		t.Errorf("CheckIntervalHours = %d, want %d", loaded.CheckIntervalHours, cfg.CheckIntervalHours)
	}
}
