package column

import (
	"fmt"
	"wfmon/pkg/widgets/sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// Base column defintion.
type Column interface {
	Key() string
	Width() int
	Style() lipgloss.Style
	Sorter() sort.FncSorter
}

// Converts columns definition with applied sorting direction in the title to ordered array of @table.Column.
func Converter(columns []Column) func(sort Sort) []table.Column {
	// columns can be copied in generator
	return func(sort Sort) []table.Column {
		cols := make([]table.Column, len(columns))
		for i := range columns {
			col := columns[i]
			key := col.Key()
			title := col.Key()
			width := col.Width()

			if sort.Key() == key {
				title = fmt.Sprintf("%s %s", key, sort.Order())
			}

			cols[i] = table.NewColumn(
				key,
				title,
				width,
			).WithStyle(col.Style())
		}

		return cols
	}
}
