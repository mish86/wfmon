package wifitable

import (
	"sync"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultTableHeight = 10
)

// Returns default viewport style.
func getDefaultViewportStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderStyle(lipgloss.Border{}).
		BorderForeground(lipgloss.Color("240"))
}

// Returns default table styles.
func getDefaultTableStyles() table.Styles {
	// cell foreground overrides selected foreground, bug?
	return table.Styles{
		Header: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1, 0, 0), // keep space after column
		Cell: lipgloss.NewStyle().
			Bold(false).
			Padding(0, 1, 0, 0), // keep space after column
		Selected: lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1, 0, 0). // keep space after column
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")),
	}
}

// Network table view.
type TableView struct {
	table         table.Model
	cols          ColumnViewSlice
	colsByNames   ColumnViewMap
	viewportStyle lipgloss.Style
	tableStyles   table.Styles

	rowsLock sync.Mutex
}

// Returns new network table view.
func NewTableView(cols ColumnViewSlice) *TableView {
	// cols, colsByName := getDefaultCols()
	colsByName := cols.Map()
	viewportStyle := getDefaultViewportStyle()
	tableStyles := getDefaultTableStyles()

	t := table.New(
		table.WithColumns(cols.TableColumns()),
		table.WithHeight(defaultTableHeight),
		table.WithFocused(true),
	)
	t.SetStyles(tableStyles)

	t.UpdateViewport()

	return &TableView{
		table:         t,
		cols:          cols,
		colsByNames:   colsByName,
		viewportStyle: viewportStyle,
		tableStyles:   tableStyles,
	}
}

// Renderes table.
func (v *TableView) View() string {
	v.rowsLock.Lock()
	defer v.rowsLock.Unlock()

	return v.viewportStyle.Render(v.table.View()) + "\n"
}

// Updates rows in table.
// Does not invoke table redraw.
func (v *TableView) OnData(rows []table.Row) {
	v.rowsLock.Lock()
	defer v.rowsLock.Unlock()

	v.table.SetRows(rows)
}

// Applies sorting in table header.
func (v *TableView) OnSort(col ColumnViewer) {
	// v.cols[col.Index()] = &col
	columns := v.cols.SortBy(col)
	v.table.SetColumns(columns)
}

// Returns column viewer by number.
func (v *TableView) GetColumnByNum(num int) (ColumnViewer, bool) {
	idx := num - 1
	if idx < 0 || idx >= len(v.cols) {
		return nil, false
	}
	col := v.cols[idx]
	return col, true
}

// Returns column viewer by title.
func (v *TableView) GetColumnByTitle(title string) (ColumnViewer, bool) {
	col, found := v.colsByNames[title]
	return col, found
}

func (v *TableView) TableColumns() TableColumns {
	return v.cols.TableColumns()
}
