package wifi

// Returns width in MHz.
func GetChannelWidth(frame Frame) uint16 {
	band := GetBandByChan(frame.Channel)
	// Unknown bandwidth
	if band == Unknown {
		return 0
	}

	// No bonding
	var bonding uint16 = 1

	if offsetEnum := SecondaryChannelOffset(frame.SecondaryChannelOffset); offsetEnum == SCA || offsetEnum == SCB {
		// bonding
		bonding = 2
	}

	// ISM HT bonding, 20 Mhz or 40 Mhz
	if band == ISM {
		return bonding * getISMWidth(frame.Channel)
	}

	vhtChannelWidthOperation := GetChannelWidthOperation(frame.ChannelWidth)

	// UNII VHT bonding, 20 Mhz or 40 Mhz
	if vhtChannelWidthOperation == WidthOperation20Or40 {
		// return bonding * getUNIIWidth(frame.Channel)
		return bonding * 20
	}

	// UNII bonding, VHT, 80 Mhz or 160 Mhz or 80 + 80 Mhz
	// return getUNIIWidth(frame.ChannelCenterSegment0)
	return vhtChannelWidthOperation.Width()
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
