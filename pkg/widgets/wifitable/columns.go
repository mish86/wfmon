package wifitable

import (
	"fmt"
	"sort"
	"strconv"

	"wfmon/pkg/widgets"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"golang.org/x/exp/constraints"
)

const (
	ColumnSSIDKey      = "Network"
	ColumnBSSIDKey     = "BSSID"
	ColumnManufKey     = "Manuf"
	ColumnManufLongKey = "Manufactor"
	ColumnChanKey      = "Chan"
	ColumnWidthKey     = "Width"
	ColumnBandKey      = "Band"
	ColumnRSSIKey      = "RSSI"
	ColumnQualityKey   = "Quality"
	ColumnBarsKey      = "Bars"
	ColumnNoiseKey     = "Noise"
	ColumnSNRKey       = "SNR"
)

// Provides cycling of view for multi column viewers.
type Cycler[T constraints.Integer] interface {
	Next() Cycler[T]
	Current() T
	Prev() Cycler[T]
}

type defaultCycler[T constraints.Integer] struct {
	start, current, end T
}

func (c *defaultCycler[T]) Next() Cycler[T] {
	c.current++
	if c.current > c.end {
		c.current = c.start
	}

	return c
}

func (c *defaultCycler[T]) Current() T {
	return c.current
}

func (c *defaultCycler[T]) Prev() Cycler[T] {
	c.current--
	if c.current < c.start {
		c.current = c.end
	}

	return c
}

func Cycle[T constraints.Integer](start, current, end T) Cycler[T] {
	return &defaultCycler[T]{
		start:   start,
		current: current,
		end:     end,
	}
}

// Provides all keys supproted by multi column viewer.
type Enumerator interface {
	Keys() []string
}

// Provides key presentation for current state of multi column viewer.
type Keyer interface {
	Enumerator
	Key() string
}

// Generates @table.Column depending on current state of multi column viewer.
type MultiViewColumnGenerator interface {
	Enumerator
	Column(widths map[string]int, sort Sort) func() table.Column
}

// Multi column viewer for signal.
// Supports RSSI, Quality, Bars.
type SignalViewMode uint8

const (
	RSSIViewMode SignalViewMode = iota
	QualityViewMode
	BarsViewMode
)

func (v SignalViewMode) Keys() []string {
	return []string{
		ColumnRSSIKey,
		ColumnQualityKey,
		ColumnBarsKey,
	}
}

func (v SignalViewMode) Cycle() Cycler[SignalViewMode] {
	return Cycle(RSSIViewMode, v, BarsViewMode)
}

func (v SignalViewMode) Key() string {
	switch v {
	case RSSIViewMode:
		return ColumnRSSIKey
	case QualityViewMode:
		return ColumnQualityKey
	case BarsViewMode:
		return ColumnBarsKey
	default:
		return ColumnRSSIKey
	}
}

func (v SignalViewMode) Column(widths map[string]int, sort Sort) func() table.Column {
	return func() table.Column {
		switch v {
		case RSSIViewMode:
			return newColumn(widths, sort)(ColumnRSSIKey)
		case QualityViewMode:
			return newColumn(widths, sort)(ColumnQualityKey)
		case BarsViewMode:
			return newColumn(widths, sort)(ColumnBarsKey)
		default:
			return newColumn(widths, sort)(ColumnRSSIKey)
		}
	}
}

// Multi column viewer for station.
// Supports BSSID, Manuf, Manufactor.
type StationViewMode uint8

const (
	BSSIDViewMode StationViewMode = iota
	ManufViewMode
	ManufLongViewMode
)

func (v StationViewMode) Keys() []string {
	return []string{
		ColumnBSSIDKey,
		ColumnManufKey,
		ColumnManufLongKey,
	}
}

func (v StationViewMode) Cycle() Cycler[StationViewMode] {
	return Cycle(BSSIDViewMode, v, ManufLongViewMode)
}

func (v StationViewMode) Key() string {
	switch v {
	case BSSIDViewMode:
		return ColumnBSSIDKey
	case ManufViewMode:
		return ColumnManufKey
	case ManufLongViewMode:
		return ColumnManufLongKey
	default:
		return ColumnBSSIDKey
	}
}

func (v StationViewMode) Column(widths map[string]int, sort Sort) func() table.Column {
	style := lipgloss.NewStyle().
		Align(lipgloss.Left).
		PaddingLeft(1)
	return func() table.Column {
		switch v {
		case BSSIDViewMode:
			return newColumn(widths, sort)(ColumnBSSIDKey).WithStyle(style)
		case ManufViewMode:
			return newColumn(widths, sort)(ColumnManufKey).WithStyle(style)
		case ManufLongViewMode:
			return newColumn(widths, sort)(ColumnManufLongKey).WithStyle(style)
		default:
			return newColumn(widths, sort)(ColumnBSSIDKey).WithStyle(style)
		}
	}
}

