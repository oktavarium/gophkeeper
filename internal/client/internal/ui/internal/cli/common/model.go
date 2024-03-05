package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Model interface {
	View() string
	Update(tea.Msg) tea.Cmd
	Reset()
	Focus()
	Blur()
	Focused() bool
}
