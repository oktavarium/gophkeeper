package common

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

type RemoteClient interface {
	Init(ctx context.Context, addr string) error
	Register(ctx context.Context, in models.UserInfo) error
	Login(ctx context.Context, in models.UserInfo) error
	Sync(ctx context.Context) error
}
