package cli

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/internal/remote"
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
	loginLocalStoreState
)

// model saves states and commands for them
type model struct {
	ctx          context.Context
	states       map[state]tea.Model
	currentState state
	help         string
	helpShown    bool
	storage      storage.Storage
	remoteClient remote.Client
	currentUser  dto.UserInfo
	serverAddr   string
}

func newModel(ctx context.Context, s storage.Storage, c remote.Client) model {
	states := map[state]tea.Model{
		mainState:            newMainStateModel(),
		loginState:           newLoginStateModel(),
		registerState:        newRegisterStateModel(),
		workState:            newWorkStateModel(),
		settingsState:        newSettingsStateModel(),
		localStoreState:      newLocalStoreStateModel(),
		loginLocalStoreState: newLoginStoreStateModel(),
	}

	return model{
		ctx:          ctx,
		states:       states,
		currentState: mainState,
		help:         "\n\nNavigation: Tab, Arrows;\nBack: Left;\nSelect command: Enter, Space;\nExit: Ctrl+C",
		storage:      s,
		remoteClient: c,
	}
}

func Run(ctx context.Context, s storage.Storage, c remote.Client) error {
	model := newModel(ctx, s, c)
	if _, err := tea.NewProgram(model, tea.WithContext(ctx)).Run(); err != nil {
		return fmt.Errorf("could not start program: %s", err)
	}

	return nil
}

// Init optionally returns an initial command we should run.
func (m model) Init() tea.Cmd {
	return checkStore
}

// Update is called when messages are received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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
	case checkStoreMsg:
		if err := m.storage.Check(); err != nil {
			cmds = append(cmds, changeState(localStoreState))
		} else {
			cmds = append(cmds, changeState(loginLocalStoreState))
		}
		cmds = append(cmds, textinput.Blink)
	case stateMsg:
		m.currentState = state(msg)
		cmds = append(cmds, setServerAddr(m.serverAddr))
		cmds = append(cmds, textinput.Blink)
		cmds = append(cmds, tea.ClearScreen)
	case loginMsg:
		if err := m.login(msg.login, msg.password); err != nil {
			cmds = append(cmds, makeError(err))
		} else {
			m.currentUser.Login = msg.login
			m.currentUser.Password = msg.password
			cmds = append(cmds, changeState(workState), tea.ClearScreen)
		}
	case registerMsg:
		if err := m.register(msg.login, msg.password); err != nil {
			cmds = append(cmds, makeError(err))
		} else {
			m.currentUser.Login = msg.login
			m.currentUser.Password = msg.password
			cmds = append(cmds, makeReset)
		}
	// case createLocalStoreMsg:
	// 	if err := m.storage.Open(string(msg)); err != nil {
	// 		cmds = append(cmds, makeError(err))
	// 	} else {
	// 		cmds = append(cmds, changeState(mainState))
	// 	}
	case loginLocalStoreMsg:
		if serverAddr, userInfo, err := m.loginLocalStore(string(msg)); err != nil {
			cmds = append(cmds, makeError(err))
		} else {
			m.serverAddr = serverAddr
			m.currentUser = userInfo
			if len(m.serverAddr) != 0 {
				if err := m.initClient(m.serverAddr); err != nil {
					cmds = append(cmds, makeError(err))
				}
			}
			cmds = append(cmds, changeState(mainState))
		}
	case serverAddrMsg:
		m.serverAddr = string(msg)
		if err := m.initClient(m.serverAddr); err != nil {
			cmds = append(cmds, makeError(err))
		} else if err := m.storage.SetServerAddr(m.serverAddr); err != nil {
			cmds = append(cmds, makeError(err))
		}
	case saveMsg:
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

	if len(m.currentUser.Login) != 0 {
		view += fmt.Sprintf("\n\nYour current user is \"%s\"", m.currentUser.Login)
	}

	return view
}

func (m model) loginLocalStore(password string) (string, dto.UserInfo, error) {
	var serverAddr string
	var userInfo dto.UserInfo
	if err := m.storage.Open(password); err != nil {
		return serverAddr, userInfo, fmt.Errorf("error on opening storage: %w", err)
	}

	serverAddr, err := m.storage.GetServerAddr()
	if err != nil {
		return serverAddr, userInfo, fmt.Errorf("error on reading server addr: %w", err)
	}

	login, pass, err := m.storage.GetLoginAndPass()
	if err != nil {
		return serverAddr, userInfo, fmt.Errorf("error on reading login and pass: %w", err)
	}

	userInfo.Login = login
	userInfo.Password = pass

	return serverAddr, userInfo, nil
}

func (m model) register(login, password string) error {
	if err := m.remoteClient.Register(m.ctx, dto.UserInfo{
		Login:    login,
		Password: password,
	}); err != nil {
		return fmt.Errorf("error on registering user on server: %w", err)
	}

	if err := m.storage.SetLoginAndPass(login, password); err != nil {
		return fmt.Errorf("error on saving login and password in local storage: %w", err)
	}

	return nil
}

func (m model) login(login, password string) error {
	if err := m.remoteClient.Login(m.ctx, dto.UserInfo{
		Login:    login,
		Password: password,
	}); err != nil {
		return fmt.Errorf("error on loging user on server: %w", err)
	}

	if err := m.storage.SetLoginAndPass(login, password); err != nil {
		return fmt.Errorf("error on saving login and password in local storage: %w", err)
	}

	return nil
}

func (m model) initClient(addr string) error {
	if err := m.remoteClient.Init(m.ctx, addr); err != nil {
		return fmt.Errorf("error on client init: %w", err)
	}

	return nil
}

func (m model) saveData(name, data string) error {
	return nil
}
