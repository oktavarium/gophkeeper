package grpcclient

import (
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type storage interface {
	GetToken() (string, time.Time, error)
	UpdateToken(string, time.Time) error
	GetCardsEncrypted() (map[string]dto.SimpleDataEncrypted, error)
	UpdateCardsEncrypted(map[string]dto.SimpleDataEncrypted) error
}
