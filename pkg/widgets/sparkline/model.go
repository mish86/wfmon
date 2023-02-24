package sparkline

import (
	"time"
	netdata "wfmon/pkg/data/net"
	"wfmon/pkg/ds"
	"wfmon/pkg/ts"

	"github.com/charmbracelet/bubbles/viewport"
	tui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	defaultHeight          = 10
	defaultWidth           = 95
	defaultRefreshInterval = time.Second
)

type Model struct {
	viewport viewport.Model

	sparkline      *widgets.Sparkline
	sparklineGroup *widgets.SparklineGroup

	fieldKey   string
	netKey     netdata.Key
	vec        ts.Vector
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
		m.SetWidth(w)
		m.SetHeight(h)
	}
}

func WithMaxVal(val float64) Option {
	return func(m *Model) {
		m.SetMaxVal(val)
	}
}

func New(opts ...Option) *Model {
	sparkline := widgets.NewSparkline()
	sparkline.LineColor = tui.ColorGreen

	sparklineGroup := widgets.NewSparklineGroup(sparkline)
	sparklineGroup.Title = ""
	sparklineGroup.Border = false
	sparklineGroup.SetRect(0, 0, defaultWidth, defaultHeight)

	m := &Model{
		viewport:       viewport.New(defaultWidth, defaultHeight),
		sparkline:      sparkline,
		sparklineGroup: sparklineGroup,
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

func (m *Model) SetFieldKey(key string) {
	m.fieldKey = key
}

func (m *Model) SetWidth(w int) {
	m.viewport.Width = w
	m.sparklineGroup.SetRect(0, 0, m.viewport.Width, m.viewport.Height)
}

func (m *Model) SetHeight(h int) {
	m.viewport.Height = h
	m.sparklineGroup.SetRect(0, 0, m.viewport.Width, m.viewport.Height)
}

func (m *Model) SetMaxVal(val float64) {
	m.sparkline.MaxVal = val
}

// Views data redered by @refresh in viewport.
func (m *Model) View() string {
	return m.viewport.View()
}
