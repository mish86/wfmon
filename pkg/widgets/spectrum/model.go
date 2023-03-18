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
	focused        bool
	band           wifi.Band
	selected       netdata.Key
	data           []Wave
	minVal, maxVal float64
	stepVal        float64
}

type Option func(*Model)

func WithSelected(key netdata.Key) Option {
	return func(m *Model) {
		m.Selected(key)
	}
}

func WithFocused(focus bool) Option {
	return func(m *Model) {
		m.Focused(focus)
	}
}

func WithMinVal(val float64) Option {
	return func(m *Model) {
		m.SetMinVal(val)
	}
}

func WithMaxVal(val float64) Option {
	return func(m *Model) {
		m.SetMaxVal(val)
	}
}

func New(opts ...Option) *Model {
	m := &Model{
		viewport: viewport.New(defaultWidth, defaultHeight),
		focused:  true,
		band:     wifi.ISM,
		data:     []Wave{},
		minVal:   -100,
		maxVal:   0,
		stepVal:  10,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) Selected(key netdata.Key) {
	m.selected = key
}

func (m *Model) GetSelected() netdata.Key {
	return m.selected
}

func (m *Model) Focused(focus bool) {
	m.focused = focus
}

func (m *Model) GetFocused() bool {
	return m.focused
}

func (m *Model) SetWidth(w int) {
	m.viewport.Width = w
}

func (m *Model) SetMinVal(val float64) {
	m.minVal = val
}

func (m *Model) SetMaxVal(val float64) {
	m.maxVal = val
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
			fmt.Sprintf("%3.f", m.maxVal),
			m.viewport.View(),
			m.viewAxeX(),
			fmt.Sprintf("%3.f", m.minVal-m.stepVal),
		),
	)
}
