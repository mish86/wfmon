package wifitable

import "github.com/charmbracelet/bubbles/table"

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
		table.Column{Title: ColumnRSSITitle, Width: 5 + defaultSortingOrderWidth},
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
		table.Column{Title: ColumnSNRTitle, Width: 5 + defaultSortingOrderWidth},
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
