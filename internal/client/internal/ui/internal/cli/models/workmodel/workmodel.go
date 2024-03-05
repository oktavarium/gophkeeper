package workmodel

import (
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/card"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/simpledata"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

// model saves states and commands for them
type Model struct {
	table   table.Model
	card    card.Model
	simple  simpledata.Model
	focus   bool
	cards   map[string]models.SimpleCardData
	simples map[string]models.SimpleData
}

// newModel create new model for cli
func NewModel() Model {
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

	return Model{
		table:  t,
		card:   card.NewModel(),
		simple: simpledata.NewModel(),
	}
}

// Init optionally returns an initial command we should run.
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.Focused() {
		return m, nil
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
	case UpdateCardsCmd:
		m.updateCards(msg)
	case tea.KeyMsg:
		//nolint:exhaustive // too many unused cased
		switch msg.Type {
		case tea.KeyEnter:
			if m.table.Focused() {
				m.table.Blur()
				m.card.Focus()
				row := m.table.SelectedRow()
				if row != nil {
					m.card.SetData(
						row[3],
						m.cards[row[3]].Data.Name,
						m.cards[row[3]].Data.Number,
						m.cards[row[3]].Data.ValidUntil.Format("01/06"),
						strconv.FormatUint(uint64(m.cards[row[3]].Data.CVV), 10),
					)
				}
				return m, nil
			}
		case tea.KeyCtrlD:
			row := m.table.SelectedRow()
			if row != nil {
				return m, DeleteCard(row[0])
			}

		case tea.KeyCtrlS:
			return m, Sync()
		case tea.KeyCtrlN:
			m.table.Blur()
			m.card.Focus()
			m.card.Reset()
			return m, nil
		case tea.KeyCtrlE:
			m.table.Blur()
			m.card.Focus()
			return m, nil
		}
	}

	var tableCmd, cardCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	m.card, cardCmd = m.card.Update(msg)
	cmds = append(cmds, tableCmd, cardCmd)

	return m, tea.Batch(cmds...)
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

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n"
}

func (m *Model) updateCards(cards map[string]models.SimpleCardData) {
	m.cards = cards
	rows := make([]table.Row, 0, len(m.cards))
	for k, v := range m.cards {
		rows = append(rows, []string{models.DataTypeToString(v.Common.Type), v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate), k})
		sort.Slice(rows, func(i, j int) bool {
			return rows[i][2] < rows[j][2]
		})
	}
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
