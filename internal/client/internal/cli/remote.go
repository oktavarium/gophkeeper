package cli

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type remoteClient interface {
	Init(ctx context.Context, addr string) error
	Register(ctx context.Context, in dto.UserInfo) error
	Login(ctx context.Context, in dto.UserInfo) error
	Save(ctx context.Context, in dto.SaveData) error
}
