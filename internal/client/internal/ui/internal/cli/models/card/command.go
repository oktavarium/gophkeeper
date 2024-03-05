package card

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type (
	ErrorMsg   error
	NewCardCmd struct {
		CurrentCardID string
		Name          string
		Ccn           string
		Exp           time.Time
		CVV           uint32
	}
	BlureCmd struct{}
)

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

func makeBlur() tea.Cmd {
	return func() tea.Msg {
		return BlureCmd{}
	}
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
