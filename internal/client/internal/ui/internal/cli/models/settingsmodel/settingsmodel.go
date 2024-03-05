package settingsmodel

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

// NewModel create new model for cli
func NewModel() Model {
	inputs := make([]textinput.Model, 1)

	inputs[0] = textinput.New()
	inputs[0].Placeholder = "localhost:8888"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = "Server address: "

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
	case SetServerAddrMsg:
		m.inputs[0].SetValue(string(msg))
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
			err := common.ValidateInputs(m.inputs[0].Value())
			if err != nil {
				cmds = append(cmds, common.MakeError(err))
			} else {
				cmds = append(cmds, SaveServerAddr(m.inputs[0].Value()))
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
func (m Model) View() string {
	view := fmt.Sprintf(
		`%s

%s
`, "[Settings form]\n\nPlease enter your server settings.", m.inputs[0].View())

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

func (m *Model) Reset() {}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m Model) Focused() bool {
	return m.focus
}
