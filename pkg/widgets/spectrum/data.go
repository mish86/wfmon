package spectrum

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/widgets/events"
	"wfmon/pkg/wifi"

	"github.com/charmbracelet/lipgloss"
)

const (
	wave20Mhz            = 20 // wave width in Mhz
	wave20MhzWidth       = 4  // number of channels in a wave of 20Mhz width
	halfOfWave20MhzWidth = 2  // number of channels in half of a wave of 20Mhz width
)

type Wave struct {
	Key            netdata.Key    // network key
	Band           wifi.Band      // ISM or UNII
	Value          float64        // RSSI or Quality
	PrimaryChannel uint8          // primary channel
	Sign           int8           // secondary channel: +1 above / -1 below
	Width          uint8          // channels number in a wave
	Color          lipgloss.Color // spectrum color
}

func (wave *Wave) LowerChannel() uint8 {
	// if wave.SecondaryChannel == 0 {
	// 	return wave.PrimaryChannel
	// }

	// return cmp.Min(wave.PrimaryChannel, wave.SecondaryChannel)
	return cmp.Min(wave.PrimaryChannel, uint8(int8(wave.PrimaryChannel)+wave.Sign*wave20MhzWidth*int8(wave.Width-1)))
}

func (wave *Wave) UpperChannel() uint8 {
	// return cmp.Max(wave.PrimaryChannel, wave.SecondaryChannel)
	return cmp.Max(wave.PrimaryChannel, uint8(int8(wave.PrimaryChannel)+wave.Sign*wave20MhzWidth*int8(wave.Width-1)))
}

type Waver netdata.Network

func (c Waver) Wave() Wave {
	net := netdata.Network(c)

	var sign int8
	//nolint:exhaustive // ignore
	switch net.Offset {
	case wifi.SCA:
		sign = 1
	case wifi.SCB:
		sign = -1
	default:
		sign = 0
	}

	return Wave{
		Key:            net.Key(),
		Band:           net.Band,
		Value:          float64(net.RSSI),
		PrimaryChannel: net.Channel,
		Sign:           sign,
		Width:          uint8(net.ChannelWidth / wave20Mhz),
	}
}

type MultiWaver events.NetworksOnScreen

func (c MultiWaver) Waves() []Wave {
	nets := c.Networks
	colors := c.Colors

	waves := make([]Wave, len(nets))

	for i := range nets {
		w := Waver(nets[i]).Wave()
		w.Color = colors[i].Lipgloss()
		waves[i] = w
	}

	return waves
}
