package sparkline

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/widgets"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return refreshTick(defaultRefreshInterval)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case widgets.NetworkKeyMsg:
		m.SetNetworkKey(netdata.Key(msg.Key))
		m.SetColor(msg.Color.Lipgloss())
		m.sparkline.Data = m.getData()
		m.refresh()

	case widgets.FieldMsg:
		m.SetFieldKey(string(msg))
		m.sparkline.Data = m.getData()
		m.refresh()

	case widgets.TableWidthMsg:
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
