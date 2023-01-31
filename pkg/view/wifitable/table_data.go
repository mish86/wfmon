package wifitable

import (
	"sync"
	"wfmon/pkg/utils/cmp"
)

type NetworkData struct {
	BSSID        string
	NetworkName  string
	Channel      int
	ChannelWidth int
	Band         string // Bandwidth 2.4 or 5 Ghz
	RSSI         int8   // Received Signal Strength Indicator (RSSI)
	Noise        int8   // Noise level
	SNR          int8   // Signal to Noise Ratio (SNR)
	// Seen
	// Rate
}

func (data NetworkData) Key() NetworkTableKey {
	return NetworkTableKey{
		data.BSSID, data.NetworkName,
	}
}

type NetworkTableKey struct {
	BSSID       string
	NetworkName string
}

// Compares table keys.
func (key NetworkTableKey) Compare(other NetworkTableKey) int {
	var res int
	switch {
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

type NetworkTable map[NetworkTableKey]*NetworkData
type NetworkSlice []NetworkData

// Returns slice of NetworkData copied from NetworkTable.
func (t NetworkTable) ToSlice() NetworkSlice {
	s := make(NetworkSlice, len(t))

	idx := 0
	for _, data := range t {
		s[idx] = *data
		idx++
	}

	return s
}

type TableData struct {
	table     NetworkTable
	tableLock sync.RWMutex
}

func NewTableData() *TableData {
	return &TableData{
		table: make(NetworkTable),
	}
}

func (t *TableData) Add(data *NetworkData) {
	t.tableLock.Lock()
	defer t.tableLock.Unlock()

	key := NetworkTableKey{
		data.BSSID,
		data.NetworkName,
	}

	var thisData *NetworkData
	var found bool
	if thisData, found = t.table[key]; !found {
		// Copy data
		t.table[key] = data

		return
	}

	// merge network with existing
	thisData.Channel = data.Channel
	thisData.ChannelWidth = data.ChannelWidth
	thisData.Band = data.Band
	thisData.RSSI = data.RSSI
	thisData.Noise = data.Noise
	thisData.SNR = data.SNR
}

func (t *TableData) GetNetworkSlice() NetworkSlice {
	t.tableLock.RLock()
	defer t.tableLock.RUnlock()

	return t.table.ToSlice()
}
