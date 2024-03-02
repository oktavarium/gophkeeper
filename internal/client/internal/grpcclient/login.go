package grpcclient

import (
	"context"
	"fmt"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *GrpcClient) Login(ctx context.Context, in dto.UserInfo) error {
	if err := s.isInited(); err != nil {
		return fmt.Errorf("error on login: %w", err)
	}

	request := &pbapi.LoginRequest{
		UserInfo: dtoUserInfoToGrpcUserInfo(in),
	}

	_, err := s.client.Login(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling login: %w", err)
	}

	return nil
}
