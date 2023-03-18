package spectrum

import (
	"image"
	"math"
	"strings"
	"wfmon/pkg/widgets/buffer"
	"wfmon/pkg/wifi"

	"github.com/charmbracelet/lipgloss"
)

func (m *Model) viewAxeXISM() string {
	rows := strings.Builder{}
	rows.WriteString("────┰───┰───┰───┰───┰───┰───┰───┰───┰───┰───┰───┰───┰───┰────────┰────┤\n")
	rows.WriteString("        1   2   3   4   5   6   7   8   9  10  11  12  13       14")
	return rows.String()
}

func (m *Model) viewAxeXUNII1() string {
	rows := strings.Builder{}
	rows.WriteString("───┰─────┰─────┰─────┰─────┰─────┰─────┰─────┰─────┰─────┤\n")
	rows.WriteString("        36    38    40    42    44    46    48    50")
	return rows.String()
}

func (m *Model) viewAxeXUNII2A() string {
	rows := strings.Builder{}
	rows.WriteString("───┰─────┰─────┰─────┰─────┰─────┰─────┰─────┰─────┰─────┤\n")
	rows.WriteString("        50    52    54    56    58    60    62    64")
	return rows.String()
}

func (m *Model) viewAxeXUNII2C() string {
	rows := strings.Builder{}
	rows.WriteString("────┰───────┰───────┰───────┰───────┰───────┰───────┰───────┤\n")
	rows.WriteString("           100     108     116     124     132     140")
	return rows.String()
}

func (m *Model) viewAxeXUNII3() string {
	rows := strings.Builder{}
	rows.WriteString("────┰───────┰───────┰───────┰───────┰───────┤\n")
	rows.WriteString("           149     157     165     173")
	return rows.String()
}

func (m *Model) viewAxeX() string {
	//nolint:exhaustive // ignore
	switch m.band {
	case wifi.ISM:
		return m.viewAxeXISM()
	case wifi.UNII1:
		return m.viewAxeXUNII1()
	case wifi.UNII2A:
		return m.viewAxeXUNII2A()
	case wifi.UNII2B:
		return ""
	case wifi.UNII2C:
		return m.viewAxeXUNII2C()
	case wifi.UNII3:
		return m.viewAxeXUNII3()
	default:
		return ""
	}
}

// TODO: consider RSSI and Quality values
func (m *Model) viewAxeY() string {
	rows := strings.Builder{}
	rows.WriteString("┍\n")
	rows.WriteString(strings.Repeat("│\n", m.viewport.Height))
	rows.WriteString("├\n│\n┕\n")
	return rows.String()
}

// Immediately renders data to viewport.
func (m *Model) refresh() {
	m.refreshByBand()
}

func (m *Model) refreshByBand() {
	filtered := []Wave{}
	for i := range m.data {
		if m.band == m.data[i].Band {
			filtered = append(filtered, m.data[i])
		}
	}

	if len(filtered) == 0 {
		m.viewport.SetContent("")
		return
	}

	// returns zero point (X,Y) offset in symbols, and X channel step scale in symbols
	var bandParams = func(b wifi.Band) (image.Point, int) {
		//nolint:exhaustive // ignore
		switch b {
		case wifi.ISM:
			//nolint:gomnd // ignore
			return image.Point{X: 4, Y: 0}, 4
		case wifi.UNII1:
			//nolint:gomnd // ignore
			return image.Point{X: -99, Y: 0}, 3
		case wifi.UNII2A:
			//nolint:gomnd // ignore
			return image.Point{X: -141, Y: 0}, 3
		case wifi.UNII2C:
			return image.Point{X: -88, Y: 0}, 1
		case wifi.UNII3:
			return image.Point{X: -137, Y: 0}, 1
		default:
			return image.Point{}, 0
		}
	}

	buf := buffer.New(m.viewport.Width, m.viewport.Height)

	var selected *Wave
	view := image.Point{X: m.viewport.Width, Y: m.viewport.Height}

	if len(filtered) == 1 {
		selected = &filtered[0]
	} else {
		for i := range m.data {
			w := m.data[i]
			if m.band != w.Band {
				continue
			}

			xy0, step := bandParams(w.Band)

			if m.selected.Compare(w.Key) == 0 {
				selected = &w
				continue
			}
			r1, r2 := w.rects(xy0, view, step, m.minVal)
			w.renderBordered(buf, r1, r2)
		}
	}

	if selected != nil {
		xy0, step := bandParams(selected.Band)

		r1, r2 := selected.rects(xy0, view, step, m.minVal)
		selected.renderSolid(buf, r1, r2)
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, buf.Rows()...))
}

