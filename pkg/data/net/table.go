package netdata

// Network data map.
type Table map[Key]*Network

// Network data slice.
type Slice []Network

// Returns slice of NetworkData copied from NetworkTable.
func (t Table) Slice() Slice {
	s := make(Slice, len(t))

	idx := 0
	for _, data := range t {
		s[idx] = *data
		idx++
	}

	return s
}
