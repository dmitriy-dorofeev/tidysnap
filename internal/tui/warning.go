package tui

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
	"github.com/dmitriy-dorofeev/tidysnap/internal/i18n"
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
		switch {
		case msg.String() == "enter" || keyMatches(msg, 'y'):
			m.cfg.WarningAck = true
			if err := config.Save(m.cfg); err != nil {
				m.err = err
				return m, nil
			}
			return m.saveAndGoToStatus()
		case keyMatches(msg, 'n') || msg.String() == "esc":
			m.screen = screenSetup
			m.setupModel = newSetupModel(m.width, m.height, m.cfg)
			return m, m.setupModel.Init()
		}
	}
	return m, nil
}

func (m model) saveAndGoToStatus() (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		if err := config.Save(m.cfg); err != nil {
			return saveAndGoToStatusMsg{err: err}
		}
		logDir, err := config.LogDir()
		if err != nil {
			return saveAndGoToStatusMsg{err: err}
		}
		if err := os.MkdirAll(logDir, 0750); err != nil {
			return saveAndGoToStatusMsg{err: err}
		}
		binary := daemon.BinaryPath()
		if err := daemon.Install(binary, m.cfg.CheckIntervalHours); err != nil {
			return saveAndGoToStatusMsg{err: err}
		}
		return saveAndGoToStatusMsg{}
	}
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
		i18n.T("warning_folder")+"\n\n"+
			i18n.T("warning_text1")+"\n"+
			"%s\n"+
			i18n.T("warning_text2")+"\n\n"+
			i18n.T("warning_text3")+"\n\n"+
			"%s\n"+
			i18n.T("warning_text4")+"\n"+
			i18n.T("warning_text5"),
		boldStyle.Render(exts),
		m.cfg.RetentionDays,
		boldStyle.Render(i18n.T("warning_recommend")),
	)

	box := boxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render(i18n.T("warning_title")),
			content,
		),
	)

	hints := hintStyle.Render(i18n.T("warning_hints"))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, box, hints),
	)
}
