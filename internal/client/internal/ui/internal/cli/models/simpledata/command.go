package simpledata

import (
	tea "github.com/charmbracelet/bubbletea"
)

type (
	ErrorMsg         error
	NewSimpleDataCmd struct {
		CurrentID string
		Name      string
		Data      string
	}
	BlureCmd struct{}
)

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

func makeBlur() tea.Cmd {
	return func() tea.Msg {
		return BlureCmd{}
	}
}

func NewSimpleData(currentID, name, data string) tea.Cmd {
	return func() tea.Msg {
		return NewSimpleDataCmd{
			CurrentID: currentID,
			Name:      name,
			Data:      data,
		}
	}
}
