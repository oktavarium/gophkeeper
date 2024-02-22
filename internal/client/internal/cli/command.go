package cli

import (
	tea "github.com/charmbracelet/bubbletea"
)

type stateCmd state
type errorCmd error
type msgCmd string
type resetCmd struct{}
type actionCmd struct{}
type serverAddrCmd string

type loginCmd struct {
	login    string
	password string
}

type registerCmd struct {
	login    string
	password string
}

type saveCmd struct {
	name string
	data string
}

func changeState(st state) tea.Cmd {
	return func() tea.Msg {
		return stateCmd(st)
	}
}

func makeLogin(login, password string) tea.Cmd {
	return func() tea.Msg {
		return loginCmd{
			login:    login,
			password: password,
		}
	}
}

func makeRegister(login, password string) tea.Cmd {
	return func() tea.Msg {
		return registerCmd{
			login:    login,
			password: password,
		}
	}
}

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return errorCmd(err)
	}
}

func makeSaveData(name, data string) tea.Cmd {
	return func() tea.Msg {
		return saveCmd{
			name: name,
			data: data,
		}
	}
}

func saveServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return serverAddrCmd(addr)
	}
}

func makeMsg(m string) tea.Cmd {
	return func() tea.Msg {
		return msgCmd(m)
	}
}

func makeReset() tea.Msg {
	return resetCmd{}
}

func makeAction() tea.Msg {
	return actionCmd{}
}
