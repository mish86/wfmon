package wifitable

import (
	"time"

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
		WithColumns(GenerateColumns(m.sort, m.stationViewMode.Current(), m.signalViewMode.Current()))
}

// Returns table rows from networks.
// Networks already sorted in @onRefreshMsg.
func (m *Model) getRows() []table.Row {
	// get registered column keys
	columns := ColumnsKeys()
	// get registered cell viewers
	viewers := GenerateCellViewers(m.associated)

	rows := make([]table.Row, len(m.networks))
	for rowID, e := range m.networks {
		entry := e

		row := make(table.RowData, len(columns))
		for _, key := range columns {
			row[key] = viewers[key](&entry)
		}

		rows[rowID] = table.NewRow(row)
	}

	return rows
}

// Handles refresh tick.
// Fetches networks from data source.
// Sorts networks as per current column and order.
// Invokes @refresh to redraw the table.
func (m *Model) onRefreshMsg(msg refreshMsg) {
	m.networks = m.dataSource.Networks()

	// apply current sorting
	m.sort.Sort(m.networks)

	// apply columns and rows to table
	m.refresh()
}
