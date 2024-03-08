package storage

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/oktavarium/gophkeeper/internal/client/internal/crypto"
	"github.com/oktavarium/gophkeeper/internal/client/internal/storage/internal/jsonfile"
)

const storagePath = "./file.storage"

type Storage struct {
	store  *jsonfile.JSONFile[storageModel]
	crypto *crypto.Crypto
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Check() error {
	if _, err := os.Stat(storagePath); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file not exists")
	}

	return nil
}

func (s *Storage) Open(pass string) error {
	var store *jsonfile.JSONFile[storageModel]
	var err error

	if err = s.initCrypto(pass); err != nil {
		return fmt.Errorf("error on initting crypto module: %w", err)
	}

	if s.Check() != nil {
		store, err = jsonfile.New[storageModel](storagePath, s.crypto)
		if err != nil {
			return fmt.Errorf("error on checking storage: %w", err)
		}

		if err = store.Write(func(data *storageModel) error {
			data.MasterPass = crypto.HashPassword(pass)
			data.SimpleCardData = make(map[string]simpleCardData)
			data.SimpleData = make(map[string]simpleData)
			data.SimpleBinaryData = make(map[string]simpleBinaryData)
			return nil
		}); err != nil {
			return fmt.Errorf("error on saving master password; %w", err)
		}
	} else {
		store, err = jsonfile.Load[storageModel](storagePath, s.crypto)
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

	s.store = store

	return nil
}

func (s *Storage) initCrypto(pass string) error {
	c, err := crypto.NewCrypto(pass)
	if err != nil {
		return fmt.Errorf("error on creating crypto provider: %w", err)
	}

	s.crypto = c
	return nil
}

func (s *Storage) isInited() bool {
	return s.store != nil
}
