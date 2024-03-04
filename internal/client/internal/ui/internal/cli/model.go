package cli

type Model interface {
	View() string
	Reset()
	Focus()
	Blur()
	Focused() bool
}
