package sparkline

import (
	"fmt"
	"strings"
	"time"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"
	"wfmon/pkg/widgets/events"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultHeight          = 10
	defaultWidth           = 95
	axeYWidth              = 1
	defaultColor           = lipgloss.Color("#EE6FF8")
	defaultRefreshInterval = time.Second
)

var (
	axeYstyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		return lipgloss.NewStyle().Border(b, false, true, false, false)
	}
)

type Model struct {
	viewport viewport.Model
	focused  bool

	data      []float64
	minVal    float64
	maxVal    float64
	color     lipgloss.Color
	axesShown bool

	fieldKey   string
	netKey     netdata.Key
	dataSource ds.TimeSeriesProvider
}

type Option func(*Model)

func WithDataSource(dataSource ds.TimeSeriesProvider) Option {
	return func(m *Model) {
		m.SetDataSource(dataSource)
	}
}

func WithNetwork(key netdata.Key) Option {
	return func(m *Model) {
		m.SetNetworkKey(key)
	}
}

func WithField(key string) Option {
	return func(m *Model) {
		m.SetFieldKey(key)
	}
}

func WithDimention(w int, h int) Option {
	return func(m *Model) {
		m.SetDimension(w, h)
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

func WithYAxe(shown bool) Option {
	return func(m *Model) {
		m.ShowYAxe(shown)
	}
}

func WithFocused(focus bool) Option {
	return func(m *Model) {
		m.Focused(focus)
	}
}

func New(opts ...Option) *Model {
	m := &Model{
		viewport:   viewport.New(defaultWidth, defaultHeight),
		focused:    true,
		data:       []float64{},
		minVal:     0,
		maxVal:     0,
		color:      defaultColor,
		axesShown:  false,
		dataSource: ds.EmptyProvider{},
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *Model) SetDataSource(dataSource ds.TimeSeriesProvider) {
	m.dataSource = dataSource
}

func (m *Model) SetNetworkKey(key netdata.Key) {
	m.netKey = key
}

func (m *Model) NetworkKey() netdata.Key {
	return m.netKey
}

func (m *Model) SetColor(c lipgloss.Color) {
	m.color = c
}

func (m *Model) ShowYAxe(show bool) {
	m.axesShown = show
	m.SetDimension(m.viewport.Width, m.viewport.Height)
}

func (m *Model) SetFieldKey(key string) {
	m.fieldKey = key
}

func (m *Model) FieldKey() string {
	return m.fieldKey
}

func (m *Model) SetDimension(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h
}

func (m *Model) SetWidth(w int) {
	m.SetDimension(w, m.viewport.Height)
}

func (m *Model) Width() int {
	return m.viewport.Width
}

func (m *Model) SetHeight(h int) {
	m.SetDimension(m.viewport.Width, h)
}

func (m *Model) SetMinVal(val float64) {
	m.minVal = val
}

func (m *Model) SetMaxVal(val float64) {
	m.maxVal = val
}

func (m *Model) Focused(focus bool) {
	m.focused = focus
	// m.viewport.SetContent("")
}

func (m *Model) GetFocused() bool {
	return m.focused
}

func (m *Model) Title() string {
	return m.fieldKey
}

// Views data redered by @refresh in viewport.
// Axe Y takes extra 2 lines to viewport height.
func (m *Model) View() string {
	// do not display widget when no data
	if len(m.data) == 0 {
		return ""
	}

	content := strings.Builder{}
	if m.axesShown {
		content.WriteString(lipgloss.JoinVertical(lipgloss.Right,
			fmt.Sprintf("%.f┑", m.maxVal),
			axeYstyle().Render(m.viewport.View()),
			fmt.Sprintf("%.f┙", m.minVal),
		))
	} else {
		content.WriteString(m.viewport.View())
	}
	return content.String()
}
