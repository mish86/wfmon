package dashboard

import (
	"wfmon/pkg/ds"
	"wfmon/pkg/widgets/sparkline"
	"wfmon/pkg/widgets/wifitable"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	dataSource ds.Provider
	table      *wifitable.Model
	sparkline  *sparkline.Model
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

func New(opts ...Option) *Model {
	help := help.New()
	help.ShowAll = true

	m := &Model{
		table:     wifitable.New(),
		sparkline: sparkline.New(),
		help:      &help,
	}

	for _, opt := range opts {
		opt(m)
	}

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
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	tableKeys := m.table.Keys()
	helpKey, quitKey := tableKeys.Help, tableKeys.Quit

	// TODO check focused in table inside implementation
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	// TODO check focused in sparkline inside implementation
	m.sparkline, cmd = m.sparkline.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, helpKey):
			m.helpShown = !m.helpShown

		case key.Matches(msg, quitKey):
			cmds = append(cmds, tea.Quit)
		}
	}

	// Bubble up the cmds
	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	if m.helpShown {
		keys := m.table.Keys()
		return m.help.View(&keys)
	}

	return m.table.View() + "\n" + m.sparkline.View()
}
