package workmodel

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

type UpdateDataCmd struct {
	Cards  map[string]models.SimpleCardData
	Simple map[string]models.SimpleData
	Binary map[string]models.SimpleBinaryData
}

func UpdateData(
	cards map[string]models.SimpleCardData,
	simple map[string]models.SimpleData,
	binary map[string]models.SimpleBinaryData,
) tea.Cmd {
	return func() tea.Msg {
		return UpdateDataCmd{
			Cards:  cards,
			Simple: simple,
			Binary: binary,
		}
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

type NewFileCmd struct {
	CurrentID string
	Name      string
	Path      string
}

func NewFile(currentID, name, path string) tea.Cmd {
	return func() tea.Msg {
		return NewFileCmd{
			CurrentID: currentID,
			Name:      name,
			Path:      path,
		}
	}
}

type SaveFileCmd string

func SaveFile(id string) tea.Cmd {
	return func() tea.Msg {
		return SaveFileCmd(id)
	}
}

type DeleteDataMsg string

func DeleteData(id string) tea.Cmd {
	return func() tea.Msg {
		return DeleteDataMsg(id)
	}
}

type SyncMsg struct{}

func Sync() tea.Cmd {
	return func() tea.Msg {
		return SyncMsg{}
	}
}