// Returns rects for wave of primary and secondary channels.
// Secondary channel wave rect is optional.
func (wave *Wave) rects(xy0, view image.Point, xScale int, yMinVal float64) (image.Rectangle, image.Rectangle) {
	w := int(wave.Width) * wave20MhzWidth * xScale
	w1 := wave20MhzWidth * xScale
	h := view.Y - int(math.Floor(float64(view.Y)*wave.Value/yMinVal))
	leftMargin := xy0.X + (int(wave.LowerChannel())-halfOfWave20MhzWidth)*xScale
	bottomMargin := xy0.Y

	if wave.Sign < 0 {
		return image.Rect(leftMargin+w1, bottomMargin, leftMargin+w, bottomMargin+h-1),
			image.Rect(leftMargin, bottomMargin, leftMargin+w1-1, bottomMargin+h-1)
	} else if wave.Sign > 0 {
		return image.Rect(leftMargin, bottomMargin, leftMargin+w1, bottomMargin+h-1),
			image.Rect(leftMargin+w1+1, bottomMargin, leftMargin+w, bottomMargin+h-1)
	}

	return image.Rect(leftMargin, bottomMargin, leftMargin+w, bottomMargin+h-1), image.Rectangle{}
}

// Renders a wave in a buffer by rectangle using fill and border runes.
func (wave *Wave) render(buf *buffer.Buffer, r image.Rectangle, b lipgloss.Border, fill rune) {
	var fncR0 = func(s string) rune {
		if s == "" {
			return 0
		}
		return []rune(s)[0]
	}

	buf.SetCell(r.Min.X, r.Max.Y, fncR0(b.TopLeft), wave.Color)
	for x := r.Min.X + 1; x < r.Max.X; x++ {
		buf.SetCell(x, r.Max.Y, fncR0(b.Top), wave.Color)
	}
	buf.SetCell(r.Max.X, r.Max.Y, fncR0(b.TopRight), wave.Color)

	for y := r.Min.Y; y < r.Max.Y; y++ {
		buf.SetCell(r.Min.X, y, fncR0(b.Left), wave.Color)
		if fill != 0 {
			for x := r.Min.X + 1; x < r.Max.X; x++ {
				buf.SetCell(x, y, fill, wave.Color)
			}
		}
		buf.SetCell(r.Max.X, y, fncR0(b.Right), wave.Color)
	}
}

// Renders a bordered wave for primary and secondary channels with no fill.
// Secondary channel rect is optional.
func (wave *Wave) renderBordered(buf *buffer.Buffer, r ...image.Rectangle) {
	if len(r) == 0 {
		return
	}

	var r1 image.Rectangle
	if len(r) > 0 {
		r1 = r[0]
	}
	if wave.Sign != 0 && len(r) > 1 {
		r1 = r1.Union(r[1])
	}

	wave.render(buf, r1, lipgloss.RoundedBorder(), 0)
}

// Renders a solid wave for primary and secondary channels.
// Secondary channel rect is optional.
func (wave *Wave) renderSolid(buf *buffer.Buffer, r ...image.Rectangle) {
	if len(r) == 0 {
		return
	}

	if wave.Sign != 0 && len(r) > 1 {
		wave.renderSecondary(buf, r[1])
	}
	if len(r) > 0 {
		wave.renderPrimary(buf, r[0])
	}
}

// Renders a solid wave for primary channel.
func (wave *Wave) renderPrimary(buf *buffer.Buffer, r image.Rectangle) {
	var b lipgloss.Border
	switch {
	case wave.Sign > 0:
		b = lipgloss.Border{
			Top:      "▄",
			Left:     "▐",
			Right:    "█",
			TopLeft:  "▗",
			TopRight: "▄",
		}
	case wave.Sign < 0:
		b = lipgloss.Border{
			Top:      "▄",
			Left:     "█",
			Right:    "▌",
			TopLeft:  "▄",
			TopRight: "▖",
		}
	default:
		b = lipgloss.Border{
			Top:      "▄",
			Left:     "▐",
			Right:    "▌",
			TopLeft:  "▗",
			TopRight: "▖",
		}
	}
	wave.render(buf, r, b, '█')
}

// Renders a solid wave for secondary channel.
func (wave *Wave) renderSecondary(buf *buffer.Buffer, r image.Rectangle) {
	if wave.Sign == 0 {
		return
	}
	var b lipgloss.Border
	if wave.Sign > 0 {
		b = lipgloss.Border{
			Top:      "▖",
			Left:     "▒",
			Right:    "▌",
			TopLeft:  "▖",
			TopRight: "▖",
		}
	} else if wave.Sign < 0 {
		b = lipgloss.Border{
			Top:      "▗",
			Left:     "▐",
			Right:    "▒",
			TopLeft:  "▗",
			TopRight: "▗",
		}
	}
	wave.render(buf, r, b, '▒')
}
