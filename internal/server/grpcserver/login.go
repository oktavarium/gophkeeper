package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

func (s *GrpcServer) Login(ctx context.Context, req *pbapi.LoginRequest) (*pbapi.LoginResponse, error) {
	resp := &pbapi.LoginResponse{}
	if err := s.storage.Login(ctx, req.GetUserInfo().GetLogin(), req.GetUserInfo().GetPassword()); err != nil {
		return resp, status.Error(codes.Internal, "login failed")
	}

	return resp, nil
}
