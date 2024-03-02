package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
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
	var userID string
	var err error
	peerInfo, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("can't obtain ip address of request")
	}

	remoteIP, _, _ := net.SplitHostPort(peerInfo.Addr.String())

	switch r := req.(type) {
	case *pbapi.SyncRequest:
		validIP, validUntil, err := s.storage.GetToken(ctx, r.GetTokenID())
		if err != nil {
			return &pbapi.SyncResponse{}, status.Errorf(codes.NotFound, "invalid token")
		}

		if validIP != remoteIP {
			return &pbapi.SyncResponse{}, status.Errorf(codes.PermissionDenied, "token is used for another ip")
		}

		if validUntil.Before(time.Now().UTC()) {
			return &pbapi.SyncResponse{}, status.Errorf(codes.PermissionDenied, "token expired")
		}

		userID, err = s.storage.GetUserIDByToken(ctx, r.GetTokenID())
		if err != nil {
			return &pbapi.SyncResponse{}, status.Errorf(codes.Internal, "user not found")
		}

	case *pbapi.LoginRequest:
		userID, err = s.storage.GetUserIDByLogin(ctx, r.GetUserInfo().GetLogin())
		if err != nil {
			return &pbapi.LoginResponse{}, status.Errorf(codes.Internal, "user not found")
		}
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
		userID, err = s.storage.GetUserIDByLogin(ctx, req.(*pbapi.RegisterRequest).GetUserInfo().GetLogin())
		if err != nil {
			return &pbapi.RegisterResponse{}, status.Errorf(codes.Internal, "user not found")
		}
	case *pbapi.SyncResponse:
		r.Token = newToken
		resp = r
	}

	if err := s.storage.UpdateToken(ctx, userID, newTokenId, remoteIP, newTokenValidUntil); err != nil {
		return resp, status.Errorf(codes.Internal, "error on updating token")
	}

	return resp, err
}
