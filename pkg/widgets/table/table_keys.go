package wifitable

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
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

func NewKeyMap() *KeyMap {
	return &KeyMap{
		KeyMap: table.DefaultKeyMap(),
		SignalView: key.NewBinding(
			key.WithKeys("ctrl+^"),
			key.WithHelp("ctrl+^", "swap RSSI/Quality/Bars"),
		),
		BSSIDView: key.NewBinding(
			key.WithKeys("ctrl+@"),
			key.WithHelp("ctrl+@", "swap BSSID/Vendor"),
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
	return []key.Binding{k.LineUp, k.LineDown, k.HalfPageUp, k.HalfPageDown, k.PageUp, k.PageDown, k.GotoTop, k.GotoBottom}
}

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Sort, k.ResetSort, k.SignalView, k.BSSIDView, k.Help, k.Quit}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.MoveBindings(),
		{k.Sort, k.ResetSort, k.SignalView, k.BSSIDView},
		{k.Help, k.Quit},
	}
}
