package sparkline

import (
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
		Range(m.viewport.Width, m.modifier).
		Reverse()
}

// Immediately renders data to viewport.
func (m *Model) refresh() {
	buf := buffer.New(m.viewport.Width, m.viewport.Height)

	maxVal := m.maxVal
	// if maxVal == 0 {
	// 	maxVal, _ = GetMaxFloat64FromSlice(sl.Data)
	// }
	sparklineHeight := m.viewport.Height

	bars := widgets.VBars()

	// draw line
	data := m.data
	for i := 0; i < len(data); i++ {
		x := m.viewport.Width - i - 1

		fh := (data[i] / maxVal) * float64(sparklineHeight)
		height := int(fh)

		if height == 0 {
			sparkChar := bars[1]
			buf.SetCell(x, 0, sparkChar, m.color)
			continue
		}

		sparkChar := bars[len(bars)-1]
		for y := 0; y < height-1; y++ {
			buf.SetCell(x, y, sparkChar, m.color)
		}

		spike := bars[len(bars)-1]
		r := int(fh*float64(10)) % 10
		if r > 0 {
			spike = bars[cmp.Min(r, len(bars)-1)]
		}
		buf.SetCell(x, height-1, spike, m.color)
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
