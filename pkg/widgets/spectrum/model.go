package spectrum

import (
	"fmt"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"
	"wfmon/pkg/widgets/events"
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
	data           []Wave
	selected       netdata.Key
	fieldKey       string
	minVal, maxVal float64
	dataSource     ds.TimeSeriesProvider
}

type Option func(*Model)

func WithDataSource(dataSource ds.TimeSeriesProvider) Option {
	return func(m *Model) {
		m.SetDataSource(dataSource)
	}
}

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

func WithField(key string) Option {
	return func(m *Model) {
		m.SetFieldKey(key)
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

func WithSignalField(msg events.SignalFieldMsg) Option {
	return func(m *Model) {
		m.SetFieldKey(msg.Key)
		m.SetMinVal(msg.MinVal)
		m.SetMaxVal(msg.MaxVal)
	}
}

func New(opts ...Option) *Model {
	m := &Model{
		viewport: viewport.New(defaultWidth, defaultHeight),
		focused:  true,
		band:     wifi.ISM,
		data:     []Wave{},
		minVal:   0,
		maxVal:   0,
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) SetDataSource(dataSource ds.TimeSeriesProvider) {
	m.dataSource = dataSource
}

func (m *Model) Focused(focus bool) {
	m.focused = focus
	// m.viewport.SetContent("")
}

func (m *Model) GetFocused() bool {
	return m.focused
}

func (m *Model) SetWidth(w int) {
	m.viewport.Width = w
}

func (m *Model) Width() int {
	return m.viewport.Width
}

func (m *Model) Selected(key netdata.Key) {
	m.selected = key
}

func (m *Model) GetSelected() netdata.Key {
	return m.selected
}

func (m *Model) SetFieldKey(key string) {
	m.fieldKey = key
}

func (m *Model) FieldKey() string {
	return m.fieldKey
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
}

func (m *Model) Title() string {
	return fmt.Sprintf("%s / %s", m.fieldKey, m.band)
}

func (m *Model) View() string {
	maxVal := m.maxVal
	// yStep := (m.maxVal - m.minVal) / float64(m.viewport.Height)
	// minVal := m.minVal-yStep
	minVal := m.minVal

	return lipgloss.JoinHorizontal(lipgloss.Top,
		m.viewAxeY(),
		lipgloss.JoinVertical(lipgloss.Left,
			fmt.Sprintf("%3.f", maxVal),
			m.viewport.View(),
			m.viewAxeX(),
			fmt.Sprintf("%3.f", minVal),
		),
	)
}
