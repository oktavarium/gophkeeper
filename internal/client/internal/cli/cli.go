package cli

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/internal/storage"
	"github.com/oktavarium/gophkeeper/internal/shared/buildinfo"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type state int

const (
	quitState state = iota
	mainState
	loginState
	registerState
	workState
	settingsState
	localStoreState
)

// model saves states and commands for them
type model struct {
	ctx          context.Context
	states       map[state]tea.Model
	currentState state
	help         string
	helpShown    bool
	storage      storage.Storage
	currentUser  dto.UserInfo
	serverAddr   string
}

func newModel(ctx context.Context, s storage.Storage) model {
	states := map[state]tea.Model{
		mainState:       newMainStateModel(),
		loginState:      newLoginStateModel(),
		registerState:   newRegisterStateModel(),
		workState:       newWorkStateModel(),
		settingsState:   newSettingsStateModel(),
		localStoreState: newLocalStoreStateModel(),
	}

	return model{
		ctx:          ctx,
		states:       states,
		currentState: mainState,
		help:         "\n\nNavigation: Tab, Arrows;\nBack: Left;\nSelect command: Enter, Space;\nExit: Ctrl+C",
		storage:      s,
	}
}

func Run(ctx context.Context, s storage.Storage) error {
	if _, err := tea.NewProgram(newModel(ctx, s), tea.WithContext(ctx)).Run(); err != nil {
		return fmt.Errorf("could not start program: %s", err)
	}

	return nil
}

// Init optionally returns an initial command we should run.
func (m model) Init() tea.Cmd {
	return tea.ClearScreen
}

// Update is called when messages are received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		case tea.KeyLeft:
			cmds = append(cmds, makeReset)
		case tea.KeyCtrlH:
			m.helpShown = !m.helpShown
		case tea.KeyEnter, tea.KeySpace:
			cmds = append(cmds, makeAction)
		}
	case stateCmd:
		m.currentState = state(msg)
		return m, tea.ClearScreen
	case loginCmd:
		if err := m.login(msg.login, msg.password); err != nil {
			return m, makeError(err)
		}
		m.currentUser = dto.UserInfo{
			Login:    msg.login,
			Password: msg.password,
		}
		m.currentState = workState
		return m, tea.ClearScreen
	case registerCmd:
		if err := m.register(msg.login, msg.password); err != nil {
			cmds = append(cmds, makeError(err))
		} else {
			cmds = append(cmds, makeReset)
		}
	case serverAddrCmd:
		m.serverAddr = string(msg)
		if err := m.initClient(m.serverAddr); err != nil {
			cmds = append(cmds, makeError(err))
		}
	case saveCmd:
		if err := m.saveData(msg.name, msg.data); err != nil {
			cmds = append(cmds, makeError(err))
		} else {
			cmds = append(cmds, makeMsg("Data saved"))
		}
	case tea.QuitMsg:
		return m, tea.Quit
	}

	updatedModel, cmd := m.states[m.currentState].Update(msg)
	m.states[m.currentState] = updatedModel
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m model) View() string {
	view := m.states[m.currentState].View()
	view += fmt.Sprintf("\n\nPress Ctrl+H to show/hide help.")
	if m.helpShown {
		view += m.help
		view += fmt.Sprintf("\n\nVersion: %s\nBuild date: %s", buildinfo.Version, buildinfo.BuildDate)
	}

	return view
}

func (m model) login(login, password string) error {
	return m.storage.Login(
		m.ctx,
		dto.UserInfo{
			Login:    login,
			Password: password,
		},
	)
}

func (m model) register(login, password string) error {
	return m.storage.Register(
		m.ctx,
		dto.UserInfo{
			Login:    login,
			Password: password,
		},
	)
}

func (m model) initClient(addr string) error {
	if err := m.storage.Init(m.ctx, addr); err != nil {
		return fmt.Errorf("error on client init: %w", err)
	}

	return nil
}

func (m model) saveData(name, data string) error {
	return m.storage.Save(
		m.ctx,
		dto.SaveData{
			UserInfo: m.currentUser,
			Name:     name,
			Data:     data,
		},
	)
}
