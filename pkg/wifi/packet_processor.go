package wifi

// https://mrncciew.com/2014/10/04/my-cwap-study-notes/

import (
	"net"
	"regexp"
	log "wfmon/pkg/logger"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// Wrapper for gopacket.Packet.
type PacketDiscover struct {
	gopacket.Packet
}

// Wraps gopacket.Packet into PacketDiscover.
func FromPacket(packet gopacket.Packet) *PacketDiscover {
	return &PacketDiscover{Packet: packet}
}

// Management layers constraint for tryLayer func.
type mgmtLayers = interface {
	layers.Dot11MgmtBeacon |
		layers.Dot11MgmtProbeResp |
		layers.Dot11MgmtAssociationResp |
		layers.Dot11MgmtReassociationResp
}

// Overall supported layers constraint for tryLayer func.
type supportedLayers = interface {
	layers.RadioTap | layers.Dot11 | mgmtLayers
}

// Tries to extract required layer type from packet and cast it to gopacket.Layer structure.
func tryLayer[T supportedLayers](p *PacketDiscover, layerType gopacket.LayerType) (*T, bool) {
	layer := p.Layer(layerType)
	if layer == nil {
		return nil, false
	}

	casted, ok := any(layer).(*T)
	if !ok {
		log.Warnf("failed to cast layer to %s in packet %v", layerType, p.Metadata().Timestamp)
		return nil, false
	}

	return casted, true
}

// Discoveres radio frame info from packet.
func (p *PacketDiscover) DiscoverRadioFrame() *RadioFrame {
	radio, ok := tryLayer[layers.RadioTap](p, layers.LayerTypeRadioTap)
	if !ok {
		return nil
	}

	// phyType := radio.VHT.Bandwidth
	// Channel Frequency
	freq := int(radio.ChannelFrequency)
	// Received Signal Strength Indicator (RSSI)
	rssi := radio.DBMAntennaSignal
	// Noise level
	noise := radio.DBMAntennaNoise

	return &RadioFrame{
		Frequency: freq,
		RSSI:      rssi,
		Noise:     noise,
	}
}

// Discovers wifi and radio frame information from packet.
// https://mrncciew.com/2014/09/28/cwap-mac-headeraddresses/
func (p *PacketDiscover) DiscoverDot11Frame() *Dot11Frame {
	dot11, ok := tryLayer[layers.Dot11](p, layers.LayerTypeDot11)
	if !ok {
		return nil
	}

	var frame *Dot11Frame

	switch {
	// ToDS == 0 and FromDS == 0
	case !dot11.Flags.ToDS() && !dot11.Flags.FromDS():
		frame = NewDot11Frame(dot11.Type,
			dot11.Address2, dot11.Address1, dot11.Address2, dot11.Address1, dot11.Address3)
	// ToDS == 0 and FromDS == 1
	case !dot11.Flags.ToDS() && dot11.Flags.FromDS():
		frame = NewDot11Frame(dot11.Type,
			dot11.Address3, dot11.Address1, net.HardwareAddr{}, dot11.Address1, dot11.Address1)
	// ToDS == 1 and FromDS == 0
	case dot11.Flags.ToDS() && !dot11.Flags.FromDS():
		frame = NewDot11Frame(dot11.Type,
			dot11.Address2, dot11.Address3, dot11.Address2, dot11.Address1, dot11.Address1)
	// ToDS == 1 and FromDS == 1
	case dot11.Flags.ToDS() && dot11.Flags.FromDS():
		frame = NewDot11Frame(dot11.Type,
			dot11.Address4, dot11.Address3, dot11.Address2, dot11.Address1, net.HardwareAddr{})
	}

	if frame == nil {
		return frame
	}

	if radio := p.DiscoverRadioFrame(); radio != nil {
		frame.RadioFrame = *radio
	}

	return frame
}

// Discovers Information Elements from packet.
func (p *PacketDiscover) DiscoverIEs() *InformationElements {
	var ie *InformationElements

	for _, layer := range p.Layers() {
		if layer.LayerType() != layers.LayerTypeDot11InformationElement {
			continue
		}

		var ok bool
		dot11info, ok := layer.(*layers.Dot11InformationElement)
		if !ok {
			continue
		}

		//nolint:exhaustive // process only 3 IE
		switch dot11info.ID {
		// case layers.Dot11InformationElementIDSSID:
		// 	if ie == nil {
		// 		ie = &InformationElements{}
		// 	}
		// 	ie.discoverSSIDIE(dot11info)

		// Operation Elements can be discovered from
		// Beacon, Reassociation Response & Probe Response frames transmitted by an AP.
		// https://mrncciew.com/2014/11/04/cwap-ht-operations-ie/
		case layers.Dot11InformationElementIDHTInfo:
			if ie == nil {
				ie = &InformationElements{}
			}
			ie.discoverHTIE(dot11info)

		case layers.Dot11InformationElementIDVHTOperation:
			if ie == nil {
				ie = &InformationElements{}
			}
			ie.discoverVHTIE(dot11info)

		case layers.Dot11InformationElementIDDSSet:
			if ie == nil {
				ie = &InformationElements{}
			}
			ie.discoverDSSetIE(dot11info)
		}
	}

	if ie != nil && ie.Channel == 0 && ie.PrimaryChannel != 0 {
		ie.Channel = ie.PrimaryChannel
	}

	return ie
}

// Discovers SSID from Information Element.
// func (ie *InformationElements) discoverSSIDIE(dot11info *layers.Dot11InformationElement) {
// 	if len(dot11info.Info) > 0 {
// 		ie.SSIDIE = SSIDIE{
// 			SSID: string(dot11info.Info),
// 		}
// 	}
// }

// Discovers HT Operation from Information Element.
func (ie *InformationElements) discoverHTIE(dot11info *layers.Dot11InformationElement) {
	// check malformed packet
	if len(dot11info.Contents) >= 4 {
		ie.HTOperationIE = HTOperationIE{
			PrimaryChannel:         dot11info.Contents[2],
			SecondaryChannelOffset: dot11info.Contents[3] & 0b00000011,        //nolint:gomnd // ignore
			SupportedChannelWidth:  (dot11info.Contents[3] & 0b00000100) >> 2, //nolint:gomnd // ignore
		}
	}
}

// Discovers VHT Operation from Information Element.
func (ie *InformationElements) discoverVHTIE(dot11info *layers.Dot11InformationElement) {
	// check malformed packet
	if len(dot11info.Contents) >= 5 {
		ie.VHTOperationIE = VHTOperationIE{
			ChannelWidth:          dot11info.Contents[2],
			ChannelCenterSegment0: dot11info.Contents[3],
			ChannelCenterSegment1: dot11info.Contents[4],
		}
	}
}

// Discovers DS Set from Information Element.
func (ie *InformationElements) discoverDSSetIE(dot11info *layers.Dot11InformationElement) {
	// check malformed packet
	if len(dot11info.Info) >= 1 {
		ie.DSSetIE = DSSetIE{
			Channel: dot11info.Info[0],
		}
	}
}

// Discovers wifi frame.
// Then traverses packet and discovers a management frame that contains Information Elements.
func (p *PacketDiscover) DiscoverMgmtFrame() *MgmtFrame {
	var dot11 *Dot11Frame
	if dot11 = p.DiscoverDot11Frame(); dot11 == nil {
		// skip non wifi package
		return nil
	}

	for _, discover := range []func() *MgmtFrame{
		p.DiscoverMgmtBeaconFrame,
		p.DiscoverMgmtProbeRespFrame,
		p.DiscoverMgmtAssociationRespFrame,
		p.DiscoverMgmtReassociationRespFrame,
	} {
		if frame := discover(); frame != nil {
			frame.Dot11Frame = *dot11

			if ie := p.DiscoverIEs(); ie != nil {
				frame.InformationElements = *ie
				// if len(ie.SSID) > 0 {
				// frame.SSID = ie.SSID
				// }
			}

			return frame
		}
	}

	return nil
}

// Discovers Management Beacon frame from packet.
// https://mrncciew.com/2014/10/08/802-11-mgmt-beacon-frame/
func (p *PacketDiscover) DiscoverMgmtBeaconFrame() *MgmtFrame {
	beacon, ok := tryLayer[layers.Dot11MgmtBeacon](p, layers.LayerTypeDot11MgmtBeacon)
	if !ok {
		return nil
	}

	// malformed packet
	if len(beacon.BaseLayer.Contents) < 14 {
		return nil
	}
	ssIDLen := int(beacon.BaseLayer.Contents[13])
	// malformed packet
	if ssIDLen > len(beacon.BaseLayer.Contents)-14 {
		return nil
	}
	ssID := string(beacon.BaseLayer.Contents[14 : 14+ssIDLen])
	re := regexp.MustCompile(`[[:cntrl:]]`)
	ssID = re.ReplaceAllString(ssID, "?")
	frame := &MgmtFrame{
		SSID: ssID,
	}

	return frame
}

// Discovers Management Probe Response frame from packet.
// https://mrncciew.com/2014/10/27/cwap-802-11-probe-requestresponse/
func (p *PacketDiscover) DiscoverMgmtProbeRespFrame() *MgmtFrame {
	resp, ok := tryLayer[layers.Dot11MgmtProbeResp](p, layers.LayerTypeDot11MgmtProbeResp)
	if !ok {
		return nil
	}

	// malformed packet
	if len(resp.BaseLayer.Contents) < 14 {
		return nil
	}
	ssIDLen := int(resp.BaseLayer.Contents[13])
	// malformed packet
	if ssIDLen > len(resp.BaseLayer.Contents)-14 {
		return nil
	}

	ssIDLen = min(ssIDLen, len(resp.BaseLayer.Contents)-14)
	ssID := string(resp.BaseLayer.Contents[14 : 14+ssIDLen])
	re := regexp.MustCompile(`[[:cntrl:]]`)
	ssID = re.ReplaceAllString(ssID, "?")
	frame := &MgmtFrame{
		SSID: ssID,
	}

	return frame
}

// Discovers Management Association Response frame from packet.
// https://mrncciew.com/2014/10/28/802-11-mgmt-association-reqresponse/
func (p *PacketDiscover) DiscoverMgmtAssociationRespFrame() *MgmtFrame {
	_, ok := tryLayer[layers.Dot11MgmtAssociationResp](p, layers.LayerTypeDot11MgmtAssociationResp)
	if !ok {
		return nil
	}

	return &MgmtFrame{}
}

// Discovers Management Reassociation Response frame from packet.
// https://mrncciew.com/2014/10/28/cwap-reassociation-reqresponse/
func (p *PacketDiscover) DiscoverMgmtReassociationRespFrame() *MgmtFrame {
	_, ok := tryLayer[layers.Dot11MgmtReassociationResp](p, layers.LayerTypeDot11MgmtReassociationResp)
	if !ok {
		return nil
	}

	return &MgmtFrame{}
}
