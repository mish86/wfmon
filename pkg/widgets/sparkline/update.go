package sparkline

import (
	"wfmon/pkg/widgets/events"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return refreshTick(defaultRefreshInterval)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case events.NetworkKeyMsg:
		m.SetNetworkKey(msg.Key)
		m.SetColor(msg.Color.Lipgloss())
		m.data = m.getData()
		m.refresh()

	case events.FieldMsg:
		m.SetFieldKey(string(msg))
		m.data = m.getData()
		m.refresh()

	case events.TableWidthMsg:
		m.SetWidth(int(msg))
		m.refresh()

	case refreshMsg:
		// Apply refresh data to viewport
		m.onRefreshMsg(msg)

		// schedule next refresh tick
		cmds = append(cmds, refreshTick(defaultRefreshInterval))
	}

	// Bubble up the cmds
	return m, tea.Batch(cmds...)
}
