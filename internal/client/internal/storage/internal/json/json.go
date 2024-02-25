package json

import (
	"bytes"
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/client/internal/crypto"
	"github.com/oktavarium/gophkeeper/internal/jsonfile"
)

const storagePath = "./file.storage"

type JsonStorage struct {
	store  *jsonfile.JSONFile[storageModel]
	crypto *crypto.Crypto
}

func NewStorage() *JsonStorage {
	return &JsonStorage{}
}

func (s *JsonStorage) Check() error {
	if _, err := jsonfile.Load[storageModel](storagePath); err != nil {
		return fmt.Errorf("error on checking storage: %w", err)
	}

	return nil
}

func (s *JsonStorage) Open(pass string) error {
	var store *jsonfile.JSONFile[storageModel]
	var err error
	if err = s.Check(); err != nil {
		store, err = jsonfile.New[storageModel](storagePath)
		if err != nil {
			return fmt.Errorf("error on checking storage: %w", err)
		}

		if err = store.Write(func(data *storageModel) error {
			data.MasterPass = crypto.HashPassword(pass)
			return nil
		}); err != nil {
			return fmt.Errorf("error on saving master password; %w", err)
		}

	} else {
		store, err = jsonfile.Load[storageModel](storagePath)
		if err != nil {
			return fmt.Errorf("error on loading storage: %w", err)
		}

		var masterPass [32]byte
		store.Read(func(data *storageModel) {
			masterPass = data.MasterPass
		})

		userPassHash := crypto.HashPassword(pass)
		if !bytes.Equal(masterPass[:], userPassHash[:]) {
			return fmt.Errorf("error on checking paster password")
		}
	}

	if err = s.initCrypto(pass); err != nil {
		return fmt.Errorf("error on initting crypto module: %w", err)
	}

	s.store = store

	return nil
}

// func (s *JsonStorage) Register(in dto.UserInfo) error {
// 	return nil
// }
// func (s *JsonStorage) Login(in dto.UserInfo) error {
// 	return nil
// }
// func (s *JsonStorage) Save(in dto.SaveData) error {
// 	return nil
// }

func (s *JsonStorage) GetServerAddr() (string, error) {
	if !s.isInited() {
		return "", fmt.Errorf("storage is not inited")
	}

	var serverAddrEnrypted []byte
	s.store.Read(func(data *storageModel) {
		serverAddrEnrypted = data.ServerAddr
	})

	serverAddr, err := s.crypto.DecryptData(serverAddrEnrypted)
	if err != nil {
		return serverAddr, fmt.Errorf("error on decrypting server addr: %w", err)
	}

	return serverAddr, nil
}

func (s *JsonStorage) GetLoginAndPass() (string, string, error) {
	if !s.isInited() {
		return "", "", fmt.Errorf("storage is not inited")
	}

	var loginEnrypted []byte
	var passEnrypted []byte
	s.store.Read(func(data *storageModel) {
		loginEnrypted = data.Login
		passEnrypted = data.Password
	})

	login, err := s.crypto.DecryptData(loginEnrypted)
	if err != nil {
		return "", "", fmt.Errorf("error on decrypting login: %w", err)
	}

	pass, err := s.crypto.DecryptData(passEnrypted[:])
	if err != nil {
		return "", "", fmt.Errorf("error on decrypting pass: %w", err)
	}

	return login, pass, nil
}

func (s *JsonStorage) SetServerAddr(serveAddr string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	serverAddrEncrypted, err := s.crypto.EncryptData(serveAddr)
	if err != nil {
		return fmt.Errorf("error on ecnrypting server addr: %w", err)
	}

	if err = s.store.Write(func(data *storageModel) error {
		data.ServerAddr = serverAddrEncrypted
		return nil
	}); err != nil {
		return fmt.Errorf("error on saving server addr; %w", err)
	}

	return nil
}

func (s *JsonStorage) SetLoginAndPass(login, pass string) error {
	if !s.isInited() {
		return fmt.Errorf("storage is not inited")
	}

	loginEncrypted, err := s.crypto.EncryptData(login)
	if err != nil {
		return fmt.Errorf("error on ecnrypting login: %w", err)
	}

	passEncrypted, err := s.crypto.EncryptData(pass)
	if err != nil {
		return fmt.Errorf("error on ecnrypting pass: %w", err)
	}

	if err = s.store.Write(func(data *storageModel) error {
		data.Login = loginEncrypted
		data.Password = passEncrypted
		return nil
	}); err != nil {
		return fmt.Errorf("error on saving login and pass; %w", err)
	}

	return nil
}

func (s *JsonStorage) initCrypto(pass string) error {
	c, err := crypto.NewCrypto(pass)
	if err != nil {
		return fmt.Errorf("error on creating crypto provider: %w", err)
	}

	s.crypto = c
	return nil
}

func (s *JsonStorage) isInited() bool {
	return s.store != nil
}
