package wifitable

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"
	log "wfmon/pkg/logger"
	"wfmon/pkg/wifi"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	defaultRefreshInterval = time.Second
)

type TableCtrl struct {
	ctx      context.Context
	stop     context.CancelFunc
	framesCh <-chan wifi.Frame

	help      *help.Model
	helpShown bool
	keys      *KeyMap

	data            *TableData
	getters         map[string]FncRowGetter
	view            *TableView
	cursor          *NetworkTableKey
	cursorLock      sync.Mutex
	sort            Sort
	refreshInterval time.Duration
}

type RefreshMsg time.Time

// Returns new network table controller.
func NewTableCtrl(frames <-chan wifi.Frame) *TableCtrl {
	columns := GenerateDefaultColumns()
	getters := GenerateRowGetters()

	data := NewTableData()
	view := NewTableView(columns)
	keys := NewKeyMap()
	help := help.New()
	help.ShowAll = true
	view.table.KeyMap = keys.KeyMap

	ctrl := &TableCtrl{
		framesCh:        frames,
		data:            data,
		getters:         getters,
		view:            view,
		help:            &help,
		keys:            keys,
		refreshInterval: defaultRefreshInterval,
	}

	// default sorting
	ctrl.sortByTitle(ColumnSSIDTitle, None)

	return ctrl
}

// Starts processing incomming frames from packets.
func (c *TableCtrl) Start(ctx context.Context) error {
	c.ctx, c.stop = context.WithCancel(ctx)

	for {
		select {
		case frame, ok := <-c.framesCh:
			if !ok {
				return fmt.Errorf("frames source closed, stopping updating table")
			}

			data := NetworkDataConverter(frame).NetworkData()
			c.data.Add(data)

		case <-c.ctx.Done():
			return nil
		}
	}
}

// Network data field getter. Returns string value to pass in tea.Row.
type FncRowGetter func(data *NetworkData) string

// Returns sorted rows to redraw tick.
func (c *TableCtrl) GetRows() ([]table.Row, int) {
	networks := c.data.NetworkSlice()

	// Sorts in ASC/DESC order depending on current Order value.
	sort.Sort(c.sort.Sorter(networks))

	// lock cursor update
	c.cursorLock.Lock()
	defer c.cursorLock.Unlock()

	if c.cursor == nil && len(networks) > 0 {
		c.cursor = networks[0].Key()
	}

	rows := make([]table.Row, len(networks))
	cursor := 0
	for rowNum, data := range networks {
		d := data
		key := d.Key()
		if key.Compare(c.cursor) == 0 {
			cursor = rowNum
		}

		cols := c.view.TableColumns()
		row := make([]string, len(cols))

		for i, col := range cols {
			if getter, found := c.getters[col.Title]; !found {
				log.Errorf("data source for %s not found", col.Title)
			} else {
				row[i] = getter(&d)
			}
		}
		rows[rowNum] = row
	}

	return rows, cursor
}

// Inits redraw tick.
func (c *TableCtrl) Init() tea.Cmd {
	return tea.Batch(
		c.redrawTick(),
	)
}

// Process @tea.Msg.
func (c *TableCtrl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cursorChanged bool
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keys.MoveBindings()...):
			cursorChanged = true
		case key.Matches(msg, c.keys.SignalView):
			// swap signal view
			c.swapSignalView()
		case key.Matches(msg, c.keys.ResetSort):
			// clear sorting flag
			c.sortByTitle(c.sort.Title(), None)
			// restore default sorting
			c.sortByTitle(ColumnSSIDTitle, None)
		case key.Matches(msg, c.keys.Sort):
			c.swapSortByNum(msg.String())
		case key.Matches(msg, c.keys.Help):
			c.helpShown = !c.helpShown
		case key.Matches(msg, c.keys.Quit):
			return c, tea.Quit
		}
	case RefreshMsg:
		//  Return redraw tick command again to loop.
		// get []@table.Row presentation
		rows, cursor := c.GetRows()

		// set []@table.Row in @table.Model
		c.view.OnData(rows)

		// update cursor position in @table.Model
		if cursor > 0 && cursor < len(rows) {
			c.view.table.SetCursor(cursor)
		}

		return c, c.redrawTick()
	default:
		log.Debugf("msg: %v", msg)
	}

	// traverse message down to table view
	c.view.table, cmd = c.view.table.Update(msg)

	// read updated cursor position from view
	if cursorChanged {
		c.onCursortChanged()
	}

	return c, cmd
}

