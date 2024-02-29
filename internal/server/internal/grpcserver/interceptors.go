package grpcserver

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/server/internal/crypto"
)

func (s *GrpcServer) cryptoUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	var reqTokenId string
	var userID string
	var newLogin string
	var err error

	switch r := req.(type) {
	case *pbapi.SyncRequest:
		reqTokenId = r.GetTokenID()
		_, validUntil, err := s.storage.GetToken(ctx, reqTokenId)
		if err != nil {
			return &pbapi.SyncResponse{}, status.Errorf(codes.NotFound, "invalid token")
		}
		if validUntil.Before(time.Now()) {
			return &pbapi.SyncResponse{}, status.Errorf(codes.PermissionDenied, "token expired")
		}
		userID, err = s.storage.GetUserIDByToken(ctx, reqTokenId)
		if err != nil {
			return &pbapi.SyncResponse{}, status.Errorf(codes.Internal, "user not found")
		}

	case *pbapi.LoginRequest:
		reqTokenId, err = s.storage.GetTokenIDByLogin(ctx, r.GetUserInfo().GetLogin())
		if err != nil {
			return &pbapi.LoginResponse{}, status.Errorf(codes.PermissionDenied, "no such user")
		}
		userID, err = s.storage.GetUserIDByToken(ctx, reqTokenId)
		if err != nil {
			return &pbapi.LoginResponse{}, status.Errorf(codes.Internal, "user not found")
		}
	case *pbapi.RegisterRequest:
		newLogin = r.GetUserInfo().GetLogin()
	}
	resp, err := handler(ctx, req)
	if err != nil {
		return resp, err
	}

	newTokenId, newTokenValidUntil := crypto.GenerateToken()
	newToken := &pbapi.Token{
		Id:         newTokenId,
		ValidUntil: timestamppb.New(newTokenValidUntil),
	}

	switch r := resp.(type) {
	case *pbapi.LoginResponse:
		r.Token = newToken
		resp = r
	case *pbapi.RegisterResponse:
		r.Token = newToken
		resp = r
		userID, err = s.storage.GetUserIDByLogin(ctx, newLogin)
		if err != nil {
			return &pbapi.RegisterResponse{}, status.Errorf(codes.Internal, "user not found")
		}
	case *pbapi.SyncResponse:
		r.Token = newToken
		resp = r
	}

	if err := s.storage.UpdateToken(ctx, userID, reqTokenId, newTokenId, newTokenValidUntil); err != nil {
		return resp, status.Errorf(codes.Internal, "error on updating token")
	}

	return resp, err
}
