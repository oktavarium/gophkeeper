package server

import (
	"context"

	"github.com/oktavarium/gophkeeper/internal/server/internal/grpcserver"
	storage2 "github.com/oktavarium/gophkeeper/internal/server/internal/storage"
)

func Run() error {
	ctx := context.Background()
	config := loadFlags()

	storage := storage2.NewMemoryStorage()
	server := grpcserver.NewGrpcServer(ctx, storage, config.serverAddr)

	return server.ListenAndServe()
}
