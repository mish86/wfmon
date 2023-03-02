package row

import (
	netdata "wfmon/pkg/data/net"
	column "wfmon/pkg/widgets/wifitable/col"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
)

// View property key.
type propKey int

const (
	rowStyle  propKey = iota // style for each cell in a row (default, associated network, etc)
	hashColor                // first column (#) with uniq color per network
)

type props map[propKey]any

// Row network data with view properties @propKey.
type Data struct {
	netdata.Network
	opts props
}

func (r *Data) init() {
	if r.opts == nil {
		r.opts = make(props)
	}
}

func (r *Data) set(key propKey, value any) {
	r.init()
	r.opts[key] = value
}

func (r Data) getAsColor(k propKey) lipgloss.TerminalColor {
	var (
		ok   bool
		prop any
		c    lipgloss.TerminalColor
	)
	if prop, ok = r.opts[k]; !ok {
		return lipgloss.NoColor{}
	}
	if c, ok = prop.(lipgloss.TerminalColor); ok {
		return c
	}

	return lipgloss.NoColor{}
}

func (r Data) getAsStyle(k propKey) lipgloss.Style {
	var (
		ok   bool
		prop any
		s    lipgloss.Style
	)
	if prop, ok = r.opts[k]; !ok {
		return lipgloss.Style{}
	}
	if s, ok = prop.(lipgloss.Style); ok {
		return s
	}
	return lipgloss.Style{}
}

func (r Data) HashColor(c lipgloss.Color) Data {
	r.set(hashColor, c)
	return r
}

func (r Data) GetHashColor() lipgloss.TerminalColor {
	return r.getAsColor(hashColor)
}

func (r Data) Style(s lipgloss.Style) Data {
	r.set(rowStyle, s)
	return r
}

func (r Data) GetRowStyle() lipgloss.Style {
	return r.getAsStyle(rowStyle)
}

// Cell viewer.
// Accepts row data and returns string, @table.StyledCell, averything that @table.RowData accepts.
type FncCellViewer func(row *Data) any

// Generates row data converter to @table.Row.
// Takes ordered columns definitions and cells viewers.
// If there is no viewer for a column then no data displayed for a cell of a row.
func Converter(columns []column.Column, viewers map[string]FncCellViewer) func(data *Data) table.Row {
	// columns and viewers can be copied in generator
	return func(data *Data) table.Row {
		row := make(table.RowData, len(columns))
		for i := range columns {
			col := columns[i]
			key := col.Key()
			if viewer, found := viewers[key]; found {
				row[key] = viewer(data)
			}
		}

		return table.NewRow(row)
	}
}
