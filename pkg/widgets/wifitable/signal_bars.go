package wifitable

import (
	netdata "wfmon/pkg/data/net"

	"github.com/charmbracelet/lipgloss"
)

// Bars presentation of signal quality.
type Bars netdata.Quality

// Returns style per each bar (signal quality) value.
func (q Bars) Style() lipgloss.Style {
	//nolint:gomnd // ignore
	switch {
	case q >= 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#77dd77")) // green
	case 60 <= q && q < 80:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#a7c7e7")) // blue
	case 40 <= q && q < 60:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ffb347")) // orange
	case 20 <= q && q < 40:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")) // red
	case q < 20:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")) // red
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#ff6961")) // red
	}
}

// Returns bars string presentation of signal quality.
func (q Bars) String() string {
	//nolint:gomnd // ignore
	switch {
	case q >= 80:
		return "▂▄▆█"
	case 60 <= q && q < 80:
		return "▂▄▆▁"
	case 40 <= q && q < 60:
		return "▂▄▁▁"
	case 20 <= q && q < 40:
		return "▂▁▁▁"
	case q < 20:
		return "▁▁▁▁"
	default:
		return "▁▁▁▁"
	}
}
