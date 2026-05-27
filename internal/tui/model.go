package tui

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
)

type screen int

const (
	screenWelcome screen = iota
	screenFolderPicker
	screenSetup
	screenWarning
	screenStatus
	screenPreview
	screenLogView
	screenResetConfirm
)

type model struct {
	width  int
	height int
	screen screen
	cfg    *config.Config
	err    error

	spinner spinner.Model

	// Sub-models
	welcomeModel      welcomeModel
	folderPickerModel folderPickerModel
	setupModel        setupModel
	warningModel      warningModel
	statusModel       statusModel
	previewModel      previewModel
	logViewModel      logViewModel
	resetModel        resetModel
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		func() tea.Msg {
			cfg, err := config.Load()
			if err != nil {
				return errMsg{err}
			}
			return configLoadedMsg{cfg}
		},
	)
}

type configLoadedMsg struct {
	cfg *config.Config
}

type errMsg struct {
	err error
}

type cleanupDoneMsg struct {
	stats *scanner.CleanupStats
}

type scanDoneMsg struct {
	files []scanner.ScanResult
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "esc" {
			return m, tea.Quit
		}

	case configLoadedMsg:
		m.cfg = msg.cfg
		if config.Exists() {
			m.screen = screenStatus
			m.statusModel = newStatusModel(m.width, m.height, m.cfg)
		} else {
			m.screen = screenWelcome
			m.welcomeModel = newWelcomeModel(m.width, m.height)
		}
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil
	}

	switch m.screen {
	case screenWelcome:
		return m.updateWelcome(msg)
	case screenFolderPicker:
		return m.updateFolderPicker(msg)
	case screenSetup:
		return m.updateSetup(msg)
	case screenWarning:
		return m.updateWarning(msg)
	case screenStatus:
		return m.updateStatus(msg)
	case screenPreview:
		return m.updatePreview(msg)
	case screenLogView:
		return m.updateLogView(msg)
	case screenResetConfirm:
		return m.updateReset(msg)
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Ошибка: %v\n\nНажмите Esc или Ctrl+C для выхода.", m.err)
	}

	switch m.screen {
	case screenWelcome:
		return m.welcomeView()
	case screenFolderPicker:
		return m.folderPickerView()
	case screenSetup:
		return m.setupView()
	case screenWarning:
		return m.warningView()
	case screenStatus:
		return m.statusView()
	case screenPreview:
		return m.previewView()
	case screenLogView:
		return m.logView()
	case screenResetConfirm:
		return m.resetView()
	}

	return "Загрузка..."
}

func containsSystemDir(path string) bool {
	home, _ := os.UserHomeDir()
	systemDirs := []string{
		"/",
		"/System",
		"/Users",
		home,
	}
	for _, d := range systemDirs {
		if path == d {
			return true
		}
	}
	return false
}

func needsWarning(dir string) bool {
	home, _ := os.UserHomeDir()
	warnDirs := []string{
		filepath.Join(home, "Desktop"),
		filepath.Join(home, "Downloads"),
		filepath.Join(home, "Documents"),
		filepath.Join(home, "Movies"),
		filepath.Join(home, "Music"),
		filepath.Join(home, "Pictures"),
		filepath.Join(home, "Public"),
		filepath.Join(home, "Library"),
	}
	for _, d := range warnDirs {
		if dir == d {
			return true
		}
	}
	return false
}
