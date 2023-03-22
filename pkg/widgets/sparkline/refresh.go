package sparkline

import (
	"math"
	"time"
	"wfmon/pkg/ts"
	"wfmon/pkg/utils/cmp"
	"wfmon/pkg/widgets"
	"wfmon/pkg/widgets/buffer"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type refreshMsg time.Time

// Invokes refresh chart by refreshInterval.
// Fresh data obtained on timer end.
func refreshTick(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return refreshMsg(t)
	})
}

// Returns vector as a range of timeseries by network and field keys.
// Vector is reversed to render chart from right to left having new values on right.
func (m *Model) getData() ts.Vector {
	return m.dataSource.
		TimeSeries(m.netKey)(m.fieldKey).
		Range(m.viewport.Width).
		Reverse()
}

// Immediately renders data to viewport.
func (m *Model) refresh() {
	if !m.focused {
		return
	}

	buf := buffer.New(m.viewport.Width, m.viewport.Height)

	// absolut max value
	absMaxVal := cmp.Max(math.Abs(m.minVal), math.Abs(m.maxVal))
	viewHeight := float64(m.viewport.Height)

	bars := widgets.VBars()

	// draw line
	data := m.data
	for i := 0; i < len(data); i++ {
		x := m.viewport.Width - i - 1

		// max val of a range
		maxVal := math.Copysign(absMaxVal, data[i])
		// height value
		fh := data[i] * viewHeight / maxVal
		// values in a range less than zero
		if math.Signbit(maxVal) {
			// reverse
			fh = cmp.Max(0, viewHeight-fh)
		}
		// height in chars
		height := int(fh)

		// zero value spark
		if height == 0 {
			sparkChar := bars[1]
			buf.SetCell(x, 0, sparkChar, m.color)
			continue
		}

		// value spark
		sparkChar := bars[len(bars)-1]
		for y := 0; y < height; y++ {
			buf.SetCell(x, y, sparkChar, m.color)
		}

		// add spark spike
		if fh > float64(height) {
			r := int(fh*float64(10)) % 10
			spike := bars[cmp.Min(r, len(bars)-1)]
			buf.SetCell(x, height, spike, m.color)
		}
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, buf.Rows()...))
}

// Handles refresh tick.
// Fetches data from data source.
// Applies in the chart.
func (m *Model) onRefreshMsg(msg refreshMsg) {
	m.data = m.getData()

	m.refresh()
}
