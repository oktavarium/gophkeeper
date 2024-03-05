package common

import (
	tea "github.com/charmbracelet/bubbletea"
)

type StateMsg State
type ErrorMsg error
type MsgMsg string
type ResetMsg struct{}
type SaveServerAddrMsg string
type SetServerAddrMsg string
type DeleteCardMsg string
type SyncMsg struct{}

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

func SaveServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return SaveServerAddrMsg(addr)
	}
}

func SetServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return SetServerAddrMsg(addr)
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

func DeleteCard(id string) tea.Cmd {
	return func() tea.Msg {
		return DeleteCardMsg(id)
	}
}

func Sync() tea.Cmd {
	return func() tea.Msg {
		return SyncMsg{}
	}
}
