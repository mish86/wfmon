package wifi

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
