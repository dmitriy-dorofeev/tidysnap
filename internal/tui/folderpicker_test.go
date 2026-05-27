package tui

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

func TestNewFolderPickerModel(t *testing.T) {
	tmp := t.TempDir()
	m := newFolderPickerModel(80, 24, tmp, screenWelcome)
	if m.cwd != tmp {
		t.Errorf("cwd = %q, want %q", m.cwd, tmp)
	}
	if m.returnScreen != screenWelcome {
		t.Errorf("returnScreen = %d, want screenWelcome", m.returnScreen)
	}
}

func TestFolderPickerLoad(t *testing.T) {
	tmp := setTestHome(t)
	// Create subdirectories inside home so cwd == home and no parent item is added
	os.MkdirAll(filepath.Join(tmp, "alpha"), 0755)
	os.MkdirAll(filepath.Join(tmp, "beta"), 0755)

	m := newFolderPickerModel(80, 24, tmp, screenWelcome)
	cmd := m.Init()
	msg := cmd()

	loaded, ok := msg.(folderPickerLoadedMsg)
	if !ok {
		t.Fatalf("expected folderPickerLoadedMsg, got %T", msg)
	}
	if len(loaded.items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(loaded.items))
	}
	if loaded.items[0].name != "alpha" || loaded.items[1].name != "beta" {
		t.Errorf("items = %v", loaded.items)
	}
}

func TestFolderPickerLoad_WithParent(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	os.MkdirAll(sub, 0755)

	m := newFolderPickerModel(80, 24, sub, screenWelcome)
	cmd := m.Init()
	msg := cmd()

	loaded, ok := msg.(folderPickerLoadedMsg)
	if !ok {
		t.Fatalf("expected folderPickerLoadedMsg, got %T", msg)
	}
	foundParent := false
	for _, it := range loaded.items {
		if it.isParent {
			foundParent = true
		}
	}
	if !foundParent {
		t.Error("expected parent item '..' when not in home")
	}
}

func TestUpdateFolderPicker_Navigation(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	os.MkdirAll(sub, 0755)

	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)
	m.folderPickerModel.items = []dirItem{
		{name: "sub", path: sub},
	}

	// down
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	m2 := newM.(model)
	if m2.folderPickerModel.cursor != 0 {
		// only one item, cursor stays at 0
		t.Logf("cursor = %d", m2.folderPickerModel.cursor)
	}

	// up should not go below 0
	newM, _ = m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	m3 := newM.(model)
	if m3.folderPickerModel.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m3.folderPickerModel.cursor)
	}

	// enter should navigate into sub
	newM, cmd := m3.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m4 := newM.(model)
	if m4.folderPickerModel.cwd != sub {
		t.Errorf("cwd = %q, want %q", m4.folderPickerModel.cwd, sub)
	}
	if cmd == nil {
		t.Error("expected load command after entering directory")
	}
}

func TestUpdateFolderPicker_Navigation_Russian(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	os.MkdirAll(sub, 0755)

	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)
	m.folderPickerModel.items = []dirItem{
		{name: "sub", path: sub},
	}

	// down (Russian 'о')
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'о'}})
	m2 := newM.(model)
	if m2.folderPickerModel.cursor != 0 {
		t.Logf("cursor = %d", m2.folderPickerModel.cursor)
	}

	// up should not go below 0 (Russian 'л')
	newM, _ = m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'л'}})
	m3 := newM.(model)
	if m3.folderPickerModel.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m3.folderPickerModel.cursor)
	}

	// enter should navigate into sub
	newM, cmd := m3.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m4 := newM.(model)
	if m4.folderPickerModel.cwd != sub {
		t.Errorf("cwd = %q, want %q", m4.folderPickerModel.cwd, sub)
	}
	if cmd == nil {
		t.Error("expected load command after entering directory")
	}
}

func TestUpdateFolderPicker_Back(t *testing.T) {
	tmp := t.TempDir()
	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m2 := newM.(model)
	if m2.screen != screenWelcome {
		t.Errorf("screen = %d, want screenWelcome", m2.screen)
	}
}

func TestUpdateFolderPicker_Back_Russian(t *testing.T) {
	tmp := t.TempDir()
	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)

	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'й'}})
	m2 := newM.(model)
	if m2.screen != screenWelcome {
		t.Errorf("screen = %d, want screenWelcome", m2.screen)
	}
}

