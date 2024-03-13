package grpcclient

import (
	"context"
	"fmt"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (c *GrpcClient) Register(ctx context.Context, in models.UserInfo) error {
	if err := c.isInited(); err != nil {
		return fmt.Errorf("error on register: %w", err)
	}

	request := &pbapi.RegisterRequest{
		UserInfo: dtoUserInfoToGrpcUserInfo(in),
	}

	_, err := c.client.Register(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling register: %w", err)
	}

	return nil
}
