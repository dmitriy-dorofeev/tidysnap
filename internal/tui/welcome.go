package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/i18n"
)

type welcomeModel struct {
	width  int
	height int
}

func newWelcomeModel(width, height int) welcomeModel {
	return welcomeModel{width: width, height: height}
}

func (m welcomeModel) Init() tea.Cmd { return nil }

func (m model) updateWelcome(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "enter" || keyMatches(msg, 's'):
			m.screen = screenFolderPicker
			startDir := m.cfg.TargetDir
			if startDir == "" {
				home, _ := os.UserHomeDir()
				startDir = home
			}
			m.folderPickerModel = newFolderPickerModel(m.width, m.height, startDir, screenWelcome)
			return m, m.folderPickerModel.Init()
		case keyMatches(msg, 'q') || msg.String() == "esc":
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) welcomeView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		MarginBottom(1).
		MarginTop(2)

	subtitleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("250")).
		MarginBottom(2)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("86")).
		Padding(1, 2).
		MarginBottom(1)

	tipStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("226")).
		Italic(true).
		MarginTop(1)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2)

	title := titleStyle.Render(i18n.T("welcome_title"))
	subtitle := subtitleStyle.Render(i18n.T("welcome_subtitle"))

	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		i18n.T("welcome_desc1"),
		i18n.T("welcome_desc2"),
		i18n.T("welcome_desc3"),
		tipStyle.Render(i18n.T("welcome_tip")),
	)

	box := boxStyle.Render(content)
	hints := hintStyle.Render(i18n.T("welcome_hints"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, box, hints),
	)
}
