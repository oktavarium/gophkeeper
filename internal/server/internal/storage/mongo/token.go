package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *Storage) GetToken(ctx context.Context, tokenID string) (string, time.Time, error) {
	coll := s.client.Database("keeper").Collection("tokens")
	var result Token
	if err := coll.FindOne(ctx, bson.D{{"token_id", tokenID}}).Decode(&result); err != nil {
		return "", time.Now(), fmt.Errorf("error on seaching token: %w", err)
	}

	return result.TokenID, result.ValidUntil, nil
}

func (s *Storage) GetTokenIDByLogin(ctx context.Context, login string) (string, error) {
	userID, err := s.GetUserIDByLogin(ctx, login)
	if err != nil {
		return "", fmt.Errorf("error on getting user: %w", err)
	}
	bsonUserId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return "", fmt.Errorf("wrong user id: %w", err)
	}
	coll := s.client.Database("keeper").Collection("tokens")
	var token Token
	if err := coll.FindOne(ctx, bson.D{{"user_id", bsonUserId}}).Decode(&token); err != nil {
		return "", fmt.Errorf("error on seaching token: %w", err)
	}

	return token.TokenID, nil
}

func (s *Storage) UpdateToken(ctx context.Context, userID string, oldId string, newId string, validUntil time.Time) error {
	bsonUserId, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return fmt.Errorf("wrong user id: %w", err)
	}

	coll := s.client.Database("keeper").Collection("tokens")
	filter := bson.D{{"token_id", oldId}}
	update := bson.D{{"$set",
		bson.D{
			{"token_id", newId},
			{"user_id", bsonUserId},
			{"valid_until", validUntil},
		}},
	}
	opts := options.Update().SetUpsert(true)
	if _, err := coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("error on updating token: %w", err)
	}

	return nil
}

func (s *Storage) GetUserIDByToken(ctx context.Context, tokenID string) (string, error) {
	coll := s.client.Database("keeper").Collection("tokens")
	var token Token
	if err := coll.FindOne(ctx, bson.D{{"token_id", tokenID}}).Decode(&token); err != nil {
		return "", fmt.Errorf("error on seaching user: %w", err)
	}

	return token.UserID.Hex(), nil
}
