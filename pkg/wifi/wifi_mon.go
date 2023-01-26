package wifi

import (
	"context"
	"fmt"
	"net"
	"time"
	log "wfmon/pkg/logger"
	"wfmon/pkg/network"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

const (
	DefaultTimeout = 500 * time.Millisecond
)

type Monitor struct {
	ctx  context.Context
	stop context.CancelFunc

	iface  *net.Interface
	handle *pcap.Handle
}

type Config struct {
	IFace *net.Interface
}

func NewMonitor(cfg *Config) *Monitor {
	return &Monitor{
		iface: cfg.IFace,
	}
}

// Disconnects interface from network (AP) and creates active pcap.Handle for further packets sniffering.
func (mon *Monitor) Configure() error {
	var err error

	log.Debugf("diassociate %s from any network", mon.iface.Name)
	if err = network.DisassociateFromNetwork(mon.iface.Name); err != nil {
		return err
	}

	log.Debugf("activate monitor on %s", mon.iface.Name)
	if mon.handle, err = network.CaptureWithTimeout(mon.iface.Name, DefaultTimeout); err != nil {
		return err
	}

	return nil
}

// Closes pcap.Handle.
func (mon *Monitor) Close() {
	if mon.handle != nil {
		log.Debugf("closing pcap handle")
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
		case <-mon.ctx.Done():
			log.Infof("stopping accepting inbound packets on %s", mon.iface.Name)
			return nil

		case packet := <-packets:
			p := FromPacket(packet)
			frame := p.DiscoverMgmtFrame()
			// frame := p.DiscoverIEs()
			if frame != nil {
				log.Debugf("%+v", frame)
			}

		default:
			continue
		}
	}
}

// Shutdowns sniffering service.
func (mon *Monitor) Shutdown() error {
	log.Infof("shutting down WiFi Monitor")
	mon.stop()
	return nil
}
