package grpcserver

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type Storage interface {
	Register(ctx context.Context, in dto.UserInfo) error
	Login(ctx context.Context, in dto.UserInfo) error
	Sync(ctx context.Context) error
}
