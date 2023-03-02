package wifitable

import (
	"strconv"
	netdata "wfmon/pkg/data/net"
	column "wfmon/pkg/widgets/wifitable/col"
	"wfmon/pkg/widgets/wifitable/row"
	"wfmon/pkg/widgets/wifitable/sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// All known simple column keys.
const (
	HashKey       = "#"
	SSIDKey       = netdata.SSIDKey
	BSSIDKey      = netdata.BSSIDKey
	ManufKey      = netdata.ManufKey
	ManufactorKey = netdata.ManufLongKey
	ChanKey       = netdata.ChanKey
	WidthKey      = netdata.WidthKey
	BandKey       = netdata.BandKey
	RSSIKey       = netdata.RSSIKey
	QualityKey    = netdata.QualityKey
	BarsKey       = netdata.BarsKey
	NoiseKey      = netdata.NoiseKey
	SNRKey        = netdata.SNRKey
)

// Returns predefined columns width.
func widths() map[string]int {
	//nolint:gomnd // ignore
	return map[string]int{
		HashKey:       2,
		SSIDKey:       20,
		BSSIDKey:      20,
		ManufKey:      10,
		ManufactorKey: 30,
		ChanKey:       7,
		WidthKey:      8,
		BandKey:       7,
		RSSIKey:       7,
		QualityKey:    10,
		BarsKey:       7,
		NoiseKey:      8,
		SNRKey:        5,
	}
}

// Returns new simple column with key, sorter and width taken from @widths.
func newColumn(key string, sorter sort.FncSorter) column.Simple {
	return column.NewSimple(key, widths()[key]).
		WithSorter(sorter)
}

func HashColumn() column.Simple {
	return newColumn(HashKey, BySSIDSorter())
}

func SSIDColumn() column.Simple {
	return newColumn(SSIDKey, BySSIDSorter()).
		WithStyle(lipgloss.NewStyle().
			Align(lipgloss.Left))
}

func BSSIDColumn() column.Simple {
	return newColumn(BSSIDKey, ByBSSIDSorter())
}

func ManufColumn() column.Simple {
	return newColumn(ManufKey, ByManufSorter())
}

func ManufactorColumn() column.Simple {
	return newColumn(ManufactorKey, ByManufLongSorter())
}

func ChannelColumn() column.Simple {
	return newColumn(ChanKey, ByChannelSorter())
}

func WidthColumn() column.Simple {
	return newColumn(WidthKey, ByChannelWidthSorter())
}

func BandColumn() column.Simple {
	return newColumn(BandKey, ByBandwidthSorter())
}

func RSSIColumn() column.Simple {
	return newColumn(RSSIKey, ByRSSISorter())
}

func QualityColumn() column.Simple {
	return newColumn(QualityKey, ByQualitySorter())
}

func BarsColumn() column.Simple {
	return newColumn(BarsKey, ByQualitySorter())
}

func NoiseColumn() column.Simple {
	return newColumn(NoiseKey, ByNoiseSorter())
}

func SNRColumn() column.Simple {
	return newColumn(SNRKey, BySNRSorter())
}

func SignalColumn() column.Multiple {
	return column.NewMultiple(BarsColumn(), RSSIColumn(), QualityColumn())
}

func StationColumn() column.Multiple {
	return column.NewMultiple(BSSIDColumn(), ManufColumn(), ManufactorColumn())
}

// Index of MultiColumns in @columns array.
// Hash column is not registered in hot keys for sorting.
const (
	StationMColumnIdx = 2
	SignalMColumnIdx  = 6
)

// Returns an ordered array of columns to view in a table.
// Column numbering starts from 1 and from Network (SSID).
// Hash column should not be registered in hot keys for sorting and swaping of multi column view.
func columns() []column.Column {
	return []column.Column{
		HashColumn(),
		SSIDColumn(),
		StationColumn(),
		ChannelColumn(),
		WidthColumn(),
		BandColumn(),
		SignalColumn(),
		NoiseColumn(),
		SNRColumn(),
	}
}

// Returns all known simple columns as map with sorters and widths.
func simpleColumns() map[string]column.Simple {
	return map[string]column.Simple{
		HashKey:       HashColumn(),
		SSIDKey:       SSIDColumn(),
		BSSIDKey:      BSSIDColumn(),
		ManufKey:      ManufColumn(),
		ManufactorKey: ManufactorColumn(),
		ChanKey:       ChannelColumn(),
		WidthKey:      WidthColumn(),
		BandKey:       BandColumn(),
		RSSIKey:       RSSIColumn(),
		QualityKey:    QualityColumn(),
		BarsKey:       BarsColumn(),
		NoiseKey:      NoiseColumn(),
		SNRKey:        SNRColumn(),
	}
}

// Returns registered cell viewers for all simple column keys.
func cellViewers() map[string]row.FncCellViewer {
	return map[string]row.FncCellViewer{
		HashKey: func(row *row.Data) any {
			return table.NewStyledCell(" ", lipgloss.NewStyle().Background(row.GetHashColor()))
		},
		SSIDKey: func(row *row.Data) any {
			// The goal is to keep space between columns.
			// Setup border.left/border.right with ' ' does not work and has side effects.
			// Padding/Margin does not work properly on this column or right after this one.
			// ref. https://github.com/Evertras/bubble-table/issues/130
			// Thus manually truncate string.
			// style := associatedStyle(data)
			// return table.NewStyledCell(reflow.StringWithTail(data.NetworkName, widths[ColumnSSIDKey]-1), style)
			return table.NewStyledCell(row.NetworkName, row.GetRowStyle())
		},
		BSSIDKey: func(row *row.Data) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(row.GetRowStyle())
			return table.NewStyledCell(row.BSSID, style)
		},
		ManufKey: func(row *row.Data) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(row.GetRowStyle())
			return table.NewStyledCell(row.Manuf, style)
		},
		ManufactorKey: func(row *row.Data) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(row.GetRowStyle())
			return table.NewStyledCell(row.ManufLong, style)
		},
		ChanKey: func(row *row.Data) any {
			return table.NewStyledCell(strconv.Itoa(int(row.Channel)), row.GetRowStyle())
		},
		WidthKey: func(row *row.Data) any {
			return table.NewStyledCell(strconv.Itoa(int(row.ChannelWidth)), row.GetRowStyle())
		},
		BandKey: func(row *row.Data) any {
			style := lipgloss.NewStyle().AlignHorizontal(lipgloss.Left).Inherit(row.GetRowStyle())
			return table.NewStyledCell(row.Band.String(), style)
		},
		RSSIKey: func(row *row.Data) any {
			return table.NewStyledCell(strconv.Itoa(int(row.RSSI)), row.GetRowStyle())
		},
		QualityKey: func(row *row.Data) any {
			return table.NewStyledCell(row.Quality.String(), row.GetRowStyle())
		},
		BarsKey: func(row *row.Data) any {
			bars := Bars(row.Quality)
			return table.NewStyledCell(bars.String(), bars.Style().Inherit(row.GetRowStyle()))
		},
		NoiseKey: func(row *row.Data) any {
			return table.NewStyledCell(strconv.Itoa(int(row.Noise)), row.GetRowStyle())
		},
		SNRKey: func(row *row.Data) any {
			return table.NewStyledCell(strconv.Itoa(int(row.SNR)), row.GetRowStyle())
		},
	}
}
