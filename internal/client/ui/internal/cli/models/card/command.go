package card

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type ErrorMsg error

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
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

type BackMsg struct{}

func Back() tea.Cmd {
	return func() tea.Msg {
		return BackMsg{}
	}
}
