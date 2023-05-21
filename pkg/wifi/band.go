package wifi

// https://mrncciew.com/2014/10/15/cwap-2-4ghz-vs-5ghz/
type Band uint8

const (
	Unknown Band = iota
	ISM          // 2.4GHz (industrial, scientific, and medical - ISM bands)
	UNII1        // 5 GHz (Unlicensed National Information Infrastructure – UNII bands)
	UNII2A
	UNII2B
	UNII2C
	UNII3
)

const (
	MinBand = ISM
	MaxBand = UNII3
)

func (b Band) String() string {
	return []string{
		Unknown: "",
		ISM:     "ISM",
		UNII1:   "U-NII-1",
		UNII2A:  "U-NII-2A",
		UNII2B:  "U-NII-2B",
		UNII2C:  "U-NII-2C",
		UNII3:   "U-NII-3",
	}[b]
}

func (b Band) Range() string {
	return []string{Unknown: "", ISM: "2.4", UNII1: "5", UNII2A: "5", UNII2B: "5", UNII2C: "5", UNII3: "5"}[b]
}

// Returns '2.4GHz' or '5GHz'.
// TODO: review and actualize bounds.
func GetBandByChan(channel uint8) Band {
	switch {
	case channel >= 1 && channel <= 14:
		return ISM
	case channel >= 32 && channel <= 48:
		return UNII1
	case channel >= 50 && channel <= 68:
		return UNII2A
	case channel >= 96 && channel <= 144:
		return UNII2C
	case channel >= 149 && channel <= 173:
		return UNII3
	default:
		return Unknown
	}
}
