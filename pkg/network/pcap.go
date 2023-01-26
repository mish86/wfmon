package network

import (
	"fmt"
	"time"
	log "wfmon/pkg/logger"

	"github.com/google/gopacket/pcap"
)

const (
	DefaultRadioMonitorMode = true
	DefaultSnapLen          = 65536
	DefaultBufSize          = 2_097_152
	DefaultPromisc          = true
	DefaultTimeout          = pcap.BlockForever
)

// Packet capture options for pcap lib.
type CaptureOptions struct {
	Monitor bool // sudo required for true
	Snaplen int
	Bufsize int
	Promisc bool
	Timeout time.Duration
}

// Returns default packet capture options.
func DefaultOptions() CaptureOptions {
	return CaptureOptions{
		Monitor: DefaultRadioMonitorMode,
		Snaplen: DefaultSnapLen,
		Bufsize: DefaultBufSize,
		Promisc: DefaultPromisc,
		Timeout: DefaultTimeout,
	}
}

// Returns active pcap.Handle for interface and capture options.
// Sudo preveliges required for setting interface in monitoring mode.
func CaptureWithOptions(ifName string, options CaptureOptions) (*pcap.Handle, error) {
	log.Debugf("creating capture for '%s' with options: %+v", ifName, options)

	var (
		err     error
		ihandle *pcap.InactiveHandle
	)
	ihandle, err = pcap.NewInactiveHandle(ifName)
	if err != nil {
		return nil, fmt.Errorf("error while opening interface %s: %w", ifName, err)
	}
	defer ihandle.CleanUp()

	if err = ihandle.SetRFMon(options.Monitor); err != nil {
		return nil, fmt.Errorf("error while setting interface %s in monitor mode: %w", ifName, err)
	}
	if err = ihandle.SetSnapLen(options.Snaplen); err != nil {
		return nil, fmt.Errorf("error while setting snapshot length: %w", err)
	}
	if err = ihandle.SetBufferSize(options.Bufsize); err != nil {
		return nil, fmt.Errorf("error while setting buffer size: %w", err)
	}
	if err = ihandle.SetPromisc(options.Promisc); err != nil {
		return nil, fmt.Errorf("error while setting promiscuous mode to %v: %w", options.Promisc, err)
	}
	if err = ihandle.SetTimeout(options.Timeout); err != nil {
		return nil, fmt.Errorf("error while setting timeout: %w", err)
	}

	return ihandle.Activate()
}

// Returns active pcap.Handle for interface with default capture options.
// Sudo preveliges required for setting interface in monitoring mode.
func Capture(ifName string) (*pcap.Handle, error) {
	return CaptureWithOptions(ifName, DefaultOptions())
}

// Returns active pcap.Handle for interface with default capture options and given timeout for handle.
// Sudo preveliges required for setting interface in monitoring mode.
func CaptureWithTimeout(ifName string, timeout time.Duration) (*pcap.Handle, error) {
	opts := DefaultOptions()
	opts.Timeout = timeout
	return CaptureWithOptions(ifName, opts)
}
