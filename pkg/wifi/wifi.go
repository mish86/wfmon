package wifi

// https://mrncciew.com/2014/10/15/cwap-2-4ghz-vs-5ghz/
type Band uint8

const (
	Unknown Band = iota
	ISM          // 2.4GHz (industrial, scientific, and medical - ISM bands)
	UNII1        // 5 GHz (Unlicensed National Information Infrastructure – UNII bands)
	UNII2
	UNII3
)

func (b Band) String() string {
	return []string{Unknown: "", ISM: "2.4", UNII1: "5", UNII2: "5", UNII3: "5"}[b]
}

// Returns '2.4GHz' or '5GHz'.
// TODO: review and actualize bounds.
func GetBandByChan(channel uint8) Band {
	switch {
	case channel >= 1 && channel <= 14:
		return ISM
	case channel >= 32 && channel <= 48:
		return UNII1
	case channel >= 50 && channel <= 142:
		return UNII2
	case channel >= 142 && channel <= 177:
		return UNII3
	default:
		return Unknown
	}
}

// https://mrncciew.com/2014/11/04/cwap-ht-operations-ie/
type SecondaryChannelOffset uint8

const (
	SCN      SecondaryChannelOffset = 0 // no secondary channel is present
	SCA      SecondaryChannelOffset = 1 // secondary channel is above the primary channel
	Reserved SecondaryChannelOffset = 2 // reserved
	SCB      SecondaryChannelOffset = 3 // secondary channel is below the primary channel
)

func (o SecondaryChannelOffset) String() string {
	return []string{SCN: "SCN", SCA: "SCA", Reserved: "RSRVD", SCB: "SCB"}[o]
}

// Returns width in MHz.
func GetBondingWidth(frame Frame) uint16 {
	band := GetBandByChan(frame.Channel)
	// Unknown bandwidth
	if band == Unknown {
		return 0
	}

	// No bonding
	var bondingMultiplier uint16 = 1

	if offsetEnum := SecondaryChannelOffset(frame.SecondaryChannelOffset); offsetEnum == SCA || offsetEnum == SCB {
		// bonding
		bondingMultiplier = 2
	}

	// ISM bonding, 20 Mhz or 40 Mhz
	if band == ISM {
		return bondingMultiplier * getISMWidth(frame.Channel)
	}

	// UNII bonding, HT, 20 Mhz or 40 Mhz
	if frame.ChannelWidth == 0 {
		return bondingMultiplier * getUNIIWidth(frame.Channel)
	}

	// UNII bonding, VHT, 80 Mhz or 160 Mhz or 80 + 80 Mhz
	return getUNIIWidth(frame.ChannelCenterSegment0)
}

// ISM bonding width.
// Returns width in MHz.
func getISMWidth(channel uint8) uint16 {
	//nolint:gomnd // ignore
	return 20
}

// UNII bonding width.
// Returns width in MHz.
func getUNIIWidth(channel uint8) uint16 {
	//nolint:gomnd // ignore
	switch channel {
	case 50, 114, 163:
		return 160
	case 42, 58, 106, 122, 138, 155, 171:
		return 80
	case 34, 38, 46, 54, 62, 102, 110, 118, 126, 134, 142, 151, 159, 167, 175:
		return 40
	default:
		return 20
	}
}