// Renders table view.
func (c *TableCtrl) View() string {
	if c.helpShown {
		return c.help.View(c.keys)
	}
	return c.view.View()
}

// Invokes update tea to redraw table by interval.
func (c *TableCtrl) redrawTick() tea.Cmd {
	return tea.Tick(c.refreshInterval, func(t time.Time) tea.Msg {
		return RefreshMsg(t)
	})
}

// Set sorting order for a column by display title.
func (c *TableCtrl) sortByTitle(title string, ord Order) {
	var col ColumnViewer
	var ok bool
	// clear sorting flag
	if col, ok = c.view.GetColumnByTitle(title); !ok {
		log.Errorf("failed to get column with title %s", title)
		return
	}

	col.(ColumnOrder).SetOrder(ord)
	c.sort = col.(ColumnOrder).Sort()
	c.view.OnSort(col)
}

// Change sorting order for a column by display number.
func (c *TableCtrl) swapSortByNum(sNum string) {
	var col ColumnViewer
	var num int
	var err error
	var ok bool

	if num, err = strconv.Atoi(sNum); err != nil {
		log.Errorf("failed to get column with number %s", sNum)
		return
	}
	if col, ok = c.view.GetColumnByNum(num); !ok {
		log.Errorf("failed to get column with number %s", sNum)
		return
	}

	// reset current sorting
	if c.sort.Title() != col.Title() {
		c.sortByTitle(c.sort.Title(), None)
	}

	// swaps sorting order
	col.(ColumnOrder).SwapOrder()
	sort := col.(ColumnOrder).Sort()
	c.sortByTitle(col.Title(), sort.Order())
}

// Swaps RSSI/Quality/Bars.
func (c *TableCtrl) swapSignalView() {
	var col ColumnViewer
	var swapper ColumnViewSwaper
	var ok bool

	// get current column view
	if col, ok = c.view.GetColumnByTitle(ColumnRSSITitle); !ok {
		log.Errorf("failed to get column with title %s", ColumnRSSITitle)
		return
	}

	// get sort of current column view
	sort := col.(ColumnOrder).Sort()

	// get swapper for column
	if swapper, ok = col.(ColumnViewSwaper); !ok {
		log.Errorf("failed to swap %s column view", ColumnRSSITitle)
		return
	}

	// swap column view
	col = swapper.Next()
	// restore sorting order
	col.(ColumnOrder).SetOrder(sort.Order())
	// store sorting information
	c.sort = col.(ColumnOrder).Sort()
	// update headers view
	c.view.OnSort(col)
}

// Reads updated cursor position from view.
func (c *TableCtrl) onCursortChanged() {
	c.cursorLock.Lock()
	defer c.cursorLock.Unlock()

	if c.view.table.Cursor() < 0 || c.view.table.Cursor() >= len(c.view.table.Rows()) {
		return
	}

	row := c.selectedRow()
	bssIDCol, _ := c.view.GetColumnByTitle(ColumnBSSIDTitle)
	ssIDCol, _ := c.view.GetColumnByTitle(ColumnSSIDTitle)
	c.cursor = NewNetworkTableKey(row[bssIDCol.Index()], row[ssIDCol.Index()])
}

// Returns selected row.
func (c *TableCtrl) selectedRow() table.Row {
	return c.view.table.SelectedRow()
}
