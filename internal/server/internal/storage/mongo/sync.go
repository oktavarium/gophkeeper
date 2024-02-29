package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *Storage) Sync(ctx context.Context, userID string, cards map[string]dto.SimpleDataEncrypted) (map[string]dto.SimpleDataEncrypted, error) {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, fmt.Errorf("wrong user id: %w", err)
	}

	idsToUpdate := make([]string, 0, len(cards))
	idsToUpdateMap := make(map[string]struct{}, len(cards))
	for k := range cards {
		idsToUpdate = append(idsToUpdate, k)
		idsToUpdateMap[k] = struct{}{}
	}

	updatedCards := make(map[string]dto.SimpleDataEncrypted)

	// получаем данные, которых нет у клиента
	coll := s.client.Database("keeper").Collection("cards")
	filter := bson.D{{"data_id", bson.D{{"$nin", idsToUpdate}}}, {"user_id", objectUserID}}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error on find: %w", err)
	}
	for cursor.Next(context.TODO()) {
		var result CardData
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error on decoding: %w", err)
		}
		updatedCards[result.DataID] = dto.SimpleDataEncrypted{
			Common: dto.CommonData{
				IsDeleted: result.IsDeleted,
				Modified:  result.Modified,
			},
			Data: result.Data,
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error on fetching data: %w", err)
	}

	// проверяем данные, которых есть у клиента
	filter = bson.D{{"data_id", bson.D{{"$in", idsToUpdate}}}, {"user_id", objectUserID}}
	cursor, err = coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("error on find: %w", err)
	}
	for cursor.Next(context.TODO()) {
		var result CardData
		if err := cursor.Decode(&result); err != nil {
			return nil, fmt.Errorf("error on decoding: %w", err)
		}

		if cards[result.DataID].Common.Modified.After(result.Modified) {
			// если на клиенте более актуальные данные, обновим их у себя
			if err := s.UpsertData(
				ctx,
				userID,
				CardData{
					DataID:    result.DataID,
					Modified:  cards[result.DataID].Common.Modified,
					IsDeleted: cards[result.DataID].Common.IsDeleted,
					Data:      cards[result.DataID].Data,
				}); err != nil {
				return nil, fmt.Errorf("error on updating data: %w", err)
			}
		} else if cards[result.DataID].Common.Modified.Before(result.Modified) {
			// если на сервере более актуальные данные, отправим их клиенту
			updatedCards[result.DataID] = dto.SimpleDataEncrypted{
				Common: dto.CommonData{
					IsDeleted: result.IsDeleted,
					Modified:  result.Modified,
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
			CardData{
				DataID:    k,
				Modified:  cards[k].Common.Modified,
				IsDeleted: cards[k].Common.IsDeleted,
				Data:      cards[k].Data,
			}); err != nil {
			return nil, fmt.Errorf("error on updating data: %w", err)
		}
	}

	return updatedCards, nil
}

func (s *Storage) UpsertData(ctx context.Context, userID string, data CardData) error {
	objectUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("wrong user id: %w", err)
	}
	coll := s.client.Database("keeper").Collection("cards")
	filter := bson.D{{"data_id", data.DataID}, {"user_id", objectUserID}}
	update := bson.D{{"$set",
		bson.D{
			{"user_id", objectUserID},
			{"is_deleted", data.IsDeleted},
			{"modified", data.Modified},
			{"data", data.Data},
		}},
	}
	opts := options.Update().SetUpsert(true)
	if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("error on updating data: %w", err)
	}

	return nil
}
