package wifi

type ChannelWidthOperation uint8

const (
	WidthOperation20Or40  ChannelWidthOperation = 0 // 20MHz or 40MHz
	WidthOperation80      ChannelWidthOperation = 1 // 80MHz
	WidthOperation160     ChannelWidthOperation = 2 // 160MHz
	WidthOperation80And80 ChannelWidthOperation = 3 // 80+80MHz
)

func (o ChannelWidthOperation) String() string {
	return []string{
		WidthOperation20Or40:  "20/40",
		WidthOperation80:      "80",
		WidthOperation160:     "160",
		WidthOperation80And80: "80+80",
	}[o]
}

func (o ChannelWidthOperation) Width() uint16 {
	//nolint:gomnd // ignore
	return []uint16{
		WidthOperation20Or40:  20,
		WidthOperation80:      80,
		WidthOperation160:     160,
		WidthOperation80And80: 160,
	}[o]
}

// 0 - 20MHz or 40MHz; 1 - 80MHz; 2 - 160MHz; 3 - 80+80MHz; others - reserved.
func GetChannelWidthOperation(w uint8) ChannelWidthOperation {
	knownOperations := []ChannelWidthOperation{
		WidthOperation20Or40,
		WidthOperation80,
		WidthOperation160,
		WidthOperation80And80,
	}

	if int(w) >= len(knownOperations) {
		return WidthOperation20Or40
	}

	return knownOperations[w]
}
