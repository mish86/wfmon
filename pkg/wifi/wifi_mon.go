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
	handle   *pcap.Handle
	framesCh chan Frame
}

type Config struct {
	IFace *net.Interface
}

func NewMonitor(cfg *Config) *Monitor {
	return &Monitor{
		iface:    cfg.IFace,
		framesCh: make(chan Frame, defaultFramesBuffer),
	}
}

// Disconnects interface from network (AP) and creates active pcap.Handle for further packets sniffering.
func (mon *Monitor) Configure() error {
	var err error

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

	log.Infof("👀 monitoring %s", mon.iface.Name)
	// inbound packets
	packetSource := gopacket.NewPacketSource(mon.handle, mon.handle.LinkType())
	packets := packetSource.Packets()
	// start serv inbound packets
	for {
		select {
		case packet, ok := <-packets:
			if !ok {
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
			log.Infof("stopping accepting inbound packets on %s", mon.iface.Name)
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
