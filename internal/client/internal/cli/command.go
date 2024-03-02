package cli

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type checkStoreMsg struct{}
type loginLocalStoreMsg string
type stateMsg state
type errorMsg error
type msgMsg string
type resetMsg struct{}
type saveServerAddrMsg string
type setServerAddrMsg string
type deleteCardMsg string
type syncMsg struct{}

type loginMsg struct {
	login    string
	password string
}

type registerMsg struct {
	login    string
	password string
}

func changeState(st state) tea.Cmd {
	return func() tea.Msg {
		return stateMsg(st)
	}
}

func makeLogin(login, password string) tea.Cmd {
	return func() tea.Msg {
		return loginMsg{
			login:    login,
			password: password,
		}
	}
}

func makeRegister(login, password string) tea.Cmd {
	return func() tea.Msg {
		return registerMsg{
			login:    login,
			password: password,
		}
	}
}

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return errorMsg(err)
	}
}

func saveServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return saveServerAddrMsg(addr)
	}
}

func setServerAddr(addr string) tea.Cmd {
	return func() tea.Msg {
		return setServerAddrMsg(addr)
	}
}

func makeMsg(m string) tea.Cmd {
	return func() tea.Msg {
		return msgMsg(m)
	}
}

func makeReset() tea.Msg {
	return resetMsg{}
}

func checkStore() tea.Msg {
	return checkStoreMsg{}
}

func loginLocalStore(pass string) tea.Cmd {
	return func() tea.Msg {
		return loginLocalStoreMsg(pass)
	}
}

type newCardCmd struct {
	CurrentCardID string
	Name          string
	Ccn           string
	Exp           time.Time
	CVV           uint32
}

func newCard(currentCardID, name, ccn string, exp time.Time, cvv uint32) tea.Cmd {
	return func() tea.Msg {
		return newCardCmd{
			CurrentCardID: currentCardID,
			Name:          name,
			Ccn:           ccn,
			Exp:           exp,
			CVV:           cvv,
		}
	}
}

func deleteCard(id string) tea.Cmd {
	return func() tea.Msg {
		return deleteCardMsg(id)
	}
}

func sync() tea.Cmd {
	return func() tea.Msg {
		return syncMsg{}
	}
}
