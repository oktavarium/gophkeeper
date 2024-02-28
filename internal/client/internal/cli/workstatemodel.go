package cli

import (
	"fmt"

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
	cursor  int
	table   table.Model
	card    card.Model
	buttons []string
	err     error
	cards   map[string]dto.SimpleCardData
}

// newModel create new model for cli
func newWorkStateModel() workStateModel {
	buttons := []string{
		"Create new card",
		"Sync with server",
	}
	columns := []table.Column{
		{Title: "Id", Width: 4},
		{Title: "Name", Width: 30},
		{Title: "Modified", Width: 30},
	}

	rows := []table.Row{
		{"1", "Tokyo", "Japan"},
		{"2", "Delhi", "India"},
		{"3", "Shanghai", "China"},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
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
		cursor:  len(buttons),
		buttons: buttons,
		table:   t,
		card:    card.InitialModel(),
		err:     nil,
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
		cmds = append(cmds, newCard(msg.Name, msg.Ccn, msg.Exp, msg.CVV))
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.cursor == 0 {
				m.table.Blur()
				m.card.Focus()
			} else {
				if m.table.Focused() {
					m.table.Blur()
					m.card.Focus()
					m.card.SetData(m.table.SelectedRow()[1])
				}
			}

		case tea.KeyTab:
			m.cursor++
			m.card.Blur()
			m.table.Blur()
			if m.cursor == len(m.buttons) {
				m.table.Focus()
			} else if m.cursor > len(m.buttons) {
				m.cursor = 0
			}
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

	var buttonViews []string
	for i, name := range m.buttons {
		var cursor string = " "
		if m.cursor == i {
			cursor = ">"
		}
		buttonViews = append(buttonViews, fmt.Sprintf("%s %s", cursor, name))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, views...) + "\n\n" +
		lipgloss.JoinVertical(lipgloss.Left, buttonViews...) + "\n\n"
}

func (m *workStateModel) UpdateCards(cards map[string]dto.SimpleCardData) {
	m.cards = cards
	rows := make([]table.Row, 0, len(m.cards))
	for k, v := range m.cards {
		rows = append(rows, []string{k, v.Data.Number, v.Common.Modified.String()})
	}
	m.table.SetRows(rows)
}
