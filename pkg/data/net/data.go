package netdata

import (
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/wifi"
)

const (
	SSIDKey      = "Network"
	BSSIDKey     = "BSSID"
	ManufKey     = "Manuf"
	ManufLongKey = "Manufactor"
	ChanKey      = "Chan"
	WidthKey     = "Width"
	BandKey      = "Band"
	RSSIKey      = "RSSI"
	QualityKey   = "Quality"
	BarsKey      = "Bars"
	NoiseKey     = "Noise"
	SNRKey       = "SNR"
)

// Aggragated network data.
type Network struct {
	BSSID            string                      // Station MAC address
	Manuf            string                      // Short vendor' name
	ManufLong        string                      // Long vendor' name
	NetworkName      string                      // SSID
	Channel          uint8                       // Primary channel number
	Offset           wifi.SecondaryChannelOffset // Secondary channel direction (2.5/5Ghz HT)
	FrequencyCenter0 uint8                       // Lower frequency segment center (5GHz VHT)
	FrequencyCenter1 uint8                       // Second frequency segment center (5GHz VHT)
	ChannelWidth     uint16                      // Channel width, MHz
	WidthOperation   wifi.ChannelWidthOperation  // Channel width operation (5GHz VHT)
	Band             wifi.Band                   // Bandwidth 2.4/5, Ghz
	RSSI             int8                        // Received Signal Strength Indicator (RSSI), dBm
	Quality          Quality                     // Signal Quality, %
	Noise            int8                        // Noise level, dBm
	SNR              int8                        // Signal to Noise Ratio (SNR), dBm
	// Seen
	// Rate
}

// Returns network data key.
func (data *Network) Key() Key {
	return NewKey(data.BSSID, data.NetworkName)
}

// Network data uniq key in table.
// Used for sorting in table.
type Key struct {
	BSSID       string
	NetworkName string
}

// Returns new network data key.
func NewKey(bssID, ssID string) Key {
	return Key{
		BSSID:       bssID,
		NetworkName: ssID,
	}
}

func Empty() Key {
	return NewKey("", "")
}

// Compares network data keys.
func (key Key) Compare(other Key) int {
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