// Returns all registered columns keys.
func ColumnsKeys() []string {
	return []string{
		ColumnSSIDKey,
		ColumnBSSIDKey,
		ColumnManufKey,
		ColumnManufLongKey,
		ColumnChanKey,
		ColumnWidthKey,
		ColumnBandKey,
		ColumnRSSIKey,
		ColumnQualityKey,
		ColumnBarsKey,
		ColumnNoiseKey,
		ColumnSNRKey,
	}
}

// Returns visible columns.
// Keyers required for multiview columns. The proper one selected by key columns.
// Supoprted multi columns viewers:
// - station: BSSID (Default), Manuf, Manufactor
// - signal RSSI (Default), Quality, Bars
func VisibleColumnsKeys(enums ...Keyer) []string {
	var station Keyer = BSSIDViewMode
	var signal Keyer = RSSIViewMode

	for _, enum := range enums {
		switch enum.(type) {
		case StationViewMode:
			station = enum
		case SignalViewMode:
			signal = enum
		}
	}

	return []string{
		ColumnSSIDKey,
		station.Key(),
		ColumnChanKey,
		ColumnWidthKey,
		ColumnBandKey,
		signal.Key(),
		ColumnNoiseKey,
		ColumnSNRKey,
	}
}

// Returns predefined columns width
func columnsWidth() map[string]int {
	//nolint:gomnd // ignore
	return map[string]int{
		ColumnSSIDKey:      20,
		ColumnBSSIDKey:     20,
		ColumnManufKey:     10,
		ColumnManufLongKey: 30,
		ColumnChanKey:      7,
		ColumnWidthKey:     8,
		ColumnBandKey:      7,
		ColumnRSSIKey:      7,
		ColumnQualityKey:   10,
		ColumnBarsKey:      7,
		ColumnNoiseKey:     8,
		ColumnSNRKey:       5,
	}
}

// Returns @table.Column with applied sorting in title and width.
func newColumn(widths map[string]int, sort Sort) func(key string) table.Column {
	return func(key string) table.Column {
		title := key
		if sort.key == key {
			title = fmt.Sprintf("%s %s", key, sort.ord)
		}
		return table.NewColumn(key, title, widths[key])
	}
}

// Returns default columns.
// Multiview column viewers can be passed in any order. The proper one selected by key columns.
// Supoprted multi columns viewers:
// - station: BSSID (Default), Manuf, Manufactor
// - signal RSSI (Default), Quality, Bars
func GenerateColumns(sort Sort, enums ...MultiViewColumnGenerator) []table.Column {
	widths := columnsWidth()

	var station MultiViewColumnGenerator = BSSIDViewMode
	var signal MultiViewColumnGenerator = RSSIViewMode

	for _, enum := range enums {
		switch enum.(type) {
		case StationViewMode:
			station = enum
		case SignalViewMode:
			signal = enum
		}
	}

	return []table.Column{
		newColumn(widths, sort)(ColumnSSIDKey).WithStyle(
			lipgloss.NewStyle().
				Align(lipgloss.Left),
		),
		station.Column(widths, sort)(),
		newColumn(widths, sort)(ColumnChanKey),
		newColumn(widths, sort)(ColumnWidthKey),
		newColumn(widths, sort)(ColumnBandKey),
		signal.Column(widths, sort)(),
		newColumn(widths, sort)(ColumnNoiseKey),
		newColumn(widths, sort)(ColumnSNRKey),
	}
}

