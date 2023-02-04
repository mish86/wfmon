package wifi

import (
	"fmt"
	"net"

	"github.com/google/gopacket/layers"
)

// Radio Frame.
type RadioFrame struct {
	Frequency int  // Channel Frequency
	RSSI      int8 // Received Signal Strength Indicator (RSSI)
	Noise     int8 // Noise level
}

func (f *RadioFrame) String() string {
	return fmt.Sprintf("Frequency:%d RSSI:%d Noise:%d",
		f.Frequency,
		f.RSSI,
		f.RSSI,
	)
}

// Wifi Frame.
type Dot11Frame struct {
	RadioFrame

	Dot11Type          layers.Dot11Type
	SourceAddress      net.HardwareAddr
	DestinationAddress net.HardwareAddr
	TransmitterAddress net.HardwareAddr
	ReceiverAddress    net.HardwareAddr
	BSSID              net.HardwareAddr // L2 ID of BSS (Basic Service Set)
}

// Prints only Radio frame, frame type, Source and Destination MAC addresses and BSSID.
func (f *Dot11Frame) String() string {
	return fmt.Sprintf("Radio:%+v, Dot11:%v Src:%s Dst:%s BSSID:%s",
		f.RadioFrame,
		f.Dot11Type,
		f.SourceAddress,
		f.DestinationAddress,
		f.BSSID,
	)
}

// Creates Dot11Frame with given parameters in the order:
// Dot11Type, Source, Destination, Transmitter, Receiver, BSSID.
// Use net.HardwareAddr{} for empty address.
func NewDot11Frame(dot11Type layers.Dot11Type, addrs ...net.HardwareAddr) *Dot11Frame {
	const addrNum = 5
	if len(addrs) < addrNum {
		addrs = append(addrs, make(net.HardwareAddr, addrNum-len(addrs)))
	}

	return &Dot11Frame{
		Dot11Type:          dot11Type,
		SourceAddress:      addrs[0],
		DestinationAddress: addrs[1],
		TransmitterAddress: addrs[2],
		ReceiverAddress:    addrs[3],
		BSSID:              addrs[4],
	}
}

// High Throughput Operations Information Element.
type HTOperationsIE struct {
	PrimaryChannel         uint8
	SecondaryChannelOffset uint8
	SupportedChannelWidth  uint8
}

type DSSetIE struct {
	Channel uint8
}

type SSIDIE struct {
	SSID string
}

type InformationElements struct {
	HTOperationsIE // optional
	DSSetIE        // optional
	// SSIDIE         // optional
}

func (ie *InformationElements) String() string {
	// return fmt.Sprintf("HT:%+v DS:%+v SSID:%+v", ie.HTOperationsIE, ie.DSSetIE, ie.SSIDIE)
	return fmt.Sprintf("HT:%+v DS:%+v", ie.HTOperationsIE, ie.DSSetIE)
}

// Management frame.
type MgmtFrame struct {
	Dot11Frame
	InformationElements
	SSID string // optional
}

func (f *MgmtFrame) String() string {
	return fmt.Sprintf("Dot11:%+v, SSID:%s IE:%+v", f.Dot11Frame, f.SSID, f.InformationElements)
}

// Generic frame.
type Frame MgmtFrame
