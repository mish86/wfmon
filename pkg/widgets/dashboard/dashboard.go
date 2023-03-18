package dashboard

import (
	"strings"
	"wfmon/pkg/ds"
	log "wfmon/pkg/logger"
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/widgets"
	"wfmon/pkg/widgets/events"
	"wfmon/pkg/widgets/sparkline"
	"wfmon/pkg/widgets/spectrum"
	"wfmon/pkg/widgets/wifitable"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()
)

type Model struct {
	dataSource ds.Provider
	width      int
	table      *wifitable.Model
	sparkline  *sparkline.Model
	spectrum   *spectrum.Model
	chart      tea.Model
	keys       KeyMap
	help       *help.Model
	helpShown  bool
}

type Option func(*Model)

func WithDataSource(dataSource ds.Provider) Option {
	return func(m *Model) {
		m.dataSource = dataSource

		m.table.SetDataSource(dataSource)
		m.sparkline.SetDataSource(dataSource)
	}
}

func WithTable(t *wifitable.Model) Option {
	return func(m *Model) {
		m.table = t
	}
}

func WithSparkline(sl *sparkline.Model) Option {
	return func(m *Model) {
		m.sparkline = sl
	}
}

func WithSpectrum(s *spectrum.Model) Option {
	return func(m *Model) {
		m.spectrum = s
	}
}

func New(opts ...Option) *Model {
	help := help.New()
	help.ShowAll = true

	m := &Model{
		table:     wifitable.New(),
		sparkline: sparkline.New(),
		spectrum:  spectrum.New(),
		help:      &help,
		keys:      NewKeyMap(),
	}

	for _, opt := range opts {
		opt(m)
	}

	m.width = m.table.TableWidth()
	m.chart = m.sparkline

	return m
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		m.table.Init(),
		m.sparkline.Init(),
	)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		model tea.Model
		ok    bool
		cmd   tea.Cmd
		cmds  []tea.Cmd
	)

	model, cmd = m.table.Update(msg)
	if m.table, ok = model.(*wifitable.Model); !ok {
		log.Fatalf("wifi table update method returned unexpected model %v", model)
	}
	cmds = append(cmds, cmd)

	model, cmd = m.sparkline.Update(msg)
	if m.sparkline, ok = model.(*sparkline.Model); !ok {
		log.Fatalf("sparkline update method returned unexpected model %v", model)
	}
	cmds = append(cmds, cmd)

	model, cmd = m.spectrum.Update(msg)
	if m.spectrum, ok = model.(*spectrum.Model); !ok {
		log.Fatalf("spectrum update method returned unexpected model %v", model)
	}
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case events.TableWidthMsg:
		m.width = int(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Sparkline):
			m.chart = m.sparkline

		case key.Matches(msg, m.keys.Spectrum):
			if m.chart == m.spectrum {
				m.spectrum.NextBandView()
			}
			m.chart = m.spectrum

		case key.Matches(msg, m.keys.Help):
			m.helpShown = !m.helpShown

		case key.Matches(msg, m.keys.Quit):
			cmds = append(cmds, tea.Quit)
		}
	}

	// Bubble up the cmds
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.helpShown {
		return m.help.View(&m.keys)
	}

	return m.table.View() + "\n" + m.viewChartTitle() + "\n" + m.chart.View()
}

func (m *Model) viewChartTitle() string {
	var title string
	if t, ok := m.chart.(widgets.WithTitle); ok {
		title = t.Title()
	}
	title = titleStyle.Render(title)
	gaps := strings.Repeat("─", cmp.Max(0, (m.width-lipgloss.Width(title)))/2)
	return lipgloss.JoinHorizontal(lipgloss.Center, gaps, title, gaps)
}
