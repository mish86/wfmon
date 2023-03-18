package events

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

// Event signal strength measurement (RSSI, Quality, Bars) selected in wifi table.
// TODO: Sent by wifi table.
// TODO: Handled by signal sparkline chart and spectrum chart.
type FieldMsg string

// Event with current wifi table width.
// Sent by wifi table.
// Other widgets should be aligned with this width.
type TableWidthMsg int

// Event with networks currently displayed in wifi table.
// Sent by wifi table on pagination change.
// Handled by spectrum chart.
type NetworksOnScreen struct {
	Networks []netdata.Network
	Colors   []color.HexColor
}
