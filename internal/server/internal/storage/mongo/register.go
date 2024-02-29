package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

func (s *Storage) Register(ctx context.Context, in dto.UserInfo) error {
	coll := s.client.Database("keeper").Collection("users")
	var result UserInfo
	if err := coll.FindOne(ctx, bson.D{{"login", in.Login}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), 8)
			_, err = coll.InsertOne(
				ctx,
				&UserInfo{
					ID:       primitive.NewObjectID(),
					Login:    in.Login,
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
