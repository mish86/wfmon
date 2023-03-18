package column

import (
	"wfmon/pkg/widgets/sort"

	"github.com/charmbracelet/lipgloss"
)

// Simple column defintion.
type Simple struct {
	key    string         // header key and title
	width  int            // column width
	style  lipgloss.Style // style for column header
	sorter sort.FncSorter // network data rows sorter by the column
}

type Simples []Simple

// Returns @Simple column.
func NewSimple(key string, width int) Simple {
	return Simple{
		key:    key,
		width:  width,
		sorter: sort.ByKeySorter(),
	}
}

func (c Simple) Key() string {
	return c.key
}

func (c Simple) Width() int {
	return c.width
}

func (c Simple) Style() lipgloss.Style {
	return c.style
}

func (c Simple) Sorter() sort.FncSorter {
	return c.sorter
}

func (c Simple) Copy() Simple {
	var replica = c
	return replica
}

func (c Simple) WithStyle(s lipgloss.Style) Simple {
	c = c.Copy()
	c.style = s
	return c
}

func (c Simple) WithSorter(s sort.FncSorter) Simple {
	c = c.Copy()
	c.sorter = s
	return c
}
