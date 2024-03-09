package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (s *Storage) Login(ctx context.Context, login string, password string) error {
	coll := s.client.Database("keeper").Collection("users")
	var result UserInfo
	if err := coll.FindOne(ctx, bson.D{{"login", login}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("no such user: %w", err)
		}

		return fmt.Errorf("internal error: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword(result.Password, []byte(password)); err != nil {
		return fmt.Errorf("wrong password: %w", err)
	}

	return nil
}
