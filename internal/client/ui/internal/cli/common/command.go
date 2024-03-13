package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type StateMsg State
type ErrorMsg error
type MsgMsg string
type ResetMsg struct{}

func ChangeState(st State) tea.Cmd {
	return func() tea.Msg {
		return StateMsg(st)
	}
}

func MakeError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

func MakeMsg(m string) tea.Cmd {
	return func() tea.Msg {
		return MsgMsg(m)
	}
}

func MakeReset() tea.Msg {
	return ResetMsg{}
}

type LoginStoreMsg string

func LoginStore(pass string) tea.Cmd {
	return func() tea.Msg {
		return LoginStoreMsg(pass)
	}
}
