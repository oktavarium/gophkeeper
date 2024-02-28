package server

import (
	"context"
	"fmt"

	"github.com/oktavarium/gophkeeper/internal/server/internal/grpcserver"
	"github.com/oktavarium/gophkeeper/internal/server/internal/storage/mongo"
)

func Run() error {
	ctx := context.Background()
	config := loadFlags()

	storage, err := mongo.NewStorage(ctx)
	if err != nil {
		return fmt.Errorf("error on creating storage: %w", err)
	}
	server := grpcserver.NewGrpcServer(ctx, storage, config.serverAddr)

	return server.ListenAndServe()
}
