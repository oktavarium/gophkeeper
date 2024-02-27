package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/oktavarium/gophkeeper/internal/client/internal/crypto"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

type GrpcClient struct {
	conn    *grpc.ClientConn
	client  pbapi.GophKeeperClient
	storage storage
	crypto  *crypto.Crypto
}

func NewGrpcClient(s storage) (*GrpcClient, error) {
	return &GrpcClient{
		storage: s,
	}, nil
}

func (s *GrpcClient) Register(ctx context.Context, in dto.UserInfo) error {
	if err := s.isInited(); err != nil {
		return fmt.Errorf("error on register: %w", err)
	}

	request := &pbapi.RegisterRequest{
		UserInfo: dtoUserInfoToGrpcUserInfo(in),
	}

	resp, err := s.client.Register(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling register: %w", err)
	}

	if err := s.storage.UpdateToken(dto.Token{
		Id:         resp.GetToken().GetId(),
		ValidUntil: resp.GetToken().GetValidUntil().AsTime(),
	}); err != nil {
		return fmt.Errorf("error on updating token")
	}

	return nil
}

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

func (s *GrpcClient) Save(ctx context.Context, in dto.SaveData) error {
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

func (s *GrpcClient) Init(ctx context.Context, addr string) error {
	if err := s.isInited(); err == nil {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("error on closing current conn: %w", err)
		}

		s.conn = nil
	}

	conn, err := grpc.DialContext(ctx,
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(s.cryptoUnaryInterceptor))
	if err != nil {
		return fmt.Errorf("error on dialing: %s: %w", addr, err)
	}

	s.conn = conn
	s.client = pbapi.NewGophKeeperClient(conn)

	return nil
}

func (s *GrpcClient) isInited() error {
	if s.conn == nil {
		return fmt.Errorf("client not inited")
	}

	return nil
}
