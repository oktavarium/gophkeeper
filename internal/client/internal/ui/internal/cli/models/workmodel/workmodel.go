package workmodel

import (
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/card"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/simpledata"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

// model saves states and commands for them
type Model struct {
	table    table.Model
	card     *card.Model
	simple   *simpledata.Model
	focus    bool
	cards    map[string]models.SimpleCardData
	simples  map[string]models.SimpleData
	binaries map[string]models.SimpleBinaryData
}

// newModel create new model for cli
func NewModel() *Model {
	columns := []table.Column{
		{Title: "Type", Width: 8},
		{Title: "Name", Width: 30},
		{Title: "Modified", Width: 30},
		{},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(nil),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	t.Focus()

	return &Model{
		table:  t,
		card:   card.NewModel(),
		simple: simpledata.NewModel(),
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
		m.table.Focus()
		m.card.Blur()
		cmds = append(cmds, NewCard(msg.CurrentCardID, msg.Name, msg.Ccn, msg.Exp, msg.CVV))
	case simpledata.NewSimpleDataCmd:
		m.table.Focus()
		m.card.Blur()
		cmds = append(cmds, NewSimpleData(msg.CurrentID, msg.Name, msg.Data))
	case card.BlureCmd:
		m.table.Focus()
		m.card.Blur()
	case UpdateDataCmd:
		m.updateData(msg)
	case tea.KeyMsg:
		//nolint:exhaustive // too many unused cased
		switch msg.Type {
		case tea.KeyEnter:
			if m.table.Focused() {
				m.table.Blur()
				row := m.table.SelectedRow()
				if row == nil {
					return nil
				}

				switch models.DataTypeFromString(row[0]) {
				case models.Card:
					m.card.Focus()
					m.card.SetData(
						row[3],
						m.cards[row[3]].Data.Name,
						m.cards[row[3]].Data.Number,
						m.cards[row[3]].Data.ValidUntil.Format("01/06"),
						strconv.FormatUint(uint64(m.cards[row[3]].Data.CVV), 10),
					)
				case models.Simple:
					m.simple.Focus()
					m.simple.SetData(
						row[3],
						m.simples[row[3]].Data.Name,
						m.simples[row[3]].Data.Data,
					)
				case models.Binary:
				}

				return nil
			}
		case tea.KeyCtrlD:
			row := m.table.SelectedRow()
			if row != nil {
				return DeleteData(row[3])
			}

		case tea.KeyCtrlS:
			return Sync()
		case tea.KeyCtrlN:
			m.table.Blur()
			m.simple.Blur()
			m.card.Focus()
			m.card.Reset()
			return nil
		case tea.KeyCtrlJ:
			m.table.Blur()
			m.card.Blur()
			m.simple.Focus()
			m.simple.Reset()
			return nil
		case tea.KeyCtrlE:
			m.table.Blur()
			m.card.Focus()
			return nil
		}
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd, m.card.Update(msg), m.simple.Update(msg))

	return tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	var views []string
	baseStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240"))
	views = append(views, baseStyle.Render(m.table.View()))
	views = append(views, baseStyle.Render(m.card.View()))
	views = append(views, baseStyle.Render(m.simple.View()))

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n"
}

func (m *Model) updateData(data UpdateDataCmd) {
	m.cards = data.Cards
	m.simples = data.Simple
	m.binaries = data.Binary
	rows := make([]table.Row, 0, len(m.cards)+len(m.simples)+len(m.binaries))
	for k, v := range m.cards {
		rows = append(rows, []string{models.DataTypeToString(v.Common.Type), v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate), k})
	}
	for k, v := range m.simples {
		rows = append(rows, []string{models.DataTypeToString(v.Common.Type), v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate), k})
	}
	for k, v := range m.binaries {
		rows = append(rows, []string{models.DataTypeToString(v.Common.Type), v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate), k})
	}

	sort.Slice(rows, func(i, j int) bool {
		return rows[i][2] < rows[j][2]
	})
	m.table.SetRows(rows)
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
