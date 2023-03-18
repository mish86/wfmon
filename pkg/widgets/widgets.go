package widgets

type WithTitle interface {
	Title() string
}

type WithFocus interface {
	Focused(focus bool)
	GetFocused() bool
}
