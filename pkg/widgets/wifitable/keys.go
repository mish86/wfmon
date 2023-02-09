package wifitable

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/evertras/bubble-table/table"
)

type KeyMap struct {
	table.KeyMap
	SignalView key.Binding
	BSSIDView  key.Binding
	Sort       key.Binding
	ResetSort  key.Binding
	Help       key.Binding
	Quit       key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		KeyMap: table.KeyMap{
			RowUp: key.NewBinding(
				key.WithKeys("up", "k"),
				key.WithHelp("↑/k", "up"),
			),
			RowDown: key.NewBinding(
				key.WithKeys("down", "j"),
				key.WithHelp("↓/j", "down"),
			),
			PageUp: key.NewBinding(
				key.WithKeys("left", "b", "pgup"),
				key.WithHelp("←/b/pgup", "page up"),
			),
			PageDown: key.NewBinding(
				key.WithKeys("right", "f", "pgdown"),
				key.WithHelp("→/f/pgdown", "page down"),
			),
			PageFirst: key.NewBinding(
				key.WithKeys("home", "g"),
				key.WithHelp("g/home", "go to start"),
			),
			PageLast: key.NewBinding(
				key.WithKeys("end", "G"),
				key.WithHelp("G/end", "go to end"),
			),
		},
		BSSIDView: key.NewBinding(
			key.WithKeys("ctrl+@"),
			key.WithHelp("ctrl+@", "swap BSSID/Vendor"),
		),
		SignalView: key.NewBinding(
			key.WithKeys("ctrl+^"),
			key.WithHelp("ctrl+^", "swap RSSI/Quality/Bars"),
		),
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
	return []key.Binding{k.RowUp, k.RowDown, k.PageUp, k.PageDown, k.PageFirst, k.PageLast}
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Sort, k.ResetSort, k.BSSIDView, k.SignalView, k.Help, k.Quit}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.MoveBindings(),
		{k.Sort, k.ResetSort, k.BSSIDView, k.SignalView},
		{k.Help, k.Quit},
	}
}
