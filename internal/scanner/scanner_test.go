package scanner

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	exts := []string{".PNG", ".MOV", ".MP4"}
	s := New(exts, 30)

	if len(s.extensions) != 3 {
		t.Fatalf("expected 3 extensions, got %d", len(s.extensions))
	}
	for i, want := range []string{".png", ".mov", ".mp4"} {
		if s.extensions[i] != want {
			t.Errorf("extension[%d] = %q, want %q", i, s.extensions[i], want)
		}
	}

	expectedRetention := time.Duration(30) * 24 * time.Hour
	if s.retention != expectedRetention {
		t.Errorf("retention = %v, want %v", s.retention, expectedRetention)
	}
}

func TestScan_EmptyDir(t *testing.T) {
	tmp := t.TempDir()
	s := New([]string{".png"}, 0)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestScan_MatchingFiles(t *testing.T) {
	tmp := t.TempDir()

	oldFile := filepath.Join(tmp, "old.png")
	if err := os.WriteFile(oldFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-100 * 24 * time.Hour)
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	recentFile := filepath.Join(tmp, "recent.png")
	if err := os.WriteFile(recentFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}

	s := New([]string{".png"}, 30)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Path != oldFile {
		t.Errorf("path = %q, want %q", results[0].Path, oldFile)
	}
	if results[0].Size != 4 {
		t.Errorf("size = %d, want 4", results[0].Size)
	}
}

func TestScan_SkipsRecentFiles(t *testing.T) {
	tmp := t.TempDir()
	recentFile := filepath.Join(tmp, "recent.png")
	if err := os.WriteFile(recentFile, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	// File is brand new, retention is 30 days
	s := New([]string{".png"}, 30)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for recent file, got %d", len(results))
	}
}

func TestScan_NestedDirs(t *testing.T) {
	tmp := t.TempDir()
	nested := filepath.Join(tmp, "a", "b")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}

	f := filepath.Join(nested, "file.jpg")
	if err := os.WriteFile(f, []byte("nested"), 0644); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-100 * 24 * time.Hour)
	if err := os.Chtimes(f, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	s := New([]string{".jpg"}, 30)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Path != f {
		t.Errorf("path = %q, want %q", results[0].Path, f)
	}
}

func TestScan_SkipsNonMatchingExtensions(t *testing.T) {
	tmp := t.TempDir()
	files := []string{"a.txt", "b.pdf", "c.go"}
	for _, name := range files {
		path := filepath.Join(tmp, name)
		if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
			t.Fatal(err)
		}
		oldTime := time.Now().Add(-100 * 24 * time.Hour)
		os.Chtimes(path, oldTime, oldTime)
	}

	s := New([]string{".png"}, 0)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestScan_SkipsDirs(t *testing.T) {
	tmp := t.TempDir()
	dirPath := filepath.Join(tmp, "subdir")
	if err := os.Mkdir(dirPath, 0755); err != nil {
		t.Fatal(err)
	}

	s := New([]string{".png"}, 0)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestScan_NonExistentDir(t *testing.T) {
	s := New([]string{".png"}, 0)
	_, err := s.Scan("/nonexistent/path/12345")
	if err == nil {
		t.Fatal("expected error for non-existent directory")
	}
}

func TestScan_SkipsInaccessibleNestedDir(t *testing.T) {
	tmp := t.TempDir()
	nested := filepath.Join(tmp, "nested")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatal(err)
	}
	// Create a file inside nested dir
	f := filepath.Join(nested, "file.png")
	if err := os.WriteFile(f, []byte("data"), 0644); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().Add(-100 * 24 * time.Hour)
	if err := os.Chtimes(f, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}
	// Remove read permission from nested dir so walk encounters an error
	if err := os.Chmod(nested, 0000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(nested, 0755) // cleanup

	s := New([]string{".png"}, 0)
	results, err := s.Scan(tmp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// The inaccessible dir should be skipped gracefully
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
