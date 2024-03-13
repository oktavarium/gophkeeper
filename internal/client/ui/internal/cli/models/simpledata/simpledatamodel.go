package simpledata

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/common"
)

const (
	name = iota
	data
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(hotPink)
)

type Model struct {
	inputs    []textinput.Model
	focused   int
	focus     bool
	currentID string
}

func (m Model) Focused() bool {
	return m.focus
}

func (m *Model) Focus() {
	m.focus = true
}

func (m *Model) Blur() {
	m.focus = false
}

func (m *Model) Reset() {
	for i := range m.inputs {
		m.inputs[i].Reset()
	}
	m.focused = 0
	m.currentID = ""
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func NewModel() *Model {
	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[name] = textinput.New()
	inputs[name].Placeholder = "Record name"
	inputs[name].Focus()
	inputs[name].CharLimit = 20
	inputs[name].Width = 30
	inputs[name].Prompt = ""

	inputs[data] = textinput.New()
	inputs[data].Placeholder = "content"
	inputs[data].CharLimit = 20
	inputs[data].Width = 30
	inputs[data].Prompt = ""

	return &Model{
		inputs:  inputs,
		focused: 0,
	}
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.focus {
		return nil
	}
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if err := common.ValidateInputs(m.inputs[0].Value(), m.inputs[1].Value()); err != nil {
				return makeError(err)
			}
			return NewSimpleData(m.currentID, m.inputs[0].Value(), m.inputs[1].Value())
		case tea.KeyEsc:
			m.Reset()
			return Back()
		case tea.KeyTab:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.focus {
		return ""
	}

	return fmt.Sprintf(
		` %s
 %s

 %s
 %s
`,
		inputStyle.Width(30).Render("Name:"),
		m.inputs[name].View(),
		inputStyle.Width(30).Render("Data:"),
		m.inputs[data].View(),
	)
}

// nextInput focuses the next input field
func (m *Model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *Model) SetData(currentID, nameValue, dataValue string) {
	m.currentID = currentID
	m.inputs[name].SetValue(nameValue)
	m.inputs[data].SetValue(dataValue)
}
