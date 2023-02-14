//go:build darwin

package network

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

const airPortPath = "/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport"

var (
	channelFrequencyRegex = regexp.MustCompile(`^(\d{1,3})\s\([2,5]GHz\)$`)
	bssIDRegex            = regexp.MustCompile(`BSSID:\s((?:[A-Fa-f0-9]{2}[\.:\-]){5}[A-Fa-f0-9]{2})`)
	ssIDRegex             = regexp.MustCompile(`\s+SSID:\s(.*)`)
	channelRegex          = regexp.MustCompile(`\s+channel:\s(\d{1,3})`)
)

const (
	numGroups2 = 2
)

type systemProfiler struct {
	SPAirPortDataType []spAirportDataType `json:"SPAirPortDataType"`
}

type spAirportDataType struct {
	IFaces []spAirportAirportInterface `json:"spairport_airport_interfaces"`
}

type spNetworkInformation struct {
	Name        string `json:"_name"`
	Channel     string `json:"spairport_network_channel"`
	CountryCode string `json:"spairport_network_country_code"`
	MCS         int    `json:"spairport_network_mcs"`
	PHYMode     string `json:"spairport_network_phymode"`
	Rate        int    `json:"spairport_network_rate"`
	Security    string `json:"spairport_security_mode"`
}

type spAirportAirportInterface struct {
	Name      string               `json:"_name"`
	Channel   int                  `json:"spairport_airdrop_channel"`
	Supported []string             `json:"spairport_supported_channels"`
	Network   spNetworkInformation `json:"spairport_current_network_information"`
}

// Disconnects interface from network. Required before for setting interface in monitoring mode.
func DisassociateFromNetwork(ifaceName string) error {
	return exec.Command(airPortPath, ifaceName, "-z").Run()
}

func getSPAirportAirportInterface(ifaceName string) (*spAirportAirportInterface, error) {
	var (
		err error
		raw []byte
	)

	if raw, err = exec.Command("system_profiler", "SPAirPortDataType", "-json").Output(); err != nil {
		return nil, err
	}

	var sysInfo systemProfiler
	if err = json.Unmarshal(raw, &sysInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal system_profiler output: %w", err)
	}

	if len(sysInfo.SPAirPortDataType) == 0 {
		return nil, fmt.Errorf("no SPAirPortDataType in system_profiler output")
	}

	for _, iface := range sysInfo.SPAirPortDataType[0].IFaces {
		if iface.Name == ifaceName {
			return &iface, nil
		}
	}

	return nil, fmt.Errorf("no spairport_airport_interfaces found in system_profiler output by '%s'", ifaceName)
}

// Returns channels supported by given interface.
func GetSupportedChannels(ifaceName string) ([]int, error) {
	iface, err := getSPAirportAirportInterface(ifaceName)
	if err != nil {
		return []int{}, err
	}

	if len(iface.Supported) == 0 {
		return []int{}, fmt.Errorf("no channels supported by %s", ifaceName)
	}

	channels := make([]int, len(iface.Supported))
	for i, chFreq := range iface.Supported {
		chA := channelFrequencyRegex.FindStringSubmatch(chFreq)
		if len(chA) != numGroups2 {
			return nil, fmt.Errorf("failed to parse frequency '%s', expected pattern'%s'", chFreq, channelFrequencyRegex)
		}

		var chI int
		if chI, err = strconv.Atoi(chA[1]); err != nil {
			return nil, fmt.Errorf("failed to parse frequency '%s': %w", chFreq, err)
		}

		channels[i] = chI
	}

	return channels, nil
}

func GetAssociatedNetwork(ifaceName string) (Network, error) {
	return getAssociatedNetworkAirport()
}

func getAssociatedNetworkAirport() (Network, error) {
	var (
		err error
		raw []byte
	)

	if raw, err = exec.Command(airPortPath, "-I").Output(); err != nil {
		return Network{}, err
	}

	output := string(raw)

	bssIDA := bssIDRegex.FindStringSubmatch(output)
	if len(bssIDA) != numGroups2 {
		return Network{}, fmt.Errorf("failed to parse BSSID '%s', expected pattern '%s'", output, bssIDRegex)
	}

	ssIDA := ssIDRegex.FindStringSubmatch(output)
	if len(ssIDA) != numGroups2 {
		return Network{}, fmt.Errorf("failed to parse SSID '%s', expected pattern '%s'", output, ssIDRegex)
	}

	chA := channelRegex.FindStringSubmatch(output)
	if len(chA) != numGroups2 {
		return Network{}, fmt.Errorf("failed to parse channel '%s', expected pattern '%s'", output, channelRegex)
	}

	var chI int
	if chI, err = strconv.Atoi(chA[1]); err != nil {
		return Network{}, fmt.Errorf("failed to parse channel '%s': %w", chA, err)
	}

	return Network{
		BSSID:   bssIDA[1],
		SSID:    ssIDA[1],
		Channel: uint8(chI),
	}, nil
}

// Change radio channel on given interface.
func SetInterfaceChannel(ifaceName string, channel int) error {
	return exec.Command(airPortPath, ifaceName, fmt.Sprintf("-c%d", channel)).Run() //nolint:gosec // ignore
}
