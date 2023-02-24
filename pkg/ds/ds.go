package ds

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ts"
)

type NetworkProvider interface {
	Networks() netdata.Slice
}

type TimeSeriesProvider interface {
	TimeSeries(netKey netdata.Key) func(colKey string) ts.TimeSeries
}

type Provider interface {
	NetworkProvider
	TimeSeriesProvider
}
