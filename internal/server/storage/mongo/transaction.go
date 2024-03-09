package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

func (s *Storage) transaction(
	ctx context.Context,
	fn func(ctx mongo.SessionContext) (interface{}, error),
) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	session, err := s.client.StartSession()
	if err != nil {
		return fmt.Errorf("error on starting session: %w", err)
	}
	defer session.EndSession(context.TODO())

	_, err = session.WithTransaction(ctx, fn, txnOptions)
	if err != nil {
		return fmt.Errorf("error on making transaction: %w", err)
	}

	return nil
}
