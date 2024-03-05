package common

type Model interface {
	View() string

	Reset()
	Focus()
	Blur()
	Focused() bool
}
