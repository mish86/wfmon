package events

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/widgets/color"
)

// Event with currently highlighted network in wifi table.
// Sent by wifi table on refresh, sort, (filter).
type NetworkKeyMsg struct {
	Key   netdata.Key
	Color color.HexColor
}

// Event with currently highlited network in wifi table.
// Sent by wifi table on cursor move.
type SelectedNetworkKeyMsg NetworkKeyMsg

// Event with toggled network in wifi table.
// Sent by wifi table on toggle a row.
type ToggledNetworkKeyMsg NetworkKeyMsg

// Event with type of signal measurement (RSSI, Quality, Bars) selected in wifi table.
type SignalFieldMsg struct {
	Key            string
	MinVal, MaxVal float64
}

// Event with current wifi table width.
// Sent by wifi table.
// Other widgets should be aligned with this width.
type TableWidthMsg int

// Event with networks currently displayed in wifi table.
// Sent by wifi table on pagination.
type NetworksOnScreen struct {
	Networks []netdata.Network
	Colors   []color.HexColor
}
