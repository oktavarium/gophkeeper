package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func (s *Storage) Register(ctx context.Context, login string, password string) error {
	coll := s.client.Database("keeper").Collection("users")
	var result UserInfo
	if err := coll.FindOne(ctx, bson.D{{"login", login}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
			_, err = coll.InsertOne(
				ctx,
				&UserInfo{
					ID:       primitive.NewObjectID(),
					Login:    login,
					Password: hashedPassword,
				},
			)
			if err != nil {
				return fmt.Errorf("error on saving user: %w", err)
			}
			return nil
		}

		return fmt.Errorf("error on seaching user: %w", err)
	}

	return fmt.Errorf("user already exists")
}
