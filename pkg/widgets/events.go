package widgets

import netdata "wfmon/pkg/data/net"

// Event with currently selected network in wifi table.
// Sent by wifi table on cursor move.
// Handled by signal sparkline chart.
type NetworkKeyMsg netdata.Key

// Event indicates change measurement in signal sparkline chart.
type FieldMsg string

// Event with current wifi table width.
// Other widgets should be aligned with this widht.
type TableWidthMsg int
