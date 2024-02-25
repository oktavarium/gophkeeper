package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type command struct {
	name  string
	state state
}

// model saves states and commands for them
type mainStateModel struct {
	commands []command
	cursor   int
	err      error
}

// newMainStateModel create new model for cli
func newMainStateModel() mainStateModel {
	return mainStateModel{
		commands: []command{
			{
				name:  "Login",
				state: loginState,
			},
			{
				name:  "Register",
				state: registerState,
			},
			{
				name:  "Settings",
				state: settingsState,
			},
		},
		cursor: 0,
	}
}

// Init optionally returns an initial command we should run.
func (m mainStateModel) Init() tea.Cmd {
	return nil
}

// Update is called when messages are received.
func (m mainStateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case resetMsg:
		m.reset()
		return m, nil
	case actionMsg:
		return m, changeState(m.commands[m.cursor].state)
	case errorMsg:
		m.err = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab, tea.KeyDown:
			m.nextInput()
		case tea.KeyUp:
			m.prevInput()
		}
	}

	return m, nil
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m mainStateModel) View() string {
	view := "Hello! This is gophkeeper.\nTo start using keeper you have to login or register.\n\n"
	for i, cmd := range m.commands {
		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		view += fmt.Sprintf("%s %s\n", cursor, cmd.name)
	}

	if m.err != nil {
		view += fmt.Sprintf("\n\nError: %s", m.err)
	}

	return view
}

func (m *mainStateModel) nextInput() {
	m.cursor = (m.cursor + 1) % len(m.commands)
}

// prevInput focuses the previous input field
func (m *mainStateModel) prevInput() {
	m.cursor--
	if m.cursor < 0 {
		m.cursor = len(m.commands) - 1
	}
}

func (m *mainStateModel) reset() {
	m.cursor = 0
}
