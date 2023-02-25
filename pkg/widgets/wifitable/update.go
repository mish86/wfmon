package wifitable

import (
	"strconv"
	log "wfmon/pkg/logger"
	"wfmon/pkg/widgets"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	// return refreshTick(m.dataSource, defaultRefreshInterval)
	return refreshTick(defaultRefreshInterval)
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	// returns event with selected network key
	var onSelectedCmd = func() func() tea.Msg {
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
		return func() tea.Msg {
			return widgets.NetworkKeyMsg(m.networks[cursor].Key())
		}
	}

	// returns event with new table width
	var onResizeCmd = func() func() tea.Msg {
		enums := []Keyer{
			m.signalViewMode.Current(),
			m.stationViewMode.Current(),
		}
		return func() tea.Msg {
			return widgets.TableWidthMsg(tableWidth(enums...))
		}
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
			prevKey := m.signalViewMode.Current().Key()
			m.signalViewMode = m.signalViewMode.Next()
			if m.sort.key == prevKey {
				m.sort = SortBy(m.signalViewMode.Current().Key())(m.sort.ord)
			}
			// apply current sorting
			m.sort.Sort(m.networks)
			// refresh table
			m.refresh()
			// send event about table width and cursort change
			cmds = append(cmds, onResizeCmd(), onSelectedCmd())

		case key.Matches(msg, m.keys.StationView):
			prevKey := m.stationViewMode.Current().Key()
			m.stationViewMode = m.stationViewMode.Next()
			if m.sort.key == prevKey {
				m.sort = SortBy(m.stationViewMode.Current().Key())(m.sort.ord)
			}
			// apply current sorting
			m.sort.Sort(m.networks)
			// refresh table
			m.refresh()
			// send event about table width and cursort change
			cmds = append(cmds, onResizeCmd(), onSelectedCmd())

		case key.Matches(msg, m.keys.Reset):
			m.sort, m.signalViewMode, m.stationViewMode = defaultViewMode()
			// apply current sorting
			m.sort.Sort(m.networks)
			// refresh table
			m.refresh()
			// send event about table width and cursort change
			cmds = append(cmds, onResizeCmd(), onSelectedCmd())

		case key.Matches(msg, m.keys.Sort):
			if m.sortBy(msg) {
				m.refresh()
				cmds = append(cmds, onSelectedCmd())
			}
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
