package storage

import "fmt"

func (s *Storage) GetServerAddr() (string, error) {
	if !s.isInited() {
		return "", fmt.Errorf("storage is not inited")
	}

	var serverAddr string
	s.store.Read(func(data *storageModel) {
		serverAddr = data.ServerAddr
	})

	return serverAddr, nil
}

func (s *Storage) GetLoginAndPass() (string, string, error) {
	if !s.isInited() {
		return "", "", fmt.Errorf("storage is not inited")
	}

	var login string
	var password string
	s.store.Read(func(data *storageModel) {
		login = data.Login
		password = data.Password
	})

	return login, password, nil
}

func (s *Storage) SetServerAddr(serverAddr string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		data.ServerAddr = serverAddr
		return nil
	}); err != nil {
		return fmt.Errorf("error on saving server addr; %w", err)
	}

	return nil
}

func (s *Storage) SetLoginAndPass(login, password string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	if err := s.store.Write(func(data *storageModel) error {
		data.Login = login
		data.Password = password
		return nil
	}); err != nil {
		return fmt.Errorf("error on saving login and pass; %w", err)
	}

	return nil
}
