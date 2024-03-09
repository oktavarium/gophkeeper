package registermodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type RegisterMsg struct {
	Login    string
	Password string
}

func MakeRegister(login, password string) tea.Cmd {
	return func() tea.Msg {
		return RegisterMsg{
			Login:    login,
			Password: password,
		}
	}
}
