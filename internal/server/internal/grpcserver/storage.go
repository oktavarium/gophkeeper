package grpcserver

import (
	"context"
	"time"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type Storage interface {
	Register(context.Context, string, string) error
	Login(context.Context, string, string) error
	Sync(context.Context, string, map[string]dto.SimpleDataEncrypted) (map[string]dto.SimpleDataEncrypted, error)
	GetToken(context.Context, string) (string, time.Time, error)
	UpdateToken(context.Context, string, string, string, time.Time) error
	GetTokenIDByLogin(context.Context, string) (string, error)
	GetUserIDByToken(context.Context, string) (string, error)
	GetUserIDByLogin(context.Context, string) (string, error)
}
