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
}

// newModel create new model for cli
func newSettingsStateModel() settingsStateModel {
	inputs := make([]textinput.Model, 1)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "localhost:8888"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = "Server address: "

	return settingsStateModel{
		cursor: 0,
		inputs: inputs,
	}
}

// Init optionally returns an initial command we should run.
func (m settingsStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m settingsStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case setServerAddrMsg:
		m.inputs[0].SetValue(string(msg))
	case resetMsg:
		m.reset()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyDown:
			m.nextInput()
		case tea.KeyUp:
			m.prevInput()
		case tea.KeyEnter:
			err := validateInputs(m.inputs[0].Value())
			if err != nil {
				cmds = append(cmds, makeError(err))
			} else {
				cmds = append(cmds, saveServerAddr(m.inputs[0].Value()), changeState(mainState))
			}
		}
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
func (m settingsStateModel) View() string {
	view := fmt.Sprintf(
		`%s

%s
`, "[Settings form]\n\nPlease enter your server settings.", m.inputs[0].View())

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
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
}
