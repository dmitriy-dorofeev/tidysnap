package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dirItem struct {
	name     string
	path     string
	isParent bool
}

type folderPickerModel struct {
	width  int
	height int
	cwd    string
	home   string
	items  []dirItem
	cursor int
}

type folderPickerLoadedMsg struct {
	items []dirItem
}

type folderPickerErrorMsg struct {
	err error
}

type folderSelectedMsg struct {
	path string
}

type folderPickerBackMsg struct{}

func newFolderPickerModel(width, height int, startDir string) folderPickerModel {
	home, _ := os.UserHomeDir()
	if startDir == "" {
		startDir = home
	}
	return folderPickerModel{
		width:  width,
		height: height,
		cwd:    startDir,
		home:   home,
		items:  nil,
		cursor: 0,
	}
}

func (m folderPickerModel) Init() tea.Cmd {
	return m.load()
}

func (m folderPickerModel) load() tea.Cmd {
	return func() tea.Msg {
		entries, err := os.ReadDir(m.cwd)
		if err != nil {
			return folderPickerErrorMsg{err}
		}

		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
		})

		var items []dirItem
		if m.cwd != m.home {
			items = append(items, dirItem{name: "..", path: filepath.Dir(m.cwd), isParent: true})
		}
		for _, e := range entries {
			if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
				items = append(items, dirItem{
					name: e.Name(),
					path: filepath.Join(m.cwd, e.Name()),
				})
			}
		}
		return folderPickerLoadedMsg{items}
	}
}

func (m model) updateFolderPicker(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case folderPickerLoadedMsg:
		m.folderPickerModel.items = msg.items
		m.folderPickerModel.cursor = 0
		return m, nil

	case folderPickerErrorMsg:
		m.err = msg.err
		return m, nil

	case folderSelectedMsg:
		m.cfg.TargetDir = msg.path
		m.screen = screenSetup
		m.setupModel = newSetupModel(m.width, m.height, m.cfg)
		return m, m.setupModel.Init()

	case folderPickerBackMsg:
		m.screen = screenWelcome
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			m.screen = screenWelcome
			return m, nil
		case "up", "k":
			if m.folderPickerModel.cursor > 0 {
				m.folderPickerModel.cursor--
			}
			return m, nil
		case "down", "j":
			if m.folderPickerModel.cursor < len(m.folderPickerModel.items)-1 {
				m.folderPickerModel.cursor++
			}
			return m, nil
		case "left", "h":
			parent := filepath.Dir(m.folderPickerModel.cwd)
			if parent != m.folderPickerModel.cwd && m.folderPickerModel.cwd != m.folderPickerModel.home {
				m.folderPickerModel.cwd = parent
				return m, m.folderPickerModel.load()
			}
			return m, nil
		case "right", "l", "enter":
			if len(m.folderPickerModel.items) > 0 {
				item := m.folderPickerModel.items[m.folderPickerModel.cursor]
				m.folderPickerModel.cwd = item.path
				return m, m.folderPickerModel.load()
			}
			return m, nil
		case " ":
			if len(m.folderPickerModel.items) > 0 {
				item := m.folderPickerModel.items[m.folderPickerModel.cursor]
				m.cfg.TargetDir = item.path
				m.screen = screenSetup
				m.setupModel = newSetupModel(m.width, m.height, m.cfg)
				return m, m.setupModel.Init()
			}
			return m, nil
		}
	}

	return m, nil
}

func (m model) folderPickerView() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("86")).MarginBottom(1)
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).MarginTop(1)
	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("250")).MarginBottom(1)
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
	itemStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	parentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	emptyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true)

	title := titleStyle.Render("📁 Выберите папку для очистки")
	currentPath := pathStyle.Render("Текущая папка: " + m.folderPickerModel.cwd)

	var lines []string
	if len(m.folderPickerModel.items) == 0 {
		lines = append(lines, emptyStyle.Render("(пустая папка)"))
	} else {
		for i, item := range m.folderPickerModel.items {
			prefix := "  "
			style := itemStyle
			if item.isParent {
				style = parentStyle
			}
			if i == m.folderPickerModel.cursor {
				prefix = cursorStyle.Render("> ")
				style = style.Copy().Bold(true)
			}
			suffix := "/"
			if item.isParent {
				suffix = ""
			}
			lines = append(lines, prefix+style.Render(item.name+suffix))
		}
	}

	// Ограничиваем высоту списка
	maxItems := m.folderPickerModel.height - 10
	if maxItems < 3 {
		maxItems = 3
	}
	if len(lines) > maxItems {
		lines = lines[:maxItems]
	}

	listView := strings.Join(lines, "\n")
	hints := hintStyle.Render("↑/↓ навигация • ← назад • →/Enter открыть • Пробел выбрать • q назад")

	content := lipgloss.JoinVertical(lipgloss.Left, currentPath, listView)
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Center, title, content, hints),
	)
}
