package wifitable

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultTableHeight = 10
)

const (
	ColumnBSSIDTitle = "BSSID"
	ColumnBSSIDWidth = 17
	ColumnSSIDTitle  = "Network"
	ColumnSSIDWidth  = 20
	ColumnChanTitle  = "Chan"
	ColumnChanWidth  = 6
	ColumnWidthTitle = "Width"
	ColumnWidthWidth = 7
	ColumnBandTitle  = "Band"
	ColumnBandWidth  = 6
	ColumnRSSITitle  = "RSSI"
	ColumnRSSIWidth  = 6
	ColumnNoiseTitle = "Noise"
	ColumnNoiseWidth = 7
	ColumnSNRTitle   = "SNR"
	ColumnSNRWidth   = 5
)

// Column Definition.
type ColumnDef struct {
	table.Column
	Num  int // starts from 1
	Sort SortDef
}

// Returns column defintion.
func NewColumnDef(num int, col table.Column) *ColumnDef {
	return &ColumnDef{
		Num:    num,
		Sort:   ColumnSorterGenerator(col.Title),
		Column: col,
	}
}

// Returns 0 based column index.
func (def ColumnDef) Index() int {
	return def.Num - 1
}

// Slice of Column Definitions.
type ColumnDefs []*ColumnDef

// Map of Column Definitions.
type ColumnDefsByName map[string]*ColumnDef

// Returns @table.Column slice with applies sorting in column view.
func (defs ColumnDefs) SortBy(def *ColumnDef) Columns {
	cols := defs.ToColumns()

	cols[def.Index()] = table.Column{
		Title: fmt.Sprintf("%s %s", def.Title, def.Sort.ord),
		Width: def.Width,
	}

	return cols
}

// Returns @table.Column slice.
func (defs ColumnDefs) ToColumns() Columns {
	cols := make(Columns, len(defs))

	for i, col := range defs {
		cols[i] = col.Column
	}

	return cols
}

// Alias for table.Column slice.
type Columns []table.Column

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

// Returns default columns definitions.
func getDefaultCols() (ColumnDefs, ColumnDefsByName) {
	cols := []table.Column{
		{Title: ColumnBSSIDTitle, Width: ColumnBSSIDWidth},
		{Title: ColumnSSIDTitle, Width: ColumnSSIDWidth},
		{Title: ColumnChanTitle, Width: ColumnChanWidth},
		{Title: ColumnWidthTitle, Width: ColumnWidthWidth},
		{Title: ColumnBandTitle, Width: ColumnBandWidth},
		{Title: ColumnRSSITitle, Width: ColumnRSSIWidth},
		{Title: ColumnNoiseTitle, Width: ColumnNoiseWidth},
		{Title: ColumnSNRTitle, Width: ColumnSNRWidth},
	}

	defs := make(ColumnDefs, len(cols))
	defsByName := make(ColumnDefsByName, len(cols))
	for i, col := range cols {
		def := NewColumnDef(i+1, col)
		defs[i] = def
		defsByName[col.Title] = def
	}

	return defs, defsByName
}

// Reponsible for table view.
type View struct {
	table         table.Model
	cols          ColumnDefs
	colsByNames   ColumnDefsByName
	viewportStyle lipgloss.Style
	tableStyles   table.Styles
}

// Returns new table view.
func NewView() *View {
	cols, colsByName := getDefaultCols()
	viewportStyle := getDefaultViewportStyle()
	tableStyles := getDefaultTableStyles()

	t := table.New(
		table.WithColumns(cols.ToColumns()),
		table.WithHeight(defaultTableHeight),
		table.WithFocused(true),
	)
	t.SetStyles(tableStyles)

	t.UpdateViewport()

	return &View{
		table:         t,
		cols:          cols,
		colsByNames:   colsByName,
		viewportStyle: viewportStyle,
		tableStyles:   tableStyles,
	}
}

// Renderes table.
func (v *View) View() string {
	return v.viewportStyle.Render(v.table.View()) + "\n"
}

// Updates rows in table.
// Does not invoke table redraw.
func (v *View) OnData(rows []table.Row) {
	v.table.SetRows(rows)
}

// Applies sort view in table header.
func (v *View) OnSort(def ColumnDef) {
	v.cols[def.Index()] = &def
	columns := v.cols.SortBy(&def)
	v.table.SetColumns(columns)
}

// Returns copy of column definition by number.
func (v *View) GetColumnByNum(num int) (ColumnDef, bool) {
	idx := num - 1
	if idx < 0 || idx >= len(v.cols) {
		return ColumnDef{}, false
	}
	col := v.cols[idx]
	return *col, true
}

// Returns copy of column definition by title.
func (v *View) GetColumnByTitle(title string) (ColumnDef, bool) {
	col, found := v.colsByNames[title]
	return *col, found
}
