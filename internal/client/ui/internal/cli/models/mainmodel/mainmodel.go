package mainmodel

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/common"
)

type command struct {
	name  string
	state common.State
}

// model saves states and commands for them
type Model struct {
	commands []command
	cursor   int
	focus    bool
}

// newModel create new model for cli
func NewModel() *Model {
	return &Model{
		commands: []command{
			{
				name:  "Login",
				state: common.LoginState,
			},
			{
				name:  "Register",
				state: common.RegisterState,
			},
			{
				name:  "Settings",
				state: common.SettingsState,
			},
		},
		cursor: 0,
	}
}

// Update is called when messages are received.
func (m *Model) Update(msg tea.Msg) tea.Cmd {
	if !m.Focused() {
		return nil
	}

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
			return common.ChangeState(m.commands[m.cursor].state)
		}
	}

	return nil
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m Model) View() string {
	view := "Hello! This is gophkeeper.\nTo start using keeper you have to login or register.\n\n"
	for i, cmd := range m.commands {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		view += fmt.Sprintf("%s %s\n", cursor, cmd.name)
	}

	return view
}

func (m *Model) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.commands)
}

// prevInput focuses the previous input field
func (m *Model) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.commands) - 1
	}
}

func (m *Model) Reset() {
	m.cursor = 0
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
