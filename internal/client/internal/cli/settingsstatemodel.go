package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// model saves states and commands for them
type settingsStateModel struct {
	cursor int
	inputs []textinput.Model
	err    error
}

// newModel create new model for cli
func newSettingsStateModel() settingsStateModel {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "localhost:8888"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = "Server address: "

	return settingsStateModel{
		cursor: 0,
		inputs: inputs,
		err:    nil,
	}
}

// Init optionally returns an initial command we should run.
func (m settingsStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m settingsStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case resetMsg:
		m.reset()
		cmds = append(cmds, changeState(mainState))
	case errorMsg:
		m.err = msg
	case actionMsg:
		err := validateInputs(m.inputs[0].Value())
		if err != nil {
			m.err = err
		} else {
			cmds = append(cmds, saveServerAddr(m.inputs[0].Value()), changeState(mainState))
		}
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
func (m settingsStateModel) View() string {
	view := fmt.Sprintf(
		`%s

%s
%s
`, "[Settings form]\n\nPlease enter your server settings.", m.inputs[0].View(), m.inputs[1].View())

	if m.err != nil {
		view += fmt.Sprintf("\n\nError: %s", m.err)
	}

	return view
}

// nextInput focuses the next input field
func (m *settingsStateModel) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *settingsStateModel) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
}

func (m *settingsStateModel) reset() {
	m.err = nil
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
}
