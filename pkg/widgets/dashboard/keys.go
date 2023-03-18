package dashboard

import (
	"wfmon/pkg/widgets/wifitable"

	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	TableKeyMap wifitable.KeyMap
	Spectrum    key.Binding
	Sparkline   key.Binding
	Help        key.Binding
	Quit        key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		TableKeyMap: wifitable.NewKeyMap(),
		Spectrum: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "wifi spectrum"),
		),
		Sparkline: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "signal sparkline"),
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

func (k *KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.TableKeyMap.Sort,
		k.TableKeyMap.Reset,
		k.TableKeyMap.StationView,
		k.TableKeyMap.SignalView,
		k.Spectrum,
		k.Sparkline,
		k.Help,
		k.Quit,
	}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		k.TableKeyMap.MoveBindings(),
		k.TableKeyMap.ViewBindings(),
		{k.Spectrum, k.Sparkline},
		{k.Help, k.Quit},
	}
}
