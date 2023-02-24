package sparkline

import (
	"image"
	"strings"
	"time"
	"wfmon/pkg/ts"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	tui "github.com/gizak/termui/v3"
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
	if m.dataSource == nil {
		return ts.Vector{}
	}

	return m.dataSource.
		TimeSeries(m.netKey)(m.fieldKey).
		Range(m.viewport.Width, func(val float64) float64 { return m.sparkline.MaxVal + val }).
		Reverse()
}

// Immediately renders data to viewport.
func (m *Model) refresh() {
	item := m.sparklineGroup
	buf := tui.NewBuffer(item.GetRect())
	item.Lock()
	item.Draw(buf)
	item.Unlock()

	inlineStyle := lipgloss.NewStyle().Width(item.GetRect().Dx()).MaxWidth(item.GetRect().Dx()).Inline(true)
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#EE6FF8"))

	renderedRows := make([]string, item.GetRect().Dy())
	for y := 0; y < item.GetRect().Dy(); y++ {
		row := strings.Builder{}
		// reverse order: from right to left
		for x := item.GetRect().Dx() - 1; x >= 0; x-- {
			cell := buf.GetCell(image.Point{x, y})
			row.WriteRune(cell.Rune)
		}
		rowText := row.String()
		renderedRows[y] = rowText
		renderedRows[y] = inlineStyle.Render(renderedRows[y])
		renderedRows[y] = rowStyle.Render(renderedRows[y])
	}

	m.viewport.SetContent(
		lipgloss.JoinVertical(lipgloss.Right, renderedRows...),
	)
}

// Handles refresh tick.
// Fetches data from data source.
// Applies in the chart.
func (m *Model) onRefreshMsg(msg refreshMsg) {
	m.sparkline.Data = m.getData()

	m.refresh()
}
