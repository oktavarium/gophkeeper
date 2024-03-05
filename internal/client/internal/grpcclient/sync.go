package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/types/known/timestamppb"

	pbapi "github.com/oktavarium/gophkeeper/api"
	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (c *GrpcClient) Sync(ctx context.Context) error {
	if err := c.isInited(); err != nil {
		return fmt.Errorf("error on sync: %w", err)
	}

	cards, err := c.storage.GetDataEncrypted()
	if err != nil {
		return fmt.Errorf("error on getting data: %w", err)
	}

	req := &pbapi.SyncRequest{}
	for k, v := range cards {
		req.SyncData = append(req.GetSyncData(), &pbapi.SyncData{
			Uid:      k,
			Modified: timestamppb.New(v.Common.Modified),
			Deleted:  v.Common.Deleted,
			Type:     pbapi.DataTypes_Card,
			Data:     v.Data,
		})
	}

	resp, err := c.client.Sync(ctx, req)
	if err != nil {
		return fmt.Errorf("error on data sync: %w", err)
	}

	if resp.GetSyncData() != nil {
		cards := make(map[string]models.SimpleDataEncrypted, len(resp.GetSyncData()))
		for _, v := range resp.GetSyncData() {
			cards[v.GetUid()] = models.SimpleDataEncrypted{
				Common: models.CommonData{
					Modified: v.GetModified().AsTime(),
					Deleted:  v.GetDeleted(),
					Type:     models.DataType(v.GetType().Number()),
				},
				Data: v.GetData(),
			}
		}

		if err := c.storage.UpdateDataEncrypted(cards); err != nil {
			return fmt.Errorf("error on updating data after sync: %w", err)
		}
	}

	return nil
}
