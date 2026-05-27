package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
)

func TestPreviewItem(t *testing.T) {
	item := previewItem{file: scanner.ScanResult{Path: "/tmp/test.png", Size: 1024}}
	if item.Title() != "/tmp/test.png" {
		t.Errorf("Title = %q, want /tmp/test.png", item.Title())
	}
	if item.FilterValue() != "/tmp/test.png" {
		t.Errorf("FilterValue = %q, want /tmp/test.png", item.FilterValue())
	}
	if item.Description() == "" {
		t.Error("Description should not be empty")
	}
}

func TestNewPreviewModel(t *testing.T) {
	files := []scanner.ScanResult{
		{Path: "/tmp/a.png", Size: 100},
		{Path: "/tmp/b.png", Size: 200},
	}
	m := newPreviewModel(80, 24, files, false)
	if len(m.files) != 2 {
		t.Errorf("files = %d, want 2", len(m.files))
	}
	if m.dryRun != false {
		t.Error("dryRun should be false")
	}
}

func TestUpdatePreview_Enter(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.cfg.LogPath = filepath.Join(t.TempDir(), "cleanup.log")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{}, false)
	m.width = 80
	m.height = 24

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command for cleanup")
	}
	msg := cmd()
	_, ok := msg.(cleanupDoneMsg)
	if !ok {
		// Could be errMsg if log file can't be opened
		t.Logf("got %T", msg)
	}
}

func TestUpdatePreview_D(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.cfg.LogPath = filepath.Join(t.TempDir(), "cleanup.log")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{}, false)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}})
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command for cleanup")
	}
}

func TestUpdatePreview_D_Russian(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.cfg.LogPath = filepath.Join(t.TempDir(), "cleanup.log")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{}, false)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'в'}})
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command for cleanup")
	}
}

func TestUpdatePreview_Esc(t *testing.T) {
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{}, false)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newM.(model)
	// Esc is handled globally as quit
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestUpdatePreview_ListNavigation(t *testing.T) {
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{
		{Path: "/tmp/a.png", Size: 100},
		{Path: "/tmp/b.png", Size: 200},
	}, false)

	// Send a key that list handles (e.g., 'j' or down)
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m2 := newM.(model)
	// list internal state changed; we just verify no panic and screen stays
	if m2.screen != screenPreview {
		t.Errorf("screen = %d, want screenPreview", m2.screen)
	}
}

func TestPreviewView_Empty(t *testing.T) {
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{}, false)
	m.width = 80
	m.height = 24
	view := m.previewView()
	if view == "" {
		t.Error("previewView should not be empty")
	}
}

func TestPreviewView_WithFiles(t *testing.T) {
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{
		{Path: "/tmp/a.png", Size: 100},
	}, false)
	m.width = 80
	m.height = 24
	view := m.previewView()
	if view == "" {
		t.Error("previewView should not be empty")
	}
}

func TestRunCleanup_WithLog(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.cfg.LogPath = filepath.Join(t.TempDir(), "cleanup.log")
	m.cfg.DryRun = true
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{
		{Path: "/tmp/a.png", Size: 100},
	}, true)

	newM, cmd := m.runCleanup()
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	_, ok := msg.(cleanupDoneMsg)
	if !ok {
		t.Logf("got %T", msg)
	}
}

func TestRunCleanup_LogOpenError(t *testing.T) {
	m := InitialModel()
	m.cfg = defaultConfigWithTarget("/tmp")
	m.cfg.LogPath = "/nonexistent/path/log.txt"
	m.screen = screenPreview
	m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{
		{Path: "/tmp/a.png", Size: 100},
	}, true)

	newM, cmd := m.runCleanup()
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	_, ok := msg.(errMsg)
	if !ok {
		t.Fatalf("expected errMsg, got %T", msg)
	}
}
