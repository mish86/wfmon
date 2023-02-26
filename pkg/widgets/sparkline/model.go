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
	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	defaultHeight          = 10
	defaultWidth           = 95
	axeYWidth              = 1
	defaultColor           = lipgloss.Color("#EE6FF8")
	defaultRefreshInterval = time.Second
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		b.Left = "┤"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	axeYstyle = func() lipgloss.Style {
		b := lipgloss.NormalBorder()
		return lipgloss.NewStyle().Border(b, false, true, false, false)
	}
)

type Model struct {
	viewport viewport.Model

	sparkline      *widgets.Sparkline
	sparklineGroup *widgets.SparklineGroup
	color          lipgloss.Color
	axesShown      bool

	fieldKey   string
	netKey     netdata.Key
	dataSource ds.TimeSeriesProvider

	minVal   float64
	maxVal   float64
	modifier func(float64) float64
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

func New(opts ...Option) *Model {
	sparkline := widgets.NewSparkline()
	sparkline.LineColor = tui.ColorGreen
	sparkline.MaxVal = 0

	sparklineGroup := widgets.NewSparklineGroup(sparkline)
	sparklineGroup.Title = ""
	sparklineGroup.Border = false
	sparklineGroup.SetRect(0, 0, defaultWidth, defaultHeight)

	m := &Model{
		viewport:       viewport.New(defaultWidth, defaultHeight),
		sparkline:      sparkline,
		sparklineGroup: sparklineGroup,
		color:          defaultColor,
		axesShown:      false,
		dataSource:     ds.EmptyProvider{},
		minVal:         0,
		maxVal:         sparkline.MaxVal,
		modifier:       func(v float64) float64 { return v },
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
	m.sparklineGroup.SetRect(
		0, 0,
		cmp.Nvl(m.axesShown, m.viewport.Width-1, m.viewport.Width),
		// @see refresh
		m.viewport.Height+1,
	)
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
	m.sparkline.MaxVal = val
}

func (m *Model) SetModifier(f func(float64) float64) {
	m.modifier = f
}

// Views data redered by @refresh in viewport.
// Axe Y takes extra 2 lines to viewport height.
func (m *Model) View() string {
	// do not display widget when no data
	if len(m.sparkline.Data) == 0 {
		return ""
	}

	content := strings.Builder{}
	content.WriteString(m.viewTitle())
	content.WriteRune('\n')
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

// TODO: move tabs in dashboard
func (m *Model) viewTitle() string {
	title := titleStyle.Render(m.fieldKey)
	gaps := strings.Repeat("─", cmp.Max(0, (m.viewport.Width-lipgloss.Width(title)))/2)
	return lipgloss.JoinHorizontal(lipgloss.Center, gaps, title, gaps)
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
