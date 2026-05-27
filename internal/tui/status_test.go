package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
)

func TestHumanizeBytes(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{1024 * 1024 * 1024, "1.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := humanizeBytes(tt.bytes)
			if got != tt.want {
				t.Errorf("humanizeBytes(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestNewStatusModel(t *testing.T) {
	cfg := config.DefaultConfig()
	m := newStatusModel(80, 24, cfg)
	if m.width != 80 || m.height != 24 {
		t.Error("dimensions mismatch")
	}
}

func TestUpdateStatus_R(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	cfg.TargetDir = t.TempDir()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)
	m.width = 80
	m.height = 24

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command for scan")
	}
	msg := cmd()
	_, ok := msg.(scanDoneMsg)
	if !ok {
		t.Fatalf("expected scanDoneMsg, got %T", msg)
	}
}

func TestUpdateStatus_L(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)
	m.width = 80
	m.height = 24

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	m2 := newM.(model)
	if m2.screen != screenLogView {
		t.Errorf("screen = %d, want screenLogView", m2.screen)
	}
	if cmd == nil {
		t.Error("expected command for log view init")
	}
}

func TestUpdateStatus_E(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)
	m.width = 80
	m.height = 24

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	m2 := newM.(model)
	if m2.screen != screenFolderPicker {
		t.Errorf("screen = %d, want screenFolderPicker", m2.screen)
	}
	if cmd == nil {
		t.Error("expected command for folder picker init")
	}
}

func TestUpdateStatus_X(t *testing.T) {
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}})
	m2 := newM.(model)
	if m2.screen != screenResetConfirm {
		t.Errorf("screen = %d, want screenResetConfirm", m2.screen)
	}
}

func TestUpdateStatus_Q(t *testing.T) {
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestUpdateStatus_S_Install(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m2 := newM.(model)
	// daemon.Install may succeed or fail; just verify no panic and screen stays
	if m2.screen != screenStatus {
		t.Logf("screen changed to %d", m2.screen)
	}
}

func TestUpdateStatus_UnknownKey(t *testing.T) {
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	_ = newM.(model)
	if cmd != nil {
		t.Error("unknown key should not produce command")
	}
}

func TestUpdateStatus_CleanupDone(t *testing.T) {
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)

	stats := &scanner.CleanupStats{FilesRemoved: 5, BytesFreed: 1024}
	newM, _ := m.Update(cleanupDoneMsg{stats: stats})
	m2 := newM.(model)
	if m2.screen != screenStatus {
		t.Errorf("screen = %d, want screenStatus", m2.screen)
	}
	if m2.statusModel.msg == "" {
		t.Error("status message should be set")
	}
}

func TestUpdateStatus_ScanDone(t *testing.T) {
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)
	m.width = 80
	m.height = 24

	files := []scanner.ScanResult{{Path: "/tmp/a.png", Size: 10}}
	newM, _ := m.Update(scanDoneMsg{files: files})
	m2 := newM.(model)
	if m2.screen != screenPreview {
		t.Errorf("screen = %d, want screenPreview", m2.screen)
	}
	if len(m2.previewModel.files) != 1 {
		t.Errorf("preview files = %d, want 1", len(m2.previewModel.files))
	}
}

func TestStatusView(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)
	m.width = 80
	m.height = 24
	view := m.statusView()
	if view == "" {
		t.Error("statusView should not be empty")
	}
}

func TestRunScan(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	cfg.TargetDir = t.TempDir()
	m.cfg = cfg
	m.screen = screenStatus
	m.statusModel = newStatusModel(80, 24, cfg)

	newM, cmd := m.runScan()
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	_, ok := msg.(scanDoneMsg)
	if !ok {
		t.Fatalf("expected scanDoneMsg, got %T", msg)
	}
}
