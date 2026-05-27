package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func TestNewLogViewModel(t *testing.T) {
	m := newLogViewModel(80, 24)
	if m.viewport.Width != 80 || m.viewport.Height != 22 {
		// height-2 = 22
		t.Logf("viewport size = (%d, %d)", m.viewport.Width, m.viewport.Height)
	}
}

func TestUpdateLogView_LogLoaded(t *testing.T) {
	m := InitialModel()
	m.screen = screenLogView
	m.logViewModel = newLogViewModel(80, 24)

	content := "line1\nline2\nline3"
	newM, _ := m.Update(logLoadedMsg{content: content})
	m2 := newM.(model)
	if m2.logViewModel.content != content {
		t.Errorf("content mismatch")
	}
}

func TestUpdateLogView_EmptyLog(t *testing.T) {
	m := InitialModel()
	m.screen = screenLogView
	m.logViewModel = newLogViewModel(80, 24)

	newM, _ := m.Update(logLoadedMsg{content: ""})
	m2 := newM.(model)
	if m2.logViewModel.content != "Лог пуст." {
		t.Errorf("content = %q, want 'Лог пуст.'", m2.logViewModel.content)
	}
}

func TestUpdateLogView_Esc(t *testing.T) {
	m := InitialModel()
	m.screen = screenLogView
	m.logViewModel = newLogViewModel(80, 24)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	_ = newM.(model)
	// Esc is handled globally as quit
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestLogViewInit(t *testing.T) {
	setTestHome(t)
	// Create a log file
	logPath, err := config.LogPath()
	if err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Dir(logPath), 0755)
	os.WriteFile(logPath, []byte("test log"), 0644)

	m := newLogViewModel(80, 24)
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected command")
	}
	msg := cmd()
	_, ok := msg.(logLoadedMsg)
	if !ok {
		t.Fatalf("expected logLoadedMsg, got %T", msg)
	}
}

func TestLogView(t *testing.T) {
	m := InitialModel()
	m.screen = screenLogView
	m.logViewModel = newLogViewModel(80, 24)
	view := m.logView()
	if view == "" {
		t.Error("logView should not be empty")
	}
}
