package settingsmodel

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SaveServerAddrMsg string

func SaveServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return SaveServerAddrMsg(addr)
	}
}

type SetServerAddrMsg string

func SetServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return SetServerAddrMsg(addr)
	}
}

type BackMsg struct{}

func BackCmd() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}
