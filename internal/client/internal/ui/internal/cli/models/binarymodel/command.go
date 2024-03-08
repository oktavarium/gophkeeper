package binarymodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type NewFileCmd struct {
	CurrentID string
	Name      string
	Path      string
}

func NewFile(currentID, name, path string) tea.Cmd {
	return func() tea.Msg {
		return NewFileCmd{
			CurrentID: currentID,
			Name:      name,
			Path:      path,
		}
	}
}
