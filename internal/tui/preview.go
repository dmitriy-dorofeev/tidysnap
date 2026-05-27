package tui

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dmitriy-dorofeev/tidysnap/internal/cleaner"
	"github.com/dmitriy-dorofeev/tidysnap/internal/scanner"
	"github.com/dustin/go-humanize"
)

type previewItem struct {
	file scanner.ScanResult
}

func (i previewItem) Title() string { return i.file.Path }
func (i previewItem) Description() string {
	size := uint64(i.file.Size)
	if i.file.Size < 0 {
		size = 0
	}
	return fmt.Sprintf("%s | %s", i.file.ModTime.Format("2006-01-02"), humanize.Bytes(size))
}
func (i previewItem) FilterValue() string { return i.file.Path }

type previewModel struct {
	list   list.Model
	files  []scanner.ScanResult
	dryRun bool
}

func newPreviewModel(width, height int, files []scanner.ScanResult, dryRun bool) previewModel {
	items := make([]list.Item, 0, len(files))
	for _, f := range files {
		items = append(items, previewItem{file: f})
	}

	l := list.New(items, list.NewDefaultDelegate(), width-4, height-8)
	l.Title = "Файлы на удаление"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.KeyMap.Quit = key.NewBinding(key.WithKeys("esc", "q"))

	return previewModel{list: l, files: files, dryRun: dryRun}
}

func (m previewModel) Init() tea.Cmd { return nil }

func (m model) updatePreview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "d":
			return m.runCleanup()
		case "esc", "q":
			m.screen = screenStatus
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.previewModel.list, cmd = m.previewModel.list.Update(msg)
	return m, cmd
}

func (m model) runCleanup() (tea.Model, tea.Cmd) {
	return m, func() tea.Msg {
		logFile, err := os.OpenFile(m.cfg.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return errMsg{err}
		}
		defer logFile.Close()

		logger := log.New(logFile, "", log.LstdFlags)
		c := cleaner.New(m.cfg.DryRun, logger)
		stats, err := c.Clean(context.Background(), m.previewModel.files)
		if err != nil {
			return errMsg{err}
		}
		return cleanupDoneMsg{stats}
	}
}

func (m model) previewView() string {
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
	action := "Удалить"
	if m.cfg.DryRun {
		action = "Показать (Dry Run)"
	}
	hints := hintStyle.Render(fmt.Sprintf("[d/Enter] %s  [Esc/q] Назад", action))

	if len(m.previewModel.files) == 0 {
		emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Center, emptyStyle.Render("Нет файлов для удаления"), hints))
	}

	return lipgloss.JoinVertical(lipgloss.Left, m.previewModel.list.View(), hints)
}
