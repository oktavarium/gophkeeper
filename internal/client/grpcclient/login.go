package grpcclient

import (
	"context"
	"fmt"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (c *GrpcClient) Login(ctx context.Context, in models.UserInfo) error {
	if err := c.isInited(); err != nil {
		return fmt.Errorf("error on login: %w", err)
	}

	request := &pbapi.LoginRequest{
		UserInfo: dtoUserInfoToGrpcUserInfo(in),
	}

	_, err := c.client.Login(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling login: %w", err)
	}

	return nil
}
