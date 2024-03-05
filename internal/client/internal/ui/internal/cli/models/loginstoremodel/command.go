package loginstoremodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LoginStoreMsg string

func LoginStore(pass string) tea.Cmd {
	return func() tea.Msg {
		return LoginStoreMsg(pass)
	}
}
