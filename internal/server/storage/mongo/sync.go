package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/oktavarium/gophkeeper/internal/shared/models"
)

func (s *Storage) Sync(ctx context.Context, userID string, records map[string]models.SimpleDataEncrypted) (map[string]models.SimpleDataEncrypted, error) {
	idsToUpdate := make([]string, 0, len(records))
	idsToUpdateMap := make(map[string]struct{}, len(records))
	for k := range records {
		idsToUpdate = append(idsToUpdate, k)
		idsToUpdateMap[k] = struct{}{}
	}

	updatedData := make(map[string]models.SimpleDataEncrypted)

	// получаем данные, которых нет у клиента
	coll := s.client.Database("keeper").Collection(userID)
	filter := bson.D{{"data_id", bson.D{{"$nin", idsToUpdate}}}}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error on find: %w", err)
	}
	for cursor.Next(context.TODO()) {
		var result CommonData
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error on decoding: %w", err)
		}
		updatedData[result.DataID] = models.SimpleDataEncrypted{
			Common: models.CommonData{
				Deleted:  result.Deleted,
				Modified: result.Modified,
				Type:     models.DataType(result.DataType),
			},
			Data: result.Data,
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error on fetching data: %w", err)
	}

	// проверяем данные, которых есть у клиента
	filter = bson.D{{"data_id", bson.D{{"$in", idsToUpdate}}}}
	cursor, err = coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error on find: %w", err)
	}
	for cursor.Next(context.TODO()) {
		var result CommonData
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error on decoding: %w", err)
		}

		if records[result.DataID].Common.Modified.After(result.Modified) {
			// если на клиенте более актуальные данные, обновим их у себя
			if err := s.UpsertData(
				ctx,
				userID,
				CommonData{
					DataID:   result.DataID,
					DataType: result.DataType,
					Modified: records[result.DataID].Common.Modified,
					Deleted:  records[result.DataID].Common.Deleted,
					Data:     records[result.DataID].Data,
				}); err != nil {
				return nil, fmt.Errorf("error on updating data: %w", err)
			}
		} else if records[result.DataID].Common.Modified.Before(result.Modified) {
			// если на сервере более актуальные данные, отправим их клиенту
			updatedData[result.DataID] = models.SimpleDataEncrypted{
				Common: models.CommonData{
					Type:     models.DataType(result.DataType),
					Deleted:  result.Deleted,
					Modified: result.Modified,
				},
				Data: result.Data,
			}
		}
		delete(idsToUpdateMap, result.DataID)
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error on fetching data: %w", err)
	}

	for k := range idsToUpdateMap {
		if err := s.UpsertData(
			ctx,
			userID,
			CommonData{
				DataID:   k,
				DataType: int(records[k].Common.Type),
				Modified: records[k].Common.Modified,
				Deleted:  records[k].Common.Deleted,
				Data:     records[k].Data,
			}); err != nil {
			return nil, fmt.Errorf("error on updating data: %w", err)
		}
	}

	return updatedData, nil
}

func (s *Storage) UpsertData(ctx context.Context, userID string, data CommonData) error {
	coll := s.client.Database("keeper").Collection(userID)
	filter := bson.D{{"data_id", data.DataID}}
	update := bson.D{{"$set",
		bson.D{
			{"is_deleted", data.Deleted},
			{"modified", data.Modified},
			{"data_type", data.DataType},
			{"data", data.Data},
		}},
	}
	opts := options.Update().SetUpsert(true)
	if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("error on updating data: %w", err)
	}

	return nil
}
