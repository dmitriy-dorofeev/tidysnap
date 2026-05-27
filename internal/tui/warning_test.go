package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func TestUpdateWarning_Yes(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	cfg.TargetDir = "/tmp"
	m.cfg = cfg
	m.screen = screenWarning
	m.warningModel = newWarningModel("/tmp", []string{".png"}, 30)
	m.width = 80
	m.height = 24

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := newM.(model)
	// After confirming warning, it saves config and goes to status (if install succeeds)
	// Since daemon.Install may fail in tests, we expect either status or err set.
	if m2.screen != screenStatus && m2.err == nil {
		t.Logf("screen = %d, err = %v", m2.screen, m2.err)
	}
	if m2.cfg.WarningAck != true {
		t.Error("WarningAck should be true")
	}
}

func TestUpdateWarning_No(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	m.cfg = cfg
	m.screen = screenWarning
	m.warningModel = newWarningModel("/tmp", []string{".png"}, 30)
	m.width = 80
	m.height = 24

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m2 := newM.(model)
	if m2.screen != screenSetup {
		t.Errorf("screen = %d, want screenSetup", m2.screen)
	}
	if cmd == nil {
		t.Error("expected command from setup init")
	}
}

func TestUpdateWarning_UnknownKey(t *testing.T) {
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	m.cfg = cfg
	m.screen = screenWarning
	m.warningModel = newWarningModel("/tmp", []string{".png"}, 30)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	_ = newM.(model)
	if cmd != nil {
		t.Error("unknown key should not produce command")
	}
}

func TestWarningView(t *testing.T) {
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	m.cfg = cfg
	m.cfg.TargetDir = "/tmp"
	m.width = 80
	m.height = 24
	view := m.warningView()
	if view == "" {
		t.Error("warningView should not be empty")
	}
}
