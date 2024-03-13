package grpcclient

import (
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

type storage interface {
	GetToken() (string, time.Time, error)
	UpdateToken(string, time.Time) error
	GetDataEncrypted() (map[string]models.SimpleDataEncrypted, error)
	UpdateDataEncrypted(map[string]models.SimpleDataEncrypted) error
}
