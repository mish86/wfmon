package sparkline

import (
	"fmt"
	"strings"
	"time"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"
	"wfmon/pkg/utils/cmp"

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
	modifier  func(float64) float64
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

func WithModifier(f func(float64) float64) Option {
	return func(m *Model) {
		m.SetModifier(f)
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
		modifier:   func(v float64) float64 { return v },
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

func (m *Model) SetDimension(w, h int) {
	m.viewport.Width = w
	m.viewport.Height = h
}

func (m *Model) SetWidth(w int) {
	m.SetDimension(w, m.viewport.Height)
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

func (m *Model) SetModifier(f func(float64) float64) {
	m.modifier = f
}

func (m *Model) Focused(focus bool) {
	m.focused = focus
}

func (m *Model) GetFocused() bool {
	return m.focused
}

func (m *Model) NetworkKey() netdata.Key {
	return m.netKey
}

func (m *Model) FieldKey() string {
	return m.fieldKey
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
			m.viewAxeYStart(),
			axeYstyle().Render(m.viewport.View()),
			m.viewAxeYEnd(),
		))
	} else {
		content.WriteString(m.viewport.View())
	}
	return content.String()
}

// Returns axe Y start point.
// Y axe vector from bottom to top.
func (m *Model) viewAxeYStart() string {
	return fmt.Sprintf("%.f┐", cmp.Max(0.0, m.maxVal-m.modifier(m.maxVal)))
}

// Returns axe Y start point.
// Y axe vector from bottom to top.
func (m *Model) viewAxeYEnd() string {
	return fmt.Sprintf("%.f┘", cmp.Min(0.0, m.minVal-m.modifier(m.minVal)))
}
