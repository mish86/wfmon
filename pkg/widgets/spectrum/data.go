package spectrum

import (
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/widgets/events"
	"wfmon/pkg/wifi"

	"github.com/charmbracelet/lipgloss"
)

const (
	wave20Mhz                          = 20 // wave width in Mhz
	wave20MhzWidth                     = 4  // number of channels in a wave of 20Mhz width
	halfOfWave80MhzWidthWithoutCenter  = 6  // number of channels in a wave of 80Mhz width excluding center segment
	halfOfWave160MhzWidthWithoutCenter = 14 // number of channels in a wave of 160Mhz width excluding center segment
	halfOfWave20MhzWidth               = 2  // number of channels in half of a wave of 20Mhz width
)

type Wave struct {
	Key            netdata.Key                // network key
	Band           wifi.Band                  // ISM or UNII
	Value          float64                    // RSSI or Quality
	Channel        uint8                      // primary channel
	Sign           int8                       // HT secondary channel location: +1 above / -1 below
	Center         [2]uint8                   // VHT frequency centers. lower is mandatory for VHT and second is mandatory for VHT and 80+80
	WidthOperation wifi.ChannelWidthOperation // VHT Channel Width Operaiton
	Width          uint8                      // 20Mhz channels count in a wave
	Color          lipgloss.Color             // spectrum color
}

func (wave *Wave) LowerChannel() uint8 {
	// if wave.SecondaryChannel == 0 {
	// 	return wave.PrimaryChannel
	// }

	// return cmp.Min(wave.PrimaryChannel, wave.SecondaryChannel)

	// HT
	if wave.Center[0] == 0 {
		return cmp.Min(wave.Channel, uint8(int8(wave.Channel)+wave.Sign*wave20MhzWidth*int8(wave.Width-1)))
	}

	// VHT
	switch wave.WidthOperation {
	case wifi.WidthOperation80, wifi.WidthOperation80And80:
		return cmp.Min(wave.Channel, wave.Center[0]-halfOfWave80MhzWidthWithoutCenter)
	case wifi.WidthOperation160:
		return cmp.Min(wave.Channel, wave.Center[0]-halfOfWave160MhzWidthWithoutCenter)
	default:
		return cmp.Min(wave.Channel, uint8(int8(wave.Channel)+wave.Sign*wave20MhzWidth*int8(wave.Width-1)))
	}
}

// func (wave *Wave) UpperChannel() uint8 {
// 	// return cmp.Max(wave.PrimaryChannel, wave.SecondaryChannel)
// 	return cmp.Max(wave.Channel, uint8(int8(wave.Channel)+wave.Sign*wave20MhzWidth*int8(wave.Width-1)))
// }

type Waver struct {
	netdata.Network

	fieldKey string
	ts       ds.TimeSeriesProvider
}

func (c Waver) Wave() Wave {
	net := netdata.Network(c.Network)

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

	val, _ := c.ts.TimeSeries(net.Key())(c.fieldKey).Last()

	return Wave{
		Key:            net.Key(),
		Band:           net.Band,
		Value:          val,
		Channel:        net.Channel,
		Sign:           sign,
		Width:          uint8(net.ChannelWidth / wave20Mhz),
		WidthOperation: net.WidthOperation,
		Center:         [2]uint8{c.FrequencyCenter0, c.FrequencyCenter1},
	}
}

type MultiWaver struct {
	events.NetworksOnScreen

	fieldKey string
	ts       ds.TimeSeriesProvider
}

func (c MultiWaver) Waves() []Wave {
	nets := c.Networks
	colors := c.Colors

	waves := make([]Wave, len(nets))

	for i := range nets {
		w := Waver{
			Network:  nets[i],
			fieldKey: c.fieldKey,
			ts:       c.ts,
		}.Wave()
		w.Color = colors[i].Lipgloss()
		waves[i] = w
	}

	return waves
}
