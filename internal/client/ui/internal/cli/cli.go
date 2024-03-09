package cli

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/common"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/loginmodel"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/loginstoremodel"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/mainmodel"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/registermodel"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/settingsmodel"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/storemodel"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/workmodel"
	"github.com/oktavarium/gophkeeper/internal/shared/buildinfo"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

// model saves states and commands for them
type model struct {
	ctx          context.Context
	states       map[common.State]common.Model
	currentState common.State
	help         string
	helpShown    bool
	storage      Storage
	remoteClient common.RemoteClient
	currentUser  models.UserInfo
	serverAddr   string
	err          error
}

func newModel(ctx context.Context, s Storage, c common.RemoteClient) model {
	states := map[common.State]common.Model{
		common.MainState:       mainmodel.NewModel(),
		common.LoginState:      loginmodel.NewModel(),
		common.RegisterState:   registermodel.NewModel(),
		common.WorkState:       workmodel.NewModel(),
		common.SettingsState:   settingsmodel.NewModel(),
		common.StoreState:      storemodel.NewModel(),
		common.LoginStoreState: loginstoremodel.NewModel(),
	}

	return model{
		ctx:          ctx,
		states:       states,
		help:         "\n\nNavigation: Arrows;\nBack: Esc;\nSelect command: Enter;\nExit: Ctrl+Q",
		storage:      s,
		remoteClient: c,
	}
}

func Run(ctx context.Context, s Storage, c common.RemoteClient) error {
	m := newModel(ctx, s, c)
	if _, err := tea.NewProgram(m, tea.WithContext(ctx)).Run(); err != nil {
		return fmt.Errorf("could not start program: %w", err)
	}

	return nil
}

// Init optionally returns an initial command we should run.
func (m model) Init() tea.Cmd {
	if err := m.storage.Check(); err != nil {
		return common.ChangeState(common.StoreState)
	}

	return common.ChangeState(common.LoginStoreState)
}

// Update is called when messages are received.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case common.ErrorMsg:
		m.err = msg
	case tea.KeyMsg:
		//nolint:exhaustive // too many unused cased
		switch msg.Type {
		case tea.KeyCtrlQ:
			cmds = append(cmds, tea.Quit)
		case tea.KeyCtrlH:
			m.helpShown = !m.helpShown
		}
	case common.StateMsg:
		m.err = nil
		m.currentState = common.State(msg)
		m.Focus()
		cmds = append(cmds, settingsmodel.SetServerAddr(m.serverAddr), textinput.Blink, tea.ClearScreen)
	case settingsmodel.SetServerAddrMsg:
		settingsmodel.SetServerAddr(m.serverAddr)
	case loginmodel.LoginMsg:
		if err := m.login(msg.Login, msg.Password); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			m.currentUser.Login, m.currentUser.Password = msg.Login, msg.Password
			cmds = append(cmds, common.ChangeState(common.WorkState), tea.ClearScreen, m.updateCards())
		}
	case registermodel.RegisterMsg:
		if err := m.register(msg.Login, msg.Password); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			m.currentUser.Login, m.currentUser.Password = msg.Login, msg.Password
			cmds = append(cmds, common.MakeReset, common.ChangeState(common.MainState))
		}
	case common.LoginStoreMsg:
		if serverAddr, userInfo, err := m.loginLocalStore(string(msg)); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			m.serverAddr, m.currentUser = serverAddr, userInfo
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
			cmds = append(cmds, common.MakeMsg("data saved"), m.updateCards())
		}
	case workmodel.NewSimpleDataCmd:
		if err := m.storage.UpsertSimple(msg.CurrentID, msg.Name, msg.Data); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cmds = append(cmds, common.MakeMsg("data saved"), m.updateCards())
		}
	case workmodel.NewFileCmd:
		if err := m.storage.UpsertBinary(msg.CurrentID, msg.Name, msg.Path); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cmds = append(cmds, common.MakeMsg("data saved"), m.updateCards())
		}
	case workmodel.SaveFileCmd:
		if err := m.storage.SaveFile(string(msg)); err != nil {
			cmds = append(cmds, common.MakeError(err))
		}
	case workmodel.SyncMsg:
		if err := m.remoteClient.Sync(m.ctx); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cmds = append(cmds, m.updateCards())
		}
	case workmodel.DeleteDataMsg:
		if err := m.storage.DeleteData(string(msg)); err != nil {
			cmds = append(cmds, common.MakeError(err))
		} else {
			cmds = append(cmds, m.updateCards())
		}
	case tea.QuitMsg:
		return m, tea.Quit
	}

	cmds = append(cmds, m.states[m.currentState].Update(msg))
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

func (m *model) Blur() {
	m.states[m.currentState].Blur()
}

func (m *model) Focus() {
	for _, v := range m.states {
		v.Blur()
	}
	m.states[m.currentState].Focus()
}
