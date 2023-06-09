package ds

import (
	"context"
	"fmt"
	"sync"
	"time"

	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/manuf" //nolint
	"wfmon/pkg/ts"
	"wfmon/pkg/wifi"
)

const (
	defaultTimeSeriesSize = 200
)

// Wraps networks table.
type DataSource struct {
	table     netdata.Table
	tableLock sync.RWMutex

	ts     map[netdata.Key]map[string]ts.TimeSeries
	tsLock sync.RWMutex

	ctx      context.Context
	stop     context.CancelFunc
	framesCh <-chan wifi.Frame
}

// Returns new networks table.
func New(framesCh <-chan wifi.Frame) *DataSource {
	const defaultInitTableSize = 20

	return &DataSource{
		table:    make(netdata.Table, defaultInitTableSize),
		ts:       make(map[netdata.Key]map[string]ts.TimeSeries),
		framesCh: framesCh,
	}
}

// Starts processing incomming frames from packets.
func (ds *DataSource) Start(ctx context.Context) error {
	ds.ctx, ds.stop = context.WithCancel(ctx)

	for {
		select {
		case frame, ok := <-ds.framesCh:
			if !ok {
				return fmt.Errorf("frames source closed, stopping updating table")
			}

			ds.Add(frameConverter(frame).Network())

		case <-ds.ctx.Done():
			return nil
		}
	}
}

// Appends or merges new data in networks table.
func (ds *DataSource) Add(newData *netdata.Network) {
	ds.tableLock.Lock()
	defer ds.tableLock.Unlock()

	addMetric := func(netKey netdata.Key, fieldKey string, val float64, timestamp time.Time) {
		ds.tsLock.Lock()
		defer ds.tsLock.Unlock()

		if _, found := ds.ts[netKey]; !found {
			ds.ts[netKey] = map[string]ts.TimeSeries{
				fieldKey: ts.New(defaultTimeSeriesSize).Add(val, timestamp),
			}

			return
		}

		ds.ts[netKey][fieldKey] = ds.ts[netKey][fieldKey].Add(val, timestamp)
	}

	key := newData.Key()

	var entry *netdata.Network
	var found bool
	if entry, found = ds.table[key]; !found {
		// Copy data
		ds.table[key] = newData

		// new timeseries
		addMetric(key, netdata.RSSIKey, float64(newData.RSSI), time.Now())
		addMetric(key, netdata.QualityKey, float64(newData.Quality), time.Now())

		return
	}

	// merge network with existing
	{
		entry = &*newData
		// TODO: use avarage for quality and noise
		ds.table[key] = entry

		// append timeseries
		addMetric(key, netdata.RSSIKey, float64(newData.RSSI), time.Now())
		addMetric(key, netdata.QualityKey, float64(newData.Quality), time.Now())

		return
	}
}

// Returns network data slice.
func (ds *DataSource) Networks() netdata.Slice {
	ds.tableLock.RLock()
	defer ds.tableLock.RUnlock()

	return ds.table.Slice()
}

func (ds *DataSource) TimeSeries(netKey netdata.Key) func(colKey string) ts.TimeSeries {
	ds.tsLock.RLock()
	defer ds.tsLock.RUnlock()

	if timeSeries, found := ds.ts[netKey]; found {
		copied := make(map[string]ts.TimeSeries, len(timeSeries))
		for key, ts := range timeSeries {
			copied[key] = ts
		}
		return func(colKey string) ts.TimeSeries {
			return copied[colKey].Copy()
		}
	}

	return func(colKey string) ts.TimeSeries {
		return ts.Empty()
	}
}

// Alias for network data converter.
type frameConverter wifi.Frame

// Converts wifi frame to network netdata.
func (frame frameConverter) Network() *netdata.Network {
	entry := &netdata.Network{
		BSSID:            frame.BSSID.String(),
		NetworkName:      frame.SSID,
		Channel:          frame.Channel,
		FrequencyCenter0: frame.ChannelCenterSegment0,
		FrequencyCenter1: frame.ChannelCenterSegment1,
		Offset:           wifi.SecondaryChannelOffset(frame.SecondaryChannelOffset),
		RSSI:             frame.RSSI,
		Noise:            frame.Noise,
		SNR:              frame.RSSI - frame.Noise,
	}

	entry.Manuf, entry.ManufLong = manuf.Lookup(frame.BSSID.String())
	entry.Quality = netdata.QualityConverter{
		RSSI: entry.RSSI,
		SNR:  entry.SNR,
	}.SignalQuality()
	entry.Band = wifi.GetBandByChan(frame.Channel)
	entry.WidthOperation = wifi.GetChannelWidthOperation(frame.ChannelWidth)
	entry.ChannelWidth = wifi.GetChannelWidth(wifi.Frame(frame))
	entry.WidthOperation = wifi.GetChannelWidthOperation(frame.ChannelWidth)

	return entry
}
