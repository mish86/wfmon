package wifitable

import (
	"strconv"
	netdata "wfmon/pkg/data/net"
	log "wfmon/pkg/logger"
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/widgets/color"
	"wfmon/pkg/widgets/events"
	column "wfmon/pkg/widgets/wifitable/col"
	order "wfmon/pkg/widgets/wifitable/ord"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return refreshTick(defaultRefreshInterval)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// returns event with selected network key
	var onSelectedCmd = func() tea.Cmd {
		net := m.GetSelectedNetwork()

		// broadcast event to other widgets about changed selection
		k := net.Key()
		c := m.colors[k]
		return func() tea.Msg {
			// key := m.networks[cursor].Key()
			return events.NetworkKeyMsg{
				Key:   k,
				Color: c,
			}
		}
	}

	// returns event with new table width
	var onResizeCmd = func() tea.Cmd {
		w := m.Width()
		return func() tea.Msg {
			return events.TableWidthMsg(w)
		}
	}

	var onPageUpdate = func() tea.Cmd {
		from := cmp.Max(0, (m.CurrentPage()-1)*m.PageSize())
		to := cmp.Min(len(m.networks), m.CurrentPage()*m.PageSize())
		n := make([]netdata.Network, to-from)
		c := make([]color.HexColor, len(n))
		copy(n, m.networks[from:to])
		var found bool
		for i := range n {
			if c[i], found = m.colors[n[i].Key()]; !found {
				c[i] = color.Black()
			}
		}
		return func() tea.Msg {
			return events.NetworksOnScreen{Networks: n, Colors: c}
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

		return tea.Batch(onSelectedCmd(), onPageUpdate())
	}

	m.Model, cmd = m.Model.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.RowUp):
			m.moveRowUp()
			cmds = append(cmds, onSelectedCmd(), onPageUpdate())

		case key.Matches(msg, m.keys.RowDown):
			m.moveRowDown()
			cmds = append(cmds, onSelectedCmd(), onPageUpdate())

		case key.Matches(msg, m.keys.PageUp):
			m.Model = m.PageUp()
			cmds = append(cmds, onSelectedCmd(), onPageUpdate())

		case key.Matches(msg, m.keys.PageDown):
			m.Model = m.PageDown()
			cmds = append(cmds, onSelectedCmd(), onPageUpdate())

		case key.Matches(msg, m.keys.GotoTop):
			m.gotoTop()
			cmds = append(cmds, onSelectedCmd(), onPageUpdate())

		case key.Matches(msg, m.keys.GotoBottom):
			m.gotoBottom()
			cmds = append(cmds, onSelectedCmd(), onPageUpdate())

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
			cmds = append(cmds, onResizeCmd(), onSelectedCmd(), onPageUpdate())

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

		cmds = append(cmds, onSelectedCmd(), onPageUpdate(), refreshTick(defaultRefreshInterval))
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
