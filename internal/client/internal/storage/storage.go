package storage

import (
	"github.com/oktavarium/gophkeeper/internal/client/internal/storage/internal/json"
)

type Storage interface {
	Check() error
	Open(string) error
	// Register(dto.UserInfo) error
	// Login(dto.UserInfo) error
	// Save(dto.SaveData) error
	GetServerAddr() (string, error)
	GetLoginAndPass() (string, string, error)
	SetServerAddr(string) error
	SetLoginAndPass(string, string) error
}

func NewStorage() Storage {
	return json.NewStorage()
}
