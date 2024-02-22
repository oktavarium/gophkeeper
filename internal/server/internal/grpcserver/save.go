package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbapi "github.com/oktavarium/gophkeeper/api"
)

func (s *GrpcServer) Save(ctx context.Context, req *pbapi.SaveRequest) (*pbapi.SaveResponse, error) {
	resp := &pbapi.SaveResponse{}
	if err := s.storage.Save(ctx, grpcSaveDataToDtoSavaData(req)); err != nil {
		return resp, status.Error(codes.Internal, "save failed")
	}

	return resp, nil
}
