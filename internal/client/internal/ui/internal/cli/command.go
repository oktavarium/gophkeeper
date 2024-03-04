package cli

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type LoginStoreMsg string
type StateMsg State
type ErrorMsg error
type MsgMsg string
type ResetMsg struct{}
type SaveServerAddrMsg string
type SetServerAddrMsg string
type DeleteCardMsg string
type SyncMsg struct{}

type LoginMsg struct {
	Login    string
	Password string
}

type RegisterMsg struct {
	Login    string
	Password string
}

func ChangeState(st State) tea.Cmd {
	return func() tea.Msg {
		return StateMsg(st)
	}
}

func MakeLogin(login, password string) tea.Cmd {
	return func() tea.Msg {
		return LoginMsg{
			Login:    login,
			Password: password,
		}
	}
}

func MakeRegister(login, password string) tea.Cmd {
	return func() tea.Msg {
		return RegisterMsg{
			Login:    login,
			Password: password,
		}
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

func LoginStore(pass string) tea.Cmd {
	return func() tea.Msg {
		return LoginStoreMsg(pass)
	}
}

type NewCardCmd struct {
	CurrentCardID string
	Name          string
	Ccn           string
	Exp           time.Time
	CVV           uint32
}

func NewCard(currentCardID, name, ccn string, exp time.Time, cvv uint32) tea.Cmd {
	return func() tea.Msg {
		return NewCardCmd{
			CurrentCardID: currentCardID,
			Name:          name,
			Ccn:           ccn,
			Exp:           exp,
			CVV:           cvv,
		}
	}
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
