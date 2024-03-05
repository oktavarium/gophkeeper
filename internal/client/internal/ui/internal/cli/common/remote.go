package common

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type RemoteClient interface {
	Init(ctx context.Context, addr string) error
	Register(ctx context.Context, in dto.UserInfo) error
	Login(ctx context.Context, in dto.UserInfo) error
	Sync(ctx context.Context) error
}
