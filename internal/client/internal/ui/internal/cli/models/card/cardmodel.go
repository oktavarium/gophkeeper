package card

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	name = iota
	ccn
	exp
	cvv
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle = lipgloss.NewStyle().Foreground(hotPink)
)

type Model struct {
	inputs        []textinput.Model
	focused       int
	focus         bool
	currentCardID string
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
	m.currentCardID = ""
	for i := range m.inputs {
		m.inputs[i].Blur()
	}
}

func NewModel() Model {
	var inputs []textinput.Model = make([]textinput.Model, 4)
	inputs[name] = textinput.New()
	inputs[name].Placeholder = "Credit card name"
	inputs[name].Focus()
	inputs[name].CharLimit = 20
	inputs[name].Width = 30
	inputs[name].Prompt = ""

	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4505 **** **** 1234"
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY"
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 3
	inputs[cvv].Width = 5
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	return Model{
		inputs:  inputs,
		focused: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focus {
		return m, nil
	}
	cmds := make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if err := m.validateFields(); err != nil {
				return m, makeError(err)
			}
			return m, NewCard(m.currentCardID, m.inputs[0].Value(), getCcn(m.inputs[1].Value()), getExp(m.inputs[2].Value()), getCvv(m.inputs[3].Value()))
		case tea.KeyEsc:
			m.Reset()
			return m, makeBlur()
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
	return m, tea.Batch(cmds...)
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

 %s  %s
 %s  %s
`,
		inputStyle.Width(30).Render("Name:"),
		m.inputs[name].View(),
		inputStyle.Width(30).Render("Card Number"),
		m.inputs[ccn].View(),
		inputStyle.Width(6).Render("EXP"),
		inputStyle.Width(6).Render("CVV"),
		m.inputs[exp].View(),
		m.inputs[cvv].View(),
	)
}

// nextInput focuses the next input field
func (m *Model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *Model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m *Model) SetData(currentCardID, nameValue, ccnValue, expValue, cvvValue string) {
	m.currentCardID = currentCardID
	m.inputs[name].SetValue(nameValue)
	m.inputs[ccn].SetValue(ccnValue)
	m.inputs[exp].SetValue(expValue)
	m.inputs[cvv].SetValue(cvvValue)
}
