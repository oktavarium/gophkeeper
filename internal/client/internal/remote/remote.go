package remote

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/client/internal/remote/internal/grpcclient"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

type Client interface {
	Init(ctx context.Context, addr string) error
	Register(ctx context.Context, in dto.UserInfo) error
	Login(ctx context.Context, in dto.UserInfo) error
	Save(ctx context.Context, in dto.SaveData) error
}

func NewGrpcClient() (Client, error) {
	return grpcclient.NewGrpcClient()
}
