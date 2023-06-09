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
		m.spectrum.SetDataSource(dataSource)
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

	m.width = m.table.Width()
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
		ok bool
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	var chartFocused = func(focus bool) {
		if chart, ok := m.chart.(widgets.WithFocus); ok {
			chart.Focused(focus)
		}
	}

	var focusChart = func(chart tea.Model) func() tea.Cmd {
		chartFocused(false)
		if m.chart == chart && m.chart == m.spectrum {
			m.spectrum.NextBandView()
		}
		m.chart = chart
		chartFocused(true)

		w := m.width
		net, color := m.table.GetSelectedNetwork()

		return func() tea.Cmd {
			return tea.Batch(
				func() tea.Msg {
					return events.TableWidthMsg(w)
				},
				func() tea.Msg {
					return events.SelectedNetworkKeyMsg{
						Key:   net.Key(),
						Color: color,
					}
				},
			)
		}
	}

	{
		model, cmd := m.table.Update(msg)
		if m.table, ok = model.(*wifitable.Model); !ok {
			log.Fatalf("wifi table update method returned unexpected model %v", model)
		}
		cmds = append(cmds, cmd)
	}

	{
		model, cmd := m.sparkline.Update(msg)
		if m.sparkline, ok = model.(*sparkline.Model); !ok {
			log.Fatalf("sparkline update method returned unexpected model %v", model)
		}
		cmds = append(cmds, cmd)
	}

	{
		model, cmd := m.spectrum.Update(msg)
		if m.spectrum, ok = model.(*spectrum.Model); !ok {
			log.Fatalf("spectrum update method returned unexpected model %v", model)
		}
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {
	case events.TableWidthMsg:
		m.width = int(msg)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Sparkline):
			focusChart(m.sparkline)
			// chartFocused(false)
			// m.chart = m.sparkline
			// chartFocused(true)
			// cmds = append(cmds, onChartRefresh())

		case key.Matches(msg, m.keys.Spectrum):
			focusChart(m.spectrum)
			// chartFocused(false)
			// if m.chart == m.spectrum {
			// 	m.spectrum.NextBandView()
			// }
			// m.chart = m.spectrum
			// chartFocused(true)
			// cmds = append(cmds, onChartRefresh())

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
