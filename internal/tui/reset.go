package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
)

type resetModel struct{}

func newResetModel() resetModel {
	return resetModel{}
}

func (m resetModel) Init() tea.Cmd { return nil }

func (m model) updateReset(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "y":
			_ = daemon.Uninstall()
			_ = config.Cleanup()
			m.cfg = config.DefaultConfig()
			m.screen = screenWelcome
			m.welcomeModel = newWelcomeModel(m.width, m.height)
			return m, nil
		case "esc", "n":
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

	content := fmt.Sprintf(
		"Будут удалены все файлы и настройки TidySnap:\n\n"+
			"• %s\n"+
			"• %s\n"+
			"• %s\n\n"+
			"Демон будет остановлен и выгружен.\n\n"+
			"Это действие %s.",
		boldStyle.Render(config.ConfigPath()),
		boldStyle.Render(config.LogDir()),
		boldStyle.Render(config.PlistPath()),
		boldStyle.Render("нельзя отменить"),
	)

	box := boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("⚠️  Удаление TidySnap"),
			content,
		),
	)
	hints := hintStyle.Render("[y/Enter] Удалить всё  [n/Esc] Отмена")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, box, hints),
	)
}
