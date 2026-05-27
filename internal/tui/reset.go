package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
	"github.com/dmitriy-dorofeev/tidysnap/internal/i18n"
)

type resetModel struct{}

func newResetModel() resetModel {
	return resetModel{}
}

func (m resetModel) Init() tea.Cmd { return nil }

func (m model) updateReset(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "enter" || keyMatches(msg, 'y'):
			_ = daemon.Uninstall()
			_ = config.Cleanup()
			cfg, err := config.DefaultConfig()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.cfg = cfg
			m.screen = screenWelcome
			m.welcomeModel = newWelcomeModel(m.width, m.height)
			return m, nil
		case keyMatches(msg, 'n') || msg.String() == "esc":
			m.screen = screenStatus
			return m, nil
		}
	}
	return m, nil
}

func (m model) resetView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")).MarginBottom(1)
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2).
		Width(60)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
	boldStyle := lipgloss.NewStyle().Bold(true)

	configPath, _ := config.ConfigPath()
	logDir, _ := config.LogDir()
	plistPath, _ := config.PlistPath()

	content := fmt.Sprintf(
		i18n.T("reset_text1")+"\n\n"+
			"• %s\n"+
			"• %s\n"+
			"• %s\n\n"+
			i18n.T("reset_text2")+"\n\n"+
			i18n.T("reset_text3"),
		boldStyle.Render(configPath),
		boldStyle.Render(logDir),
		boldStyle.Render(plistPath),
		boldStyle.Render(i18n.T("reset_irreversible")),
	)

	box := boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(i18n.T("reset_title")),
			content,
		),
	)
	hints := hintStyle.Render(i18n.T("reset_hints"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, box, hints),
	)
}
