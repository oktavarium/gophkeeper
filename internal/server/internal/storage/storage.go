package storage

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/server/internal/storage/internal/memory"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type Storage interface {
	Register(ctx context.Context, in dto.UserInfo) error
	Login(ctx context.Context, in dto.UserInfo) error
	Save(ctx context.Context, in dto.SaveData) error
}

func NewMemoryStorage() Storage {
	return memory.NewStorage()
}
