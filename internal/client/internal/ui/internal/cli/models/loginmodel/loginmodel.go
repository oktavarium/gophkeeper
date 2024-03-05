package loginmodel

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/common"
)

// model saves states and commands for them
type Model struct {
	cursor int
	inputs []textinput.Model
	focus  bool
}

// newModel create new model for cli
func NewModel() *Model {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "username"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = "Login: "

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "password"
	inputs[1].CharLimit = 8
	inputs[1].Width = 30
	inputs[1].EchoMode = textinput.EchoPassword
	inputs[1].Prompt = "Password: "

	return &Model{
		cursor: 0,
		inputs: inputs,
	}
}

// Update is called when messages are received.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.Focused() {
		return nil
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case common.ResetMsg:
		m.Reset()
	case tea.KeyMsg:
		//nolint:exhaustive // too many unused cased
		switch msg.Type {
		case tea.KeyTab, tea.KeyDown:
			m.nextInput()
		case tea.KeyUp:
			m.prevInput()
		case tea.KeyEnter:
			err := common.ValidateInputs(m.inputs[0].Value(), m.inputs[1].Value())
			if err != nil {
				cmds = append(cmds, common.MakeError(err))
			} else {
				cmds = append(cmds, MakeLogin(m.inputs[0].Value(), m.inputs[1].Value()))
			}
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.cursor].Focus()
	}

	var cmd tea.Cmd
	for i := range m.inputs {
		m.inputs[i], cmd = m.inputs[i].Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	view := fmt.Sprintf(
		`%s

%s
%s
`, "[Login form]\n\nPlease enter your login and password to login.", m.inputs[0].View(), m.inputs[1].View())

	return view
}

// nextInput focuses the next input field
func (m *Model) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *Model) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
}

func (m *Model) Reset() {
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
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
