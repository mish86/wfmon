package wifitable

import (
	"sort"
	"strconv"

	"github.com/charmbracelet/bubbles/table"
)

const (
	ColumnSSIDTitle    = "Network"
	ColumnBSSIDTitle   = "BSSID"
	ColumnChanTitle    = "Chan"
	ColumnWidthTitle   = "Width"
	ColumnBandTitle    = "Band"
	ColumnRSSITitle    = "RSSI"
	ColumnQualityTitle = "Quality"
	ColumnBarsTitle    = "Bars"
	ColumnNoiseTitle   = "Noise"
	ColumnSNRTitle     = "SNR"
)

const (
	defaultSortingOrderWidth = 2
)

func ColumnSSID(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnSSIDTitle, Width: 17},
	)
}

func ColumnBSSID(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnBSSIDTitle, Width: 20},
	)
}

func ColumnChan(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnChanTitle, Width: 4 + defaultSortingOrderWidth},
	)
}

func ColumnWidth(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnWidthTitle, Width: 5 + defaultSortingOrderWidth},
	)
}

func ColumnBand(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnBandTitle, Width: 4 + defaultSortingOrderWidth},
	)
}

func ColumnSignal(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnRSSITitle, Width: 4 + defaultSortingOrderWidth},
		table.Column{Title: ColumnQualityTitle, Width: 7 + defaultSortingOrderWidth},
		table.Column{Title: ColumnBarsTitle, Width: 4 + defaultSortingOrderWidth},
	)
}

func ColumnNoise(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnNoiseTitle, Width: 5 + defaultSortingOrderWidth},
	)
}

func ColumnSNR(num int) ColumnViewer {
	//nolint:gomnd // ignore
	return NewColumnViewer(num,
		table.Column{Title: ColumnSNRTitle, Width: 4 + defaultSortingOrderWidth},
	)
}

type FncColumnGenerator func(num int) ColumnViewer

func GenerateColumns(generators ...FncColumnGenerator) ColumnViewSlice {
	columns := make(ColumnViewSlice, len(generators))
	for i, generator := range generators {
		// column number starts from 1
		columns[i] = generator(i + 1)
	}

	return columns
}

// Returns default columns viewers.
func GenerateDefaultColumns() ColumnViewSlice {
	return GenerateColumns(
		ColumnSSID,
		ColumnBSSID,
		ColumnChan,
		ColumnWidth,
		ColumnBand,
		ColumnSignal,
		ColumnNoise,
		ColumnSNR,
	)
}

// Returns a sort regeistered for given column.
// Default is @BySSIDSorter.
func ColumnSorterGenerator(column string) Sort {
	sorters := map[string]Sort{
		ColumnBSSIDTitle: {
			title:  ColumnBSSIDTitle,
			sorter: ByBSSIDSorter(),
		},
		ColumnSSIDTitle: {
			title: ColumnSSIDTitle,
			sorter: func(networks NetworkSlice) sort.Interface {
				return BySSIDSorter(networks)
			},
		},
		ColumnChanTitle: {
			title:  ColumnChanTitle,
			sorter: ByChannelSorter(),
		},
		ColumnWidthTitle: {
			title:  ColumnWidthTitle,
			sorter: ByChannelWidthSorter(),
		},
		ColumnBandTitle: {
			title:  ColumnBandTitle,
			sorter: ByBandwidthSorter(),
		},
		ColumnRSSITitle: {
			title:  ColumnRSSITitle,
			sorter: ByRSSISorter(),
		},
		ColumnQualityTitle: {
			title:  ColumnQualityTitle,
			sorter: ByQualitySorter(),
		},
		ColumnBarsTitle: {
			title:  ColumnBarsTitle,
			sorter: ByBarsSorter(),
		},
		ColumnNoiseTitle: {
			title:  ColumnNoiseTitle,
			sorter: ByNoiseSorter(),
		},
		ColumnSNRTitle: {
			title:  ColumnSNRTitle,
			sorter: BySNRSorter(),
		},
	}

	sorter, ok := sorters[column]
	if !ok {
		// default sorter
		sorter = sorters[ColumnSSIDTitle]
	}

	return sorter
}

// Returns string presentation of cell by column title.
func GenerateRowGetters() map[string]FncRowGetter {
	getters := map[string]FncRowGetter{
		ColumnSSIDTitle:    func(data *NetworkData) string { return data.NetworkName },
		ColumnBSSIDTitle:   func(data *NetworkData) string { return data.BSSID },
		ColumnChanTitle:    func(data *NetworkData) string { return strconv.Itoa(int(data.Channel)) },
		ColumnWidthTitle:   func(data *NetworkData) string { return strconv.Itoa(int(data.ChannelWidth)) },
		ColumnBandTitle:    func(data *NetworkData) string { return data.Band.String() },
		ColumnRSSITitle:    func(data *NetworkData) string { return strconv.Itoa(int(data.RSSI)) },
		ColumnQualityTitle: func(data *NetworkData) string { return data.Quality.String() },
		ColumnBarsTitle:    func(data *NetworkData) string { return data.Quality.Bars() },
		ColumnNoiseTitle:   func(data *NetworkData) string { return strconv.Itoa(int(data.Noise)) },
		ColumnSNRTitle:     func(data *NetworkData) string { return strconv.Itoa(int(data.SNR)) },
	}

	return getters
}
