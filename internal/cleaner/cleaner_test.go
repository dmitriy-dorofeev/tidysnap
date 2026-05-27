package cleaner

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := New(true, logger)
	if c.dryRun != true {
		t.Errorf("dryRun = %v, want true", c.dryRun)
	}
	if c.logger != logger {
		t.Error("logger mismatch")
	}
}

func TestClean_DryRun(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "file.png")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := New(true, logger)

	files := []scanner.ScanResult{
		{Path: f, Size: 5},
	}

	stats, err := c.Clean(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.FilesRemoved != 1 {
		t.Errorf("FilesRemoved = %d, want 1", stats.FilesRemoved)
	}
	if stats.BytesFreed != 5 {
		t.Errorf("BytesFreed = %d, want 5", stats.BytesFreed)
	}

	// file should still exist
	if _, err := os.Stat(f); os.IsNotExist(err) {
		t.Error("file should not be deleted in dry run")
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "DRY RUN") {
		t.Errorf("log should contain DRY RUN, got: %s", logOutput)
	}
}

func TestClean_RealDelete(t *testing.T) {
	tmp := t.TempDir()
	f := filepath.Join(tmp, "file.png")
	if err := os.WriteFile(f, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := New(false, logger)

	files := []scanner.ScanResult{
		{Path: f, Size: 5},
	}

	stats, err := c.Clean(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.FilesRemoved != 1 {
		t.Errorf("FilesRemoved = %d, want 1", stats.FilesRemoved)
	}
	if stats.BytesFreed != 5 {
		t.Errorf("BytesFreed = %d, want 5", stats.BytesFreed)
	}

	if _, err := os.Stat(f); !os.IsNotExist(err) {
		t.Error("file should be deleted")
	}

	logOutput := buf.String()
	if !strings.Contains(logOutput, "Deleted") {
		t.Errorf("log should contain Deleted, got: %s", logOutput)
	}
}

func TestClean_MultipleFiles(t *testing.T) {
	tmp := t.TempDir()
	f1 := filepath.Join(tmp, "a.png")
	f2 := filepath.Join(tmp, "b.png")
	os.WriteFile(f1, []byte("a"), 0644)
	os.WriteFile(f2, []byte("bb"), 0644)

	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := New(false, logger)

	files := []scanner.ScanResult{
		{Path: f1, Size: 1},
		{Path: f2, Size: 2},
	}

	stats, err := c.Clean(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.FilesRemoved != 2 {
		t.Errorf("FilesRemoved = %d, want 2", stats.FilesRemoved)
	}
	if stats.BytesFreed != 3 {
		t.Errorf("BytesFreed = %d, want 3", stats.BytesFreed)
	}
}

func TestClean_ErrorOnMissingFile(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := New(false, logger)

	files := []scanner.ScanResult{
		{Path: "/nonexistent/file.png", Size: 0},
	}

	stats, err := c.Clean(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats.FilesRemoved != 0 {
		t.Errorf("FilesRemoved = %d, want 0", stats.FilesRemoved)
	}
	if len(stats.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(stats.Errors))
	}
}

func TestClean_EmptyList(t *testing.T) {
	var buf bytes.Buffer
	logger := log.New(&buf, "", 0)
	c := New(false, logger)

	stats, err := c.Clean(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.FilesRemoved != 0 {
		t.Errorf("FilesRemoved = %d, want 0", stats.FilesRemoved)
	}
	if stats.BytesFreed != 0 {
		t.Errorf("BytesFreed = %d, want 0", stats.BytesFreed)
	}
	if !stats.Timestamp.IsZero() {
		// timestamp should be set
		if time.Since(stats.Timestamp) > time.Minute {
			t.Error("timestamp seems incorrect")
		}
	}
}
