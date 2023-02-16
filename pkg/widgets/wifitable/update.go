package wifitable

import (
	"strconv"
	log "wfmon/pkg/logger"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	// refreshCmd := func() tea.Msg { return RefreshMsg(m.networks) }

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.RowUp):
			m.moveRowUp()

		case key.Matches(msg, m.keys.RowDown):
			m.moveRowDown()

		case key.Matches(msg, m.keys.GotoTop):
			m.gotoTop()

		case key.Matches(msg, m.keys.GotoBottom):
			m.gotoBottom()

		case key.Matches(msg, m.keys.SignalView):
			prevKey := m.signalViewMode.Current().Key()
			m.signalViewMode = m.signalViewMode.Next()
			if m.sort.key == prevKey {
				m.sort = SortBy(m.signalViewMode.Current().Key())(m.sort.ord)
			}
			m.refresh()

		case key.Matches(msg, m.keys.StationView):
			prevKey := m.stationViewMode.Current().Key()
			m.stationViewMode = m.stationViewMode.Next()
			if m.sort.key == prevKey {
				m.sort = SortBy(m.stationViewMode.Current().Key())(m.sort.ord)
			}
			m.refresh()

		case key.Matches(msg, m.keys.Reset):
			m.sort, m.signalViewMode, m.stationViewMode = defaultViewMode()
			m.refresh()

		case key.Matches(msg, m.keys.Sort):
			if m.sortBy(msg) {
				m.refresh()
			}

		case key.Matches(msg, m.keys.Help):
			m.helpShown = !m.helpShown

		case key.Matches(msg, m.keys.Quit):
			cmds = append(cmds, tea.Quit)
		}

	case RefreshMsg:
		// Apply rows and columns in table
		m.onRefreshMsg(msg)

		// schedule next refresh tick
		cmds = append(cmds, refreshTick(m.dataSource, defaultRefreshInterval))
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) moveRowUp() {
	rowIdx := m.table.GetHighlightedRowIndex() - 1

	if rowIdx < 0 {
		rowIdx = 0
	}

	m.table = m.table.WithHighlightedRow(rowIdx)
}

func (m *Model) moveRowDown() {
	rowIdx := m.table.GetHighlightedRowIndex() + 1

	if rowIdx >= len(m.table.GetVisibleRows()) {
		rowIdx = len(m.table.GetVisibleRows()) - 1
	}

	m.table = m.table.WithHighlightedRow(rowIdx)
}

func (m *Model) gotoTop() {
	m.table = m.table.
		WithCurrentPage(0).
		WithHighlightedRow(0)
}

func (m *Model) gotoBottom() {
	m.table = m.table.
		WithCurrentPage(m.table.MaxPages() - 1).
		WithHighlightedRow(len(m.table.GetVisibleRows()) - 1)
}

func (m *Model) sortBy(msg tea.KeyMsg) bool {
	var num, idx int
	var err error
	if num, err = strconv.Atoi(msg.String()); err != nil {
		log.Warnf("failed to sort, %w", err)
		return false
	}

	idx = num - 1
	keys := VisibleColumnsKeys(
		m.stationViewMode.Current(),
		m.signalViewMode.Current(),
	)
	if idx < 0 || idx >= len(keys) {
		log.Warnf("unsupported sort key, %d", num)
		return false
	}

	key := keys[idx]
	// swap order for current column
	if m.sort.key == key {
		m.sort.SwapOrder()
	} else {
		// ASC order for new column
		m.sort = SortBy(key)(ASC)
	}

	// Immediately apply current sorting for networks
	m.sort.Sort(m.networks)

	return true
}
