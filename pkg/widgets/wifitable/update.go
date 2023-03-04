package wifitable

import (
	"strconv"
	log "wfmon/pkg/logger"
	"wfmon/pkg/widgets"
	column "wfmon/pkg/widgets/wifitable/col"
	order "wfmon/pkg/widgets/wifitable/ord"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return refreshTick(defaultRefreshInterval)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// returns event with selected network key
	var onSelectedCmd = func() tea.Cmd {
		cursor := m.GetHighlightedRowIndex()

		// no data
		if len(m.networks) == 0 {
			return nil
		}

		// out of bounds
		if cursor < 0 || cursor >= len(m.networks) {
			log.Errorf("cursor %d out of bounds networks: %v", cursor, m.networks)
			return nil
		}

		// cursor not changed
		if m.selected.Compare(m.networks[cursor].Key()) == 0 {
			return nil
		}

		// broadcast event to other widgets about changed selection
		k := m.networks[cursor].Key()
		c := m.colors[k]
		return func() tea.Msg {
			// key := m.networks[cursor].Key()
			return widgets.NetworkKeyMsg{
				Key:   k,
				Color: c,
			}
		}
	}

	// returns event with new table width
	var onResizeCmd = func() tea.Cmd {
		w := m.tableWidth()
		return func() tea.Msg {
			return widgets.TableWidthMsg(w)
		}
	}

	// Rotates column in @Multiple column view.
	// Refresh table and send resize and select events.
	var onCycleColumn = func(colIdx int) tea.Cmd {
		col := m.columns[colIdx]
		prevKey := col.Key()

		if c, ok := col.(column.Cycler); !ok {
			return nil
		} else {
			col = c.Next()
		}

		if m.sort.Key() == prevKey {
			m.sort = sortBy(col.Key())(m.sort.Order())
		}

		m.columns[colIdx] = col

		// apply current sorting and refresh table
		m.sort.Sort(m.networks)
		m.refresh()

		// send event about table width and cursort change
		return tea.Batch(onResizeCmd(), onSelectedCmd())
		// return onResizeCmd()
	}

	// Sorts table by column index.
	// Numbering starts from SSID column.
	var onSortColumn = func(msg tea.KeyMsg) tea.Cmd {
		var num, idx int
		var err error
		if num, err = strconv.Atoi(msg.String()); err != nil {
			log.Warnf("failed to sort, %w", err)
			return nil
		}

		// Column number starts from 1
		// Hash column is not registered for sorting
		idx = num

		keys := visibleColumnKeys(m.columns)
		if idx < 0 || idx >= len(keys) {
			log.Warnf("unsupported sort key, %d", num)
			return nil
		}

		key := keys[idx]
		// swap order for current column
		if m.sort.Key() == key {
			m.sort = m.sort.SwapOrder()
		} else {
			// ASC order for new column
			m.sort = sortBy(key)(order.ASC)
		}

		// apply current sorting and refresh table
		m.sort.Sort(m.networks)
		m.refresh()

		return onSelectedCmd()
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.RowUp):
			m.moveRowUp()
			cmds = append(cmds, onSelectedCmd())

		case key.Matches(msg, m.keys.RowDown):
			m.moveRowDown()
			cmds = append(cmds, onSelectedCmd())

		case key.Matches(msg, m.keys.PageUp):
			m.Model = m.PageUp()
			cmds = append(cmds, onSelectedCmd())

		case key.Matches(msg, m.keys.PageDown):
			m.Model = m.PageDown()
			cmds = append(cmds, onSelectedCmd())

		case key.Matches(msg, m.keys.GotoTop):
			m.gotoTop()
			cmds = append(cmds, onSelectedCmd())

		case key.Matches(msg, m.keys.GotoBottom):
			m.gotoBottom()
			cmds = append(cmds, onSelectedCmd())

		case key.Matches(msg, m.keys.SignalView):
			cmds = append(cmds, onCycleColumn(SignalMColumnIdx))

		case key.Matches(msg, m.keys.StationView):
			cmds = append(cmds, onCycleColumn(StationMColumnIdx))

		case key.Matches(msg, m.keys.Reset):
			// reset columns view
			m.columns = columns()
			// reset sorting
			m.sort = defaultSort()
			// apply current sorting
			m.sort.Sort(m.networks)
			// refresh table
			m.refresh()
			// send event about table width and cursort change
			cmds = append(cmds, onResizeCmd(), onSelectedCmd())

		case key.Matches(msg, m.keys.Sort):
			cmds = append(cmds, onSortColumn(msg))
		}

	case refreshMsg:
		broadcastFirstEvent := len(m.networks) == 0

		// fetch fresh data from data source and apply it to the table.
		m.onRefreshMsg(msg)

		broadcastFirstEvent = broadcastFirstEvent && len(m.networks) > 0

		if broadcastFirstEvent {
			cmds = append(cmds, onResizeCmd())
		}

		cmds = append(cmds, onSelectedCmd(), refreshTick(defaultRefreshInterval))
	}

	// Bubble up the cmds
	return m, tea.Batch(cmds...)
}

func (m *Model) moveRowUp() {
	rowIdx := m.GetHighlightedRowIndex() - 1

	if rowIdx < 0 {
		rowIdx = 0
	}

	m.Model = m.WithHighlightedRow(rowIdx)
}

func (m *Model) moveRowDown() {
	rowIdx := m.GetHighlightedRowIndex() + 1

	if rowIdx >= len(m.GetVisibleRows()) {
		rowIdx = len(m.GetVisibleRows()) - 1
	}

	m.Model = m.WithHighlightedRow(rowIdx)
}

func (m *Model) gotoTop() {
	m.Model = m.
		WithCurrentPage(0).
		WithHighlightedRow(0)
}

func (m *Model) gotoBottom() {
	m.Model = m.
		WithCurrentPage(m.MaxPages() - 1).
		WithHighlightedRow(len(m.GetVisibleRows()) - 1)
}
