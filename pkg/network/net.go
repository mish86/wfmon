package network

import (
	"fmt"
	"net"
)

type Network struct {
	SSID    string
	BSSID   string
	Channel uint8 // optional
}

func InterfaceByName(name string) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range ifaces {
		if iface.Name == name {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("no interface '%s' found", name)
}
