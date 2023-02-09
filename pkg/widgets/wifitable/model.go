package wifitable

import (
	"strconv"
	"strings"
	"time"
	log "wfmon/pkg/logger"
	"wfmon/pkg/network"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	defaultRefreshInterval = time.Second
	defaultTableHeight     = 10
)

var (
	defaultBaseStyle       = lipgloss.NewStyle()
	defaultHeaderStyle     = lipgloss.NewStyle().Foreground(lipgloss.NoColor{}).Bold(true)
	defaultSelectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
	defaultAssociatedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")).Bold(true)
)

type Model struct {
	table           table.Model
	sort            Sort
	signalViewMode  SignalViewMode
	dataSource      *DataSource
	networks        NetworkSlice
	associated      *NetworkKey
	refreshInterval time.Duration
	keys            KeyMap
	help            help.Model
	helpShown       bool
}

func NewTable(ds *DataSource, network network.Network) *Model {
	help := help.New()
	help.ShowAll = true

	signalViewMode := BarsViewMode
	sort := SortBy(ColumnBarsKey)(DESC)
	columns := GenerateColumns(sort, uint8(signalViewMode))

	keys := NewKeyMap()
	t := table.New(columns).
		Border(table.Border{}).
		WithPageSize(defaultTableHeight).
		WithKeyMap(keys.KeyMap).
		WithBaseStyle(defaultBaseStyle).
		HeaderStyle(defaultHeaderStyle).
		HighlightStyle(defaultSelectedStyle).
		Focused(true)

	return &Model{
		table:          t,
		keys:           keys,
		help:           help,
		signalViewMode: signalViewMode,
		sort:           sort,
		dataSource:     ds,
		networks:       NetworkSlice{},
		associated:     NewNetworkKey(network.BSSID, network.SSID),
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		refreshTick(m.dataSource, m.refreshInterval),
	)
}

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
		case key.Matches(msg, m.keys.SignalView):
			m.signalViewMode = m.signalViewMode.Next()
			m.sort = SortBy(m.signalViewMode.Key())(m.sort.ord)
			m.refresh()
			// cmds = append(cmds, refreshCmd)

		case key.Matches(msg, m.keys.ResetSort):
			m.sort = SortBy(ColumnSSIDKey)(None)
			m.refresh()
			// cmds = append(cmds, refreshCmd)

		case key.Matches(msg, m.keys.Sort):
			if m.sortBy(msg) {
				m.refresh()
				// cmds = append(cmds, refreshCmd)
			}

		case key.Matches(msg, m.keys.Help):
			m.helpShown = !m.helpShown

		case key.Matches(msg, m.keys.Quit):
			cmds = append(cmds, tea.Quit)
		}

	case RefreshMsg:
		// get networks from last tick
		m.networks = NetworkSlice(msg)

		// apply current sorting
		m.sort.Sort(m.networks)

		// apply columns and rows to table
		m.refresh()

		// schedule next refresh tick
		cmds = append(cmds, refreshTick(m.dataSource, m.refreshInterval))
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) View() string {
	body := strings.Builder{}

	if m.helpShown {
		body.WriteString(m.help.View(&m.keys))
	} else {
		// style := lipgloss.NewStyle()
		// body.WriteString(style.Render(m.table.View()))
		body.WriteString(m.table.View())
	}

	return body.String()
}

// Reapplies data and columns.
func (m *Model) refresh() {
	m.table = m.table.
		WithRows(getRows(m.networks, m.associated)).
		WithColumns(GenerateColumns(m.sort, uint8(m.signalViewMode)))
}

func (m *Model) sortBy(msg tea.KeyMsg) bool {
	var num, idx int
	var err error
	if num, err = strconv.Atoi(msg.String()); err != nil {
		log.Warnf("failed to sort, %w", err)
		return false
	}

	idx = num - 1
	keys := VisibleColumnsKeys(m.signalViewMode)
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

	// Immediately apply sorting for networks
	m.sort.Sort(m.networks)

	return true
}

// Returns sorted rows to redraw tick.
func getRows(networks NetworkSlice, associated *NetworkKey) []table.Row {
	// get registered column keys
	columns := ColumnsKeys()
	// get registered cell viewers
	viewers := GenerateCellViewers(associated)

	rows := make([]table.Row, len(networks))
	for rowID, e := range networks {
		entry := e

		row := make(table.RowData, len(columns))
		for _, key := range columns {
			row[key] = viewers[key](&entry)
		}

		rows[rowID] = table.NewRow(row)
	}

	return rows
}

type RefreshMsg NetworkSlice

// Invokes refresh table by refreshInterval.
func refreshTick(ds *DataSource, interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		// Copy networks stats
		networks := ds.NetworkSlice()

		// return RefreshMsg
		return RefreshMsg(networks)
	})
}
