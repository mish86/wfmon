//go:build !windows | !darwin

package radionet

// Disconnects interface from network. Required before for setting interface in monitoring mode.
func DisassociateFromNetwork(ifaceName string) error {
	log.Error("Unimplemented")
}

// Returns channels supported by given interface.
func GetSupportedChannels(ifaceName string) ([]int, error) {
	log.Error("Unimplemented")
}

func GetAssociatedNetwork(ifaceName string) (Network, error) {
	log.Error("Unimplemented")
}

// Change radio channel on given interface.
func SetInterfaceChannel(ifaceName string, channel int) error {
	log.Error("Unimplemented")
}
