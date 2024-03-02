package grpcclient

import (
	"context"
	"fmt"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *GrpcClient) Register(ctx context.Context, in dto.UserInfo) error {
	if err := s.isInited(); err != nil {
		return fmt.Errorf("error on register: %w", err)
	}

	request := &pbapi.RegisterRequest{
		UserInfo: dtoUserInfoToGrpcUserInfo(in),
	}

	_, err := s.client.Register(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling register: %w", err)
	}

	return nil
}
