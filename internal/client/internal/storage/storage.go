package storage

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/client/internal/storage/internal/remote"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type Storage interface {
	Init(ctx context.Context, addr string) error
	Register(ctx context.Context, in dto.UserInfo) error
	Login(ctx context.Context, in dto.UserInfo) error
	Save(ctx context.Context, in dto.SaveData) error
}

func NewRemoteStorage() Storage {
	return remote.NewStorage()
}
