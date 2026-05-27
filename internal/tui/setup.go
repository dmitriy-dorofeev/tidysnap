package tui

import (
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/dmitriy-dorofeev/tidysnap/internal/config"
	"github.com/dmitriy-dorofeev/tidysnap/internal/i18n"
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
				Key("extensions").
				Title(i18n.T("setup_extensions")).
				Description(i18n.T("setup_extensions_desc")).
				Value(&extStr),

			huh.NewInput().
				Key("retention").
				Title(i18n.T("setup_retention")).
				Description(i18n.T("setup_retention_desc")).
				Value(&retentionStr),

			huh.NewInput().
				Key("interval").
				Title(i18n.T("setup_interval")).
				Description(i18n.T("setup_interval_desc")).
				Value(&intervalStr),

			huh.NewConfirm().
				Key("dryrun").
				Title(i18n.T("setup_dryrun")).
				Description(i18n.T("setup_dryrun_desc")).
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
		extStr := m.setupModel.form.GetString("extensions")
		if extStr == "" {
			extStr = ".png, .jpg, .jpeg, .mov, .mp4, .gif"
		}
		m.cfg.Extensions = parseExtensions(extStr)

		// Parse retention
		retStr := m.setupModel.form.GetString("retention")
		if retStr == "" {
			retStr = "30"
		}
		m.cfg.RetentionDays = parseInt(retStr, 30)

		// Parse interval
		intStr := m.setupModel.form.GetString("interval")
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
