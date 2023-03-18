package sort

import (
	"sort"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/utils/cmp"

	"golang.org/x/exp/constraints"
)

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
