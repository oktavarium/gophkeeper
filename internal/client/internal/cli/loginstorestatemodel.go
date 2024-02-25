package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// model saves states and commands for them
type loginStoreStateModel struct {
	cursor int
	inputs []textinput.Model
	err    error
}

// newModel create new model for cli
func newLoginStoreStateModel() loginStoreStateModel {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "password"
	inputs[0].CharLimit = 8
	inputs[0].Width = 30
	inputs[0].EchoMode = textinput.EchoPassword
	inputs[0].Prompt = "Master password: "
	inputs[0].Focus()

	return loginStoreStateModel{
		cursor: 0,
		inputs: inputs,
		err:    nil,
	}
}

// Init optionally returns an initial command we should run.
func (m loginStoreStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m loginStoreStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case resetMsg:
		m.reset()
		cmds = append(cmds, changeState(mainState))
	case errorMsg:
		m.err = msg
	case actionMsg:
		if err := validateInputs(m.inputs[0].Value()); err != nil {
			m.err = err
			break
		}
		cmds = append(cmds, loginLocalStore(m.inputs[0].Value()))

	case tea.KeyMsg:
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.cursor].Focus()
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
func (m loginStoreStateModel) View() string {
	view := fmt.Sprintf(
		`%s

%s
`, "[Local storage form]\n\nPlease enter your master password to log in.", m.inputs[0].View())

	if m.err != nil {
		view += fmt.Sprintf("\n\nError: %s", m.err)
	}

	return view
}

func (m *loginStoreStateModel) reset() {
	m.err = nil
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
}
