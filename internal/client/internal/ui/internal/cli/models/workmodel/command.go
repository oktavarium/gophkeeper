package workmodel

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

type UpdateCardsCmd map[string]models.SimpleCardData

func UpdateCards(cards map[string]models.SimpleCardData) tea.Cmd {
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

type NewSimpleDataCmd struct {
	CurrentID string
	Name      string
	Data      string
}

func NewSimpleData(currentID, name, data string) tea.Cmd {
	return func() tea.Msg {
		return NewSimpleDataCmd{
			CurrentID: currentID,
			Name:      name,
			Data:      data,
		}
	}
}

type DeleteCardMsg string

func DeleteCard(id string) tea.Cmd {
	return func() tea.Msg {
		return DeleteCardMsg(id)
	}
}

type DeleteSimpleMsg string

func DeleteSimple(id string) tea.Cmd {
	return func() tea.Msg {
		return DeleteSimpleMsg(id)
	}
}

type DeleteBinaryMsg string

func DeleteBinary(id string) tea.Cmd {
	return func() tea.Msg {
		return DeleteBinaryMsg(id)
	}
}

type SyncMsg struct{}

func Sync() tea.Cmd {
	return func() tea.Msg {
		return SyncMsg{}
	}
}
