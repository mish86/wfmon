package wifitable

import (
	"sort"
	"wfmon/pkg/utils/cmp"

	"golang.org/x/exp/constraints"
)

// Sort order enum.
type Order uint8

const (
	None Order = iota
	ASC
	DESC
)

func (o Order) String() string {
	//nolint:exhaustive // ignore
	switch o {
	case ASC:
		return "↓"
	case DESC:
		return "↑"
	default:
		return ""
	}
}

// Swaps sorting order.
func (o Order) swap() Order {
	if o == ASC {
		return DESC
	}

	return ASC
}

// Sorting definition.
type SortDef struct {
	col    string
	ord    Order
	sorter FncSorter
}

// Change sorting order.
func (s *SortDef) ChangeOrder() {
	s.ord = s.ord.swap()
}

// Returns @sort.Interface depending on order value.
func (s *SortDef) Sorter(networks NetworkSlice) sort.Interface {
	if s.ord == ASC || s.ord == None {
		return s.sorter(networks)
	}

	return &Inverser{s.sorter(networks)}
}

// Inverses resul of @sort.Interface.Less.
type Inverser struct {
	sort.Interface
}

// Inverses resul of @sort.Interface.Less.
func (a Inverser) Less(i, j int) bool {
	return !a.Interface.Less(i, j)
}

// Wraps NetworkSlice with @sort.Interface.
type FncSorter func(networks NetworkSlice) sort.Interface

// Returns Sort Definition regeisterd for given column.
// Default is @BySSIDByBSSIDSorter.
func ColumnSorterGenerator(column string) SortDef {
	sorters := map[string]SortDef{
		ColumnBSSIDTitle: {
			col:    ColumnBSSIDTitle,
			sorter: ByBSSIDSorter(),
		},
		ColumnSSIDTitle: {
			col: ColumnSSIDTitle,
			sorter: func(networks NetworkSlice) sort.Interface {
				return BySSIDSorter(networks)
			},
		},
		ColumnChanTitle: {
			col:    ColumnChanTitle,
			sorter: ByChannelSorter(),
		},
		ColumnWidthTitle: {
			col:    ColumnWidthTitle,
			sorter: ByChannelWidthSorter(),
		},
		ColumnBandTitle: {
			col:    ColumnBandTitle,
			sorter: ByBandwidthSorter(),
		},
		ColumnRSSITitle: {
			col:    ColumnRSSITitle,
			sorter: ByRSSISorter(),
		},
		ColumnNoiseTitle: {
			col:    ColumnNoiseTitle,
			sorter: ByNoiseSorter(),
		},
		ColumnSNRTitle: {
			col:    ColumnSNRTitle,
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

type DefaultSorter struct {
	len  func() int
	swap func(i, j int)
	less func(i, j int) bool
}

func Sorter[T constraints.Ordered](fncGet func(n NetworkSlice, i int) T) FncSorter {
	return func(n NetworkSlice) sort.Interface {
		return &DefaultSorter{
			len:  func() int { return len(n) },
			swap: func(i, j int) { n[i], n[j] = n[j], n[i] },
			less: func(i, j int) bool {
				cmp := cmp.Compare(fncGet(n, i), fncGet(n, j))
				if cmp == 0 {
					cmp = n[i].Key().Compare(n[j].Key())
				}

				return cmp < 0
			},
		}
	}
}
func (s DefaultSorter) Len() int           { return s.len() }
func (s DefaultSorter) Swap(i, j int)      { s.swap(i, j) }
func (s DefaultSorter) Less(i, j int) bool { return s.less(i, j) }

// Sort by BSSID asc.
func ByBSSIDSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) string { return n[i].BSSID })
}

// Sort by SSID asc.
type BySSIDSorter NetworkSlice

func (a BySSIDSorter) Len() int           { return len(a) }
func (a BySSIDSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySSIDSorter) Less(i, j int) bool { return a[i].Key().Compare(a[j].Key()) < 0 }

// Sort by Channel asc.
func ByChannelSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return n[i].Channel })
}

// Sort by Channel Width asc.
func ByChannelWidthSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return n[i].ChannelWidth })
}

// Sort by Bandwidth asc.
func ByBandwidthSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) string { return n[i].Band })
}

// Sort by RSSI asc.
func ByRSSISorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].RSSI) })
}

// Sort by Noise asc.
func ByNoiseSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].Noise) })
}

// Sort by SNR asc.
func BySNRSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].SNR) })
}
