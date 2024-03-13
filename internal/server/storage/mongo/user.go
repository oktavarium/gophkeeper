package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (s *Storage) GetUserIDByLogin(ctx context.Context, login string) (string, error) {
	coll := s.client.Database("keeper").Collection("users")
	var userInfo UserInfo
	if err := coll.FindOne(ctx, bson.D{{"login", login}}).Decode(&userInfo); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", fmt.Errorf("no such user: %w", err)
		}

		return "", fmt.Errorf("internal error: %w", err)
	}

	return userInfo.ID.Hex(), nil
}
