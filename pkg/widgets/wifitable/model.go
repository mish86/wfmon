package wifitable

import (
	"strings"
	"time"
	"wfmon/pkg/network"

	"github.com/charmbracelet/bubbles/help"
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

func defaultViewMode() (SignalViewMode, Sort) {
	return BarsViewMode, SortBy(ColumnBarsKey)(DESC)
}

func NewTable(ds *DataSource, network network.Network) *Model {
	help := help.New()
	help.ShowAll = true

	signalViewMode, sort := defaultViewMode()
	columns := GenerateColumns(sort, uint8(signalViewMode))

	keys := NewKeyMap()
	t := table.New(columns).
		Border(table.Border{}).
		WithPageSize(defaultTableHeight).
		WithPaginationWrapping(false).
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

func (m *Model) onRefreshMsg(msg RefreshMsg) {
	// get networks from last tick
	m.networks = NetworkSlice(msg)

	// apply current sorting
	m.sort.Sort(m.networks)

	// apply columns and rows to table
	m.refresh()
}
