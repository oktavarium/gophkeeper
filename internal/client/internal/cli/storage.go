package cli

import (
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type storage interface {
	Check() error
	Open(string) error
	GetServerAddr() (string, error)
	GetLoginAndPass() (string, string, error)
	SetServerAddr(string) error
	SetLoginAndPass(string, string) error
	SaveNewCard(string, string, uint32, time.Time) error
	GetCards() (map[string]dto.SimpleCardData, error)
}
