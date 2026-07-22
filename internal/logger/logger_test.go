package logger

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestPrune_KeepsRecentLines(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cleanup.log")

	recent := time.Now().Format(LogTimeFormat) + " Recent entry\n"
	old := time.Now().Add(-31*24*time.Hour).Format(LogTimeFormat) + " Old entry\n"
	noDate := "Some line without date\n"

	if err := os.WriteFile(path, []byte(old+recent+noDate), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Prune(path, 30); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)
	if !strings.Contains(content, "Recent entry") {
		t.Error("expected recent entry to be kept")
	}
	if strings.Contains(content, "Old entry") {
		t.Error("expected old entry to be removed")
	}
	if !strings.Contains(content, "Some line without date") {
		t.Error("expected line without date to be kept")
	}
}

func TestPrune_RemovesAllOldLines(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cleanup.log")

	lines := []string{
		time.Now().Add(-100*24*time.Hour).Format(LogTimeFormat) + " very old\n",
		time.Now().Add(-60*24*time.Hour).Format(LogTimeFormat) + " old\n",
		time.Now().Add(-31*24*time.Hour).Format(LogTimeFormat) + " slightly old\n",
	}

	if err := os.WriteFile(path, []byte(strings.Join(lines, "")), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Prune(path, 30); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) != 0 {
		t.Errorf("expected empty log, got: %q", string(data))
	}
}

func TestPrune_NoFile(t *testing.T) {
	if err := Prune(filepath.Join(t.TempDir(), "missing.log"), 30); err != nil {
		t.Fatal(err)
	}
}

func TestPrune_EmptyFile(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cleanup.log")
	if err := os.WriteFile(path, []byte{}, 0600); err != nil {
		t.Fatal(err)
	}

	if err := Prune(path, 30); err != nil {
		t.Fatal(err)
	}
}

func TestPrune_ZeroRetention(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cleanup.log")
	content := time.Now().Add(-100*24*time.Hour).Format(LogTimeFormat) + " Old entry\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}

	if err := Prune(path, 0); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "Old entry") {
		t.Error("expected no pruning when retention is 0")
	}
}

func TestPrune_PreservesFileMode(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "cleanup.log")
	content := time.Now().Add(-31*24*time.Hour).Format(LogTimeFormat) + " Old entry\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	if err := Prune(path, 30); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0644 {
		t.Errorf("expected mode 0644, got %v", info.Mode().Perm())
	}
}
