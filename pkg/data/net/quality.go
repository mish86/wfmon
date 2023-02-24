package netdata

import "fmt"

// Alias for quality field in network data.
type Quality uint8

// Quality converter based on RSSI and SNR values.
type QualityConverter struct {
	RSSI int8
	SNR  int8
}

// Determines signal quality in pecents (0-100%).
func (c QualityConverter) SignalQuality() Quality {
	// rssiQuality := c.quadRSSI()
	// snrQuality := c.linerSNR()

	// if rssiQuality < snrQuality {
	// 	return rssiQuality
	// }

	// return snrQuality

	return c.quadRSSI()
}

// Calculates signal quality based on RSSI using quadratic model.
// ref. https://github.com/torvalds/linux/blob/master/drivers/net/wireless/intel/ipw2x00/ipw2200.c#L4304-L4317
func (c QualityConverter) quadRSSI() Quality {
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

	return Quality(rssiQuality)
}

// Calculates signal quality based on SNR using liner model.
// ref. https://gist.github.com/senseisimple/002cdba344de92748695a371cef0176a
func (c QualityConverter) linerSNR() Quality {
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

	return Quality(snrQuality)
}

// Returns percent presentation of signal quality.
func (q Quality) String() string {
	return fmt.Sprintf("%d%%", q)
}
