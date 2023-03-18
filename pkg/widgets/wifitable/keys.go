package wifitable

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/evertras/bubble-table/table"
)

type KeyMap struct {
	table.KeyMap
	PageUp      key.Binding
	PageDown    key.Binding
	RowUp       key.Binding
	RowDown     key.Binding
	GotoTop     key.Binding
	GotoBottom  key.Binding
	SignalView  key.Binding
	StationView key.Binding
	Sort        key.Binding
	Reset       key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		KeyMap: table.KeyMap{
			ScrollLeft: key.NewBinding(
				key.WithKeys("shift+left"),
				key.WithHelp("shift+←", "scroll left"),
			),
			ScrollRight: key.NewBinding(
				key.WithKeys("shift+right"),
				key.WithHelp("shift+→", "scroll right"),
			),
			RowSelectToggle: key.NewBinding(
				key.WithKeys(" ", "enter"),
				key.WithHelp("⏎", "select row"),
			),
		},
		PageUp: key.NewBinding(
			key.WithKeys("left", "b", "pgup"),
			key.WithHelp("←/b/pgup", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("right", "f", "pgdown"),
			key.WithHelp("→/f/pgdown", "page down"),
		),
		GotoTop: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GotoBottom: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		RowUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		RowDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		StationView: key.NewBinding(
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
		Reset: key.NewBinding(
			key.WithKeys("0"),
			key.WithHelp("0", "reset view"),
		),
	}
}

func (k *KeyMap) MoveBindings() []key.Binding {
	return []key.Binding{k.RowUp, k.RowDown, k.PageUp, k.PageDown, k.ScrollLeft, k.ScrollRight, k.GotoTop, k.GotoBottom}
}

func (k *KeyMap) ViewBindings() []key.Binding {
	return []key.Binding{k.Sort, k.Reset, k.StationView, k.SignalView, k.RowSelectToggle}
}
