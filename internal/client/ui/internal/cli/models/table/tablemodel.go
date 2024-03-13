package table

import (
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

// model saves states and commands for them
type Model struct {
	table table.Model
	focus bool
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
		table: t,
	}
}

// Update is called when messages are received.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.Focused() {
		return nil
	}
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case UpdateDataCmd:
		m.updateData(msg)
	case tea.KeyMsg:
		//nolint:exhaustive // too many unused cased
		switch msg.Type {
		case tea.KeyEnter:
			row := m.table.SelectedRow()
			if row != nil {
				return SelectData(row[0], row[3])
			}
		case tea.KeyBackspace:
			row := m.table.SelectedRow()
			if row != nil {
				return DeleteData(row[3])
			}
		case tea.KeyEsc:
			return Back()
		}
	}

	var tableCmd tea.Cmd
	m.table, tableCmd = m.table.Update(msg)
	cmds = append(cmds, tableCmd)

	return tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	if !m.Focused() {
		return ""
	}
	return m.table.View()
}

func (m *Model) updateData(data UpdateDataCmd) {
	rows := make([]table.Row, 0, len(data.Cards)+len(data.Simple)+len(data.Binary))
	for k, v := range data.Cards {
		rows = append(rows, []string{models.DataTypeToString(v.Common.Type), v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate), k})
	}
	for k, v := range data.Simple {
		rows = append(rows, []string{models.DataTypeToString(v.Common.Type), v.Data.Name, v.Common.Modified.UTC().Format(time.UnixDate), k})
	}
	for k, v := range data.Binary {
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
