package simpledata

import (
	tea "github.com/charmbracelet/bubbletea"
)

type ErrorMsg error

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

type BackMsg struct{}

func Back() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}

type NewSimpleDataMsg struct {
	CurrentID string
	Name      string
	Data      string
}

func NewSimpleData(currentID, name, data string) tea.Cmd {
	return func() tea.Msg {
		return NewSimpleDataMsg{
			CurrentID: currentID,
			Name:      name,
			Data:      data,
		}
	}
}
