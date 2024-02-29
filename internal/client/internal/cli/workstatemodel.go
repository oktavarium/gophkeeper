package cli

import (
	"sort"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oktavarium/gophkeeper/internal/client/internal/cli/internal/card"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

// model saves states and commands for them
type workStateModel struct {
	table table.Model
	card  card.Model
	err   error
	cards map[string]dto.SimpleCardData
}

// newModel create new model for cli
func newWorkStateModel() workStateModel {
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Name", Width: 30},
		{Title: "Modified", Width: 30},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(nil),
		table.WithFocused(true),
		table.WithHeight(11),
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

	return workStateModel{
		table: t,
		card:  card.InitialModel(),
		err:   nil,
	}
}

// Init optionally returns an initial command we should run.
func (m workStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m workStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case card.NewCardCmd:
		m.table.Focus()
		m.card.Blur()
		cmds = append(cmds, newCard(msg.CurrentCardID, msg.Name, msg.Ccn, msg.Exp, msg.CVV))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.table.Focused() {
				m.table.Blur()
				m.card.Focus()
				row := m.table.SelectedRow()
				if row != nil {
					m.card.SetData(
						row[0],
						m.cards[row[0]].Data.Name,
						m.cards[row[0]].Data.Number,
						m.cards[row[0]].Data.ValidUntil.Format("01/06"),
						strconv.FormatUint(uint64(m.cards[row[0]].Data.CVV), 10),
					)
				}
				return m, nil
			}
		case tea.KeyCtrlD:
			row := m.table.SelectedRow()
			if row != nil {
				return m, deleteCard(row[0])
			}

		case tea.KeyCtrlS:
			return m, sync()
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
func (m workStateModel) View() string {
	// return baseStyle.Render(m.table.View()) + "\n"

	var views []string
	views = append(views, baseStyle.Render(m.table.View()))
	views = append(views, baseStyle.Render(m.card.View()))

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n"
}

func (m *workStateModel) UpdateCards(cards map[string]dto.SimpleCardData) {
	m.cards = cards
	rows := make([]table.Row, 0, len(m.cards))
	for k, v := range m.cards {
		rows = append(rows, []string{k, v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate)})
		sort.Slice(rows, func(i, j int) bool {
			return rows[i][2] < rows[j][2]
		})
	}
	m.table.SetRows(rows)
}
