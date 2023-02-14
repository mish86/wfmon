package wifitable

import (
	"context"
	"fmt"
	"sync"
	"wfmon/pkg/manuf" //nolint
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/wifi"
)

// Aggragated network data.
type NetworkData struct {
	BSSID        string
	Manuf        string // Short vendor' name
	ManufLong    string // Long vendor' name
	NetworkName  string
	Channel      uint8
	ChannelWidth uint16
	Band         wifi.Band // Bandwidth 2.4 or 5 Ghz
	RSSI         int8      // Received Signal Strength Indicator (RSSI)
	Quality      Quality   // Signal Quality
	Noise        int8      // Noise level
	SNR          int8      // Signal to Noise Ratio (SNR)
	// Seen
	// Rate
}

// Returns network data key.
func (data *NetworkData) Key() *NetworkKey {
	return NewNetworkKey(data.BSSID, data.NetworkName)
}

// Network data uniq key in table.
// Used for sorting in table.
type NetworkKey struct {
	BSSID       string
	NetworkName string
}

// Returns new network data key.
func NewNetworkKey(bssID, ssID string) *NetworkKey {
	return &NetworkKey{
		BSSID:       bssID,
		NetworkName: ssID,
	}
}

// Compares network data keys.
func (key *NetworkKey) Compare(other *NetworkKey) int {
	var res int
	switch {
	case len(key.NetworkName) == 0 && len(other.NetworkName) == 0:
		res = cmp.Compare(key.BSSID, other.BSSID)
	case len(key.NetworkName) == 0:
		res = 1
	case len(other.NetworkName) == 0:
		res = -1
	default:
		res = cmp.Compare(key.NetworkName, other.NetworkName)
		if res == 0 {
			res = cmp.Compare(key.BSSID, other.BSSID)
		}
	}

	return res
}

// Network data map.
type NetworkTable map[NetworkKey]*NetworkData

// Network data slice.
type NetworkSlice []NetworkData

// Returns slice of NetworkData copied from NetworkTable.
func (t NetworkTable) Slice() NetworkSlice {
	s := make(NetworkSlice, len(t))

	idx := 0
	for _, data := range t {
		s[idx] = *data
		idx++
	}

	return s
}

// Wraps networks table.
type DataSource struct {
	table     NetworkTable
	tableLock sync.RWMutex

	ctx      context.Context
	stop     context.CancelFunc
	framesCh <-chan wifi.Frame
}

// Returns new networks table.
func NewDataSource(framesCh <-chan wifi.Frame) *DataSource {
	const defaultInitTableSize = 20

	return &DataSource{
		table:    make(NetworkTable, defaultInitTableSize),
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

			data := NetworkDataConverter(frame).NetworkData()
			ds.Add(data)

		case <-ds.ctx.Done():
			return nil
		}
	}
}

// Appends or merges new data in networks table.
func (ds *DataSource) Add(data *NetworkData) {
	ds.tableLock.Lock()
	defer ds.tableLock.Unlock()

	key := data.Key()

	var thisData *NetworkData
	var found bool
	if thisData, found = ds.table[*key]; !found {
		// Copy data
		ds.table[*key] = data

		return
	}

	// merge network with existing
	thisData.Manuf = data.Manuf
	thisData.Channel = data.Channel
	thisData.ChannelWidth = data.ChannelWidth
	thisData.Band = data.Band
	thisData.RSSI = data.RSSI
	thisData.Noise = data.Noise
	thisData.SNR = data.SNR
}

// Returns network data slice.
func (ds *DataSource) NetworkSlice() NetworkSlice {
	ds.tableLock.RLock()
	defer ds.tableLock.RUnlock()

	return ds.table.Slice()
}

// Alias for network data converter.
type NetworkDataConverter wifi.Frame

// Converts wifi frame to network data.
func (frame NetworkDataConverter) NetworkData() *NetworkData {
	data := &NetworkData{
		BSSID:       frame.BSSID.String(),
		NetworkName: frame.SSID,
		Channel:     frame.Channel,
		RSSI:        frame.RSSI,
		Noise:       frame.Noise,
		SNR:         frame.RSSI - frame.Noise,
	}

	data.Manuf, data.ManufLong = manuf.Lookup(frame.BSSID.String())
	data.Quality = QualityConverter{data.RSSI, data.SNR}.SignalQuality()
	data.Band = wifi.GetBandByChan(frame.Channel)
	data.ChannelWidth = wifi.GetBondingWidth(frame.Channel, frame.SecondaryChannelOffset)

	return data
}

// Alias for quality field in network data.
type Quality uint8

// Quality converter based on RSSI and SNR values.
type QualityConverter struct {
	RSSI int8
	SNR  int8
}

// Determines signal quality in pecents (0-100%).
// Calculates signal quality based on RSSI using quadratic model,
// based on SNR using liner model and selects a returns value of them.
// ref. https://github.com/torvalds/linux/blob/master/drivers/net/wireless/intel/ipw2x00/ipw2200.c#L4304-L4317
// ref. https://gist.github.com/senseisimple/002cdba344de92748695a371cef0176a
func (c QualityConverter) SignalQuality() Quality {
	const (
		expAvgRSSI  = -60
		perfectRSSI = -20
		worstRSSI   = -85
	)

	rssi := int(c.RSSI)

	rssiQuality :=
		(100*
			(perfectRSSI-worstRSSI)*(perfectRSSI-worstRSSI) -
			(perfectRSSI-rssi)*(15*(perfectRSSI-worstRSSI)+62*(perfectRSSI-rssi))) /
			((perfectRSSI - worstRSSI) * (perfectRSSI - worstRSSI))

	//nolint:gomnd // ignore
	if rssiQuality > 100 {
		rssiQuality = 100
	} else if rssiQuality < 1 {
		rssiQuality = 0
	}

	snr := int(c.SNR)

	//nolint:gomnd // ignore
	snrQuality := func() int {
		switch {
		case snr < 0:
			return 0
		case 0 <= snr || snr < 40:
			return 5.0 * snr / 2.0
		case snr >= 40:
			return 100
		default:
			return 0
		}
	}()

	if rssiQuality < snrQuality {
		return Quality(rssiQuality)
	}

	return Quality(snrQuality)
}

// Returns percent presentation of signal quality.
func (q Quality) String() string {
	return fmt.Sprintf("%d%%", q)
}
