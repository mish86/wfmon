package spectrum

import (
	"fmt"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/wifi"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeight = 10
	defaultWidth  = 95
)

type Model struct {
	viewport       viewport.Model
	waves          []Wave
	band           wifi.Band
	selected       netdata.Key
	minVal, maxVal int
	stepVal        int
}

type Option func(*Model)

func New(opts ...Option) *Model {
	m := &Model{
		viewport: viewport.New(defaultWidth, defaultHeight),
		waves:    []Wave{},
		band:     wifi.ISM,
		minVal:   -100,
		stepVal:  10,
		maxVal:   0,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) SetWidth(w int) {
	m.viewport.Width = w
}

func (m *Model) SetBandView(b wifi.Band) {
	if b == wifi.UNII2B {
		return
	}

	m.band = b
}

func (m *Model) NextBandView() {
	// next band
	b := uint8(m.band) + 1

	// unsupported
	if b == uint8(wifi.UNII2B) {
		b++
	}

	// cycle
	if b > uint8(wifi.MaxBand) {
		b = uint8(wifi.MinBand)
	}

	m.band = wifi.Band(b)

	// TODO move to event
	m.refresh()
}

func (m *Model) Title() string {
	return fmt.Sprintf("%s / %s", "RSSI", m.band)
}

func (m *Model) View() string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.viewAxeY(),
		lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("%d", m.maxVal),
			m.viewport.View(),
			m.viewAxeX(),
			fmt.Sprintf("%d", m.minVal-m.stepVal),
		),
	)
}
