package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// model saves states and commands for them
type registerStateModel struct {
	cursor int
	inputs []textinput.Model
	err    error
}

// newModel create new model for cli
func newRegisterStateModel() registerStateModel {
	inputs := make([]textinput.Model, 3)

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

	inputs[2] = textinput.New()
	inputs[2].Placeholder = "password"
	inputs[2].CharLimit = 8
	inputs[2].Width = 30
	inputs[2].EchoMode = textinput.EchoPassword
	inputs[2].Prompt = "Password again: "

	return registerStateModel{
		cursor: 0,
		inputs: inputs,
		err:    nil,
	}
}

// Init optionally returns an initial command we should run.
func (m registerStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m registerStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case resetMsg:
		m.reset()
		cmds = append(cmds, changeState(mainState))
	case errorMsg:
		m.err = msg
	case actionMsg:
		if err := validateInputs(m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value()); err != nil {
			m.err = err
			break
		}
		if err := validatePasswords(m.inputs[1].Value(), m.inputs[2].Value()); err != nil {
			m.err = err
			break
		}

		cmds = append(cmds, makeRegister(m.inputs[0].Value(), m.inputs[1].Value()))

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyDown:
			m.nextInput()
		case tea.KeyUp:
			m.prevInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.cursor].Focus()

	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m registerStateModel) View() string {
	view := fmt.Sprintf(
		`%s

%s
%s
%s
`, "[Registration form]\n\nPlease enter your login and password to register as new user.", m.inputs[0].View(), m.inputs[1].View(), m.inputs[2].View())

	if m.err != nil {
		view += fmt.Sprintf("\n\nError: %s", m.err)
	}

	return view
}

// nextInput focuses the next input field
func (m *registerStateModel) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *registerStateModel) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
}

func (m *registerStateModel) reset() {
	m.err = nil
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
}