func TestUpdateFolderPicker_Select(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	os.MkdirAll(sub, 0755)

	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	m.cfg = cfg
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)
	m.folderPickerModel.items = []dirItem{{name: "sub", path: sub}}

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	m2 := newM.(model)
	if m2.screen != screenSetup {
		t.Errorf("screen = %d, want screenSetup", m2.screen)
	}
	if m2.cfg.TargetDir != sub {
		t.Errorf("TargetDir = %q, want %q", m2.cfg.TargetDir, sub)
	}
	if cmd == nil {
		t.Error("expected command from setup init")
	}
}

func TestUpdateFolderPicker_Left(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	os.MkdirAll(sub, 0755)

	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, sub, screenWelcome)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}})
	m2 := newM.(model)
	if m2.folderPickerModel.cwd != tmp {
		t.Errorf("cwd = %q, want %q", m2.folderPickerModel.cwd, tmp)
	}
	if cmd == nil {
		t.Error("expected load command after going up")
	}
}

func TestUpdateFolderPicker_Left_Russian(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	os.MkdirAll(sub, 0755)

	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, sub, screenWelcome)

	newM, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'р'}})
	m2 := newM.(model)
	if m2.folderPickerModel.cwd != tmp {
		t.Errorf("cwd = %q, want %q", m2.folderPickerModel.cwd, tmp)
	}
	if cmd == nil {
		t.Error("expected load command after going up")
	}
}

func TestUpdateFolderPicker_FolderPickerLoadedMsg(t *testing.T) {
	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, "/tmp", screenWelcome)
	m.folderPickerModel.cursor = 5

	items := []dirItem{{name: "a", path: "/tmp/a"}}
	newM, _ := m.Update(folderPickerLoadedMsg{items: items})
	m2 := newM.(model)
	if m2.folderPickerModel.cursor != 0 {
		t.Errorf("cursor = %d, want 0", m2.folderPickerModel.cursor)
	}
	if len(m2.folderPickerModel.items) != 1 {
		t.Errorf("items = %d, want 1", len(m2.folderPickerModel.items))
	}
}

func TestUpdateFolderPicker_FolderPickerErrorMsg(t *testing.T) {
	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, "/tmp", screenWelcome)

	newM, _ := m.Update(folderPickerErrorMsg{err: os.ErrInvalid})
	m2 := newM.(model)
	if m2.err != os.ErrInvalid {
		t.Error("err not set correctly")
	}
}

func TestUpdateFolderPicker_FolderSelectedMsg(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	m.cfg = cfg
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, "/tmp", screenWelcome)

	newM, cmd := m.Update(folderSelectedMsg{path: "/tmp/selected"})
	m2 := newM.(model)
	if m2.screen != screenSetup {
		t.Errorf("screen = %d, want screenSetup", m2.screen)
	}
	if m2.cfg.TargetDir != "/tmp/selected" {
		t.Errorf("TargetDir = %q, want /tmp/selected", m2.cfg.TargetDir)
	}
	if cmd == nil {
		t.Error("expected command")
	}
}

func TestUpdateFolderPicker_FolderPickerBackMsg(t *testing.T) {
	m := InitialModel()
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, "/tmp", screenWelcome)

	newM, _ := m.Update(folderPickerBackMsg{})
	m2 := newM.(model)
	if m2.screen != screenWelcome {
		t.Errorf("screen = %d, want screenWelcome", m2.screen)
	}
}

func TestFolderPickerView_Empty(t *testing.T) {
	tmp := t.TempDir()
	m := InitialModel()
	m.width = 80
	m.height = 24
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)
	m.folderPickerModel.items = []dirItem{}
	view := m.folderPickerView()
	if view == "" {
		t.Error("folderPickerView should not be empty")
	}
}

func TestFolderPickerView_Clipping(t *testing.T) {
	tmp := t.TempDir()
	m := InitialModel()
	m.width = 80
	m.height = 10 // very small height to trigger clipping
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 10, tmp, screenWelcome)
	// Add many items
	for i := 0; i < 20; i++ {
		m.folderPickerModel.items = append(m.folderPickerModel.items, dirItem{name: "dir", path: tmp})
	}
	view := m.folderPickerView()
	if view == "" {
		t.Error("folderPickerView should not be empty")
	}
}

func TestFolderPickerView(t *testing.T) {
	tmp := t.TempDir()
	m := InitialModel()
	m.width = 80
	m.height = 24
	m.screen = screenFolderPicker
	m.folderPickerModel = newFolderPickerModel(80, 24, tmp, screenWelcome)
	m.folderPickerModel.items = []dirItem{{name: "foo", path: filepath.Join(tmp, "foo")}}
	view := m.folderPickerView()
	if view == "" {
		t.Error("folderPickerView should not be empty")
	}
}
