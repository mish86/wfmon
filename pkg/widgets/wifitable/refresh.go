package wifitable

import (
	"time"
	"wfmon/pkg/widgets/color"
	column "wfmon/pkg/widgets/wifitable/col"
	"wfmon/pkg/widgets/wifitable/row"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/evertras/bubble-table/table"
)

// Event to refresh and redraw table.
type refreshMsg time.Time

// Invokes refresh table by interval.
// Fresh data obtained on timer end.
func refreshTick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return refreshMsg(t)
	})
}

// Immediately reapplies data and columns.
func (m *Model) refresh() {
	m.Model = m.
		WithRows(m.getRows()).
		WithColumns(column.Converter(m.columns)(m.sort))
}

// Returns table rows from networks.
// Networks already sorted in @onRefreshMsg.
func (m *Model) getRows() []table.Row {
	viewer := row.Converter(m.columns, cellViewers())

	rows := make([]table.Row, len(m.networks))
	for rowID, e := range m.networks {
		entry := e

		rowStyle := defaultBaseStyle
		if entry.Key().Compare(m.associated) == 0 {
			rowStyle = defaultAssociatedStyle
		}

		data := row.Data{Network: entry}.
			HashColor(m.colors[entry.Key()].Lipgloss()).
			Style(rowStyle)

		rows[rowID] = viewer(&data)
	}

	return rows
}

// Handles refresh tick.
// Fetches networks from data source.
// Sorts networks as per current column and order.
// Invokes @refresh to redraw the table.
func (m *Model) onRefreshMsg(_ refreshMsg) {
	selectedNetwork, _ := m.GetSelectedNetwork()

	// FIXME: race at access to networks and colors in update.
	m.networks = m.dataSource.Networks()

	iter := color.Random()
	// preserve row colors
	for _, network := range m.networks {
		key := network.Key()
		if _, found := m.colors[key]; !found {
			m.colors[key], iter = iter()
		}
	}

	// apply current sorting
	m.sort.Sort(m.networks)

	// preserve selected row
	currendRowID := 0
	for rowID, network := range m.networks {
		key := network.Key()
		if key.Compare(selectedNetwork.Key()) == 0 {
			currendRowID = rowID
		}
	}

	// apply columns and rows to table
	m.refresh()

	// preserve selected row
	m.Model = m.WithHighlightedRow(currendRowID)
}
