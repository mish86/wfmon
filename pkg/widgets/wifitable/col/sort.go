package column

import (
	"sort"
	netdata "wfmon/pkg/data/net"
	s "wfmon/pkg/widgets/sort"
	order "wfmon/pkg/widgets/wifitable/ord"
)

// Describes sorting order for a column and sorting function.
type Sort struct {
	key    string
	ord    order.Dir
	sorter s.FncSorter
}

// Returns new Sort object.
func NewSort(key string, sorter s.FncSorter) Sort {
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
