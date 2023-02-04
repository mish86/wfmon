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
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.Border{}).
		BorderBottom(true).
		Bold(false)
	s.Cell = s.Cell.
		BorderStyle(lipgloss.Border{})
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)

	return s
}

// Returns default columns as slice and map.
func getDefaultCols() (ColumnViewSlice, ColumnViewMap) {
	cols := GenerateColumns(
		ColumnSSID,
		ColumnBSSID,
		ColumnChan,
		ColumnWidth,
		ColumnBand,
		ColumnSignal,
		ColumnNoise,
		ColumnSNR,
	)

	colsMap := make(ColumnViewMap, len(cols))
	for _, col := range cols {
		if swapper, ok := col.(ColumnSwapper); ok {
			for _, title := range swapper.Titles() {
				colsMap[title] = col
			}
		} else {
			colsMap[col.Title()] = col
		}
	}

	return cols, colsMap
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
func NewTableView() *TableView {
	cols, colsByName := getDefaultCols()
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
