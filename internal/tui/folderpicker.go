package tui

import (
	"os"

	"github.com/charmbracelet/bubbles/filepicker"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type folderPickerModel struct {
	width  int
	height int
	picker filepicker.Model
}

func newFolderPickerModel(width, height int, startDir string) folderPickerModel {
	fp := filepicker.New()
	fp.DirAllowed = true
	fp.FileAllowed = false
	fp.ShowHidden = false
	fp.AutoHeight = false
	fp.SetHeight(height - 8)
	fp.KeyMap.Select = key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "select"))
	fp.KeyMap.Open = key.NewBinding(key.WithKeys("l", "right", "enter"), key.WithHelp("l", "open"))

	if startDir == "" {
		home, _ := os.UserHomeDir()
		startDir = home
	}
	fp.CurrentDirectory = startDir

	return folderPickerModel{
		width:  width,
		height: height,
		picker: fp,
	}
}

func (m folderPickerModel) Init() tea.Cmd {
	return m.picker.Init()
}

func (m model) updateFolderPicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "esc" {
			m.screen = screenWelcome
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.folderPickerModel.picker, cmd = m.folderPickerModel.picker.Update(msg)

	if didSelect, path := m.folderPickerModel.picker.DidSelectFile(msg); didSelect {
		m.cfg.TargetDir = path
		m.screen = screenSetup
		m.setupModel = newSetupModel(m.width, m.height, m.cfg)
		return m, m.setupModel.Init()
	}

	return m, cmd
}

func (m model) folderPickerView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginBottom(1)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).MarginBottom(1)

	title := titleStyle.Render("📁 Выберите папку для очистки")
	currentPath := pathStyle.Render("Текущая папка: " + m.folderPickerModel.picker.CurrentDirectory)
	view := m.folderPickerModel.picker.View()
	hints := hintStyle.Render("↑/↓ навигация • ← назад • → открыть папку • Пробел — выбрать папку • q — назад")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, currentPath, view, hints),
	)
}
