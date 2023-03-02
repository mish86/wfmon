package sort

import (
	"sort"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/utils/cmp"
	order "wfmon/pkg/widgets/wifitable/ord"

	"golang.org/x/exp/constraints"
)

// Describes sorting order for a column and sorting function.
type Sort struct {
	key    string
	ord    order.Dir
	sorter FncSorter
}

// Returns new Sort object.
func New(key string, sorter FncSorter) Sort {
	return Sort{
		key:    key,
		ord:    order.None,
		sorter: sorter,
	}
}

// Change sorting order.
func (s Sort) SwapOrder() Sort {
	s.ord = s.ord.Swap()
	return s
}

// Returns @sort.Interface depending on order value.
func (s Sort) Sorter(networks netdata.Slice) sort.Interface {
	if s.ord == order.ASC || s.ord == order.None {
		return s.sorter(networks)
	}

	return &Inverser{s.sorter(networks)}
}

func (s Sort) Sort(networks netdata.Slice) {
	sort.Sort(s.Sorter(networks))
}

// Returns sorting order.
func (s Sort) Order() order.Dir {
	return s.ord
}

// Returns sorting order.
func (s Sort) WithOrder(ord order.Dir) Sort {
	s.ord = ord
	return s
}

// Returns column title.
func (s Sort) Key() string {
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
type FncSorter func(networks netdata.Slice) sort.Interface

// Implements default sorter behaviour for all columns.
type defaultSorter struct {
	len  func() int
	swap func(i, j int)
	less func(i, j int) bool
}

// Takes an implementation of getter and returns sorter as high order func.
func Sorter[T constraints.Ordered](fncGet func(n netdata.Slice, i int) T) FncSorter {
	return func(n netdata.Slice) sort.Interface {
		return &defaultSorter{
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
func (s defaultSorter) Len() int           { return s.len() }
func (s defaultSorter) Swap(i, j int)      { s.swap(i, j) }
func (s defaultSorter) Less(i, j int) bool { return s.less(i, j) }

// Default sorter by network key (SSID, BSSID).
func ByKeySorter() FncSorter {
	return func(n netdata.Slice) sort.Interface {
		return &defaultSorter{
			len:  func() int { return len(n) },
			swap: func(i, j int) { n[i], n[j] = n[j], n[i] },
			less: func(i, j int) bool {
				return n[i].Key().Compare(n[j].Key()) < 0
			},
		}
	}
}
