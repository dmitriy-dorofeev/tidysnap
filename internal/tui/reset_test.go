package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func TestNewResetModel(t *testing.T) {
	m := newResetModel()
	_ = m
}

func TestUpdateReset_Yes(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenResetConfirm
	m.resetModel = newResetModel()
	m.width = 80
	m.height = 24

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m2 := newM.(model)
	if m2.screen != screenWelcome {
		t.Errorf("screen = %d, want screenWelcome", m2.screen)
	}
}

func TestUpdateReset_No(t *testing.T) {
	m := InitialModel()
	cfg := config.DefaultConfig()
	m.cfg = cfg
	m.screen = screenResetConfirm
	m.resetModel = newResetModel()

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
	m2 := newM.(model)
	if m2.screen != screenStatus {
		t.Errorf("screen = %d, want screenStatus", m2.screen)
	}
}

func TestResetView(t *testing.T) {
	m := InitialModel()
	m.cfg = config.DefaultConfig()
	m.screen = screenResetConfirm
	m.width = 80
	m.height = 24
	view := m.resetView()
	if view == "" {
		t.Error("resetView should not be empty")
	}
}
