package wifitable

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
)

type ColumnOrder interface {
	Sort() Sort
	SetOrder(ord Order)
	SwapOrder()
}

// Column Viewer with number in table, title, width and sorting order.
type ColumnViewer interface {
	Index() int
	Title() string
	Width() int
}

// Swaper of column view.
type ColumnViewSwaper interface {
	Next() ColumnViewer
	Prev() ColumnViewer
}

// Column View with number in table, title, width and sorting order.
type ColumnView struct {
	table.Column
	num  int // starts from 1
	sort Sort
}

// Returns column view.
func NewColumnView(num int, col table.Column) *ColumnView {
	view := &ColumnView{
		num:    num,
		sort:   ColumnSorterGenerator(col.Title),
		Column: col,
	}

	return view
}

// Returns 0 based column index.
func (view *ColumnView) Index() int {
	return view.num - 1
}

// Returns column title.
func (view *ColumnView) Title() string {
	return view.Column.Title
}

// returns column width.
func (view *ColumnView) Width() int {
	return view.Column.Width
}

// Returns column sort.
func (view *ColumnView) Sort() Sort {
	return view.sort
}

// Sets sorting order.
func (view *ColumnView) SetOrder(ord Order) {
	view.sort.SetOrder(ord)
}

// Swaps sorting order.
func (view *ColumnView) SwapOrder() {
	view.sort.SwapOrder()
}

// Swaper of column view.
type ColumnSwapper interface {
	Titles() []string
	Next() ColumnViewer
	Prev() ColumnViewer
}

// Aggregates several column views in one.
type MultiColumnView struct {
	selected int
	cols     []*ColumnView
}

// Returns new swappable column view.
// Requires column number in table and columns to swap.
func NewMultiColumnView(num int, columns ...table.Column) *MultiColumnView {
	cols := make([]*ColumnView, len(columns))

	view := &MultiColumnView{
		selected: 0,
		cols:     cols,
	}

	for i, col := range columns {
		view.cols[i] = NewColumnView(num, col)
	}

	return view
}

// Next column viewer.
func (view *MultiColumnView) Next() ColumnViewer {
	view.selected++
	if view.selected >= len(view.cols) {
		view.selected = 0
	}

	return view.cols[view.selected]
}

// Previous column viewer.
func (view *MultiColumnView) Prev() ColumnViewer {
	view.selected--
	if view.selected < 0 {
		view.selected = len(view.cols) - 1
	}

	return view.cols[view.selected]
}

func (view *MultiColumnView) Titles() []string {
	names := make([]string, len(view.cols))

	for i, col := range view.cols {
		names[i] = col.Title()
	}

	return names
}

// func (view *MultiColumnView) columns() []*ColumnView {
// return view.cols
// }

// Returns 0 based column index.
func (view *MultiColumnView) Index() int {
	return view.cols[view.selected].Index()
}

// Returns column title.
func (view *MultiColumnView) Title() string {
	return view.cols[view.selected].Column.Title
}

// returns column width.
func (view *MultiColumnView) Width() int {
	return view.cols[view.selected].Column.Width
}

// Returns column sort.
func (view *MultiColumnView) Sort() Sort {
	return view.cols[view.selected].sort
}

// Sets sorting order.
func (view *MultiColumnView) SetOrder(ord Order) {
	view.cols[view.selected].SetOrder(ord)
}

// Swaps sorting order.
func (view *MultiColumnView) SwapOrder() {
	view.cols[view.selected].SwapOrder()
}

// Slice of Column Views.
type ColumnViewSlice []ColumnViewer

// Map of Column Views.
type ColumnViewMap map[string]ColumnViewer

// Alias for @table.Column slice.
type TableColumns []table.Column

// Returns @table.Column slice with applied sorting in title.
func (cols ColumnViewSlice) SortBy(col ColumnViewer) TableColumns {
	columns := cols.TableColumns()

	title := col.Title()
	sort := col.(ColumnOrder).Sort()
	columns[col.Index()] = table.Column{
		Title: fmt.Sprintf("%s %s", title, sort.Order()),
		Width: col.Width(),
	}

	return columns
}

// Returns @table.Column slice.
func (cols ColumnViewSlice) TableColumns() TableColumns {
	columns := make(TableColumns, len(cols))

	for i, col := range cols {
		columns[i] = table.Column{
			Title: col.Title(),
			Width: col.Width(),
		}
	}

	return columns
}

func NewColumnViewer(num int, cols ...table.Column) ColumnViewer {
	if len(cols) > 1 {
		return NewMultiColumnView(num, cols...)
	}

	return NewColumnView(num, cols[0])
}
