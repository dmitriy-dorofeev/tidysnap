package tui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
	"github.com/dmitriy-dorofeev/tidysnap/internal/i18n"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
	"github.com/dustin/go-humanize"
)

type daemonStatus struct {
	installed  bool
	running    bool
	loaded     bool
	nextRun    time.Time
	hasNextRun bool
}

type statusModel struct {
	width  int
	height int
	cfg    *config.Config
	msg    string
	status daemonStatus
}

func newStatusModel(width, height int, cfg *config.Config) statusModel {
	return statusModel{width: width, height: height, cfg: cfg, status: getDaemonStatus(cfg.CheckIntervalHours)}
}

func (sm statusModel) Init() tea.Cmd {
	if sm.status.running {
		return pollDaemonStatusCmd(sm.cfg.CheckIntervalHours)
	}
	return nil
}

func pollDaemonStatusCmd(intervalHours int) tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return daemonStatusMsg{status: getDaemonStatus(intervalHours)}
	})
}

func (m model) updateStatus(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case keyMatches(msg, 'r'):
			return m.runScan()
		case keyMatches(msg, 'l'):
			m.screen = screenLogView
			m.logViewModel = newLogViewModel(m.width, m.height)
			return m, m.logViewModel.Init()
		case keyMatches(msg, 'e'):
			m.screen = screenFolderPicker
			m.folderPickerModel = newFolderPickerModel(m.width, m.height, m.cfg.TargetDir, screenStatus)
			return m, m.folderPickerModel.Init()
		case keyMatches(msg, 's'):
			return m, tea.Batch(m.daemonActionCmd(), pollDaemonStatusCmd(m.cfg.CheckIntervalHours))
		case keyMatches(msg, 'x'):
			m.screen = screenResetConfirm
			m.resetModel = newResetModel()
			return m, nil
		case keyMatches(msg, 'q') || msg.String() == "esc":
			return m, tea.Quit
		}
	case scanDoneMsg:
		m.screen = screenPreview
		m.previewModel = newPreviewModel(m.width, m.height, msg.files, m.cfg.DryRun)
		return m, m.previewModel.Init()
	case cleanupDoneMsg:
		m.statusModel.msg = fmt.Sprintf("Очистка завершена: удалено %d файлов, освобождено %s",
			msg.stats.FilesRemoved, humanize.Bytes(safeUint64(msg.stats.BytesFreed)))
		m.screen = screenStatus
		return m, nil
	case daemonStatusMsg:
		m.statusModel.status = msg.status
		if msg.errMsg != "" {
			m.statusModel.msg = msg.errMsg
		}
		if m.statusModel.status.running {
			return m, pollDaemonStatusCmd(m.cfg.CheckIntervalHours)
		}
		return m, nil
	}
	return m, nil
}

func (m model) runScan() (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		s := scanner.New(m.cfg.Extensions, m.cfg.RetentionDays)
		files, err := s.Scan(context.Background(), m.cfg.TargetDir)
		if err != nil {
			return errMsg{err}
		}
		return scanDoneMsg{files}
	}
}

func (m model) daemonActionCmd() tea.Cmd {
	return func() tea.Msg {
		var errMsg string
		switch {
		case !daemon.IsInstalled():
			if err := daemon.Install(daemon.BinaryPath(), m.cfg.CheckIntervalHours); err != nil {
				errMsg = fmt.Sprintf("Ошибка установки: %v", err)
			}
		case daemon.IsRunning():
			if err := daemon.Stop(); err != nil {
				errMsg = fmt.Sprintf("Ошибка остановки: %v", err)
			}
		case daemon.IsLoaded():
			if err := daemon.Start(); err != nil {
				errMsg = fmt.Sprintf("Ошибка запуска: %v", err)
			}
		default:
			if err := daemon.Load(); err != nil {
				errMsg = fmt.Sprintf("Ошибка загрузки: %v", err)
			}
		}
		return daemonStatusMsg{status: getDaemonStatus(m.cfg.CheckIntervalHours), errMsg: errMsg}
	}
}

func getDaemonStatus(intervalHours int) daemonStatus {
	s := daemonStatus{
		installed: daemon.IsInstalled(),
		running:   daemon.IsRunning(),
		loaded:    daemon.IsLoaded(),
	}
	if nextRun, ok := daemon.NextRunTime(intervalHours); ok {
		s.nextRun = nextRun
		s.hasNextRun = true
	}
	return s
}

type daemonStatusMsg struct {
	status daemonStatus
	errMsg string
}

func (m model) statusView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginBottom(1)
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("86")).Padding(1, 2).Width(70)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	valueStyle := lipgloss.NewStyle().Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
	msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).MarginTop(1)

	status := m.statusModel.status

	var statusStr string
	switch {
	case !status.installed:
		statusStr = i18n.T("status_not_installed")
	case status.running:
		statusStr = i18n.T("status_running")
	case status.loaded:
		statusStr = i18n.T("status_loaded")
	default:
		statusStr = i18n.T("status_unloaded")
	}

	lines := []string{
		fmt.Sprintf("%s %s", labelStyle.Render(i18n.T("label_folder")), valueStyle.Render(m.cfg.TargetDir)),
		fmt.Sprintf("%s %s", labelStyle.Render(i18n.T("label_extensions")), valueStyle.Render(strings.Join(m.cfg.Extensions, ", "))),
		fmt.Sprintf("%s %s", labelStyle.Render(i18n.T("label_retention")), valueStyle.Render(fmt.Sprintf(i18n.T("label_retention_days"), m.cfg.RetentionDays))),
		fmt.Sprintf("%s %s", labelStyle.Render(i18n.T("label_dryrun")), valueStyle.Render(map[bool]string{true: i18n.T("yes"), false: i18n.T("no")}[m.cfg.DryRun])),
		fmt.Sprintf("%s %s", labelStyle.Render(i18n.T("label_daemon")), valueStyle.Render(statusStr)),
	}

	if status.hasNextRun {
		lines = append(lines, fmt.Sprintf("%s %s", labelStyle.Render(i18n.T("label_next_run")), valueStyle.Render(status.nextRun.Format("02.01.2006 15:04"))))
	}

	if m.statusModel.msg != "" {
		lines = append(lines, "", msgStyle.Render(m.statusModel.msg))
	}

	box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	var daemonHint string
	switch {
	case !status.installed:
		daemonHint = i18n.T("hint_install")
	case status.running:
		daemonHint = i18n.T("hint_stop")
	case status.loaded:
		daemonHint = i18n.T("hint_start")
	default:
		daemonHint = i18n.T("hint_load")
	}
	hints := hintStyle.Render(fmt.Sprintf(i18n.T("status_hints"), daemonHint))

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, titleStyle.Render(i18n.T("status_title")), box, hints),
	)
}

func safeUint64(n int64) uint64 {
	if n < 0 {
		return 0
	}
	return uint64(n)
}
