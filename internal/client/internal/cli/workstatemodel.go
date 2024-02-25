package cli

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// model saves states and commands for them
type workStateModel struct {
	cursor  int
	inputs  []textinput.Model
	err     error
	message string
}

// newModel create new model for cli
func newWorkStateModel() workStateModel {
	inputs := make([]textinput.Model, 2)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Name"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = "Name: "

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "value"
	inputs[1].CharLimit = 8
	inputs[1].Width = 30
	inputs[1].Prompt = "Value: "

	return workStateModel{
		cursor: 0,
		inputs: inputs,
		err:    nil,
	}
}

// Init optionally returns an initial command we should run.
func (m workStateModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update is called when messages are received.
func (m workStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	switch msg := msg.(type) {
	case resetMsg:
		m.reset()
		cmds = append(cmds, changeState(mainState))
	case errorMsg:
		m.err = msg
	case msgMsg:
		m.message = string(msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			err := validateInputs(m.inputs[0].Value(), m.inputs[1].Value())
			if err != nil {
				m.err = err
			} else {
				cmds = append(cmds, makeSaveData(m.inputs[0].Value(), m.inputs[1].Value()))
			}
		case tea.KeyTab:
			m.nextInput()
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
func (m workStateModel) View() string {
	view := fmt.Sprintf(
		`%s

%s
%s

%s
`, "Please enter data you want to save.", m.inputs[0].View(), m.inputs[1].View(), "Press enter to continue ->")

	if m.err != nil {
		view += fmt.Sprintf("\n\nError: %s", m.err)
	}

	if len(m.message) != 0 {
		view += fmt.Sprintf("\n\nStatus: %s", m.message)
	}

	return view
}

// nextInput focuses the next input field
func (m *workStateModel) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *workStateModel) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.inputs) - 1
	}
}

func (m *workStateModel) reset() {
	m.err = nil
	m.cursor = 0
	for i, input := range m.inputs {
		input.Reset()
		m.inputs[i] = input
	}
}
