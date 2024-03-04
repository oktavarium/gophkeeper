package loginstoremodel

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli"
)

// model saves states and commands for them
type Model struct {
	cursor int
	inputs []textinput.Model
	focus  bool
}

// newModel create new model for cli
func NewModel() Model {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "password"
	inputs[0].CharLimit = 8
	inputs[0].Width = 30
	inputs[0].EchoMode = textinput.EchoPassword
	inputs[0].Prompt = "Master password: "
	inputs[0].Focus()

	return Model{
		cursor: 0,
		inputs: inputs,
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
	case cli.ResetMsg:
		m.Reset()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if err := cli.ValidateInputs(m.inputs[0].Value()); err != nil {
				cmds = append(cmds, cli.MakeError(err))
				break
			}
			cmds = append(cmds, cli.LoginStore(m.inputs[0].Value()))
		default:
			for i := range m.inputs {
				m.inputs[i].Blur()
			}
			m.inputs[m.cursor].Focus()
		}

	}

	for i := range m.inputs {
		var cmd tea.Cmd
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	view := fmt.Sprintf(
		`%s

%s
`, "[Local storage form]\n\nPlease enter your master password to log in.", m.inputs[0].View())

	return view
}

func (m *Model) Reset() {
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
	m.inputs[0].Focus()
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
