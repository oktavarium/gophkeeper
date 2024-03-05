package cli

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/common"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/loginmodel"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/loginstoremodel"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/mainmodel"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/registermodel"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/settingsmodel"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/storemodel"
	"github.com/oktavarium/gophkeeper/internal/client/internal/ui/internal/cli/models/workmodel"
	"github.com/oktavarium/gophkeeper/internal/shared/buildinfo"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

// model saves states and commands for them
type model struct {
	ctx             context.Context
	states          map[common.State]common.Model
	currentState    common.State
	mainState       mainmodel.Model
	loginState      loginmodel.Model
	registerState   registermodel.Model
	workState       workmodel.Model
	settingsState   settingsmodel.Model
	storeState      storemodel.Model
	loginStoreState loginstoremodel.Model
	help            string
	helpShown       bool
	storage         Storage
	remoteClient    common.RemoteClient
	currentUser     dto.UserInfo
	serverAddr      string
	err             error
}

func newModel(ctx context.Context, s Storage, c common.RemoteClient) model {
	mainState := mainmodel.NewModel()
	loginState := loginmodel.NewModel()
	registerState := registermodel.NewModel()
	workState := workmodel.NewModel()
	settingsState := settingsmodel.NewModel()
	storeState := storemodel.NewModel()
	loginStoreState := loginstoremodel.NewModel()

	return model{
		ctx:             ctx,
		mainState:       mainState,
		loginState:      loginState,
		registerState:   registerState,
		workState:       workState,
		settingsState:   settingsState,
		storeState:      storeState,
		loginStoreState: loginStoreState,
		help:            "\n\nNavigation: Tab, Arrows;\nBack: Left;\nSelect command: Enter, Space;\nExit: Ctrl+C",
		storage:         s,
		remoteClient:    c,
	}
}

func Run(ctx context.Context, s Storage, c common.RemoteClient) error {
	model := newModel(ctx, s, c)
	if _, err := tea.NewProgram(model, tea.WithContext(ctx)).Run(); err != nil {
		return fmt.Errorf("could not start program: %s", err)
	}

	return nil
}

// Init optionally returns an initial command we should run.
func (m model) Init() tea.Cmd {
	if err := m.storage.Check(); err != nil {
		return common.ChangeState(common.StoreState)
	} else {
		return common.ChangeState(common.LoginStoreState)
	}
}

// Update is called when messages are received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case common.ErrorMsg:
		m.err = msg
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			cmds = append(cmds, tea.Quit)
		case tea.KeyLeft:
			cmds = append(cmds, common.MakeReset)
			cmds = append(cmds, common.ChangeState(common.MainState))
		case tea.KeyCtrlH:
			m.helpShown = !m.helpShown
		}
	case common.StateMsg:
		m.err = nil
		m.currentState = common.State(msg)
		cmds = append(cmds, settingsmodel.SetServerAddr(m.serverAddr))
		cmds = append(cmds, textinput.Blink)
		cmds = append(cmds, tea.ClearScreen)
	case settingsmodel.SetServerAddrMsg:
		settingsmodel.SetServerAddr(m.serverAddr)
	case loginmodel.LoginMsg:
		if err := m.login(msg.Login, msg.Password); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			m.currentUser.Login = msg.Login
			m.currentUser.Password = msg.Password
			cmds = append(cmds, common.ChangeState(common.WorkState), tea.ClearScreen)
			cards, err := m.storage.GetCards()
			if err != nil {
				cmds = append(cmds, common.MakeError(err))
			} else {
				cmds = append(cmds, workmodel.UpdateCards(cards))
			}
		}
	case registermodel.RegisterMsg:
		if err := m.register(msg.Login, msg.Password); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			m.currentUser.Login = msg.Login
			m.currentUser.Password = msg.Password
			cmds = append(cmds, common.MakeReset)
			cmds = append(cmds, common.ChangeState(common.MainState))
		}
	case common.LoginStoreMsg:
		if serverAddr, userInfo, err := m.loginLocalStore(string(msg)); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			m.serverAddr = serverAddr
			m.currentUser = userInfo
			if len(m.serverAddr) != 0 {
				if err := m.initClient(m.serverAddr); err != nil {
					cmds = append(cmds, common.MakeError(err))
				}
			}
			cmds = append(cmds, common.ChangeState(common.MainState))
		}
	case settingsmodel.SaveServerAddrMsg:
		m.serverAddr = string(msg)
		if err := m.initClient(m.serverAddr); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else if err := m.storage.SetServerAddr(m.serverAddr); err != nil {
			cmds = append(cmds, common.MakeError(err), common.ChangeState(common.MainState))
		}
	case workmodel.NewCardCmd:
		if err := m.storage.UpsertCard(msg.CurrentCardID, msg.Name, msg.Ccn, msg.CVV, msg.Exp); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cmds = append(cmds, common.MakeMsg("data saved"))
			cards, err := m.storage.GetCards()
			if err != nil {
				cmds = append(cmds, common.MakeError(err))
			} else {
				cmds = append(cmds, workmodel.UpdateCards(cards))
			}
		}
	case workmodel.SyncMsg:
		if err := m.remoteClient.Sync(m.ctx); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cards, err := m.storage.GetCards()
			if err != nil {
				cmds = append(cmds, common.MakeError(err))
			} else {
				cmds = append(cmds, workmodel.UpdateCards(cards))
			}
		}
	case workmodel.DeleteCardMsg:
		if err := m.storage.DeleteCard(string(msg)); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cards, err := m.storage.GetCards()
			if err != nil {
				cmds = append(cmds, common.MakeError(err))
			} else {
				cmds = append(cmds, workmodel.UpdateCards(cards))
			}
		}
	case tea.QuitMsg:
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.mainState, cmd = m.mainState.Update(msg)
	cmds = append(cmds, cmd)
	m.loginState, cmd = m.loginState.Update(msg)
	cmds = append(cmds, cmd)
	m.registerState, cmd = m.registerState.Update(msg)
	cmds = append(cmds, cmd)
	m.workState, cmd = m.workState.Update(msg)
	cmds = append(cmds, cmd)
	m.settingsState, cmd = m.settingsState.Update(msg)
	cmds = append(cmds, cmd)
	m.storeState, cmd = m.storeState.Update(msg)
	cmds = append(cmds, cmd)
	m.loginStoreState, cmd = m.loginStoreState.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View returns a string based on data in the model. That string which will be
// rendered to the terminal.
func (m model) View() string {
	view := m.states[m.currentState].View()
	view += "\n\nPress Ctrl+H to show/hide help."
	if m.helpShown {
		view += m.help
		view += fmt.Sprintf("\n\nVersion: %s\nBuild date: %s", buildinfo.Version, buildinfo.BuildDate)
	}

	if len(m.currentUser.Login) != 0 {
		view += fmt.Sprintf("\n\nYour current user is \"%s\"", m.currentUser.Login)
	}

	if m.err != nil {
		view += fmt.Sprintf("\n\nError: %s", m.err)
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
