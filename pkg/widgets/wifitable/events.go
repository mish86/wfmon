package wifitable

import "wfmon/pkg/widgets/events"

func SignalFieldMsges() map[string]events.SignalFieldMsg {
	return map[string]events.SignalFieldMsg{
		RSSIKey:    RSSIFieldMsg(),
		QualityKey: QualityFieldMsg(),
		BarsKey:    BarsFieldMsg(),
	}
}

func RSSIFieldMsg() events.SignalFieldMsg {
	return events.SignalFieldMsg{
		Key:    RSSIKey,
		MinVal: -100,
		MaxVal: 0,
	}
}

func QualityFieldMsg() events.SignalFieldMsg {
	return events.SignalFieldMsg{
		Key:    QualityKey,
		MinVal: 0,
		MaxVal: 100,
	}
}

func BarsFieldMsg() events.SignalFieldMsg {
	return events.SignalFieldMsg{
		Key:    QualityKey,
		MinVal: 0,
		MaxVal: 100,
	}
}
