package remote

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

type Storage struct {
	inited bool
	client pbapi.GophKeeperClient
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Register(ctx context.Context, in dto.UserInfo) error {
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

func (s *Storage) Login(ctx context.Context, in dto.UserInfo) error {
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

func (s *Storage) Save(ctx context.Context, in dto.SaveData) error {
	if err := s.isInited(); err != nil {
		return fmt.Errorf("error on save: %w", err)
	}

	request := dtoSavaDataToGrpcSaveData(in)

	_, err := s.client.Save(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling save: %w", err)
	}

	return nil
}

func (s *Storage) Init(ctx context.Context, addr string) error {
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("error on dialing: %s: %w", addr, err)
	}

	s.client = pbapi.NewGophKeeperClient(conn)
	s.inited = true

	return nil
}

func (s *Storage) isInited() error {
	if !s.inited {
		return fmt.Errorf("client not inited")
	}

	return nil
}
