package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

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

	_, err := s.client.Register(ctx, request)
	if err != nil {
		return fmt.Errorf("error on calling register: %w", err)
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

func (s *GrpcClient) Sync(ctx context.Context) error {
	if err := s.isInited(); err != nil {
		return fmt.Errorf("error on save: %w", err)
	}

	cards, err := s.storage.GetCardsEncrypted()
	if err != nil {
		return fmt.Errorf("error on getting data: %w", err)
	}

	req := &pbapi.SyncRequest{}
	for k, v := range cards {
		req.SyncData = append(req.SyncData, &pbapi.SyncData{
			Uid:      k,
			Modified: timestamppb.New(v.Common.Modified),
			Deleted:  v.Common.IsDeleted,
			Data:     v.Data,
		})
	}

	resp, err := s.client.Sync(ctx, req)
	if err != nil {
		return fmt.Errorf("error on data sync: %w", err)
	}

	if resp.GetSyncData() != nil {
		cards := make(map[string]dto.SimpleCardDataEncrypted, len(resp.GetSyncData()))
		for _, v := range resp.GetSyncData() {
			cards[v.Uid] = dto.SimpleCardDataEncrypted{
				Common: dto.CommonData{
					Modified:  v.Modified.AsTime(),
					IsDeleted: v.Deleted,
				},
				Data: v.GetData(),
			}
		}

		if err := s.storage.UpdateCardsEncrypted(cards); err != nil {
			return fmt.Errorf("error on updating data after sync: %w", err)
		}
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
