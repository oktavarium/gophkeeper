package workmodel

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/common"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/binary"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/card"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/simpledata"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/table"
	tablemodel "github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/table"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

// model saves states and commands for them
type Model struct {
	table    *tablemodel.Model
	card     *card.Model
	simple   *simpledata.Model
	binary   *binary.Model
	focus    bool
	cards    map[string]models.SimpleCardData
	simples  map[string]models.SimpleData
	binaries map[string]models.SimpleBinaryData
}

// newModel create new model for cli
func NewModel() *Model {
	t := table.NewModel()
	t.Focus()
	return &Model{
		table:  t,
		card:   card.NewModel(),
		simple: simpledata.NewModel(),
		binary: binary.NewModel(),
	}
}

// Update is called when messages are received.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.Focused() {
		return nil
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case card.NewCardCmd:
		m.BlurAll()
		m.table.Focus()
		cmds = append(cmds, NewCard(msg.CurrentCardID, msg.Name, msg.Ccn, msg.Exp, msg.CVV))
	case simpledata.NewSimpleDataMsg:
		m.BlurAll()
		m.table.Focus()
		cmds = append(cmds, NewSimpleData(msg.CurrentID, msg.Name, msg.Data))
	case binary.NewFileCmd:
		m.table.Focus()
		m.binary.Blur()
		cmds = append(cmds, NewFile(msg.CurrentID, msg.Name, msg.Path))
	case card.BackMsg, simpledata.BackMsg, binary.BackMsg:
		m.BlurAll()
		m.table.Focus()
	case table.BackMsg:
		return common.ChangeState(common.MainState)
	case table.SelectDataCmd:
		switch models.DataTypeFromString(msg.Type) {
		case models.Card:
			m.BlurAll()
			m.card.Focus()
			m.card.SetData(
				msg.ID,
				m.cards[msg.ID].Data.Name,
				m.cards[msg.ID].Data.Number,
				m.cards[msg.ID].Data.ValidUntil.Format("01/06"),
				strconv.FormatUint(uint64(m.cards[msg.ID].Data.CVV), 10),
			)
		case models.Simple:
			m.BlurAll()
			m.simple.Focus()
			m.simple.SetData(
				msg.ID,
				m.simples[msg.ID].Data.Name,
				m.simples[msg.ID].Data.Data,
			)
		case models.Binary:
			cmds = append(cmds, SaveFile(msg.ID))
		}
	case table.DeleteDataCmd:
		cmds = append(cmds, DeleteData(msg.ID))
	case UpdateDataCmd:
		cmds = append(cmds, m.updateData(msg))
	case tea.KeyMsg:
		//nolint:exhaustive // too many unused cased
		switch msg.Type {
		case tea.KeyCtrlS:
			return Sync()
		case tea.KeyF5:
			m.BlurAll()
			m.card.Focus()
			m.card.Reset()
		case tea.KeyF6:
			m.BlurAll()
			m.simple.Focus()
			m.simple.Reset()
		case tea.KeyF7:
			m.BlurAll()
			m.binary.Focus()
			cmds = append(cmds, m.binary.Init())
		}
	}

	cmds = append(cmds, m.table.Update(msg), m.card.Update(msg), m.simple.Update(msg), m.binary.Update(msg))
	return tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	var view string
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))

	switch {
	case m.table.Focused():
		view = m.table.View()
	case m.card.Focused():
		view = m.card.View()
	case m.simple.Focused():
		view = m.simple.View()
	case m.binary.Focused():
		view = m.binary.View()
	}

	return baseStyle.Render(view) + "\n\n"
}

func (m *Model) updateData(data UpdateDataCmd) tea.Cmd {
	m.cards = data.Cards
	m.simples = data.Simple
	m.binaries = data.Binary

	return table.UpdateData(data.Cards, data.Simple, data.Binary)
}

func (m *Model) Reset() {
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m Model) Focused() bool {
	return m.focus
}

func (m *Model) BlurAll() {
	m.card.Blur()
	m.table.Blur()
	m.simple.Blur()
	m.binary.Blur()
}
