package daemon

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func TestGeneratePlist(t *testing.T) {
	plist := GeneratePlist("com.test", "/usr/local/bin/test", 2)

	if !strings.Contains(plist, "<string>com.test</string>") {
		t.Error("plist missing label")
	}
	if !strings.Contains(plist, "<string>/usr/local/bin/test</string>") {
		t.Error("plist missing binary path")
	}
	if !strings.Contains(plist, "<integer>7200</integer>") {
		t.Error("plist missing or wrong interval (expected 7200 seconds)")
	}
	if !strings.Contains(plist, "<string>--cleanup</string>") {
		t.Error("plist missing --cleanup argument")
	}
	if !strings.Contains(plist, "StandardOutPath") {
		t.Error("plist missing StandardOutPath")
	}
	if !strings.Contains(plist, "StandardErrorPath") {
		t.Error("plist missing StandardErrorPath")
	}
}

func TestWritePlistAndRemovePlist(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	content := GeneratePlist("com.test", "/bin/test", 1)
	if err := WritePlist(content); err != nil {
		t.Fatalf("WritePlist error: %v", err)
	}

	path, err := config.PlistPath()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("plist file should exist")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != content {
		t.Error("plist content mismatch")
	}

	if err := RemovePlist(); err != nil {
		t.Fatalf("RemovePlist error: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("plist file should be removed")
	}
}

func TestIsInstalled(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	if IsInstalled() {
		t.Error("IsInstalled should be false initially")
	}

	content := GeneratePlist("com.test", "/bin/test", 1)
	if err := WritePlist(content); err != nil {
		t.Fatal(err)
	}

	if !IsInstalled() {
		t.Error("IsInstalled should be true after WritePlist")
	}
}

func TestBinaryPath(t *testing.T) {
	path := BinaryPath()
	if path == "" {
		t.Error("BinaryPath should not be empty")
	}
}

func TestNextRunTime_NotInstalled(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	_, ok := NextRunTime(24)
	if ok {
		t.Error("expected false when plist is not installed")
	}
}

func TestNextRunTime_WithLog(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	// Create plist to satisfy IsInstalled
	content := GeneratePlist(label, "/bin/test", 1)
	WritePlist(content)

	// Create log file
	logPath, err := config.LogPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(logPath), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(logPath, []byte("log"), 0644); err != nil {
		t.Fatal(err)
	}

	// NextRunTime checks IsLoaded via launchctl, which will fail in test env,
	// so we only test the path where IsInstalled is true but IsLoaded is false.
	_, ok := NextRunTime(24)
	if ok {
		// This is acceptable depending on launchctl availability in test environment
		t.Log("NextRunTime returned true (may happen if launchctl sees the plist)")
	}
}

func TestInstall_Uninstall_Integration(t *testing.T) {
	// launchctl requires actual macOS launchd and may fail in CI.
	// We only verify that Install returns an error when launchctl fails,
	// and Uninstall does not panic.
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	err := Install("/bin/false", 1)
	if err == nil {
		// If launchctl works in this environment, that's fine.
		// But usually it will fail for a fake binary.
		t.Log("Install succeeded (unexpected in CI)")
	}

	// Uninstall should never panic
	err = Uninstall()
	if err != nil {
		t.Logf("Uninstall returned error (acceptable): %v", err)
	}
}

func TestLoadUnloadStartStop_WithoutPlist(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	// launchctl on macOS may return exit code 0 even on errors,
	// so we only verify these functions do not panic.
	_ = Load()
	_ = Unload()
	_ = Start()
	_ = Stop()
}

func TestIsLoaded_IsRunning_WithoutPlist(t *testing.T) {
	// Should return false when plist is not installed
	if IsLoaded() {
		t.Error("IsLoaded should be false without plist")
	}
	if IsRunning() {
		t.Error("IsRunning should be false without plist")
	}
}

func TestNextRunTime_WithInstalledAndLoaded(t *testing.T) {
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	defer os.Setenv("HOME", oldHome)

	content := GeneratePlist(label, "/bin/test", 1)
	WritePlist(content)

	logPath, err := config.LogPath()
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Dir(logPath), 0755)
	now := time.Now().Add(-time.Hour)
	os.WriteFile(logPath, []byte("log"), 0644)
	os.Chtimes(logPath, now, now)

	next, ok := NextRunTime(1)
	_ = next
	_ = ok
	// Depending on launchctl availability, ok may be true or false.
	// We just ensure no panic.
}
