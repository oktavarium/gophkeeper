package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *GrpcServer) Sync(ctx context.Context, req *pbapi.SyncRequest) (*pbapi.SyncResponse, error) {
	userID, err := s.storage.GetUserIDByToken(ctx, req.GetTokenID())
	if err != nil {
		return &pbapi.SyncResponse{}, status.Errorf(codes.Internal, "error on getting token user")
	}

	datagrams := make(map[string]dto.SimpleDataEncrypted, len(req.GetSyncData()))
	for _, v := range req.GetSyncData() {
		datagrams[v.GetUid()] = dto.SimpleDataEncrypted{
			Common: dto.CommonData{
				Deleted:  v.GetDeleted(),
				Modified: v.GetModified().AsTime(),
			},
			Data: v.GetData(),
		}
	}

	updatedDatagrams, err := s.storage.Sync(ctx, userID, datagrams)
	if err != nil {
		return &pbapi.SyncResponse{}, status.Errorf(codes.Internal, "error on syncing data")
	}

	syncData := make([]*pbapi.SyncData, 0, len(updatedDatagrams))
	for k, v := range updatedDatagrams {
		syncData = append(syncData, &pbapi.SyncData{
			Uid:      k,
			Modified: timestamppb.New(v.Common.Modified),
			Deleted:  v.Common.Deleted,
			Type:     pbapi.DataTypes(v.Common.Type),
			Data:     v.Data,
		})
	}
	return &pbapi.SyncResponse{
		SyncData: syncData,
	}, nil
}
