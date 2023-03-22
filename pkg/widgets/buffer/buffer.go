package buffer

import (
	"strings"
	log "wfmon/pkg/logger"

	"github.com/charmbracelet/lipgloss"
)

type cellStyle struct {
	fg lipgloss.Color
	bg lipgloss.Color
}

type Buffer struct {
	w, h   int
	cells  []rune
	styles []cellStyle
}

func New(w, h int) *Buffer {
	return &Buffer{
		w:      w,
		h:      h,
		cells:  make([]rune, w*h),
		styles: make([]cellStyle, w*h),
	}
}

func (b *Buffer) Width() int {
	return b.w
}

func (b *Buffer) Height() int {
	return b.h
}

// (0,0) at left bottom corner.
func (b *Buffer) idx(x, y int) int {
	return b.w*(b.h-y-1) + x
}

// (0,0) at left bottom corner.
// s - foreground and background colors.
func (b *Buffer) SetCell(x, y int, r rune, s ...lipgloss.Color) {
	idx := b.idx(x, y)

	if idx >= len(b.cells) {
		log.Fatalf("Index %d (%d, %d) out of bounds buffer dimensions (%d)", idx, x, y, len(b.cells))
		return
	}

	b.cells[idx] = r

	if len(s) > 0 {
		b.styles[idx].fg = s[0]
	}

	if len(s) > 1 {
		b.styles[idx].bg = s[1]
	}
}

// (0,0) at left bottom corner.
// Returns rune, foreground and background color.
func (b *Buffer) GetCell(x, y int) (rune, lipgloss.Color, lipgloss.Color) {
	idx := b.idx(x, y)
	return b.cells[idx], b.styles[idx].fg, b.styles[idx].bg
}

// Returns rendered rows.
func (b *Buffer) Rows() []string {
	rows := make([]string, b.h)

	rowID, start, end := 0, 0, 0
	s := b.styles[start]
	style := lipgloss.NewStyle().Foreground(s.fg).Background(s.bg)
	row := strings.Builder{}
	for ; end < len(b.cells); end++ {
		if end != 0 && end%b.w == 0 {
			style.Width(end - start)
			row.WriteString(style.Render(string(b.cells[start:end])))
			rows[rowID] = row.String()
			row = strings.Builder{}
			start = end
			s = b.styles[start]
			style = lipgloss.NewStyle().Foreground(s.fg).Background(s.bg)
			rowID++
			continue
		}
		if s != b.styles[end] {
			style.Width(end - start)
			row.WriteString(style.Render(string(b.cells[start:end])))
			start = end
			s = b.styles[start]
			style = lipgloss.NewStyle().Foreground(s.fg).Background(s.bg)
		}
	}

	style.Width(end - start)
	row.WriteString(style.Render(string(b.cells[start:end])))
	rows[rowID] = row.String()

	return rows
}
