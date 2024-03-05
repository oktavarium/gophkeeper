package workmodel

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type UpdateCardsCmd map[string]dto.SimpleCardData

func UpdateCards(cards map[string]dto.SimpleCardData) tea.Cmd {
	return func() tea.Msg {
		return UpdateCardsCmd(cards)
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
