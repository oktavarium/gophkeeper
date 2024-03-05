package loginmodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type LoginMsg struct {
	Login    string
	Password string
}

func MakeLogin(login, password string) tea.Cmd {
	return func() tea.Msg {
		return LoginMsg{
			Login:    login,
			Password: password,
		}
	}
}
