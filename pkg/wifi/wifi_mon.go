package wifi

import (
	"context"
	"fmt"
	"net"
	"time"
	log "wfmon/pkg/logger"
	"wfmon/pkg/network"
	radionet "wfmon/pkg/network/radio"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	defaultTimeout      = 500 * time.Millisecond
	defaultFramesBuffer = 100
)

type Monitor struct {
	ctx  context.Context
	stop context.CancelFunc

	iface    *net.Interface
	file     string
	handle   *pcap.Handle
	framesCh chan Frame
}

type Config struct {
	IFace *net.Interface
	File  string
}

func NewMonitor(cfg *Config) *Monitor {
	return &Monitor{
		iface:    cfg.IFace,
		file:     cfg.File,
		framesCh: make(chan Frame, defaultFramesBuffer),
	}
}

// Disconnects interface from network (AP) and creates active pcap.Handle for further packets sniffering.
func (mon *Monitor) Configure() error {
	var err error

	if mon.isFromFile() {
		log.Infof("pcap file provided %s", mon.file)

		if mon.handle, err = network.CaptureFromFile(mon.file); err == nil {
			// loaded from file
			return nil
		} else if !mon.isFromIFace() {
			// failed to load from file and iface not provided
			return err
		}
	}

	if !mon.isFromIFace() {
		return fmt.Errorf("no interface provided to monitor")
	}

	log.Infof("iface provided %s", mon.iface.Name)

	log.Debugf("diassociate %s from any network", mon.iface.Name)
	if err = radionet.DisassociateFromNetwork(mon.iface.Name); err != nil {
		return err
	}

	log.Debugf("activate monitor on %s", mon.iface.Name)
	if mon.handle, err = network.CaptureWithTimeout(mon.iface.Name, defaultTimeout); err != nil {
		return err
	}

	return nil
}

// Closes pcap.Handle.
func (mon *Monitor) Close() {
	if mon.handle != nil {
		log.Info("closing pcap handle")
		mon.handle.Close()
	}
}

// Starts sniffering packets until shutdown.
func (mon *Monitor) Start(ctx context.Context) error {
	if mon.handle == nil {
		return fmt.Errorf("handle is nil")
	}

	mon.ctx, mon.stop = context.WithCancel(ctx)

	if mon.isFromIFace() {
		log.Infof("👀 monitoring %s", mon.iface.Name)
	}
	if mon.isFromFile() {
		log.Infof("👀 reading %s", mon.file)
	}

	// inbound packets
	packetSource := gopacket.NewPacketSource(mon.handle, mon.handle.LinkType())
	packets := packetSource.Packets()
	// start serv inbound packets
	for {
		select {
		case packet, ok := <-packets:
			// channel is closed or packets from file is over
			if !ok {
				if mon.isFromFile() {
					<-mon.ctx.Done()
					return nil
				}
				return fmt.Errorf("packet source closed, stopping monitoring")
			}

			p := FromPacket(packet)
			frame := p.DiscoverMgmtFrame()
			if frame != nil {
				log.Debugf("%+v", frame)
				// send a copy of frame to output channel
				// if len(frame.SSID) > 0 {
				mon.framesCh <- Frame(*frame)
				// }
			}

		case <-mon.ctx.Done():
			if mon.isFromIFace() {
				log.Infof("stopping accepting inbound packets on %s", mon.iface.Name)
			}
			if mon.isFromFile() {
				log.Infof("stopping reading packets from file %s", mon.file)
			}
			return nil
		}
	}
}

// Shutdowns sniffering service.
func (mon *Monitor) Stop() error {
	log.Infof("stopping WiFi Monitor")
	mon.stop()
	return nil
}

// Returns frames output channel.
func (mon *Monitor) GetFrames() <-chan Frame {
	return mon.framesCh
}

func (mon *Monitor) isFromFile() bool {
	return len(mon.file) > 0
}

func (mon *Monitor) isFromIFace() bool {
	return mon.iface != nil
}
