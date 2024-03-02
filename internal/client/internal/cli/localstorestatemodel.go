package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// model saves states and commands for them
type localStoreStateModel struct {
	cursor int
	inputs []textinput.Model
}

// newModel create new model for cli
func newLocalStoreStateModel() localStoreStateModel {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "password"
	inputs[0].CharLimit = 8
	inputs[0].Width = 30
	inputs[0].EchoMode = textinput.EchoPassword
	inputs[0].Prompt = "Master password: "

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "password"
	inputs[1].CharLimit = 8
	inputs[1].Width = 30
	inputs[1].EchoMode = textinput.EchoPassword
	inputs[1].Prompt = "Master password again: "

	return localStoreStateModel{
		cursor: 0,
		inputs: inputs,
	}
}

// Init optionally returns an initial command we should run.
func (m localStoreStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m localStoreStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case resetMsg:
		m.reset()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyDown:
			m.nextInput()
		case tea.KeyUp:
			m.prevInput()
		case tea.KeyEnter:
			if err := validateInputs(m.inputs[0].Value(), m.inputs[1].Value()); err != nil {
				cmds = append(cmds, makeError(err))
				break
			}
			if err := validatePasswords(m.inputs[0].Value(), m.inputs[1].Value()); err != nil {
				cmds = append(cmds, makeError(err))
				break
			}
			cmds = append(cmds, loginLocalStore(m.inputs[0].Value()))
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
func (m localStoreStateModel) View() string {
	view := fmt.Sprintf(
		`%s

%s
%s
`, "[Creating local storage form]\n\nPlease enter your master password for data encryption.", m.inputs[0].View(), m.inputs[1].View())

	return view
}

// nextInput focuses the next input field
func (m *localStoreStateModel) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *localStoreStateModel) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
}

func (m *localStoreStateModel) reset() {
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
}
