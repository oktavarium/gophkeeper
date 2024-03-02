package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *GrpcClient) Sync(ctx context.Context) error {
	if err := s.isInited(); err != nil {
		return fmt.Errorf("error on sync: %w", err)
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
			Deleted:  v.Common.Deleted,
			Type:     pbapi.DataTypes_Card,
			Data:     v.Data,
		})
	}

	resp, err := s.client.Sync(ctx, req)
	if err != nil {
		return fmt.Errorf("error on data sync: %w", err)
	}

	if resp.GetSyncData() != nil {
		cards := make(map[string]dto.SimpleDataEncrypted, len(resp.GetSyncData()))
		for _, v := range resp.GetSyncData() {
			cards[v.Uid] = dto.SimpleDataEncrypted{
				Common: dto.CommonData{
					Modified: v.Modified.AsTime(),
					Deleted:  v.Deleted,
					Type:     dto.DataType(v.Type.Number()),
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
