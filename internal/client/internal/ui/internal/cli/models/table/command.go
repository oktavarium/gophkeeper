package table

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

type ErrorMsg error

func makeError(err error) tea.Cmd {
	return func() tea.Msg {
		return ErrorMsg(err)
	}
}

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

type SelectDataCmd struct {
	Type string
	ID   string
}

func SelectData(dataType, ID string) tea.Cmd {
	return func() tea.Msg {
		return SelectDataCmd{
			Type: dataType,
			ID:   ID,
		}
	}
}

type DeleteDataCmd struct {
	ID string
}

func DeleteData(ID string) tea.Cmd {
	return func() tea.Msg {
		return DeleteDataCmd{
			ID: ID,
		}
	}
}
