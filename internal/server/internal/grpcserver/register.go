package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

func (s *GrpcServer) Register(ctx context.Context, req *pbapi.RegisterRequest) (*pbapi.RegisterResponse, error) {
	resp := &pbapi.RegisterResponse{}
	if err := s.storage.Register(ctx, grpcUserInfoToDtoUserInfo(req.GetUserInfo())); err != nil {
		return resp, status.Error(codes.Internal, "register failed")
	}

	return resp, nil
}
