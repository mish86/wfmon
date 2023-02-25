package widgets

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/widgets/color"
)

// Event with currently selected network in wifi table.
// Sent by wifi table on cursor move.
// Handled by signal sparkline chart.
type NetworkKeyMsg struct {
	Key   netdata.Key
	Color color.HexColor
}

// Event indicates change measurement in signal sparkline chart.
type FieldMsg string

// Event with current wifi table width.
// Other widgets should be aligned with this widht.
type TableWidthMsg int
