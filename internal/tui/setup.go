package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
)

type setupModel struct {
	form *huh.Form
	cfg  *config.Config
}

func newSetupModel(width, height int, cfg *config.Config) setupModel {
	extStr := strings.Join(cfg.Extensions, ", ")
	retentionStr := strconv.Itoa(cfg.RetentionDays)
	intervalStr := strconv.Itoa(cfg.CheckIntervalHours)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Папка для очистки").
				Description("Укажите путь к папке со скриншотами").
				Value(&cfg.TargetDir),

			huh.NewInput().
				Title("Расширения файлов").
				Description("Через запятую, например: .png, .mov, .mp4").
				Value(&extStr),

			huh.NewInput().
				Title("Срок хранения (дней)").
				Description("Файлы старше этого срока будут удалены").
				Value(&retentionStr),

			huh.NewInput().
				Title("Интервал проверки (часов)").
				Description("Как часто запускать очистку в фоне").
				Value(&intervalStr),

			huh.NewConfirm().
				Title("Тестовый режим (Dry Run)").
				Description("Показывать файлы, но не удалять их").
				Value(&cfg.DryRun),
		),
	)

	return setupModel{form: form, cfg: cfg}
}

func (m setupModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m model) updateSetup(msg tea.Msg) (tea.Model, tea.Cmd) {
	formModel, cmd := m.setupModel.form.Update(msg)
	if f, ok := formModel.(*huh.Form); ok {
		m.setupModel.form = f
	}

	if m.setupModel.form.State == huh.StateCompleted {
		// Parse extensions
		extStr := m.setupModel.form.GetString("Расширения файлов")
		if extStr == "" {
			extStr = ".png, .mov, .mp4, .gif"
		}
		m.cfg.Extensions = parseExtensions(extStr)

		// Parse retention
		retStr := m.setupModel.form.GetString("Срок хранения (дней)")
		if retStr == "" {
			retStr = "30"
		}
		m.cfg.RetentionDays = parseInt(retStr, 30)

		// Parse interval
		intStr := m.setupModel.form.GetString("Интервал проверки (часов)")
		if intStr == "" {
			intStr = "24"
		}
		m.cfg.CheckIntervalHours = parseInt(intStr, 24)

		if !m.cfg.WarningAck && needsWarning(m.cfg.TargetDir) {
			m.screen = screenWarning
			m.warningModel = newWarningModel(m.cfg.TargetDir, m.cfg.Extensions, m.cfg.RetentionDays)
			return m, nil
		}

		return m.saveAndGoToStatus()
	}

	return m, cmd
}

func (m model) setupView() string {
	return m.setupModel.form.View()
}

func parseExtensions(s string) []string {
	parts := strings.Split(s, ",")
	var res []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			res = append(res, strings.ToLower(p))
		}
	}
	return res
}

func parseInt(s string, def int) int {
	v, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil || v <= 0 {
		return def
	}
	return v
}
