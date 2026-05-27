package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
)

type warningModel struct {
	targetDir  string
	extensions []string
	retention  int
}

func newWarningModel(targetDir string, extensions []string, retention int) warningModel {
	return warningModel{
		targetDir:  targetDir,
		extensions: extensions,
		retention:  retention,
	}
}

func (m model) updateWarning(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "y":
			m.cfg.WarningAck = true
			if err := config.Save(m.cfg); err != nil {
				m.err = err
				return m, nil
			}
			return m.saveAndGoToStatus()
		case "esc", "n":
			m.screen = screenSetup
			m.setupModel = newSetupModel(m.width, m.height, m.cfg)
			return m, m.setupModel.Init()
		}
	}
	return m, nil
}

func (m model) saveAndGoToStatus() (tea.Model, tea.Cmd) {
	if err := config.Save(m.cfg); err != nil {
		m.err = err
		return m, nil
	}
	if err := os.MkdirAll(config.LogDir(), 0750); err != nil {
		m.err = err
		return m, nil
	}
	binary := daemon.BinaryPath()
	if err := daemon.Install(binary, m.cfg.CheckIntervalHours); err != nil {
		m.err = err
		return m, nil
	}
	m.screen = screenStatus
	m.statusModel = newStatusModel(m.width, m.height, m.cfg)
	return m, nil
}

func (m model) warningView() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("196")).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("196")).
		Padding(1, 2).
		Width(60)

	hintStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		MarginTop(1)

	boldStyle := lipgloss.NewStyle().Bold(true)

	exts := strings.Join(m.cfg.Extensions, ", ")

	content := fmt.Sprintf(
		"Указана папка: %s\n\n"+
			"Утилита будет удалять ВСЕ файлы с расширениями:\n"+
			"%s\n"+
			"старше %d дней, независимо от имени файла.\n\n"+
			"Если в этой папке есть важные файлы — они БУДУТ УДАЛЕНЫ.\n\n"+
			"%s\n"+
			"Используйте отдельную папку только для скриншотов,\n"+
			"например: ~/Screenshots",
		boldStyle.Render(m.cfg.TargetDir),
		boldStyle.Render(exts),
		m.cfg.RetentionDays,
		boldStyle.Render("РЕКОМЕНДАЦИЯ:"),
	)

	box := boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("⚠️  ВНИМАНИЕ"),
			content,
		),
	)

	hints := hintStyle.Render("[y/Enter] Я понимаю риск, продолжить  [n/Esc] Изменить папку")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, box, hints),
	)
}
