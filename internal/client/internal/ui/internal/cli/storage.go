package cli

import (
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

type Storage interface {
	Check() error
	Open(string) error
	GetServerAddr() (string, error)
	GetLoginAndPass() (string, string, error)
	SetServerAddr(string) error
	SetLoginAndPass(string, string) error
	UpsertCard(string, string, string, uint32, time.Time) error
	UpsertSimple(string, string, string) error
	GetData() (map[string]models.SimpleCardData, map[string]models.SimpleData, map[string]models.SimpleBinaryData, error)
	DeleteCard(string) error
	DeleteSimple(string) error
	DeleteBinary(string) error
}
