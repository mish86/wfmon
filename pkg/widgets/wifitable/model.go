package wifitable

import (
	"time"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	defaultRefreshInterval = time.Second
	defaultTableHeight     = 10
	defaultTableWidth      = 120
)

var (
	defaultBaseStyle       = lipgloss.NewStyle()
	defaultHeaderStyle     = lipgloss.NewStyle().Foreground(lipgloss.NoColor{}).Bold(true)
	defaultSelectedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("229")).Background(lipgloss.Color("57")).Bold(true)
	defaultAssociatedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")).Bold(true)
)

type Model struct {
	table.Model
	viewport        viewport.Model
	dataSource      ds.NetworkProvider
	networks        netdata.Slice
	associated      netdata.Key
	selected        netdata.Key
	stationViewMode Cycler[StationViewMode]
	signalViewMode  Cycler[SignalViewMode]
	sort            Sort
	keys            KeyMap
}

func defaultViewMode() (Sort, Cycler[SignalViewMode], Cycler[StationViewMode]) {
	return SortBy(ColumnBarsKey)(DESC),
		BarsViewMode.Cycle(),
		BSSIDViewMode.Cycle()
}

type Option func(*Model)

func WithDataSource(dataSource ds.NetworkProvider) Option {
	return func(m *Model) {
		m.dataSource = dataSource
	}
}

func WithAssociated(key netdata.Key) Option {
	return func(m *Model) {
		m.associated = key
	}
}

func New(opts ...Option) *Model {
	sort, signalViewMode, stationViewMode := defaultViewMode()
	columns := GenerateColumns(sort, signalViewMode.Current())

	keys := NewKeyMap()
	t := table.New(columns).
		Border(table.Border{}).
		WithPageSize(defaultTableHeight).
		WithPaginationWrapping(false).
		WithMaxTotalWidth(defaultTableWidth).
		WithKeyMap(keys.KeyMap).
		WithBaseStyle(defaultBaseStyle).
		HeaderStyle(defaultHeaderStyle).
		HighlightStyle(defaultSelectedStyle).
		Focused(true)

	m := &Model{
		Model:           t,
		viewport:        viewport.New(defaultTableWidth, defaultTableHeight+2),
		keys:            keys,
		stationViewMode: stationViewMode,
		signalViewMode:  signalViewMode,
		sort:            sort,
		networks:        netdata.Slice{},
		dataSource:      ds.EmptyProvider{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) SetDataSource(dataSource ds.NetworkProvider) {
	m.dataSource = dataSource
}

func (m *Model) Keys() KeyMap {
	return m.keys
}

func (m *Model) View() string {
	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Left, m.Model.View()),
	)
	return m.viewport.View()
}
