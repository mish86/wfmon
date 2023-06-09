package spectrum

import (
	"reflect"
	netdata "wfmon/pkg/data/net"
	log "wfmon/pkg/logger"
	"wfmon/pkg/widgets/color"
	"wfmon/pkg/widgets/events"

	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		// cmd  tea.Cmd
		cmds []tea.Cmd
	)

	var changeViewOnSelect = func() {
		for i := range m.data {
			if m.selected.Compare(m.data[i].Key) == 0 {
				m.SetBandView(m.data[i].Band)
			}
		}
	}

	switch msg := msg.(type) {
	case events.ToggledNetworkKeyMsg:
		// nothing

	case events.SelectedNetworkKeyMsg:
		m.selected = msg.Key
		// auto-change band view for selected network
		changeViewOnSelect()
		// render data to viewport
		m.refresh()

	case events.NetworkKeyMsg:
		m.selected = msg.Key
		// auto-change band view for selected network
		// changeViewOnSelect()
		// render data to viewport
		m.refresh()

	case events.SignalFieldMsg:
		WithSignalField(msg)(m)

		m.refresh()

	case events.TableWidthMsg:
		m.SetWidth(int(msg))
		m.refresh()

	case events.NetworksOnScreen:
		// copy slices from origin message
		nets := make(netdata.Slice, len(msg.Networks))
		copy(nets, msg.Networks)
		colors := make([]color.HexColor, len(msg.Colors))
		copy(colors, msg.Colors)
		if len(colors) < len(nets) {
			log.Warnf("%s mailformed len(colors) < len(networks)", reflect.TypeOf(msg))
			colors = append(colors, make([]color.HexColor, len(nets)-len(colors))...) //nolint:makezero // ignore
		}

		// convert to waves to display
		m.data = MultiWaver{
			NetworksOnScreen: events.NetworksOnScreen{
				Networks: nets,
				Colors:   colors,
			},
			fieldKey: m.fieldKey,
			ts:       m.dataSource,
		}.Waves()

		// auto-change band view for selected network
		// changeViewOnSelect()
		// render data to viewport
		m.refresh()
	}

	// Bubble up the cmds
	return m, tea.Batch(cmds...)
}
