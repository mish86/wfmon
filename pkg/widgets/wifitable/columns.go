package wifitable

import (
	"fmt"
	"sort"
	"strconv"

	"wfmon/pkg/widgets"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

const (
	ColumnSSIDKey    = "Network"
	ColumnBSSIDKey   = "BSSID"
	ColumnChanKey    = "Chan"
	ColumnWidthKey   = "Width"
	ColumnBandKey    = "Band"
	ColumnRSSIKey    = "RSSI"
	ColumnQualityKey = "Quality"
	ColumnBarsKey    = "Bars"
	ColumnNoiseKey   = "Noise"
	ColumnSNRKey     = "SNR"
)

type SignalViewMode uint8

const (
	RSSIViewMode SignalViewMode = iota
	QualityViewMode
	BarsViewMode
)

func (v SignalViewMode) Next() SignalViewMode {
	mode := uint8(v) + 1
	if mode > uint8(BarsViewMode) {
		mode = uint8(RSSIViewMode)
	}

	return SignalViewMode(mode)
}

func (v SignalViewMode) Prev() SignalViewMode {
	mode := uint8(v) - 1
	if mode > uint8(RSSIViewMode) {
		mode = uint8(BarsViewMode)
	}

	return SignalViewMode(mode)
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

// Returns columns keys.
func ColumnsKeys() []string {
	return []string{
		ColumnSSIDKey,
		ColumnBSSIDKey,
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

func VisibleColumnsKeys(signal SignalViewMode) []string {
	return []string{
		ColumnSSIDKey,
		ColumnBSSIDKey,
		ColumnChanKey,
		ColumnWidthKey,
		ColumnBandKey,
		signal.Key(),
		ColumnNoiseKey,
		ColumnSNRKey,
	}
}

func columnsWidth() map[string]int {
	//nolint:gomnd // ignore
	return map[string]int{
		ColumnSSIDKey:    20,
		ColumnBSSIDKey:   20,
		ColumnChanKey:    7,
		ColumnWidthKey:   8,
		ColumnBandKey:    7,
		ColumnRSSIKey:    7,
		ColumnQualityKey: 10,
		ColumnBarsKey:    7,
		ColumnNoiseKey:   8,
		ColumnSNRKey:     5,
	}
}

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
func GenerateColumns(modes ...any) []table.Column {
	widths := columnsWidth()

	sort := SortBy(ColumnSSIDKey)(None)
	if len(modes) > 0 {
		sort, _ = modes[0].(Sort)
	}

	signal := RSSIViewMode
	if len(modes) > 1 {
		signal = SignalViewMode(modes[1].(uint8))
	}

	return []table.Column{
		newColumn(widths, sort)(ColumnSSIDKey).WithStyle(
			lipgloss.NewStyle().
				Align(lipgloss.Left),
		),
		newColumn(widths, sort)(ColumnBSSIDKey).WithStyle(
			lipgloss.NewStyle().
				Align(lipgloss.Left).
				PaddingLeft(1),
		),
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
