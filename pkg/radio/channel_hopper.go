package radio

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"
	log "wfmon/pkg/logger"
	radionet "wfmon/pkg/network/radio"
	"wfmon/pkg/repeater"
)

const (
	DefaultHopInterval = 250 * time.Millisecond // Delay before next channel hop
)

type ChannelHopperServ struct {
	ctx  context.Context
	stop context.CancelFunc

	iface *net.Interface

	idx         int
	channels    []int
	hopInterval time.Duration
	chLock      sync.RWMutex
}

type ChannelHopperConfig struct {
	IFace       *net.Interface
	HopInterval time.Duration
}

func NewChannelHopperServ(cfg *ChannelHopperConfig) *ChannelHopperServ {
	return &ChannelHopperServ{
		iface:       cfg.IFace,
		idx:         0,
		hopInterval: cfg.HopInterval,
	}
}

// Loads supported channels from configured interface.
func (h *ChannelHopperServ) Configure() error {
	var err error
	log.Infof("Loading supported channel on '%s'", h.iface.Name)
	h.channels, err = radionet.GetSupportedChannels(h.iface.Name)

	return err
}

func (h *ChannelHopperServ) Close() {
	//
}

// Locks on writing, shifts current channel, sets it on interface and unlocks.
func (h *ChannelHopperServ) hop() error {
	h.chLock.Lock()
	defer h.chLock.Unlock()

	if h.idx++; h.idx >= len(h.channels) {
		h.idx = 0
	}

	log.Debugf("Interface %s hopping to channel %d", h.iface.Name, h.channels[h.idx])
	return radionet.SetInterfaceChannel(h.iface.Name, h.channels[h.idx])
}

// Returns current channel number.
func (h *ChannelHopperServ) Channel() int {
	h.chLock.RLock()
	defer h.chLock.RUnlock()

	return h.channels[h.idx]
}

// Start hopping until shutdown.
func (h *ChannelHopperServ) Start(ctx context.Context) error {
	if len(h.channels) == 0 {
		return fmt.Errorf("no supported channels for hopping")
	}

	h.ctx, h.stop = context.WithCancel(ctx)

	log.Infof("🐇 hopping %s", h.iface.Name)

	repeater.Default(h.ctx, h.hopInterval, func() {
		if err := h.hop(); err != nil {
			log.Errorf("failed to hop, got %w", err)
		}
	}, func() {
		log.Infof("stopping hopping on %s", h.iface.Name)
	})

	return nil
}

// Shutdowns hopper service.
func (h *ChannelHopperServ) Stop() error {
	log.Infof("stopping Channel Hopper")
	if h.stop != nil {
		h.stop()
	}
	return nil
}
