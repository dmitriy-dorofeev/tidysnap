package tui

import (
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

type logViewModel struct {
	viewport viewport.Model
	content  string
}

func newLogViewModel(width, height int) logViewModel {
	vp := viewport.New(width, height-2)
	return logViewModel{viewport: vp}
}

func (m logViewModel) Init() tea.Cmd {
	return func() tea.Msg {
		logPath, err := config.LogPath()
		if err != nil {
			return errMsg{err}
		}
		// #nosec G304 — logPath is an internal log path, not user input.
		data, err := os.ReadFile(logPath)
		if err != nil {
			return errMsg{err}
		}
		return logLoadedMsg{content: string(data)}
	}
}

type logLoadedMsg struct {
	content string
}

func (m model) updateLogView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "q":
			m.screen = screenStatus
			return m, nil
		}
	case logLoadedMsg:
		m.logViewModel.content = msg.content
		if m.logViewModel.content == "" {
			m.logViewModel.content = "Лог пуст."
		}
		lines := strings.Split(m.logViewModel.content, "\n")
		if len(lines) > m.logViewModel.viewport.Height {
			m.logViewModel.viewport.SetYOffset(len(lines) - m.logViewModel.viewport.Height)
		}
		m.logViewModel.viewport.SetContent(m.logViewModel.content)
		return m, nil
	case errMsg:
		m.logViewModel.content = "Лог пуст или недоступен."
		m.logViewModel.viewport.SetContent(m.logViewModel.content)
		return m, nil
	}

	var cmd tea.Cmd
	m.logViewModel.viewport, cmd = m.logViewModel.viewport.Update(msg)
	return m, cmd
}

func (m model) logView() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	hints := hintStyle.Render("[↑/↓] Прокрутка  [Esc/q] Назад")
	return lipgloss.JoinVertical(lipgloss.Left, m.logViewModel.viewport.View(), hints)
}
