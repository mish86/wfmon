package wifitable

import (
	"time"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"
	log "wfmon/pkg/logger"
	"wfmon/pkg/widgets/color"
	column "wfmon/pkg/widgets/wifitable/col"
	order "wfmon/pkg/widgets/wifitable/ord"

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
	viewport   viewport.Model
	dataSource ds.NetworkProvider
	networks   netdata.Slice
	colors     map[netdata.Key]color.HexColor
	associated netdata.Key
	// selected   netdata.Key
	columns []column.Column
	sort    column.Sort
	keys    KeyMap
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

func WithFocused(focus bool) Option {
	return func(m *Model) {
		m.Focused(focus)
	}
}

func New(opts ...Option) *Model {
	sort := defaultSort()
	cols := columns()

	keys := NewKeyMap()
	t := table.New(column.Converter(cols)(sort)).
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
		Model:      t,
		viewport:   viewport.New(defaultTableWidth, defaultTableHeight+2),
		keys:       keys,
		columns:    cols,
		sort:       sort,
		networks:   netdata.Slice{},
		colors:     map[netdata.Key]color.HexColor{},
		dataSource: ds.EmptyProvider{},
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

func (m *Model) SetWidth(w int) {
	m.Model.WithMaxTotalWidth(w)
	m.viewport.Width = w
}

// Returns table width calculated by width of visible simple columns.
func (m *Model) Width() int {
	width := 0
	for _, col := range m.columns {
		width += col.Width()
	}

	return width
}

// Searches a simple column with requested key. If not found uses default @SSIDKey.
// Returns generator which accepts sorting order to build Sort definition.
func sortBy(key string) func(ord order.Dir) column.Sort {
	cols := simpleColumns()

	col, found := cols[key]
	if !found {
		log.Warnf("column key %s not found to sort table, using default %s", key, SSIDKey)
		col = cols[SSIDKey]
	}

	return func(ord order.Dir) column.Sort {
		return column.NewSort(key, col.Sorter()).WithOrder(ord)
	}
}

// Returns default Sort for the table.
func defaultSort() column.Sort {
	return sortBy(BarsKey)(order.DESC)
}

// Returns all keys of simple columns visible as of now.
func visibleColumnKeys(cols []column.Column) []string {
	keys := make([]string, len(cols))

	i := 0
	for _, col := range cols {
		keys[i] = col.Key()
		i++
	}

	return keys
}

func (m *Model) GetSelectedNetwork() netdata.Network {
	cursor := m.GetHighlightedRowIndex()

	// FIXME: race at access to networks

	// no data
	if len(m.networks) == 0 {
		return netdata.Network{}
	}

	// out of bounds
	if cursor < 0 || cursor >= len(m.networks) {
		log.Errorf("cursor %d out of bounds networks: %v", cursor, m.networks)
		return netdata.Network{}
	}

	return m.networks[cursor]
}
