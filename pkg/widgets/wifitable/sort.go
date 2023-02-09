package wifitable

import (
	"sort"
	"wfmon/pkg/utils/cmp"

	"golang.org/x/exp/constraints"
)

// Sorting order enum.
type Order uint8

const (
	None Order = iota
	ASC
	DESC
)

// Returns string presentation of sorting order.
// Used in column title view.
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

// Describes sorting order for a column and sorting function.
type Sort struct {
	key    string
	ord    Order
	sorter FncSorter
}

func SortBy(key string) func(ord Order) Sort {
	sort := GenerateSorters(key)

	return func(ord Order) Sort {
		sort.ord = ord
		return sort
	}
}

// Change sorting order.
func (s *Sort) SwapOrder() {
	s.ord = s.ord.swap()
}

// Returns @sort.Interface depending on order value.
func (s *Sort) Sorter(networks NetworkSlice) sort.Interface {
	if s.ord == ASC || s.ord == None {
		return s.sorter(networks)
	}

	return &Inverser{s.sorter(networks)}
}

func (s *Sort) Sort(networks NetworkSlice) {
	sort.Sort(s.Sorter(networks))
}

// Returns sorting order.
func (s *Sort) Order() Order {
	return s.ord
}

// Returns sorting order.
func (s *Sort) SetOrder(ord Order) {
	s.ord = ord
}

// Returns column title.
func (s *Sort) Key() string {
	return s.key
}

// Inverses resul of @sort.Interface.Less.
type Inverser struct {
	sort.Interface
}

// Inverses result of @sort.Interface.Less.
func (a Inverser) Less(i, j int) bool {
	return !a.Interface.Less(i, j)
}

// Wraps NetworkSlice with @sort.Interface.
type FncSorter func(networks NetworkSlice) sort.Interface

// Implements default sorter behaviour for all columns.
type DefaultSorter struct {
	len  func() int
	swap func(i, j int)
	less func(i, j int) bool
}

// Takes an implementation of getter and returns sorter as high order func.
func Sorter[T constraints.Ordered](fncGet func(n NetworkSlice, i int) T) FncSorter {
	return func(n NetworkSlice) sort.Interface {
		return &DefaultSorter{
			len:  func() int { return len(n) },
			swap: func(i, j int) { n[i], n[j] = n[j], n[i] },
			less: func(i, j int) bool {
				// first sort by table field
				cmp := cmp.Compare(fncGet(n, i), fncGet(n, j))
				// then sort by table key
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
// Sorts only by network table key.
type BySSIDSorter NetworkSlice

func (a BySSIDSorter) Len() int           { return len(a) }
func (a BySSIDSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySSIDSorter) Less(i, j int) bool { return a[i].Key().Compare(a[j].Key()) < 0 }

// Sort by Channel asc.
func ByChannelSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].Channel) })
}

// Sort by Channel Width asc.
func ByChannelWidthSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].ChannelWidth) })
}

// Sort by Bandwidth asc.
func ByBandwidthSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].Band) })
}

// Sort by RSSI asc.
func ByRSSISorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].RSSI) })
}

// Sort by Quality asc.
func ByQualitySorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].Quality) })
}

// Sort by Quality/Bars asc.
func ByBarsSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].Quality) })
}

// Sort by Noise asc.
func ByNoiseSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].Noise) })
}

// Sort by SNR asc.
func BySNRSorter() FncSorter {
	return Sorter(func(n NetworkSlice, i int) int { return int(n[i].SNR) })
}
