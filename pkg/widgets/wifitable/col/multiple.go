package column

import (
	"wfmon/pkg/widgets/wifitable/sort"

	"github.com/charmbracelet/lipgloss"
)

// Cycler of @Multiple column view.
type Cycler interface {
	Next() Multiple
	Current() Simple
	Prev() Multiple
}

// Multiple column view.
type Multiple struct {
	Simples     // swappable columns
	current int // currently viewed
}

// Returns @Multiple column.
func NewMultiple(cols ...Simple) Multiple {
	return Multiple{
		Simples: cols,
		current: 0,
	}
}

func (c Multiple) Next() Multiple {
	col := c.Clone()
	col.current++
	if col.current >= len(c.Simples) {
		col.current = 0
	}

	return col
}

func (c Multiple) Current() Simple {
	return c.Simples[c.current]
}

func (c Multiple) Prev() Multiple {
	col := c.Clone()
	col.current--
	if col.current < 0 {
		col.current = len(c.Simples) - 1
	}

	return col
}

func (c Multiple) Clone() Multiple {
	var cols = make([]Simple, len(c.Simples))
	copy(cols, c.Simples)

	return Multiple{
		Simples: cols,
		current: c.current,
	}
}

func (c Multiple) Key() string {
	return c.Current().Key()
}

func (c Multiple) Width() int {
	return c.Current().Width()
}

func (c Multiple) Style() lipgloss.Style {
	return c.Current().Style()
}

func (c Multiple) Sorter() sort.FncSorter {
	return c.Current().Sorter()
}
