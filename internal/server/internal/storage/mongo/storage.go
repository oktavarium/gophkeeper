package mongo

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"github.com/oktavarium/gophkeeper/internal/shared/dto"
)

const uri = "mongodb://root:example@localhost:27017/"

type Storage struct {
	client *mongo.Client
}

func NewStorage(ctx context.Context) (*Storage, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("error on mongo connect: %w", err)
	}

	// defer func() {
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	return &Storage{
		client: client,
	}, nil
}

func (s *Storage) Register(ctx context.Context, in dto.UserInfo) error {
	coll := s.client.Database("keeper").Collection("users")
	var result UserInfo
	if err := coll.FindOne(ctx, bson.D{{"login", in.Login}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), 8)
			_, err = coll.InsertOne(ctx, &UserInfo{Login: in.Login, Password: hashedPassword})
			if err != nil {
				return fmt.Errorf("error on saving user: %w", err)
			}
			return nil
		}

		return fmt.Errorf("error on seaching user: %w", err)
	}

	return fmt.Errorf("user already exists")
}

func (s *Storage) Login(ctx context.Context, in dto.UserInfo) error {
	coll := s.client.Database("keeper").Collection("users")
	var result UserInfo
	if err := coll.FindOne(ctx, bson.D{{"login", in.Login}}).Decode(&result); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), 8)
			_, err = coll.InsertOne(ctx, &UserInfo{Login: in.Login, Password: hashedPassword})
			if err != nil {
				return fmt.Errorf("error on saving user: %w", err)
			}
			return nil
		}

		return fmt.Errorf("error on seaching user: %w", err)
	}

	return fmt.Errorf("user already exists")
}

func (s *Storage) Sync(ctx context.Context) error {
	return nil
}
