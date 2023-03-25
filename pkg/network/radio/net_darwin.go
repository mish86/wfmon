//go:build darwin

package radionet

import (
	"wfmon/pkg/network"
	"wfmon/pkg/network/radio/darwin/airport"
	"wfmon/pkg/network/radio/darwin/corewlan"
)

// Returns default WiFi interface name.
func GetDefaultWiFiInterface() (string, error) {
	return corewlan.GetDefaultWiFiInterface()
}

// Returns network associated with the given interface.
// Should be invoked before setting interface in monitoring mode.
func GetAssociatedNetwork(ifaceName string) (network.Network, error) {
	return airport.GetAssociatedNetwork(ifaceName)
}

// Disconnects interface from network.
// Should be invoked before setting interface in monitoring mode.
func DisassociateFromNetwork(ifaceName string) error {
	return airport.DisassociateFromNetwork(ifaceName)
}

// Returns channels supported by given interface.
func GetSupportedChannels(ifaceName string) ([]int, error) {
	return corewlan.GetSupportedChannels(ifaceName)
}

// Change radio channel on given interface.
func SetInterfaceChannel(ifaceName string, channel int) error {
	return corewlan.SetInterfaceChannel(ifaceName, channel)
}
