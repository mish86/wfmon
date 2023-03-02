package order

// Sorting order enum.
type Dir uint8

const (
	None Dir = iota
	ASC
	DESC
)

// Returns string presentation of sorting order.
// Used in column title view.
func (o Dir) String() string {
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
func (o Dir) Swap() Dir {
	if o == ASC {
		return DESC
	}

	return ASC
}
