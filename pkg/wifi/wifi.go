package wifi

type Band int

// https://mrncciew.com/2014/10/15/cwap-2-4ghz-vs-5ghz/
const (
	Unknown Band = iota
	ISM          // 2.4GHz (industrial, scientific, and medical - ISM bands)
	UNII         // 5 GHz (Unlicensed National Information Infrastructure – UNII bands)
)

func (b Band) String() string {
	return []string{Unknown: "", ISM: "2.4", UNII: "5"}[b]
}

// Returns '2.4GHz' or '5GHz'.
func GetBandByChan(channel int) Band {
	switch {
	case channel >= 1 && channel <= 14:
		return ISM
	case channel >= 36 && channel <= 64:
		return UNII
	case channel >= 100 && channel <= 140:
		return UNII
	default:
		return Unknown
	}
}