// Returns a sorter regeistered for given column.
// Default is BySSIDSorter.
func GenerateSorters(column string) Sort {
	sorters := map[string]Sort{
		ColumnSSIDKey: {
			key: ColumnSSIDKey,
			sorter: func(networks NetworkSlice) sort.Interface {
				return BySSIDSorter(networks)
			},
		},
		ColumnBSSIDKey: {
			key:    ColumnBSSIDKey,
			sorter: ByBSSIDSorter(),
		},
		ColumnManufKey: {
			key:    ColumnManufKey,
			sorter: ByBSSIDSorter(),
		},
		ColumnManufLongKey: {
			key:    ColumnManufLongKey,
			sorter: ByBSSIDSorter(),
		},
		ColumnChanKey: {
			key:    ColumnChanKey,
			sorter: ByChannelSorter(),
		},
		ColumnWidthKey: {
			key:    ColumnWidthKey,
			sorter: ByChannelWidthSorter(),
		},
		ColumnBandKey: {
			key:    ColumnBandKey,
			sorter: ByBandwidthSorter(),
		},
		ColumnRSSIKey: {
			key:    ColumnRSSIKey,
			sorter: ByRSSISorter(),
		},
		ColumnQualityKey: {
			key:    ColumnQualityKey,
			sorter: ByQualitySorter(),
		},
		ColumnBarsKey: {
			key:    ColumnBarsKey,
			sorter: ByBarsSorter(),
		},
		ColumnNoiseKey: {
			key:    ColumnNoiseKey,
			sorter: ByNoiseSorter(),
		},
		ColumnSNRKey: {
			key:    ColumnSNRKey,
			sorter: BySNRSorter(),
		},
	}

	sorter, ok := sorters[column]
	if !ok {
		// default sorter
		sorter = sorters[ColumnSSIDKey]
	}

	return sorter
}

// Network data field getter. Returns cell presentation in accetable as table.RowData.
type FncCellViewer func(data *NetworkData) any

// Returns string presentation of cell by column title.
func GenerateCellViewers(associated *NetworkKey) map[any]FncCellViewer {
	widths := columnsWidth()

	associatedStyle := func(data *NetworkData) lipgloss.Style {
		if data.Key().Compare(associated) == 0 {
			return defaultAssociatedStyle
		}

		return defaultBaseStyle
	}

	getters := map[any]FncCellViewer{
		ColumnSSIDKey: func(data *NetworkData) any {
			// The goal is to keep space between columns.
			// Setup border.left/border.right with ' ' does not work and has side effects.
			// Padding/Margin does not work properly on this column or right after this one.
			// ref. https://github.com/Evertras/bubble-table/issues/130
			// Thus manually truncate string.
			style := associatedStyle(data)
			return table.NewStyledCell(widgets.StringWithTail(data.NetworkName, widths[ColumnSSIDKey]-1), style)
		},
		ColumnBSSIDKey: func(data *NetworkData) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(associatedStyle(data))
			return table.NewStyledCell(data.BSSID, style)
		},
		ColumnManufKey: func(data *NetworkData) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(associatedStyle(data))
			return table.NewStyledCell(data.Manuf, style)
		},
		ColumnManufLongKey: func(data *NetworkData) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(associatedStyle(data))
			return table.NewStyledCell(data.ManufLong, style)
		},
		ColumnChanKey: func(data *NetworkData) any {
			style := associatedStyle(data)
			return style.Render(strconv.Itoa(int(data.Channel)))
		},
		ColumnWidthKey: func(data *NetworkData) any {
			style := associatedStyle(data)
			return style.Render(strconv.Itoa(int(data.ChannelWidth)))
		},
		ColumnBandKey: func(data *NetworkData) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(associatedStyle(data))
			return table.NewStyledCell(data.Band.String(), style)
		},
		ColumnRSSIKey: func(data *NetworkData) any {
			style := associatedStyle(data)
			return style.Render(strconv.Itoa(int(data.RSSI)))
		},
		ColumnQualityKey: func(data *NetworkData) any {
			style := associatedStyle(data)
			return style.Render(data.Quality.String())
		},
		ColumnBarsKey: func(data *NetworkData) any {
			bars := Bars(data.Quality)
			style := bars.Style().Inherit(associatedStyle(data))
			return style.Render(bars.String())
		},
		ColumnNoiseKey: func(data *NetworkData) any {
			style := associatedStyle(data)
			return style.Render(strconv.Itoa(int(data.Noise)))
		},
		ColumnSNRKey: func(data *NetworkData) any {
			style := associatedStyle(data)
			return style.Render(strconv.Itoa(int(data.SNR)))
		},
	}

	return getters
}

type Bars Quality

func (q Bars) Style() lipgloss.Style {
	//nolint:gomnd // ignore
	switch {
	case q >= 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#77dd77")) // green
	case 60 <= q && q < 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#a7c7e7")) // blue
	case 40 <= q && q < 60:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb347")) // orange
	case 20 <= q && q < 40:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")) // red
	case q < 20:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")) // red
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")) // red
	}
}

// Returns bars presentation of signal quality.
func (q Bars) String() string {
	//nolint:gomnd // ignore
	switch {
	case q >= 80:
		return "▂▄▆█"
	case 60 <= q && q < 80:
		return "▂▄▆_"
	case 40 <= q && q < 60:
		return "▂▄__"
	case 20 <= q && q < 40:
		return "▂___"
	case q < 20:
		return "____"
	default:
		return "____"
	}
}
