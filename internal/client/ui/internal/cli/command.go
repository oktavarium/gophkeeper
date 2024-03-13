package cli

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/common"
	"github.com/oktavarium/gophkeeper/internal/client/ui/internal/cli/models/workmodel"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (m model) loginLocalStore(password string) (string, models.UserInfo, error) {
	var serverAddr string
	var userInfo models.UserInfo
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
	if err := m.remoteClient.Register(m.ctx, models.UserInfo{
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
	if err := m.remoteClient.Login(m.ctx, models.UserInfo{
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

func (m model) updateCards() tea.Cmd {
	cards, simple, binary, err := m.storage.GetData()
	if err != nil {
		return common.MakeError(err)
	}
	return workmodel.UpdateData(cards, simple, binary)
}
