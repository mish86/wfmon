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

var channelFrequencyRegex = regexp.MustCompile(`^(\d{1,3})\s\([2,5]GHz\)$`)

const numchFreqRegexGroups = 2

type systemProfiler struct {
	SPAirPortDataType []spAirportDataType `json:"SPAirPortDataType"`
}

type spAirportDataType struct {
	IFaces []spAirportAirportInterface `json:"spairport_airport_interfaces"`
}

type spAirportAirportInterface struct {
	Name      string   `json:"_name"`
	Channel   int      `json:"spairport_airdrop_channel"`
	Supported []string `json:"spairport_supported_channels"`
}

// Disconnects interface from network. Required before for setting interface in monitoring mode.
func DisassociateFromNetwork(ifaceName string) error {
	return exec.Command(airPortPath, ifaceName, "-z").Run()
}

// Returns channels supported by given interface.
func GetSupportedChannels(ifaceName string) ([]int, error) {
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
		if iface.Name != ifaceName {
			continue
		}

		if len(iface.Supported) == 0 {
			return []int{}, fmt.Errorf("no channels supported by %s", ifaceName)
		}

		channels := make([]int, len(iface.Supported))
		for i, chFreq := range iface.Supported {
			chA := channelFrequencyRegex.FindStringSubmatch(chFreq)
			if len(chA) != numchFreqRegexGroups {
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

	return []int{}, fmt.Errorf("no channels supported by %s", ifaceName)
}

// Change radio channel on given interface.
func SetInterfaceChannel(ifaceName string, channel int) error {
	return exec.Command(airPortPath, ifaceName, fmt.Sprintf("-c%d", channel)).Run() //nolint:gosec // ignore
}
