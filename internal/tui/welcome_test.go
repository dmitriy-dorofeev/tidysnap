package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func TestUpdateWelcome_Enter(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	newM, _ := m.Update(configLoadedMsg{cfg: defaultConfigWithTarget("/tmp")})
	mod := newM.(model)
	mod.screen = screenWelcome
	mod.welcomeModel = newWelcomeModel(80, 24)

	newM2, cmd := mod.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := newM2.(model)
	if m2.screen != screenFolderPicker {
		t.Errorf("screen = %d, want screenFolderPicker", m2.screen)
	}
	if cmd == nil {
		t.Error("expected command from folder picker init")
	}
}

func TestUpdateWelcome_S(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	newM, _ := m.Update(configLoadedMsg{cfg: defaultConfigWithTarget("/tmp")})
	mod := newM.(model)
	mod.screen = screenWelcome
	mod.welcomeModel = newWelcomeModel(80, 24)

	newM2, cmd := mod.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	m2 := newM2.(model)
	if m2.screen != screenFolderPicker {
		t.Errorf("screen = %d, want screenFolderPicker", m2.screen)
	}
	if cmd == nil {
		t.Error("expected command")
	}
}

func TestUpdateWelcome_Q(t *testing.T) {
	m := InitialModel()
	m.screen = screenWelcome
	m.welcomeModel = newWelcomeModel(80, 24)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	_ = newM.(model)
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestUpdateWelcome_UnknownKey(t *testing.T) {
	m := InitialModel()
	m.screen = screenWelcome
	m.welcomeModel = newWelcomeModel(80, 24)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'z'}})
	_ = newM.(model)
	if cmd != nil {
		t.Error("unknown key should not produce command")
	}
}

func TestWelcomeView(t *testing.T) {
	m := newWelcomeModel(80, 24)
	view := m.Init()
	_ = view
	// Just ensure it doesn't panic; visual testing is hard for TUI.
}

func defaultConfigWithTarget(target string) *config.Config {
	cfg := config.DefaultConfig()
	cfg.TargetDir = target
	return cfg
}
