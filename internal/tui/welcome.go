package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
		switch msg.String() {
		case "enter", "s":
			m.screen = screenSetup
			m.setupModel = newSetupModel(m.width, m.height, m.cfg)
			return m, m.setupModel.Init()
		case "q", "esc":
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

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(2)

	title := titleStyle.Render("🧹 Добро пожаловать в TidySnap!")
	subtitle := subtitleStyle.Render("Автоочистка скриншотов и записей экрана")

	content := fmt.Sprintf(
		"%s\n\n%s\n\n%s",
		"Эта утилита поможет вам автоматически удалять старые скриншоты и видеозаписи экрана.",
		"При первом запуске нужно выбрать папку и настроить параметры.",
		"После настройки TidySnap будет работать в фоне через launchd.",
	)

	box := boxStyle.Render(content)
	hints := hintStyle.Render("[s/Enter] Начать настройку  [q] Выход")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, subtitle, box, hints),
	)
}
