package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
)

func setTestHome(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmp)
	t.Cleanup(func() { os.Setenv("HOME", oldHome) })
	return tmp
}

func TestInitialModel(t *testing.T) {
	m := InitialModel()
	if m.screen != 0 {
		t.Errorf("initial screen = %d, want 0", m.screen)
	}
	if m.width != 0 || m.height != 0 {
		t.Error("initial dimensions should be 0")
	}
}

func TestInit(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("Init should return a command")
	}
	// Init returns tea.Batch which produces tea.BatchMsg containing inner commands.
	msg := cmd()
	if msg == nil {
		t.Fatal("Init command should return a message")
	}
	batch, ok := msg.(tea.BatchMsg)
	if !ok {
		t.Fatalf("expected tea.BatchMsg, got %T", msg)
	}
	if len(batch) == 0 {
		t.Fatal("BatchMsg should contain commands")
	}
	// Execute the second command (config load)
	configMsg := batch[1]()
	if _, ok := configMsg.(configLoadedMsg); !ok {
		t.Fatalf("expected configLoadedMsg, got %T", configMsg)
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	m := InitialModel()
	newM, cmd := m.Update(tea.WindowSizeMsg{Width: 100, Height: 50})
	model := newM.(model)
	if model.width != 100 || model.height != 50 {
		t.Errorf("size = (%d, %d), want (100, 50)", model.width, model.height)
	}
	if cmd != nil {
		t.Error("WindowSize should not return a command")
	}
}

func TestUpdate_WindowSize_TooSmall(t *testing.T) {
	m := InitialModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 79, Height: 24})
	mm := newM.(model)
	if !mm.tooSmall {
		t.Error("expected tooSmall to be true for width < minWidth")
	}

	m = InitialModel()
	newM, _ = m.Update(tea.WindowSizeMsg{Width: 80, Height: 23})
	mm = newM.(model)
	if !mm.tooSmall {
		t.Error("expected tooSmall to be true for height < minHeight")
	}
}

func TestUpdate_WindowSize_Adequate(t *testing.T) {
	m := InitialModel()
	newM, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	mm := newM.(model)
	if mm.tooSmall {
		t.Error("expected tooSmall to be false for exact minimum size")
	}
}

func TestView_TooSmall(t *testing.T) {
	m := InitialModel()
	m.tooSmall = true
	m.width = 40
	m.height = 10
	view := m.View()
	if view == "" {
		t.Error("View should not be empty when tooSmall is set")
	}
	if !strings.Contains(view, "80") || !strings.Contains(view, "24") {
		t.Error("tooSmall view should mention minimum dimensions")
	}
}

func TestUpdate_CtrlC(t *testing.T) {
	m := InitialModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("expected quit command")
	}
	// Verify it quits
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("expected QuitMsg")
	}
}

func TestUpdate_Esc(t *testing.T) {
	m := InitialModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	if cmd == nil {
		t.Fatal("expected quit command")
	}
}

func TestUpdate_ConfigLoaded_WithoutConfig(t *testing.T) {
	setTestHome(t)
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	newM, cmd := m.Update(configLoadedMsg{cfg: cfg})
	model := newM.(model)
	if model.screen != screenWelcome {
		t.Errorf("screen = %d, want screenWelcome", model.screen)
	}
	if cmd != nil {
		t.Error("should not return a command")
	}
}

func TestUpdate_ConfigLoaded_WithConfig(t *testing.T) {
	tmp := setTestHome(t)
	m := InitialModel()
	cfg, err := config.DefaultConfig()
	if err != nil {
		t.Fatal(err)
	}
	// Save config so Exists returns true
	if err := config.Save(cfg); err != nil {
		t.Fatal(err)
	}
	newM, cmd := m.Update(configLoadedMsg{cfg: cfg})
	model := newM.(model)
	if model.screen != screenStatus {
		t.Errorf("screen = %d, want screenStatus", model.screen)
	}
	if cmd != nil {
		t.Error("should not return a command")
	}
	_ = tmp
}

func TestUpdate_ErrMsg(t *testing.T) {
	m := InitialModel()
	err := os.ErrInvalid
	newM, _ := m.Update(errMsg{err: err})
	model := newM.(model)
	if model.err != err {
		t.Error("err not set correctly")
	}
}

func TestView_Error(t *testing.T) {
	m := InitialModel()
	m.err = os.ErrInvalid
	view := m.View()
	if view == "" {
		t.Error("View should not be empty when err is set")
	}
}

func TestView_Loading(t *testing.T) {
	m := InitialModel()
	view := m.View()
	if view == "" {
		t.Error("View should not be empty")
	}
}

func TestView_AllScreens(t *testing.T) {
	screens := []screen{
		screenWelcome,
		screenFolderPicker,
		screenSetup,
		screenWarning,
		screenStatus,
		screenPreview,
		screenLogView,
		screenResetConfirm,
	}
	for _, s := range screens {
		m := InitialModel()
		m.screen = s
		m.width = 80
		m.height = 24
		cfg, err := config.DefaultConfig()
		if err != nil {
			t.Fatal(err)
		}
		m.cfg = cfg
		m.folderPickerModel = newFolderPickerModel(80, 24, "/tmp", screenWelcome)
		m.folderPickerModel.items = []dirItem{{name: "tmp", path: "/tmp"}}
		m.welcomeModel = newWelcomeModel(80, 24)
		m.setupModel = newSetupModel(80, 24, m.cfg)
		m.warningModel = newWarningModel("/tmp", []string{".png"}, 30)
		m.statusModel = newStatusModel(80, 24, m.cfg)
		m.previewModel = newPreviewModel(80, 24, []scanner.ScanResult{}, false)
		m.logViewModel = newLogViewModel(80, 24)
		m.logViewModel.content = "test"
		m.logViewModel.viewport.SetContent("test")
		m.resetModel = newResetModel()

		view := m.View()
		if view == "" {
			t.Errorf("View for screen %d should not be empty", s)
		}
	}
}

func TestContainsSystemDir(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		path string
		want bool
	}{
		{"/", true},
		{"/System", true},
		{"/Users", true},
		{home, true},
		{"/tmp", false},
		{"/Users/foo/Downloads", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := containsSystemDir(tt.path)
			if got != tt.want {
				t.Errorf("containsSystemDir(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNeedsWarning(t *testing.T) {
	home, _ := os.UserHomeDir()
	tests := []struct {
		dir  string
		want bool
	}{
		{filepath.Join(home, "Desktop"), true},
		{filepath.Join(home, "Downloads"), true},
		{filepath.Join(home, "Documents"), true},
		{filepath.Join(home, "Movies"), true},
		{filepath.Join(home, "Music"), true},
		{filepath.Join(home, "Pictures"), true},
		{filepath.Join(home, "Public"), true},
		{filepath.Join(home, "Library"), true},
		{filepath.Join(home, "Screenshots"), false},
		{"/tmp", false},
	}

	for _, tt := range tests {
		t.Run(tt.dir, func(t *testing.T) {
			got := needsWarning(tt.dir)
			if got != tt.want {
				t.Errorf("needsWarning(%q) = %v, want %v", tt.dir, got, tt.want)
			}
		})
	}
}
