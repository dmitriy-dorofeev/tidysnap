package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/daemon"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
)

type statusModel struct {
	width  int
	height int
	cfg    *config.Config
	msg    string
}

func newStatusModel(width, height int, cfg *config.Config) statusModel {
	return statusModel{width: width, height: height, cfg: cfg}
}

func (m model) updateStatus(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			return m.runScan()
		case "l":
			m.screen = screenLogView
			m.logViewModel = newLogViewModel(m.width, m.height)
			return m, m.logViewModel.Init()
		case "e":
			m.screen = screenSetup
			m.setupModel = newSetupModel(m.width, m.height, m.cfg)
			return m, m.setupModel.Init()
		case "s":
			if daemon.IsRunning() {
				if err := daemon.Stop(); err != nil {
					m.statusModel.msg = fmt.Sprintf("Ошибка остановки: %v", err)
				}
			} else if daemon.IsInstalled() {
				if err := daemon.Start(); err != nil {
					m.statusModel.msg = fmt.Sprintf("Ошибка запуска: %v", err)
				}
			} else {
				if err := daemon.Install(daemon.BinaryPath(), m.cfg.CheckIntervalHours); err != nil {
					m.statusModel.msg = fmt.Sprintf("Ошибка установки: %v", err)
				}
			}
			return m, nil
		case "x":
			m.screen = screenResetConfirm
			m.resetModel = newResetModel()
			return m, nil
		case "q", "esc":
			return m, tea.Quit
		}
	case scanDoneMsg:
		m.screen = screenPreview
		m.previewModel = newPreviewModel(m.width, m.height, msg.files, m.cfg.DryRun)
		return m, m.previewModel.Init()
	case cleanupDoneMsg:
		m.statusModel.msg = fmt.Sprintf("Очистка завершена: удалено %d файлов, освобождено %s",
			msg.stats.FilesRemoved, humanizeBytes(msg.stats.BytesFreed))
		m.screen = screenStatus
		return m, nil
	}
	return m, nil
}

func (m model) runScan() (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		s := scanner.New(m.cfg.Extensions, m.cfg.RetentionDays)
		files, err := s.Scan(m.cfg.TargetDir)
		if err != nil {
			return errMsg{err}
		}
		return scanDoneMsg{files}
	}
}

func (m model) statusView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginBottom(1)
	boxStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("86")).Padding(1, 2).Width(70)
	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	valueStyle := lipgloss.NewStyle().Bold(true)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
	msgStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("220")).MarginTop(1)

	status := "❌ Не установлен"
	if daemon.IsInstalled() {
		if daemon.IsRunning() {
			status = "✅ Активен"
		} else {
			status = "⚠️  Установлен, но не запущен"
		}
	}

	lines := []string{
		fmt.Sprintf("%s %s", labelStyle.Render("Папка:"), valueStyle.Render(m.cfg.TargetDir)),
		fmt.Sprintf("%s %s", labelStyle.Render("Расширения:"), valueStyle.Render(strings.Join(m.cfg.Extensions, ", "))),
		fmt.Sprintf("%s %d дней", labelStyle.Render("Срок хранения:"), m.cfg.RetentionDays),
		fmt.Sprintf("%s %s", labelStyle.Render("Тестовый режим:"), valueStyle.Render(map[bool]string{true: "Да", false: "Нет"}[m.cfg.DryRun])),
		fmt.Sprintf("%s %s", labelStyle.Render("Демон:"), valueStyle.Render(status)),
	}

	if m.statusModel.msg != "" {
		lines = append(lines, "", msgStyle.Render(m.statusModel.msg))
	}

	box := boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	var daemonHint string
	if daemon.IsRunning() {
		daemonHint = "[s] Остановить демон"
	} else if daemon.IsInstalled() {
		daemonHint = "[s] Запустить демон"
	} else {
		daemonHint = "[s] Установить демон"
	}
	hints := hintStyle.Render("[r] Запустить очистку  [l] Логи  [e] Настройки  " + daemonHint + "  [x] Удалить  [q] Выход")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, titleStyle.Render("📊 TidySnap — Статус"), box, hints),
	)
}

func humanizeBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
