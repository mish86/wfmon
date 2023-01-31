package wifitable

import (
	"context"
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

	help     *help.Model
	showHelp bool
	keys     *KeyMap

	data            *TableData
	view            *View
	cursor          *NetworkTableKey
	cursorLock      sync.Mutex
	sort            SortDef
	refreshInterval time.Duration
}

type RefreshMsg time.Time

type KeyMap struct {
	table.KeyMap
	Sort      key.Binding
	ResetSort key.Binding
	Help      key.Binding
	Quit      key.Binding
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		KeyMap: table.DefaultKeyMap(),
		Sort: key.NewBinding(
			key.WithKeys("1", "2", "3", "4", "5", "6", "7", "8"),
			key.WithHelp("[1:8]", "sort"),
		),
		ResetSort: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "reset sort"),
		),
		Help: key.NewBinding(
			key.WithKeys("h", "?"),
			key.WithHelp("h", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k *KeyMap) MoveBindings() []key.Binding {
	return []key.Binding{k.LineUp, k.LineDown, k.HalfPageUp, k.HalfPageDown, k.PageUp, k.PageDown, k.GotoTop, k.GotoBottom}
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.MoveBindings(),
		{k.Sort, k.ResetSort},
		{k.Help, k.Quit},
	}
}

func NewTableCtrl(frames <-chan wifi.Frame) *TableCtrl {
	data := NewTableData()
	view := NewView()
	help := help.New()
	help.ShowAll = true
	keys := NewKeyMap()
	view.table.KeyMap = keys.KeyMap

	// default sorting
	col, _ := view.GetColumnByTitle(ColumnSSIDTitle)
	view.OnSort(col)

	ctrl := &TableCtrl{
		framesCh:        frames,
		data:            data,
		view:            view,
		help:            &help,
		keys:            keys,
		sort:            col.Sort,
		refreshInterval: defaultRefreshInterval,
	}

	return ctrl
}

// Starts processing incomming frames from packets.
func (c *TableCtrl) Start(ctx context.Context) error {
	c.ctx, c.stop = context.WithCancel(ctx)

	for {
		select {
		case <-c.ctx.Done():
			return nil

		case frame := <-c.framesCh:
			data := &NetworkData{
				BSSID:       frame.BSSID.String(),
				NetworkName: frame.SSID,
				RSSI:        frame.RSSI,
				Noise:       frame.Noise,
				SNR:         frame.RSSI - frame.Noise,
			}

			if frame.Channel != 0 {
				data.Channel = int(frame.Channel)
				data.Band = wifi.GetBandByChan(data.Channel).String()
			}

			if frame.SupportedChannelWidth != 0 {
				data.ChannelWidth = int(frame.SupportedChannelWidth)
			}

			c.data.Add(data)

		default:
			continue
		}
	}
}

// Returns sorted rows to redraw tick.
func (c *TableCtrl) GetRows() ([]table.Row, int) {
	networks := c.data.GetNetworkSlice()

	sort.Sort(c.sort.Sorter(networks))

	// lock cursor update
	c.cursorLock.Lock()
	defer c.cursorLock.Unlock()

	if c.cursor == nil && len(networks) > 0 {
		c.cursor = &NetworkTableKey{networks[0].BSSID, networks[0].NetworkName}
	}

	rows := make([]table.Row, len(networks))
	cursor := 0
	for rowNum, data := range networks {
		key := NetworkTableKey{data.BSSID, data.NetworkName}
		if key == *c.cursor {
			cursor = rowNum
		}
		rows[rowNum] = table.Row{
			data.BSSID,
			data.NetworkName,
			strconv.Itoa(data.Channel),
			strconv.Itoa(data.ChannelWidth),
			data.Band,
			strconv.Itoa(int(data.RSSI)),
			strconv.Itoa(int(data.Noise)),
			strconv.Itoa(int(data.SNR)),
		}
	}

	return rows, cursor
}

// Inits redraw tick.
func (c *TableCtrl) Init() tea.Cmd {
	return tea.Batch(
		c.redrawTick(),
	)
}

func (c *TableCtrl) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cursorChanged bool
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keys.MoveBindings()...):
			cursorChanged = true
		case key.Matches(msg, c.keys.ResetSort):
			// clear flag
			prevCol, _ := c.view.GetColumnByTitle(c.sort.col)
			prevCol.Sort.ord = None
			c.view.OnSort(prevCol)

			// restore default sorting
			col, _ := c.view.GetColumnByTitle(ColumnSSIDTitle)
			c.sort = col.Sort
			c.view.OnSort(col)
		case key.Matches(msg, c.keys.Sort):
			num, _ := strconv.Atoi(msg.String())
			col, _ := c.view.GetColumnByNum(num)
			col.Sort.ChangeOrder()
			if c.sort.col != col.Sort.col {
				prevCol, _ := c.view.GetColumnByTitle(c.sort.col)
				prevCol.Sort.ord = None
				c.view.OnSort(prevCol)
			}
			c.sort = col.Sort
			c.view.OnSort(col)
		case key.Matches(msg, c.keys.Help):
			c.showHelp = !c.showHelp
		case key.Matches(msg, c.keys.Quit):
			return c, tea.Quit
		}
	case RefreshMsg:
		//  Return redraw tick command again to loop.
		return c, c.redrawTick()
	default:
		log.Debugf("msg: %v", msg)
	}

	// traverse message down to table view
	c.view.table, cmd = c.view.table.Update(msg)

	// read updated cursor position from view
	if cursorChanged {
		c.cursorLock.Lock()
		// cursor := c.view.table.Cursor()
		row := c.view.table.SelectedRow()
		c.cursor = &NetworkTableKey{row[0], row[1]}
		c.cursorLock.Unlock()
	}

	return c, cmd
}

// Renders table view.
func (c *TableCtrl) View() string {
	if c.showHelp {
		return c.help.View(c.keys)
	}
	return c.view.View()
}

// Invokes update tea to redraw table by interval.
func (c *TableCtrl) redrawTick() tea.Cmd {
	return tea.Tick(c.refreshInterval, func(t time.Time) tea.Msg {
		rows, cursor := c.GetRows()
		c.view.OnData(rows)

		if cursor > 0 && cursor < len(rows) {
			c.view.table.SetCursor(cursor)
		}

		return RefreshMsg(t)
	})
}
